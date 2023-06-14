package generators

import (
	"os"
	"strings"

	"github.com/blinkops/blink-steampipe/scripts/consts"
)

func setPluginVersion(dataAsString, defaultVersion string) string {
	pluginVersion := os.Getenv(consts.SteampipePluginVersionEnvVar)
	if pluginVersion == "" {
		pluginVersion = defaultVersion
	}

	return strings.ReplaceAll(dataAsString, "{{VERSION}}", pluginVersion)
}
