DOCKER_DIR=~/balgass/docker
PROMTAIL_DIR=$DOCKER_DIR/promtail

docker run \
--restart always \
-d \
--name promtail \
-e LANG=C.UTF-8 \
-e TZ=UTC \
-v $PROMTAIL_DIR/config.yml:/etc/promtail/config.yml:ro \
-v $DOCKER_DIR/nginx-fail2ban/logs:/var/log/nginx-fail2ban:ro \
-p 9080:9080 \
grafana/promtail:3.5.1