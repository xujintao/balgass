USER=root
PASSWORD=1234
DEFAULT_DB=django
PGSQL_DATA=~/balgass/docker/pgsql/data
PGSQL_INIT=~/balgass/docker/pgsql/initdb.d

docker volume create pgsql_data

docker run \
--restart always \
-d \
--name pgsql \
-e TZ=Asia/Shanghai \
-e POSTGRES_USER=$USER \
-e POSTGRES_PASSWORD=$PASSWORD \
-e POSTGRES_DB=$DEFAULT_DB \
-e ALLOW_IP_RANGE=0.0.0.0/0 \
-v pgsql_data:/var/lib/postgresql/data \
-p 5432:5432 \
postgres:14.2-alpine
