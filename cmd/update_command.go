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

// NewUpdateCmd creates a new update command using Cobra and CommandBase
func NewUpdateCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"update",
		"Update Ballerina current distribution",
		`Update the current Ballerina distribution to the latest version.

This command automatically checks for and downloads the latest patch version
of your current Ballerina distribution and sets it as the active distribution.`,
		printStream,
	)

	// Define flags
	var testFlag bool

	// Add flags
	cmd.Flags().BoolVarP(&testFlag, "test", "t", false, "Update with a test distribution")

	// Add example
	cmd.Example = `  # Update to the latest patch version
  bal dist update

  # Update with test flag
  bal dist update --test`

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// Check for too many arguments
		if len(args) > 0 {
			panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", constants.BallerinaCliCommands.UPDATE))
		}

		// Handle permissions
		utils.ToolUtil.HandleInstallDirPermission()

		// Execute update command
		executeUpdate(cmdBase.GetPrintStream(), testFlag)
	}

	return cmd
}

// executeUpdate updates the Ballerina distribution to the latest version
func executeUpdate(printStream io.Writer, testFlag bool) {
	if !testFlag {
		// Check and update the tool if any latest version available
		toolDetails := utils.ToolUtil.UpdateTool(printStream)
		if toolDetails.Compatibility != "true" {
			return
		}
	}

	version := utils.ToolUtil.GetCurrentBallerinaVersion()
	fmt.Fprintln(printStream, "Fetching the latest distribution from the remote server...")

	latestVersion := utils.ToolUtil.GetLatest(version, "patch")
	if latestVersion == "" {
		fmt.Fprintln(printStream, "Failed to find the latest Ballerina distribution")
		return
	}

	if latestVersion != version {
		utils.ToolUtil.DownloadDistribution(printStream, latestVersion, utils.ToolUtil.GetType(latestVersion), latestVersion, testFlag)
		utils.ToolUtil.UseBallerinaVersion(printStream, latestVersion)
		fmt.Fprintf(printStream, "Successfully set the distribution '%s' as the active distribution\n", latestVersion)
		return
	}

	fmt.Fprintf(printStream, "The latest distribution '%s' is already the active distribution\n", latestVersion)
}

// For backward compatibility with the existing command system

// UpdateCommandStruct is a wrapper for backward compatibility
type UpdateCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	UpdateCommands  []string
	HelpFlag        bool
	TestFlag        bool
	ParentCmdParser interface{}
}

// NewUpdate creates a new UpdateCommand for backward compatibility
func NewUpdate(printStream io.Writer) *UpdateCommandStruct {
	// Create the Cobra command
	cobraCmd := NewUpdateCmd(printStream)

	// Create the wrapper
	cmd := &UpdateCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *UpdateCommandStruct) Execute() {
	if cmd.HelpFlag {
		cmd.cmdBase.PrintUsageInfo(constants.CommandToolConstants.CliHelpFilePrefix + cmd.GetName())
		return
	}

	// Apply flags to the cobra command
	cmd.cobraCmd.Flags().Set("test", fmt.Sprintf("%v", cmd.TestFlag))

	// Check arguments
	if cmd.UpdateCommands == nil {
		utils.ToolUtil.HandleInstallDirPermission()
		executeUpdate(cmd.cmdBase.GetPrintStream(), cmd.TestFlag)
		return
	}

	if len(cmd.UpdateCommands) > 0 {
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
	}
}

// GetName returns the name of the command
func (cmd *UpdateCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.UPDATE
}

// PrintLongDesc prints the long description of the command
func (cmd *UpdateCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *UpdateCommandStruct) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal dist update\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *UpdateCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *UpdateCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *UpdateCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}
