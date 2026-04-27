# Security

This setup can expose your laptop to the internet. Treat it like a real server.

## Database Ports

Do not expose PostgreSQL, MySQL, MongoDB, Redis, or other database ports publicly.

Prefer `expose` instead of `ports` for internal backend and database services:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    expose:
      - "5432"
    networks:
      - arch-vps-net
```

Only containers on the shared Docker network should reach database containers.

## Guest SSH Access

A dedicated user `vps-guest` has been created to allow friends to access the server. This user is restricted to the `/mnt/storage/arch-vps-server` directory and belongs to the `docker` group.

To allow a friend to connect:

1.  Get their public SSH key.
2.  Run the following command:
    ```bash
    avps-ssh add 'ssh-rsa AAA...'
    ```
3.  Check the exact public IP and SSH command:
    ```bash
    avps-ip
    ```
4.  Your friend can then connect using the exact IP printed by `avps-ip`:
    ```bash
    ssh vps-guest@your_public_ip_address
    ```

The same SSH command works from these client shells:

- Windows Command Prompt
- Windows PowerShell
- Git Bash
- Linux/macOS Terminal

If your friend uses Git Bash, Command Prompt, or PowerShell, the command is
still the same:

```bash
ssh vps-guest@your_public_ip_address
```

If the connection times out, the problem is network reachability, not the SSH
key. Confirm that:

- Your friend typed the public IP exactly.
- Your router forwards TCP port `22` to this laptop's LAN IP address.
- Your router WAN / Internet IP matches the public IP shown by `avps-ip`.
- Any local firewall allows TCP port `22`.

On Windows, your friend can test the route with:

```powershell
Test-NetConnection your_public_ip_address -Port 22
```

If this returns `TcpTestSucceeded: False`, fix router forwarding, firewall, or
CGNAT before debugging SSH keys.

## General Notes

- Use `.env` files for secrets.
- Do not commit `.env` files.
- Keep Docker images updated.
- Keep Arch updated.
- Use strong passwords.
- Back up databases.
- Be careful because this exposes your laptop to the internet.
