FROM alpine:3.17
LABEL maintainer "xujintao@51aes.com"
RUN set -ex; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' \
    /etc/apk/repositories; \
    apk update; \
    apk add tzdata
ADD server_game /usr/local/bin/
CMD ["server_game"]