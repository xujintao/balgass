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

2 Run image

```
./start.sh
```
