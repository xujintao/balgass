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

All service settings live in `.env`. Set the update host for your deployment. Game server addresses and visibility come from each `server-game` registration. Keep `.env` private and edit `.env.example` when changing shared defaults.

For local Go runs, source `.env` before starting the service. The service no longer needs a config directory mount.

3 Run image from any directory:

```
./docker/server-connect/start.sh
```
