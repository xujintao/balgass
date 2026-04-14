#!/bin/bash

TAG=${1:-latest}

CONFIG_PATH=~/balgass/config/server-connect
ENV_FILE=~/balgass/config/server-connect/.env

docker run -d \
--name server-connect \
--restart always \
-v $CONFIG_PATH:/etc/server-connect \
-e TZ=Asia/Shanghai \
--env-file $ENV_FILE \
-p 44405:44405 \
-p 55667:55667/udp \
xujintao/server-connect:$TAG
