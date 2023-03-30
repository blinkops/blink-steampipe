package generators

import (
	"fmt"
	"github.com/blinkops/blink-steampipe/scripts/consts"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	awsConnectionIdentifier     = "AWS_CONNECTION"
	awsAccessKeyId              = "ACCESS_KEY_ID"
	awsSecretAccessKey          = "SECRET_ACCESS_KEY"
	awsSessionToken             = "aws_session_token"
	awsRoleArn                  = "ROLE_ARN"
	awsExternalID               = "EXTERNAL_ID"
	awsWebIdentityTokenFile     = "AWS_WEB_IDENTITY_TOKEN_FILE"
	awsDefaultSessionRegion     = "us-east-1"
	awsRegionEnvVariable        = "AWS_REGION"
	awsDefaultRegionEnvVariable = "AWS_DEFAULT_REGION"
	awsRegionsListParam         = "AWS_REGIONS_PARAM"

	steampipeAwsConfigurationFile = consts.SteampipeSpcConfigurationPath + "aws.spc"
)

const (
	awsTrustedIdentityCredsAccessKeyId     = "TRUSTED_IDENTITY_ACCESS_KEY_ID"
	awsTrustedIdentityCredsAccessSecretKey = "TRUSTED_IDENTITY_ACCESS_SECRET_KEY"
	awsTrustedIdentityCredsSessionToken    = "TRUSTED_IDENTITY_SESSION_TOKEN"
)

const (
	awsUserBased       = "user_based"
	awsRoleBased       = "role_based"
	assumeCrossAccount = "assume_cross_account"
	assumeIdentity     = "assume_identity"
)

type AWSCredentialGenerator struct{}

func (gen AWSCredentialGenerator) Generate() error {
	if err := gen.generate(); err != nil {
		log.Debugf("failed resolving aws credentials, will try without credentials: %v", err)
		if err = replaceSpcConfigs("", "", ""); err != nil {
			log.Errorf("failed repalce aws credentials %v", err)
		}
	}
	return nil
}

func (gen AWSCredentialGenerator) generate() error {
	if _, ok := os.LookupEnv(awsConnectionIdentifier); !ok {

		// we need to replace the configs in case of AWS Query without connection
		// such as with EC2 runner
		if err := replaceSpcConfigs("", "", ""); err != nil {
			log.Errorf("failed repalce aws credentials %v", err)
		}
		return nil
	}

	awsCredentials := gen.initCredentialsMap()

	var access, secret, sessionToken string

	base, subBase := gen.detect(awsCredentials)
	switch base {
	case awsUserBased: // Implemented automatically via aws cli
		access, secret, sessionToken = awsCredentials[awsAccessKeyId], awsCredentials[awsSecretAccessKey], ""
	case awsRoleBased:
		var err error
		access, secret, sessionToken, err = gen.assumeRole(subBase, awsCredentials)
		if err != nil {
			return fmt.Errorf("unable to assume role with error: %w", err)
		}

	default:
		return errors.New("invalid aws connection was provided")
	}

	if err := replaceSpcConfigs(access, secret, sessionToken); err != nil {
		return err
	}

	return nil
}

func (gen AWSCredentialGenerator) getSessionRegion() string {
	if region := os.Getenv(awsRegionEnvVariable); region != "" {
		return region
	}
	if region := os.Getenv(awsDefaultRegionEnvVariable); region != "" {
		return region
	}
	return awsDefaultSessionRegion
}

func (gen AWSCredentialGenerator) detect(awsCredentials map[string]string) (base, sub string) {
	if awsCredentials[awsAccessKeyId] != "" && awsCredentials[awsSecretAccessKey] != "" {
		base = awsUserBased
	}
	if awsCredentials[awsRoleArn] == "" {
		return base, sub
	}
	if base == awsUserBased {
		return awsRoleBased, assumeCrossAccount
	}

	return awsRoleBased, assumeIdentity
}

func (gen AWSCredentialGenerator) getSession(awsCredentials map[string]string) (*session.Session, error) {
	sessConfig := aws.Config{
		Region: aws.String(awsDefaultSessionRegion),
	}
	if awsCredentials != nil {
		sessConfig.Credentials = credentials.NewStaticCredentials(awsCredentials[awsAccessKeyId], awsCredentials[awsSecretAccessKey], awsCredentials[awsSessionToken])
	}

	return session.NewSession(&sessConfig)
}

func (gen AWSCredentialGenerator) assumeRole(subBase string, awsCredentials map[string]string) (access, secret, sessionToken string, err error) {
	sessionName := strconv.Itoa(rand.Int())

	switch subBase {
	case assumeIdentity:
		return gen.assumeRoleWithIdentity(sessionName, awsCredentials)
	case assumeCrossAccount:
		return gen.assumeRoleCrossAccounts(sessionName, awsCredentials)
	}
	return "", "", "", errors.New("invalid assume role type was provided")
}

// assumeRoleWithIdentity first tries to assume a trusted identity using the role and external id, and if it doesn't
// succeed, it falls back to assuming a web identity using only the role
func (gen AWSCredentialGenerator) assumeRoleWithIdentity(sessionName string, awsCredentials map[string]string) (string, string, string, error) {
	log.Debugf("assuming role with identity. Trying to assume role with trusted identity first and falling back to web identity")
	accessKey, secretAccessKey, sessionToken, err := gen.assumeRoleWithTrustedIdentity(sessionName, awsCredentials)
	if err != nil {
		log.Errorf("error assuming role with trusted identity: %v", err)
		return gen.assumeRoleWithWebIdentity(sessionName, awsCredentials)
	}

	return accessKey, secretAccessKey, sessionToken, nil
}

func (gen AWSCredentialGenerator) assumeRoleWithWebIdentity(sessionName string, awsCredentials map[string]string) (string, string, string, error) {
	log.Debug("assuming role with web identity")

	sess, err := gen.getSession(awsCredentials)
	if err != nil {
		return "", "", "", errors.Wrap(err, "initializing aws session")
	}

	svc := sts.New(sess)
	tokenFile, ok := os.LookupEnv(awsWebIdentityTokenFile)
	if !ok {
		log.Debug("token file for irsa not found. try assume role")
		result, err := svc.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         aws.String(awsCredentials[awsRoleArn]),
			RoleSessionName: &sessionName,
		})
		if err != nil {
			return "", "", "", err
		}
		return *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken, err
	}

	data, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", "", "", fmt.Errorf("unable to open web identity token file with error: %w", err)
	}

	input := &sts.AssumeRoleWithWebIdentityInput{
		DurationSeconds:  aws.Int64(3600),
		RoleArn:          aws.String(awsCredentials[awsRoleArn]),
		RoleSessionName:  aws.String(sessionName),
		WebIdentityToken: aws.String(string(data)),
	}

	result, err := svc.AssumeRoleWithWebIdentity(input)
	if err != nil {
		return "", "", "", err
	}
	return *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken, err
}

func (gen AWSCredentialGenerator) assumeRoleWithTrustedIdentity(sessionName string, awsCredentials map[string]string) (string, string, string, error) {
	log.Debug("assuming role with trusted entity")

	trustedIdentityCreds := gen.getTrustedIdentityCreds(awsCredentials)

	sess, err := gen.getSession(trustedIdentityCreds)
	if err != nil {
		return "", "", "", errors.Wrap(err, "initializing aws session")
	}

	svc := sts.New(sess)

	assumeInput := &sts.AssumeRoleInput{
		RoleArn:         aws.String(awsCredentials[awsRoleArn]),
		RoleSessionName: &sessionName,
		ExternalId:      aws.String(awsCredentials[awsExternalID]),
	}

	if result, err := svc.AssumeRole(assumeInput); err == nil {
		return *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken, err
	}

	log.Errorf("error assuming role with trusted identity: %v", err)
	log.Debug("try assuming role without Blink trusted entity")

	sess, err = gen.getSession(nil)
	if err != nil {
		return "", "", "", errors.Wrap(err, "initializing aws session")
	}

	svc = sts.New(sess)
	result, err := svc.AssumeRole(assumeInput)
	if err != nil {
		return "", "", "", err
	}

	return *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken, nil
}

func (gen AWSCredentialGenerator) assumeRoleCrossAccounts(sessionName string, awsCredentials map[string]string) (string, string, string, error) {
	sess, err := gen.getSession(awsCredentials)
	if err != nil {
		return "", "", "", errors.Wrap(err, "initializing aws session")
	}

	svc := sts.New(sess)

	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(awsCredentials[awsRoleArn]),
		RoleSessionName: &sessionName,
	})
	if err != nil {
		return "", "", "", err
	}
	return *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken, err
}

func (gen AWSCredentialGenerator) initCredentialsMap() map[string]string {
	return map[string]string{
		awsAccessKeyId:                         os.Getenv(awsAccessKeyId),
		awsSecretAccessKey:                     os.Getenv(awsSecretAccessKey),
		awsRoleArn:                             os.Getenv(awsRoleArn),
		awsExternalID:                          os.Getenv(awsExternalID),
		awsTrustedIdentityCredsAccessKeyId:     os.Getenv(awsTrustedIdentityCredsAccessKeyId),
		awsTrustedIdentityCredsAccessSecretKey: os.Getenv(awsTrustedIdentityCredsAccessSecretKey),
		awsTrustedIdentityCredsSessionToken:    os.Getenv(awsTrustedIdentityCredsSessionToken),
	}
}

func (gen AWSCredentialGenerator) getTrustedIdentityCreds(credentials map[string]string) map[string]string {
	return map[string]string{
		awsAccessKeyId:     credentials[awsTrustedIdentityCredsAccessKeyId],
		awsSecretAccessKey: credentials[awsTrustedIdentityCredsAccessSecretKey],
		awsSessionToken:    credentials[awsTrustedIdentityCredsSessionToken],
	}
}

func replaceSpcConfigs(access, secret, sessionToken string) error {
	data, err := os.ReadFile(steampipeAwsConfigurationFile)
	if err != nil {
		return fmt.Errorf("unable to prepare aws credentials on configuration: %w", err)
	}

	var accessReplace, secretReplace, sessionReplace string
	dataAsString := string(data)
	if access != "" {
		accessReplace = fmt.Sprintf(`access_key = "%s"`, access)
	}
	dataAsString = strings.ReplaceAll(dataAsString, "{{ACCESS_KEY}}", accessReplace)

	if secret != "" {
		secretReplace = fmt.Sprintf(`secret_key = "%s"`, secret)
	}
	dataAsString = strings.ReplaceAll(dataAsString, "{{SECRET_KEY}}", secretReplace)

	if sessionToken != "" {
		sessionReplace = fmt.Sprintf(`session_token = "%s"`, sessionToken)
	}
	dataAsString = strings.ReplaceAll(dataAsString, "{{SESSION_TOKEN}}", sessionReplace)

	regionsEnvValue := os.Getenv(awsRegionsListParam)
	separatedRegions := strings.Split(regionsEnvValue, ",")
	regions := make([]string, len(separatedRegions))

	for i, region := range regions {
		if i < len(regions)-1 {
			regions[i] = fmt.Sprintf(`"%s",`, region)
		} else {
			regions[i] = fmt.Sprintf(`"%s"`, region)
		}
	}
	dataAsString = strings.ReplaceAll(dataAsString, "{{REGIONS}}", fmt.Sprintf(`regions = %s`, regions))

	if err = os.WriteFile(steampipeAwsConfigurationFile, []byte(dataAsString), 0o600); err != nil {
		return fmt.Errorf("unable to prepare aws config file: %w", err)
	}
	return nil
}
