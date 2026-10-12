# Alloy logs to Grafana Cloud

Alloy collects Caddy and Fail2ban log files plus Docker logs from `pgsql`,
`server-game`, and `server-connect`. It sends them to the Grafana Cloud Loki
endpoint configured in `config.alloy`.

1. In the Grafana Cloud Portal, open your stack's Loki details and confirm that
   the endpoint URL and User ID match `config.alloy`. Create a Cloud Access
   Policy token with `logs:write` permission for that stack.
2. Copy `.env.example` to `.env` in this directory and replace the token
   placeholder. `.env` is ignored by Git; do not commit the token.
3. On the Docker host, run `./start.sh`. If an `alloy` container already exists,
   stop and remove that container before rerunning the script. Its `alloy_data`
   volume keeps file positions across container replacement.
4. Check `docker logs alloy` for authentication or delivery errors. In Grafana
   Cloud Explore, choose your stack's Loki data source and query `{job="caddy"}`
   and `{job="docker"}` for recent logs.

The logs are stored in Grafana Cloud Loki.
