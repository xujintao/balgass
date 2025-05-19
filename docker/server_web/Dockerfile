FROM python:3.10-slim-bookworm
LABEL maintainer "xujintao@126.com"
ENV SERVER_WEB=/server_web
COPY ./src/server_web ${SERVER_WEB}
WORKDIR ${SERVER_WEB}
RUN set -ex; \
    apt-get update; \
    apt-get install -y tzdata; \
    apt-get install -y gcc libpq-dev; \
    rm -rf /var/lib/apt/lists/*; \
    pip install --no-cache-dir -r ${SERVER_WEB}/requirements.txt
CMD ["sh","-c","python manage.py migrate && python manage.py runserver 0.0.0.0:8000"]
