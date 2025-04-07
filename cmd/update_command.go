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

// UpdateCommand represents the "update" command and holds arguments and flags specified by the user
// Command name: "command", description: "Update Ballerina current distribution"
type UpdateCommand struct {
	*Command
	UpdateCommands  []string // Command name
	HelpFlag        bool     // --help, -h, ?
	TestFlag        bool     // --test, -t
	ParentCmdParser interface{}
}

// NewUpdate creates a new UpdateCommand
func NewUpdate(printStream io.Writer) *UpdateCommand {
	cmd := &UpdateCommand{}
	cmd.Command = NewWithWriter(printStream)
	return cmd
}

// Execute runs the command
func (cmd *UpdateCommand) Execute() {
	if cmd.HelpFlag {
		cmd.PrintUsageInfo(constants.CommandToolConstants.CliHelpFilePrefix + cmd.GetName())
		return
	}

	if cmd.UpdateCommands == nil {
		utils.ToolUtil.HandleInstallDirPermission()
		Update(cmd.GetPrintStream(), cmd.TestFlag)
		return
	}

	if len(cmd.UpdateCommands) > 0 {
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
	}
}

// GetName returns the name of the command
func (cmd *UpdateCommand) GetName() string {
	return constants.BallerinaCliCommands.UPDATE
}

// PrintLongDesc prints the long description of the command
func (cmd *UpdateCommand) PrintLongDesc(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// PrintUsage prints the usage of the command
func (cmd *UpdateCommand) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal dist command\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *UpdateCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// Update updates the Ballerina distribution to the latest version
func Update(printStream io.Writer, testFlag bool) {
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
