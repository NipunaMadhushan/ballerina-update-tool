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

// NewDistCmd creates a new distribution command using Cobra and CommandBase
func NewDistCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"dist",
		"Ballerina distribution commands",
		`Manage Ballerina distributions with various subcommands.

This command allows you to list, download, remove, update, and use different
Ballerina distributions.`,
		printStream,
	)

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// If no arguments, print help
		if len(args) == 0 {
			cmdBase.PrintUsageInfo(constants.BallerinaCliCommands.DIST)
			return
		}

		// If too many arguments
		if len(args) > 1 {
			panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("too many arguments", constants.BallerinaCliCommands.DIST))
		}

		// Unknown subcommand
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("unknown command '"+args[0]+"'", constants.BallerinaCliCommands.DIST))
	}

	// Add subcommands - these would be added when using this command
	// cmd.AddCommand(NewListCmd(printStream))
	// cmd.AddCommand(NewPullCmd(printStream))
	// cmd.AddCommand(NewRemoveCmd(printStream))
	// cmd.AddCommand(NewUpdateCmd(printStream))
	// cmd.AddCommand(NewUseCmd(printStream))

	return cmd
}

// For backward compatibility with the existing command system

// DistributionCommandStruct is a wrapper for backward compatibility
type DistributionCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	DistCommands    []string
	HelpFlag        bool
	ParentCmdParser interface{}
}

// NewDistribution creates a new DistributionCommand for backward compatibility
func NewDistribution(printStream io.Writer) *DistributionCommandStruct {
	// Create the Cobra command
	cobraCmd := NewDistCmd(printStream)

	// Create the wrapper
	cmd := &DistributionCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *DistributionCommandStruct) Execute() {
	// Handle help flag or no arguments
	if cmd.HelpFlag || cmd.DistCommands == nil {
		cmd.cmdBase.PrintUsageInfo(constants.BallerinaCliCommands.DIST)
		return
	}

	// Handle too many arguments
	if len(cmd.DistCommands) > 1 {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("too many arguments", cmd.GetName()))
	}

	// Handle unknown command
	// This would only happen if the subcommand isn't registered with Cobra
	panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("unknown command '"+cmd.DistCommands[0]+"'", cmd.GetName()))
}

// GetName returns the name of the command
func (cmd *DistributionCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.DIST
}

// PrintLongDesc prints the long description of the command
func (cmd *DistributionCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *DistributionCommandStruct) PrintUsage(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// SetParentCmdParser sets the parent command parser
func (cmd *DistributionCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *DistributionCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *DistributionCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}

// AddCommand adds a subcommand to this command
func (cmd *DistributionCommandStruct) AddCommand(subCmd interface{}) {
	// Extract cobra command from different subcommand types
	var cobraSubCmd *cobra.Command

	switch c := subCmd.(type) {
	case *cobra.Command:
		cobraSubCmd = c
	case interface{ GetCobraCommand() *cobra.Command }:
		cobraSubCmd = c.GetCobraCommand()
	default:
		// For command types that don't have a GetCobraCommand method,
		// log a warning or handle appropriately
		fmt.Fprintf(cmd.GetPrintStream(), "Warning: Cannot add command of type %T\n", c)
		return
	}

	// Add the subcommand to the cobra command
	cmd.cobraCmd.AddCommand(cobraSubCmd)
}
