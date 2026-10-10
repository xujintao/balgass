## Deploy

<img src="deploy.drawio.svg">

## Container environment files

Before starting Grafana, Loki, pgAdmin, PostgreSQL, or Promtail, copy its
`docker/<container>/.env.example` to `.env` in the same directory. Fill the
blank credentials in the Grafana, pgAdmin, and PostgreSQL files with your
existing credentials for persistent installations, or new strong credentials
for new installations. The `.env` files are gitignored and read by each
container's `start.sh`. Loki and Promtail only need the copied timezone file;
Promtail's log collection settings remain in `promtail/config.yml`.

## ~~Use UFW~~

### 1. iptables

```
# add config to /etc/docker/daemon.json and systemctl restart docker
{
    "iptables": false
}
./iptables-ufw.sh
./ufw.sh
```

### 2. make it persistent

```
sudo cp iptables-ufw-boot.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable iptables-ufw-boot.service
```

## Use Docker iptables

```
    ┌────────────────────────────────────────────────────────────────────┐
    │                                                                    │
    │                                                                    │
    │  ┌─────────────► PREROUTING                              INPUT ────┼─────────┐
    │  │                   │                                     ▲       │         │
    │  │                   │                                     │       │         │
    │  │                   │                                     │       │         │
    │  │                   │            ┌─────────┐              │       │         │
    │  │                   └──────────► │ Routing │ ─────────────┘       │         |
    │                                   │         │                      │         ▼
interfaces                 ┌─────────── │ decision│ ◄────────────┐       │    Local Process
 ▲  │                      │            └────┬────┘              │       │         │
 │  │                      │                 │                   │       │         │
 │  │                      │                 │                   │       │         │
 │  │                      │                 │                   │       │         │
 │  │                      │                 ▼                   │       │         │
 │  │                      │              FORWARD                │       │         │
 │  │                      │                 │                   │       │         │
 │  │                      ▼                 │                           │         │
 └──┼───────────────── POSTROUTING ◄─────────┘                 OUTPUT ◄──┼─────────┘
    │                                                                    │
    │                                                                    │
    └────────────────────────────────────────────────────────────────────┘
```

### 1. iptables

```
./iptables-docker.sh
```

### 2. make it persistent

```
sudo cp iptables-docker-boot.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable iptables-docker-boot.service
```
