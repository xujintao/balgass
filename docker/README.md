## Deploy

<img src="deploy.drawio.svg">

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
