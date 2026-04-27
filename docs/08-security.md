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
3.  Your friend can then connect using:
    ```bash
    ssh vps-guest@your_ip_address
    ```

## General Notes

- Use `.env` files for secrets.
- Do not commit `.env` files.
- Keep Docker images updated.
- Keep Arch updated.
- Use strong passwords.
- Back up databases.
- Be careful because this exposes your laptop to the internet.
