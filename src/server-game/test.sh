#!/bin/sh

set -a
. "${HOME}/balgass/config/server-game/.env"
set +a

# Add each completed system's stable tests here.
go test ./game/bot
go test -race ./game/bot
go test ./game/fixture
