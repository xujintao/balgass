package main

import (
	"encoding/binary"
	"unsafe"
)

type t11 struct {
	f00 *t11
	f04 *t11
	f08 *t11
	f10 uint32
	f0C *conn
	f15 uint8
}

var v08C8D014encrypt packetEncrypt // {f00:0x0ECA3300, f18:0x0ECA82E0}
var v08C88FB4decrypt packetEncrypt

type packetEncrypt struct {
	f00 *t12
	f18 *t11
}
type t12 struct {
	encrypt *packetEncrypt
	f04     *t11
}

// f00B98ED0encrypt 加密报文 被花了
func (t *packetEncrypt) f00B98ED0encrypt(conn *conn, bufdst []uint8, bufsrc []uint8, lensrc int) int {
	// // 0x0AC34017
	// // 8字节局部变量
	// var ebp8 t12
	// // 0x0A327C79 0x09FE6629 0x0A561513 0x0AF9A9FC 0x0A059A62 0x0A88BAD0 0x09FD2349 0x0AD981DC 0x0A561A1B
	// // push ebx
	// // push ebp
	// // push esi
	// // push edi
	// // push ecx
	// // pop edi // edi = ecx
	// // push &conn
	// // push &ebp8
	// // push edi
	// // pop ecx // ecx = edi
	// // 0x00B98EE5
	// t.f00B98D90(&ebp8, &conn) // ebp8:{0x08C8D014, 0}
	// // 0x0AFD452E
	// esi := ebp8.encrypt
	// ebp := t.f18
	// eax := t.f00
	// if esi != nil {
	// 	// 0x0A8FA15D
	// 	if esi == eax {
	// 		// 0x00B98F00
	// 	}
	// }
	// // 0x00B98EFB
	// // f00DE84C8()
	// // 0x0A9D324B
	// ebx := ebp8.f04
	// if ebx != ebp {
	// 	// 0x00B98F15 0x00B98ED5 0x0A43C4FF
	// 	if esi == nil {
	// 		// 0x00B98F19
	// 		// f00DE84C8()
	// 	} else {
	// 		// 0x00B98F8A 0x00B98F0C 0x0AD81E47
	// 		esi := esi.f00
	// 		// 0x00B98F1E
	// 	}
	// 	// 0x0A83E939 0x0A748151 0x0A4E9023 0x09E30AF0 0x0A7ABDD1
	// 	if ebx == esi.f18 {
	// 		// 0x00B98F23
	// 		// f00DE84C8()
	// 		// 0x0A5FAB24
	// 	} else {
	// 		// 0x00B98F28 0x0A5FAB24
	// 	}
	// 	// 0x0A5FAB24
	// 	edx := lensrc
	// 	esi := ebx.f10
	// 	eax := bufsrc
	// 	ebp := bufdst
	// 	// push edx
	// 	// push eax
	// 	// push ebp
	// 	// ebx = &esi.f90
	// 	// ecx = ebx
	// 	// 0x00B98F42
	// 	// ebp8 := xx.f00BA1B70(ebp, eax, edx) // (bufdst, bufsrc, lensrc)
	// 	// 0x09EAF588 0x0AF7C87D 0x0A38CE99 0x0A05D16B 0x0AF7C2BF 0x0A4E15DD 0x09FE8547
	// 	if ebp == nil {
	// 		// 0x00B9BFA7
	// 	} else {
	// 		// 0x09FBA820
	// 	}
	// }
	// // 0x0A603AE1
	// // pop edi
	// // pop esi
	return 0
}

// f00B98FC0decrypt 解密报文 被花了
func (t *packetEncrypt) f00B98FC0decrypt(conn *conn, bufdst []uint8, bufsrc []uint8, lensrc int) int {
	// 8字节局部变量
	var ebp8 t12

	// 0x09FE83A3
	if conn == nil {
		// 0x0A8FB4A8
	}
	// 0x00B98FD7 0x0A05EBFC 0x0A8881B6 0x0A5FC1CD 0x0AD2B88E
	// push esi
	// push edi
	// push &conn
	// push &ebp8
	// 0x00B98FE5
	t.f00B98D90(&ebp8, &conn)
	// 0x09FC76D8
	return 0
}

func (t *packetEncrypt) f00B98D90(unk1 *t12, connp **conn) {
	// // 0x10字节局部变量
	// ebp := connp
	// esi := t
	// ecx := esi.f18
	// eax := ecx.f04
	// edi := ecx
	// for eax.f15 == 0 {
	// 	if eax.f0C >= (*connp) {
	// 		edi = eax
	// 		eax = eax.f00
	// 	} else {
	// 		eax = eax.f08
	// 	}
	// }
	// eax = esi
	// ebx = eax.f18
	// v2 := edi
	// v1 := eax
	// f00DE84C8()
	// if edi != ebx {
	// 	if (*connp) >= edi.f0C {
	// 		ecx = v2
	// 	}
	// } else {
	// 	v4 = esi.f18
	// 	v3 = esi.f00
	// 	ecx = v3
	// }
	// edx = ecx.f00
	// ecx = ecx.f04

	// eax = unk1
	// eax.f00 = edx
	// eax.f04 = ecx
	*unk1 = *t.f00
}

var v086A3C28 int

// hook到 igc.dll了
func f00760292parse(conn *conn, unk1, unk2 int) {
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
				ebp182C = v08C8D014encrypt.f00B98FC0decrypt(conn, ebp181Cbuf[2:], ebp18buf[2:], ebp10len-2)

			// 0AF964D8
			default: // 0xC4
				// ebp1834 := ebp18buf
				ebp10len = int(ebp18buf[1]<<8 + ebp18buf[2])
				ebp182C = v08C8D014encrypt.f00B98FC0decrypt(conn, ebp181Cbuf[2:], ebp18buf[3:], ebp10len-3)
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
		if unk1 == 1 {
			f00735DC7(unk2, ebp14code, ebp18buf, ebp10len, ebp1820enc)
		} else {
			f0075C3B2handlecmd(uint8(ebp14code), ebp18buf, ebp10len, func() bool { return ebp1820enc == 1 }())
		}
	}
}

// pb
type pb struct {
	len uint16
	buf [5210]uint8 // 0x145A
	// buf   [5212]uint8 // 0x145C
	f145C t1
}

var v012E4034conn *conn = &v08C88FF0conn

// s9 f00439106
func (t *pb) f00439178init() {
	// SEH
	t.f145C.f007BF63B()
	t.f145C.f007BF68F()
	t.f0043921B()
}
func (t *pb) f004391CF() {}
func (t *pb) f0043921B() {
	t.len = 0
}

// s9 f004391BA
func (t *pb) f0043922CwriteHead(flag uint8, code uint8) {
	t.f0043921B()

	// 写flag
	var buf [1]uint8
	buf[0] = flag
	t.f00439298writeBuf(buf[:], 1, false)

	// 写len
	switch flag {
	case 0xC1:
		t.f00439298writeBuf(t.buf[:], 1, false)
	case 0xC2:
		t.f00439298writeBuf(t.buf[:], 2, false)
	}

	// 写code
	buf[0] = code
	t.f00439298writeBuf(buf[:], 1, false)
}

func (t *pb) f00439298writeBuf(buf []uint8, len int, needXor bool) {
	tlen := int(t.len)
	if tlen+len > 0x145A {
		return
	}
	f00DE7C90memcpy(t.buf[tlen:], buf, len)
	if needXor {
		t.f0043930Cxor(tlen, tlen+len, 1) // (3,4,1), (4,8,1)
		// xor.Enc(t.buf[:], tlen, tlen+len)
	}
	t.len += uint16(len)
}

func (t *pb) f0043930Cxor(begin, end, interval int) {
	var buf = [32]uint8{
		0xAB, 0x11, 0xCD, 0xFE,
		0x18, 0x23, 0xC5, 0xA3,
		0xCA, 0x33, 0xC1, 0xCC,
		0x66, 0x67, 0x21, 0xF3,
		0x32, 0x12, 0x15, 0x35,
		0x29, 0xFF, 0xFE, 0x1D,
		0x44, 0xEF, 0xCD, 0x41,
		0x26, 0x3C, 0x4E, 0x4D,
	}
	ebp24 := begin
	for {
		if ebp24 == end {
			break
		}
		t.buf[ebp24+2] ^= (t.buf[ebp24+1] ^ buf[ebp24]%32) // t.buf[5] ^= t.buf[4] ^ buf[3%32]
		ebp24 += interval
	}
}

// s9 f00439378
func (t *pb) f004393EAsend(needEnc, isC2 bool) {
	// t.f00439612() 写len
	func() []uint8 {
		switch t.buf[0] {
		case 0xC1:
			f00DE7C90memcpy(t.buf[1:], ((*[1]uint8)(unsafe.Pointer(&t.len)))[:], 1)
		case 0xC2:
			f00DE7C90memcpy(t.buf[1:], ((*[2]uint8)(unsafe.Pointer(&t.len)))[:], 2)
		}
		return t.buf[:]
	}()

	t.f0043968Frandom() // hook到 igc.dll了，直接ret

	// f00439420
	func(buf []uint8, len int, needEnc, isC2 bool) {
		// hook到 IGC.dll SendPacket
		f00DE8A70chkstk() // 0x3124

		if !needEnc {
			v012E4034conn.f004397E3write(buf, len)
			return
		}

		var ebp311C int
		var ebp3118 [1000]uint8 // 原始缓存副本
		var ebp1914 int
		var ebp1910 [1000]uint8 // C4编码缓存
		var ebp108 [1000]uint8  // C3编码缓存

		// 编码数据
		f00DE7C90memcpy(ebp3118[:], buf, len)
		ebp3118[len] = uint8(f00DE8AADrand() % 256) //
		if isC2 {
			ebp3118[0] = 0xC2
		}
		if ebp3118[0] != 0xC1 {
			ebp311C = 1
		}
		ebp311C += 2
		ebp3118[ebp311C-1] = v08C88F60 // 流水号？
		v08C88F60++
		ebp311C--

		ebp1914 = v08C88FB4decrypt.f00B98ED0encrypt(v012E4034conn, nil, ebp3118[ebp311C:], len-ebp311C) // 得到len为0x11

		// len<256且C1 编码方案
		if ebp1914 < 0x100 && !isC2 {
			ebp3120 := ebp1914 + 2 // 0x13，超过255怎么办？
			ebp108[0] = 0xC3
			ebp108[1] = uint8(ebp3120)
			v08C88FB4decrypt.f00B98ED0encrypt(v012E4034conn, ebp108[2:], ebp3118[ebp311C:], len-ebp311C) // 编码数据
			v012E4034conn.f004397E3write(ebp108[:], ebp3120)
			return
		}
		// len>=256或者C2 编码方案
		ebp3124 := ebp1914 + 3
		ebp1910[0] = 0xC4
		binary.BigEndian.PutUint16(ebp1910[1:], uint16(ebp3124))
		v08C88FB4decrypt.f00B98ED0encrypt(v012E4034conn, ebp1910[3:], ebp3118[ebp311C:], len-ebp311C)
		v012E4034conn.f004397E3write(ebp1910[:], ebp3124)
	}(t.buf[2:], int(*(*uint16)(unsafe.Pointer(&t.buf[0]))), needEnc, isC2)
}

// hook到 igc.dll了，直接ret
func (t *pb) f0043968Frandom() {
	var ebp4len uint16
	var ebp8buf []uint8
	var ebp9code uint8
	if t.buf[2] == 0xC1 {
		ebp4len = uint16(t.buf[1])
		ebp9code = t.buf[2]
		ebp8buf = t.buf[3:]
		ebp4len -= 3

	}
	if t.buf[2] == 0xC2 {
		ebp4len = uint16(t.buf[1]<<8 + t.buf[2])
		ebp9code = t.buf[3]
		ebp8buf = t.buf[4:]
		ebp4len -= 4
	}

	t.f145C.f007C2763(ebp9code, ebp8buf, ebp4len)
}

func (t *pb) f0043974FwriteZero(x int) {
	var buf [5216]uint8
	f00DE8A70chkstk()

	f00DE8100memset(buf[:], 0, x)
	t.f00439298writeBuf(buf[:], x, true)
}

func (t *pb) f004397B1writeUint8(data uint8) *pb {
	var buf [1]uint8
	buf[0] = data
	t.f00439298writeBuf(buf[:], 1, true)
	return t
}

func (t *pb) f004C65EFwriteUint16(data uint16) *pb {
	var buf [2]uint8
	binary.LittleEndian.PutUint16(buf[:], data)
	t.f00439298writeBuf(buf[:], 2, true)
	return t
}

func (t *pb) f0043EDF5writeUint32(data uint32) {
	var buf [4]uint8
	binary.LittleEndian.PutUint32(buf[:], data)
	t.f00439298writeBuf(buf[:], 4, true)
}
