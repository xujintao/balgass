#!/usr/bin/env bash
set -euo pipefail

TAG=${1:-latest}

REPO_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)
COMMON_PATH="$REPO_DIR/config/server-game-common"
ENV_FILE="$REPO_DIR/docker/server-game/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "Copy $REPO_DIR/docker/server-game/.env.example to $ENV_FILE before starting server-game." >&2
  exit 1
fi

set -a
. "$ENV_FILE"
set +a

docker run -d \
--name server-game \
--restart always \
--add-host host.docker.internal:host-gateway \
-v "$COMMON_PATH:/etc/server-game-common" \
--env-file "$ENV_FILE" \
-p "$HTTP_PORT:$HTTP_PORT" \
-p "$GAME_SERVER_PORT:$GAME_SERVER_PORT" \
"xujintao/server-game:$TAG"
