#!/bin/bash

TAG=${1:-latest}

CONFIG_PATH=~/balgass/config/server-game
COMMON_PATH=~/balgass/config/server-game-common
ENV_FILE=~/balgass/config/server-game/.env

docker run -d \
--name server-game \
--restart always \
-v $CONFIG_PATH:/etc/server-game \
-v $COMMON_PATH:/etc/server-game-common \
-e TZ=UTC \
--env-file $ENV_FILE \
-p 8080:8080 \
-p 56900:56900 \
xujintao/server-game:$TAG
