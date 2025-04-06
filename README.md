# SMP - Secure MCP Manager

SMP is a command-line tool for managing MCPs.

## Installation

To install SMP, run:

```bash
go install github.com/lvrach/smp@latest
```

## Basic Commands

### `smp install [name]`
Installs an MCP with the specified name. This command sets up the necessary configuration and environment for the MCP.

### `smp uninstall [name]`
Uninstalls an MCP with the specified name. This removes the MCP's configuration and cleans up associated resources.


### `smp run [name]`
Runs the container for the specified MCP. This command starts the MCP with its configured environment and settings.

### `smp list`
Lists all available MCPs and their current status.

### `smp host`
Manages host-specific settings and configurations for MCPs.
