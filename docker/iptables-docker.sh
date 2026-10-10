#!/bin/bash

while ! iptables -L DOCKER-USER -n &>/dev/null; do
    sleep 1
done

# pgsql
sudo iptables -I DOCKER-USER 1 -s 172.17.0.0/16 -p tcp --dport 5432 -j ACCEPT
sudo iptables -I DOCKER-USER 2 -p tcp --dport 5432 -j DROP

# pgadmin
sudo iptables -I DOCKER-USER 3 -s 172.17.0.0/16 -p tcp --dport 8084 -j ACCEPT
sudo iptables -I DOCKER-USER 4 -p tcp --dport 8084 -j DROP

# server-game
sudo iptables -I DOCKER-USER 5 -s 172.17.0.0/16 -p tcp --dport 8080 -j ACCEPT
sudo iptables -I DOCKER-USER 6 -p tcp --dport 8080 -j DROP

# server-connect
sudo iptables -I DOCKER-USER 7 -s 172.17.0.0/16 -p udp --dport 55667 -j ACCEPT
sudo iptables -I DOCKER-USER 8 -p udp --dport 55667 -j DROP

# loki
sudo iptables -I DOCKER-USER 9 -s 172.17.0.0/16 -p tcp --dport 3100 -j ACCEPT
sudo iptables -I DOCKER-USER 10 -p tcp --dport 3100 -j DROP

# grafana
sudo iptables -I DOCKER-USER 11 -s 172.17.0.0/16 -p tcp --dport 3000 -j ACCEPT
sudo iptables -I DOCKER-USER 12 -p tcp --dport 3000 -j DROP
