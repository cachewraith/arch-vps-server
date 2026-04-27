#!/usr/bin/env bash
set -euo pipefail

if [ "${1:-}" = "" ]; then
    echo "Usage: $0 /path/to/project"
    exit 1
fi

PROJECT_DIR="$1"

if [ ! -d "$PROJECT_DIR" ]; then
    echo "Error: project directory does not exist: $PROJECT_DIR"
    exit 1
fi

if [ ! -d "$PROJECT_DIR/.git" ]; then
    echo "Error: project directory is not a Git repository: $PROJECT_DIR"
    exit 1
fi

if [ -f "$PROJECT_DIR/docker-compose.yml" ]; then
    COMPOSE_FILE="docker-compose.yml"
elif [ -f "$PROJECT_DIR/compose.yml" ]; then
    COMPOSE_FILE="compose.yml"
else
    echo "Error: no docker-compose.yml or compose.yml found in: $PROJECT_DIR"
    exit 1
fi

cd "$PROJECT_DIR"
echo "Updating Git repository in $PROJECT_DIR..."
git pull

echo "Rebuilding and starting services using $COMPOSE_FILE..."
docker compose -f "$COMPOSE_FILE" up -d --build
