#!/usr/bin/env bash

# Source this file from your shell:
# source /mnt/storage/arch-vps-server/scripts/aliases.sh

export ARCH_VPS_ROOT="/mnt/storage/arch-vps-server"

alias avps='cd "$ARCH_VPS_ROOT"'
alias avps-net='"$ARCH_VPS_ROOT/scripts/create-network.sh"'
alias avps-up='"$ARCH_VPS_ROOT/scripts/start-proxy.sh"'
alias avps-down='"$ARCH_VPS_ROOT/scripts/stop-proxy.sh"'
alias avps-restart='"$ARCH_VPS_ROOT/scripts/restart-proxy.sh"'
alias avps-ps='"$ARCH_VPS_ROOT/scripts/list-containers.sh"'
alias avps-ip='"$ARCH_VPS_ROOT/scripts/check-public-ip.sh"'
alias avps-hosts='"$ARCH_VPS_ROOT/scripts/add-local-domain-example.sh"'

avps-update() {
    if [ "${1:-}" = "" ]; then
        echo "Usage: avps-update /path/to/project"
        return 1
    fi

    "$ARCH_VPS_ROOT/scripts/update-project.sh" "$1"
}

avps-logs() {
    docker logs -f arch-vps-caddy
}

avps-caddy() {
    "${EDITOR:-nano}" "$ARCH_VPS_ROOT/proxy/Caddyfile"
}

avps-host-backend() {
    if [ "$#" -ne 4 ]; then
        cat <<'EOF'
Usage:
  avps-host-backend DOMAIN CONTAINER_NAME INTERNAL_PORT PROJECT_DIR

Example:
  avps-host-backend api.yourdomain.com my-backend 3000 ./projects/my-backend

This prints safe snippets to add to your backend docker-compose.yml and proxy/Caddyfile.
It does not edit files automatically.
EOF
        return 1
    fi

    local domain="$1"
    local container_name="$2"
    local internal_port="$3"
    local project_dir="$4"

    cat <<EOF
Backend hosting checklist
=========================

1. Make sure the project exists:

   $project_dir

2. In the backend docker-compose.yml, make sure the service has:

services:
  $container_name:
    build: .
    container_name: $container_name
    restart: unless-stopped
    expose:
      - "$internal_port"
    env_file:
      - .env
    networks:
      - arch-vps-net

networks:
  arch-vps-net:
    external: true

3. Start or rebuild the backend:

   cd "$project_dir"
   docker compose up -d --build

4. Add this to $ARCH_VPS_ROOT/proxy/Caddyfile:

$domain {
    reverse_proxy $container_name:$internal_port
}

5. Restart Caddy:

   avps-restart

6. Test:

   curl -I https://$domain
EOF
}

avps-help() {
    cat <<'EOF'
Arch VPS aliases
================

Navigation:
  avps                 cd to /mnt/storage/arch-vps-server

Proxy:
  avps-net             create Docker network arch-vps-net
  avps-up              start Caddy proxy
  avps-down            stop Caddy proxy
  avps-restart         restart Caddy proxy
  avps-caddy           edit proxy/Caddyfile
  avps-logs            follow Caddy logs

Status:
  avps-ps              list Docker containers
  avps-ip              check public IP and LAN IPs
  avps-hosts           print local /etc/hosts examples

Projects:
  avps-update PATH     git pull and docker compose up -d --build
  avps-host-backend DOMAIN CONTAINER_NAME INTERNAL_PORT PROJECT_DIR
                       print backend Compose and Caddy hosting snippets

Example:
  avps-host-backend api.yourdomain.com my-backend 3000 ./projects/my-backend
EOF
}
