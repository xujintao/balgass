## server_game

1 Build image

```
cd ~/r2f2/server_game
cp -r ~/github.com/xujintao/balgass .
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
-e COMMON_PATH=/etc/server_game_common \
-v ~/r2f2/config/server_game:/etc/server_game \
-v ~/r2f2/config/server_game_common:/etc/server_game_common \
-p 8080:8080 \
-p 56900:56900 \
server_game:latest
```
