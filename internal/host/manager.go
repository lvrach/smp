package host

import (
	"fmt"
	"os"
	"path/filepath"
)

type host interface {
	Name() string
	Available() bool
	Connect(binaryPath string, server string) error
	Disconnect(server string) (bool, error)
}

var hosts = map[string]host{
	"claude_desktop": &claudeDesktop{},
	"cline":          &cline{},
	"cursor":         &cursor{},
}

type Manager struct {
	BinaryPath string
	Hosts      map[string]host
}

func DefaultManager() *Manager {
	execPath, err := os.Executable()
	if err != nil {
		// If we can't get the path, return empty string
		execPath = ""
	}
	// Get the absolute path
	absPath, err := filepath.EvalSymlinks(execPath)
	if err == nil {
		execPath = absPath
	}

	return &Manager{
		BinaryPath: execPath,
		Hosts:      hosts,
	}
}

// List available MCP hosts installed in the system
func (m *Manager) List() ([]string, error) {
	hosts := []string{}
	for _, h := range m.Hosts {
		if h.Available() {
			hosts = append(hosts, h.Name())
		}
	}
	return hosts, nil
}

// Connect an MCP server to a MCP host
func (m *Manager) Connect(host string, server string) error {
	if h, exists := m.Hosts[host]; exists {
		return h.Connect(m.BinaryPath, server)
	}
	return fmt.Errorf("host %s not found", host)
}

// Disconnect an MCP server from a MCP host
// Returns true if the server was configured in the host, false if it wasn't
func (m *Manager) Disconnect(host string, server string) (bool, error) {
	if h, exists := m.Hosts[host]; exists {
		return h.Disconnect(server)
	}
	return false, fmt.Errorf("host %s not found", host)
}
