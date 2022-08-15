package generators

import (
	"fmt"
	"os"
)

const (
	envFilePath = "/home/steampipe/.env"
)

type Variable struct {
	Key   string
	Value string
}

func WriteEnvFile(variable ...Variable) error {
	envFile := ""
	for _, info := range variable {
		envFile += fmt.Sprintf("export %s=%s\n", info.Key, info.Value)
	}

	if err := os.WriteFile(envFilePath, []byte(envFile), 0o600); err != nil {
		return fmt.Errorf("unable to prepare environment variables: %w", err)
	}

	return nil
}
