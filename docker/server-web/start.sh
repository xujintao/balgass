#!/bin/bash
set -e

TAG=${1:-latest}

GAME_CONFIG_PATH=~/balgass/config/server-game-common
ENV_FILE=~/balgass/config/server-web/.env

docker run \
--restart always \
-d \
--name server-web \
-v $GAME_CONFIG_PATH:/etc/server-game-common \
--env-file $ENV_FILE \
-p 8000:8000 \
xujintao/server-web:$TAG
