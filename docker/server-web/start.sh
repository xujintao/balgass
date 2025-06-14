#!/bin/bash
set -e

TAG=${1:-latest}
ENV_FILE=~/balgass/config/server-web/.env

docker run \
--restart always \
-d \
--name server-web \
--env-file $ENV_FILE \
-p 8000:8000 \
xujintao/server-web:$TAG
