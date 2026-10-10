## server-game

1 Build image

```
REPO=~/balgass
FILE=./docker/server-game/Dockerfile

cd $REPO

docker build \
-t xujintao/server-game:latest \
-f $FILE \
.
```

2 Copy the environment example to the config directory and set `GAME_API_TOKEN`:

```
cp docker/server-game/.env.example docker/server-game/.env
```

On the VPS, keep the example's `/etc/server-game` and `/etc/server-game-common` paths in `.env`. For local Go runs, set `PATH_CONFIG` and `PATH_COMMON` in your local copy of `.env` to the absolute paths of `docker/server-game` and `config/server-game-common`. Docker mounts those directories at the `/etc` paths and passes the VPS `.env` to the container without overriding its variables.

3 Run image from any directory:

```
./docker/server-game/start.sh
```
