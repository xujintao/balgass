package main

import (
	"encoding/binary"
	"log"
	"unsafe"
)

// cmd
func f0075C3B2handlecmd(code uint8, buf []uint8, len int, enc bool) {
	// 0x09FB655C
	if h, ok := cmds[int(code)]; ok {
		h(code, buf, len, enc)
		return
	}
	// 0x0075CAFD 0x0075FCAC
	log.Printf("invalid cmd, code:%02dx, buf: %v", code, buf)
}

// cmd table
// key: 0x0075FF6A
// value: 0x0075FCB2
var cmds = map[int]func(code uint8, buf []uint8, len int, enc bool){
	0x00: f006FA9EBhandle00,
	0x01: handle01,
	0x02: handle02,
	0x0D: f007087BFhandle0D,
	0xD2: f0075F7D0handleD2,    // cash shop
	0xD7: handleD7position,     // hook D4->D7
	0xD9: handleD9normalAttack, // hook D9->11
	0xF1: handleF1,
	0xF3: f0075C794handleF3, // character
	0xF4: handleF4,
	0xF6: f006C18B6handleF6,
}

func f006FA9EBhandle00(code uint8, buf []uint8, len int, enc bool) {
	// SEH
	f00DE8A70chkstk() // 0x4B34
	if v012E2340 == 2 {
		// 请求服务器列表
		v01319E08log.f00B38AE4printf("Send Request Server List.\r\n")
		// 0x006FAA2B 压缩
		var reqServerList pb // [c1 04 f4 06]
		reqServerList.f00439178init()
		reqServerList.buf[0] = 0xC1
		reqServerList.buf[2] = 0xF4
		reqServerList.buf[3] = 0x06
		reqServerList.len = 4
		reqServerList.buf[1] = uint8(reqServerList.len)
		// reqServerList.f004393EAsend(false, false) // 直接api发送了
		// ...
	}
	// ...
}
func handle01(code uint8, buf []uint8, len int, enc bool) {}
func handle02(code uint8, buf []uint8, len int, enc bool) {}
func f007087BFhandle0D(code uint8, buf []uint8, len int, enc bool) {
	subcode := buf[3]
	switch subcode {
	case 0:
	case 1:
		if v012E2340 == 4 {
			ebp14 := f004A7D34getServiceManager()
			ebp14.f004A9F3B(buf[13:]) // "you are logged in as Free Player"
		}
	}
}

// s9 f00673854
func f0075F7D0handleD2(code uint8, buf []uint8, len int, enc bool) {
	// 0x0A5FBCA6 0x0A5D4C5A 0x0A4381F1
	subcode := buf[3]
	subcode--

	// 0x00760216 table, s9 0x00674292
	switch subcode {
	case 0: // s9 0x00673894
		// 0x0075F810 0x0075F797 0x0A0C1E4B 0x0075F813
		// f006C28D9, s9 f0066B85C
		func(buf []uint8) {
			// 0x0A6014ED 0x09F8B6AA 0x0A9044AC 0x09FBB3E7 0x0ABE308B
			// totalCash := *(*float64)(unsafe.Pointer(&buf[5]))
			// 0x006C28EE f009F0787
			// 0x006C287C 0x0ABF9A98 0x006C28F5
			// v09D8D548.f009F1748(totalCash)

			// 0x0A55EA23 0x0A1938CA
			// cashC := *(*float64)(unsafe.Pointer(&buf[13]))
			// 0x006C2905 f009F0787
			// 0x006C2898 0x0A9F59D7 0x006C290C
			// v09D8D548.f009F175F(cashC)

			// 0x0A5FED76
			// cashP := *(*float64)(unsafe.Pointer(&buf[21]))
			// 0x006C291C f009F0787
			// 0x006C28A3 0x0AF7CA9A 0x006C2923
			// v09D8D548.f009F1776(cashP)

			// 0x006C28BA 0x0AF76971
			return
		}(buf)
	case 1:
		// f0075F81E
	}
}
func handleD7position(code uint8, buf []uint8, len int, enc bool)     {}
func handleD9normalAttack(code uint8, buf []uint8, len int, enc bool) {}
func handleF1(code uint8, buf []uint8, len int, enc bool) {
	subcode := buf[3]
	switch subcode {
	case 0: // version match
		// f006C75A7 被花到 0x0A4E7A49 // version match, buf: 01 2E 87 "10525"
		func(buf []uint8) {
			ebp4 := buf
			ebp10 := buf[4]
			if ebp10 == 1 {
				// ebp8 := &v0130F728
				// ebp8.f004A9123(&ebp8.f4880)
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
	case 4:
	}
}

func f0075C794handleF3(code uint8, buf []uint8, len int, enc bool) {
	if buf[0] != 0xC1 {
		// 0x0075C7B3

	}
	// 0x09FB56EA 0x0075C7C3 0x0AD3CA7D
	subcode := buf[3]
	switch subcode {
	case 0: // character list
		// 0x0075C97D
	case 1:
		// 0x0075C98B
	case 2:
		// 0x0075C999
	case 3: // character info
		// 0x0075C9A7
		f00701DF4(buf, enc)
	case 4:
		// 0x0075C9C4
	case 5:
		// 0x0075C9EF
	case 0x13:
		// 0x0075CA35
	case 0x14:
		// 0x0075CA43
	case 0x15:
		// 0x0075CAE9
	case 0x20:
		// 0x0075CA51
	case 0x22:
		// 0x0075CA6D
	case 0x23:
		// 0x0075CA91
	case 0x24:
		// 0x0075CA7B
	case 0x25:
		// 0x0075CA86
	case 0x26:
		// 0x0075CA5F
	case 0x30:
		// 0x0075CA9C
	case 0x32:
		// 0x0075CAA7
	case 0x40:
		// 0x0075CAB2
	case 0x50:
		// 0x0075CABD
	case 0x51:
		// 0x0075CAC8
	case 0x52:
		// 0x0075CAD3
	case 0x53:
		// 0x0075CADE
	}
}

func handleF4(code uint8, buf []uint8, len int, enc bool) {
	subcode := buf[3]
	switch subcode {
	case 3: // server info
	// close连接服务器
	// f006BF89ADial(ip, port) // 拨号游戏服务器
	// set state
	case 5:
	case 6: // server list
	}
}

func handleFA(code uint8, buf []uint8, len int, enc bool) {}
func f006C18B6handleF6(code uint8, buf []uint8, len int, enc bool) {
	// c1 05 f6 1a 00
	// 0x0A8FB4BE hook to hide f0A9F6026, jump complicated shell logic, which would send trap message
	var label1 uint32 = 0x00FF20B1
	// push label1
	// push 0x09E29528
	// ret

	// 0x09E29528
	var label2 uint32 = 0x0042D75C
	// push 0x0042D75C
	// push 0x09E27458
	// ret

	// 0x09E27458
	var label3 uint32 = 0x0A88AE46
	// push 0x0A88AE46
	// push 0x012DDF16
	// ret

	// 0x012DDF16 0x0ABF9D31 0x0ABE1253
	// push edi esi ebx ecx fd edx

	// 0x09FDE30A
	// push 0x0A32AC32
	// push 0x0A4425A9
	// ret

	// f0A4425A9
	func() {
		label1 = 0x0A9F6026
		label2 = 0x0AD70BD7
		// f0A4EC883 send trap message
		func(unk uint32) {
			var ebp148Ctrap pb
			ebp148Ctrap.f00439178init()
			ebp148Ctrap.f0043922CwritePrefix(0xC1, 0xF3)
			ebp148Ctrap.f004397B1writeUint8(0x31)
			// big endian
			ebp148Ctrap.f004397B1writeUint8(uint8(unk >> 24))
			ebp148Ctrap.f004397B1writeUint8(uint8(unk >> 16))
			ebp148Ctrap.f004397B1writeUint8(uint8(unk >> 8))
			ebp148Ctrap.f004397B1writeUint8(uint8(unk))
			ebp148Ctrap.f004393EAsend(true, false)
			ebp148Ctrap.f004391CF()

			// for {
			// 	ebp1490 := 0
			// 	if ebp1490 >= 0x190 {
			// 		break // 0x007C8001
			// 	}
			// 	if ebp1490 == v08C88CAC {
			// 		break
			// 	}
			// 	ebp1494 := f004373C5(ebp1490).f00A38D5BgetObject()
			// 	if ebp1494.m5E {
			// 		break
			// 	}
			// 	ebp1490++
			// }
			// ebp1494.m404 = 0x2710
		}(0x00007533)
		// 0x0A9F747B
		label3 = 0x0A4DD48C
	}()

	// 0x0A32AC32
	// pop edx fd ecx ebx esi edi
	// label3
	// label2
	// label1

	// f0A9F6026 隐藏函数
}

func f00705035resourceGuardError() {
	//
	// 0x00705035: f006BD5BBclose()
	// 0x0070503A: jmp 0x09EB027B, ResourceGuard Error, jmp 0x00705044
	// 0x00705044: f00B38AE4() // v01319E08log.f00B38AE4printf()
	// 0x00705067: f00B60690()
	// 0x00706761: return // 0x0075C9B2
}

var v09FD6693 = [...]uint{
	0x7C847811,
	0x9A04CD68,
	0x21295446,
	0xECDDBC71,
	0xBA156216,
}
var v0A05044A = [...]uint{
	0x1BC15B10,
	0x75C966E2,
	0xB99388B8,
	0xFCEFE807,
	0x79C783E6,
}

// cipher key
var v0A846872 = [4]uint{
	0x6D6DA521,
	0x50FF5BCA,
	0x51F68227,
	0x45C5A1EC,
}

// block 0: (v0A8FCD55, 4)
var v0A8FCD55 = [4]uint8{0x9E, 0x01, 0x1D, 0xC6}

// block 1~n, cipher text
var v0A7445E4 = [485]block{
	{0x4D1250B3, 0x91A3D473}, // block 1: v00B98ED0, 10
	{0x02EAB896, 0x966E5E5B}, // block 2: v00B98EE5, 15
	// ...
	{0xEDF20833, 0x61320F98}, // 0x0A7455FC, block 484, ?
	{0x09609CAB, 0xDCB5C892}, // 0x0A745504, block 485: 0xFFFFFFFF, 0xFFFFFFFF
}

// xor key
var v0A935B1D uint = 0xAE5D1C94
var v0AD3C20F uint = 0x306A652D

// xor key
var v0B10C78C = [4]uint{}
var v09EAC72E = [4]uint{}

var v0AF7C3C4label1 = 0x0ABF8601
var v0AFD83CAlabel1 = 8
var v0AD7BED9label2 = 0x0A55FB91
var v0ABD891Flabel2 = 7
var v09FE4BFFlabel3 = 0x09E733E8
var v0AF836B5label3 = 6

func f00701DF4(buf []uint8, enc bool) { // (v0018C80C, true)
	// 0x0AC3581C
	var label1 = 0x0A9368DD
	// push label1
	// push 0x0AD33DDE
	// ret

	// 0x0AD33DDE
	var label2 = 0x09FDA253
	// push label2
	// push 0x0A5D4A14
	// ret

	// 0x0A5D4A14
	var label3 = 0x00BF60A2
	// push label3
	// push 0x0AD3B967
	// ret

	// 0x0AD3B967
	// push edi esi ebx edx ecx fd
	// 0x0AC34158
	// push 0x0AF80858
	// push 0x0A604CD4
	// ret

	// 0x0A604CD4 f0A604CD4, 完整性校验
	func() {
		// 0x84局部变量
		v0018C068 := struct {
			m00  [5]uint // 0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476, 0xC3D2E1F0
			m14  uint
			m18  uint
			size uint16
			m1E  [64]uint8
			m60  uint
			m64  uint
		}{}
		var ebp8 block

		// push esi edi eax
		// push 0x0AF76043
		// push 0x0AD91F40
		// ret

		// 0x0AD91F40
		eax := &v0018C068 // eax为什么为0x0018C068?
		// push esi
		if eax == nil {
			// 0x0A392626
			// eax = 1
			// pop esi
			// 0x0AF76043
		} else {
			// 0x09E8C84C 0x09EB665D 0x0AD3DBA8
			edx := 0
			v0018C068.m14 = 0
			v0018C068.m18 = 0
			v0018C068.size = 0
			ecx := 0
			// push edi
			// 0x09E8C4D7
			for edx < 5 {
				// v0018C068.data[edx] = v09FD6693[ecx] ^ v0A05044A[ecx]
				// 0x0A4DD423 0x0A9010BC 0x09F8A6F5 0x09E9201B 0x0AD8957E 0x09FC3702 0x0AD56A99
				edx++
				ecx++
			}
			// 0x0A55D14E 0x09FDE315 0x09FD31D0 0x0AD71C58
			v0018C068.m60 = 0
			v0018C068.m64 = 0
			// pop edi
			// eax = 0
			// pop esi
			// 0x0AF76043
		}

		f0A4424B3calc := func() int {
			// 计算v0018C068.m00[5]
			// 0x15C局部变量
			var ebp1C = [4]uint{}   // 0x5A827999, 0x6ED9EBA1, 0x8F1BBCDC, 0xCA62C1D6
			var ebp15C = [80]uint{} // 0x9E011DC6,
			// push ebx esi edi
			ecx := 0
			eax := 0
			edi := 4
			// 0x09EB0A9C
			for {
				ecx++
				ebp1C[ecx] = v0B10C78C[eax] ^ v09EAC72E[eax] // *(ebp + ecx*4 - 20) = v0B10C78C[eax] ^ v09EAC72E[eax]
				eax++
				if ecx >= edi {
					break // 0x0A43FBE5
				}
			}
			// 0x0A43FBE5
			esi := 0
			data := v0018C068.m1E[:]
			// 0x09F886D1
			for {
				ebp15C[esi] = uint(binary.BigEndian.Uint32(data)) // *(ebp + esi*4 - 0x15C) = uint(binary.BigEndian.Uint32(data))
				esi++
				data = data[edi:]
				if esi > 16 {
					break // 0x09FE811D
				}
			}
			// 0x09FE811D
			esi = 64
			edx := 0
			// 0x0A84AD4E
			for {
				ecx = int(ebp15C[edx+13] ^ ebp15C[edx+8] ^ ebp15C[edx+0] ^ ebp15C[edx+2])
				ebp15C[edx+17] = uint(ecx>>31 | ecx*2)
				edx++
				esi--
				if esi <= 0 {
					break // 0x0A557E99
				}
			}
			// 0x0A557E99
			i := 0
			// 0x0AB4CE8F
			for {
				ecx = int(v0018C068.m00[0])
				edx = int(v0018C068.m00[1])
				esi = int(v0018C068.m00[2])
				edi = int(v0018C068.m00[3])
				ebp1C[1] = uint(edi)
				ebp1C[2] = v0018C068.m00[4]
				ebp1C[3] = uint(ecx)
				ecx = ^edx & edi
				edi = esi & edi
				ecx |= edi
				// ...
				i++
				if i >= 20 {
					break // 0x09EBA091
				}
			}
			// 0x09EBA091
			// for {
			// 	// ...
			// }
			// 0x0A000237
			// 0x0A33B671
			// for{
			// 	// ...
			// }
			// 0x09E037DC
			// 0x0AFDB396
			// for{
			// 	// ...
			// }
			// 0x09E314D1
			v0018C068.size = 0
			// pop edi esi ebx
			return 80
		}
		f0A04D0E9copy := func(buf []uint8, size int) {
			// 0x0AFD91A1 ... 0x0A38C729 0x0AF98DB7
			if size == 0 {
				// 0x0A84B8A2
				// eax = 0
				// return
			}
			// 0x0A0501C2 0x09E70719 0x0AA064A9 0x09FDB3C6 0x09FFEA27
			// push esi edi
			if buf == nil {
				// 0x0AF127E7
				// eax = 1
				// return 0x0B285753
			} else {
				// 0x0AA2DD14 0x0AD7F88C 0x0A888838 0x0A32F731
				if v0018C068.m60 != 0 {
					// 0x0A4447E1
					v0018C068.m64 = 3
					// return 0x0B285753
				}
				// 0x0AF999CA 0x0AAB2477 0x0AF12655 0x0A4E99B7 0x0A05B63E
				if v0018C068.m64 != 0 {
					// return 0x0B285753
				}
				// 0x0A9D5444
				for {
					size--
					// 0x0AC36ABB 0x0ABE21DB 0x0AA30D88 0x0A9F8A08 0x0A041A45 0x0A8F9951 0x0A0584FC
					if v0018C068.m64 != 0 {
						// break 0x0A333884
					}
					// 0x0A742D43 0x0AF80694 0x0ABD8DA9
					v0018C068.m1E[v0018C068.size] = buf[0]
					v0018C068.size++
					v0018C068.m14 += 8
					// 0x0A39139D
					if v0018C068.size == 64 {
						// 0x0ABDA41F
						// push esi
						// push 0x0A8912FD
						// push 0x0A4424B3
						// ret
						// 0x0A4424B3 f0A4424B3()
						f0A4424B3calc()
						// 0x0A8912FD 0x09E2FA6E
					}
					// 0x09E2FA6E
					buf = buf[1:]
					if size == 0 {
						break // 0x0A333884
					}
				} // for loop end 0x0A9D5444
				// 0x0A333884
				// eax = 0
				// return 0x0B285753
			}
			// 0x0B285753
			// pop edi
			// pop esi
			// return
		}

		f09FC0148dec := func(key []uint) {
			// 0x18局部变量
			ebp4 := ebp8.addr
			edx := v0A935B1D ^ v0AD3C20F
			eax := edx << 5
			// push esi
			esi := ebp8.size
			// push edi
			// ebp10 := key[3]
			// ebpC := key[2]
			// ebp18 := key[1]
			// ebp14 := key[0]
			cnt := 32

			// 0x0A436A1C
			for {
				// esi -= (ebp4<<4 + ebpC) ^ (ebp4>>5 + ebp10) ^ (ebp4 + eax)
				// ecx := (esi<<4 + ebp14) ^ (esi>>5 + ebp18) ^ (esi + eax)
				// ebp4 -= ecx
				eax -= edx
				cnt--
				if cnt == 0 {
					break // 0x09E2C971
				}
			}
			// 0x09E2C971
			ebp8.addr = ebp4
			ebp8.size = esi

			// pop edi
			// pop esi
			return
		}

		// 计算大块校验和
		// 0x0AF76043, copy section 0 -------------------------------------------
		// push 4
		// push v0A8FCD55[:]
		// push &v0018C068
		// push 0x0A952311
		// push f0A04D0E9
		// ret
		// f0A04D0E9(&v0018C068, v0A8FCD55[:], 4)
		f0A04D0E9copy(v0A8FCD55[:], 4)

		// 0x0A952311, decode and copy section 1 -------------------------------
		// var ebp8 = [2]uint{v0A7445E4[0], v0A7445E4[1]} // 密文，需要异或解密
		ebp8.addr = v0A7445E4[0].addr
		ebp8.size = v0A7445E4[0].size
		// push v0A846872[:]
		// push ebp8[:]
		esi := v0A7445E4[:]
		// push 0x0A32877C
		// push 0x09FC0148
		// ret
		// 0x09FC0148 f09FC0148(&ebp8, v0A846872[:]) // 解密
		f09FC0148dec(v0A846872[:])
		// 0x0A32877C
		label1 = v0AF7C3C4label1 // *(ebp + v0AFD83CAlabel1*4 + 8) = v0AF7C3C4label1
		// 0x0A915FFE
		for {
			if ebp8.addr == uintptr(^uint(0)) {
				break // 0x0A887B7F
			} else {
				// 0x0A338AC2
				if ebp8.addr == uintptr(^uint(1)) {
					// 0x0A0C5AFC
				}
			}
			// 0x09EB2544
			// push ebp8.size                push 10
			// push 0x00400000+ebp8.data     push 0x00B98ED0
			// push &v0018C068
			// push 0x09EB153E
			// push f0A04D0E9
			// ret
			f0A04D0E9copy((*[100]uint8)((unsafe.Pointer(ebp8.addr)))[:], int(ebp8.size))

			// 0x09EB153E, decode and copy section 2 --------------------------------
			esi = esi[2:]
			ebp8.addr = uintptr(esi[0].addr)
			ebp8.size = esi[0].size
			// 0x0A7465E3 0x0AAB49E0 0x09E90C91
			// push v0A846872[:]
			// push &ebp8
			// push 0x0A84AAEA
			// push 0x09FC0148
			// ret
			// 0x09FC0148 f09FC0148(&ebp8, v0A846872[:])
			f09FC0148dec(v0A846872[:])
			// 0x0A84AAEA
			// 0x0A915FFE
		}

		// 0x0A887B7F 开始校验
		// *(ebp+v0ABD891Flabel2*4+8) = v0AD7BED9label2
		label2 = v0AD7BED9label2
		var ebp1C [20]uint8
		// push &ebp1C
		// push &v0018C068
		// push 0x0A32D2F8
		// push 0x0AD3DA6E
		// ret
		// 0x0AD3DA6E f0AD3DA6E(&v0018C068, ebp1C[:]) 计算最后一帧校验和
		func() {
			// push esi
			// esi := &v0018C068
			// if esi == nil || ebp1C == nil {
			// 	// 0x0A9F8D92
			// 	// eax = 1
			// 	// pop esi
			// 	// 0x0A32D2F8
			// }
			// 0x09FDDBAB
			if v0018C068.m64 != 0 {
				// 0x0A4390F3
				// pop esi
				return // 0x0A32D2F8
			}
			// 0x09FFE6EA
			if v0018C068.m60 == 0 {
				// 0x0AFDA74C
				// push edi esi
				// push 0x0A56CE1C
				// push 0x0A842E15
				// ret
				// 0x0A842E15
				// push esi
				// esi = &v0018C068
				v0018C068.m1E[v0018C068.size] = 0x80 // 0x35
				if v0018C068.size <= 0x37 {
					// 0x09E8F0A1
					v0018C068.size++ // 0x36
					if v0018C068.size < 0x38 {
						// 0x0AFD8921
						for {
							v0018C068.m1E[v0018C068.size] = 0
							v0018C068.size++ // 0x37, 0x38
							if v0018C068.size >= 0x38 {
								break // 0x09EB8CE2
							}
						}
					}
					// 0x09EB8CE2 0x0AA0711A 0x0A0C7CA1 0x0AFDE4B8 0x09FD360D 0x0A84CCAB 0x0A90432C
					// 0x0A56D736 0x09FDB4F4 0x0A841B27 0x0AD831E8 0x0AF77ADB
					// little endian
					v0018C068.m1E[56] = uint8(v0018C068.m18 >> 24)
					v0018C068.m1E[57] = uint8(v0018C068.m18 >> 16)
					v0018C068.m1E[58] = uint8(v0018C068.m18 >> 8)
					v0018C068.m1E[59] = uint8(v0018C068.m18)
					// big endian
					v0018C068.m1E[60] = uint8(v0018C068.m14 >> 24)
					v0018C068.m1E[61] = uint8(v0018C068.m14 >> 16)
					v0018C068.m1E[62] = uint8(v0018C068.m14 >> 8)
					v0018C068.m1E[63] = uint8(v0018C068.m14)
					// push 0x0A5FDFCF
					// push 0x0A4424B3
					// ret
					f0A4424B3calc()
					// 0x0A5FDFCF
					// pop ecx esi
				} else {
					// 0x0AC3860F
					v0018C068.size++
					// 0x0A88E950 0x09FD9899 0x0AF9A5FA 0x0B0739E4
					if v0018C068.size < 0x40 {
						// 0x0AF897D5
						for {
							v0018C068.m1E[v0018C068.size] = 0
							v0018C068.size++
							// 0x0AA0D4A4 0x0ABD9808 0x0B10F6EB 0x09E6A6D1 0x0ABB65BD 0x0A4E2DEE 0x0A847593
							if v0018C068.size > 0x40 {
								break // 0x09F8AC4E
							}
						}
					}
					// 0x09F8AC4E 0x09DEB5D0
					// push esi
					// push 0x0AA0514D
					// push 0x0A4424B3
					// ret
					f0A4424B3calc()
					// 0x0AA0514D
					// pop ecx
					// 0x0A0C7A91
					if v0018C068.size < 0x38 {
						// 0x09F89A0A
					} else {
						// 0x09FD73A0
					}
				}
				// 0x0A56CE1C
				// pop ecx
				// eax := 0
				// ecx := 16
				// memset(v0018C068.m14, 0, 16*4)
				v0018C068.m14 = 0
				v0018C068.m18 = 0
				v0018C068.m60 = 1
				// pop edi
			}
			// 0x0A84F492
			eax := 0
			// 0x09FDBFF8, change byte endians
			for {
				edx := eax
				edx &= 3
				ecx := 3
				ecx -= edx
				edx = eax
				edx >>= 2
				edx = int(v0018C068.m00[edx])
				ecx <<= 3
				edx >>= uint(ecx)
				ebp1C[eax] = uint8(edx)
				eax++
				if eax >= 20 {
					break // 0x09FDCEBE
				}
			}
			// 0x09FDCEBE
			eax = 0
			// pop esi
			// 0x0B2846F4
			return // 0x0A32D2F8
		}()

		// 校验
		// 0x0A32D2F8 0x0AD7F532 0x0A56C5AD 0x0AD56B3D
		eax2 := 5 // eax2 := v0A5FA71Clabel // 5
		// pop ecx
		// pop ecx
		eax2--
		// pop edi
		ecx := eax2 // 4
		// pop esi
		// 0x0A84B91E 0x0A5FA4B4 0x09E259D4 0x0A903F4E 0x0A0017AA 0x0A5D6D14 0x09F87D9F 0x0AD3CD2E
		if eax2 > 0 {
			// 0x0AA2FE2E
			for {
				ecx-- // 3, 2, 1, 0, -1
				// v0AFE0A39[ecx] = ebp1C[ecx+1] // v0AFE0A39[3]=ebp1C[4], v0AFE0A39[2]=ebp1C[3], v0AFE0A39[1]=ebp1C[2], v0AFE0A39[0]=ebp1C[1], v0AFE0A39[-1]=ebp1C[0]
				// copy(v0AFE0A35[:], ebp1C[:])
				if ecx < 0 {
					break // 0x0B074B24
				}
			}
		}
		// 0x0B074B24
		if eax2 >= 0 { // eax2 = 4
			// 0x0A43DFB1
			for {
				// ecx := *(ebp + eax2*4 - 0x1C)
				ecx := ebp1C[eax2]
				// ecx -= v0A391F0E[eax2] // v0A391F0E
				if ecx != 0 {
					// 0x0A9FE958, 不一致
				}
				// 0x0AF13005
				eax2--
				if eax2 < 0 {
					// 0x09FCB8D5, 校验通过
					break // 0x0A9F773B
				}
				// 0x0A43DFB1
			}
		}
		// 0x0A9F773B
		// v09FE5B85 = v0A047F75
		label3 = v09FE4BFFlabel3 // *(ebp+v0AF836B5*4+8) = v09FE4BFF
		// return 0x0AF80858
	}()

	// 0x0AF80858
	// pop fd ecx edx ebx esi edi
	// 0x09FB8F6A 0x0AAB2DEA
	// label3 ret
	// label2 0x012DDD4C ret
	// label1 0x0ABF8601
	f00DE8A70chkstk() // 0xDB80局部变量
	// 0x0A049CFF
	if enc == false {
		//
	}
}

// hijack cmd
func handlecmdhook(code uint8, subcode uint8, buf []uint8, len int) {
	if h, ok := cmdhooks[int(code)]; ok {
		h(buf)
		return
	}
	codes := []byte{code, subcode}
	code16 := binary.BigEndian.Uint16(codes)
	if h, ok := cmdhooks[int(code16)]; ok {
		h(buf)
		return
	}
	log.Printf("invalid cmd, code:%02dx, buf: %v", code, buf)
}

// hijack cmd table
var cmdhooks = map[int]func(buf []uint8){
	0xF100: joinResultSend,
	0xF303: charMapJoinResult,
	0xF304: charRegen,
	0xF305: levelUpSend,
	0xF306: levelUpPointAdd,
	0xF351: masterLevelUpSend,
	0xBF51: ansMuBotRecvUse,
	0xF403: reconnect,
	0xCA:   friendRoomCreate,
	0xD4:   AHCheckGetTickHook,
	0xFA0D: offTradeReq,
	0xFA10: customPost,
	0xFA11: customPost,
	0xFA12: setChatColors,
	0xFA90: IGCStatAdd,
	0xFAA0: enableSiegeSkills,
	0xFAA2: setAgilityFix,
	0xFAA3: setCashItemMoveEnable,
	0xFAA4: luckyItemFix, // 区分C1和C2
	0xFAA5: disableReconnect,
	0xFAA6: handleExitProcess,
	0xFAA7: disableReconnectSystem,
	0xFAA8: alterPShopVault,
	0xFAB0: dropSellMod,
	0xFAF8: setBattleZoneData,
	0x11:   attackResult,
	0x17:   diePlayerSend,
	0x26:   refillSend,
	0x27:   manaSend,
}

// hijack cmd handle
func joinResultSend(buf []uint8)         {}
func charMapJoinResult(buf []uint8)      {}
func charRegen(buf []uint8)              {}
func levelUpSend(buf []uint8)            {}
func levelUpPointAdd(buf []uint8)        {}
func masterLevelUpSend(buf []uint8)      {}
func ansMuBotRecvUse(buf []uint8)        {}
func reconnect(buf []uint8)              {}
func friendRoomCreate(buf []uint8)       {}
func AHCheckGetTickHook(buf []uint8)     {}
func offTradeReq(buf []uint8)            {}
func customPost(buf []uint8)             {}
func setChatColors(buf []uint8)          {}
func IGCStatAdd(buf []uint8)             {}
func enableSiegeSkills(buf []uint8)      {}
func setAgilityFix(buf []uint8)          {}
func setCashItemMoveEnable(buf []uint8)  {}
func luckyItemFix(buf []uint8)           {}
func disableReconnect(buf []uint8)       {}
func handleExitProcess(buf []uint8)      {}
func disableReconnectSystem(buf []uint8) {}
func alterPShopVault(buf []uint8)        {}
func dropSellMod(buf []uint8)            {}
func setBattleZoneData(buf []uint8)      {}
func attackResult(buf []uint8)           {}
func diePlayerSend(buf []uint8)          {}
func refillSend(buf []uint8)             {}
func manaSend(buf []uint8)               {}

func f00735DC7(uk1 int, code int, buf []uint8, len int, enc int) {

}

// s9 f0067084A
func f0075C3B2(code int, buf []uint8, len int, x int) {
	f0075C3B2handlecmd(uint8(code), buf, len, func() bool { return x == 1 }())
}
