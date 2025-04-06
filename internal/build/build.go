package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lvrach/smp/definitions"
	"github.com/lvrach/smp/internal/config"
	"github.com/lvrach/smp/internal/git"
	"github.com/lvrach/smp/internal/state"
)

const tagPrefix = "mcp-"

// Builder handles the build process for MCPs
type Builder struct {
	Config  *config.MCPConfig
	TempDir string
	State   *state.MCPServer
}

// NewBuilder creates a new builder for the given MCP config
func NewBuilder(mcpConfig *config.MCPConfig) (*Builder, error) {
	// Create a temporary directory for the build
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("mcp-%s-", mcpConfig.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Initialize state manager
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	stateManager, err := state.NewStore(filepath.Join(homeDir, ".mcp"))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize state manager: %w", err)
	}

	// Load or create MCP state
	mcpState, err := stateManager.Load(mcpConfig.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to load MCP state: %w", err)
	}

	return &Builder{
		Config:  mcpConfig,
		TempDir: tempDir,
		State:   mcpState,
	}, nil
}

// CleanUp removes the temporary directory
func (b *Builder) CleanUp() error {
	return os.RemoveAll(b.TempDir)
}

// BuildFromRepo clones a repository and builds a Docker image
func (b *Builder) BuildFromRepo() (string, error) {
	// Clone the repository to a temporary directory
	fmt.Printf("Cloning repository %s...\n", b.Config.Repository)
	tempDir, err := git.CloneRepository(b.Config.Repository, b.Config.Branch)
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}
	// Clean up the temporary directory when done
	defer os.RemoveAll(tempDir)

	// Update the builder's path to use the cloned repository
	b.TempDir = tempDir

	// Build the Docker image
	return b.buildImage()
}

// BuildImage builds the Docker image using the Docker package
func (b *Builder) buildImage() (string, error) {
	// Get embedded repository
	repo := definitions.NewRepository()
	tag := fmt.Sprintf("mcp-%s:latest", b.Config.Name)

	// Try to get Dockerfile from embedded definitions first
	dockerfileContent, err := repo.Dockerfile(b.Config.Dockerfile)
	if err != nil {
		return tag, fmt.Errorf("failed to get Dockerfile from embedded definitions: %w", err)
	}

	// Write Dockerfile to temporary directory
	outputPath := filepath.Join(b.TempDir, "Dockerfile")
	if err := os.WriteFile(outputPath, dockerfileContent, 0644); err != nil {
		return tag, fmt.Errorf("failed to write Dockerfile to temp dir: %w", err)
	}

	fmt.Println("overriding Dockerfile: ", outputPath)

	// Build the image
	fmt.Printf("Building Docker image for MCP '%s'...\n", b.Config.Name)

	// Prepare build args
	args := []string{"build", "-t", tagPrefix + b.Config.Name}

	// Add Dockerfile path and context
	args = append(args, b.TempDir)

	// Execute docker build command
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return tag, fmt.Errorf("failed to build Docker image: %w", err)
	}

	// Update state with the new image tag
	b.State.SetLocalImageTag(fmt.Sprintf("mcp-%s:latest", b.Config.Name))

	// Save the updated state
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return tag, fmt.Errorf("failed to get user home directory: %w", err)
	}

	stateManager, err := state.NewStore(filepath.Join(homeDir, ".mcp"))
	if err != nil {
		return tag, fmt.Errorf("failed to initialize state manager: %w", err)
	}

	if err := stateManager.Save(b.State); err != nil {
		return tag, fmt.Errorf("failed to save MCP state: %w", err)
	}

	fmt.Printf("Docker image '%s' built successfully\n", b.Config.Name)
	return tag, nil
}
