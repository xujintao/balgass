#!/bin/bash
set -euo pipefail

LOKI_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
ENV_FILE="$LOKI_DIR/.env"

if [[ ! -f "$ENV_FILE" ]]; then
    echo "Copy $LOKI_DIR/.env.example to $ENV_FILE before starting Loki." >&2
    exit 1
fi

docker volume create loki_data

docker run \
    --restart always \
    -d \
    --name loki \
    -e LANG=C.UTF-8 \
    --env-file "$ENV_FILE" \
    -v loki_data:/loki \
    -p 3100:3100 \
    grafana/loki:3.5.1
