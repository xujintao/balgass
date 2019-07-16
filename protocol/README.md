##### Connect protocol
express
```
ip.addr==192.168.0.100 and (tcp.port==44405 or tcp.port==56900)
```

这是什么？
```
client <- server(4):
c1 04 00 01

client(4) -> server:
c1 04 f4 06

client <- server(204):
c1 cc fa 00 "Server News" 00  
cc cc cc cc 0f 10 08 e8 28 f3 f4 09 66 2a 32 01  
01 00 00 00 04 f4 f4 09 60 f5 f4 09 80 e3 4b 01  
cc cc cc ...  

client <- server(236)
c2 00 63 fa 01 0a 0c dd 07 ff 00 9b 00 d7 00 64
00 87 00 "Info" ...  3f 00 "When no News added News Windows in Client will not be displayed"
c2 00 0b f4 06 00 01 00 00 00 00
```

服务器列表
```
client(4) -> server:
c1 04 f4 06

client <- server(11):
c2 00 0b f4 06 00 01 00 00 00 00
```

code为0的服务器ip和端口
```
client(6) -> server:
c1 06 f4 03 00 00

client <- server(22):
c1 16 f4 03 192.168.0.100 00 00 00 44 de
```
客户端RST掉CS连接并和GS建立连接

##### Game protocol
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

client(9) -> server:
c1 09 fa 15 86 5a ca e2 c4
```

version
```
client <- server(12):
c1 0c f1 00 01 2e 7e "10444"

客户端判断服务器版本，匹配就继续，不匹配就RST连接并退出程序  
(想办法禁用客户端的版本匹配)
```

login
```
client(179) -> server: 应该是username和passwd的加密信息
c3 b3 92 af ... cb 0f

client <- server(5): 登录成功
c1 05 f1 01 01
```

character
```
client(4) -> server:
c1 04 f3 0d

client <- server(49):
c1 31 0d 01 00 00 00 00 cd 00 00 00 00 "You are logged in as Free Player" 00 00 00 f8

client <- server(5):
c1 05 de 00 1f

client <- server(541): 角色信息
c1 0e fa 0a ...

```


##### end