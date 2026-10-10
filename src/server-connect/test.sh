#!/usr/bin/env bash
set -euo pipefail

REPO_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)
set -a
. "$REPO_DIR/docker/server-connect/.env"
set +a

go test ./...
