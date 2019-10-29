#### Frame Format
```
<flag + len + code + [data]> +
<flag + len + code + [data]> +
...

note:
flag=c1, single byte len
flag=c3, aes encrypt

flag=c2, double bytes len, big endian
flag=c4, aes encrypt
```

#### Display filter express
```
(tcp.port==44405 or tcp.port==56900) and !(icmp or tcp.analysis.retransmission or tcp.analysis.out_of_order or tcp.analysis.duplicate_ack)
```

#### Connect protocol
```
code:00   response connection result

code:f33c keep alive 27s ???

code:f406 server list
code:fa00 news title
code:fa01 news
code:f403 server ip and port
```

client reset CS and connect to GS

#### Game protocol
```
code:f100 response connection result with server's version

code:fa11 anti hack 20s
code:fa15 hit hack 3s
code:0e   keep alive
code:f331 trap message ???

code:f101 login
code:f300 character list
code:7399 invalid packet ???

code:f315 check character
code:f303 character map
code:fa0d file crc
code:f326 popup
code:f61a progress quest list
code:f621 event quest list
code:4e11 ride select
code:18   action
code:40   party
code:d725 invalid packet
```

#### End