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

#### wireshark capture loopback
```
route add 192.168.0.100 mask 255.255.255.255 192.168.0.1
```

#### Display filter express
```
(tcp.port==44405 or tcp.port==56900) and !(icmp or tcp.analysis.retransmission or tcp.analysis.out_of_order or tcp.analysis.duplicate_ack)
```

#### Connect protocol
```
code:00   response connection result

code:f406 server list, 第一次请求服务器列表其实是为了得到news，第二次请求服务器列表才会真正解析并渲染
code:fa00 news title
code:fa01 news

code:f33d end

code:f403 server ip and port
// client reset CS and connect to GS
```

```
code:f33c keep alive 27s, 触发不稳定，实际上是f331的xor，发向connect服务器的消息到底要不要xor?
```

#### Game protocol
```
code:f100 response connection result with server's version
code:f101 login
code:f300 character list

code:fa11 anti hack 20s, begin at f100
code:fa15 hit hack 3s
code:0e   keep alive 20s, begin at f101

code:f315 check character
code:f303 character map
code:fa0d file crc
code:f326 popup
code:f61a progress quest list
code:f621 event quest list
code:4e11 ride select
code:18   action
code:40   party
```

```
code:f331 trap message
code:7399 invalid packet
code:d725 invalid packet
code:d74b invalid packet [c1 07 d7 4b bb 72 77]
[c3 0b 25 cf 0f 39 e7 e9 66 60 01], invalid packet

code:5d3f invalid packet [c1 0d 5d 3f 6b d6 08 3f e1 b0 06 12 82]
code:cf07 invalid packet [c1 0d cf 07 49 aa 5c 56 d9 46 01 df e0]
code:ab59 invalid packet [c1 0d ab 59 d3 59 7d c9 6f ee 26 21 86]
lost anti-hack packet, GetTickCount problem
```

#### End