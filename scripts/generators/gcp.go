package generators

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	gcpConnectionIdentifier          = "GCP_CONNECTION"
	gcpJsonCredential                = "GOOGLE_CREDENTIALS"
	gcpProjectIdKey                  = "project_id"
	cloudSdkProjectEnvVariable       = "CLOUDSDK_CORE_PROJECT"
	gcpCredentialDirectoryPathFormat = "/home/steampipe/.config/gcloud/"
	gcpCredentialFileName            = "application_default_credentials.json"
)

type GCPCredentialGenerator struct{}

func (gen GCPCredentialGenerator) Generate() error {
	if _, ok := os.LookupEnv(gcpConnectionIdentifier); !ok {
		return nil
	}

	return gen.generateJSONCredentials()
}

func (gen GCPCredentialGenerator) generateJSONCredentials() error {
	jsonValue, ok := os.LookupEnv(gcpJsonCredential)
	if !ok {
		return fmt.Errorf("invalid gcp connection was provided")
	}

	jsonData := []byte(jsonValue)

	credentials := map[string]any{}
	if err := json.Unmarshal(jsonData, &credentials); err != nil {
		return fmt.Errorf("unable to parse gcp credentials with error: %w", err)
	}

	projectId, ok := credentials[gcpProjectIdKey]
	if !ok {
		return fmt.Errorf("unable to fetch project id from provided connection")
	}

	projectIdAsString, ok := projectId.(string)
	if !ok {
		return fmt.Errorf("invalid project id fetched from provided connection")
	}

	if err := os.MkdirAll(gcpCredentialDirectoryPathFormat, 0o770); err != nil {
		return fmt.Errorf("unable to prepare gcp credentials path: %v", err)
	}

	filePath := gcpCredentialDirectoryPathFormat + gcpCredentialFileName
	if err := os.WriteFile(filePath, jsonData, 0o600); err != nil {
		return fmt.Errorf("unable to prepare gcp credentials: %w", err)
	}

	variables := []Variable{
		{
			Key:   cloudSdkProjectEnvVariable,
			Value: projectIdAsString,
		},
	}
	return WriteEnvFile(variables...)
}
