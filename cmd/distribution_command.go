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

// DistributionCommand represents the "dist" command and holds arguments and flags specified by the user
// Command name: "dist", description: "Ballerina distribution commands"
type DistributionCommand struct {
	*Command
	DistCommands    []string // Command name
	HelpFlag        bool     // --help, -h, ?
	ParentCmdParser interface{}
}

// NewDistribution creates a new DistributionCommand
func NewDistribution(printStream io.Writer) *DistributionCommand {
	cmd := &DistributionCommand{}
	cmd.Command = NewWithWriter(printStream)
	return cmd
}

// Execute runs the command
func (cmd *DistributionCommand) Execute() {
	if cmd.HelpFlag || cmd.DistCommands == nil {
		cmd.PrintUsageInfo(constants.BallerinaCliCommands.DIST)
		return
	}

	if len(cmd.DistCommands) > 1 {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("too many arguments", cmd.GetName()))
	}

	panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("unknown command '"+cmd.DistCommands[0]+"'",
		cmd.GetName()))
}

// GetName returns the name of the command
func (cmd *DistributionCommand) GetName() string {
	return constants.BallerinaCliCommands.DIST
}

// PrintLongDesc prints the long description of the command
func (cmd *DistributionCommand) PrintLongDesc(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// PrintUsage prints the usage of the command
func (cmd *DistributionCommand) PrintUsage(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// SetParentCmdParser sets the parent command parser
func (cmd *DistributionCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}
