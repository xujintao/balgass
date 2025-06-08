DOCKER_DIR=~/balgass/docker
PROMTAIL_DIR=$DOCKER_DIR/promtail
NGINX_FAIL2BAN_LOG_DIR=$DOCKER_DIR/nginx-fail2ban/logs
MAIL_LOG_DIR=$DOCKER_DIR/mail/logs

docker run \
--restart always \
-d \
--name promtail \
-e LANG=C.UTF-8 \
-e TZ=UTC \
-v $PROMTAIL_DIR/config.yml:/etc/promtail/config.yml:ro \
-v $NGINX_FAIL2BAN_LOG_DIR:/var/log/nginx-fail2ban:ro \
-v $MAIL_LOG_DIR:/var/log/mail:ro \
-p 9080:9080 \
grafana/promtail:3.5.1
