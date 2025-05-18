#!/bin/bash

TAG=${1:-latest}

CONFIG_PATH=~/balgass/config/server_connect

docker run -d \
--name server_connect \
--restart always \
-v $CONFIG_PATH:/etc/server_connect \
-e TZ=Asia/Shanghai \
-e CONFIG_PATH=/etc/server_connect \
-p 44405:44405 \
-p 55667:55667/udp \
xujintao/server_connect:$TAG