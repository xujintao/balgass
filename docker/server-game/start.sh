#!/usr/bin/env bash
set -euo pipefail

TAG=${1:-latest}

REPO_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)
CONFIG_PATH="$REPO_DIR/docker/server-game"
COMMON_PATH="$REPO_DIR/config/server-game-common"
ENV_FILE="$CONFIG_PATH/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "Copy $CONFIG_PATH/.env.example to $ENV_FILE before starting server-game." >&2
  exit 1
fi

docker run -d \
--name server-game \
--restart always \
--add-host host.docker.internal:host-gateway \
-v "$CONFIG_PATH:/etc/server-game" \
-v "$COMMON_PATH:/etc/server-game-common" \
--env-file "$ENV_FILE" \
-p 8080:8080 \
-p 56900:56900 \
"xujintao/server-game:$TAG"
