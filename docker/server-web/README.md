## server-web

1 Build image

```
REPO=~/balgass
FILE=./docker/server-web/Dockerfile

cd $REPO

docker build \
-t xujintao/server-web:latest \
-f $FILE \
.
```

2 Run image

```
./start.sh
```
