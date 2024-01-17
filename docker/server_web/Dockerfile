FROM python:3.10-alpine3.17
LABEL maintainer "xujintao@126.com"
ENV SERVER_WEB=/server_web
ADD ./balgass/src/server_web ${SERVER_WEB}
WORKDIR ${SERVER_WEB}
RUN set -ex; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' \
    /etc/apk/repositories; \
    apk update; \
    apk add tzdata; \
    pip install -r ${SERVER_WEB}/requirements.txt
CMD ["sh","-c","python manage.py migrate && python manage.py runserver 0.0.0.0:8000"]
