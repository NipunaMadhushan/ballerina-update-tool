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
	"os"
	"path/filepath"
	"strings"
)

// NewRemoveCmd creates a new remove command using Cobra and CommandBase
func NewRemoveCmd(printStream io.Writer) *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"remove",
		"Remove Ballerina distribution",
		`Remove a specific Ballerina distribution from your local environment.

This command removes a previously installed Ballerina distribution. You cannot
remove the currently active distribution. Use the --all flag to remove all
non-active distributions at once.`,
		printStream,
	)

	// Define flags
	var allFlag bool

	// Add flags
	cmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Remove all non-active distributions")

	// Add example
	cmd.Example = `  # Remove a specific distribution
  bal dist remove 2.0.0

  # Remove all non-active distributions
  bal dist remove --all`

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Handle panic recovery
		defer HandlePanic()

		// Check permissions first
		utils.ToolUtil.HandleInstallDirPermission()

		// Handle different scenarios based on flags and arguments
		if allFlag {
			// When using --all flag
			if len(args) > 0 {
				panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", constants.BallerinaCliCommands.REMOVE))
			}

			// Execute removeAll
			executeRemoveAll(cmdBase.GetPrintStream())
			return
		}

		// When not using --all flag, a distribution argument is required
		if len(args) == 0 {
			panic(utils.ErrorUtil.CreateDistributionRequiredExceptionWithFlags("remove", "--all, -a"))
		}

		// Check for too many arguments
		if len(args) > 1 {
			panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", constants.BallerinaCliCommands.REMOVE))
		}

		// Execute remove for a specific distribution
		executeRemove(args[0], cmdBase.GetPrintStream())
	}

	return cmd
}

// executeRemove removes a specific distribution version
func executeRemove(version string, printStream io.Writer) {
	if isCurrentVersion(version) {
		panic(utils.ErrorUtil.CreateCommandException("The active Ballerina distribution cannot be removed"))
	} else {
		file := utils.ToolUtil.GetType(version) + "-" + version
		directory := filepath.Join(utils.ToolUtil.GetDistributionsPath(), file)

		if _, err := os.Stat(directory); !os.IsNotExist(err) {
			err := utils.OSUtils.DeleteFiles(directory)
			if err != nil {
				panic(utils.ErrorUtil.CreateCommandException("error occurred while removing '" + version + "': " + err.Error()))
			}

			err = utils.OSUtils.DeleteCaches(version, printStream)
			if err != nil {
				panic(utils.ErrorUtil.CreateCommandException("error occurred while removing caches for '" + version + "': " + err.Error()))
			}

			fmt.Fprintf(printStream, "Distribution '%s' successfully removed\n", version)
			utils.ToolUtil.RemoveUnusedDependencies(version, printStream)
		} else {
			panic(utils.ErrorUtil.CreateCommandException("distribution '" + version + "' not found"))
		}
	}
}

// executeRemoveAll removes all distributions except the active one
func executeRemoveAll(printStream io.Writer) {
	folder := utils.ToolUtil.GetDistributionsPath()
	listOfFiles, err := os.ReadDir(folder)
	if err != nil {
		panic(utils.ErrorUtil.CreateCommandException("error occurred while reading distributions: " + err.Error()))
	}

	// checking for 2 files for zip pack and 3 files for installers
	if len(listOfFiles) == 2 || (len(listOfFiles) == 3 &&
		fileExists(filepath.Join(folder, "installer-version"))) {
		fmt.Fprintln(printStream, "There is nothing to remove. Only active distribution is remaining")
		return
	}

	for _, file := range listOfFiles {
		if file.IsDir() {
			version := ""
			fileName := file.Name()
			parts := strings.Split(fileName, "-")
			if len(parts) == 2 {
				version = parts[1]
			}

			directory := filepath.Join(utils.ToolUtil.GetDistributionsPath(), fileName)
			if fileExists(directory) && (!isCurrentVersion(version) || version == "") {
				err := utils.OSUtils.DeleteFiles(directory)
				if err != nil {
					panic(utils.ErrorUtil.CreateCommandException("error occurred while removing distribution '" + fileName + "': " + err.Error()))
				}
			}
		}
	}

	fmt.Fprintln(printStream, "All non-active distributions are successfully removed")

	activeDistribution := utils.ToolUtil.GetCurrentBallerinaVersion()
	dependencyForActiveDistribution := utils.ToolUtil.GetDependency(
		printStream,
		activeDistribution,
		utils.ToolUtil.GetType(activeDistribution),
		activeDistribution)

	dependencies, err := os.ReadDir(utils.ToolUtil.GetDependencyPath())
	if err != nil {
		panic(utils.ErrorUtil.CreateCommandException("error occurred while reading dependencies: " + err.Error()))
	}

	if len(dependencies) > 1 {
		fmt.Fprintln(printStream, "Removing unused dependencies")
		for _, dependency := range dependencies {
			if dependency.IsDir() && dependency.Name() != dependencyForActiveDistribution {
				depPath := filepath.Join(utils.ToolUtil.GetDependencyPath(), dependency.Name())
				err := utils.OSUtils.DeleteFiles(depPath)
				if err != nil {
					panic(utils.ErrorUtil.CreateCommandException("error occurred while removing dependency '" + dependency.Name() + "': " + err.Error()))
				}
			}
		}
	} else if len(dependencies) == 0 {
		panic(utils.ErrorUtil.CreateCommandException("No dependencies found"))
	}
}

// isCurrentVersion checks if the given version is the current active version
func isCurrentVersion(version string) bool {
	return version == utils.ToolUtil.GetCurrentBallerinaVersion()
}

// fileExists checks if a file or directory exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// For backward compatibility with the existing command system

// RemoveCommandStruct is a wrapper for backward compatibility
type RemoveCommandStruct struct {
	cmdBase         *CommandBase
	cobraCmd        *cobra.Command
	RemoveCommands  []string
	HelpFlag        bool
	AllFlag         bool
	ParentCmdParser interface{}
}

// NewRemove creates a new RemoveCommand for backward compatibility
func NewRemove(printStream io.Writer) *RemoveCommandStruct {
	// Create the Cobra command
	cobraCmd := NewRemoveCmd(printStream)

	// Create the wrapper
	cmd := &RemoveCommandStruct{
		cobraCmd: cobraCmd,
		cmdBase:  NewCommandBase(cobraCmd, printStream),
	}

	return cmd
}

// Execute runs the command (for backward compatibility)
func (cmd *RemoveCommandStruct) Execute() {
	if cmd.HelpFlag {
		cmd.cmdBase.PrintUsageInfo(constants.CommandToolConstants.CliHelpFilePrefix + cmd.GetName())
		return
	}

	// Apply flags to the cobra command
	cmd.cobraCmd.Flags().Set("all", fmt.Sprintf("%v", cmd.AllFlag))

	// Handle permissions
	utils.ToolUtil.HandleInstallDirPermission()

	// Handle different scenarios based on flags and arguments
	if cmd.AllFlag {
		// When using --all flag
		if cmd.RemoveCommands != nil && len(cmd.RemoveCommands) > 0 {
			panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
		}

		// Execute removeAll
		executeRemoveAll(cmd.cmdBase.GetPrintStream())
		return
	}

	// When not using --all flag, a distribution argument is required
	if cmd.RemoveCommands == nil || len(cmd.RemoveCommands) == 0 {
		panic(utils.ErrorUtil.CreateDistributionRequiredExceptionWithFlags("remove", "--all, -a"))
	}

	// Check for too many arguments
	if len(cmd.RemoveCommands) > 1 {
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
	}

	// Execute remove for a specific distribution
	executeRemove(cmd.RemoveCommands[0], cmd.cmdBase.GetPrintStream())
}

// GetName returns the name of the command
func (cmd *RemoveCommandStruct) GetName() string {
	return constants.BallerinaCliCommands.REMOVE
}

// PrintLongDesc prints the long description of the command
func (cmd *RemoveCommandStruct) PrintLongDesc(out *strings.Builder) {
	// No implementation needed, Cobra handles this
}

// PrintUsage prints the usage of the command
func (cmd *RemoveCommandStruct) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal dist remove\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *RemoveCommandStruct) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// GetCobraCommand returns the underlying cobra command
func (cmd *RemoveCommandStruct) GetCobraCommand() *cobra.Command {
	return cmd.cobraCmd
}

// GetPrintStream returns the print stream
func (cmd *RemoveCommandStruct) GetPrintStream() io.Writer {
	return cmd.cmdBase.GetPrintStream()
}
