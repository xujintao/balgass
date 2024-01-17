USER=root
PASSWORD=1234
DEFAULT_DB=django
PGSQL_DATA=~/r2f2/pgsql/data
PGSQL_INIT=~/r2f2/pgsql/initdb.d

docker run \
--restart always \
-d \
--name pgsql \
-e TZ=Asia/Shanghai \
-e POSTGRES_USER=$USER \
-e POSTGRES_PASSWORD=$PASSWORD \
-e POSTGRES_DB=$DEFAULT_DB \
-e ALLOW_IP_RANGE=0.0.0.0/0 \
-v $PGSQL_DATA:/var/lib/postgresql/data \
-v $PGSQL_INIT:/docker-entrypoint-initdb.d \
-p 5432:5432 \
postgres:14.2-alpine
