package generators

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	log "github.com/sirupsen/logrus"
)

const (
	awsConnectionIdentifier     = "AWS_CONNECTION"
	awsAccessKeyId              = "ACCESS_KEY_ID"
	awsSecretAccessKey          = "SECRET_ACCESS_KEY"
	awsRoleArn                  = "ROLE_ARN"
	awsExternalID               = "EXTERNAL_ID"
	awsSessionToken             = "AWS_SESSION_TOKEN"
	awsWebIdentityTokenFile     = "AWS_WEB_IDENTITY_TOKEN_FILE"
	awsDefaultSessionRegion     = "eu-west-1"
	awsRegionEnvVariable        = "AWS_REGION"
	awsDefaultRegionEnvVariable = "AWS_DEFAULT_REGION"
	awsAccessKeyIdEnv           = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKeyEnv       = "AWS_SECRET_ACCESS_KEY"

	steampipeAwsConfigurationFile = "/home/steampipe/.steampipe/config/aws.spc"
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
		log.Tracef("failed resolving aws credentials, will try without credentials: %v", err)
	}
	return nil
}

func (gen AWSCredentialGenerator) generate() error {
	if _, ok := os.LookupEnv(awsConnectionIdentifier); !ok {
		return nil
	}

	var access, secret, sessionToken string

	base, subBase := gen.detect()
	switch base {
	case awsUserBased: // Implemented automatically via aws cli
		access, secret, sessionToken = os.Getenv(awsAccessKeyId), os.Getenv(awsSecretAccessKey), ""
	case awsRoleBased:
		sessionRegion := gen.getSessionRegion()
		roleArn, externalId := os.Getenv(awsRoleArn), os.Getenv(awsExternalID)

		svc := gen.initSTSClient(subBase, sessionRegion)

		var err error
		access, secret, sessionToken, err = gen.assumeRole(svc, subBase, roleArn, externalId)
		if err != nil {
			return fmt.Errorf("unable to assume role with error: %w", err)
		}

		//variables := []Variable{
		//	{
		//		Key:   awsAccessKeyIdEnv,
		//		Value: access,
		//	},
		//	{
		//		Key:   awsSecretAccessKeyEnv,
		//		Value: secret,
		//	},
		//	{
		//		Key:   awsSessionToken,
		//		Value: sessionToken,
		//	},
		//}
		//
		//return WriteEnvFile(variables...)
	default:
		return errors.New("invalid aws connection was provided")
	}

	data, err := os.ReadFile(steampipeAwsConfigurationFile)
	if err != nil {
		return fmt.Errorf("unable to prepare aws credentials on configuration: %w", err)
	}

	dataAsString := strings.ReplaceAll(string(data), "{{ACCESS_KEY}}", access)
	dataAsString = strings.ReplaceAll(dataAsString, "{{SECRET_KEY}}", secret)
	dataAsString = strings.ReplaceAll(dataAsString, "{{SESSION_TOKEN}}", sessionToken)

	if err = os.WriteFile(steampipeAwsConfigurationFile, []byte(dataAsString), 0o600); err != nil {
		return fmt.Errorf("unable to prepare aws config file: %w", err)
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

func (gen AWSCredentialGenerator) detect() (base, sub string) {
	if accessKeyId, secretAccessKey := os.Getenv(awsAccessKeyId), os.Getenv(awsSecretAccessKey); accessKeyId != "" && secretAccessKey != "" {
		base = awsUserBased
	}
	if roleArn := os.Getenv(awsRoleArn); roleArn == "" {
		return base, sub
	}
	if base == awsUserBased {
		return awsRoleBased, assumeCrossAccount
	}

	return awsRoleBased, assumeIdentity
}

func (gen AWSCredentialGenerator) initSTSClient(subBase, region string) stsiface.STSAPI {
	sessConfig := aws.Config{
		Region: aws.String(region),
	}
	if subBase == assumeCrossAccount {
		accessKeyId, secretAccessKey := os.Getenv(awsAccessKeyId), os.Getenv(awsSecretAccessKey)
		sessConfig.Credentials = credentials.NewStaticCredentials(accessKeyId, secretAccessKey, "")
	}

	sess, _ := session.NewSession(&sessConfig)
	return sts.New(sess)
}

func (gen AWSCredentialGenerator) assumeRole(svc stsiface.STSAPI, subBase, role, externalID string) (access, secret, sessionToken string, err error) {
	sessionName := strconv.Itoa(rand.Int())

	switch subBase {
	case assumeIdentity:
		return gen.assumeRoleWithIdentity(svc, role, externalID, sessionName)
	case assumeCrossAccount:
		return gen.assumeRoleCrossAccounts(svc, role, sessionName)
	}
	return "", "", "", errors.New("invalid assume role type was provided")
}

// assumeRoleWithIdentity first tries to assume a trusted identity using the role and external id, and if it doesn't
// succeed, it falls back to assuming a web identity using only the role
func (gen AWSCredentialGenerator) assumeRoleWithIdentity(svc stsiface.STSAPI, role, externalId, sessionName string) (string, string, string, error) {
	log.Debugf("assuming role with identity. Trying to assume role with trusted identity first and falling back to web identity")
	accessKey, secretAccessKey, sessionToken, err := gen.assumeRoleWithTrustedIdentity(svc, role, externalId, sessionName)
	if err != nil {
		log.Errorf("error assuming role with trusted identity: %v", err)
		return gen.assumeRoleWithWebIdentity(svc, role, sessionName)
	}

	return accessKey, secretAccessKey, sessionToken, nil
}

func (gen AWSCredentialGenerator) assumeRoleWithWebIdentity(svc stsiface.STSAPI, role, sessionName string) (string, string, string, error) {
	tokenFile, ok := os.LookupEnv(awsWebIdentityTokenFile)
	if !ok {
		log.Debug("token file for irsa not found. try assume role")
		result, err := svc.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         &role,
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
		RoleArn:          aws.String(role),
		RoleSessionName:  aws.String(sessionName),
		WebIdentityToken: aws.String(string(data)),
	}

	result, err := svc.AssumeRoleWithWebIdentity(input)
	if err != nil {
		return "", "", "", err
	}
	return *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken, err
}

func (gen AWSCredentialGenerator) assumeRoleWithTrustedIdentity(svc stsiface.STSAPI, role, externalID, sessionName string) (string, string, string, error) {
	input := &sts.AssumeRoleInput{
		RoleArn:         &role,
		RoleSessionName: &sessionName,
		ExternalId:      &externalID,
	}
	result, err := svc.AssumeRole(input)
	if err != nil {
		return "", "", "", err
	}
	return *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken, nil
}

func (gen AWSCredentialGenerator) assumeRoleCrossAccounts(svc stsiface.STSAPI, role, sessionName string) (string, string, string, error) {
	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         &role,
		RoleSessionName: &sessionName,
	})
	if err != nil {
		return "", "", "", err
	}
	return *result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken, err
}
