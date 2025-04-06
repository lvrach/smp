package commands

import (
	"fmt"

	"github.com/lvrach/smp/internal/docker"
	"github.com/lvrach/smp/internal/host"
	"github.com/lvrach/smp/internal/state"
	"github.com/urfave/cli/v2"
)

// UninstallCommand returns the command for uninstalling an MCP
func UninstallCommand() *cli.Command {
	return &cli.Command{
		Name:      "uninstall",
		Usage:     "Uninstall an MCP server and disconnect from all hosts",
		ArgsUsage: "[name]",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("missing required argument: name")
			}

			name := c.Args().Get(0)

			stateManager, err := state.NewHomeStore()
			if err != nil {
				return fmt.Errorf("creating state manager: %w", err)
			}

			mcpState, err := stateManager.Load(name)
			if err != nil {
				return fmt.Errorf("loading MCP state: %w", err)
			}

			if mcpState.LocalImageTag != "" {
				if err := docker.DeleteImage(mcpState.LocalImageTag); err != nil {
					return fmt.Errorf("deleting Docker image %q: %w", mcpState.LocalImageTag, err)
				}
			}

			hostManager := host.DefaultManager()
			for _, h := range mcpState.ConfiguredHosts {
				wasConfigured, err := hostManager.Disconnect(h, name)
				if err != nil {
					return fmt.Errorf("disconnecting from host %q: %w", h, err)
				}
				if !wasConfigured {
					fmt.Printf("Warning: server %q was not configured in host %q\n", name, h)
				}
			}

			// Delete the state
			if err := stateManager.Delete(name); err != nil {
				return fmt.Errorf("deleting MCP state: %w", err)
			}

			fmt.Printf("MCP '%s' uninstalled successfully\n", name)
			return nil
		},
	}
}
