DOCKER_DIR=~/balgass/docker
MAIL_DIR=$DOCKER_DIR/mail
MAIL_DATA=$MAIL_DIR/data
MAIL_STATE=$MAIL_DIR/state
MAIL_LOGS=$MAIL_DIR/logs
MAIL_CONFIG=$MAIL_DIR/config

NGINX_SSL_DIR=$DOCKER_DIR/nginx/ssl

docker run -d \
--name mailserver \
--restart always \
--hostname mail \
--domainname r2f2.com \
-p 25:25 \
-p 465:465 \
-p 587:587 \
-p 143:143 \
-p 993:993 \
-v $MAIL_DATA:/var/mail/ \
-v $MAIL_STATE:/var/mail-state/ \
-v $MAIL_LOGS:/var/log/mail/ \
-v $MAIL_CONFIG:/tmp/docker-mailserver/ \
-v $NGINX_SSL_DIR:/tmp/ssl/ \
-e TZ=Asia/Shanghai \
-e ENABLE_CLAMAV=0 \
-e ENABLE_FAIL2BAN=1 \
-e ENABLE_POSTGREY=1 \
-e ENABLE_SPAMASSASSIN=1 \
-e SSL_TYPE=manual \
-e SSL_CERT_PATH=/tmp/ssl/cert.pem \
-e SSL_KEY_PATH=/tmp/ssl/key.pem \
--cap-add=NET_ADMIN \
--cap-add=SYS_PTRACE \
--tty \
mailserver/docker-mailserver:15.0.2
