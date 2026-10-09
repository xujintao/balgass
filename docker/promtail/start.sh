DOCKER_DIR=~/balgass/docker
PROMTAIL_DIR=$DOCKER_DIR/promtail
CADDY_FAIL2BAN_LOG_DIR=$DOCKER_DIR/caddy-fail2ban/logs

docker run \
--restart always \
-d \
--name promtail \
--add-host host.docker.internal:host-gateway \
-e LANG=C.UTF-8 \
-e TZ=UTC \
-v $PROMTAIL_DIR/config.yml:/etc/promtail/config.yml:ro \
-v $CADDY_FAIL2BAN_LOG_DIR:/var/log/caddy-fail2ban:ro \
-v /var/run/docker.sock:/var/run/docker.sock:ro \
-v /var/lib/docker/containers:/var/lib/docker/containers:ro \
-p 9080:9080 \
grafana/promtail:3.5.1
