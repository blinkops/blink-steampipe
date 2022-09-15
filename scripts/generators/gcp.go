package generators

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	gcpConnectionIdentifier          = "GCP_CONNECTION"
	gcpJsonCredential                = "GOOGLE_CREDENTIALS"
	gcpProjectIdKey                  = "project_id"
	steampipeGcpConfigurationFile    = "/home/steampipe/.steampipe/config/gcp.spc"
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

	data, err := os.ReadFile(steampipeGcpConfigurationFile)
	if err != nil {
		return fmt.Errorf("unable to prepare gcp credentials on configuration: %w", err)
	}

	dataAsString := strings.ReplaceAll(string(data), "{{PROJECT}}", projectIdAsString)
	if err = os.WriteFile(steampipeGcpConfigurationFile, []byte(dataAsString), 0o600); err != nil {
		return fmt.Errorf("unable to prepare gcp config file: %w", err)
	}

	if err = os.MkdirAll(gcpCredentialDirectoryPathFormat, 0o770); err != nil {
		return fmt.Errorf("unable to prepare gcp credentials path: %v", err)
	}

	filePath := gcpCredentialDirectoryPathFormat + gcpCredentialFileName
	if err = os.WriteFile(filePath, jsonData, 0o600); err != nil {
		return fmt.Errorf("unable to prepare gcp credentials: %w", err)
	}

	return nil
}
