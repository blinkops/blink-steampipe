package generators

import (
	"fmt"
	"os"
)

const (
	crowdstrikeDomain = "FALCON_CLOUD"
)

var addressToDomainMap = map[string]string{
	"https://api.crowdstrike.com":            "us-1",
	"https://api.us-2.crowdstrike.com":       "us-2",
	"https://api.laggar.gcw.crowdstrike.com": "us-gov-1",
	"https://api.eu-1.crowdstrike.com":       "eu-1",
}

type CrowdstrikeCredentialGenerator struct{}

func (gen CrowdstrikeCredentialGenerator) Generate() error {
	apiAddress, ok := os.LookupEnv(crowdstrikeDomain)
	if !ok {
		return fmt.Errorf("invalid crowdstrike connection was provided")
	}
	domain, ok := addressToDomainMap[apiAddress]
	if !ok {
		return fmt.Errorf("invalid crowdstrike connection was provided")
	}
	return os.Setenv(crowdstrikeDomain, domain)
}
