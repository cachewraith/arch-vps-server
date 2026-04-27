# Lab CLI Manager

The `lab` CLI is a unified tool written in Go to manage the Arch VPS Server lab. It replaces the legacy bash scripts and aliases with a more robust and automated interface.

## Commands

### Proxy Management
*   `lab proxy up`: Starts the Caddy proxy.
*   `lab proxy down`: Stops the Caddy proxy.
*   `lab proxy restart`: Restarts the Caddy proxy.
*   `lab proxy logs`: Follows Caddy logs.

### Network Management
*   `lab network create`: Ensures the `arch-vps-net` Docker network exists.

### Project Management
*   `lab project update <name>`: Pulls the latest changes from Git and rebuilds the project container.
*   `lab project add <name> --domain <domain> --port <port>`: 
    *   Creates a new directory in `projects/`.
    *   Generates a default `docker-compose.yml`.
    *   Automatically appends the domain and reverse proxy configuration to the `Caddyfile`.

### Status
*   `lab status`: Lists all running containers attached to the `arch-vps-net` network.

## Installation / Usage
To use the CLI from anywhere, add the binary to your path or create an alias:

```bash
alias lab='/mnt/storage/arch-vps-server/lab-cli/lab'
```

## Legacy
The previous bash scripts in `scripts/` and aliases in `aliases.sh` are now deprecated.
