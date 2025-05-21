## Install acme.sh online

https://github.com/acmesh-official/acme.sh

```
curl https://get.acme.sh | sh -s email=my@example.com
```

## Use Aliyun domain API to automatically issue cert

https://github.com/acmesh-official/acme.sh/wiki/dnsapi#11-use-aliyun-domain-api-to-automatically-issue-cert  
https://f-e-d.club/topic/use-acme-sh-deployment-let-s-encrypt-by-ali-cloud-dns-generic-domain-https-authentication.article

```
export Ali_Key="LTAI5tAaH1qEmJUREcktWDQg"
export Ali_Secret="JYTOpXRkE88igrSTnhIc1sz2gE2JYj"

acme.sh --issue --dns dns_ali -d r2f2.com -d *.r2f2.com
```

## Install the cert to Apache/Nginx etc.

```
acme.sh --install-cert -d r2f2.com \
--key-file       ~/balgass/docker/nginx/ssl/key.pem  \
--fullchain-file ~/balgass/docker/nginx/ssl/cert.pem \
--reloadcmd     "docker restart nginx mailserver"
```

## Start nginx

```
./start.sh
```
