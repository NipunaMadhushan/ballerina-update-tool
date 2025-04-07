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

package main

import (
	"ballerina-update-tool/cmd"
	"ballerina-update-tool/constants"
	"ballerina-update-tool/exceptions"
	"ballerina-update-tool/utils"
	"fmt"
	"os"
)

// Main entry point for the Ballerina CLI tool
func main() {
	outStream := os.Stdout
	errStream := os.Stderr

	// Set up the command parser and execute the appropriate command
	exitCode := executeCommand(os.Args[1:], outStream, errStream)
	os.Exit(exitCode)
}

// executeCommand parses the command line arguments and executes the appropriate command
func executeCommand(args []string, outStream, errStream *os.File) int {
	defer func() {
		if r := recover(); r != nil {
			if cmdException, ok := r.(exceptions.CommandException); ok {
				utils.ErrorUtil.PrintLauncherException(cmdException, errStream)
			} else {
				fmt.Fprintln(errStream, r)
			}
		}
	}()

	invokedCmd := getInvokedCmd(args, outStream, errStream)
	if invokedCmd != nil {
		invokedCmd.Execute()
	}

	return 0
}

// getInvokedCmd parses the command line arguments and returns the corresponding command
func getInvokedCmd(args []string, outStream, errStream *os.File) cmd.BCommand {
	// Create default command
	defaultCmd := cmd.NewDefaultCommand(outStream)

	// Define command parser (in Go, we'll use a simple map-based approach)
	cmdParser := make(map[string]cmd.BCommand)

	// Add help command
	helpCmd := cmd.NewHelp(outStream)
	cmdParser[constants.BallerinaCliCommands.HELP] = helpCmd

	// Distribution command and subcommands
	distCmd := cmd.NewDistribution(outStream)
	distCmdParser := make(map[string]cmd.BCommand)

	listCmd := cmd.NewList(outStream)
	distCmdParser[constants.BallerinaCliCommands.LIST] = listCmd

	pullCmd := cmd.NewPull(outStream)
	distCmdParser[constants.BallerinaCliCommands.PULL] = pullCmd

	removeCmd := cmd.NewRemove(outStream)
	distCmdParser[constants.BallerinaCliCommands.REMOVE] = removeCmd

	updateCmd := cmd.NewUpdate(outStream)
	distCmdParser[constants.BallerinaCliCommands.UPDATE] = updateCmd

	useCmd := cmd.NewUse(outStream)
	distCmdParser[constants.BallerinaCliCommands.USE] = useCmd

	// Add dist command and its subcommands
	cmdParser[constants.BallerinaCliCommands.DIST] = distCmd

	// Add update tool command
	updateToolCmd := cmd.NewUpdateTool(outStream)
	cmdParser[constants.BallerinaCliCommands.UPDATE] = updateToolCmd

	// Add build command
	buildCmd := cmd.NewBuild(outStream)
	cmdParser[constants.BallerinaCliCommands.BUILD] = buildCmd

	// Add version command
	versionCmd := cmd.NewVersion(outStream)
	cmdParser[constants.BallerinaCliCommands.VERSION] = versionCmd

	// Parse arguments and return the appropriate command
	if len(args) == 0 {
		return defaultCmd
	}

	// Check for top-level commands
	if cmd, exists := cmdParser[args[0]]; exists {
		// Handle subcommands for dist
		if args[0] == constants.BallerinaCliCommands.DIST && len(args) > 1 {
			if subCmd, exists := distCmdParser[args[1]]; exists {
				// Set command arguments and flags
				setCommandArgs(subCmd, args[2:])
				return subCmd
			}
		}
		// Set command arguments and flags
		setCommandArgs(cmd, args[1:])
		return cmd
	}
	if containsFlag(args, "--version") || containsFlag(args, "-v") {
		setCommandArgs(versionCmd, args)
		return versionCmd
	}

	// If no command matches, return default
	return defaultCmd
}

// setCommandArgs sets the arguments and flags for a command
// This is a simplified approach; in a real implementation, you would
// use a proper flag parsing library like "flag" or "github.com/spf13/pflag"
func setCommandArgs(command cmd.BCommand, args []string) {
	// In this simplified implementation, we're just setting the
	// command's arguments directly based on the command type
	switch c := command.(type) {
	case *cmd.DefaultCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		if len(args) > 0 {
			c.HelpCommands = args
		}
	case *cmd.HelpCommand:
		if len(args) > 0 {
			c.HelpCommands = args
		}
	case *cmd.DistributionCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		if len(args) > 0 {
			c.DistCommands = args
		}
	case *cmd.ListCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		if containsFlag(args, "--all") || containsFlag(args, "-a") {
			c.AllFlag = true
		}
		if containsFlag(args, "--pre-releases") || containsFlag(args, "-p") {
			c.PreReleasesFlag = true
		}
		// Filter out flags to get commands
		c.ListCommands = getCommandsWithoutFlags(args)
	case *cmd.PullCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		if containsFlag(args, "--test") || containsFlag(args, "-t") {
			c.TestFlag = true
		}
		// Filter out flags to get commands
		c.PullCommands = getCommandsWithoutFlags(args)
	case *cmd.RemoveCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		if containsFlag(args, "--all") || containsFlag(args, "-a") {
			c.AllFlag = true
		}
		// Filter out flags to get commands
		c.RemoveCommands = getCommandsWithoutFlags(args)
	case *cmd.UpdateCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		if containsFlag(args, "--test") || containsFlag(args, "-t") {
			c.TestFlag = true
		}
		// Filter out flags to get commands
		c.UpdateCommands = getCommandsWithoutFlags(args)
	case *cmd.UpdateToolCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		// Filter out flags to get commands
		c.UpdateCommands = getCommandsWithoutFlags(args)
	case *cmd.UseCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		// Filter out flags to get commands
		c.UseCommands = getCommandsWithoutFlags(args)
	case *cmd.BuildCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		// Filter out flags to get commands
		c.BuildCommands = getCommandsWithoutFlags(args)
	case *cmd.VersionCommand:
		if containsHelpFlag(args) {
			c.HelpFlag = true
		}
		// Filter out flags to get commands
		c.VersionCommands = getCommandsWithoutFlags(args)
	}
}

// containsHelpFlag checks if the arguments contain a help flag
func containsHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" || arg == "?" {
			return true
		}
	}
	return false
}

// containsFlag checks if the arguments contain a specific flag
func containsFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}

// getCommandsWithoutFlags returns the arguments without any flags
func getCommandsWithoutFlags(args []string) []string {
	var result []string
	for _, arg := range args {
		// Skip arguments that start with - or --
		if len(arg) > 0 && arg[0] != '-' {
			result = append(result, arg)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
