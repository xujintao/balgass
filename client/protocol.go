package main

// t08C8D014
var v08C8D014encrypt packetEncrypt

type packetEncrypt struct{}

func (t *packetEncrypt) f00B98FC0(conn *conn, bufdst []uint8, bufsrc []uint8, len int) int { return 0 }

func f00735DC7(uk1 int, code int, buf []uint8, len int, enc int) {

}

func f0075C3B2(flag int, buf []uint8, len int, x int) {
	ebp5D8 := flag
	if ebp5D8 > 0xFD {
	}

	// 寻找处理器
	// v0075FCB2cmd[4*v0075FF6Acmd[ebp5D8]]

	// 0x00 处理器
	// _006FA9EB
	func(buf []uint8) {
		// 业务逻辑很复杂
		// send c1 04 f4 06
	}(buf)

	// 0xF1 处理器
	ebp4 := buf
	ebp5DC := ebp4[3]
	if ebp5DC > 4 {
	}
	switch ebp5DC {
	case 4:
	case 0: // version match, buf: 01 2E 87 "10525"
		// _006C75A7 被花到 0x0A4E7A49
		func(buf []uint8) {
			ebp4 := buf
			ebp10 := buf[4]
			if ebp10 == 1 {
				ebp8 := &v0130F728
				ebp8.f004A9123(&ebp8.f4880)
			}
			v08C88E0C = int(ebp4[5]<<8 + ebp4[6])
			v08C88E08 = 2
			var ebpC uint8
			if ebpC >= 5 {
			}
			for {
				if v012E4018[ebpC]-ebpC-1 != ebp4[7+ebpC] { // 改为jmp
					break
				}
				ebpC++
				if ebpC >= 5 {
					break
				}
			}

			// _004A9146
			func() {}()

			// _004A9EEB
			func() {}()
			v01319E08log.f00B38AE4printf("Version dismatch - Join server.\r\n")
			// _0051FB61
			func() {}()
		}(buf)
	}

	// 0xFA 处理器(没有)

	// 0xF4 处理器
	if buf[0] == 0xC1 {
	}
	ebp1C := buf
	ebp14 := ebp1C[4]
	ebp5EC := ebp14
	switch ebp5EC {
	case 3:
	case 5:
	case 6:
		// _006BFA3E
		func(buf []uint8) {
			ebp4 := 6
			// ebp9 := buf[ebp4]
			ebp4++
			if ebp14 >= buf[5]<<8|1 {
			}
			// ebp18 := buf[ebp4:]
		}(buf)
	}
}

// hook到 igc.dll了
func f00760292parse(conn *conn, x, y int) {
	// SEH
	f00DE8A70chkstk() // 0x5F5C

	for {
		// 0AD3CD7D
		ebp18buf := conn.f006BDC33() // 从连接中获取buf
		if ebp18buf == nil {
			break
		}

		var ebp14code int
		var ebp10len int
		var ebp181Cbuf [7024]uint8
		var ebp1820enc int

		switch ebp18buf[0] {
		case 0xC1:
			// ebp1824 := ebp18buf
			ebp14code = int(ebp18buf[2])
			ebp10len = int(ebp18buf[1])
		case 0xC2:
			// ebp1828 := ebp18buf
			ebp14code = int(ebp18buf[3])
			ebp10len = int(ebp18buf[1]<<8 + ebp18buf[2])

		// 0A443ED9 0AD2AFB1
		case 0xC3, 0xC4:
			var ebp182C int

			switch ebp18buf[0] {
			// 0A0595EE
			case 0xC3:
				// ebp1830 := ebp18buf
				ebp10len = int(ebp18buf[1])
				ebp182C = v08C8D014encrypt.f00B98FC0(conn, ebp181Cbuf[2:], ebp18buf[2:], ebp10len-2)

			// 0AF964D8
			default: // 0xC4
				// ebp1834 := ebp18buf
				ebp10len = int(ebp18buf[1]<<8 + ebp18buf[2])
				ebp182C = v08C8D014encrypt.f00B98FC0(conn, ebp181Cbuf[2:], ebp18buf[3:], ebp10len-3)
			}
			ebp1820enc = 1

			// 0A049505
			if ebp182C < 0 {
				// var ebp2CB4 pb
				// ebp2CB4.f00439178init()
				// ebp2D6E := (uint8)0xF1
				// ebp2D6D := (uint8)0xC1
				// (uint16)&ebp2CB4[0] = 0
				// ebp5F60 := ebp2D6D
				// if ebp5F60 == 0xC1
			}

			// 0A83F8C1
			if v08C88F61 == ebp181Cbuf[2] {
				v08C88F61++
			} else {
				ebp1820enc = 0
				v08C88F61 = ebp181Cbuf[2]
				v01319E08log.f00B38AE4printf("Decrypt error : g_byPacketSerialRecv(0x%02X), byDec(0x%02X)\r\n", v08C88F61, ebp181Cbuf[2])
				v01319E08log.f00B38AE4printf("Dump : \r\n")
				v01319E08log.f00B38B43(ebp18buf, ebp10len)
				v01319E08log.f00B38AE4printf("\r\n")
				v01319E08log.f00B38B43(ebp181Cbuf[2:], ebp182C)
				// ...
			}

			switch ebp18buf[0] {
			// 0A43FECF
			case 0xC3:
				// ebp2CB8 := ebp181Cbuf[1:]
				ebp181Cbuf[1] = 0xC1
				ebp181Cbuf[2] = uint8(ebp182C)
				ebp14code = int(ebp181Cbuf[3])
				ebp18buf = ebp181Cbuf[1:]

			// 0A4417A0
			default:
				// ebp2CBC := ebp181Cbuf[:]
				ebp181Cbuf[0] = 0xC2
				ebp181Cbuf[1] = uint8(ebp182C >> 8)
				ebp181Cbuf[2] = uint8(ebp182C & 0xFF)
				ebp14code = int(ebp181Cbuf[3])
				ebp18buf = ebp181Cbuf[:]
			}
			ebp10len = ebp182C
		}

		// 0ABE2DAF
		v086A3C28 += ebp10len // 累计接收字节数

		//
		if x == 1 {
			f00735DC7(y, ebp14code, ebp18buf, ebp10len, ebp1820enc)
		} else {
			f0075C3B2(ebp14code, ebp18buf, ebp10len, ebp1820enc)
		}
	}
}
