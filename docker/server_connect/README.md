## server_connect

1 Build image

```
cd ~/r2f2/server_connect
cp -r ~/github.com/xujintao/balgass .
docker build -t server_connect .
```

2 Run image

```
docker run \
--restart always \
-d \
--name server_connect \
-e TZ=Asia/Shanghai \
-e CONFIG_PATH=/etc/server_connect \
-v ~/r2f2/config/server_connect:/etc/server_connect \
-p 44405:44405 \
-p 55667:55667/udp \
server_connect:latest
```
