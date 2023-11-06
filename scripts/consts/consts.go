package consts

const (
	SteampipeBasePath             = "/home/steampipe/"
	SteampipeSpcConfigurationPath = SteampipeBasePath + ".steampipe/config/"
)

const (
	SteampipePluginVersionEnvVar           = "PLUGIN_VERSION"
	FileIdentifierParentDirEnvVar          = "FILE_IDENTIFIER_PARENT_DIR"
	FileIdentifierEnvVar                   = "FILE_IDENTIFIER"
	SteampipeReportCustomModLocationEnvVar = "CUSTOM_MOD_LOCATION"
	FileOutputOverrideFormat               = "output written to file with identifier '%s'"
	FileOutputOverrideOnErrorFormat        = FileOutputOverrideFormat + " on error"
)

const (
	CommandCheck = "check"
)
