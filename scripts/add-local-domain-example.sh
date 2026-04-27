#!/usr/bin/env bash
set -euo pipefail

cat <<'EOF'
This script does not edit /etc/hosts automatically.

To add local test domains, manually add lines like these to /etc/hosts:

127.0.0.1 archvps.local
127.0.0.1 api.localhost.test
127.0.0.1 app.localhost.test
127.0.0.1 admin.localhost.test

On Arch, you can edit it with:
sudo nano /etc/hosts
EOF
