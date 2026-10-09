ufw reset

# 1. ssh
ufw allow 22/tcp comment 'SSH'

# 2. caddy
ufw allow 80/tcp comment 'HTTP'
ufw allow 443/tcp comment 'HTTPS'

# 3. pgsql
ufw allow from 172.17.0.0/16 to any port 5432 proto tcp comment 'pgsql <- docker'

# 4. pgadmin
ufw allow from 172.17.0.0/16 to any port 8084 proto tcp comment 'pgadmin <- docker'

# 5. server-game
ufw allow from 172.17.0.0/16 to any port 8080 proto tcp comment 'server-game <- docker'
ufw allow 56900/tcp comment 'server-game'

# 6. server-connect
ufw allow 44405/tcp comment 'server-connect'
ufw allow from 172.17.0.0/16 to any port 55667 proto udp comment 'server-connect <- docker'

# 7. promtail
ufw allow from 172.17.0.0/16 to any port 9080 proto tcp comment 'promtail <- docker'

# 8. loki
ufw allow from 172.17.0.0/16 to any port 3100 proto tcp comment 'loki <- docker'

# 9. grafana
ufw allow from 172.17.0.0/16 to any port 3000 proto tcp comment 'grafana <- docker'

ufw enable
ufw status numbered
