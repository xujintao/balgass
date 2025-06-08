ufw reset

# 1. ssh
ufw allow 22/tcp comment 'SSH'

# 2. nginx
ufw allow 80/tcp comment 'HTTP'
ufw allow 443/tcp comment 'HTTPS'

# 3. mail
ufw allow 25/tcp comment 'SMTP'
ufw allow 465/tcp comment 'SMTPS'
ufw allow 587/tcp comment 'SMTP-Auth'
ufw allow 143/tcp comment 'IMAP'
ufw allow 993/tcp comment 'IMAPS'

# 4. pgsql
ufw allow from 172.17.0.0/16 to any port 5432 proto tcp comment 'pgsql <- docker'

# 5. pgadmin
ufw allow from 172.17.0.0/16 to any port 8084 proto tcp comment 'pgadmin <- docker'

# 6. server_web
ufw allow from 172.17.0.0/16 to any port 8000 proto tcp comment 'server_web <- docker'

# 7. server_game
ufw allow from 172.17.0.0/16 to any port 8080 proto tcp comment 'server_game <- docker'
ufw allow 56900/tcp comment 'server_game'

# 8. server_connect
ufw allow 44405/tcp comment 'server_connect'
ufw allow from 172.17.0.0/16 to any port 55667 proto udp comment 'server_connect <- docker'

# 9. promtail
ufw allow from 172.17.0.0/16 to any port 9080 proto tcp comment 'promtail <- docker'

# 10. loki
ufw allow from 172.17.0.0/16 to any port 3100 proto tcp comment 'loki <- docker'

# 11. grafana
ufw allow from 172.17.0.0/16 to any port 3000 proto tcp comment 'grafana <- docker'

ufw enable
ufw status numbered
