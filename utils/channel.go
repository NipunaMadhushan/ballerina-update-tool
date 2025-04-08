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

// Channel represents a distribution channel, corresponding to Channel.java
type Channel struct {
	Name          string
	Distributions []Distribution
}

// NewChannel creates a new Channel instance
func NewChannel() *Channel {
	return &Channel{
		Name:          "",
		Distributions: []Distribution{},
	}
}

// NewChannelWithName creates a new Channel instance with a specified name
func NewChannelWithName(name string) *Channel {
	return &Channel{
		Name:          name,
		Distributions: []Distribution{},
	}
}

// NewChannelWithNameAndDistributions creates a new Channel instance with specified name and distributions
func NewChannelWithNameAndDistributions(name string, distributions []Distribution) *Channel {
	return &Channel{
		Name:          name,
		Distributions: distributions,
	}
}

// GetName returns the channel name
func (c *Channel) GetName() string {
	return c.Name
}

// SetName sets the channel name
func (c *Channel) SetName(name string) {
	c.Name = name
}

// GetDistributions returns the distributions in this channel
func (c *Channel) GetDistributions() []Distribution {
	return c.Distributions
}

// SetDistributions sets the distributions in this channel
func (c *Channel) SetDistributions(distributions []Distribution) {
	c.Distributions = distributions
}

// AddDistribution adds a distribution to the front of the channel
func (c *Channel) AddDistribution(distribution Distribution) {
	c.Distributions = append([]Distribution{distribution}, c.Distributions...)
}
