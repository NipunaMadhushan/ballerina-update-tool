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
	"ballerina-update-tool/utils"
	"io"
	"strings"
)

// VersionCommand represents the "version" command and holds arguments and flags specified by the user
// Command name: "version", description: "Prints Ballerina version"
type VersionCommand struct {
	*Command
	VersionCommands []string // Command name
	HelpFlag        bool     // --help, -h, ?
	ParentCmdParser interface{}
}

// NewVersion creates a new VersionCommand
func NewVersion(printStream io.Writer) *VersionCommand {
	cmd := &VersionCommand{}
	cmd.Command = NewWithWriter(printStream)
	return cmd
}

// Execute runs the command
func (cmd *VersionCommand) Execute() {
	if cmd.HelpFlag {
		// Ignore since we have nothing to print here.
		return
	}

	if cmd.VersionCommands == nil {
		cmd.PrintVersionInfo()
		return
	}

	if len(cmd.VersionCommands) > 0 {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("too many arguments", cmd.GetName()))
	}
}

// GetName returns the name of the command
func (cmd *VersionCommand) GetName() string {
	return constants.BallerinaCliCommands.VERSION
}

// PrintLongDesc prints the long description of the command
func (cmd *VersionCommand) PrintLongDesc(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// PrintUsage prints the usage of the command
func (cmd *VersionCommand) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal version \n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *VersionCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}
