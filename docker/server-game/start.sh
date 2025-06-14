#!/bin/bash

TAG=${1:-latest}

CONFIG_PATH=~/balgass/config/server-game
COMMON_PATH=~/balgass/config/server-game-common

docker run -d \
--name server-game \
--restart always \
-v $CONFIG_PATH:/etc/server-game \
-v $COMMON_PATH:/etc/server-game-common \
-e TZ=Asia/Shanghai \
-e CONFIG_PATH=/etc/server-game \
-e COMMON_PATH=/etc/server-game-common \
-p 8080:8080 \
-p 56900:56900 \
xujintao/server-game:$TAG
