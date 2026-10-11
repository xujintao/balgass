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

All service specific settings live in `.env`. Set `DB_PASSWORD` and `GAME_API_TOKEN` for your deployment. Keep the example's `/etc/server-game-common` value for `PATH_COMMON` on the VPS; for local Go runs, set it to the absolute path of `config/server-game-common`. Docker mounts only that common directory and passes `.env` to the container.

3 Run image from any directory:

```
./docker/server-game/start.sh
```
