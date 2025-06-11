#!/bin/sh

# equivalent to:
# docker run \
# --env-file "${HOME}/balgass/config/server_web/.env" \
# ...

set -a
. "${HOME}/balgass/config/server_web/.env"
set +a

exec python manage.py runserver
# exec ${HOME}/balgass/docker/server_web/entrypoint.sh
