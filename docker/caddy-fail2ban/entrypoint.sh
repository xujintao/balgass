#!/bin/bash
set -euo pipefail

touch /var/log/caddy/access.log /var/log/caddy/patch-http.log \
    /var/log/caddy/patch-https.log /var/log/caddy/fail2ban.log
chmod 0644 /var/log/caddy/access.log /var/log/caddy/patch-http.log \
    /var/log/caddy/patch-https.log /var/log/caddy/fail2ban.log

fail2ban-server -f -c /etc/fail2ban -x start &
fail2ban_pid=$!
caddy_pid=

stop() {
    for pid in "$fail2ban_pid" "$caddy_pid"; do
        if [[ -n "$pid" ]]; then
            kill "$pid" 2>/dev/null || true
            wait "$pid" 2>/dev/null || true
        fi
    done
}
trap stop EXIT
trap 'exit 0' TERM INT

for attempt in {1..30}; do
    if fail2ban-client ping >/dev/null 2>&1; then
        break
    fi
    kill -0 "$fail2ban_pid"
    sleep 0.2
done
fail2ban-client ping >/dev/null

caddy run --config /etc/caddy/Caddyfile --adapter caddyfile &
caddy_pid=$!

wait -n "$fail2ban_pid" "$caddy_pid"
exit 1
