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

package cmd

import (
	"strings"
)

// BCommand represents a Ballerina launcher command.
type BCommand interface {
	// Execute runs the command.
	Execute()

	// GetName retrieves the command name.
	GetName() string

	// PrintLongDesc prints the detailed description of the command.
	PrintLongDesc(out *strings.Builder)

	// PrintUsage prints usage info for the command.
	PrintUsage(out *strings.Builder)

	// SetParentCmdParser sets the parent command line parser.
	// In Go we use interface{} since we don't have a direct equivalent of CommandLine.
	SetParentCmdParser(parentCmdParser interface{})
}
