## server_game

1 Build image

```
REPO=~/balgass
FILE=./docker/server_game/Dockerfile

cd $REPO

docker build \
-t xujintao/server_game:latest \
-f $FILE \
.
```

2 Run image

```
./start.sh
```
