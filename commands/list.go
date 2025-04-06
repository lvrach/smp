package commands

import (
	"fmt"

	"github.com/lvrach/smp/definitions"
	"github.com/urfave/cli/v2"
)

// ListCommand returns the command for listing installed MCPs
func ListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List installed MCPs",
		Action: func(c *cli.Context) error {
			repo := definitions.NewRepository()
			mcpConfigs, err := repo.ListMCPs()
			if err != nil {
				return fmt.Errorf("failed to list MCPs: %w", err)
			}

			for _, mcpConfig := range mcpConfigs {
				fmt.Printf("  %s\n", mcpConfig)
			}

			return nil
		},
	}
}
