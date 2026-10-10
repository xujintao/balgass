#!/bin/bash
set -euo pipefail

ALLOY_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
DOCKER_DIR=$(dirname -- "$ALLOY_DIR")
ENV_FILE="$ALLOY_DIR/.env"
CADDY_FAIL2BAN_LOG_DIR="$DOCKER_DIR/caddy-fail2ban/logs"

if [[ ! -f "$ENV_FILE" ]]; then
    echo "Copy $ALLOY_DIR/.env.example to $ENV_FILE before starting Alloy." >&2
    exit 1
fi

docker volume create alloy_data

docker run \
    --restart always \
    -d \
    --name alloy \
    --add-host host.docker.internal:host-gateway \
    -e LANG=C.UTF-8 \
    --env-file "$ENV_FILE" \
    -v "$ALLOY_DIR/config.alloy:/etc/alloy/config.alloy:ro" \
    -v "$CADDY_FAIL2BAN_LOG_DIR:/var/log/caddy-fail2ban:ro" \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -v alloy_data:/var/lib/alloy/data \
    grafana/alloy:v1.20.1 \
    run --storage.path=/var/lib/alloy/data /etc/alloy/config.alloy
