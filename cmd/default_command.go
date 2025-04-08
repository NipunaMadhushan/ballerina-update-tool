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
	"io"
	"strings"

	"github.com/spf13/cobra"
)

// NewDefaultCmd creates a new default command using Cobra and CommandBase
func NewDefaultCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"bal",
		"Ballerina CLI tool",
		`Ballerina programming language command line tool.

This CLI helps you manage Ballerina distributions and update the tool.
Run "bal help <command>" for more information about a command.`,
		printStream,
	)

	var debugPort string
	var versionFlag bool

	// Add flags
	cmd.Flags().StringVar(&debugPort, "debug", "", "Start Ballerina in remote debugging mode")
	cmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print version information")

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// If version flag is set, print version info
		if versionFlag {
			cmdBase.PrintVersionInfo()
			return
		}

		// Otherwise, print help/usage info
		cmdBase.PrintUsageInfo(constants.BallerinaCliCommands.HELP)
	}

	return cmd
}

// For backward compatibility with the existing command system

// DefaultCommandStruct is a wrapper for backward compatibility
type DefaultCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	HelpFlag        bool
	DebugPort       string
	VersionFlag     bool
	HelpCommands    []string
	ParentCmdParser interface{}
}

// NewDefaultCommand creates a new DefaultCommand for backward compatibility
func NewDefaultCommand(printStream io.Writer) *DefaultCommandStruct {
	// Create the Cobra command
	cobraCmd := NewDefaultCmd(printStream)

	// Create the wrapper
	cmd := &DefaultCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *DefaultCommandStruct) Execute() {
	// Apply flags to the cobra command
	if cmd.DebugPort != "" {
		cmd.cobraCmd.Flags().Set("debug", cmd.DebugPort)
	}

	if cmd.VersionFlag {
		// If version flag is set, print version info directly
		cmd.cmdBase.PrintVersionInfo()
		return
	}

	// Otherwise, print help/usage info
	cmd.cmdBase.PrintUsageInfo(constants.BallerinaCliCommands.HELP)
}

// GetName returns the name of the command
func (cmd *DefaultCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.DEFAULT
}

// PrintLongDesc prints the long description of the command
func (cmd *DefaultCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *DefaultCommandStruct) PrintUsage(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// SetParentCmdParser sets the parent command parser
func (cmd *DefaultCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *DefaultCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *DefaultCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}
