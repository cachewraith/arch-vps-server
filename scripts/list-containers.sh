#!/usr/bin/env bash
set -euo pipefail

echo "Running containers:"
docker ps

echo
echo "All containers:"
docker ps -a
