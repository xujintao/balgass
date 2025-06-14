FROM golang:1.21.6-alpine3.18 as builder
COPY . /balgass
WORKDIR /balgass/src/server-game
RUN CGO_ENABLED=0 go install .

FROM alpine:3.18
LABEL maintainer "xujintao@126.com"
RUN set -ex; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' \
    /etc/apk/repositories; \
    apk update; \
    apk add tzdata
COPY --from=builder /go/bin/server-game /usr/local/bin/
CMD ["server-game"]
