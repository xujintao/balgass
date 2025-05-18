## server_connect

1 Build image

```
REPO=~/balgass
FILE=./docker/server_connect/Dockerfile

cd $REPO

docker build \
-t xujintao/server_connect:latest \
-f $FILE \
.
```

2 Run image

```
./start.sh
```
