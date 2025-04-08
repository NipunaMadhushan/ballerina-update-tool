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

// NewUpdateToolCmd creates a new update tool command using Cobra and CommandBase
func NewUpdateToolCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"update",
		"Update Ballerina CLI tool",
		`Update the Ballerina CLI tool to the latest version.

This command checks for and downloads the latest version of the Ballerina CLI
tool itself. To update your Ballerina distribution, use 'bal dist update' instead.`,
		printStream,
	)

	// Add example
	cmd.Example = `  # Update the CLI tool to the latest version
  bal update`

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// Check for too many arguments
		if len(args) > 0 {
			panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("too many arguments", constants.BallerinaCliCommands.UPDATE))
		}

		// Handle permissions
		utils.ToolUtil.HandleInstallDirPermission()

		// Execute update tool command
		updateCommands(cmdBase.GetPrintStream())
	}

	return cmd
}

// updateCommands updates the CLI tool commands to the latest version
func updateCommands(printStream io.Writer) {
	version := utils.ToolUtil.GetCurrentToolsVersion()
	fmt.Fprintln(printStream, "Fetching the latest update tool version from the remote server...")

	latestVersionDetails := utils.ToolUtil.GetLatestToolVersion()
	latestVersion := latestVersionDetails.Version
	if latestVersion == "" {
		fmt.Fprintln(printStream, "Failed to find the latest update tool version")
		return
	}

	if latestVersion == version {
		fmt.Fprintf(printStream, "The latest update tool version '%s' is already in use\n", latestVersion)
		fmt.Fprintln(printStream, "\nIf you want to update the Ballerina distribution, use `bal dist update`")
		return
	}

	utils.ToolUtil.DownloadTool(printStream, latestVersion)
}

// For backward compatibility with the existing command system

// UpdateToolCommandStruct is a wrapper for backward compatibility
type UpdateToolCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	UpdateCommands  []string
	HelpFlag        bool
	ParentCmdParser interface{}
}

// NewUpdateTool creates a new UpdateToolCommand for backward compatibility
func NewUpdateTool(printStream io.Writer) *UpdateToolCommandStruct {
	// Create the Cobra command
	cobraCmd := NewUpdateToolCmd(printStream)

	// Create the wrapper
	cmd := &UpdateToolCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *UpdateToolCommandStruct) Execute() {
	if cmd.HelpFlag {
		cmd.cmdBase.PrintUsageInfo(constants.BallerinaCliCommands.UPDATE)
		return
	}

	// Check arguments
	if cmd.UpdateCommands == nil {
		utils.ToolUtil.HandleInstallDirPermission()
		updateCommands(cmd.cmdBase.GetPrintStream())
		return
	}

	if len(cmd.UpdateCommands) > 0 {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("too many arguments", cmd.GetName()))
	}
}

// GetName returns the name of the command
func (cmd *UpdateToolCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.UPDATE
}

// PrintLongDesc prints the long description of the command
func (cmd *UpdateToolCommandStruct) PrintLongDesc(out *strings.Builder) {
	out.WriteString("Updates the Ballerina tool to the latest version.\n")
}

// PrintUsage prints the usage of the command
func (cmd *UpdateToolCommandStruct) PrintUsage(out *strings.Builder) {
	out.WriteString("  ballerina tool\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *UpdateToolCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *UpdateToolCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *UpdateToolCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}
