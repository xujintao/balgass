EMAIL=xujintao@126.com
PASSWORD=1234
PGADMIN_DATA=~/balgass/docker/pgadmin/data

docker run \
--restart always \
-d \
--name pgadmin \
-u root \
-e TZ=Asia/Shanghai \
-e PGADMIN_DEFAULT_EMAIL=$EMAIL \
-e PGADMIN_DEFAULT_PASSWORD=$PASSWORD \
-e PGADMIN_LISTEN_PORT=8084 \
-v $PGADMIN_DATA:/var/lib/pgadmin \
-p 8084:8084 \
dpage/pgadmin4:6.7
