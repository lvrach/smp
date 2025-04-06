package definitions

import (
	"embed"
	"fmt"
	"strings"

	"github.com/lvrach/smp/internal/config"
	"gopkg.in/yaml.v3"
)

//go:embed *.yaml
//go:embed */Dockerfile
var content embed.FS

// MCPRepository provides access to embedded MCP definitions and Dockerfiles
type MCPRepository struct {
	fs embed.FS
}

// NewRepository creates a new MCPRepository instance
func NewRepository() *MCPRepository {
	return &MCPRepository{
		fs: content,
	}
}

// ListMCPs returns a list of all available MCP names
func (r *MCPRepository) ListMCPs() ([]string, error) {
	var mcps []string

	entries, err := r.fs.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded definitions: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			// Remove .yaml extension to get MCP name
			mcpName := strings.TrimSuffix(entry.Name(), ".yaml")
			mcps = append(mcps, mcpName)
		}
	}

	return mcps, nil
}

// MCPConfig returns the configuration for a specific MCP
func (r *MCPRepository) MCPConfig(name string) (*config.MCPConfig, error) {
	// Construct the path to the YAML file
	yamlPath := fmt.Sprintf("%s.yaml", name)

	// Read the YAML file
	data, err := r.fs.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("MCP '%s' not found: %w", name, err)
	}

	// Parse the YAML data
	var mcpConfig config.MCPConfig
	if err := yaml.Unmarshal(data, &mcpConfig); err != nil {
		return nil, fmt.Errorf("failed to parse MCP configuration: %w", err)
	}

	return &mcpConfig, nil
}

// Dockerfile returns the Dockerfile content for a specific MCP
func (r *MCPRepository) Dockerfile(dockerfilePath string) ([]byte, error) {

	// Read the Dockerfile
	data, err := r.fs.ReadFile(dockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("Dockerfile '%s' not found: %w", dockerfilePath, err)
	}

	return data, nil
}
