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
	"sort"
	"strings"
)

// NewPullCmd creates a new pull command using Cobra and CommandBase
func NewPullCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"pull",
		"Pull Ballerina distribution",
		`Pull a specific Ballerina distribution from the remote repository.

This command downloads the specified Ballerina distribution from the remote server
and makes it available for use locally. The distribution is automatically set as 
the active distribution after download.`,
		printStream,
	)

	// Define flags
	var testFlag bool

	// Add flags
	cmd.Flags().BoolVarP(&testFlag, "test", "t", false, "Pull a test distribution")

	// Add example
	cmd.Example = `  # Pull the latest distribution
  bal dist pull latest

  # Pull a specific version
  bal dist pull 2.0.0

  # Pull a test distribution
  bal dist pull 2.0.0 --test`

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// Check for missing arguments
		if len(args) == 0 {
			panic(utils.ErrorUtil.CreateDistributionRequiredException("pull"))
		}

		// Check for too many arguments
		if len(args) > 1 {
			panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", constants.BallerinaCliCommands.PULL))
		}

		// Execute pull command
		executePull(args[0], testFlag, cmdBase.GetPrintStream())
	}

	return cmd
}

// executePull handles the execution of the pull command
func executePull(distribution string, testFlag bool, printStream io.Writer) {
	if !testFlag {
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
	alreadyAvailable := utils.ToolUtil.DownloadDistribution(printStream, distribution, utils.ToolUtil.GetType(distribution), distribution, testFlag)

	if alreadyAvailable && distribution == currentVersion {
		return
	}
	utils.ToolUtil.UseBallerinaVersion(printStream, distribution)
	fmt.Fprintf(printStream, "'%s' successfully set as the active distribution\n", distribution)
}

// For backward compatibility with the existing command system

// PullCommandStruct is a wrapper for backward compatibility
type PullCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	PullCommands    []string
	HelpFlag        bool
	TestFlag        bool
	ParentCmdParser interface{}
}

// NewPull creates a new PullCommand for backward compatibility
func NewPull(printStream io.Writer) *PullCommandStruct {
	// Create the Cobra command
	cobraCmd := NewPullCmd(printStream)

	// Create the wrapper
	cmd := &PullCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *PullCommandStruct) Execute() {
	if cmd.HelpFlag {
		cmd.cmdBase.PrintUsageInfo(constants.CommandToolConstants.CliHelpFilePrefix + cmd.GetName())
		return
	}

	// Apply flags to the cobra command
	cmd.cobraCmd.Flags().Set("test", fmt.Sprintf("%v", cmd.TestFlag))

	// Check for missing arguments
	if cmd.PullCommands == nil || len(cmd.PullCommands) == 0 {
		panic(utils.ErrorUtil.CreateDistributionRequiredException("pull"))
	}

	// Check for too many arguments
	if len(cmd.PullCommands) > 1 {
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
	}

	// Execute pull command
	executePull(cmd.PullCommands[0], cmd.TestFlag, cmd.cmdBase.GetPrintStream())
}

// GetName returns the name of the command
func (cmd *PullCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.PULL
}

// PrintLongDesc prints the long description of the command
func (cmd *PullCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *PullCommandStruct) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal dist pull\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *PullCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *PullCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *PullCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}
