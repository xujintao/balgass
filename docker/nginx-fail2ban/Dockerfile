FROM ubuntu:latest

RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx \
    fail2ban \
    python3-pyinotify \
    python3-dnspython \
    curl \
    ca-certificates \
    && rm -rf /etc/fail2ban/jail.d/* \
    && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["sh", "-c"]
CMD ["/usr/bin/fail2ban-server -xf start & nginx -g 'daemon off;'"]
