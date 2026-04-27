# Arch VPS Server Lab

This project turns this folder into a clean Docker-based VPS-style server environment for an Arch laptop.

## 🚀 Lab CLI (New!)

We now have a unified CLI tool called `lab` to manage everything. **This is the recommended way to use the lab.**

### Setup
Add the following alias to your `~/.bashrc`:
```bash
alias lab='/mnt/storage/arch-vps-server/lab-cli/lab'
```

### Core Commands
*   `lab proxy up/down/restart`: Manage the Caddy proxy.
*   `lab project add <name> --domain <domain> --port <port>`: Scaffold a new project and update Caddy automatically.
*   `lab project update <name>`: Git pull and rebuild a project.
*   `lab status`: See what's running.

**Detailed CLI Documentation: [Lab CLI Manager](docs/12-lab-cli.md)**

---

## Folder Layout

```text
/mnt/storage/arch-vps-server/
├── proxy/
│   ├── docker-compose.yml
│   └── Caddyfile
├── projects/
├── lab-cli/       <-- New Go CLI
├── databases/
├── tunnels/
├── docs/
├── scripts/       <-- Legacy scripts
├── backups/
├── logs/
└── README.md
```

## Detailed Docs

- [Overview](docs/01-overview.md)
- [Lab CLI Manager](docs/12-lab-cli.md) (Recommended)
- [Install Docker on Arch](docs/02-install-docker-arch.md)
- [Caddy Proxy](docs/03-caddy-proxy.md)
- [Projects](docs/04-projects.md)
- [Local Domains](docs/05-local-domains.md)
- [Public DNS A Record Mode](docs/06-public-dns-a-record-mode.md)
- [Cloudflare Tunnel Mode](docs/07-cloudflare-tunnel-mode.md)
- [Security](docs/08-security.md)
- [Move to a Real VPS](docs/09-move-to-real-vps.md)
- [Troubleshooting](docs/10-troubleshooting.md)
- [Aliases](docs/11-aliases.md) (Legacy)

... (rest of the original README content)
