package commands

import (
	"fmt"

	"github.com/lvrach/smp/definitions"
	"github.com/lvrach/smp/internal/build"
	"github.com/lvrach/smp/internal/host"
	"github.com/lvrach/smp/internal/prompt"
	"github.com/lvrach/smp/internal/state"
	"github.com/urfave/cli/v2"
)

// InstallCommand returns the command for installing an MCP
func InstallCommand() *cli.Command {
	return &cli.Command{
		Name:      "install",
		Usage:     "Install an MCP from an embedded definition",
		ArgsUsage: "[name]",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {

			if c.NArg() < 1 {
				return fmt.Errorf("missing required argument: name of embedded MCP")
			}

			name := c.Args().Get(0)

			repo := definitions.NewRepository()
			// Get the embedded MCP configuration
			mcpConfig, err := repo.MCPConfig(name)
			if err != nil {
				return fmt.Errorf("getting embedded MCP configuration: %w", err)
			}

			stateManager, err := state.NewHomeStore()
			if err != nil {
				return fmt.Errorf("creating state manager: %w", err)
			}

			builder, err := build.NewBuilder(mcpConfig)
			if err != nil {
				return fmt.Errorf("creating builder: %w", err)
			}
			defer builder.CleanUp() // Clean up temporary directory when done

			// Build the image from the repository
			tag, err := builder.DockerImage()
			if err != nil {
				return fmt.Errorf("building image: %w", err)
			}

			// Load or create MCP state
			mcpState, err := stateManager.Load(name)
			if err != nil {
				return fmt.Errorf("loading MCP state: %w", err)
			}

			mcpState.LocalImageTag = tag

			// Prompt for environment variables
			if err := prompt.PromptEnvironmentVariables(mcpConfig, mcpState); err != nil {
				return fmt.Errorf("getting environment variables: %w", err)
			}

			// Configure server in all available hosts
			hostManager := host.DefaultManager()
			availableHosts, err := hostManager.List()
			if err != nil {
				return fmt.Errorf("listing available hosts: %w", err)
			}

			if len(availableHosts) > 0 {
				fmt.Println("\nAvailable hosts:")
				selectedHosts := make(map[string]bool)
				for _, h := range availableHosts {
					selectedHosts[h] = true // Default to all selected
				}

				// Prompt user to select hosts
				options := make([]string, len(availableHosts))
				for i, h := range availableHosts {
					options[i] = h
				}
				selected, err := prompt.MultiSelect("Select hosts to configure (space to toggle, enter to confirm):", options, selectedHosts)
				if err != nil {
					return fmt.Errorf("selecting hosts: %w", err)
				}

				if len(selected) > 0 {
					fmt.Println("\nConfiguring server in selected hosts:")
					for _, h := range selected {
						fmt.Printf("  Configuring in %s... ", h)
						if err := hostManager.Connect(h, name); err != nil {
							fmt.Printf("failed: %v\n", err)
							continue
						}
						mcpState.ConfiguredHosts = append(mcpState.ConfiguredHosts, h)
						fmt.Println("done")
					}
				} else {
					fmt.Println("\nNo hosts selected for configuration")
				}
			}

			// Save the state
			if err := stateManager.Save(mcpState); err != nil {
				return fmt.Errorf("saving MCP state: %w", err)
			}

			fmt.Printf("\nMCP '%s' installed successfully\n", name)
			return nil
		},
	}
}
