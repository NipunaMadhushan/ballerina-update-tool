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

package utils

import (
	"ballerina-update-tool/constants"
	"ballerina-update-tool/exceptions"
	"fmt"
	"io"
)

// ErrorUtilStruct provides utility functions for error handling
type ErrorUtilStruct struct{}

// ErrorUtil is the singleton instance of ErrorUtilStruct
var ErrorUtil = ErrorUtilStruct{}

// CreateCommandException creates a command exception with the given error message
func (e ErrorUtilStruct) CreateCommandException(errorMsg string) *exceptions.CommandException {
	exception := exceptions.New()
	exception.AddMessage("ballerina: " + errorMsg)
	return exception
}

// CreateUsageExceptionWithHelp creates a command exception with help information
func (e ErrorUtilStruct) CreateUsageExceptionWithHelp(errorMsg string) *exceptions.CommandException {
	exception := exceptions.New()
	exception.AddMessage("ballerina: " + errorMsg)
	exception.AddMessage("Run 'bal help' for usage.")
	return exception
}

// CreateUsageExceptionWithHelpSubCommand creates a command exception with help for a specific subcommand
func (e ErrorUtilStruct) CreateUsageExceptionWithHelpSubCommand(errorMsg, subCommand string) *exceptions.CommandException {
	exception := exceptions.New()
	exception.AddMessage("ballerina: " + errorMsg)
	exception.AddMessage(fmt.Sprintf("Run 'bal help %s' for usage.", subCommand))
	return exception
}

// CreateDistSubCommandUsageExceptionWithHelp creates a command exception for distribution-related commands
func (e ErrorUtilStruct) CreateDistSubCommandUsageExceptionWithHelp(errorMsg, subCommand string) *exceptions.CommandException {
	return e.CreateUsageExceptionWithHelpSubCommand(errorMsg, constants.BallerinaCliCommands.DIST+" "+subCommand)
}

// CreateDistributionNotFoundException creates an exception for when a distribution is not found
func (e ErrorUtilStruct) CreateDistributionNotFoundException(distribution string) *exceptions.CommandException {
	return e.CreateCommandException(fmt.Sprintf("distribution '%s' not found", distribution))
}

// CreateDependencyNotFoundException creates an exception for when a dependency is not found
func (e ErrorUtilStruct) CreateDependencyNotFoundException(dependency string) *exceptions.CommandException {
	return e.CreateCommandException(fmt.Sprintf("dependency '%s' not found", dependency))
}

// CreateDistributionRequiredException creates an exception for operations requiring a distribution
func (e ErrorUtilStruct) CreateDistributionRequiredException(operation string) *exceptions.CommandException {
	return e.CreateDistSubCommandUsageExceptionWithHelp("a distribution must be specified to "+operation,
		operation)
}

// CreateDistributionRequiredExceptionWithFlags creates an exception for operations requiring a distribution or specific flags
func (e ErrorUtilStruct) CreateDistributionRequiredExceptionWithFlags(operation, flags string) *exceptions.CommandException {
	return e.CreateDistSubCommandUsageExceptionWithHelp(
		fmt.Sprintf("a distribution or `%s` must be specified to %s", flags, operation), operation)
}

// PrintLauncherException prints the exception messages to the given output stream
func (e ErrorUtilStruct) PrintLauncherException(exception exceptions.CommandException, outStream io.Writer) {
	for _, message := range exception.GetMessages() {
		fmt.Fprintln(outStream, message)
	}
}
