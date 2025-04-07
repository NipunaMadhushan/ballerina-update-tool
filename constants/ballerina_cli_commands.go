/*
 * Copyright (c) 2019, WSO2 Inc. (http://wso2.com) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package constants

// BallerinaCliCommands contains all the CLI command constants
var BallerinaCliCommands = struct {
	DEFAULT string
	BUILD   string
	HELP    string
	UPDATE  string
	DIST    string
	LIST    string
	PULL    string
	USE     string
	REMOVE  string
	VERSION string
}{
	DEFAULT: "default",
	BUILD:   "build",
	HELP:    "help",
	UPDATE:  "update",
	DIST:    "dist",
	LIST:    "list",
	PULL:    "pull",
	USE:     "use",
	REMOVE:  "remove",
	VERSION: "version",
}
