## Deploy

<img src="deploy.jpg">

## ~~UFW~~

```
# 1. reset
ufw reset

# 2. default
ufw default deny incoming
ufw default allow outgoing

# 3. ssh
ufw allow 22/tcp comment 'SSH'

# 4. nginx
ufw allow 80/tcp comment 'HTTP'
ufw allow 443/tcp comment 'HTTPS'

# 5. mail
ufw allow 25/tcp comment 'SMTP'
ufw allow 465/tcp comment 'SMTPS'
ufw allow 587/tcp comment 'SMTP-Auth'
ufw allow 143/tcp comment 'IMAP'
ufw allow 993/tcp comment 'IMAPS'

# 6. pgsql
ufw allow from 172.17.0.0/16 to any port 5432 proto tcp comment 'pgsql <- docker'

# 7. pgadmin
ufw allow 8084/tcp comment 'pgadmin'

# 8. server_web
ufw allow from 172.17.0.0/16 to any port 8000 proto tcp comment 'server_web <- docker'

# 9. server_game
ufw allow from 172.17.0.0/16 to any port 8080 proto tcp comment 'server_game <- docker'
ufw allow 56900/tcp comment 'server_game'

# 10. server_connect
ufw allow 44405/tcp comment 'server_connect'
ufw allow from 172.17.0.0/16 to any port 55667 proto udp comment 'server_connect <- docker'

# 11. enbale
ufw enable

# 12. status
ufw status numbered
```

## Use iptables

```
sudo iptables -L -n -v
./iptables.sh
sudo iptables -L -n -v
```

## Add iptables.sh to boot

```
sudo cp iptables-boot.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable iptables-boot.service
sudo systemctl start iptables-boot.service
```
