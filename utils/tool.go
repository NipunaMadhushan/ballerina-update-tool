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

// Tool represents tool version information, corresponding to Tool.java
type Tool struct {
	Version       string
	Compatibility string
}

// NewTool creates a new Tool instance
func NewTool() *Tool {
	return &Tool{
		Version:       "",
		Compatibility: "",
	}
}

// NewToolWithValues creates a new Tool instance with specified values
func NewToolWithValues(version, compatibility string) *Tool {
	return &Tool{
		Version:       version,
		Compatibility: compatibility,
	}
}

// GetVersion returns the tool version
func (t *Tool) GetVersion() string {
	return t.Version
}

// SetVersion sets the tool version
func (t *Tool) SetVersion(version string) {
	t.Version = version
}

// GetCompatibility returns the compatibility information
func (t *Tool) GetCompatibility() string {
	return t.Compatibility
}

// SetCompatibility sets the compatibility information
func (t *Tool) SetCompatibility(compatibility string) {
	t.Compatibility = compatibility
}
