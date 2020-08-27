2019/272 01:22:48.077408    9409 → 44405 [SYN] Seq=0 Win=8192 Len=0 MSS=1460 WS=4 SACK_PERM=1

Frame 237: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 0, Len: 0

2019/272 01:22:48.078224    44405 → 9409 [SYN, ACK] Seq=0 Ack=1 Win=8192 Len=0 MSS=1460 WS=256 SACK_PERM=1

Frame 239: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 44405, Dst Port: 9409, Seq: 0, Ack: 1, Len: 0

2019/272 01:22:48.081914    9409 → 44405 [ACK] Seq=1 Ack=1 Win=65700 Len=0

Frame 242: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 1, Ack: 1, Len: 0

2019/272 01:22:48.088616    44405 → 9409 [PSH, ACK] Seq=1 Ack=1 Win=65536 Len=4

Frame 245: 58 bytes on wire (464 bits), 58 bytes captured (464 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 44405, Dst Port: 9409, Seq: 1, Ack: 1, Len: 4
<flag:c1 + len:04 + code:0001 + data:null>

2019/272 01:22:48.125602    9409 → 44405 [PSH, ACK] Seq=1 Ack=5 Win=65696 Len=4

Frame 249: 58 bytes on wire (464 bits), 58 bytes captured (464 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 1, Ack: 5, Len: 4
<flag:c1 + len:04 + code:f406 + data:null>

2019/272 01:22:48.127296    44405 → 9409 [PSH, ACK] Seq=5 Ack=5 Win=65536 Len=16

Frame 251: 70 bytes on wire (560 bits), 70 bytes captured (560 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 44405, Dst Port: 9409, Seq: 5, Ack: 5, Len: 16
<flag:c1 + len:10 + code:fa00 + data:53 65 72 76 65 72 20 4e 65 77 ...>

2019/272 01:22:48.327074    9409 → 44405 [ACK] Seq=5 Ack=21 Win=65680 Len=0

Frame 263: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 5, Ack: 21, Len: 0

2019/272 01:22:48.327781    44405 → 9409 [PSH, ACK] Seq=21 Ack=5 Win=65536 Len=236

Frame 265: 290 bytes on wire (2320 bits), 290 bytes captured (2320 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 44405, Dst Port: 9409, Seq: 21, Ack: 5, Len: 236
<flag:c2 + len:0063 + code:fa01 + data:0a 0c dd 07 ff 00 9b 00 d7 00 ...>
<flag:c2 + len:007e + code:fa01 + data:0a 0c dd 07 ff 00 9b 00 d7 00 ...>
<flag:c2 + len:000b + code:f406 + data:00 01 00 00 00 00>

2019/272 01:22:48.527096    9409 → 44405 [ACK] Seq=5 Ack=257 Win=65444 Len=0

Frame 270: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 5, Ack: 257, Len: 0

2019/272 01:22:50.458564    9409 → 44405 [PSH, ACK] Seq=5 Ack=257 Win=65444 Len=19

Frame 274: 73 bytes on wire (584 bits), 73 bytes captured (584 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 5, Ack: 257, Len: 19
<flag:c3 + len:13><flag:c1 + len:08 + code:f33c + data:24 07 e5 53>

2019/272 01:22:50.658183    44405 → 9409 [ACK] Seq=257 Ack=24 Win=65536 Len=0

Frame 277: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 44405, Dst Port: 9409, Seq: 257, Ack: 24, Len: 0

2019/272 01:23:17.156785    9409 → 44405 [PSH, ACK] Seq=24 Ack=257 Win=65444 Len=19

Frame 306: 73 bytes on wire (584 bits), 73 bytes captured (584 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 24, Ack: 257, Len: 19
<flag:c3 + len:13><flag:c1 + len:08 + code:f33c + data:24 07 e5 53>

2019/272 01:23:17.358200    44405 → 9409 [ACK] Seq=257 Ack=43 Win=65536 Len=0

Frame 309: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 44405, Dst Port: 9409, Seq: 257, Ack: 43, Len: 0

2019/272 01:23:43.866287    9409 → 44405 [PSH, ACK] Seq=43 Ack=257 Win=65444 Len=4

Frame 343: 58 bytes on wire (464 bits), 58 bytes captured (464 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 43, Ack: 257, Len: 4
<flag:c1 + len:04 + code:f406 + data:null>

2019/272 01:23:43.868394    44405 → 9409 [PSH, ACK] Seq=257 Ack=47 Win=65536 Len=11

Frame 345: 65 bytes on wire (520 bits), 65 bytes captured (520 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 44405, Dst Port: 9409, Seq: 257, Ack: 47, Len: 11
<flag:c2 + len:000b + code:f406 + data:00 01 00 00 00 00>

2019/272 01:23:44.076242    9409 → 44405 [ACK] Seq=47 Ack=268 Win=65432 Len=0

Frame 349: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 47, Ack: 268, Len: 0

2019/272 01:23:44.300916    9409 → 44405 [PSH, ACK] Seq=47 Ack=268 Win=65432 Len=19

Frame 352: 73 bytes on wire (584 bits), 73 bytes captured (584 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 47, Ack: 268, Len: 19
<flag:c3 + len:13><flag:c1 + len:08 + code:f33c + data:24 07 e5 53>

2019/272 01:23:44.506266    44405 → 9409 [ACK] Seq=268 Ack=66 Win=65536 Len=0

Frame 355: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 44405, Dst Port: 9409, Seq: 268, Ack: 66, Len: 0

2019/272 01:23:44.757616    9409 → 44405 [PSH, ACK] Seq=66 Ack=268 Win=65432 Len=6

Frame 358: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 66, Ack: 268, Len: 6
<flag:c1 + len:06 + code:f403 + data:00 00>

2019/272 01:23:44.758433    44405 → 9409 [PSH, ACK] Seq=268 Ack=72 Win=65536 Len=22

Frame 360: 76 bytes on wire (608 bits), 76 bytes captured (608 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 44405, Dst Port: 9409, Seq: 268, Ack: 72, Len: 22
<flag:c1 + len:16 + code:f403 + data:31 39 32 2e 31 36 38 2e 30 2e ...>

2019/272 01:23:44.797608    9409 → 44405 [RST, ACK] Seq=72 Ack=290 Win=0 Len=0

Frame 364: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 72, Ack: 290, Len: 0

2019/272 01:23:44.797926    9412 → 56900 [SYN] Seq=0 Win=8192 Len=0 MSS=1460 WS=4 SACK_PERM=1

Frame 365: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 0, Len: 0

2019/272 01:23:44.798456    9409 → 44405 [RST, ACK] Seq=72 Ack=290 Win=0 Len=0

Frame 366: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32), Dst: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9409, Dst Port: 44405, Seq: 72, Ack: 290, Len: 0

2019/272 01:23:44.799071    56900 → 9412 [SYN, ACK] Seq=0 Ack=1 Win=8192 Len=0 MSS=1460 WS=256 SACK_PERM=1

Frame 369: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 0, Ack: 1, Len: 0

2019/272 01:23:44.799861    9412 → 56900 [ACK] Seq=1 Ack=1 Win=65700 Len=0

Frame 372: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 1, Ack: 1, Len: 0

2019/272 01:23:44.800890    56900 → 9412 [PSH, ACK] Seq=1 Ack=1 Win=65536 Len=12

Frame 375: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 1, Ack: 1, Len: 12
<flag:c1 + len:0c + code:f100 + data:01 2e 7d 31 30 35 32 35>

2019/272 01:23:45.000316    9412 → 56900 [ACK] Seq=1 Ack=13 Win=65688 Len=0

Frame 379: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 1, Ack: 13, Len: 0

2019/272 01:23:45.304048    9412 → 56900 [PSH, ACK] Seq=1 Ack=13 Win=65688 Len=9

Frame 382: 63 bytes on wire (504 bits), 63 bytes captured (504 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 1, Ack: 13, Len: 9
<flag:c1 + len:09 + code:fa11 + data:8b ff 55 8b ec>

2019/272 01:23:45.560307    56900 → 9412 [ACK] Seq=13 Ack=10 Win=65536 Len=0

Frame 386: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 13, Ack: 10, Len: 0

2019/272 01:23:46.652623    9412 → 56900 [PSH, ACK] Seq=10 Ack=13 Win=65688 Len=35

Frame 389: 89 bytes on wire (712 bits), 89 bytes captured (712 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 10, Ack: 13, Len: 35
<flag:c3 + len:23><flag:c1 + len:17 + code:0e3d + data:00 c4 00 00 00 00 00 00 00 39 ...>

2019/272 01:23:46.855380    56900 → 9412 [ACK] Seq=13 Ack=45 Win=65536 Len=0

Frame 393: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 13, Ack: 45, Len: 0

2019/272 01:23:50.140440    9412 → 56900 [PSH, ACK] Seq=45 Ack=13 Win=65688 Len=179

Frame 398: 233 bytes on wire (1864 bits), 233 bytes captured (1864 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 45, Ack: 13, Len: 179
<flag:c3 + len:b3><flag:c1 + len:a3 + code:f101 + data:9d ab c6 95 a1 ab fc cf ab fc ...>

2019/272 01:23:50.162934    56900 → 9412 [PSH, ACK] Seq=13 Ack=224 Win=65280 Len=5

Frame 401: 59 bytes on wire (472 bits), 59 bytes captured (472 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 13, Ack: 224, Len: 5
<flag:c1 + len:05 + code:f101 + data:01>

2019/272 01:23:50.186711    9412 → 56900 [PSH, ACK] Seq=224 Ack=18 Win=65680 Len=35

Frame 404: 89 bytes on wire (712 bits), 89 bytes captured (712 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 224, Ack: 18, Len: 35
<flag:c3 + len:23><flag:c1 + len:17 + code:0e4b + data:00 99 00 00 00 00 00 00 00 39 ...>

2019/272 01:23:50.396524    56900 → 9412 [ACK] Seq=18 Ack=259 Win=65280 Len=0

Frame 407: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 18, Ack: 259, Len: 0

2019/272 01:23:50.397469    9412 → 56900 [PSH, ACK] Seq=259 Ack=18 Win=65680 Len=4

Frame 409: 58 bytes on wire (464 bits), 58 bytes captured (464 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 259, Ack: 18, Len: 4
<flag:c1 + len:04 + code:f300 + data:null>

2019/272 01:23:50.440388    56900 → 9412 [PSH, ACK] Seq=18 Ack=263 Win=65280 Len=5

Frame 413: 59 bytes on wire (472 bits), 59 bytes captured (472 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 18, Ack: 263, Len: 5
<flag:c1 + len:05 + code:de00 + data:1f>

2019/272 01:23:50.636561    9412 → 56900 [ACK] Seq=263 Ack=23 Win=65676 Len=0

Frame 416: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 263, Ack: 23, Len: 0

2019/272 01:23:50.637587    56900 → 9412 [PSH, ACK] Seq=23 Ack=263 Win=65280 Len=662

Frame 418: 716 bytes on wire (5728 bits), 716 bytes captured (5728 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 23, Ack: 263, Len: 662
<flag:c1 + len:0e + code:fa0a + data:00 00 00 00 00 00 00 00 00 00>
<flag:c1 + len:74 + code:f300 + data:1f 00 03 00 00 6a 75 6a 75 00 ...>
<flag:c2 + len:01c4 + code:faf8 + data:00 00 00 00 00 00 00 00 00 00 ...>
<flag:c1 + len:1f + code:fa12 + data:64 96 ff ff 1e 00 cd dc ef ff ...>
<flag:c1 + len:31 + code:0d01 + data:00 00 00 00 7f 00 00 00 00 59 ...>

2019/272 01:23:50.836539    9412 → 56900 [ACK] Seq=263 Ack=685 Win=65016 Len=0

Frame 422: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 263, Ack: 685, Len: 0

2019/272 01:23:55.451709    9412 → 56900 [PSH, ACK] Seq=263 Ack=685 Win=65016 Len=14

Frame 435: 68 bytes on wire (544 bits), 68 bytes captured (544 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 263, Ack: 685, Len: 14
<flag:c1 + len:0e + code:f315 + data:6a 75 6a 75 00 00 00 00 00 00>

2019/272 01:23:55.456874    56900 → 9412 [PSH, ACK] Seq=685 Ack=277 Win=65280 Len=15

Frame 438: 69 bytes on wire (552 bits), 69 bytes captured (552 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 685, Ack: 277, Len: 15
<flag:c1 + len:0f + code:f315 + data:00 00 00 00 00 00 00 00 00 00 ...>

2019/272 01:23:55.656776    9412 → 56900 [ACK] Seq=277 Ack=700 Win=65000 Len=0

Frame 441: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 277, Ack: 700, Len: 0

2019/272 01:23:55.679610    9412 → 56900 [PSH, ACK] Seq=277 Ack=700 Win=65000 Len=6

Frame 444: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 277, Ack: 700, Len: 6
<flag:c1 + len:06 + code:4e11 + data:ff ff>

2019/272 01:23:55.876732    56900 → 9412 [ACK] Seq=700 Ack=283 Win=65280 Len=0

Frame 447: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 700, Ack: 283, Len: 0

2019/272 01:23:55.877858    9412 → 56900 [PSH, ACK] Seq=283 Ack=700 Win=65000 Len=15

Frame 449: 69 bytes on wire (552 bits), 69 bytes captured (552 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 283, Ack: 700, Len: 15
<flag:c1 + len:0f + code:f303 + data:6a 75 6a 75 00 00 00 00 00 00 ...>

2019/272 01:23:55.906532    56900 → 9412 [PSH, ACK] Seq=700 Ack=298 Win=65280 Len=16

Frame 453: 70 bytes on wire (560 bits), 70 bytes captured (560 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 700, Ack: 298, Len: 16
<flag:c1 + len:10 + code:26fe + data:04 c8 00 33 e7 00 00 00 c8 04 ...>

2019/272 01:23:55.910384    56900 → 9412 [PSH, ACK] Seq=716 Ack=298 Win=65280 Len=1460

Frame 456: 1514 bytes on wire (12112 bits), 1514 bytes captured (12112 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 716, Ack: 298, Len: 1460
<flag:c1 + len:0c + code:27fe + data:02 66 02 8b 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 44 00 66 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 33 e7 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 8b 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 44 00 66 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 33 e7 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 8b 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 44 00 66 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 33 e7 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 8b 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 44 00 66 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 34 8d 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 8b 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 44 00 66 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 34 8d 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 8b 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 44 00 66 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 34 8d 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 8b 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 44 00 66 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 34 8d 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 8b 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 44 00 66 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 34 e6 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 96 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6a 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 34 e6 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 96 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6a 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 34 e6 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 96 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6a 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:0f 00 00 00 0b 00 00 00>
<flag:c1 + len:10 + code:26fe + data:04 c8 00 35 4c 00 00 00 c8 04 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 a5 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6a 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:14 00 00 00 0f 00 00 00>
<flag:c1 + len:10 + code:26fe + data:05 8e 00 35 9b 00 00 00 8e 05 ...>
<flag:c1 + len:0c + code:27fe + data:02 66 02 b8 66 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6a 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:14 00 00 00 0f 00 00 00>
<flag:c1 + len:10 + code:26fe + data:05 8e 00 35 f4 00 00 00 8e 05 ...>
<flag:c1 + len:0c + code:27fe + data:02 b0 03 02 b0 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6a 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:14 00 00 00 0f 00 00 00>
<flag:c1 + len:10 + code:26fe + data:05 8e 00 35 f4 00 00 00 8e 05 ...>
<flag:c1 + len:0c + code:27fe + data:02 b0 03 02 b0 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6a 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:14 00 00 00 0f 00 00 00>
<flag:c1 + len:10 + code:26fe + data:05 8e 00 35 f4 00 00 00 8e 05 ...>
<flag:c1 + len:0c + code:27fe + data:02 b0 03 02 b0 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6a 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:14 00 00 00 0f 00 00 00>
<flag:c1 + len:0a + code:f311 + data:fe 00 1d 16 00 68>
<flag:c1 + len:0a + code:f311 + data:fe 00 1e 16 00 68>
<flag:c1 + len:20 + code:2d00 + data:4a 00 00 00 00 00 00 00 f6 ff ...>
<flag:c1 + len:07 + code:0701 + data:2e 7d ce>
<flag:c1 + len:0c + code:ec30 + data:58 00 00 00 53 00 00 00>
<flag:c1 + len:0a + code:f311 + data:fe 00 1f 31 00 00>
<flag:c1 + len:10 + code:26fe + data:06 78 00 37 fe 00 00 00 78 06 ...>
<flag:c1 + len:0c + code:27fe + data:02 f8 03 0d f8 02 00 00>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6b 00 04 00 01 00 ...>
<flag:c1 + len:0c + code:ec30 + data:5f 00 00 00 5a 00 00 00>
<flag:c1 + len:04 + code:d200 + data:null>
<flag:c1 + len:0a + code:d20c + data:00 02 de 07 7c 00>
<flag:c1 + len:0a + code:d215 + data:47 02 de 07 01 00>
<flag:c1 + len:06 + code:fa0b + data:00 00>
<flag:c3 + len:63><flag:c1 + len:5d + code:5cf3 + data:03 25 31 38 00 00 00 00 03 1c ...>

2019/272 01:23:55.910396    56900 → 9412 [PSH, ACK] Seq=2176 Ack=298 Win=65280 Len=1459

Frame 457: 1513 bytes on wire (12104 bits), 1513 bytes captured (12104 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 2176, Ack: 298, Len: 1459
[2 Reassembled TCP Segments (452 bytes): #456(8), #457(444)]
<flag:c4 + len:01c4><flag:c2 + len:01b4 + code:b3f3 + data:10 21 00 1a ec 75 40 00 00 07 ...>
<flag:c1 + len:28 + code:f350 + data:9e 00 00 00 00 03 1c 51 d3 ba ...>
<flag:c2 + len:010c + code:f353 + data:00 00 00 10 00 00 00 01 0a 00 ...>
<flag:c1 + len:46 + code:f311 + data:10 00 00 2c 00 04 01 43 00 03 ...>
<flag:c1 + len:06 + code:a006 + data:aa ea>
<flag:c1 + len:29 + code:e703 + data:00 00 01 00 8b 7d 48 75 6e 74 ...>
<flag:c1 + len:29 + code:e703 + data:01 00 01 00 7c 6d 48 75 6e 74 ...>
<flag:c1 + len:29 + code:e703 + data:02 00 01 00 8c 5f 48 75 6e 74 ...>
<flag:c1 + len:29 + code:e703 + data:03 00 01 00 96 6d 48 75 6e 74 ...>
<flag:c1 + len:29 + code:e703 + data:04 00 01 00 b9 bc 53 61 66 65 ...>
<flag:c1 + len:29 + code:e703 + data:05 00 01 00 38 b1 53 61 66 65 ...>
<flag:c1 + len:29 + code:e703 + data:06 00 01 00 44 34 53 61 66 65 ...>
<flag:c1 + len:29 + code:e703 + data:07 00 01 00 c5 0d 53 61 66 65 ...>
<flag:c1 + len:08 + code:4d18 + data:03 00 00 00>
<flag:c1 + len:08 + code:8e01 + data:aa 4e 3d 1b>
<flag:c1 + len:31 + code:0d00 + data:00 00 00 00 9e 7d 9b fe fe 57 ...>
<flag:c2 + len:003e + code:1201 + data:ae 7d 25 31 38 1a 1a 11 11 1d ...>
<flag:c1 + len:18 + code:f807 + data:00 00 00 00 00 00 00 00 00 00 ...>
<flag:c1 + len:05 + code:b801 + data:00>
<flag:c1 + len:10 + code:26ff + data:04 c8 00 37 fe ff ff ff c8 04 ...>
<flag:c1 + len:10 + code:26fe + data:06 78 00 37 fe ff ff ff 78 06 ...>
<flag:c1 + len:0c + code:27ff + data:02 66 01 45 66 02 00 00>
<flag:c1 + len:0c + code:27fe + data:02 66 03 0d 66 02 00 00>
<flag:c1 + len:08 + code:bf52 + data:00 00 00 00>
<flag:c1 + len:05 + code:faa2 + data:0f>
<flag:c2 + len:002c + code:faa4 + data:03 00 00 c0 c6 2d 00 80 8d 5b ...>
<flag:c1 + len:22 + code:5900 + data:01 00 47 00 6b 00 04 00 01 00 ...>

2019/272 01:23:55.911811    9412 → 56900 [ACK] Seq=298 Ack=2176 Win=65700 Len=0

Frame 459: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 298, Ack: 2176, Len: 0

2019/272 01:23:55.915847    56900 → 9412 [PSH, ACK] Seq=3635 Ack=298 Win=65280 Len=167

Frame 464: 221 bytes on wire (1768 bits), 221 bytes captured (1768 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 3635, Ack: 298, Len: 167
<flag:c1 + len:24 + code:f330 + data:ff ff 00 16 00 29 00 2b 00 30 ...>
<flag:c2 + len:003e + code:1201 + data:ae 7d 25 31 38 1a 1a 11 11 1d ...>
<flag:c1 + len:18 + code:f807 + data:00 00 00 00 00 00 00 00 00 00 ...>
<flag:c1 + len:2d + code:d201 + data:00 00 00 00 f0 54 df 01 00 00 ...>

2019/272 01:23:55.916815    9412 → 56900 [ACK] Seq=298 Ack=3802 Win=65700 Len=0

Frame 467: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 298, Ack: 3802, Len: 0

2019/272 01:23:55.920130    56900 → 9412 [PSH, ACK] Seq=3802 Ack=298 Win=65280 Len=261

Frame 471: 315 bytes on wire (2520 bits), 315 bytes captured (2520 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 3802, Ack: 298, Len: 261
<flag:c2 + len:0105 + code:ae00 + data:91 88 08 00 16 00 2b 00 00 00 ...>

2019/272 01:23:56.116733    9412 → 56900 [ACK] Seq=298 Ack=4063 Win=65436 Len=0

Frame 474: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 298, Ack: 4063, Len: 0

2019/272 01:23:56.117594    56900 → 9412 [PSH, ACK] Seq=4063 Ack=298 Win=65280 Len=378

Frame 476: 432 bytes on wire (3456 bits), 432 bytes captured (3456 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4063, Ack: 298, Len: 378
<flag:c1 + len:04 + code:f620 + data:null>
<flag:c2 + len:0007 + code:c000 + data:32 00>
<flag:c1 + len:08 + code:4d18 + data:00 00 00 00>
<flag:c4 + len:0024><flag:c2 + len:0014 + code:134e + data:02 01 02 da 08 f3 00 00 d0 01 ...>
<flag:c1 + len:18 + code:f807 + data:01 4d 00 00 01 00 00 00 0e 00 ...>
<flag:c4 + len:0014><flag:c2 + len:0007 + code:064d + data:02 00>
<flag:c1 + len:08 + code:4d18 + data:01 00 00 00>
<flag:c1 + len:05 + code:d211 + data:00>
<flag:c2 + len:00c2 + code:1309 + data:0f 84 01 bb 28 2d 27 2d 60 00 ...>
<flag:c1 + len:08 + code:d40f + data:84 25 2b 70>
<flag:c1 + len:08 + code:d40f + data:85 28 2a 30>
<flag:c1 + len:08 + code:d40f + data:86 23 31 00>
<flag:c1 + len:08 + code:d40f + data:87 28 2b 10>
<flag:c1 + len:08 + code:d40f + data:88 19 33 50>
<flag:c1 + len:08 + code:d40f + data:94 1a 36 50>
<flag:c1 + len:08 + code:d40f + data:95 26 2e 10>
<flag:c1 + len:08 + code:d40f + data:96 1b 2f 30>
<flag:c1 + len:08 + code:d40f + data:a7 16 2e 70>

2019/272 01:23:56.310032    9412 → 56900 [PSH, ACK] Seq=298 Ack=4441 Win=65060 Len=12

Frame 480: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 298, Ack: 4441, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 7c 56 4e 00>

2019/272 01:23:56.311772    56900 → 9412 [PSH, ACK] Seq=4441 Ack=310 Win=65280 Len=68

Frame 482: 122 bytes on wire (976 bits), 122 bytes captured (976 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4441, Ack: 310, Len: 68
<flag:c1 + len:10 + code:26fe + data:06 78 00 37 fe 00 00 00 78 06 ...>
<flag:c1 + len:0c + code:27fe + data:02 f8 03 0d f8 02 00 00>
<flag:c2 + len:0008 + code:ee01 + data:00 00 00>
<flag:c2 + len:0008 + code:ee01 + data:00 00 01>
<flag:c2 + len:0008 + code:ec32 + data:01 03 00>
<flag:c2 + len:0006 + code:4f00 + data:00>
<flag:c1 + len:05 + code:cd01 + data:00>
<flag:c1 + len:05 + code:dc01 + data:00>

2019/272 01:23:56.516717    9412 → 56900 [ACK] Seq=310 Ack=4509 Win=64992 Len=0

Frame 486: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 310, Ack: 4509, Len: 0

2019/272 01:23:57.044494    56900 → 9412 [PSH, ACK] Seq=4509 Ack=310 Win=65280 Len=6

Frame 489: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4509, Ack: 310, Len: 6
<flag:c1 + len:06 + code:1401 + data:0f a7>

2019/272 01:23:57.316746    9412 → 56900 [ACK] Seq=310 Ack=4515 Win=64984 Len=0

Frame 492: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 310, Ack: 4515, Len: 0

2019/272 01:23:57.319835    56900 → 9412 [PSH, ACK] Seq=4515 Ack=310 Win=65280 Len=24

Frame 494: 78 bytes on wire (624 bits), 78 bytes captured (624 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4515, Ack: 310, Len: 24
<flag:c1 + len:08 + code:d40f + data:88 1a 2d 10>
<flag:c1 + len:08 + code:d40f + data:94 19 36 60>
<flag:c1 + len:08 + code:d40f + data:96 1a 31 50>

2019/272 01:23:57.516749    9412 → 56900 [ACK] Seq=310 Ack=4539 Win=64960 Len=0

Frame 498: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 310, Ack: 4539, Len: 0

2019/272 01:23:58.043927    56900 → 9412 [PSH, ACK] Seq=4539 Ack=310 Win=65280 Len=26

Frame 501: 80 bytes on wire (640 bits), 80 bytes captured (640 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4539, Ack: 310, Len: 26
<flag:c2 + len:001a + code:1301 + data:0f a7 01 be 17 2f 1a 2f 40 00 ...>

2019/272 01:23:58.091929    9412 → 56900 [PSH, ACK] Seq=310 Ack=4565 Win=64936 Len=24

Frame 504: 78 bytes on wire (624 bits), 78 bytes captured (624 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 310, Ack: 4565, Len: 24
<flag:c1 + len:18 + code:fa0d + data:5f e7 60 96 37 ea 60 7f 17 a6 ...>

2019/272 01:23:58.092879    56900 → 9412 [PSH, ACK] Seq=4565 Ack=334 Win=65280 Len=185

Frame 506: 239 bytes on wire (1912 bits), 239 bytes captured (1912 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4565, Ack: 334, Len: 185
<flag:c1 + len:08 + code:d40f + data:84 25 2d 50>
<flag:c1 + len:08 + code:d40f + data:85 29 2d 50>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:87 24 32 50>
<flag:c1 + len:08 + code:d40f + data:88 1c 34 50>
<flag:c1 + len:08 + code:d40f + data:94 1c 33 20>
<flag:c1 + len:08 + code:d40f + data:96 1a 34 50>
<flag:c1 + len:08 + code:d40f + data:a7 1b 2e 30>
<flag:c1 + len:10 + code:110f + data:86 04 e2 00 44 00 00 40 e2 04 ...>
<flag:c1 + len:12 + code:fa05 + data:86 0f 39 db 00 00 90 e2 00 00 ...>
<flag:c1 + len:12 + code:fa05 + data:86 0f 39 db 00 00 90 e2 00 00 ...>
<flag:c1 + len:0e + code:ec10 + data:0f 86 00 e2 00 90 00 db 00 39>
<flag:c1 + len:0c + code:27ff + data:02 ce 03 0d ce 02 00 00>
<flag:c1 + len:04 + code:0f20 + data:null>

2019/272 01:23:58.286795    9412 → 56900 [ACK] Seq=334 Ack=4750 Win=64752 Len=0

Frame 510: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 334, Ack: 4750, Len: 0

2019/272 01:23:59.061326    56900 → 9412 [PSH, ACK] Seq=4750 Ack=334 Win=65280 Len=8

Frame 513: 62 bytes on wire (496 bits), 62 bytes captured (496 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4750, Ack: 334, Len: 8
<flag:c1 + len:08 + code:d40f + data:88 18 34 50>

2019/272 01:23:59.256857    9412 → 56900 [ACK] Seq=334 Ack=4758 Win=64744 Len=0

Frame 516: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 334, Ack: 4758, Len: 0

2019/272 01:23:59.259243    56900 → 9412 [PSH, ACK] Seq=4758 Ack=334 Win=65280 Len=32

Frame 518: 86 bytes on wire (688 bits), 86 bytes captured (688 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4758, Ack: 334, Len: 32
<flag:c1 + len:08 + code:d40f + data:94 19 32 10>
<flag:c1 + len:08 + code:d40f + data:95 26 30 50>
<flag:c1 + len:08 + code:d40f + data:96 19 2f 00>
<flag:c1 + len:08 + code:d40f + data:a7 16 31 50>

2019/272 01:23:59.456845    9412 → 56900 [ACK] Seq=334 Ack=4790 Win=64712 Len=0

Frame 522: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 334, Ack: 4790, Len: 0

2019/272 01:23:59.541262    56900 → 9412 [PSH, ACK] Seq=4790 Ack=334 Win=65280 Len=6

Frame 525: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4790, Ack: 334, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:23:59.736876    9412 → 56900 [ACK] Seq=334 Ack=4796 Win=64704 Len=0

Frame 528: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 334, Ack: 4796, Len: 0

2019/272 01:23:59.737816    56900 → 9412 [PSH, ACK] Seq=4796 Ack=334 Win=65280 Len=41

Frame 530: 95 bytes on wire (760 bits), 95 bytes captured (760 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4796, Ack: 334, Len: 41
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:23:59.936850    9412 → 56900 [ACK] Seq=334 Ack=4837 Win=64664 Len=0

Frame 534: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 334, Ack: 4837, Len: 0

2019/272 01:24:00.055244    56900 → 9412 [PSH, ACK] Seq=4837 Ack=334 Win=65280 Len=8

Frame 537: 62 bytes on wire (496 bits), 62 bytes captured (496 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4837, Ack: 334, Len: 8
<flag:c1 + len:08 + code:d40f + data:85 27 2f 60>

2019/272 01:24:00.256896    9412 → 56900 [ACK] Seq=334 Ack=4845 Win=64656 Len=0

Frame 540: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 334, Ack: 4845, Len: 0

2019/272 01:24:00.258014    56900 → 9412 [PSH, ACK] Seq=4845 Ack=334 Win=65280 Len=79

Frame 542: 133 bytes on wire (1064 bits), 133 bytes captured (1064 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4845, Ack: 334, Len: 79
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1d 35 40>
<flag:c1 + len:08 + code:d40f + data:94 16 36 70>
<flag:c1 + len:08 + code:d40f + data:96 1b 33 50>
<flag:c1 + len:08 + code:d40f + data:a7 1b 2c 30>

2019/272 01:24:00.310568    9412 → 56900 [PSH, ACK] Seq=334 Ack=4924 Win=64576 Len=12

Frame 546: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 334, Ack: 4924, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 7c 56 4e 00>

2019/272 01:24:00.506867    56900 → 9412 [ACK] Seq=4924 Ack=346 Win=65280 Len=0

Frame 549: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4924, Ack: 346, Len: 0

2019/272 01:24:00.556009    56900 → 9412 [PSH, ACK] Seq=4924 Ack=346 Win=65280 Len=32

Frame 552: 86 bytes on wire (688 bits), 86 bytes captured (688 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4924, Ack: 346, Len: 32
<flag:c1 + len:20 + code:2d00 + data:1c 00 14 00 00 00 00 00 0a 00 ...>

2019/272 01:24:00.806883    9412 → 56900 [ACK] Seq=346 Ack=4956 Win=64544 Len=0

Frame 559: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 346, Ack: 4956, Len: 0

2019/272 01:24:00.807978    56900 → 9412 [PSH, ACK] Seq=4956 Ack=346 Win=65280 Len=66

Frame 561: 120 bytes on wire (960 bits), 120 bytes captured (960 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 4956, Ack: 346, Len: 66
<flag:c1 + len:07 + code:0701 + data:2e 7d 38>
<flag:c1 + len:0c + code:ec30 + data:5f 00 00 00 5a 00 00 00>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:01.006888    9412 → 56900 [ACK] Seq=346 Ack=5022 Win=64480 Len=0

Frame 565: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 346, Ack: 5022, Len: 0

2019/272 01:24:01.050065    56900 → 9412 [PSH, ACK] Seq=5022 Ack=346 Win=65280 Len=8

Frame 568: 62 bytes on wire (496 bits), 62 bytes captured (496 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5022, Ack: 346, Len: 8
<flag:c1 + len:08 + code:d40f + data:88 18 36 50>

2019/272 01:24:01.107594    9412 → 56900 [PSH, ACK] Seq=346 Ack=5030 Win=64472 Len=4

Frame 571: 58 bytes on wire (464 bits), 58 bytes captured (464 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 346, Ack: 5030, Len: 4
<flag:c1 + len:04 + code:f326 + data:null>

2019/272 01:24:01.113998    56900 → 9412 [PSH, ACK] Seq=5030 Ack=350 Win=65280 Len=36

Frame 573: 90 bytes on wire (720 bits), 90 bytes captured (720 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5030, Ack: 350, Len: 36
<flag:c1 + len:08 + code:d40f + data:94 18 34 10>
<flag:c1 + len:08 + code:d40f + data:96 1c 32 40>
<flag:c1 + len:08 + code:d40f + data:a7 18 2c 10>
<flag:c1 + len:0c + code:27ff + data:02 f8 03 0d f8 02 00 00>

2019/272 01:24:01.406901    9412 → 56900 [ACK] Seq=350 Ack=5066 Win=64436 Len=0

Frame 577: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 350, Ack: 5066, Len: 0

2019/272 01:24:01.407989    56900 → 9412 [PSH, ACK] Seq=5066 Ack=350 Win=65280 Len=8

Frame 579: 62 bytes on wire (496 bits), 62 bytes captured (496 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5066, Ack: 350, Len: 8
<flag:c1 + len:08 + code:f326 + data:01 00 00 00>

2019/272 01:24:01.573059    9412 → 56900 [PSH, ACK] Seq=350 Ack=5074 Win=64428 Len=4

Frame 583: 58 bytes on wire (464 bits), 58 bytes captured (464 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 350, Ack: 5074, Len: 4
<flag:c1 + len:04 + code:f61a + data:null>

2019/272 01:24:01.575017    56900 → 9412 [PSH, ACK] Seq=5074 Ack=354 Win=65280 Len=192

Frame 585: 246 bytes on wire (1968 bits), 246 bytes captured (1968 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5074, Ack: 354, Len: 192
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1c 35 50>
<flag:c1 + len:08 + code:d40f + data:94 15 32 00>
<flag:c1 + len:08 + code:d40f + data:96 1e 33 30>
<flag:c1 + len:08 + code:d40f + data:a7 14 2f 70>
<flag:c1 + len:10 + code:110f + data:85 04 e2 00 44 00 00 40 e2 04 ...>
<flag:c1 + len:12 + code:fa05 + data:85 0f a8 dd 00 00 90 e2 00 00 ...>
<flag:c1 + len:12 + code:fa05 + data:85 0f a8 dd 00 00 90 e2 00 00 ...>
<flag:c1 + len:0e + code:ec10 + data:0f 85 00 e2 00 90 00 dd 00 a8>

2019/272 01:24:01.575710    9412 → 56900 [PSH, ACK] Seq=354 Ack=5266 Win=65700 Len=10

Frame 588: 64 bytes on wire (512 bits), 64 bytes captured (512 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 354, Ack: 5266, Len: 10
<flag:c1 + len:04 + code:f621 + data:null>
<flag:c1 + len:06 + code:4e11 + data:03 1a>

2019/272 01:24:01.577027    56900 → 9412 [PSH, ACK] Seq=5266 Ack=364 Win=65280 Len=5

Frame 591: 59 bytes on wire (472 bits), 59 bytes captured (472 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5266, Ack: 364, Len: 5
<flag:c1 + len:05 + code:f61a + data:00>

2019/272 01:24:01.806913    9412 → 56900 [ACK] Seq=364 Ack=5271 Win=65692 Len=0

Frame 595: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 364, Ack: 5271, Len: 0

2019/272 01:24:01.808002    56900 → 9412 [PSH, ACK] Seq=5271 Ack=364 Win=65280 Len=10

Frame 597: 64 bytes on wire (512 bits), 64 bytes captured (512 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5271, Ack: 364, Len: 10
<flag:c2 + len:000a + code:4e14 + data:01 2e 7d 1a 03>

2019/272 01:24:01.842114    9412 → 56900 [PSH, ACK] Seq=364 Ack=5281 Win=65684 Len=19

Frame 601: 73 bytes on wire (584 bits), 73 bytes captured (584 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 364, Ack: 5281, Len: 19
<flag:c3 + len:13><flag:c1 + len:08 + code:f331 + data:00 00 75 33>

2019/272 01:24:02.071329    56900 → 9412 [PSH, ACK] Seq=5281 Ack=383 Win=65280 Len=6

Frame 604: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5281, Ack: 383, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:02.306940    9412 → 56900 [ACK] Seq=383 Ack=5287 Win=65676 Len=0

Frame 608: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 383, Ack: 5287, Len: 0

2019/272 01:24:02.308007    56900 → 9412 [PSH, ACK] Seq=5287 Ack=383 Win=65280 Len=101

Frame 610: 155 bytes on wire (1240 bits), 155 bytes captured (1240 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5287, Ack: 383, Len: 101
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1a 31 10>
<flag:c1 + len:08 + code:d40f + data:94 1c 33 30>
<flag:c1 + len:08 + code:d40f + data:96 1b 2e 10>
<flag:c1 + len:08 + code:d40f + data:a7 17 2f 60>
<flag:c1 + len:0c + code:27ff + data:02 f8 03 0d f8 02 00 00>
<flag:c1 + len:10 + code:26ff + data:05 6a 00 37 fe 05 2e 40 6a 05 ...>

2019/272 01:24:02.506950    9412 → 56900 [ACK] Seq=383 Ack=5388 Win=65576 Len=0

Frame 614: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 383, Ack: 5388, Len: 0

2019/272 01:24:02.537121    56900 → 9412 [PSH, ACK] Seq=5388 Ack=383 Win=65280 Len=6

Frame 617: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5388, Ack: 383, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:02.806963    9412 → 56900 [ACK] Seq=383 Ack=5394 Win=65572 Len=0

Frame 620: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 383, Ack: 5394, Len: 0

2019/272 01:24:02.809038    56900 → 9412 [PSH, ACK] Seq=5394 Ack=383 Win=65280 Len=88

Frame 622: 142 bytes on wire (1136 bits), 142 bytes captured (1136 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5394, Ack: 383, Len: 88
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:03.006977    9412 → 56900 [ACK] Seq=383 Ack=5482 Win=65484 Len=0

Frame 626: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 383, Ack: 5482, Len: 0

2019/272 01:24:03.068595    56900 → 9412 [PSH, ACK] Seq=5482 Ack=383 Win=65280 Len=6

Frame 629: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5482, Ack: 383, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:03.306980    9412 → 56900 [ACK] Seq=383 Ack=5488 Win=65476 Len=0

Frame 632: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 383, Ack: 5488, Len: 0

2019/272 01:24:03.308076    56900 → 9412 [PSH, ACK] Seq=5488 Ack=383 Win=65280 Len=132

Frame 634: 186 bytes on wire (1488 bits), 186 bytes captured (1488 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5488, Ack: 383, Len: 132
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1b 32 40>
<flag:c1 + len:08 + code:d40f + data:94 1b 34 30>
<flag:c1 + len:08 + code:d40f + data:96 1a 33 50>
<flag:c1 + len:08 + code:d40f + data:a7 16 2d 10>
<flag:c1 + len:0c + code:27ff + data:02 f8 03 0d f8 02 00 00>

2019/272 01:24:03.506989    9412 → 56900 [ACK] Seq=383 Ack=5620 Win=65344 Len=0

Frame 638: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 383, Ack: 5620, Len: 0

2019/272 01:24:04.044848    56900 → 9412 [PSH, ACK] Seq=5620 Ack=383 Win=65280 Len=6

Frame 641: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5620, Ack: 383, Len: 6
<flag:c1 + len:06 + code:1401 + data:0f a7>

2019/272 01:24:04.307023    9412 → 56900 [ACK] Seq=383 Ack=5626 Win=65340 Len=0

Frame 644: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 383, Ack: 5626, Len: 0

2019/272 01:24:04.307778    56900 → 9412 [PSH, ACK] Seq=5626 Ack=383 Win=65280 Len=106

Frame 646: 160 bytes on wire (1280 bits), 160 bytes captured (1280 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5626, Ack: 383, Len: 106
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:08 + code:d40f + data:88 19 32 70>
<flag:c1 + len:08 + code:d40f + data:94 1a 31 30>
<flag:c1 + len:08 + code:d40f + data:96 1b 33 40>

2019/272 01:24:04.311703    9412 → 56900 [PSH, ACK] Seq=383 Ack=5732 Win=65232 Len=12

Frame 650: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 383, Ack: 5732, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:04.507034    56900 → 9412 [ACK] Seq=5732 Ack=395 Win=65280 Len=0

Frame 653: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5732, Ack: 395, Len: 0

2019/272 01:24:04.536534    56900 → 9412 [PSH, ACK] Seq=5732 Ack=395 Win=65280 Len=6

Frame 656: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5732, Ack: 395, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:04.807036    9412 → 56900 [ACK] Seq=395 Ack=5738 Win=65228 Len=0

Frame 659: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 395, Ack: 5738, Len: 0

2019/272 01:24:04.808170    56900 → 9412 [PSH, ACK] Seq=5738 Ack=395 Win=65280 Len=41

Frame 661: 95 bytes on wire (760 bits), 95 bytes captured (760 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5738, Ack: 395, Len: 41
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:05.007051    9412 → 56900 [ACK] Seq=395 Ack=5779 Win=65184 Len=0

Frame 665: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 395, Ack: 5779, Len: 0

2019/272 01:24:05.045436    56900 → 9412 [PSH, ACK] Seq=5779 Ack=395 Win=65280 Len=26

Frame 668: 80 bytes on wire (640 bits), 80 bytes captured (640 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5779, Ack: 395, Len: 26
<flag:c2 + len:001a + code:1301 + data:0f a7 01 be 17 2f 18 30 40 00 ...>

2019/272 01:24:05.307059    9412 → 56900 [ACK] Seq=395 Ack=5805 Win=65160 Len=0

Frame 671: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 395, Ack: 5805, Len: 0

2019/272 01:24:05.308163    56900 → 9412 [PSH, ACK] Seq=5805 Ack=395 Win=65280 Len=138

Frame 673: 192 bytes on wire (1536 bits), 192 bytes captured (1536 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5805, Ack: 395, Len: 138
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:08 + code:d40f + data:88 1d 36 50>
<flag:c1 + len:08 + code:d40f + data:94 18 33 00>
<flag:c1 + len:08 + code:d40f + data:96 1e 34 30>
<flag:c1 + len:08 + code:d40f + data:a7 19 31 50>
<flag:c1 + len:06 + code:2a07 + data:fb 00>
<flag:c1 + len:06 + code:2a0a + data:41 00>
<flag:c1 + len:06 + code:2a0b + data:46 00>
<flag:c1 + len:06 + code:2a09 + data:41 00>

2019/272 01:24:05.311788    9412 → 56900 [PSH, ACK] Seq=395 Ack=5943 Win=65020 Len=9

Frame 677: 63 bytes on wire (504 bits), 63 bytes captured (504 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 395, Ack: 5943, Len: 9
<flag:c1 + len:09 + code:fa11 + data:8b ff 55 8b ec>

2019/272 01:24:05.507068    56900 → 9412 [ACK] Seq=5943 Ack=404 Win=65280 Len=0

Frame 680: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5943, Ack: 404, Len: 0

2019/272 01:24:06.067859    56900 → 9412 [PSH, ACK] Seq=5943 Ack=404 Win=65280 Len=6

Frame 683: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5943, Ack: 404, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:06.269109    9412 → 56900 [ACK] Seq=404 Ack=5949 Win=65016 Len=0

Frame 686: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 404, Ack: 5949, Len: 0

2019/272 01:24:06.276213    56900 → 9412 [PSH, ACK] Seq=5949 Ack=404 Win=65280 Len=120

Frame 688: 174 bytes on wire (1392 bits), 174 bytes captured (1392 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 5949, Ack: 404, Len: 120
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1b 34 50>
<flag:c1 + len:08 + code:d40f + data:94 15 33 70>
<flag:c1 + len:08 + code:d40f + data:96 1b 35 60>
<flag:c1 + len:08 + code:d40f + data:a7 1b 2e 30>

2019/272 01:24:06.425775    9412 → 56900 [PSH, ACK] Seq=404 Ack=6069 Win=64896 Len=5

Frame 692: 59 bytes on wire (472 bits), 59 bytes captured (472 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 404, Ack: 6069, Len: 5
<flag:c1 + len:05 + code:1802 + data:7a>

2019/272 01:24:06.538330    56900 → 9412 [PSH, ACK] Seq=6069 Ack=409 Win=65280 Len=6

Frame 695: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6069, Ack: 409, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:06.634721    9412 → 56900 [PSH, ACK] Seq=409 Ack=6075 Win=64888 Len=35

Frame 698: 89 bytes on wire (712 bits), 89 bytes captured (712 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 409, Ack: 6075, Len: 35
<flag:c3 + len:23><flag:c1 + len:17 + code:0e8b + data:00 d4 00 5f 00 2c 01 5a 00 39 ...>

2019/272 01:24:06.636111    56900 → 9412 [PSH, ACK] Seq=6075 Ack=444 Win=65024 Len=57

Frame 700: 111 bytes on wire (888 bits), 111 bytes captured (888 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6075, Ack: 444, Len: 57
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:10 + code:26ff + data:06 78 00 37 fe cd 8b e9 78 06 ...>

2019/272 01:24:06.887123    9412 → 56900 [ACK] Seq=444 Ack=6132 Win=64832 Len=0

Frame 704: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 444, Ack: 6132, Len: 0

2019/272 01:24:07.070976    56900 → 9412 [PSH, ACK] Seq=6132 Ack=444 Win=65024 Len=6

Frame 707: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6132, Ack: 444, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:07.287133    9412 → 56900 [ACK] Seq=444 Ack=6138 Win=64828 Len=0

Frame 710: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 444, Ack: 6138, Len: 0

2019/272 01:24:07.289244    56900 → 9412 [PSH, ACK] Seq=6138 Ack=444 Win=65024 Len=120

Frame 712: 174 bytes on wire (1392 bits), 174 bytes captured (1392 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6138, Ack: 444, Len: 120
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 17 36 70>
<flag:c1 + len:08 + code:d40f + data:94 15 35 60>
<flag:c1 + len:08 + code:d40f + data:96 1b 32 10>
<flag:c1 + len:08 + code:d40f + data:a7 19 2f 50>

2019/272 01:24:07.311817    9412 → 56900 [PSH, ACK] Seq=444 Ack=6258 Win=64708 Len=12

Frame 716: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 444, Ack: 6258, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:07.587146    56900 → 9412 [ACK] Seq=6258 Ack=456 Win=65024 Len=0

Frame 720: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6258, Ack: 456, Len: 0

2019/272 01:24:08.066699    56900 → 9412 [PSH, ACK] Seq=6258 Ack=456 Win=65024 Len=6

Frame 723: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6258, Ack: 456, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:08.287175    9412 → 56900 [ACK] Seq=456 Ack=6264 Win=64700 Len=0

Frame 726: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6264, Len: 0

2019/272 01:24:08.288248    56900 → 9412 [PSH, ACK] Seq=6264 Ack=456 Win=65024 Len=120

Frame 728: 174 bytes on wire (1392 bits), 174 bytes captured (1392 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6264, Ack: 456, Len: 120
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1c 34 30>
<flag:c1 + len:08 + code:d40f + data:94 18 36 40>
<flag:c1 + len:08 + code:d40f + data:96 1f 32 30>
<flag:c1 + len:08 + code:d40f + data:a7 1a 2c 10>

2019/272 01:24:08.487180    9412 → 56900 [ACK] Seq=456 Ack=6384 Win=64580 Len=0

Frame 732: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6384, Len: 0

2019/272 01:24:08.535998    56900 → 9412 [PSH, ACK] Seq=6384 Ack=456 Win=65024 Len=16

Frame 735: 70 bytes on wire (560 bits), 70 bytes captured (560 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6384, Ack: 456, Len: 16
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>

2019/272 01:24:08.737189    9412 → 56900 [ACK] Seq=456 Ack=6400 Win=64564 Len=0

Frame 738: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6400, Len: 0

2019/272 01:24:08.738287    56900 → 9412 [PSH, ACK] Seq=6400 Ack=456 Win=65024 Len=19

Frame 740: 73 bytes on wire (584 bits), 73 bytes captured (584 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6400, Ack: 456, Len: 19
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>

2019/272 01:24:08.939196    9412 → 56900 [ACK] Seq=456 Ack=6419 Win=64544 Len=0

Frame 745: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6419, Len: 0

2019/272 01:24:09.069557    56900 → 9412 [PSH, ACK] Seq=6419 Ack=456 Win=65024 Len=6

Frame 748: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6419, Ack: 456, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:09.339215    9412 → 56900 [ACK] Seq=456 Ack=6425 Win=64540 Len=0

Frame 751: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6425, Len: 0

2019/272 01:24:09.340296    56900 → 9412 [PSH, ACK] Seq=6425 Ack=456 Win=65024 Len=136

Frame 753: 190 bytes on wire (1520 bits), 190 bytes captured (1520 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6425, Ack: 456, Len: 136
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1d 37 50>
<flag:c1 + len:08 + code:d40f + data:94 16 32 10>
<flag:c1 + len:08 + code:d40f + data:96 19 2f 10>
<flag:c1 + len:08 + code:d40f + data:a7 19 2d 00>
<flag:c1 + len:10 + code:26ff + data:06 78 00 37 fe 05 2e 40 78 06 ...>

2019/272 01:24:09.539237    9412 → 56900 [ACK] Seq=456 Ack=6561 Win=64404 Len=0

Frame 757: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6561, Len: 0

2019/272 01:24:09.542304    56900 → 9412 [PSH, ACK] Seq=6561 Ack=456 Win=65024 Len=16

Frame 759: 70 bytes on wire (560 bits), 70 bytes captured (560 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6561, Ack: 456, Len: 16
<flag:c1 + len:10 + code:26ff + data:06 78 00 37 fe cd 8b e9 78 06 ...>

2019/272 01:24:09.739253    9412 → 56900 [ACK] Seq=456 Ack=6577 Win=64388 Len=0

Frame 763: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6577, Len: 0

2019/272 01:24:10.046737    56900 → 9412 [PSH, ACK] Seq=6577 Ack=456 Win=65024 Len=32

Frame 766: 86 bytes on wire (688 bits), 86 bytes captured (688 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6577, Ack: 456, Len: 32
<flag:c1 + len:20 + code:2d00 + data:00 00 1c 00 01 00 00 00 00 00 ...>

2019/272 01:24:10.339256    9412 → 56900 [ACK] Seq=456 Ack=6609 Win=64356 Len=0

Frame 769: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6609, Len: 0

2019/272 01:24:10.340037    56900 → 9412 [PSH, ACK] Seq=6609 Ack=456 Win=65024 Len=133

Frame 771: 187 bytes on wire (1496 bits), 187 bytes captured (1496 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6609, Ack: 456, Len: 133
<flag:c1 + len:07 + code:0700 + data:2e 7d 38>
<flag:c1 + len:0c + code:ec30 + data:5f 00 00 00 5a 00 00 00>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:08 + code:d40f + data:88 19 35 70>
<flag:c1 + len:08 + code:d40f + data:94 1a 35 50>
<flag:c1 + len:08 + code:d40f + data:96 17 33 70>
<flag:c1 + len:08 + code:d40f + data:a7 1b 29 10>

2019/272 01:24:10.539268    9412 → 56900 [ACK] Seq=456 Ack=6742 Win=65700 Len=0

Frame 775: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6742, Len: 0

2019/272 01:24:10.540037    56900 → 9412 [PSH, ACK] Seq=6742 Ack=456 Win=65024 Len=98

Frame 777: 152 bytes on wire (1216 bits), 152 bytes captured (1216 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6742, Ack: 456, Len: 98
<flag:c1 + len:20 + code:2d02 + data:1c 00 14 00 00 92 c8 19 0a 00 ...>
<flag:c1 + len:07 + code:0701 + data:2e 7d 38>
<flag:c1 + len:0c + code:ec30 + data:5f 00 00 00 5a 00 00 00>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:10.739270    9412 → 56900 [ACK] Seq=456 Ack=6840 Win=65600 Len=0

Frame 781: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6840, Len: 0

2019/272 01:24:11.044884    56900 → 9412 [PSH, ACK] Seq=6840 Ack=456 Win=65024 Len=6

Frame 784: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6840, Ack: 456, Len: 6
<flag:c1 + len:06 + code:1401 + data:0f a7>

2019/272 01:24:11.311979    9412 → 56900 [PSH, ACK] Seq=456 Ack=6846 Win=65596 Len=12

Frame 787: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 456, Ack: 6846, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:11.313407    56900 → 9412 [PSH, ACK] Seq=6846 Ack=468 Win=65024 Len=118

Frame 789: 172 bytes on wire (1376 bits), 172 bytes captured (1376 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6846, Ack: 468, Len: 118
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 18 33 70>
<flag:c1 + len:08 + code:d40f + data:94 15 36 70>
<flag:c1 + len:08 + code:d40f + data:96 17 30 00>

2019/272 01:24:11.539307    9412 → 56900 [ACK] Seq=468 Ack=6964 Win=65476 Len=0

Frame 793: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 6964, Len: 0

2019/272 01:24:12.071663    56900 → 9412 [PSH, ACK] Seq=6964 Ack=468 Win=65024 Len=6

Frame 796: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6964, Ack: 468, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:12.339334    9412 → 56900 [ACK] Seq=468 Ack=6970 Win=65472 Len=0

Frame 800: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 6970, Len: 0

2019/272 01:24:12.340397    56900 → 9412 [PSH, ACK] Seq=6970 Ack=468 Win=65024 Len=116

Frame 802: 170 bytes on wire (1360 bits), 170 bytes captured (1360 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 6970, Ack: 468, Len: 116
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 18 30 10>
<flag:c1 + len:08 + code:d40f + data:94 15 35 00>
<flag:c1 + len:08 + code:d40f + data:96 1b 2f 10>
<flag:c1 + len:04 + code:0f04 + data:null>

2019/272 01:24:12.539341    9412 → 56900 [ACK] Seq=468 Ack=7086 Win=65356 Len=0

Frame 806: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7086, Len: 0

2019/272 01:24:12.541458    56900 → 9412 [PSH, ACK] Seq=7086 Ack=468 Win=65024 Len=16

Frame 808: 70 bytes on wire (560 bits), 70 bytes captured (560 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7086, Ack: 468, Len: 16
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>

2019/272 01:24:12.739350    9412 → 56900 [ACK] Seq=468 Ack=7102 Win=65340 Len=0

Frame 814: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7102, Len: 0

2019/272 01:24:12.741084    56900 → 9412 [PSH, ACK] Seq=7102 Ack=468 Win=65024 Len=19

Frame 816: 73 bytes on wire (584 bits), 73 bytes captured (584 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7102, Ack: 468, Len: 19
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>

2019/272 01:24:12.942361    9412 → 56900 [ACK] Seq=468 Ack=7121 Win=65320 Len=0

Frame 820: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7121, Len: 0

2019/272 01:24:13.045729    56900 → 9412 [PSH, ACK] Seq=7121 Ack=468 Win=65024 Len=6

Frame 823: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7121, Ack: 468, Len: 6
<flag:c1 + len:06 + code:1401 + data:0f 94>

2019/272 01:24:13.322368    9412 → 56900 [ACK] Seq=468 Ack=7127 Win=65312 Len=0

Frame 826: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7127, Len: 0

2019/272 01:24:13.323475    56900 → 9412 [PSH, ACK] Seq=7127 Ack=468 Win=65024 Len=98

Frame 828: 152 bytes on wire (1216 bits), 152 bytes captured (1216 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7127, Ack: 468, Len: 98
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:08 + code:d40f + data:88 1b 33 30>
<flag:c1 + len:08 + code:d40f + data:96 19 2e 10>

2019/272 01:24:13.522376    9412 → 56900 [ACK] Seq=468 Ack=7225 Win=65216 Len=0

Frame 832: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7225, Len: 0

2019/272 01:24:14.044088    56900 → 9412 [PSH, ACK] Seq=7225 Ack=468 Win=65024 Len=47

Frame 835: 101 bytes on wire (808 bits), 101 bytes captured (808 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7225, Ack: 468, Len: 47
<flag:c2 + len:002f + code:1302 + data:0f 94 01 bd 16 35 18 35 30 00 ...>

2019/272 01:24:14.252404    9412 → 56900 [ACK] Seq=468 Ack=7272 Win=65168 Len=0

Frame 838: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7272, Len: 0

2019/272 01:24:14.253474    56900 → 9412 [PSH, ACK] Seq=7272 Ack=468 Win=65024 Len=126

Frame 840: 180 bytes on wire (1440 bits), 180 bytes captured (1440 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7272, Ack: 468, Len: 126
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 19 35 50>
<flag:c1 + len:08 + code:d40f + data:94 17 35 30>
<flag:c1 + len:08 + code:d40f + data:96 1c 2e 30>
<flag:c1 + len:08 + code:d40f + data:a7 18 2e 60>

2019/272 01:24:14.452423    9412 → 56900 [ACK] Seq=468 Ack=7398 Win=65044 Len=0

Frame 844: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7398, Len: 0

2019/272 01:24:14.536897    56900 → 9412 [PSH, ACK] Seq=7398 Ack=468 Win=65024 Len=6

Frame 847: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7398, Ack: 468, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:14.812427    9412 → 56900 [ACK] Seq=468 Ack=7404 Win=65036 Len=0

Frame 850: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7404, Len: 0

2019/272 01:24:14.813523    56900 → 9412 [PSH, ACK] Seq=7404 Ack=468 Win=65024 Len=41

Frame 852: 95 bytes on wire (760 bits), 95 bytes captured (760 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7404, Ack: 468, Len: 41
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:15.012433    9412 → 56900 [ACK] Seq=468 Ack=7445 Win=64996 Len=0

Frame 856: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7445, Len: 0

2019/272 01:24:15.067728    56900 → 9412 [PSH, ACK] Seq=7445 Ack=468 Win=65024 Len=6

Frame 859: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7445, Ack: 468, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:15.312120    9412 → 56900 [PSH, ACK] Seq=468 Ack=7451 Win=64988 Len=12

Frame 862: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 468, Ack: 7451, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:15.312917    56900 → 9412 [PSH, ACK] Seq=7451 Ack=480 Win=65024 Len=210

Frame 864: 264 bytes on wire (2112 bits), 264 bytes captured (2112 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7451, Ack: 480, Len: 210
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1c 34 30>
<flag:c1 + len:08 + code:d40f + data:94 18 34 20>
<flag:c1 + len:08 + code:d40f + data:96 18 31 60>
<flag:c1 + len:08 + code:d40f + data:a7 18 2c 10>
<flag:c1 + len:10 + code:110f + data:95 02 58 00 04 00 00 40 58 02 ...>
<flag:c1 + len:12 + code:fa05 + data:95 0f bf c8 00 00 20 cb 00 00 ...>
<flag:c1 + len:12 + code:fa05 + data:95 0f bf c8 00 00 20 cb 00 00 ...>
<flag:c1 + len:0e + code:ec10 + data:0f 95 00 cb 00 20 00 c8 00 bf>
<flag:c1 + len:06 + code:2a07 + data:fb 00>
<flag:c1 + len:06 + code:2a0a + data:41 00>
<flag:c1 + len:06 + code:2a0b + data:46 00>
<flag:c1 + len:06 + code:2a09 + data:41 00>

2019/272 01:24:15.512452    9412 → 56900 [ACK] Seq=480 Ack=7661 Win=64780 Len=0

Frame 868: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 7661, Len: 0

2019/272 01:24:16.044704    56900 → 9412 [PSH, ACK] Seq=7661 Ack=480 Win=65024 Len=6

Frame 871: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7661, Ack: 480, Len: 6
<flag:c1 + len:06 + code:1401 + data:0f a7>

2019/272 01:24:16.312485    9412 → 56900 [ACK] Seq=480 Ack=7667 Win=64772 Len=0

Frame 874: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 7667, Len: 0

2019/272 01:24:16.313612    56900 → 9412 [PSH, ACK] Seq=7667 Ack=480 Win=65024 Len=122

Frame 876: 176 bytes on wire (1408 bits), 176 bytes captured (1408 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7667, Ack: 480, Len: 122
<flag:c1 + len:10 + code:26ff + data:06 78 00 37 fe 05 2e 40 78 06 ...>
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 16 30 00>
<flag:c1 + len:08 + code:d40f + data:94 1a 32 20>
<flag:c1 + len:08 + code:d40f + data:96 1d 32 30>

2019/272 01:24:16.512501    9412 → 56900 [ACK] Seq=480 Ack=7789 Win=64652 Len=0

Frame 880: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 7789, Len: 0

2019/272 01:24:16.536152    56900 → 9412 [PSH, ACK] Seq=7789 Ack=480 Win=65024 Len=16

Frame 883: 70 bytes on wire (560 bits), 70 bytes captured (560 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7789, Ack: 480, Len: 16
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>

2019/272 01:24:16.812503    9412 → 56900 [ACK] Seq=480 Ack=7805 Win=64636 Len=0

Frame 886: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 7805, Len: 0

2019/272 01:24:16.815585    56900 → 9412 [PSH, ACK] Seq=7805 Ack=480 Win=65024 Len=19

Frame 888: 73 bytes on wire (584 bits), 73 bytes captured (584 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7805, Ack: 480, Len: 19
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>

2019/272 01:24:17.112514    9412 → 56900 [ACK] Seq=480 Ack=7824 Win=64616 Len=0

Frame 892: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 7824, Len: 0

2019/272 01:24:17.113611    56900 → 9412 [PSH, ACK] Seq=7824 Ack=480 Win=65024 Len=218

Frame 894: 272 bytes on wire (2176 bits), 272 bytes captured (2176 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 7824, Ack: 480, Len: 218
<flag:c2 + len:001a + code:1301 + data:0f a7 01 be 17 2e 14 2f 60 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1a 37 50>
<flag:c1 + len:08 + code:d40f + data:94 1b 36 50>
<flag:c1 + len:08 + code:d40f + data:96 18 30 70>
<flag:c1 + len:08 + code:d40f + data:a7 16 2c 10>
<flag:c1 + len:10 + code:110f + data:85 02 71 00 04 00 00 40 71 02 ...>
<flag:c1 + len:12 + code:fa05 + data:85 0f 31 db 00 00 90 e2 00 00 ...>
<flag:c1 + len:12 + code:fa05 + data:85 0f 31 db 00 00 90 e2 00 00 ...>
<flag:c1 + len:0e + code:ec10 + data:0f 85 00 e2 00 90 00 db 00 31>

2019/272 01:24:17.312527    9412 → 56900 [ACK] Seq=480 Ack=8042 Win=64400 Len=0

Frame 898: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 8042, Len: 0

2019/272 01:24:18.044071    56900 → 9412 [PSH, ACK] Seq=8042 Ack=480 Win=65024 Len=6

Frame 901: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8042, Ack: 480, Len: 6
<flag:c1 + len:06 + code:1401 + data:0f a7>

2019/272 01:24:18.312562    9412 → 56900 [ACK] Seq=480 Ack=8048 Win=64392 Len=0

Frame 904: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 8048, Len: 0

2019/272 01:24:18.313339    56900 → 9412 [PSH, ACK] Seq=8048 Ack=480 Win=65024 Len=118

Frame 906: 172 bytes on wire (1376 bits), 172 bytes captured (1376 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8048, Ack: 480, Len: 118
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 16 35 70>
<flag:c1 + len:08 + code:d40f + data:94 18 34 70>
<flag:c1 + len:08 + code:d40f + data:96 16 33 60>

2019/272 01:24:18.512570    9412 → 56900 [ACK] Seq=480 Ack=8166 Win=64276 Len=0

Frame 910: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 8166, Len: 0

2019/272 01:24:18.536517    56900 → 9412 [PSH, ACK] Seq=8166 Ack=480 Win=65024 Len=6

Frame 913: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8166, Ack: 480, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:18.812582    9412 → 56900 [ACK] Seq=480 Ack=8172 Win=64268 Len=0

Frame 916: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 8172, Len: 0

2019/272 01:24:18.813414    56900 → 9412 [PSH, ACK] Seq=8172 Ack=480 Win=65024 Len=53

Frame 918: 107 bytes on wire (856 bits), 107 bytes captured (856 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8172, Ack: 480, Len: 53
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:0c + code:27ff + data:02 f8 03 0d f8 02 00 00>

2019/272 01:24:19.012592    9412 → 56900 [ACK] Seq=480 Ack=8225 Win=65700 Len=0

Frame 923: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 8225, Len: 0

2019/272 01:24:19.068754    56900 → 9412 [PSH, ACK] Seq=8225 Ack=480 Win=65024 Len=16

Frame 926: 70 bytes on wire (560 bits), 70 bytes captured (560 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8225, Ack: 480, Len: 16
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>

2019/272 01:24:19.272599    9412 → 56900 [ACK] Seq=480 Ack=8241 Win=65684 Len=0

Frame 929: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 8241, Len: 0

2019/272 01:24:19.273704    56900 → 9412 [PSH, ACK] Seq=8241 Ack=480 Win=65024 Len=90

Frame 931: 144 bytes on wire (1152 bits), 144 bytes captured (1152 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8241, Ack: 480, Len: 90
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1b 37 40>
<flag:c1 + len:08 + code:d40f + data:94 17 35 70>
<flag:c1 + len:08 + code:d40f + data:96 19 2e 10>

2019/272 01:24:19.312275    9412 → 56900 [PSH, ACK] Seq=480 Ack=8331 Win=65592 Len=12

Frame 935: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 480, Ack: 8331, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:19.574614    56900 → 9412 [ACK] Seq=8331 Ack=492 Win=65024 Len=0

Frame 938: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8331, Ack: 492, Len: 0

2019/272 01:24:20.069588    56900 → 9412 [PSH, ACK] Seq=8331 Ack=492 Win=65024 Len=6

Frame 941: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8331, Ack: 492, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:20.274637    9412 → 56900 [ACK] Seq=492 Ack=8337 Win=65588 Len=0

Frame 944: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8337, Len: 0

2019/272 01:24:20.275406    56900 → 9412 [PSH, ACK] Seq=8337 Ack=492 Win=65024 Len=155

Frame 946: 209 bytes on wire (1672 bits), 209 bytes captured (1672 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8337, Ack: 492, Len: 155
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1d 35 30>
<flag:c1 + len:08 + code:d40f + data:94 19 32 10>
<flag:c1 + len:20 + code:2da1 + data:00 00 1c 00 01 b5 1a 40 00 00 ...>
<flag:c1 + len:07 + code:0700 + data:2e 7d 38>
<flag:c1 + len:0c + code:ec30 + data:5f 00 00 00 5a 00 00 00>

2019/272 01:24:20.474648    9412 → 56900 [ACK] Seq=492 Ack=8492 Win=65432 Len=0

Frame 950: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8492, Len: 0

2019/272 01:24:20.536924    56900 → 9412 [PSH, ACK] Seq=8492 Ack=492 Win=65024 Len=32

Frame 953: 86 bytes on wire (688 bits), 86 bytes captured (688 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8492, Ack: 492, Len: 32
<flag:c1 + len:20 + code:2d40 + data:1c 00 14 00 00 5f a6 40 0a 00 ...>

2019/272 01:24:20.774657    9412 → 56900 [ACK] Seq=492 Ack=8524 Win=65400 Len=0

Frame 956: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8524, Len: 0

2019/272 01:24:20.782863    56900 → 9412 [PSH, ACK] Seq=8524 Ack=492 Win=65024 Len=66

Frame 958: 120 bytes on wire (960 bits), 120 bytes captured (960 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8524, Ack: 492, Len: 66
<flag:c1 + len:07 + code:0701 + data:2e 7d 38>
<flag:c1 + len:0c + code:ec30 + data:5f 00 00 00 5a 00 00 00>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:21.074688    9412 → 56900 [ACK] Seq=492 Ack=8590 Win=65332 Len=0

Frame 962: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8590, Len: 0

2019/272 01:24:21.075786    56900 → 9412 [PSH, ACK] Seq=8590 Ack=492 Win=65024 Len=6

Frame 964: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8590, Ack: 492, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:21.374679    9412 → 56900 [ACK] Seq=492 Ack=8596 Win=65328 Len=0

Frame 968: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8596, Len: 0

2019/272 01:24:21.375772    56900 → 9412 [PSH, ACK] Seq=8596 Ack=492 Win=65024 Len=128

Frame 970: 182 bytes on wire (1456 bits), 182 bytes captured (1456 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8596, Ack: 492, Len: 128
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1e 35 50>
<flag:c1 + len:08 + code:d40f + data:94 1a 30 10>
<flag:c1 + len:08 + code:d40f + data:96 18 2c 10>
<flag:c1 + len:10 + code:26ff + data:06 78 00 37 fe ff 00 00 78 06 ...>

2019/272 01:24:21.674692    9412 → 56900 [ACK] Seq=492 Ack=8724 Win=65200 Len=0

Frame 974: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8724, Len: 0

2019/272 01:24:22.044591    56900 → 9412 [PSH, ACK] Seq=8724 Ack=492 Win=65024 Len=6

Frame 977: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8724, Ack: 492, Len: 6
<flag:c1 + len:06 + code:1401 + data:0f 96>

2019/272 01:24:22.255729    9412 → 56900 [ACK] Seq=492 Ack=8730 Win=65192 Len=0

Frame 991: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8730, Len: 0

2019/272 01:24:22.256518    56900 → 9412 [PSH, ACK] Seq=8730 Ack=492 Win=65024 Len=71

Frame 993: 125 bytes on wire (1000 bits), 125 bytes captured (1000 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8730, Ack: 492, Len: 71
<flag:c1 + len:08 + code:d40f + data:86 25 35 50>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1e 36 30>
<flag:c1 + len:08 + code:d40f + data:94 16 34 70>

2019/272 01:24:22.555727    9412 → 56900 [ACK] Seq=492 Ack=8801 Win=65124 Len=0

Frame 999: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8801, Len: 0

2019/272 01:24:22.556806    56900 → 9412 [PSH, ACK] Seq=8801 Ack=492 Win=65024 Len=47

Frame 1001: 101 bytes on wire (808 bits), 101 bytes captured (808 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8801, Ack: 492, Len: 47
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:22.855758    9412 → 56900 [ACK] Seq=492 Ack=8848 Win=65076 Len=0

Frame 1005: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8848, Len: 0

2019/272 01:24:23.045676    56900 → 9412 [PSH, ACK] Seq=8848 Ack=492 Win=65024 Len=26

Frame 1008: 80 bytes on wire (640 bits), 80 bytes captured (640 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8848, Ack: 492, Len: 26
<flag:c2 + len:001a + code:1301 + data:0f a7 01 be 17 2f 19 30 40 00 ...>

2019/272 01:24:23.255759    9412 → 56900 [ACK] Seq=492 Ack=8874 Win=65048 Len=0

Frame 1011: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 8874, Len: 0

2019/272 01:24:23.256867    56900 → 9412 [PSH, ACK] Seq=8874 Ack=492 Win=65024 Len=146

Frame 1013: 200 bytes on wire (1600 bits), 200 bytes captured (1600 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 8874, Ack: 492, Len: 146
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1a 37 60>
<flag:c1 + len:08 + code:d40f + data:94 18 30 10>
<flag:c1 + len:08 + code:d40f + data:a7 1a 2e 30>
<flag:c1 + len:0c + code:27ff + data:02 f8 03 0d f8 02 00 00>
<flag:c1 + len:10 + code:26ff + data:06 78 00 37 fe 05 2e 40 78 06 ...>

2019/272 01:24:23.312461    9412 → 56900 [PSH, ACK] Seq=492 Ack=9020 Win=64904 Len=12

Frame 1017: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 492, Ack: 9020, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:23.555765    56900 → 9412 [ACK] Seq=9020 Ack=504 Win=65024 Len=0

Frame 1020: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9020, Ack: 504, Len: 0

2019/272 01:24:24.044132    56900 → 9412 [PSH, ACK] Seq=9020 Ack=504 Win=65024 Len=26

Frame 1023: 80 bytes on wire (640 bits), 80 bytes captured (640 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9020, Ack: 504, Len: 26
<flag:c2 + len:001a + code:1301 + data:0f 96 01 bd 19 2d 19 2d 50 00 ...>

2019/272 01:24:24.255791    9412 → 56900 [ACK] Seq=504 Ack=9046 Win=64876 Len=0

Frame 1026: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 504, Ack: 9046, Len: 0

2019/272 01:24:24.256863    56900 → 9412 [PSH, ACK] Seq=9046 Ack=504 Win=65024 Len=79

Frame 1028: 133 bytes on wire (1064 bits), 133 bytes captured (1064 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9046, Ack: 504, Len: 79
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1e 36 30>
<flag:c1 + len:08 + code:d40f + data:94 19 2f 10>
<flag:c1 + len:08 + code:d40f + data:96 18 2d 70>
<flag:c1 + len:08 + code:d40f + data:a7 17 2d 00>

2019/272 01:24:24.555801    9412 → 56900 [ACK] Seq=504 Ack=9125 Win=64800 Len=0

Frame 1032: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 504, Ack: 9125, Len: 0

2019/272 01:24:24.556907    56900 → 9412 [PSH, ACK] Seq=9125 Ack=504 Win=65024 Len=47

Frame 1034: 101 bytes on wire (808 bits), 101 bytes captured (808 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9125, Ack: 504, Len: 47
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:24.855815    9412 → 56900 [ACK] Seq=504 Ack=9172 Win=64752 Len=0

Frame 1038: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 504, Ack: 9172, Len: 0

2019/272 01:24:25.045399    56900 → 9412 [PSH, ACK] Seq=9172 Ack=504 Win=65024 Len=8

Frame 1041: 62 bytes on wire (496 bits), 62 bytes captured (496 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9172, Ack: 504, Len: 8
<flag:c1 + len:08 + code:1402 + data:0f a7 0f 96>

2019/272 01:24:25.255830    9412 → 56900 [ACK] Seq=504 Ack=9180 Win=64744 Len=0

Frame 1044: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 504, Ack: 9180, Len: 0

2019/272 01:24:25.257896    56900 → 9412 [PSH, ACK] Seq=9180 Ack=504 Win=65024 Len=134

Frame 1046: 188 bytes on wire (1504 bits), 188 bytes captured (1504 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9180, Ack: 504, Len: 134
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1a 34 00>
<flag:c1 + len:08 + code:d40f + data:94 16 2e 10>
<flag:c1 + len:06 + code:2a07 + data:fb 00>
<flag:c1 + len:06 + code:2a0a + data:41 00>
<flag:c1 + len:06 + code:2a0b + data:46 00>
<flag:c1 + len:06 + code:2a09 + data:41 00>

2019/272 01:24:25.312586    9412 → 56900 [PSH, ACK] Seq=504 Ack=9314 Win=64608 Len=9

Frame 1050: 63 bytes on wire (504 bits), 63 bytes captured (504 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 504, Ack: 9314, Len: 9
<flag:c1 + len:09 + code:fa11 + data:8b ff 55 8b ec>

2019/272 01:24:25.555843    56900 → 9412 [ACK] Seq=9314 Ack=513 Win=65024 Len=0

Frame 1053: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9314, Ack: 513, Len: 0

2019/272 01:24:26.044109    56900 → 9412 [PSH, ACK] Seq=9314 Ack=513 Win=65024 Len=26

Frame 1057: 80 bytes on wire (640 bits), 80 bytes captured (640 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9314, Ack: 513, Len: 26
<flag:c2 + len:001a + code:1301 + data:0f 96 01 bd 19 2e 19 2f 40 00 ...>

2019/272 01:24:26.255872    9412 → 56900 [ACK] Seq=513 Ack=9340 Win=64584 Len=0

Frame 1060: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 513, Ack: 9340, Len: 0

2019/272 01:24:26.257994    56900 → 9412 [PSH, ACK] Seq=9340 Ack=513 Win=65024 Len=106

Frame 1062: 160 bytes on wire (1280 bits), 160 bytes captured (1280 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9340, Ack: 513, Len: 106
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1c 33 10>
<flag:c1 + len:08 + code:d40f + data:94 15 32 70>
<flag:c1 + len:08 + code:d40f + data:96 1b 31 50>

2019/272 01:24:26.312919    9412 → 56900 [PSH, ACK] Seq=513 Ack=9446 Win=64476 Len=12

Frame 1066: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 513, Ack: 9446, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:26.535952    56900 → 9412 [PSH, ACK] Seq=9446 Ack=525 Win=65024 Len=6

Frame 1069: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9446, Ack: 525, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:26.657259    9412 → 56900 [PSH, ACK] Seq=525 Ack=9452 Win=64472 Len=19

Frame 1072: 73 bytes on wire (584 bits), 73 bytes captured (584 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 525, Ack: 9452, Len: 19
<flag:c3 + len:13><flag:c1 + len:0d + code:4012 + data:d6 b1 a9 8a 57 4b 5f 43 5f>

2019/272 01:24:26.658077    56900 → 9412 [PSH, ACK] Seq=9452 Ack=544 Win=65024 Len=41

Frame 1074: 95 bytes on wire (760 bits), 95 bytes captured (760 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9452, Ack: 544, Len: 41
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:26.955905    9412 → 56900 [ACK] Seq=544 Ack=9493 Win=64432 Len=0

Frame 1078: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9493, Len: 0

2019/272 01:24:26.956979    56900 → 9412 [PSH, ACK] Seq=9493 Ack=544 Win=65024 Len=4

Frame 1080: 58 bytes on wire (464 bits), 58 bytes captured (464 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9493, Ack: 544, Len: 4
<flag:c1 + len:04 + code:4103 + data:null>

2019/272 01:24:27.255912    9412 → 56900 [ACK] Seq=544 Ack=9497 Win=64428 Len=0

Frame 1084: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9497, Len: 0

2019/272 01:24:27.259988    56900 → 9412 [PSH, ACK] Seq=9497 Ack=544 Win=65024 Len=118

Frame 1086: 172 bytes on wire (1376 bits), 172 bytes captured (1376 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9497, Ack: 544, Len: 118
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1c 32 30>
<flag:c1 + len:08 + code:d40f + data:94 1a 34 30>
<flag:c1 + len:08 + code:d40f + data:96 18 2b 10>

2019/272 01:24:27.555922    9412 → 56900 [ACK] Seq=544 Ack=9615 Win=64308 Len=0

Frame 1090: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9615, Len: 0

2019/272 01:24:28.072782    56900 → 9412 [PSH, ACK] Seq=9615 Ack=544 Win=65024 Len=6

Frame 1093: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9615, Ack: 544, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:28.355983    9412 → 56900 [ACK] Seq=544 Ack=9621 Win=64304 Len=0

Frame 1096: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9621, Len: 0

2019/272 01:24:28.356752    56900 → 9412 [PSH, ACK] Seq=9621 Ack=544 Win=65024 Len=116

Frame 1098: 170 bytes on wire (1360 bits), 170 bytes captured (1360 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9621, Ack: 544, Len: 116
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1f 33 30>
<flag:c1 + len:08 + code:d40f + data:94 16 37 50>
<flag:c1 + len:08 + code:d40f + data:96 1a 2f 40>
<flag:c1 + len:04 + code:0f27 + data:null>

2019/272 01:24:28.557958    9412 → 56900 [ACK] Seq=544 Ack=9737 Win=65700 Len=0

Frame 1102: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9737, Len: 0

2019/272 01:24:28.559015    56900 → 9412 [PSH, ACK] Seq=9737 Ack=544 Win=65024 Len=35

Frame 1104: 89 bytes on wire (712 bits), 89 bytes captured (712 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9737, Ack: 544, Len: 35
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>

2019/272 01:24:28.757995    9412 → 56900 [ACK] Seq=544 Ack=9772 Win=65664 Len=0

Frame 1108: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9772, Len: 0

2019/272 01:24:29.068992    56900 → 9412 [PSH, ACK] Seq=9772 Ack=544 Win=65024 Len=6

Frame 1111: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9772, Ack: 544, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:29.357996    9412 → 56900 [ACK] Seq=544 Ack=9778 Win=65656 Len=0

Frame 1114: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9778, Len: 0

2019/272 01:24:29.359060    56900 → 9412 [PSH, ACK] Seq=9778 Ack=544 Win=65024 Len=100

Frame 1116: 154 bytes on wire (1232 bits), 154 bytes captured (1232 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9778, Ack: 544, Len: 100
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:08 + code:d40f + data:88 1f 34 30>
<flag:c1 + len:08 + code:d40f + data:94 18 30 10>
<flag:c1 + len:08 + code:d40f + data:96 1b 30 40>

2019/272 01:24:29.557996    9412 → 56900 [ACK] Seq=544 Ack=9878 Win=65556 Len=0

Frame 1120: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9878, Len: 0

2019/272 01:24:29.559102    56900 → 9412 [PSH, ACK] Seq=9878 Ack=544 Win=65024 Len=12

Frame 1122: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9878, Ack: 544, Len: 12
<flag:c1 + len:0c + code:27ff + data:02 f8 03 0d f8 02 00 00>

2019/272 01:24:29.758027    9412 → 56900 [ACK] Seq=544 Ack=9890 Win=65544 Len=0

Frame 1126: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9890, Len: 0

2019/272 01:24:30.069043    56900 → 9412 [PSH, ACK] Seq=9890 Ack=544 Win=65024 Len=16

Frame 1129: 70 bytes on wire (560 bits), 70 bytes captured (560 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9890, Ack: 544, Len: 16
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>

2019/272 01:24:30.313717    9412 → 56900 [PSH, ACK] Seq=544 Ack=9906 Win=65528 Len=12

Frame 1132: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 544, Ack: 9906, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:30.314485    56900 → 9412 [PSH, ACK] Seq=9906 Ack=556 Win=65024 Len=157

Frame 1134: 211 bytes on wire (1688 bits), 211 bytes captured (1688 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 9906, Ack: 556, Len: 157
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 19 33 70>
<flag:c1 + len:08 + code:d40f + data:94 19 2f 10>
<flag:c1 + len:08 + code:d40f + data:96 18 2c 10>
<flag:c1 + len:20 + code:2da1 + data:00 00 1c 00 01 b5 1a 40 00 00 ...>
<flag:c1 + len:07 + code:0700 + data:2e 7d 38>
<flag:c1 + len:0c + code:ec30 + data:5f 00 00 00 5a 00 00 00>
<flag:c1 + len:10 + code:26ff + data:06 78 00 37 fe 98 4a 40 78 06 ...>

2019/272 01:24:30.558041    9412 → 56900 [ACK] Seq=556 Ack=10063 Win=65372 Len=0

Frame 1138: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 556, Ack: 10063, Len: 0

2019/272 01:24:30.559016    56900 → 9412 [PSH, ACK] Seq=10063 Ack=556 Win=65024 Len=98

Frame 1140: 152 bytes on wire (1216 bits), 152 bytes captured (1216 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10063, Ack: 556, Len: 98
<flag:c1 + len:20 + code:2d00 + data:1c 00 14 00 00 00 00 00 0a 00 ...>
<flag:c1 + len:07 + code:0701 + data:2e 7d 38>
<flag:c1 + len:0c + code:ec30 + data:5f 00 00 00 5a 00 00 00>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:30.758045    9412 → 56900 [ACK] Seq=556 Ack=10161 Win=65276 Len=0

Frame 1144: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 556, Ack: 10161, Len: 0

2019/272 01:24:31.060619    56900 → 9412 [PSH, ACK] Seq=10161 Ack=556 Win=65024 Len=6

Frame 1147: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10161, Ack: 556, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:31.258065    9412 → 56900 [ACK] Seq=556 Ack=10167 Win=65268 Len=0

Frame 1150: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 556, Ack: 10167, Len: 0

2019/272 01:24:31.259184    56900 → 9412 [PSH, ACK] Seq=10167 Ack=556 Win=65024 Len=100

Frame 1152: 154 bytes on wire (1232 bits), 154 bytes captured (1232 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10167, Ack: 556, Len: 100
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:08 + code:d40f + data:88 19 30 00>
<flag:c1 + len:08 + code:d40f + data:94 15 34 70>
<flag:c1 + len:08 + code:d40f + data:96 19 2d 10>

2019/272 01:24:31.458069    9412 → 56900 [ACK] Seq=556 Ack=10267 Win=65168 Len=0

Frame 1156: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 556, Ack: 10267, Len: 0

2019/272 01:24:32.075144    56900 → 9412 [PSH, ACK] Seq=10267 Ack=556 Win=65024 Len=6

Frame 1159: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10267, Ack: 556, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:32.358109    9412 → 56900 [ACK] Seq=556 Ack=10273 Win=65164 Len=0

Frame 1162: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 556, Ack: 10273, Len: 0

2019/272 01:24:32.359219    56900 → 9412 [PSH, ACK] Seq=10273 Ack=556 Win=65024 Len=112

Frame 1164: 166 bytes on wire (1328 bits), 166 bytes captured (1328 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10273, Ack: 556, Len: 112
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1b 2e 10>
<flag:c1 + len:08 + code:d40f + data:94 17 2f 10>
<flag:c1 + len:08 + code:d40f + data:96 1b 32 50>

2019/272 01:24:32.422245    9412 → 56900 [PSH, ACK] Seq=556 Ack=10385 Win=65052 Len=11

Frame 1168: 65 bytes on wire (520 bits), 65 bytes captured (520 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 556, Ack: 10385, Len: 11
<flag:c3 + len:0b>
Lua Error: dissect_tcp_pdus dissect_func: C:\Program Files\Wireshark\plugins\mu.lua:112: decrypt failed! src len invalid

2019/272 01:24:32.535862    56900 → 9412 [PSH, ACK] Seq=10385 Ack=567 Win=65024 Len=6

Frame 1171: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10385, Ack: 567, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:32.758124    9412 → 56900 [ACK] Seq=567 Ack=10391 Win=65044 Len=0

Frame 1174: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 567, Ack: 10391, Len: 0

2019/272 01:24:32.759190    56900 → 9412 [PSH, ACK] Seq=10391 Ack=567 Win=65024 Len=41

Frame 1176: 95 bytes on wire (760 bits), 95 bytes captured (760 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10391, Ack: 567, Len: 41
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:32.958127    9412 → 56900 [ACK] Seq=567 Ack=10432 Win=65004 Len=0

Frame 1180: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 567, Ack: 10432, Len: 0

2019/272 01:24:33.062394    56900 → 9412 [PSH, ACK] Seq=10432 Ack=567 Win=65024 Len=6

Frame 1183: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10432, Ack: 567, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:33.258160    9412 → 56900 [ACK] Seq=567 Ack=10438 Win=64996 Len=0

Frame 1186: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 567, Ack: 10438, Len: 0

2019/272 01:24:33.267235    56900 → 9412 [PSH, ACK] Seq=10438 Ack=567 Win=65024 Len=100

Frame 1188: 154 bytes on wire (1232 bits), 154 bytes captured (1232 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10438, Ack: 567, Len: 100
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:08 + code:d40f + data:88 1b 30 10>
<flag:c1 + len:08 + code:d40f + data:94 19 36 50>
<flag:c1 + len:08 + code:d40f + data:96 17 31 70>

2019/272 01:24:33.313829    9412 → 56900 [PSH, ACK] Seq=567 Ack=10538 Win=64896 Len=12

Frame 1192: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 567, Ack: 10538, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:33.558152    56900 → 9412 [ACK] Seq=10538 Ack=579 Win=65024 Len=0

Frame 1197: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10538, Ack: 579, Len: 0

2019/272 01:24:34.043588    56900 → 9412 [PSH, ACK] Seq=10538 Ack=579 Win=65024 Len=26

Frame 1200: 80 bytes on wire (640 bits), 80 bytes captured (640 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10538, Ack: 579, Len: 26
<flag:c2 + len:001a + code:1301 + data:0f a7 01 be 19 2d 19 30 40 00 ...>

2019/272 01:24:34.258187    9412 → 56900 [ACK] Seq=579 Ack=10564 Win=64872 Len=0

Frame 1203: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 10564, Len: 0

2019/272 01:24:34.272319    56900 → 9412 [PSH, ACK] Seq=10564 Ack=579 Win=65024 Len=192

Frame 1205: 246 bytes on wire (1968 bits), 246 bytes captured (1968 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10564, Ack: 579, Len: 192
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1c 32 40>
<flag:c1 + len:08 + code:d40f + data:94 19 37 70>
<flag:c1 + len:08 + code:d40f + data:96 1b 33 50>
<flag:c1 + len:08 + code:d40f + data:a7 17 2d 70>
<flag:c1 + len:10 + code:110f + data:86 02 71 00 04 00 00 40 71 02 ...>
<flag:c1 + len:12 + code:fa05 + data:86 0f bb d8 00 00 90 e2 00 00 ...>
<flag:c1 + len:12 + code:fa05 + data:86 0f bb d8 00 00 90 e2 00 00 ...>
<flag:c1 + len:0e + code:ec10 + data:0f 86 00 e2 00 90 00 d8 00 bb>

2019/272 01:24:34.558196    9412 → 56900 [ACK] Seq=579 Ack=10756 Win=64680 Len=0

Frame 1209: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 10756, Len: 0

2019/272 01:24:34.559299    56900 → 9412 [PSH, ACK] Seq=10756 Ack=579 Win=65024 Len=87

Frame 1211: 141 bytes on wire (1128 bits), 141 bytes captured (1128 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10756, Ack: 579, Len: 87
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:0c + code:27ff + data:02 f8 03 0d f8 02 00 00>
<flag:c1 + len:10 + code:26ff + data:06 78 00 37 fe 02 00 00 78 06 ...>
<flag:c1 + len:0c + code:27ff + data:02 f8 03 0d f8 02 00 00>

2019/272 01:24:34.758198    9412 → 56900 [ACK] Seq=579 Ack=10843 Win=64592 Len=0

Frame 1215: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 10843, Len: 0

2019/272 01:24:35.045214    56900 → 9412 [PSH, ACK] Seq=10843 Ack=579 Win=65024 Len=6

Frame 1218: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10843, Ack: 579, Len: 6
<flag:c1 + len:06 + code:1401 + data:0f a7>

2019/272 01:24:35.238224    9412 → 56900 [ACK] Seq=579 Ack=10849 Win=64588 Len=0

Frame 1221: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 10849, Len: 0

2019/272 01:24:35.239316    56900 → 9412 [PSH, ACK] Seq=10849 Ack=579 Win=65024 Len=142

Frame 1223: 196 bytes on wire (1568 bits), 196 bytes captured (1568 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10849, Ack: 579, Len: 142
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 17 34 70>
<flag:c1 + len:08 + code:d40f + data:94 19 35 50>
<flag:c1 + len:08 + code:d40f + data:96 1d 32 30>
<flag:c1 + len:06 + code:2a07 + data:fb 00>
<flag:c1 + len:06 + code:2a0a + data:41 00>
<flag:c1 + len:06 + code:2a0b + data:46 00>
<flag:c1 + len:06 + code:2a09 + data:41 00>

2019/272 01:24:35.438221    9412 → 56900 [ACK] Seq=579 Ack=10991 Win=64444 Len=0

Frame 1227: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 10991, Len: 0

2019/272 01:24:36.043613    56900 → 9412 [PSH, ACK] Seq=10991 Ack=579 Win=65024 Len=26

Frame 1231: 80 bytes on wire (640 bits), 80 bytes captured (640 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 10991, Ack: 579, Len: 26
<flag:c2 + len:001a + code:1301 + data:0f a7 01 be 17 2e 17 30 60 00 ...>

2019/272 01:24:36.258257    9412 → 56900 [ACK] Seq=579 Ack=11017 Win=64420 Len=0

Frame 1234: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 11017, Len: 0

2019/272 01:24:36.259323    56900 → 9412 [PSH, ACK] Seq=11017 Ack=579 Win=65024 Len=114

Frame 1236: 168 bytes on wire (1344 bits), 168 bytes captured (1344 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11017, Ack: 579, Len: 114
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 16 32 70>
<flag:c1 + len:08 + code:d40f + data:94 16 31 00>
<flag:c1 + len:08 + code:d40f + data:96 1a 2f 10>
<flag:c1 + len:08 + code:d40f + data:a7 13 2b 70>

2019/272 01:24:36.458264    9412 → 56900 [ACK] Seq=579 Ack=11131 Win=64304 Len=0

Frame 1240: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 11131, Len: 0

2019/272 01:24:36.537442    56900 → 9412 [PSH, ACK] Seq=11131 Ack=579 Win=65024 Len=6

Frame 1243: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11131, Ack: 579, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:36.758278    9412 → 56900 [ACK] Seq=579 Ack=11137 Win=64300 Len=0

Frame 1247: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 11137, Len: 0

2019/272 01:24:36.759323    56900 → 9412 [PSH, ACK] Seq=11137 Ack=579 Win=65024 Len=41

Frame 1249: 95 bytes on wire (760 bits), 95 bytes captured (760 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11137, Ack: 579, Len: 41
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 95 00 07 2e 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>

2019/272 01:24:36.958299    9412 → 56900 [ACK] Seq=579 Ack=11178 Win=64256 Len=0

Frame 1253: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 11178, Len: 0

2019/272 01:24:37.044628    56900 → 9412 [PSH, ACK] Seq=11178 Ack=579 Win=65024 Len=6

Frame 1256: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11178, Ack: 579, Len: 6
<flag:c1 + len:06 + code:1401 + data:0f a7>

2019/272 01:24:37.248322    9412 → 56900 [ACK] Seq=579 Ack=11184 Win=64252 Len=0

Frame 1259: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 11184, Len: 0

2019/272 01:24:37.249646    56900 → 9412 [PSH, ACK] Seq=11184 Ack=579 Win=65024 Len=200

Frame 1261: 254 bytes on wire (2032 bits), 254 bytes captured (2032 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11184, Ack: 579, Len: 200
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 84 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:06 + code:2a08 + data:c4 00>
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 85 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:08 + code:d40f + data:88 1b 2f 10>
<flag:c1 + len:08 + code:d40f + data:94 1a 34 30>
<flag:c1 + len:08 + code:d40f + data:96 16 31 70>
<flag:c1 + len:10 + code:110f + data:84 02 71 00 04 00 00 40 71 02 ...>
<flag:c1 + len:12 + code:fa05 + data:84 0f 0a e0 00 00 90 e2 00 00 ...>
<flag:c1 + len:12 + code:fa05 + data:84 0f 0a e0 00 00 90 e2 00 00 ...>
<flag:c1 + len:0e + code:ec10 + data:0f 84 00 e2 00 90 00 e0 00 0a>
<flag:c1 + len:10 + code:26ff + data:06 78 00 37 fe 05 2e 40 78 06 ...>

2019/272 01:24:37.314418    9412 → 56900 [PSH, ACK] Seq=579 Ack=11384 Win=65700 Len=12

Frame 1265: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 579, Ack: 11384, Len: 12
<flag:c1 + len:0c + code:fa15 + data:00 00 00 00 00 00 00 00>

2019/272 01:24:37.537589    56900 → 9412 [PSH, ACK] Seq=11384 Ack=591 Win=65024 Len=12

Frame 1268: 66 bytes on wire (528 bits), 66 bytes captured (528 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11384, Ack: 591, Len: 12
<flag:c1 + len:0c + code:27ff + data:02 f8 03 0d f8 02 00 00>

2019/272 01:24:37.760315    9412 → 56900 [ACK] Seq=591 Ack=11396 Win=65688 Len=0

Frame 1271: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 591, Ack: 11396, Len: 0

2019/272 01:24:38.069979    56900 → 9412 [PSH, ACK] Seq=11396 Ack=591 Win=65024 Len=6

Frame 1274: 60 bytes on wire (480 bits), 60 bytes captured (480 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11396, Ack: 591, Len: 6
<flag:c1 + len:06 + code:2a08 + data:c4 00>

2019/272 01:24:38.271335    9412 → 56900 [ACK] Seq=591 Ack=11402 Win=65680 Len=0

Frame 1277: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 591, Ack: 11402, Len: 0

2019/272 01:24:38.272079    56900 → 9412 [PSH, ACK] Seq=11402 Ack=591 Win=65024 Len=100

Frame 1279: 154 bytes on wire (1232 bits), 154 bytes captured (1232 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11402, Ack: 591, Len: 100
<flag:c1 + len:06 + code:2a14 + data:dd 00>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 86 00 11 ae 7d>
<flag:c1 + len:10 + code:112e + data:7d 00 02 00 00 00 00 00 02 00 ...>
<flag:c1 + len:10 + code:112e + data:7d 00 00 00 00 00 00 e9 00 00 ...>
<flag:c3 + len:13><flag:c1 + len:0a + code:0919 + data:0f 87 00 11 ae 7d>
<flag:c1 + len:08 + code:d40f + data:88 18 30 70>
<flag:c1 + len:08 + code:d40f + data:94 1a 32 10>
<flag:c1 + len:08 + code:d40f + data:96 15 33 70>

2019/272 01:24:38.471343    9412 → 56900 [ACK] Seq=591 Ack=11502 Win=65580 Len=0

Frame 1283: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 591, Ack: 11502, Len: 0

2019/272 01:24:38.536514    56900 → 9412 [PSH, ACK] Seq=11502 Ack=591 Win=65024 Len=8

Frame 1286: 62 bytes on wire (496 bits), 62 bytes captured (496 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11502, Ack: 591, Len: 8
<flag:c1 + len:08 + code:d40f + data:95 23 2f 70>

2019/272 01:24:38.771357    9412 → 56900 [ACK] Seq=591 Ack=11510 Win=65572 Len=0

Frame 1289: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 591, Ack: 11510, Len: 0

2019/272 01:24:38.804073    9412 → 56900 [PSH, ACK] Seq=591 Ack=11510 Win=65572 Len=7

Frame 1292: 61 bytes on wire (488 bits), 61 bytes captured (488 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 591, Ack: 11510, Len: 7
<flag:c1 + len:07 + code:d725 + data:31 21 20>

2019/272 01:24:38.805550    56900 → 9412 [PSH, ACK] Seq=11510 Ack=598 Win=65024 Len=104

Frame 1295: 158 bytes on wire (1264 bits), 158 bytes captured (1264 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11510, Ack: 598, Len: 104
<flag:c1 + len:68 + code:faa5 + data:29 6b d6 eb 2c a9 03 21 bb ef ...>

2019/272 01:24:38.862944    9412 → 56900 [FIN, ACK] Seq=598 Ack=11614 Win=65468 Len=0

Frame 1298: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 598, Ack: 11614, Len: 0

2019/272 01:24:38.863734    56900 → 9412 [ACK] Seq=11614 Ack=599 Win=65024 Len=0

Frame 1300: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11614, Ack: 599, Len: 0

2019/272 01:24:38.863944    56900 → 9412 [FIN, ACK] Seq=11614 Ack=599 Win=65024 Len=0

Frame 1302: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11614, Ack: 599, Len: 0

2019/272 01:24:38.865435    [TCP Keep-Alive] 56900 → 9412 [ACK] Seq=11614 Ack=599 Win=65024 Len=0

Frame 1303: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32), Dst: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 56900, Dst Port: 9412, Seq: 11614, Ack: 599, Len: 0

2019/272 01:24:38.865510    9412 → 56900 [ACK] Seq=599 Ack=11615 Win=65468 Len=0

Frame 1306: 54 bytes on wire (432 bits), 54 bytes captured (432 bits) on interface 0
Ethernet II, Src: HonHaiPr_f7:d7:2d (7c:e9:d3:f7:d7:2d), Dst: Tp-LinkT_e8:34:32 (d0:76:e7:e8:34:32)
Internet Protocol Version 4, Src: 192.168.0.100, Dst: 192.168.0.100
Transmission Control Protocol, Src Port: 9412, Dst Port: 56900, Seq: 599, Ack: 11615, Len: 0
