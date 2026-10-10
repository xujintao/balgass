## server-connect

1 Build image

```
REPO=~/balgass
FILE=./docker/server-connect/Dockerfile

cd $REPO

docker build \
-t xujintao/server-connect:latest \
-f $FILE \
.
```

2 Copy the environment example to the config directory:

```
cp docker/server-connect/.env.example docker/server-connect/.env
```

On the VPS, keep the example's `/etc/server-connect` path in `.env`. For local Go runs, set `PATH_CONFIG` in your local copy of `.env` to the absolute path of `docker/server-connect`. Docker mounts that directory at `/etc/server-connect` and passes the VPS `.env` to the container without overriding its variables.

3 Run image from any directory:

```
./docker/server-connect/start.sh
```
