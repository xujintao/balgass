#!/bin/bash

TAG=${1:-latest}

CONFIG_PATH=~/balgass/config/server-connect

docker run -d \
--name server-connect \
--restart always \
-v $CONFIG_PATH:/etc/server-connect \
-e TZ=Asia/Shanghai \
-e CONFIG_PATH=/etc/server-connect \
-p 44405:44405 \
-p 55667:55667/udp \
xujintao/server-connect:$TAG
