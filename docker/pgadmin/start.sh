#!/bin/bash
set -euo pipefail

PGADMIN_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
ENV_FILE="$PGADMIN_DIR/.env"

for key in PGADMIN_DEFAULT_EMAIL PGADMIN_DEFAULT_PASSWORD; do
    if [[ ! -f "$ENV_FILE" ]] || ! grep -Eq "^${key}=.+$" "$ENV_FILE"; then
        echo "Set $key in $ENV_FILE before starting pgAdmin." >&2
        exit 1
    fi
done

docker volume create pgadmin_data

docker run \
    --restart always \
    -d \
    --name pgadmin \
    -u root \
    --env-file "$ENV_FILE" \
    -v pgadmin_data:/var/lib/pgadmin \
    -p 8084:8084 \
    dpage/pgadmin4:6.7
