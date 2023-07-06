# py-web

1, Build and push image to docker private registry

```
# git clone
git clone git@github.com:xujintao/py-web.git ~/github.com/xujintao/py-web
cd ~/github.com/xujintao/py-web

# docker build
docker build -t py-web:0.1 .
```

2, Run image

```
docker run \
--restart always \
-d \
--name py-web \
-e TZ=Asia/Shanghai \
-p 8000:8000 \
py-web:0.1
```