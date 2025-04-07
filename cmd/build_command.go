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
	"io"
	"strings"
)

// BuildCommand represents the "build" command and is used to notify latest distribution information
type BuildCommand struct {
	*Command
	BuildCommands   []string // Command name
	HelpFlag        bool     // --help, -h, ?
	ParentCmdParser interface{}
}

// NewBuild creates a new BuildCommand
func NewBuild(printStream io.Writer) *BuildCommand {
	return &BuildCommand{
		Command: NewWithWriter(printStream),
	}
}

// Execute runs the command
func (cmd *BuildCommand) Execute() {
	if cmd.HelpFlag {
		return
	}
	utils.ToolUtil.CheckForUpdate(cmd.GetPrintStream())
}

// GetName returns the name of the command
func (cmd *BuildCommand) GetName() string {
	return constants.BallerinaCliCommands.BUILD
}

// PrintLongDesc prints the long description of the command
func (cmd *BuildCommand) PrintLongDesc(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// PrintUsage prints the usage of the command
func (cmd *BuildCommand) PrintUsage(out *strings.Builder) {
	// Implementation is empty in the original Java code
}

// SetParentCmdParser sets the parent command parser
func (cmd *BuildCommand) SetParentCmdParser(parentCmdParser interface{}) {
	cmd.ParentCmdParser = parentCmdParser
}
