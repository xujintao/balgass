FROM golang:1.21.6-alpine3.18 as builder
ADD ./balgass /balgass
WORKDIR /balgass/src/server_game
RUN CGO_ENABLED=0 go install .

FROM alpine:3.18
LABEL maintainer "xujintao@126.com"
RUN set -ex; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' \
    /etc/apk/repositories; \
    apk update; \
    apk add tzdata
COPY --from=builder /go/bin/server_game /user/local/bin
CMD ["server_game"]