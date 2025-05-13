DOCKER_DIR=~/github.com/xujintao/balgass/docker
NGINX_DIR=$DOCKER_DIR/nginx

docker run \
--restart always \
-d \
--name nginx \
-e LANG=C.UTF-8 \
-e TZ=Asia/Shanghai \
-v $NGINX_DIR:/etc/nginx \
-p 80:80 \
-p 443:443 \
nginx:1.15.5-alpine