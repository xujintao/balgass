#!/bin/sh

# equivalent to:
# docker run \
# --env-file "${HOME}/balgass/config/server-web/.env" \
# ...

set -a
. "${HOME}/balgass/config/server-web/.env"
set +a

exec python manage.py runserver 0.0.0.0:8000
# exec ${HOME}/balgass/docker/server-web/entrypoint.sh
