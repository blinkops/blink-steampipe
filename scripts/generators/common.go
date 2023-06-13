package generators

import (
	"os"
	"strings"
)

const steampipePluginVersionEnvVar = "PLUGIN_VERSION"

func setPluginVersion(dataAsString, defaultVersion string) string {
	pluginVersion := os.Getenv(steampipePluginVersionEnvVar)
	if pluginVersion == "" {
		pluginVersion = defaultVersion
	}

	return strings.ReplaceAll(dataAsString, "{{VERSION}}", pluginVersion)
}
