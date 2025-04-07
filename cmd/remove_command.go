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
	"os"
	"path/filepath"
	"strings"
)

// RemoveCommand represents the "remove" command and holds arguments and flags specified by the user
// Command name: "remove", description: "Remove Ballerina distribution"
type RemoveCommand struct {
	*Command
	RemoveCommands  []string // Command name
	HelpFlag        bool     // --help, -h, ?
	AllFlag         bool     // --all, -a
	ParentCmdParser interface{}
}

// NewRemove creates a new RemoveCommand
func NewRemove(printStream io.Writer) *RemoveCommand {
	cmd := &RemoveCommand{}
	cmd.Command = NewWithWriter(printStream)
	return cmd
}

// Execute runs the command
func (cmd *RemoveCommand) Execute() {
	if cmd.HelpFlag {
		cmd.PrintUsageInfo(constants.CommandToolConstants.CliHelpFilePrefix + cmd.GetName())
		return
	}

	if cmd.AllFlag {
		if cmd.RemoveCommands == nil {
			utils.ToolUtil.HandleInstallDirPermission()
			cmd.removeAll()
			return
		}
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments", cmd.GetName()))
	}

	if cmd.RemoveCommands == nil || len(cmd.RemoveCommands) == 0 {
		panic(utils.ErrorUtil.CreateDistributionRequiredExceptionWithFlags("remove", "--all, -a"))
	}

	if len(cmd.RemoveCommands) > 1 {
		panic(utils.ErrorUtil.CreateDistSubCommandUsageExceptionWithHelp("too many arguments",
			constants.BallerinaCliCommands.REMOVE))
	}

	utils.ToolUtil.HandleInstallDirPermission()
	cmd.remove(cmd.RemoveCommands[0])
}

// GetName returns the name of the command
func (cmd *RemoveCommand) GetName() string {
	return constants.BallerinaCliCommands.REMOVE
}

// PrintLongDesc prints the long description of the command
func (cmd *RemoveCommand) PrintLongDesc(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// PrintUsage prints the usage of the command
func (cmd *RemoveCommand) PrintUsage(out *strings.Builder) {
	out.WriteString("  bal dist remove\n")
}

// SetParentCmdParser sets the parent command parser
func (cmd *RemoveCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}

// isCurrentVersion checks if the given version is the current active version
func (cmd *RemoveCommand) isCurrentVersion(version string) bool {
	return version == utils.ToolUtil.GetCurrentBallerinaVersion()
}

// remove removes a specific distribution version
func (cmd *RemoveCommand) remove(version string) {
	if cmd.isCurrentVersion(version) {
		panic(utils.ErrorUtil.CreateCommandException("The active Ballerina distribution cannot be removed"))
	} else {
		file := utils.ToolUtil.GetType(version) + "-" + version
		directory := filepath.Join(utils.ToolUtil.GetDistributionsPath(), file)

		if _, err := os.Stat(directory); !os.IsNotExist(err) {
			err := utils.OSUtils.DeleteFiles(directory)
			if err != nil {
				panic(utils.ErrorUtil.CreateCommandException("error occurred while removing '" + version + "': " + err.Error()))
			}

			err = utils.OSUtils.DeleteCaches(version, cmd.GetPrintStream())
			if err != nil {
				panic(utils.ErrorUtil.CreateCommandException("error occurred while removing caches for '" + version + "': " + err.Error()))
			}

			fmt.Fprintf(cmd.GetPrintStream(), "Distribution '%s' successfully removed\n", version)
			utils.ToolUtil.RemoveUnusedDependencies(version, cmd.GetPrintStream())
		} else {
			panic(utils.ErrorUtil.CreateCommandException("distribution '" + version + "' not found"))
		}
	}
}

// removeAll removes all distributions except the active one
func (cmd *RemoveCommand) removeAll() {
	folder := utils.ToolUtil.GetDistributionsPath()
	listOfFiles, err := os.ReadDir(folder)
	if err != nil {
		panic(utils.ErrorUtil.CreateCommandException("error occurred while reading distributions: " + err.Error()))
	}

	// checking for 2 files for zip pack and 3 files for installers
	if len(listOfFiles) == 2 || (len(listOfFiles) == 3 &&
		fileExists(filepath.Join(folder, "installer-version"))) {
		fmt.Fprintln(cmd.GetPrintStream(), "There is nothing to remove. Only active distribution is remaining")
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
			if fileExists(directory) && (!cmd.isCurrentVersion(version) || version == "") {
				err := utils.OSUtils.DeleteFiles(directory)
				if err != nil {
					panic(utils.ErrorUtil.CreateCommandException("error occurred while removing distribution '" + fileName + "': " + err.Error()))
				}
			}
		}
	}

	fmt.Fprintln(cmd.GetPrintStream(), "All non-active distributions are successfully removed")

	activeDistribution := utils.ToolUtil.GetCurrentBallerinaVersion()
	dependencyForActiveDistribution := utils.ToolUtil.GetDependency(
		cmd.GetPrintStream(),
		activeDistribution,
		utils.ToolUtil.GetType(activeDistribution),
		activeDistribution)

	dependencies, err := os.ReadDir(utils.ToolUtil.GetDependencyPath())
	if err != nil {
		panic(utils.ErrorUtil.CreateCommandException("error occurred while reading dependencies: " + err.Error()))
	}

	if len(dependencies) > 1 {
		fmt.Fprintln(cmd.GetPrintStream(), "Removing unused dependencies")
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

// fileExists checks if a file or directory exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
