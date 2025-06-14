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

2 Run image

```
./start.sh
```
