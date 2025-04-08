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
	"github.com/spf13/cobra"
	"io"
	"strings"
)

// NewVersionCmd creates a new version command using Cobra and CommandBase
func NewVersionCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"version",
		"Prints Ballerina version",
		`Display version information for Ballerina tool and active distribution.

This command shows the version of the Ballerina CLI tool and the currently
active Ballerina distribution.`,
		printStream,
	)

	// Add example
	cmd.Example = `  # Print version information
  bal version`

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// Check for too many arguments
		if len(args) > 0 {
			panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("too many arguments", constants.BallerinaCliCommands.VERSION))
		}

		// Print version information
		cmdBase.PrintVersionInfo()
	}

	return cmd
}

// For backward compatibility with the existing command system

// VersionCommandStruct is a wrapper for backward compatibility
type VersionCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	VersionCommands []string
	HelpFlag        bool
	ParentCmdParser interface{}
}

// NewVersion creates a new VersionCommand for backward compatibility
func NewVersion(printStream io.Writer) *VersionCommandStruct {
	// Create the Cobra command
	cobraCmd := NewVersionCmd(printStream)

	// Create the wrapper
	cmd := &VersionCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *VersionCommandStruct) Execute() {
	if cmd.HelpFlag {
		// Ignore since we have nothing to print here.
		return
	}

	if cmd.VersionCommands == nil {
		cmd.cmdBase.PrintVersionInfo()
		return
	}

	if len(cmd.VersionCommands) > 0 {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("too many arguments", cmd.GetName()))
	}
}

// GetName returns the name of the command
func (cmd *VersionCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.VERSION
}

// PrintLongDesc prints the long description of the command
func (cmd *VersionCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *VersionCommandStruct) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal version \n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *VersionCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *VersionCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *VersionCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}
