#!/bin/sh

set -e
REPO_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
set -a
. "$REPO_DIR/docker/server-game/.env"
set +a

# Add each completed system's stable tests here.
go test ./game/event/... ./game/drop ./game
go test -race ./game/event/...
go test ./game/bot
go test -race ./game/bot
go test ./game/fixture
go test ./game/object
go test ./game/object/player
go test ./game/object/monster
go test ./game/skill
go test ./game/model
go test ./game/maps
go test ./game/formula
go test ./handle
