DOCKER_DIR=~/balgass/docker
LOKI_DIR=$DOCKER_DIR/loki
LOKI_DATA=$LOKI_DIR/data

docker volume create loki_data

docker run \
--restart always \
-d \
--name loki \
-e LANG=C.UTF-8 \
-e TZ=UTC \
-v loki_data:/loki \
-p 3100:3100 \
grafana/loki:3.5.1
