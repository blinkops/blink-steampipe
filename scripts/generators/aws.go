package generators

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
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
)

const (
	awsUserBased          = "user_based"
	awsRoleBased          = "role_based"
	assumeWebIdentity     = "assume_web_identity"
	assumeCrossAccount    = "assume_cross_account"
	assumeTrustedIdentity = "assume_trusted_identity"
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

	base, subBase := gen.detect()
	switch base {
	case awsUserBased: // Implemented automatically via aws cli
	case awsRoleBased:
		sessionRegion := gen.getSessionRegion()
		roleArn, externalId := os.Getenv(awsRoleArn), os.Getenv(awsExternalID)

		svc := gen.initSTSClient(subBase, sessionRegion)

		access, secret, sessionToken, err := gen.assumeRole(svc, subBase, roleArn, externalId)
		if err != nil {
			return fmt.Errorf("unable to assume role with error: %w", err)
		}

		if err = os.Setenv(awsAccessKeyIdEnv, access); err != nil {
			return err
		}
		if err = os.Setenv(awsSecretAccessKeyEnv, secret); err != nil {
			return err
		}
		if err = os.Setenv(awsSessionToken, sessionToken); err != nil {
			return err
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

	if externalId := os.Getenv(awsExternalID); externalId != "" {
		return awsRoleBased, assumeTrustedIdentity
	}

	return awsRoleBased, assumeWebIdentity
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
	case assumeWebIdentity:
		return gen.assumeRoleWithWebIdentity(svc, role, sessionName)
	case assumeTrustedIdentity:
		return gen.assumeRoleWithTrustedIdentity(svc, role, externalID, sessionName)
	case assumeCrossAccount:
		return gen.assumeRoleCrossAccounts(svc, role, sessionName)
	}
	return "", "", "", errors.New("invalid assume role type was provided")
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
