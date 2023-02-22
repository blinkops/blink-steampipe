package consts

const (
	SteampipeBasePath             = "/home/steampipe/"
	SteampipeSpcConfigurationPath = SteampipeBasePath + ".steampipe/config/"
)

const (
	FilesMountPath                  = "/exec-files"
	ReportFileEnvVar                = "REPORT_FILE_NAME"
	FileOutputOverrideFormat        = "output written to file with identifier '%s'"
	FileOutputOverrideOnErrorFormat = FileOutputOverrideFormat + " on error"
)

const (
	CommandCheck = "check"
)
