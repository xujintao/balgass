#!/bin/bash
set -e

TAG=${1:-latest}
ENV_FILE=~/balgass/config/server_web/.env

docker run \
--restart always \
-d \
--name server_web \
--env-file $ENV_FILE \
-p 8000:8000 \
xujintao/server_web:$TAG
