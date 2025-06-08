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

## CloudFlare

https://github.com/acmesh-official/acme.sh/wiki/dnsapi#dns_cf

```
export CF_Token="Y_jpG9AnfQmuX5Ss9M_qaNab6SQwme3HWXNDzRWs"
export CF_Zone_ID="763eac4f1bcebd8b5c95e9fc50d010b4"

acme.sh --issue --dns dns_cf -d r2f2.com -d '*.r2f2.com'
```

## Install the cert to Apache/Nginx etc.

```
acme.sh --install-cert -d r2f2.com \
--key-file       ~/balgass/docker/nginx-fail2ban/nginx/ssl/key.pem  \
--fullchain-file ~/balgass/docker/nginx-fail2ban/nginx/ssl/cert.pem \
--reloadcmd     "docker restart nginx-fail2ban mail"
```

## Start nginx

```
./start.sh
```
