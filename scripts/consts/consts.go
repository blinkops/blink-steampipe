package consts

const (
	SteampipeBasePath             = "/home/steampipe/"
	SteampipeSpcConfigurationPath = SteampipeBasePath + ".steampipe/config/"
)

const (
	ReportFileParentDirEnvVar              = "REPORT_FILE_PARENT_DIR"
	ReportFilePathEnvVar                   = "REPORT_FILE_NAME"
	SteampipeReportCustomModLocationEnvVar = "CUSTOM_MOD_LOCATION"
	FileOutputOverrideFormat               = "output written to file with identifier '%s'"
	FileOutputOverrideOnErrorFormat        = FileOutputOverrideFormat + " on error"
)

const (
	CommandCheck = "check"
)
