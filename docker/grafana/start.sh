ADMIN_USER=admin
ADMIN_PASSWORD=1234
ROOT_URL=https://r2f2.com/grafana/
SERVE_FROM_SUB_PATH=true
DOCKER_DIR=~/balgass/docker
GRAFANA_DIR=$DOCKER_DIR/grafana
GRAFANA_DATA=$GRAFANA_DIR/data

docker run \
--restart always \
-d \
--name grafana \
--user root \
-e LANG=C.UTF-8 \
-e TZ=UTC \
-e GF_SECURITY_ADMIN_USER=$ADMIN_USER \
-e GF_SECURITY_ADMIN_PASSWORD=$ADMIN_PASSWORD \
-e GF_SERVER_ROOT_URL=$ROOT_URL \
-e GF_SERVER_SERVE_FROM_SUB_PATH=$SERVE_FROM_SUB_PATH \
-v $GRAFANA_DATA:/var/lib/grafana \
-p 3000:3000 \
grafana/grafana:12.0.1
