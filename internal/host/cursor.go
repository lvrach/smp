package host

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type cursor struct{}

type cursorConfig struct {
	MCPServers map[string]struct {
		Command string            `json:"command"`
		Args    []string          `json:"args"`
		Envs    map[string]string `json:"envs"`
	} `json:"mcpServers"`
}

func (c *cursor) Name() string {
	return "cursor"
}

func (c *cursor) getConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".cursor", "mcp.json")
}

func (c *cursor) Available() bool {
	// Cursor is always available since we can create the config file
	return true
}

func (c *cursor) loadConfig() (*cursorConfig, error) {
	configPath := c.getConfigPath()

	// Return empty config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &cursorConfig{
			MCPServers: make(map[string]struct {
				Command string            `json:"command"`
				Args    []string          `json:"args"`
				Envs    map[string]string `json:"envs"`
			}),
		}, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config cursorConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Initialize MCPServers map if it's nil
	if config.MCPServers == nil {
		config.MCPServers = make(map[string]struct {
			Command string            `json:"command"`
			Args    []string          `json:"args"`
			Envs    map[string]string `json:"envs"`
		})
	}

	return &config, nil
}

func (c *cursor) saveConfig(config *cursorConfig) error {
	configPath := c.getConfigPath()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return os.WriteFile(configPath, data, 0644)
}

func (c *cursor) Connect(binaryPath string, server string) error {
	config, err := c.loadConfig()
	if err != nil {
		return err
	}

	// Check if server already exists
	if existing, exists := config.MCPServers[server]; exists {
		// If it exists and has different binary path, return error
		if existing.Command != binaryPath {
			return fmt.Errorf("server %s already exists with different binary path: %s", server, existing.Command)
		}
		// If it exists with same binary path, no need to update
		return nil
	}

	// Add new server configuration
	config.MCPServers[server] = struct {
		Command string            `json:"command"`
		Args    []string          `json:"args"`
		Envs    map[string]string `json:"envs"`
	}{
		Command: binaryPath,
		Args:    []string{"run", server},
		Envs: map[string]string{
			"HOME": os.Getenv("HOME"),
		},
	}

	// Save updated config
	if err := c.saveConfig(config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func (c *cursor) Disconnect(server string) (bool, error) {
	config, err := c.loadConfig()
	if err != nil {
		return false, err
	}

	if _, exists := config.MCPServers[server]; !exists {
		// Server wasn't configured, not an error
		return false, nil
	}

	delete(config.MCPServers, server)

	// Save updated config
	if err := c.saveConfig(config); err != nil {
		return true, fmt.Errorf("saving config: %w", err)
	}

	return true, nil
}
