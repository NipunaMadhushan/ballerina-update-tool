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

// UpdateToolCommand represents the "update" command and holds arguments and flags specified by the user
// Command name: "command", description: "Update Ballerina current cli tool commands"
type UpdateToolCommand struct {
	*Command
	UpdateCommands  []string // args
	HelpFlag        bool     // --help, -h, ?
	ParentCmdParser interface{}
}

// NewUpdateTool creates a new UpdateToolCommand
func NewUpdateTool(printStream io.Writer) *UpdateToolCommand {
	cmd := &UpdateToolCommand{}
	cmd.Command = NewWithWriter(printStream)
	return cmd
}

// Execute runs the command
func (cmd *UpdateToolCommand) Execute() {
	if cmd.HelpFlag {
		cmd.PrintUsageInfo(constants.BallerinaCliCommands.UPDATE)
		return
	}

	if cmd.UpdateCommands == nil {
		utils.ToolUtil.HandleInstallDirPermission()
		updateCommands(cmd.GetPrintStream())
		return
	}

	if len(cmd.UpdateCommands) > 0 {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelpSubCommand("too many arguments", cmd.GetName()))
	}
}

// GetName returns the name of the command
func (cmd *UpdateToolCommand) GetName() string {
	return constants.BallerinaCliCommands.UPDATE
}

// PrintLongDesc prints the long description of the command
func (cmd *UpdateToolCommand) PrintLongDesc(out *strings.Builder) {
	out.WriteString("Updates the Ballerina tool to the latest version.\n")
}

// PrintUsage prints the usage of the command
func (cmd *UpdateToolCommand) PrintUsage(out *strings.Builder) {
	out.WriteString("  ballerina tool\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *UpdateToolCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
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
