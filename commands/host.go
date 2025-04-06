package commands

import (
	"fmt"

	"github.com/lvrach/smp/internal/host"
	"github.com/urfave/cli/v2"
)

// HostCommand returns the command for managing MCP hosts
func HostCommand() *cli.Command {
	return &cli.Command{
		Name:  "host",
		Usage: "Manage MCP hosts",
		Subcommands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List available MCP hosts",
				Action: func(c *cli.Context) error {
					manager := host.DefaultManager()
					hosts, err := manager.List()
					if err != nil {
						return fmt.Errorf("failed to list hosts: %w", err)
					}

					if len(hosts) == 0 {
						fmt.Println("No available hosts found")
						return nil
					}

					fmt.Println("Available hosts:")
					for _, h := range hosts {
						fmt.Printf("  %s\n", h)
					}

					return nil
				},
			},
		},
	}
}
