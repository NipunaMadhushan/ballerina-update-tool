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

// Distribution represents distribution information, corresponding to Distribution.java
type Distribution struct {
	Name       string
	Version    string
	Type       string
	Channel    string
	Dependency string
}

// NewDistribution creates a new Distribution instance
func NewDistribution() *Distribution {
	return &Distribution{
		Name:       "",
		Version:    "",
		Type:       "",
		Channel:    "",
		Dependency: "",
	}
}

// NewDistributionWithVersion creates a new Distribution instance with a specified version
func NewDistributionWithVersion(version string) *Distribution {
	return &Distribution{
		Name:       "",
		Version:    version,
		Type:       "",
		Channel:    "",
		Dependency: "",
	}
}

// NewDistributionWithAllValues creates a new Distribution instance with all values specified
func NewDistributionWithAllValues(name, version, distType, channel, dependency string) *Distribution {
	return &Distribution{
		Name:       name,
		Version:    version,
		Type:       distType,
		Channel:    channel,
		Dependency: dependency,
	}
}

// GetName returns the distribution name
func (d *Distribution) GetName() string {
	return d.Name
}

// SetName sets the distribution name
func (d *Distribution) SetName(name string) {
	d.Name = name
}

// GetVersion returns the distribution version
func (d *Distribution) GetVersion() string {
	return d.Version
}

// SetVersion sets the distribution version
func (d *Distribution) SetVersion(version string) {
	d.Version = version
}

// GetType returns the distribution type
func (d *Distribution) GetType() string {
	return d.Type
}

// SetType sets the distribution type
func (d *Distribution) SetType(distType string) {
	d.Type = distType
}

// GetChannel returns the distribution channel
func (d *Distribution) GetChannel() string {
	return d.Channel
}

// SetChannel sets the distribution channel
func (d *Distribution) SetChannel(channel string) {
	d.Channel = channel
}

// GetDependency returns the distribution dependency
func (d *Distribution) GetDependency() string {
	return d.Dependency
}

// SetDependency sets the distribution dependency
func (d *Distribution) SetDependency(dependency string) {
	d.Dependency = dependency
}
