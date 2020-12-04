export NGINX_HOME=~/volumes/nginx-mu

docker run \
--restart always \
-d \
--name nginx-mu \
-e LANG=C.UTF-8 \
-e TZ=Asia/Shanghai \
-v $NGINX_HOME:/etc/nginx \
-p 80:80 \
-p 443:443 \
nginx:1.15.5-alpine
