#!/bin/bash

TAG=${1:-latest}

CONFIG_PATH=~/balgass/config/server_game
COMMON_PATH=~/balgass/config/server_game_common

docker run -d \
--name server_game \
--restart always \
-v $CONFIG_PATH:/etc/server_game \
-v $COMMON_PATH:/etc/server_game_common \
-e TZ=Asia/Shanghai \
-e CONFIG_PATH=/etc/server_game \
-e COMMON_PATH=/etc/server_game_common \
-p 8080:8080 \
-p 56900:56900 \
xujintao/server_game:$TAG