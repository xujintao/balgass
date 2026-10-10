#!/bin/bash
set -euo pipefail

GRAFANA_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
ENV_FILE="$GRAFANA_DIR/.env"

for key in GF_SECURITY_ADMIN_USER GF_SECURITY_ADMIN_PASSWORD GF_SERVER_ROOT_URL; do
    if [[ ! -f "$ENV_FILE" ]] || ! grep -Eq "^${key}=.+$" "$ENV_FILE"; then
        echo "Set $key in $ENV_FILE before starting Grafana." >&2
        exit 1
    fi
done

docker run \
    --restart always \
    -d \
    --name grafana \
    --user root \
    -e LANG=C.UTF-8 \
    --env-file "$ENV_FILE" \
    -v "$GRAFANA_DIR/data:/var/lib/grafana" \
    -p 3000:3000 \
    grafana/grafana:12.0.1
