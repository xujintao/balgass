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

2 Copy the environment example to the config directory:

```
cp docker/server-game/.env.example docker/server-game/.env
```

All service specific settings live in `.env`. Set `DB_PASSWORD` and the client-facing `GAME_SERVER_HOST` for your deployment. `GAME_SERVER_HOST` may be an IPv4 address or a DNS name of at most 15 ASCII bytes, such as `game.r2f2.com`; enter only the host, without a scheme or port. `CONNECT_SERVER_HOST` is the IP address or DNS name used by the game server to contact `server-connect` and has no 15-byte client protocol limit. `SERVER_CODE`, `NAME`, `VISIBLE`, `GAME_SERVER_HOST`, and `GAME_SERVER_PORT` are sent to `server-connect` on every registration heartbeat. Keep the example's `/etc/server-game-common` value for `PATH_COMMON` on the VPS; for local Go runs, set it to the absolute path of `config/server-game-common`. Docker mounts only that common directory and passes `.env` to the container. Caddy authenticates `/api/command`; direct requests on port 8080 are unauthenticated, so keep the existing firewall restrictions in place.

3 Run image from any directory:

```
./docker/server-game/start.sh
```
