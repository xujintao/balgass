#!/bin/bash
set -euo pipefail

PROMTAIL_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
DOCKER_DIR=$(dirname -- "$PROMTAIL_DIR")
ENV_FILE="$PROMTAIL_DIR/.env"
CADDY_FAIL2BAN_LOG_DIR="$DOCKER_DIR/caddy-fail2ban/logs"

if [[ ! -f "$ENV_FILE" ]]; then
    echo "Copy $PROMTAIL_DIR/.env.example to $ENV_FILE before starting Promtail." >&2
    exit 1
fi

docker run \
    --restart always \
    -d \
    --name promtail \
    --add-host host.docker.internal:host-gateway \
    -e LANG=C.UTF-8 \
    --env-file "$ENV_FILE" \
    -v "$PROMTAIL_DIR/config.yml:/etc/promtail/config.yml:ro" \
    -v "$CADDY_FAIL2BAN_LOG_DIR:/var/log/caddy-fail2ban:ro" \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -v /var/lib/docker/containers:/var/lib/docker/containers:ro \
    -p 9080:9080 \
    grafana/promtail:3.5.1
