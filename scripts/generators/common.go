package generators

import (
	"os"
	"strings"
)

const steampipePluginVersionParam = "PLUGIN_VERSION_PARAM"

func setPluginVersion(dataAsString, defaultVersion string) string {
	pluginVersion := os.Getenv(steampipePluginVersionParam)
	if pluginVersion == "" {
		pluginVersion = defaultVersion
	}

	return strings.ReplaceAll(dataAsString, "{{VERSION}}", pluginVersion)
}
