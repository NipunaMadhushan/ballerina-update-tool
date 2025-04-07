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
	"sort"
	"strings"
)

// PullCommand represents the "pull" command and holds arguments and flags specified by the user
// Command name: "pull", description: "Pull Ballerina distribution"
type PullCommand struct {
	*Command
	PullCommands    []string // Command name
	HelpFlag        bool     // --help, -h, ?
	TestFlag        bool     // --test, -t
	ParentCmdParser interface{}
}

// NewPull creates a new PullCommand
func NewPull(printStream io.Writer) *PullCommand {
	cmd := &PullCommand{}
	cmd.Command = NewWithWriter(printStream)
	return cmd
}

// Execute runs the command
func (cmd *PullCommand) Execute() {
	if cmd.HelpFlag {
		cmd.PrintUsageInfo(constants.CommandToolConstants.CliHelpFilePrefix + cmd.GetName())
		return
	}

	if cmd.PullCommands == nil || len(cmd.PullCommands) == 0 {
		panic(utils.ErrorUtil.CreateDistributionRequiredException("pull"))
	}

	if len(cmd.PullCommands) > 1 {
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
	}

	printStream := cmd.GetPrintStream()
	distribution := cmd.PullCommands[0]

	if !cmd.TestFlag {
		// Check and update the tool if any latest version available
		toolDetails := utils.ToolUtil.UpdateTool(printStream)
		if toolDetails.Compatibility != "true" {
			return
		}
	}

	// To handle bal dist pull latest
	if distribution == constants.CommandToolConstants.LatestPullInput {
		fmt.Fprintln(printStream, "Fetching the latest distribution from the remote server...")
		channels := utils.ToolUtil.GetDistributions(printStream)
		if len(channels) == 0 {
			panic(utils.ErrorUtil.CreateCommandException("Failed to get distributions from server"))
		}
		// Assume channels are sorted descending
		latestChannel := channels[0]
		distributions := latestChannel.Distributions
		sort.Slice(distributions, func(i, j int) bool {
			return distributions[i].Version > distributions[j].Version
		})
		distribution = utils.ToolUtil.GetLatest(distributions[0].Version, "patch")
	}

	// To check whether the distribution is a valid one
	if distribution != constants.CommandToolConstants.LatestPullInput {
		validDist := false
		channels := utils.ToolUtil.GetDistributions(printStream)
		if len(channels) == 0 {
			panic(utils.ErrorUtil.CreateCommandException("Failed to get distributions from server"))
		}
		for _, channel := range channels {
			distributions := channel.Distributions
			for _, dist := range distributions {
				if dist.Version == distribution {
					validDist = true
					break
				}
			}
		}
		if !validDist {
			panic(utils.ErrorUtil.CreateDistributionNotFoundException(distribution))
		}
	}

	currentVersion := utils.ToolUtil.GetCurrentBallerinaVersion()
	if distribution == currentVersion {
		fmt.Fprintf(printStream, "'%s' is already the active distribution\n", distribution)
		return
	}

	utils.ToolUtil.HandleInstallDirPermission()
	alreadyAvailable := utils.ToolUtil.DownloadDistribution(printStream, distribution, utils.ToolUtil.GetType(distribution), distribution, cmd.TestFlag)

	if alreadyAvailable && distribution == currentVersion {
		return
	}
	utils.ToolUtil.UseBallerinaVersion(printStream, distribution)
	fmt.Fprintf(printStream, "'%s' successfully set as the active distribution\n", distribution)
}

// GetName returns the name of the command
func (cmd *PullCommand) GetName() string {
	return constants.BallerinaCliCommands.PULL
}

// PrintLongDesc prints the long description of the command
func (cmd *PullCommand) PrintLongDesc(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// PrintUsage prints the usage of the command
func (cmd *PullCommand) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal dist pull\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *PullCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}
