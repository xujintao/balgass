# Caddy + Fail2ban

This is the VPS entry point. It uses Docker's existing default
bridge and the existing published backend ports. It does not change host
iptables rules or create a Docker network.

## Before starting

1. Point `game.r2f2.com` and `patch.game.r2f2.com` directly to the VPS. Keep
   ports 80 and 443 available for Caddy's certificate challenges and HTTPS.
2. Configure the same `GAME_API_TOKEN` (at least 32 random characters) in
   the game service and Vercel. Configure Vercel with:

   ```text
   GAME_API_URL=https://game.r2f2.com/api/command
   GAME_WEBSOCKET_URL=wss://game.r2f2.com/api/game
   GAME_CONFIG_URL=https://game.r2f2.com/config/
   GAME_CONFIG_BASIC_AUTH_USER=nextjs
   GAME_CONFIG_BASIC_AUTH_PASSWORD=<separate random password>
   ```

3. Run `docker run --rm -it caddy:2-alpine caddy hash-password` and enter the
   **same** map-config password. Copy `.env.example` to `.env` and put the
   resulting hash in `GAME_CONFIG_AUTH_HASH`. The `.env` file is gitignored.
4. Set `HostURL = patch.game.r2f2.com` in
   `~/balgass/config/server-connect/IGCCS.ini` and restart server-connect.
   Its updater configuration documents HTTP/FTP support, so this Caddy setup
   serves the patch host over HTTP without redirecting it to HTTPS. HTTPS is
   also available for browser access.

## Build and run

From this directory, run `./start.sh`. It builds the image, validates the
Caddyfile, and starts the container. The persistent `data/`, `config/`, and
`logs/` directories are gitignored. Back up `data/` to preserve Caddy's ACME
account and certificates.

The Caddy container uses `host.docker.internal` mapped to Docker's host
gateway to reach the existing published ports 8080, 3000, and 8084. These
ports remain protected by `iptables-docker.sh`. The game config directory is
mounted read-only; only `IGC_MapList.xml` is currently in the `@gameConfigFiles`
path allowlist in `Caddyfile`. When a feature needs another IGCData file, add
its exact `/config/...` path to that matcher. It then uses the same Basic Auth,
read-only mount, and `/config` prefix removal; other files still return 404.
The old items page also used `Skills/IGC_SkillList.xml`,
`Items/IGC_ItemList.xml`, `Items/IGC_ItemSetType.xml`, and
`Items/IGC_ItemSetOption.xml`; add them only when the new site needs them.
The patch files are kept in `patch/`.

Grafana must be restarted with its updated `GF_SERVER_ROOT_URL` in
`grafana/.env`.
Alloy reads Caddy and Fail2ban logs through a read-only mount of `logs/` and
sends them to Loki through the Docker host gateway.
No Docker network or host firewall change is needed.

## Checks on the VPS

1. Check `docker logs caddy-fail2ban` and
   `docker exec caddy-fail2ban fail2ban-client status`.
2. Check `https://game.r2f2.com/config/IGC_MapList.xml` returns 401 without
   credentials and the XML with the configured Basic Auth credentials. Other
   files under `/config/` must return 404.
3. Check an authenticated `/api/command` request and a browser WebSocket at
   `wss://game.r2f2.com/api/game`; the game service must still reject requests
   without a valid API token.
4. Check `/grafana/` and `/pgadmin/` load their own login pages, including
   assets and Grafana Live. Check `http://patch.game.r2f2.com/version.wvd`
   returns the copied file without a redirect.
5. Check Caddy access logs appear in Loki. From a disposable client IP,
   trigger five probe requests, confirm `fail2ban-client status caddy-probes`
   reports the ban and that Caddy can no longer be reached from that IP. The
   host's `DOCKER-USER` rules should remain unchanged.

Fail2ban's nftables action runs in this container's network namespace. It
protects traffic reaching Caddy, while `iptables-docker.sh` continues to
protect the other published Docker ports.
