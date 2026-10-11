#!/usr/bin/env bash
set -euo pipefail

TAG=${1:-latest}

REPO_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)
ENV_FILE="$REPO_DIR/docker/server-connect/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "Copy $REPO_DIR/docker/server-connect/.env.example to $ENV_FILE before starting server-connect." >&2
  exit 1
fi

set -a
. "$ENV_FILE"
set +a

docker run -d \
--name server-connect \
--restart always \
--env-file "$ENV_FILE" \
-p "$TCP_PORT:$TCP_PORT" \
-p "$UDP_PORT:$UDP_PORT/udp" \
"xujintao/server-connect:$TAG"
