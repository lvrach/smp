package main

import (
	"log"
	"os"

	"github.com/lvrach/smp/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "smp",
		Usage:       "Secure MCP Manager",
		Description: "SMP is a tool for managing MCPs. ",
		Commands: []*cli.Command{
			commands.InstallCommand(),
			commands.UninstallCommand(),
			commands.BuildCommand(),
			commands.RunCommand(),
			commands.ListCommand(),
			commands.HostCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
