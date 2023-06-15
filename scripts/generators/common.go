package generators

import (
	"os"
	"os/exec"
	"strings"

	"github.com/blinkops/blink-steampipe/scripts/consts"
	"github.com/pkg/errors"
)

func setPluginVersion(dataAsString, defaultVersion string) (string, error) {
	version := os.Getenv(consts.SteampipePluginVersionEnvVar)
	if version == "" {
		version = defaultVersion
	}

	// if the user chooses to run with the latest version of the plugin, install it dynamically
	if strings.Contains(version, "latest") {
		cmd := exec.Command("steampipe", "install", "plugin", version)
		if output, err := cmd.CombinedOutput(); err != nil {
			return "", errors.New(string(output))
		}
	}

	return strings.ReplaceAll(dataAsString, "{{VERSION}}", version), nil
}
