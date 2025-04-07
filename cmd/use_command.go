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

// UseCommand represents the "use" command and holds arguments and flags specified by the user
// Command name: "use", description: "Use Ballerina distribution"
type UseCommand struct {
	*Command
	UseCommands     []string // Command name
	HelpFlag        bool     // --help, -h, ?
	ParentCmdParser interface{}
}

// NewUse creates a new UseCommand
func NewUse(printStream io.Writer) *UseCommand {
	cmd := &UseCommand{}
	cmd.Command = NewWithWriter(printStream)
	return cmd
}

// Execute runs the command
func (cmd *UseCommand) Execute() {
	if cmd.HelpFlag {
		cmd.PrintUsageInfo(constants.CommandToolConstants.CliHelpFilePrefix + cmd.GetName())
		return
	}

	if cmd.UseCommands == nil || len(cmd.UseCommands) == 0 {
		panic(utils.ErrorUtil.CreateDistributionRequiredException("use"))
	}

	if len(cmd.UseCommands) > 1 {
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
	}

	printStream := cmd.GetPrintStream()
	distribution := cmd.UseCommands[0]
	if distribution == utils.ToolUtil.GetCurrentBallerinaVersion() {
		fmt.Fprintf(printStream, "'%s' is the current active distribution version\n", distribution)
		return
	}

	if utils.ToolUtil.CheckDistributionAvailable(distribution) {
		utils.ToolUtil.UseBallerinaVersion(printStream, distribution)
		fmt.Fprintf(printStream, "'%s' successfully set as the active distribution\n", distribution)
		return
	}
	fmt.Fprintf(printStream, "Distribution '%s' not found\n", distribution)

	channels := utils.ToolUtil.GetDistributions(printStream)
	validDistribution := false
	for _, channel := range channels {
		for _, dist := range channel.Distributions {
			if distribution == dist.Version {
				validDistribution = true
				fmt.Fprintf(printStream, "Run 'bal dist pull %s' to fetch and set the distribution as the active distribution\n", distribution)
				break
			}
		}
	}

	if !validDistribution {
		fmt.Fprintf(printStream, "'%s' is not a valid distribution. Use 'bal dist list -a' for the available distributions list\n", distribution)
	}
}

// GetName returns the name of the command
func (cmd *UseCommand) GetName() string {
	return constants.BallerinaCliCommands.USE
}

// PrintLongDesc prints the long description of the command
func (cmd *UseCommand) PrintLongDesc(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// PrintUsage prints the usage of the command
func (cmd *UseCommand) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal dist use\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *UseCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}
