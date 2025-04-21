package docker

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lvrach/smp/internal/config"
	"github.com/lvrach/smp/internal/state"
)

const tagPrefix = "mcp-"

// Runner handles Docker build operations
type Runner struct {
	Config *config.MCPConfig
	State  *state.MCPServer
}

// NewRunner creates a new Docker runner
func NewRunner(mcpConfig *config.MCPConfig, mcpState *state.MCPServer) *Runner {
	return &Runner{
		Config: mcpConfig,
		State:  mcpState,
	}
}

// Run a container from the built image
func (b *Runner) Run() error {
	// Prepare run args
	args := []string{"run", "--rm", "-i"}

	// Add environment variables from state
	for _, envVar := range b.Config.EnvironmentVars {
		if value, exists := b.State.GetEnvironmentVariable(envVar.Name); exists {
			args = append(args, "-e", fmt.Sprintf("%s=%s", envVar.Name, value))
		}
	}

	// Add image name
	args = append(args, b.State.LocalImageTag)

	// Execute docker run command
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// DeleteImage removes the Docker image for this MCP
func DeleteImage(imageName string) error {
	// Execute docker rmi command with force flag to ignore missing images
	cmd := exec.Command("docker", "rmi", "--force", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running docker rmi: %w", err)
	}

	return nil
}
