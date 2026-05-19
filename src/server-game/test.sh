#!/bin/sh

set -a
. "${HOME}/balgass/config/server-game/.env"
set +a

go test ./...
