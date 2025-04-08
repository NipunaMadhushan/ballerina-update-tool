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

package cmd

import (
	"ballerina-update-tool/constants"
	"ballerina-update-tool/utils"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"strings"
)

// NewUseCmd creates a new use command using Cobra and CommandBase
func NewUseCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"use",
		"Use Ballerina distribution",
		`Set the active Ballerina distribution to the specified version.

This command changes which Ballerina distribution is used when you run Ballerina programs.
The distribution must be already installed locally. If not, you'll be prompted to pull it.`,
		printStream,
	)

	// Add example
	cmd.Example = `  # Set a specific distribution as active
  bal dist use 2.0.0`

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// Check for missing arguments
		if len(args) == 0 {
			panic(utils.ErrorUtil.CreateDistributionRequiredException("use"))
		}

		// Check for too many arguments
		if len(args) > 1 {
			panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", constants.BallerinaCliCommands.USE))
		}

		// Execute use command
		executeUse(args[0], cmdBase.GetPrintStream())
	}

	return cmd
}

// executeUse handles the execution of the use command
func executeUse(distribution string, printStream io.Writer) {
	// Check if this is already the current version
	if distribution == utils.ToolUtil.GetCurrentBallerinaVersion() {
		fmt.Fprintf(printStream, "'%s' is the current active distribution version\n", distribution)
		return
	}

	// Check if the distribution is available locally
	if utils.ToolUtil.CheckDistributionAvailable(distribution) {
		utils.ToolUtil.UseBallerinaVersion(printStream, distribution)
		fmt.Fprintf(printStream, "'%s' successfully set as the active distribution\n", distribution)
		return
	}

	// Distribution not found locally, inform the user
	fmt.Fprintf(printStream, "Distribution '%s' not found\n", distribution)

	// Check if it's available for download
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

	// If not a valid distribution, suggest listing available distributions
	if !validDistribution {
		fmt.Fprintf(printStream, "'%s' is not a valid distribution. Use 'bal dist list -a' for the available distributions list\n", distribution)
	}
}

// For backward compatibility with the existing command system

// UseCommandStruct is a wrapper for backward compatibility
type UseCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	UseCommands     []string
	HelpFlag        bool
	ParentCmdParser interface{}
}

// NewUse creates a new UseCommand for backward compatibility
func NewUse(printStream io.Writer) *UseCommandStruct {
	// Create the Cobra command
	cobraCmd := NewUseCmd(printStream)

	// Create the wrapper
	cmd := &UseCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *UseCommandStruct) Execute() {
	if cmd.HelpFlag {
		cmd.cmdBase.PrintUsageInfo(constants.CommandToolConstants.CliHelpFilePrefix + cmd.GetName())
		return
	}

	// Check for missing arguments
	if cmd.UseCommands == nil || len(cmd.UseCommands) == 0 {
		panic(utils.ErrorUtil.CreateDistributionRequiredException("use"))
	}

	// Check for too many arguments
	if len(cmd.UseCommands) > 1 {
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
	}

	// Execute use command with the first argument
	executeUse(cmd.UseCommands[0], cmd.cmdBase.GetPrintStream())
}

// GetName returns the name of the command
func (cmd *UseCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.USE
}

// PrintLongDesc prints the long description of the command
func (cmd *UseCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *UseCommandStruct) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal dist use\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *UseCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *UseCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *UseCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}
