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

func SetEnv(variable ...Variable) error {
	envFile := ""
	for index, info := range variable {
		if index > 0 {
			envFile += "\n"
		}
		envFile += fmt.Sprintf("export %s=%s", info.Key, info.Value)
	}

	if err := os.WriteFile(envFilePath, []byte(envFile), 0o600); err != nil {
		return fmt.Errorf("unable to prepare environment variables: %w", err)
	}

	return nil
}
