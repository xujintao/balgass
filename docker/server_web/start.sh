#!/bin/bash

TAG=${1:-latest}

ENV_FILE=~/balgass/config/server_web/.env

docker run -d \
--name server_web \
--restart always \
--env-file $ENV_FILE \
-e TZ=Asia/Shanghai \
-p 8000:8000 \
xujintao/server_web:$TAG