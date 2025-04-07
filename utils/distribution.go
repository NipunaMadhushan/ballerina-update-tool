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
