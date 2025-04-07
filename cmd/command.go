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
	"ballerina-update-tool/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Command struct represents the base command with generic methods
type Command struct {
	printStream io.Writer
}

// New creates a new Command with stdout as the default print stream
func New() *Command {
	return &Command{
		printStream: os.Stdout,
	}
}

// NewWithWriter creates a new Command with the specified print stream
func NewWithWriter(printStream io.Writer) *Command {
	return &Command{
		printStream: printStream,
	}
}

// GetPrintStream returns the print stream
func (c *Command) GetPrintStream() io.Writer {
	return c.printStream
}

// SetPrintStream sets the print stream
func (c *Command) SetPrintStream(printStream io.Writer) {
	c.printStream = printStream
}

// PrintUsageInfo prints usage information for a command
func (c *Command) PrintUsageInfo(commandName string) {
	usageInfo := c.GetCommandUsageInfo(commandName)
	fmt.Fprintln(c.printStream, usageInfo)
}

// GetCommandUsageInfo retrieves command usage info from help files
func (c *Command) GetCommandUsageInfo(commandName string) string {
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
func (c *Command) PrintVersionInfo() {
	utils.ToolUtil.GetCurrentBallerinaVersion()
	output := "Update Tool " + utils.ToolUtil.GetCurrentToolsVersion() + "\n"
	fmt.Fprint(c.printStream, output)
}
