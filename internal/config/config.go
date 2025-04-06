package config

// MCPConfig represents the configuration for a Multi-Container Platform
type MCPConfig struct {
	Name            string                `yaml:"name"`
	Repository      string                `yaml:"repository"`
	Branch          string                `yaml:"branch,omitempty"`
	Dockerfile      string                `yaml:"dockerfile,omitempty"`
	EnvironmentVars []EnvironmentVariable `yaml:"environment,omitempty"`
}

type EnvironmentVariable struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
}
