package commands

import (
	"fmt"
	"os"

	"github.com/lvrach/smp/definitions"
	"github.com/lvrach/smp/internal/docker"
	"github.com/lvrach/smp/internal/state"
	"github.com/lvrach/smp/keystore"
	"github.com/urfave/cli/v2"
)

// RunCommand returns the command for running an MCP container
func RunCommand() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Run a container for an MCP",
		ArgsUsage: "[name]",
		Flags:     []cli.Flag{},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				fmt.Fprintf(os.Stderr, "Error: missing required argument: name\n")
				return fmt.Errorf("missing required argument: name")
			}

			name := c.Args().Get(0)

			repo := definitions.NewRepository()
			mcpConfig, err := repo.MCPConfig(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to get MCP configuration: %v\n", err)
				return fmt.Errorf("failed to get MCP configuration: %w", err)
			}

			stateManager, err := state.NewHomeStore()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to create state manager: %v\n", err)
				return fmt.Errorf("failed to create state manager: %w", err)
			}

			mcpState, err := stateManager.Load(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to load MCP state: %v\n", err)
				return fmt.Errorf("failed to load MCP state: %w", err)
			}

			kc := keystore.KeyChain{}

			for accountKey, env := range mcpState.KeyChainEnvVars {
				secret, err := kc.Retrieve(keystore.AccountType(accountKey))
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: failed to retrieve %q key: %v\n", accountKey, err)
					return fmt.Errorf("failed to retrieve %q key: %w", accountKey, err)
				}
				// TODO: fix me, this should populate a different env var struct
				mcpState.SetEnvironmentVariable(env, secret)
			}

			runner := docker.NewRunner(mcpConfig, mcpState)
			if err := runner.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to run container: %v\n", err)
				return fmt.Errorf("failed to run container: %w", err)
			}

			return nil
		},
	}
}
