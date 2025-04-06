package commands

import (
	"fmt"

	"github.com/lvrach/smp/definitions"
	"github.com/lvrach/smp/internal/build"
	"github.com/lvrach/smp/internal/docker"
	"github.com/lvrach/smp/internal/prompt"
	"github.com/lvrach/smp/internal/state"
	"github.com/urfave/cli/v2"
)

// BuildCommand returns the command for building an MCP
func BuildCommand() *cli.Command {
	return &cli.Command{
		Name:      "build",
		Usage:     "Build a Docker image for an MCP",
		ArgsUsage: "[name]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "test-run",
				Usage: "Test run the container after building",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("missing required argument: name")
			}

			name := c.Args().Get(0)

			repo := definitions.NewRepository()
			mcpConfig, err := repo.MCPConfig(name)
			if err != nil {
				return fmt.Errorf("failed to get MCP configuration: %w", err)
			}

			// Create a new builder with a temporary directory
			builder, err := build.NewBuilder(mcpConfig)
			if err != nil {
				return fmt.Errorf("failed to create builder: %w", err)
			}
			defer builder.CleanUp() // Clean up temporary directory when done

			// Build the image from the repository
			tag, err := builder.BuildFromRepo()
			if err != nil {
				return err
			}

			// Run the container if requested
			if c.Bool("test-run") {
				// Create a new state for the test run
				mcpState := &state.MCPServer{
					Name:                 name,
					LocalImageTag:        tag,
					EnvironmentVariables: make(map[string]string),
				}

				// Prompt for environment variables
				if err := prompt.PromptEnvironmentVariables(mcpConfig, mcpState); err != nil {
					return fmt.Errorf("failed to get environment variables: %w", err)
				}

				// Create runner and run the container
				runner := docker.NewRunner(mcpConfig, mcpState)
				if err := runner.Run(); err != nil {
					return fmt.Errorf("failed to run container: %w", err)
				}
			}

			return nil
		},
	}
}
