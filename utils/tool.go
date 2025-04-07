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
