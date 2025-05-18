## server_web

1 Build image

```
REPO=~/balgass
FILE=./docker/server_web/Dockerfile

cd $REPO

docker build \
-t xujintao/server_web:latest \
-f $FILE \
.
```

2 Run image

```
./start.sh
```
