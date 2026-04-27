#!/usr/bin/env bash
set -euo pipefail

echo "Public IP:"
public_ip=""
if command -v curl >/dev/null 2>&1; then
    public_ip="$(curl -4 --silent --show-error ifconfig.me || true)"
    echo "$public_ip"
    echo
else
    echo "curl is not installed. Install curl or check your public IP from a browser."
fi

echo
echo "Local LAN IP address candidates:"
if command -v ip >/dev/null 2>&1; then
    ip -brief addr show | awk '$1 != "lo" { print }'
elif command -v hostname >/dev/null 2>&1; then
    hostname -I
else
    echo "Could not find ip or hostname command."
fi

cat <<'EOF'

Next check your router admin page and find its WAN / Internet IP address.

Compare:
- Public IP shown above
- Router WAN / Internet IP shown by your router

If they are the same, DNS A Record Mode may work after port forwarding.
If they are different, your ISP may be using CGNAT. In that case, inbound
port forwarding usually will not work, and Cloudflare Tunnel is the easier
fallback.
EOF

if [ -n "$public_ip" ]; then
    cat <<EOF

Friend SSH command:
ssh vps-guest@$public_ip

This SSH command works from:
- Windows Command Prompt
- Windows PowerShell
- Git Bash
- Linux/macOS Terminal

Windows PowerShell reachability test:
Test-NetConnection $public_ip -Port 22

If the test times out, check router TCP port forwarding for port 22 to this
laptop's LAN IP address and make sure the friend typed the IP exactly.
EOF
fi
