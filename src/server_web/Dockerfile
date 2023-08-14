FROM python:3.10-alpine3.17
LABEL maintainer "xujintao@126.com"
ENV PY_WEB=/py-web
ADD . ${PY_WEB}
RUN set -ex; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' \
    /etc/apk/repositories; \
    apk update; \
    apk add tzdata; \
    pip install -r ${PY_WEB}/requirements.txt
WORKDIR ${PY_WEB}
CMD ["sh","-c","python manage.py migrate && python manage.py runserver 0.0.0.0:8000"]
