#!/bin/bash
set -euo pipefail

CADDY_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
IMAGE=xujintao/caddy-fail2ban:latest
ENV_FILE="$CADDY_DIR/.env"

if [[ ! -f "$ENV_FILE" ]] || ! grep -q '^GAME_CONFIG_AUTH_HASH=.' "$ENV_FILE"; then
    echo "Create $ENV_FILE with a nonempty GAME_CONFIG_AUTH_HASH before starting Caddy." >&2
    exit 1
fi

mkdir -p "$CADDY_DIR/data" "$CADDY_DIR/config" "$CADDY_DIR/logs"
docker build -t "$IMAGE" "$CADDY_DIR"
docker run --rm \
    --entrypoint caddy \
    --env-file "$ENV_FILE" \
    -v "$CADDY_DIR/Caddyfile:/etc/caddy/Caddyfile:ro" \
    "$IMAGE" validate --config /etc/caddy/Caddyfile --adapter caddyfile

docker run \
    --restart always \
    -d \
    --name caddy-fail2ban \
    --add-host host.docker.internal:host-gateway \
    --cap-add NET_ADMIN \
    -e LANG=C.UTF-8 \
    --env-file "$ENV_FILE" \
    -v "$CADDY_DIR/Caddyfile:/etc/caddy/Caddyfile:ro" \
    -v "$CADDY_DIR/fail2ban/fail2ban.local:/etc/fail2ban/fail2ban.local:ro" \
    -v "$CADDY_DIR/fail2ban/jail.local:/etc/fail2ban/jail.local:ro" \
    -v "$CADDY_DIR/fail2ban/caddy-probes.conf:/etc/fail2ban/filter.d/caddy-probes.conf:ro" \
    -v "$CADDY_DIR/fail2ban/caddy-admin-auth.conf:/etc/fail2ban/filter.d/caddy-admin-auth.conf:ro" \
    -v "$CADDY_DIR/data:/data" \
    -v "$CADDY_DIR/config:/config" \
    -v "$CADDY_DIR/logs:/var/log/caddy" \
    -v "$CADDY_DIR/patch:/srv/patch:ro" \
    -v "$HOME/balgass/config/server-game-common/IGCData:/srv/game-config:ro" \
    -p 80:80 \
    -p 443:443 \
    "$IMAGE"
