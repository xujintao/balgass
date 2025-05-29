## System

<img src="system.jpg">

### hold ubuntu 22.04.4

```
linux-generic-hwe-22.04
 ├── linux-image-generic-hwe-22.04
 │    └── linux-image-6.5.0-18-generic
 └── linux-headers-generic-hwe-22.04
      └── linux-headers-6.5.0-18-generic
```

```
pi@k65:~$ sudo apt-mark hold linux-image-generic-hwe-22.04 linux-headers-generic-hwe-22.04 linux-generic-hwe-22.04
linux-image-generic-hwe-22.04 set on hold.
linux-headers-generic-hwe-22.04 set on hold.
linux-generic-hwe-22.04 set on hold.

pi@k65:~$ sudo apt-mark hold linux-image-6.5.0-18-generic linux-headers-6.5.0-18-generic
linux-image-6.5.0-18-generic set on hold.
linux-headers-6.5.0-18-generic set on hold.

pi@k65:~$ apt-mark showhold
linux-generic-hwe-22.04
linux-headers-6.5.0-18-generic
linux-headers-generic-hwe-22.04
linux-image-6.5.0-18-generic
linux-image-generic-hwe-22.04
```

### dummy interface

temporary

```
# add dummy interface
ip link add dummy0 type dummy
ip addr add 10.1.2.3/24 dev dummy0

# delete dummy interface
ip link delete dummy0
```

persistent

```
$ ls -l /etc/NetworkManager/system-connections/dummy0.nmconnection
-rw------- 1 root root 141 May 28 20:40 /etc/NetworkManager/system-connections/dummy0.nmconnection
```

```
$ sudo cat /etc/NetworkManager/system-connections/dummy0.nmconnection
[connection]
id=dummy0
type=dummy
interface-name=dummy0
autoconnect=true

[ipv4]
method=manual
addresses=10.1.2.3/24;

[ipv6]
method=ignore
```

```
sudo systemctl restart NetworkManager
```
