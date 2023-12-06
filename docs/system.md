## System

<img src="system.jpg">

### dummy interface

```
# add dummy interface
ip link add dummy0 type dummy
ip addr add 10.1.2.3/24 dev dummy0
```

```
# delete dummy interface
ip link delete dummy0
```
