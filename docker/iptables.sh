#!/bin/bash

while ! iptables -L DOCKER-USER -n &>/dev/null; do
    sleep 1
done

sudo iptables -I DOCKER-USER 1 -s 172.17.0.0/16 -p tcp --dport 5432 -j ACCEPT
sudo iptables -I DOCKER-USER 2 -p tcp --dport 5432 -j DROP
sudo iptables -I DOCKER-USER 3 -s 172.17.0.0/16 -p tcp --dport 8084 -j ACCEPT
sudo iptables -I DOCKER-USER 4 -p tcp --dport 8084 -j DROP
sudo iptables -I DOCKER-USER 5 -s 172.17.0.0/16 -p tcp --dport 8000 -j ACCEPT
sudo iptables -I DOCKER-USER 6 -p tcp --dport 8000 -j DROP
sudo iptables -I DOCKER-USER 7 -s 172.17.0.0/16 -p tcp --dport 8080 -j ACCEPT
sudo iptables -I DOCKER-USER 8 -p tcp --dport 8080 -j DROP
sudo iptables -I DOCKER-USER 9 -s 172.17.0.0/16 -p udp --dport 55667 -j ACCEPT
sudo iptables -I DOCKER-USER 10 -p udp --dport 55667 -j DROP