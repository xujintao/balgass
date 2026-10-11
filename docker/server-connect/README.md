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

All service settings live in `.env`. Set the update host and the server IPs in `SERVER_LIST_JSON` for your deployment. The JSON value is wrapped in single quotes so the same file works when sourced by `test.sh`; the service also accepts the enclosing quotes passed by Docker's `--env-file`. Keep `.env` private and edit `.env.example` when changing shared defaults.

For local Go runs, source `.env` before starting the service. The service no longer needs a config directory mount.

3 Run image from any directory:

```
./docker/server-connect/start.sh
```
