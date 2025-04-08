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

package main

import (
	"ballerina-update-tool/cmd"
	"ballerina-update-tool/exceptions"
	"ballerina-update-tool/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// Main entry point for the Ballerina CLI tool
func main() {
	// Execute the root command
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newRootCmd creates the root command for the Ballerina CLI
func newRootCmd() *cobra.Command {
	// Get standard output and error streams
	outStream := os.Stdout
	errStream := os.Stderr

	// Create the root command
	rootCmd := &cobra.Command{
		Use:   "bal",
		Short: "Ballerina CLI tool",
		Long: `Ballerina is a programming language for network distributed applications.
This CLI tool helps you manage Ballerina distributions and updates.`,
		// If no subcommand is provided, print help
		Run: func(command *cobra.Command, args []string) {
			command.Help()
		},
		// Handle panics and exceptions
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Add custom pre-run behavior if needed
			return nil
		},
	}

	// Set error output
	rootCmd.SetErr(errStream)
	rootCmd.SetOut(outStream)

	// Add global flags
	rootCmd.PersistentFlags().Bool("help", false, "Help for any command")
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")

	// Override help and version functions
	rootCmd.SetHelpFunc(helpFunction)

	// Check version flag specifically for root command
	versionFlag := rootCmd.Flags().Lookup("version")
	if versionFlag != nil {
		originalRun := rootCmd.Run
		rootCmd.Run = func(c *cobra.Command, args []string) {
			v, _ := c.Flags().GetBool("version")
			if v {
				// Create and execute version command
				versionCmd := cmd.NewVersionCmd(outStream)
				versionCmd.Run(versionCmd, args)
				return
			}
			originalRun(c, args)
		}
	}

	// Add all commands
	addCommands(rootCmd, outStream, errStream)

	return rootCmd
}

// addCommands adds all subcommands to the root command
func addCommands(rootCmd *cobra.Command, outStream, errStream *os.File) {
	// Add direct commands
	rootCmd.AddCommand(cmd.NewVersionCmd(outStream))
	rootCmd.AddCommand(cmd.NewBuildCmd(outStream))
	rootCmd.AddCommand(cmd.NewUpdateToolCmd(outStream))
	rootCmd.AddCommand(cmd.NewHelpCmd(outStream))

	// Create and add dist command
	distCmd := cmd.NewDistCmd(outStream)

	// Add subcommands to dist command
	distCmd.AddCommand(cmd.NewListCmd(outStream))
	distCmd.AddCommand(cmd.NewPullCmd(outStream))
	distCmd.AddCommand(cmd.NewRemoveCmd(outStream))
	distCmd.AddCommand(cmd.NewUpdateCmd(outStream))
	distCmd.AddCommand(cmd.NewUseCmd(outStream))

	// Add dist command to root
	rootCmd.AddCommand(distCmd)
}

// helpFunction is a custom help function that wraps cobra's built-in help
// and handles any panics or exceptions
func helpFunction(command *cobra.Command, args []string) {
	defer func() {
		if r := recover(); r != nil {
			if cmdException, ok := r.(exceptions.CommandException); ok {
				utils.ErrorUtil.PrintLauncherException(cmdException, os.Stderr)
			} else {
				fmt.Fprintln(os.Stderr, r)
			}
		}
	}()

	// If we're showing help for the root command and no arguments are provided,
	// show help for the help command itself
	if command.Name() == "bal" && len(args) == 0 {
		command.Help()
		return
	}

	// If arguments are provided, try to show help for the specified command
	if len(args) > 0 {
		// Try to find the command
		subCmd, _, err := command.Root().Find(args)
		if err == nil && subCmd != command.Root() {
			subCmd.Help()
			return
		}

		// If not found, use your existing help file system
		helpCmd := cmd.NewHelpCmd(os.Stdout)
		helpCmd.Run(helpCmd, args)
		return
	}

	// Default to built-in help
	command.Root().Help()
}

// RootCommand For backward compatibility with existing code
type RootCommand struct {
	cobraCmd *cobra.Command
}

// Execute runs the root command
func (r *RootCommand) Execute() error {
	return r.cobraCmd.Execute()
}

// GetCobraCommand returns the underlying cobra command
func (r *RootCommand) GetCobraCommand() *cobra.Command {
	return r.cobraCmd
}
