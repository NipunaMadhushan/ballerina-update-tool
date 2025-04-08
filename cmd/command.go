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
	"ballerina-update-tool/utils"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
)

// CommandBase provides utilities and common functionality for Cobra commands
// This is not a direct conversion of the Command struct, but rather
// a set of utilities that provide similar functionality in the Cobra context
type CommandBase struct {
	CobraCmd    *cobra.Command
	PrintStream io.Writer
}

// NewCommandBase creates a new CommandBase with the given cobra command and print stream
func NewCommandBase(cmd *cobra.Command, printStream io.Writer) *CommandBase {
	if printStream == nil {
		printStream = os.Stdout
	}

	return &CommandBase{
		CobraCmd:    cmd,
		PrintStream: printStream,
	}
}

// GetPrintStream returns the print stream
func (c *CommandBase) GetPrintStream() io.Writer {
	return c.PrintStream
}

// SetPrintStream sets the print stream
func (c *CommandBase) SetPrintStream(printStream io.Writer) {
	c.PrintStream = printStream
	c.CobraCmd.SetOut(printStream)
}

// PrintUsageInfo prints usage information for a command
func (c *CommandBase) PrintUsageInfo(commandName string) {
	// First try to use Cobra's built-in help for the command
	if c.CobraCmd != nil {
		if err := c.CobraCmd.Help(); err == nil {
			return
		}
	}

	// Fall back to the file-based help system
	usageInfo := c.GetCommandUsageInfo(commandName)
	fmt.Fprintln(c.PrintStream, usageInfo)
}

// GetCommandUsageInfo retrieves command usage info from help files
func (c *CommandBase) GetCommandUsageInfo(commandName string) string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	filePath := filepath.Join(filepath.Dir(filepath.Dir(execPath)), "resources", "cli-help", "ballerina-"+commandName+".help")
	content, err := utils.ToolUtil.ReadFileAsString(filePath)
	if err != nil {
		panic(utils.ErrorUtil.CreateUsageExceptionWithHelp("unknown help topic `" + commandName + "`"))
	}
	return content
}

// PrintVersionInfo prints version information
func (c *CommandBase) PrintVersionInfo() {
	currentVersion := utils.ToolUtil.GetCurrentBallerinaVersion()
	toolVersion := utils.ToolUtil.GetCurrentToolsVersion()

	fmt.Fprintf(c.PrintStream, "Update Tool %s\n", toolVersion)
	if currentVersion != "" {
		fmt.Fprintf(c.PrintStream, "Ballerina Distribution Version: %s\n", currentVersion)
	}
}

// Utility functions for Cobra command initialization and usage

// SetupBasicCommand configures a new cobra command with consistent settings
func SetupBasicCommand(use, short, long string, printStream io.Writer) (*cobra.Command, *CommandBase) {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		// PersistentPreRun can be used for common pre-run tasks
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Setup common pre-run tasks here if needed
		},
	}

	// Apply common settings
	SetupCommandCommon(cmd)

	// Set print stream
	cmd.SetOut(printStream)

	// Create the command base
	cmdBase := NewCommandBase(cmd, printStream)

	return cmd, cmdBase
}

// SetupCommandCommon configures common properties for all commands
func SetupCommandCommon(cmd *cobra.Command) {
	// Add common flags if needed
	cmd.Flags().BoolP("help", "h", false, "help for this command")

	// Set custom help template
	cmd.SetHelpTemplate(`{{with .Long}}{{. | trimTrailingWhitespaces}}{{end}}

Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`)
}

// HandlePanic recovers from panics and converts them to errors
func HandlePanic() {
	if r := recover(); r != nil {
		if cmdException, ok := r.(error); ok {
			fmt.Fprintln(os.Stderr, cmdException.Error())
			os.Exit(1)
		} else {
			fmt.Fprintln(os.Stderr, r)
			os.Exit(1)
		}
	}
}

// PrintVersionInfoToStream prints version information to the provided output stream
func PrintVersionInfoToStream(printStream io.Writer) {
	currentVersion := utils.ToolUtil.GetCurrentBallerinaVersion()
	toolVersion := utils.ToolUtil.GetCurrentToolsVersion()

	fmt.Fprintf(printStream, "Update Tool %s\n", toolVersion)
	if currentVersion != "" {
		fmt.Fprintf(printStream, "Ballerina Distribution Version: %s\n", currentVersion)
	}
}

// ExampleCreateCommand Example function showing how to create a command with this utility
func ExampleCreateCommand() *cobra.Command {
	// Create a command with the utility function
	cmd, cmdBase := SetupBasicCommand(
		"example",
		"Example command",
		"This is an example command that demonstrates the CommandBase usage",
		os.Stdout,
	)

	// Add command implementation
	cmd.Run = func(cobraCmd *cobra.Command, args []string) {
		// Use the command base for common functionality
		cmdBase.PrintVersionInfo()

		// Add command-specific behavior
		fmt.Fprintln(cmdBase.GetPrintStream(), "Example command executed!")
	}

	return cmd
}
