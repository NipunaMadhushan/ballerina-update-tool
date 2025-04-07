package constants

// CommandToolConstants contains constants related to the CLI
var CommandToolConstants = struct {
	ProductionUrl         string
	StagingUrl            string
	DevUrl                string
	PreRelease            string
	CliHelpFilePrefix     string
	Ballerina1XVersions   string
	BallerinaSettingsFile string
	LatestPullInput       string
}{
	ProductionUrl:         "https://api.central.ballerina.io/2.0/update-tool",
	StagingUrl:            "https://api.staging-central.ballerina.io/2.0/update-tool",
	DevUrl:                "https://api.dev-central.ballerina.io/2.0/update-tool/",
	PreRelease:            "pre-release",
	CliHelpFilePrefix:     "dist-",
	Ballerina1XVersions:   "1.0.",
	BallerinaSettingsFile: "Settings.toml",
	LatestPullInput:       "latest",
}
