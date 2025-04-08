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

package exceptions

// CommandException represents an exception that occurred in Ballerina program launcher.
// It implements the error interface in Go and provides additional functionality for
// detailed error messages.
type CommandException struct {
	detailedMessages []string
}

// New creates a new CommandException
func New() *CommandException {
	return &CommandException{
		detailedMessages: []string{},
	}
}

// Error returns the error message (required to implement the error interface)
func (e *CommandException) Error() string {
	if len(e.detailedMessages) > 0 {
		return e.detailedMessages[0]
	}
	return "command exception occurred"
}

// GetDetailedMessages returns all detailed messages
func (e *CommandException) GetDetailedMessages() []string {
	return e.detailedMessages
}

// AddMessage adds a message to the detailed messages
func (e *CommandException) AddMessage(message string) {
	e.detailedMessages = append(e.detailedMessages, message)
}

// GetMessages returns all messages (alias for GetDetailedMessages for compatibility)
func (e *CommandException) GetMessages() []string {
	return e.detailedMessages
}
