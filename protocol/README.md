#### Frame Format
```
<flag + len + type + [data]> +
<flag + len + type + [data]> + ...

note:
flag=c1, single byte len
flag=c2, double bytes len, big endian
flag=c3, data has been encoded
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
```

version
```
client <- server(12):
c1 0c f1 00 01 2e 7e "10444"
<flag:0xc1 + len:12 + type:0xf100 + data:01 2e 7e "10444">

客户端判断服务器版本，匹配就继续，不匹配就RST连接并退出程序  
(想办法禁用客户端的版本匹配)
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
<type:0x0df3>

// 响应
client <- server(49):
c1 31 0d 01 00 00 00 00 cd 00 00 00 00 "You are logged in as Free Player" 00 00 00 f8
<flag:0xc1 + len:49 + type:0x0d01 + data>

// 响应
client <- server(5):
c1 05 de 00 1f
<flag:0xc1 + len:5 + type:0xde00 + data:0x1f>

// 响应
client <- server(541): 角色信息
c1 0e fa 0a ...
<flag:0xc1 + len:14 + type:0xfa0a + data:00...> +
<flag:0xc1 + len:44 + type:0x00f3 + data:"战士编码"..."juju"...> +
<flag:0xc2 + len:452 + type:0xfaf8 + data:...> +
<flag:0xc1 + len:31 + type:0xfa12 + data:...>

```


#### End