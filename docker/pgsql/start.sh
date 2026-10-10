#!/bin/bash
set -euo pipefail

PGSQL_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
ENV_FILE="$PGSQL_DIR/.env"

for key in POSTGRES_USER POSTGRES_PASSWORD POSTGRES_DB; do
    if [[ ! -f "$ENV_FILE" ]] || ! grep -Eq "^${key}=.+$" "$ENV_FILE"; then
        echo "Set $key in $ENV_FILE before starting PostgreSQL." >&2
        exit 1
    fi
done

docker volume create pgsql_data

docker run \
    --restart always \
    -d \
    --name pgsql \
    --env-file "$ENV_FILE" \
    -v pgsql_data:/var/lib/postgresql/data \
    -p 5432:5432 \
    postgres:14.2-alpine
