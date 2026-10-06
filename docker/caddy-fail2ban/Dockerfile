FROM caddy:2-alpine

RUN apk add --no-cache bash fail2ban nftables tini
RUN rm -f /etc/fail2ban/jail.d/*

COPY entrypoint.sh /usr/local/bin/caddy-fail2ban-entrypoint
RUN chmod +x /usr/local/bin/caddy-fail2ban-entrypoint \
    && mkdir -p /var/log/caddy /var/lib/fail2ban

ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/caddy-fail2ban-entrypoint"]
