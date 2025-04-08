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

// NewHelpCmd creates a new help command using Cobra and CommandBase
func NewHelpCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"help",
		"Print usage information",
		`Display help information about Ballerina commands.

This command displays detailed help information for any Ballerina command.
You can get help on a specific command or subcommand by specifying it as an argument.`,
		printStream,
	)

	// Set example
	cmd.Example = `  # Get general help
  bal help

  # Get help for a specific command
  bal help dist

  # Get help for a subcommand
  bal help dist pull`

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// Execute help command
		executeHelp(cmdBase, cobraCmd, args)
	}

	return cmd
}

// executeHelp handles the execution of the help command
func executeHelp(cmdBase *CommandBase, cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		// If no arguments, print general help info
		cmdBase.PrintUsageInfo(constants.BallerinaCliCommands.HELP)
		return
	}

	// Check for too many arguments
	if len(args) > 2 {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelp("too many arguments"))
	}

	// Get the root command to find help topics
	rootCmd := cmd.Root()

	// Combine arguments if needed
	userCommand := args[0]
	if len(args) == 2 {
		userCommand += "-" + args[1]
	}

	// First try to find a command with this name
	var foundCmd *cobra.Command
	if len(args) == 1 {
		// Look for a single command
		foundCmd, _, _ = rootCmd.Find(args)
	} else if len(args) == 2 {
		// Look for a subcommand
		parentCmd, _, _ := rootCmd.Find([]string{args[0]})
		if parentCmd != nil && parentCmd != rootCmd {
			foundCmd, _, _ = parentCmd.Find([]string{args[1]})
		}
	}

	// If we found a command, use its help
	if foundCmd != nil && foundCmd != rootCmd {
		foundCmd.Help()
		return
	}

	// Otherwise, try to get help from a help file
	commandUsageInfo := cmdBase.GetCommandUsageInfo(userCommand)
	fmt.Fprintln(cmdBase.GetPrintStream(), commandUsageInfo)
}

// For backward compatibility with the existing command system

// HelpCommandStruct is a wrapper for backward compatibility
type HelpCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	HelpCommands    []string
	ParentCmdParser interface{}
}

// NewHelp creates a new HelpCommand for backward compatibility
func NewHelp(printStream io.Writer) *HelpCommandStruct {
	// Create the Cobra command
	cobraCmd := NewHelpCmd(printStream)

	// Create the wrapper
	cmd := &HelpCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *HelpCommandStruct) Execute() {
	// If no help commands, print general help
	if cmd.HelpCommands == nil {
		cmd.cmdBase.PrintUsageInfo(constants.BallerinaCliCommands.HELP)
		return
	}

	// Check for too many arguments
	cmdCount := len(cmd.HelpCommands)
	if cmdCount > 2 {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelp("too many arguments"))
	}

	// Combine arguments if needed
	userCommand := cmd.HelpCommands[0]
	if cmdCount == 2 {
		userCommand += "-" + cmd.HelpCommands[1]
	}

	// Get command usage info
	commandUsageInfo := cmd.cmdBase.GetCommandUsageInfo(userCommand)
	fmt.Fprintln(cmd.cmdBase.GetPrintStream(), commandUsageInfo)
}

// GetName returns the name of the command
func (cmd *HelpCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.HELP
}

// PrintLongDesc prints the long description of the command
func (cmd *HelpCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *HelpCommandStruct) PrintUsage(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// SetParentCmdParser sets the parent command parser
func (cmd *HelpCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *HelpCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *HelpCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}
