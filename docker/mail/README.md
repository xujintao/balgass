## 1, docker-mailserver

https://docker-mailserver.github.io/docker-mailserver/v15.0/

### DNS

https://docker-mailserver.github.io/docker-mailserver/v15.0/usage/#minimal-dns-setup
https://docker-mailserver.github.io/docker-mailserver/v15.0/config/best-practices/dkim_dmarc_spf/

### TLS

https://docker-mailserver.github.io/docker-mailserver/v15.0/usage/#setting-up-tls

### Fail2Ban

https://docker-mailserver.github.io/docker-mailserver/v15.0/config/security/fail2ban/

```
Fail2Ban: /etc/fail2ban/jail.conf
DMS:      /etc/fail2ban/jail.local
User:     /ect/fail2ban/jail.d/*.local <- docker-data/dms/config/fail2ban-jail.cf
```

## 2, Start

```
./start.sh
```

## 3, Add User

```
./adduser.sh
```
