<img src="game.jpg">

## deploy
1 Build image
```
# CGO_ENABLED=0 go build .
# scp server_game host
# ssh host build image
docker build -t server_game .
```

2 Run image
```
docker run \
--restart always \
-d \
--name server_game \
-e TZ=Asia/Shanghai \
-e CONFIG_PATH=/etc/server_game \
-e COMMON_PATH=/etc/server_game/IGCData \
-v ~/server_game/config:/etc/server_game \
-p 8080:8080 \
server_game
```