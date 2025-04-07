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
