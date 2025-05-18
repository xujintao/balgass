#!/bin/bash

TAG=${1:-latest}

docker run -d \
--name server_web \
--restart always \
-e TZ=Asia/Shanghai \
-p 8000:8000 \
xujintao/server_web:$TAG