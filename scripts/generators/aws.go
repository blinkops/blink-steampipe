package generators

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

const (
	awsConnectionIdentifier = "AWS_CONNECTION"
	awsAccessKeyId          = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKey      = "AWS_SECRET_ACCESS_KEY"
	awsRoleArn              = "AWS_ROLE_ARN"
	awsExternalID           = "AWS_EXTERNAL_ID"
	awsSessionToken         = "AWS_SESSION_TOKEN"
	awsWebIdentityTokenFile = "AWS_WEB_IDENTITY_TOKEN_FILE"
)

const (
	awsUserBased = "user_based"
	awsRoleBased = "role_based"
)

type AWSCredentialGenerator struct{}

func (gen AWSCredentialGenerator) Generate() error {
	if _, ok := os.LookupEnv(awsConnectionIdentifier); !ok {
		return nil
	}

	base, key, value := gen.detect()
	switch base {
	case awsUserBased: // Implemented automatically via aws cli
	case awsRoleBased:
		sess, _ := session.NewSession(&aws.Config{
			Region: aws.String("eu-west-1"),
		})

		svc := sts.New(sess)
		access, secret, sessionToken, err := gen.assumeRole(svc, key, value)
		if err != nil {
			return fmt.Errorf("unable to assume role with error: %w", err)
		}

		variables := []Variable{
			{
				Key:   awsAccessKeyId,
				Value: access,
			},
			{
				Key:   awsSecretAccessKey,
				Value: secret,
			},
			{
				Key:   awsSessionToken,
				Value: sessionToken,
			},
		}

		return WriteEnvFile(variables...)
	default:
		return errors.New("invalid aws connection was provided")
	}

	return nil
}

func (gen AWSCredentialGenerator) detect() (base, key, value string) {
	if accessKeyId, secretAccessKey := os.Getenv(awsAccessKeyId), os.Getenv(awsSecretAccessKey); accessKeyId != "" && secretAccessKey != "" {
		return awsUserBased, accessKeyId, secretAccessKey
	}
	if roleArn, externalId := os.Getenv(awsRoleArn), os.Getenv(awsExternalID); roleArn != "" {
		return awsRoleBased, roleArn, externalId
	}
	return "", "", ""
}

func (gen AWSCredentialGenerator) assumeRole(svc stsiface.STSAPI, role, externalID string) (access, secret, sessionToken string, err error) {
	sessionName := strconv.Itoa(rand.Int())
	if externalID == "" {
		return gen.assumeRoleWithWebIdentity(svc, role, sessionName)
	}
	return gen.assumeRoleWithTrustedIdentity(svc, role, externalID, sessionName)
}

func (gen AWSCredentialGenerator) assumeRoleWithWebIdentity(svc stsiface.STSAPI, role, sessionName string) (string, string, string, error) {
	tokenFile, ok := os.LookupEnv(awsWebIdentityTokenFile)
	if !ok {
		return "", "", "", fmt.Errorf("token file for irsa not found. make sure your pod is configured correctly and that your service account is created and annotated properly")
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
