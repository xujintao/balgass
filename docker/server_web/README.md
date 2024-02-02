## server_web

1 Build and push image to docker private registry

```
cd ~/r2f2/server_web
cp -r ~/github.com/xujintao/balgass .
docker build -t server_web .
```

2 Run image

```
docker run \
--restart always \
-d \
--name server_web \
-e TZ=Asia/Shanghai \
-p 8000:8000 \
server_web:latest
```
