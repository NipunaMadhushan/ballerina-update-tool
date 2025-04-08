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

// NewBuildCmd creates a new build command using Cobra and CommandBase
func NewBuildCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"build",
		"Build a Ballerina program",
		`Build a Ballerina program package to an executable jar file.

This command checks for update notifications and then passes the build
command to the active Ballerina distribution.`,
		printStream,
	)

	// Add example
	cmd.Example = `  # Build a Ballerina program
  bal build program.bal`

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// Check for updates and notify user
		utils.ToolUtil.CheckForUpdate(cmdBase.GetPrintStream())

		// Note: In a real implementation, you would forward the build command
		// to the active Ballerina distribution here
	}

	return cmd
}

// For backward compatibility with the existing command system

// BuildCommandStruct is a wrapper for backward compatibility
type BuildCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	BuildCommands   []string
	HelpFlag        bool
	ParentCmdParser interface{}
}

// NewBuild creates a new BuildCommand for backward compatibility
func NewBuild(printStream io.Writer) *BuildCommandStruct {
	// Create the Cobra command
	cobraCmd := NewBuildCmd(printStream)

	// Create the wrapper
	cmd := &BuildCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *BuildCommandStruct) Execute() {
	if cmd.HelpFlag {
		return
	}

	// Check for updates and notify user
	utils.ToolUtil.CheckForUpdate(cmd.cmdBase.GetPrintStream())

	// Note: In a real implementation, you would forward the build command
	// to the active Ballerina distribution here
}

// GetName returns the name of the command
func (cmd *BuildCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.BUILD
}

// PrintLongDesc prints the long description of the command
func (cmd *BuildCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *BuildCommandStruct) PrintUsage(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// SetParentCmdParser sets the parent command parser
func (cmd *BuildCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *BuildCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *BuildCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}
