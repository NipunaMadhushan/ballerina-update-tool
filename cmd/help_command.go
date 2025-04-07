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
	"fmt"
	"io"
	"strings"
)

// HelpCommand represents the "help" command and holds arguments and flags specified by the user
// Command name: "help", description: "print usage information"
type HelpCommand struct {
	*Command
	HelpCommands    []string // Command name
	ParentCmdParser interface{}
}

// NewHelp creates a new HelpCommand
func NewHelp(printStream io.Writer) *HelpCommand {
	cmd := &HelpCommand{}
	cmd.Command = NewWithWriter(printStream)
	return cmd
}

// Execute runs the command
func (cmd *HelpCommand) Execute() {
	if cmd.HelpCommands == nil {
		cmd.PrintUsageInfo(constants.BallerinaCliCommands.HELP)
		return
	}

	cmdCount := len(cmd.HelpCommands)

	if cmdCount > 2 {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelp("too many arguments"))
	}

	userCommand := cmd.HelpCommands[0]

	if cmdCount == 2 {
		userCommand += "-" + cmd.HelpCommands[1]
	}

	commandUsageInfo := cmd.GetCommandUsageInfo(userCommand)
	fmt.Fprintln(cmd.GetPrintStream(), commandUsageInfo)
}

// GetName returns the name of the command
func (cmd *HelpCommand) GetName() string {
	return constants.BallerinaCliCommands.HELP
}

// PrintLongDesc prints the long description of the command
func (cmd *HelpCommand) PrintLongDesc(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// PrintUsage prints the usage of the command
func (cmd *HelpCommand) PrintUsage(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// SetParentCmdParser sets the parent command parser
func (cmd *HelpCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}
