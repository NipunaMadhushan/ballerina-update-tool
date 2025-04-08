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
	"github.com/spf13/cobra"
	"io"
)

// BCommand defines the interface for all Ballerina CLI commands
// This interface provides backward compatibility with existing code
// while also supporting the new Cobra-based command structure
type BCommand interface {
	// Execute runs the command
	Execute()

	// GetName returns the name of the command
	GetName() string

	// GetPrintStream returns the print stream
	GetPrintStream() io.Writer

	// GetCobraCommand returns the underlying cobra command
	GetCobraCommand() *cobra.Command
}

// CobraCommandProvider defines a minimal interface for objects that can provide a Cobra command
type CobraCommandProvider interface {
	// GetCobraCommand returns the underlying cobra command
	GetCobraCommand() *cobra.Command
}
