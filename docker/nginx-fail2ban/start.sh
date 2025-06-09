DOCKER_DIR=~/balgass/docker
NGINX_DIR=$DOCKER_DIR/nginx-fail2ban/nginx
FAIL2BAN_DIR=$DOCKER_DIR/nginx-fail2ban/fail2ban
LOGS_DIR=$DOCKER_DIR/nginx-fail2ban/logs

docker run \
--restart always \
-d \
--name nginx-fail2ban \
-e LANG=C.UTF-8 \
-e TZ=UTC \
-v $NGINX_DIR:/etc/nginx \
-v $LOGS_DIR:/var/log/nginx \
-v $FAIL2BAN_DIR/fail2ban.local:/etc/fail2ban/fail2ban.local \
-v $FAIL2BAN_DIR/jail.local:/etc/fail2ban/jail.local \
-v $FAIL2BAN_DIR/nginx-bot-request.conf:/etc/fail2ban/filter.d/nginx-bot-request.conf \
-v $LOGS_DIR:/var/log/fail2ban \
-p 80:80 \
-p 443:443 \
--cap-add=NET_ADMIN \
xujintao/nginx-fail2ban:latest
