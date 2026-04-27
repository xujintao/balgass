#!/usr/bin/env bash
set -euo pipefail

if [[ ! -d .git || ! -d config || ! -d src || ! -d docker ]]; then
  echo "Please run this script from the repository root." >&2
  exit 1
fi

files=(
  "config/server-connect/.env"
  "config/server-connect/IGC_ServerList.xml"
  "config/server-game-common/Data/CommonServer.cfg"
  "config/server-game-common/IGCData/IGC_ExpSystem.xml"
  "config/server-game/.env"
  "config/server-web/.env"
  "docker/pgadmin/start.sh"
  "docker/pgsql/start.sh"
  "src/c1c2/tcp.go"
)

for file in "${files[@]}"; do
  git update-index --skip-worktree -- "$file"
  echo "skip-worktree: $file"
done

echo "Done. Marked ${#files[@]} files as skip-worktree."
