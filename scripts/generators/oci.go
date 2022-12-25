package generators

import (
	"fmt"
	"os"
	"strings"
)

const (
	ociConnectionIdentifier = "OCI_CONNECTION"
	ociTenancyOcid          = "TENANCY_OCID"
	ociUserOcid             = "USER_OCID"
	ociFingerprint          = "FINGERPRINT"
	ociPkey                 = "PKEY"
	ociAPIAddress           = "API_ADDRESS"

	ociSteampipeConfigurationFile = "/home/steampipe/.steampipe/config/oci.spc"

	ociPkeyFileDirPath = "/home/.ssh/"
	ociPkeyFile        = "oci_private.pem"
)

type OCICredentialGenerator struct{}

func (gen OCICredentialGenerator) Generate() error {
	if _, ok := os.LookupEnv(ociConnectionIdentifier); !ok {
		return nil
	}

	return gen.generateJSONCredentials()
}

func (gen OCICredentialGenerator) generateJSONCredentials() error {
	tenancyOcid, ok := os.LookupEnv(ociTenancyOcid)
	if !ok {
		return fmt.Errorf("invalid oci connection was provided")
	}
	userOcid, ok := os.LookupEnv(ociUserOcid)
	if !ok {
		return fmt.Errorf("invalid oci connection was provided")
	}
	fingerprint, ok := os.LookupEnv(ociFingerprint)
	if !ok {
		return fmt.Errorf("invalid oci connection was provided")
	}
	pkey, ok := os.LookupEnv(ociPkey)
	if !ok {
		return fmt.Errorf("invalid oci connection was provided")
	}
	apiAddress, ok := os.LookupEnv(ociAPIAddress)
	if !ok {
		return fmt.Errorf("invalid oci connection was provided")
	}
	splitAPIAddress := strings.Split(apiAddress, ".")
	if len(splitAPIAddress) < 2 {
		return fmt.Errorf("invalid oci connection was provided")
	}
	region := splitAPIAddress[1] // region is extracted from the API address - same behavior as HTTP

	// we write private RSA key to a file in the specified path we hard-coded to the oci.spc file
	if err := os.MkdirAll(ociPkeyFileDirPath, 0o770); err != nil {
		return fmt.Errorf("unable to prepare oci credentials path: %v", err)
	}
	if err := os.WriteFile(ociPkeyFileDirPath+ociPkeyFile, []byte(pkey), 0o600); err != nil {
		return fmt.Errorf("unable to prepare oci pkey file: %w", err)
	}

	// add all connection params to the oci.spc file before overwriting it
	data, err := os.ReadFile(ociSteampipeConfigurationFile)
	if err != nil {
		return fmt.Errorf("unable to prepare gcp credentials on configuration: %w", err)
	}

	var paramsReplacer = strings.NewReplacer("{{TENANCY_OCID}}", tenancyOcid, "{{USER_OCID}}", userOcid, "{{FINGERPRINT}}", fingerprint, "{{REGION}}", region)
	dataAsString := paramsReplacer.Replace(string(data))
	if err = os.WriteFile(ociSteampipeConfigurationFile, []byte(dataAsString), 0o600); err != nil {
		return fmt.Errorf("unable to prepare oci config file: %w", err)
	}

	return nil
}
