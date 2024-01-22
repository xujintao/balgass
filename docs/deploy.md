## System

<img src="deploy.jpg">

### 1. apply ufw rules

```
ufw allow ssh
ufw allow http
ufw allow https
ufw enable
ufw status verbose
```

Failed: server_web container can't access pgsql container

https://stackoverflow.com/questions/30383845/what-is-the-best-practice-of-docker-ufw-under-ubuntu

#### 1.1 ufw-docker

```
ufw-docker install
delete 192.168.0.0/16
changed 10.10.10.0/8 to 10.10.10.0/24
changed 172.16.0.0/12 to 172.17.0.0/16

systemctl restart ufw
```

https://stackoverflow.com/a/51741599
https://github.com/chaifeng/ufw-docker/issues/112

Failed: server_web container can't access pgsql container

#### 1.2 --iptables=false

https://stackoverflow.com/a/46266757

Failed: server_web container can't access pgsql container
