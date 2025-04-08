/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

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
