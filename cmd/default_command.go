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
	"ballerina-update-tool/constants"
	"io"
	"strings"
)

// DefaultCommand represents the "default" command required by the CLI
// Command description: "Default Command."
type DefaultCommand struct {
	*Command
	HelpFlag        bool     // --help, -h, ?
	DebugPort       string   // --debug
	VersionFlag     bool     // --version, -v
	HelpCommands    []string // Help command name
	ParentCmdParser interface{}
}

// NewDefaultCommand creates a new DefaultCommand
func NewDefaultCommand(printStream io.Writer) *DefaultCommand {
	cmd := &DefaultCommand{}
	cmd.Command = NewWithWriter(printStream)
	return cmd
}

// Execute runs the command
func (cmd *DefaultCommand) Execute() {
	if cmd.VersionFlag {
		cmd.PrintVersionInfo()
		return
	}

	cmd.PrintUsageInfo(constants.BallerinaCliCommands.HELP)
}

// GetName returns the name of the command
func (cmd *DefaultCommand) GetName() string {
	return constants.BallerinaCliCommands.DEFAULT
}

// PrintLongDesc prints the long description of the command
func (cmd *DefaultCommand) PrintLongDesc(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// PrintUsage prints the usage of the command
func (cmd *DefaultCommand) PrintUsage(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// SetParentCmdParser sets the parent command parser
func (cmd *DefaultCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}
