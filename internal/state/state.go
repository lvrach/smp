package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MCPServer represents the state of an installed MCP
type MCPServer struct {
	// Name of the MCP
	Name string `json:"name"`

	// Docker image tag for this MCP
	LocalImageTag string `json:"local_image_tag"`

	// Environment variables set for this MCP
	EnvironmentVariables map[string]string `json:"environment_variables"`

	// Environment variables that are stored in the keychain
	KeyChainEnvVars map[string]string `json:"keychain_env_vars"`

	// Hosts that are configured to run this MCP
	ConfiguredHosts []string `json:"configured_hosts"`
}

// Store handles the persistence and retrieval of MCP states
type Store struct {
	stateDir string
}

// NewStore operates a new state store
func NewStore(baseDir string) (*Store, error) {
	stateDir := filepath.Join(baseDir, "state")

	// Create state directory if it doesn't exist
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	return &Store{
		stateDir: stateDir,
	}, nil
}

// NewHomeStore operates a state store in the user's home directory
func NewHomeStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	return NewStore(filepath.Join(homeDir, ".smp"))
}

// Save the state of an MCP
func (sm *Store) Save(state *MCPServer) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal MCP state: %w", err)
	}

	filename := filepath.Join(sm.stateDir, state.Name+".json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// Load the state of an MCP
func (sm *Store) Load(mcpName string) (*MCPServer, error) {
	filename := filepath.Join(sm.stateDir, mcpName+".json")

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty state if file doesn't exist
			return &MCPServer{
				Name:                 mcpName,
				EnvironmentVariables: make(map[string]string),
			}, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state MCPServer
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal MCP state: %w", err)
	}

	return &state, nil
}

// Delete removes the state file for an MCP
func (sm *Store) Delete(mcpName string) error {
	filename := filepath.Join(sm.stateDir, mcpName+".json")

	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete state file: %w", err)
	}

	return nil
}

// List returns a list of all stored MCP states
func (sm *Store) List() ([]*MCPServer, error) {
	files, err := os.ReadDir(sm.stateDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read state directory: %w", err)
	}

	var states []*MCPServer
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		mcpName := file.Name()[:len(file.Name())-5] // Remove .json extension
		state, err := sm.Load(mcpName)
		if err != nil {
			return nil, err
		}

		states = append(states, state)
	}

	return states, nil
}

// SetEnvironmentVariable sets an environment variable in the MCP state
func (s *MCPServer) SetEnvironmentVariable(key, value string) {
	if s.EnvironmentVariables == nil {
		s.EnvironmentVariables = make(map[string]string)
	}
	s.EnvironmentVariables[key] = value
}

// GetEnvironmentVariable gets an environment variable from the MCP state
func (s *MCPServer) GetEnvironmentVariable(key string) (string, bool) {
	if s.EnvironmentVariables == nil {
		return "", false
	}
	value, exists := s.EnvironmentVariables[key]
	return value, exists
}

// SetLocalImageTag sets the local Docker image tag for this MCP
func (s *MCPServer) SetLocalImageTag(tag string) {
	s.LocalImageTag = tag
}
