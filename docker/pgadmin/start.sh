EMAIL=xujintao@126.com
PASSWORD=1234
PGADMIN_DATA=~/r2f2/pgadmin/data

docker run \
--restart always \
-d \
--name pgadmin \
-u root \
-e TZ=Asia/Shanghai \
-e PGADMIN_DEFAULT_EMAIL=$EMAIL \
-e PGADMIN_DEFAULT_PASSWORD=$PASSWORD \
-v $PGADMIN_DATA:/var/lib/pgadmin \
-p 8084:80 \
dpage/pgadmin4:6.7
