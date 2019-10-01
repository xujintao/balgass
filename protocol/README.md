#### Frame Format
```
<flag + len + code + [data]> +
<flag + len + code + [data]> + ...

note:
flag=c1, single byte len
flag=c3, data has been encoded

flag=c2, double bytes len, big endian
flag=c4, data has been encoded, big endian
```

#### Connect protocol
express
```
ip.addr==192.168.0.100 and (tcp.port==44405 or tcp.port==56900)
```

这是什么？
```
client <- server(4):
c1 04 00 01
<flag:0xc1 + len:4 + type:0x0001>

// 请求
client(4) -> server:
c1 04 f4 06
<flag:0xc1 + len:4 + type:0xf406>

// 响应
client <- server(204):
c1 cc fa 00 "Server News" 00  
cc cc cc cc 0f 10 08 e8 28 f3 f4 09 66 2a 32 01  
01 00 00 00 04 f4 f4 09 60 f5 f4 09 80 e3 4b 01  
cc cc cc ...
<flag:0xc1 + len:204 + type:0xfa00 + data:...>

// 响应
client <- server(236)
c2 00 63 fa 01 0a 0c ...
<flag:0xc2 + len:99(0x0063) + type:0xfa01 + data:0a 0c dd 07 ff 00 9b 00 d7 00 64 00 78 00 87 00 "Info" ... 00 24 00 "This news info support all languages"> +
<flag:0xc2 + len:126(0x007e) + type:0xfa01 + data:0a 0c dd 07 ff 00 9b 00 d7 00 64 00 78 00 87 00 "Info" .. 00 3f 00 "When no News added News Windows in CLient will not be displayed"> +
<flag:0xc2 + len:11(0x000b) + type:0xf406 + data: 00 01 00 00 00 00>
```

服务器列表
```
// 请求
client(4) -> server:
c1 04 f4 06
<flag:0xc1 + len:4 + type:0xf406>

// 响应
client <- server(11):
c2 00 0b f4 06 00 01 00 00 00 00
<flag:0xc2 + len:11(0x000b) + type:0xf406 + data:00 01 00 00 00 00>
```

code为0的服务器ip和端口
```
// 请求
client(6) -> server:
c1 06 f4 03 00 00
<flag:0xc1 + len:6 + type:0xf403 + data:00 00>

// 响应
client <- server(22):
c1 16 f4 03 192.168.0.100 00 00 00 44 de
<flag:0xc1 + len:22 + type:0xf403 + data:"192.168.0.100" 00 00 00 0xde44> // port是小端模式
```
客户端RST掉CS连接并和GS建立连接

#### Game protocol
express
```
ip.addr==192.168.0.100 and tcp.port==56900 and !(icmp or tcp.analysis.retransmission or tcp.analysis.out_of_order or tcp.analysis.duplicate_ack)
```

heartbeat
```
client(35) -> server:
c3 23 80 04 81 34 b4 d1 af be 41 e6 dc f8 4a 22
f5 d6 e7 b4 c7 e5 d5 6d 5a 86 e8 0c a9 15 77 fe
a3 da 0b
<flag:0xc3 + len:35 + type:0x8004 + data:...>

client(9) -> server:
c1 09 fa 15 86 5a ca e2 c4
<flag:0xc1 + len:9 + type:0xfa15 + data:86 5a ca e2 c4>
```

version
```
client <- server(12):
c1 0c f1 00 01 2e 7e "10444"
<flag:0xc1 + len:12 + type:0xf100 + data:01 2e 7e "10444">

客户端判断服务器版本，匹配就继续，不匹配就RST连接并退出程序  
(想办法禁用客户端的版本匹配)
```

unknown
```
// 20秒一次，也是心跳？
client(19) -> server:
c3 13 d4 4e 83 1e 01 45 f9 e3 c5 fb ea 32 6d 04
3b ea 05
<flag:0xc3 + len:19 + type:0xd44e + data:83 1e 01 45 f9 e3 c5 fb ea 32 6d 04 3b ea 05>

// 服务器响应的invalid packet received OP:5
client <- server(104):
c1 68 fa a5 c9 b0 ...
<flag:0xc1 + len:104 + type:0xfaa5 + data:...>
```

login
```
// 请求
client(179) -> server: 应该是username和passwd的加密信息
c3 b3 92 af ... cb 0f
<flag:0xc3 + len:179 + type:0x92af + data:...>

// 响应
client <- server(5): 登录成功
c1 05 f1 01 01
<flag:0xc1 + len:5 + type:0xf101 + data:0x01>
```

character
```
// 请求
client(4) -> server:
c1 04 f3 0d
<code:0xf30d>

// 响应
client <- server(49):
c1 31 0d 01 00 00 00 00 cd 00 00 00 00 "You are logged in as Free Player" 00 00 00 f8
<flag:0xc1 + len:49 + code:0x0d01 + data:..."You are logged in as Free Player"...>
// 响应
client <- server(5):
c1 05 de 00 1f
<flag:0xc1 + len:5 + code:0xde00 + data:0x1f>
// 响应
client <- server(662): 角色信息
c1 0e fa 0a ...
<flag:0xc1 + len:0x0e(14) + code:0xfa0a + data:00...> +
<flag:0xc1 + len:0x74(116) + code:0xf300 + data:..."juju"..."milf"..."Mature"...> +
<flag:0xc2 + len:0x01c2(450) + code:... + data:...> +
<flag:0xc1 + len:0x1f(31) + code:0xfa12 + data:...>
<flag:0xc1 + len:0x31(49) + code:0x0d01 + data:..."You are logged in as Free Player"...>
```

character scene  
character select  
main scene  
```
// 请求
client(14) -> server:
c1 0e f3 18 6a 3c 93 45 8f bc 7d b1 d7 b0
<flag:0xc1 + len:0x0e(14) + code:0xf318 + data:6a 3c 93 45 8f bc 7d b1 d7 b0>

// 响应
client <- server(15):
c1 0f f3 15 00 00 00 00 00 00 00 00 00 00 00
<flag:0xc1 + len:0x0f(15) + code:0xf315 + data:...>

// 请求
client(6) -> server:
c1 06 4e a1 46 9a
<flag:0xc1 + len:6 + code:0x4ea1 + data:9a>

// 请求
client(15) -> server:
c1 0f f3 0e 7c 2a 85 53 99 aa 6b a7 c1 a6 87
<flag:0xc1 + len:0x0f(15) + code:0xf30e + data:7c 2a 85 53 99 aa 6b a7 c1 a6 87>

// 响应
client <- server(16):
c1 10 26 fe 04 c8 00 33 e7 00 00 00 c8 04 00 00
<flag:0xc1 + len:0x10(16) + code:0x26fe + data:04 c8 00 33 e7 00 00 00 c8 04 00 00>
// 响应
client <- server(1460):
c1 0c 27 fe 02 66 02 8b 66 02 00 00
<flag:0xc1 + len:0x0c(12) + code:0x27fe + data:02 66 02 8b 66 02 00 00>
<flag:0xc1 + len:0x22(34) + code:0x5900 + data:01 00 44 00 66 00 04 00 01 00 ... 00>
<flag:0xc1 + len:0x0c(12) + code:0xec30 + data:0f 00 00 00 0b 00 00 00>

<flag:0xc1 + len:0x10(16) + code:0x26fe + data:04 c8 00 33 e7 00 00 00 c8 04 00 00>
<flag:0xc1 + len:0x0c(12) + code:0x27fe + data:02 66 02 8b 66 02 00 00>
<flag:0xc1 + len:0x22(34) + code:0x5900 + data:01 00 44 00 66 00 04 00 01 00 ... 00>
<flag:0xc1 + len:0x0c(12) + code:0xec30 + data:0f 00 00 00 0b 00 00 00>

...

<flag:0xc1 + len:0x10(16) + code:0x26fe + data:05 8e 00 35 f4 00 00 00 8e 05 00 00>
<flag:0xc1 + len:0x0c(12) + code:0x27fe + data:02 b0 03 02 b0 02 00 00>
<flag:0xc1 + len:0x22(34) + code:0x5900 + data:01 00 47 00 6a 00 04 00 01 00 ... 00>
<flag:0xc1 + len:0x0c(12) + code:0xec30 + data:14 00 00 00 0f 00 00 00>

<flag:0xc1 + len:0x0a(10) + code:0xf311 + data:fe 00 1d 16 00 68>
<flag:0xc1 + len:0x0a(10) + code:0xf311 + data:fe 00 1e 16 00 68>
<flag:0xc1 + len:0x20(32) + code:0x2d00 + data:>
<flag:0xc1 + len:0x07(07) + code:0x0701 + data:2e 7d ce>
<flag:0xc1 + len:0x0c(12) + code:0xec30 + data:58 00 00 00 53 00 00 00>

<flag:0xc1 + len:0x0a(10) + code:0xf311 + data:fe 00 1f 31 00 68>
<flag:0xc1 + len:0x10(16) + code:0x26fe + data:06 78 00 37 fe 00 00 00 78 06 00 00>
<flag:0xc1 + len:0x0c(12) + code:0x27fe + data:02 f8 03 0d f8 02 00 00>
<flag:0xc1 + len:0x22(34) + code:0x5900 + data:01 00 47 00 6b 00 04 00 01 00 ... 00>
<flag:0xc1 + len:0x0c(12) + code:0xec30 + data:5f 00 00 00 5a 00 00 00>
<flag:0xc1 + len:0x04(04) + code:0xd200>
<flag:0xc1 + len:0x0a(10) + code:0xd20c + data:00 02 de 07 7c 00>
<flag:0xc1 + len:0x0a(10) + code:0xd215 + data:47 02 de 07 01 00>
<flag:0xc1 + len:0x06(06) + code:0xfa0b + data:00 00>
<flag:0xc3 + len:0x63(99) + ...>
<flag:0xc4 + len:0x01c4(8+444)>
// 响应(分包)
client <- server(1459)
<flag:0xc1 + len:0x28(40) + code:0xf350 + data:9e 00 ...>
<flag:0xc2 + len:0x010c(268) + code:0xf353 + data:...>
<flag:0xc1 + len:0x46(70) + code:0xf311 + data:...>
<flag:0xc1 + len:0x06(06) + code:0xa006 + data:aa ea>
<flag:0xc1 + len:0x29(47) + code:0xe703 + data>
<flag:0xc1 + len:0x29(47) + code:0xe703 + data>
...
<flag:0xc1 + len:0x29(47) + code:0xe703 + data>
<flag:0xc1 + len:0x08(08) + code:0x4d18 + data:03 00 00 00>
<flag:0xc1 + len:0x08(08) + code:0x8e01 + data:aa 4e 3d 1b>
<flag:0xc1 + len:0x31(49) + code:0x0d00 + data:>
<flag:0xc2 + len:0x003e(62) + ...>
<flag:0xc1 + len:0x18(24) + code:0xf807 + data:...>
<flag:0xc1 + len:0x05(05) + code:0xb801 + data:00>
<flag:0xc1 + len:0x10(16) + code:0x26ff + data:...>
<flag:0xc1 + len:0x10(16) + code:0x26fe + data:...>
<flag:0xc1 + len:0x0c(12) + code:0x27ff + data:02 66 01 45 66 02 00 00>
<flag:0xc1 + len:0x0c(12) + code:0x27fe + data:02 66 03 0d 66 02 00 00>
<flag:0xc1 + len:0x08(08) + code:0xbf52 + data:00 00 00 00>
<flag:0xc1 + len:0x05(05) + code:0xfaa2 + data:0f>
<flag:0xc2 + len:0x002c(44) + ...>
<flag:0xc1 + len:0x22(34) + code:0x5900 + data:...>

// 响应
client <- server(167)
<flag:0xc1 + len:0x24(36) + code:0xf330 + data:...>
<flag:0xc2 + len:0x003e(62) + ...>
<flag:0xc1 + len:0x18(24) + code:0xf807 + data:...>
<flag:0xc1 + len:0x2d(45) + code:0xde01 + data:...>

// 响应
client <- server(261)
<flag:0xc2 + len:0x0105(261) + ...>

// 响应
client <- server(378)
<flag:0xc1 + len:0x04(04) + code:0xf620>
<flag:0xc2 + len:0x0007(7) + ...>
<flag:0xc1 + len:0x08(08) + code:0x4d18 + data:00 00 00 00>
<flag:0xc4 + len:0x0024(36) + ...>
<flag:0xc1 + len:0x18(24) + code:0xf807 + data:...>
<flag:0xc4 + len:0x0014(20) + ...>
<flag:0xc1 + len:0x08(08) + code:0x4d18 + data:01 00 00 00>
<flag:0xc1 + len:0x05(05) + code:0xd211 + data:00>
<flag:0xc2 + len:0x00c2() + ...>
<flag:0xc1 + len:0x08(08) + code:0xd40f + data:84 25 2b 70>
<flag:0xc1 + len:0x08(08) + code:0xd40f + data:85 28 2a 30>
<flag:0xc1 + len:0x08(08) + code:0xd40f + data:86 23 31 00>
<flag:0xc1 + len:0x08(08) + code:0xd40f + data:87 28 2b 10>
<flag:0xc1 + len:0x08(08) + code:0xd40f + data:88 19 33 50>
<flag:0xc1 + len:0x08(08) + code:0xd40f + data:94 1a 36 50>
<flag:0xc1 + len:0x08(08) + code:0xd40f + data:95 26 2e 10>
<flag:0xc1 + len:0x08(08) + code:0xd40f + data:96 1b 2f 30>
<flag:0xc1 + len:0x08(08) + code:0xd40f + data:a7 16 2e 70>
```





#### End