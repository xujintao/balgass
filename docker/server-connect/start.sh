#!/usr/bin/env bash
set -euo pipefail

TAG=${1:-latest}

REPO_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)
CONFIG_PATH="$REPO_DIR/docker/server-connect"
ENV_FILE="$CONFIG_PATH/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "Copy $CONFIG_PATH/.env.example to $ENV_FILE before starting server-connect." >&2
  exit 1
fi

docker run -d \
--name server-connect \
--restart always \
-v "$CONFIG_PATH:/etc/server-connect" \
--env-file "$ENV_FILE" \
-p 44405:44405 \
-p 55667:55667/udp \
"xujintao/server-connect:$TAG"
