package generators

import (
	"encoding/json"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"
)

const (
	gcpConnectionIdentifier      = "GCP_CONNECTION"
	gcpJsonCredential            = "GOOGLE_CREDENTIALS"
	gcpProjectIdKey              = "project_id"
	cloudSdkProjectEnvVariable   = "CLOUDSDK_CORE_PROJECT"
	gcpCredentialPathEnvVariable = "GOOGLE_APPLICATION_CREDENTIALS"
	gcpCredentialPathFormat      = "/tmp/%s/creds.json"
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

	if err := os.Setenv(cloudSdkProjectEnvVariable, projectIdAsString); err != nil {
		return fmt.Errorf("unable to set gcp project: %w", err)
	}

	path := fmt.Sprintf(gcpCredentialPathFormat, uuid.NewV4().String())
	if err := os.WriteFile(path, jsonData, 0o600); err != nil {
		return fmt.Errorf("unable to prepare gcp credentials: %w", err)
	}

	if err := os.Setenv(gcpCredentialPathEnvVariable, path); err != nil {
		return fmt.Errorf("unable to set gcp credentials path: %w", err)
	}
	return nil
}
