package main

import (
	"encoding/binary"
	"sync"
	"syscall"
	"unsafe"

	"github.com/xujintao/balgass/client/win"
)

// t3004
type t3004 struct {
	f14 uint32 // 0x0DA6D830
	f18 uint32 // 1
}

// 比较复杂，可能是解析
func (t *t3004) f006BDE09() {
	// _006BE122
	func() {

	}()

	// _006BE0DC
	func() {
		// _006BE360
		func() {

		}()

		// _006BE852
		func() {
			// _006BEEDB
			func() {

			}()
		}()
	}()

	for {
		// _006BE103
		func() {
			// _006BE852
			func() {
				// _006BEEDB
				func() {

				}()
			}()
		}()

		// _006BE2ED
		func() {
			// _006BE8D5
			func() {
				// _0042DE50
				func() {

				}()

				// _0042DBB0
				func() {

				}()
			}()
		}()

		// _006BE2A8
		func() {
			// _006BE882
			func() {
				// _0042DE50
				func() {

				}()

				// _0042DE40
				func() {

				}()

				// _006BE920
				func() {

				}()
			}()
		}()

		// _006BDE91
		func() {
			// _006BDEB8
			func() {

			}()

			// _00DE7538
			func() {

			}()
		}()

		// _006BE2B9
		func() {
			// _006BE86E
			func() {
				// _006BEF19
				func() {
					// _0042DE50
					func() {

					}()

					// _0042DE40
					func() {

					}()

					// _006BE360
					func() {

					}()
				}()
			}()
		}()
	}

	// _006BE16E
	func() {
		// _006BE360
		func() {

		}()

		// _006BE360
		func() {

		}()

		// _006BE368
		func() {

		}()

		for {
			// _006BE360
			func() {

			}()

			// _006BE5A5
			func() {
				// _006BF2C9
				func() {

				}()
			}()

			// _006BE591
			func() {
				// _00DE7538
				func() {

				}()
			}()
		}

	}()
}

// t3003
type t3003 struct{}

func (t *t3003) f006BEB1E(x uint32) []**t3002 {
	return nil
}

func (t *t3003) f006BEB35(p **t3002, pp **t3002) {
	t.f006BF322(p, pp)
}

func (t *t3003) f006BF322(p **t3002, pp **t3002) {
	// ebp4 := p

	// _00405BF0
	ebp8 := func(x uint32, p **t3002) **t3002 {
		return p
	}(4, p)
	if ebp8 == nil {
		return
	}
	*p = *pp
}

// t3002
type t3002 struct {
	buf [2000]uint8
	len int
}

func (t *t3002) f006BDFED(buf []uint8, len int) *t3002 {
	f00DE7C90memcpy(t.buf[:], buf, len)
	t.buf[len] = 0xFD
	t.len = len
	return t
}

// t3001
type t3001 struct {
	data [12]uint8
	f0C  *t3003      // 0x0D906A78
	f10  [][]**t3002 // 指针数组, 0x2817EFA0
	f14  uint32
	f18  uint32
	f1C  uint32
	f20  *t3004 // 0x0DA6FFC0
}

func (t *t3001) f006BDDF5() {
	t.f20.f006BDE09()
}

func (t *t3001) f006BDF76(buf []uint8, len int) {
	// 带SEH了
	var ebp18 *t3002 = (*t3002)(f00DE852F(2004)) // 可能是new
	var ebp20 *t3002
	if ebp18 == nil {
		ebp20 = nil
	} else {
		ebp20 = ebp18.f006BDFED(buf, len) // 可能是初始化
	}
	ebp14 := ebp20
	ebp10 := ebp14
	t.f006BE241(&ebp10) // 填表
}

func (t *t3001) f006BE241(pp **t3002) {
	t.f006BE6D9(pp)
}

func (t *t3001) f006BE6D9(pp **t3002) {
	if (t.f18+t.f1C)%4 == 0 && t.f14 <= (t.f1C+4)>>2 {
		t.f006BEC73(1)
	}
	ebp4 := t.f18 + t.f1C
	ebp8 := ebp4 >> 2
	if t.f14 <= ebp8 {
		ebp8 -= t.f14
	}
	if t.f10[4*ebp8] == nil {
		t.f10[4*ebp8] = t.f0C.f006BEB1E(1)
	}
	t.f0C.f006BEB35(t.f10[4*ebp8][ebp4%4*4], pp)
	t.f1C++
}

func (t *t3001) f006BEC73(x uint32) {}

// conn或者叫netFD
type conn struct {
	hWnd  win.HWND       // f0000
	f0004 int            // f0004
	f0008 int            // f0008
	fd    syscall.Handle // f000C
	bufw  [2000]uint8    // f0010
	w     int            // f2010
	bufr  [2000]uint8    // f2014
	r     int            // f4014
	f4020 *t3001         // f4020
	once  sync.Once      // f40EC, _08C8D0DC
}

var v08C88FF0 conn

func (t *conn) f006BD3A7init() {
	// t.once.Do(win.WSAStartup(0, nil))
	win.WSAStartup(0, nil)
}

func (t *conn) f006BD509socket(hWnd win.HWND, x int) int {
	var err error
	t.fd, err = win.Socket(2, 1, 0) // AF_INT, SOCK_STREAM, 0
	t.f0004 = x
	if err != nil { // INVALID_SOCKET
		win.WSAGetLastError()
		// log
		// MessageBoxA
		return 0
	}
	t.hWnd = hWnd
	return 1
}

func (t *conn) f006BD708dial(ip string, port int, x int) {

}

func (t *conn) f006BDA03read() uint32 {
	if t.r >= 0x2000 {
		v01319E08log.f00B38AE4printf("Receive Packet Buffer Overflow")
		return 1
	}
	ebpC, err := win.Recv(t.fd, t.bufr[t.r:], 0x2000-t.r, 0)
	if ebpC == 0 {
		// 服务器关闭连接
		return 1
	}
	if err != nil { // SOCKET_ERROR
		win.WSAGetLastError()
		//...
		return 1
	}
	t.r += ebpC
	if t.r < 3 {
		return 3 // 还没有收到完整协议报文
	}

	ebp4_len := 0
	ebp8 := 0
	var ebp10 []uint8
	var ebp14 []uint8
	for {
		switch t.bufr[ebp8] {
		case 0xC1, 0xC3:
			ebp10 = t.bufr[ebp8:]
			ebp4_len = int(ebp10[1])
		case 0xC2, 0xC4:
			ebp14 = t.bufr[ebp8:]
			ebp4_len = int(ebp14[1]<<8 + ebp14[2]) // big endian
		default:
			t.r = 0
			return 0
		}

		if ebp4_len <= 0 {
			return 0
		}
		if ebp4_len <= t.r { // 得到一个完整的协议报文了
			t.f4020.f006BDF76(t.bufr[ebp8:], ebp4_len) // 不是解析，仅仅是创建对象
			ebp8 += ebp4_len
			t.r -= ebp4_len
			if t.r > 0 {
				continue
			}
			break
		}
		if ebp8 > 0 && t.r >= 1 {
			f00DE7C90memcpy(t.bufr[:], t.bufr[ebp8:], t.r)
		}
		break
	}

	t.f4020.f006BDDF5()
	return 0
}

func (t *conn) f006BD945write() int {
	if t.w <= 0 {
		return 1
	}
	return 0
}

func (t *conn) f006BDC33() []uint8 {
	return nil
}

// 发送协议报文
func (t *conn) f004397E3write(buf []uint8, len int) int {
	// ebp10 := t // ecx也要落到栈上
	ebp4 := len
	ebp8_sum := 0
	if uintptr(t.fd) == uintptr(^uint32(0)) {
		return 0
	}
	for {
		ebpC, err := win.Send(t.fd, buf[ebp8_sum:], len-ebp8_sum, 0)
		if err != nil {
			if 0x2733 != win.WSAGetLastError() {
				v01319E08log.f00B38AE4printf("[Send Packet Error] WSAGetLastError() != WSAEWOULDBLOCK\r\n")
				//_006BD5BB()
				return 0
			}
			if t.w+len > 0x2000 {
				v01319E08log.f00B38AE4printf("[Send Packet Error] SendBuffer Overflow\r\n")
				//_006BD5BB()
				return 0
			}

			// 发送缓存满了，留着下一次发送
			f00DE7C90memcpy(t.bufw[t.w:], buf, ebp4) // buf不切一下？
			t.w += ebp4
			return 0
		}

		if ebpC == 0 {
			return 1
		}

		ebp8_sum += ebpC
		ebp4 -= ebpC
		if ebp4 <= 0 {
			return 1
		}
	}
	return 0
}

// t4001
type t4001 struct {
	fs []func()
}

// t4000
var v0130F728 t4000

type t4000 struct {
	data  [4880]uint8
	f4880 t4001 // 01313FA8
}

func (t *t4000) f004A9083(p *t4001) {}

func (t *t4000) f004A9123(p *t4001) {
	_ = p.fs[10]
	t.f004A9083(p)
}

func (t *t4000) f004A9B5B() {
	// ...

	// ebp8数组在0x377E3138地址处，应该是new出来的
	// ebp8数组有9个元素，应该是个表
	// 分别存放的是t6000的派生类型，t6000可能是接口类型
	// {_01310056, _01310598, _013180E0, _01317CD8,_01314300, _01313FA8, _0130FF40, _0130FB38, _0130F730}
	// ebp8[ebp24*4].f004CCC07()

	// f004CCC07
	func() {

		// 查表
		var t *t3
		t.f0043E60C() // v01313FA8.f0043E60C 发送login报文
	}()
}

var v08C88E08 uint32 // 可能是状态
var v08C88E0C int
var v012E4018 = "22789"            // 版本怎么会是这个？
var v012E4020 = "M7B4VM4C5i8BC49b" // serial

type t5 struct{}

func (t *t5) f007C483A(x *t6) {}
func (t *t5) f007C4B15()      {}

type t6 struct{}

func (t *t6) f007C522A(x int, y *t1) {}

type t1 struct {
	f18 int
}

func (t *t1) f007BF63B() {
	// f007C47D3
	func() {
		var ebp4 [4]uint8
		// t1.f007C51F5

		// t1.f007C4C15
		func(x []uint8, y []uint8) {
			// f007C572D

			// f007C56F5

			// f007C536A
		}(ebp4[2:], ebp4[3:])
	}()

}

// 好复杂 12500行的汇编代码
func (t *t1) f007BF68F() {}
func (t *t1) f007C2763(flag uint8, buf []uint8, len uint16) {
	var ebpC t5
	t.f007C4A65(&ebpC, &flag)
	var ebp14 t6
	ebpC.f007C483A(t.f007C4858(&ebp14))
	ebpC.f007C4B15()
	// ...
}
func (t *t1) f007C4A65(x *t5, flagp *uint8) {}
func (t *t1) f007C4858(x *t6) *t6 {
	x.f007C522A(t.f18, t)
	return x
}

// pb
type pb struct {
	// len   uint16
	// buf   [5210]uint8 // 0x145A
	buf   [5212]uint8 // 0x145C
	f145C t1
}

var v012E4034conn *conn

func (t *pb) f00439178init() {
	// SEH
	t.f145C.f007BF63B()
	t.f145C.f007BF68F()
	t.f0043921B()
}
func (t *pb) f004391CF() {}
func (t *pb) f0043921B() {
	// t.len = 0
	*(*uint16)(unsafe.Pointer(&t.buf[0])) = 0
}
func (t *pb) f0043922CwritePrefix(flag uint8, code uint8) {
	t.f0043921B()

	// 写flag
	var buf [1]uint8
	buf[0] = flag
	t.f00439298writeBuf(buf[:], 1, false)

	// 写len，没有必要
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

func (t *pb) f00439298writeBuf(buf []uint8, len int, needEnc bool) {
	tlen := (*uint16)(unsafe.Pointer(&t.buf[0]))
	if int(*tlen)+len > 0x145A {
		return
	}
	f00DE7C90memcpy(t.buf[*tlen+2:], buf, len)
	if needEnc {
		t.f0043930Cenc(int(*tlen), int(*tlen)+len, 1) // (3,4,1), (4,8,1)
	}
	*tlen += uint16(len)
}

func (t *pb) f0043930Cenc(begin, end, interval int) {
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

// s9 00439378
func (t *pb) f004393EAsend(needEnc, isC2 bool) {
	// t.f00439612() 写len
	func() {}()

	t.f0043968F() // hook到 igc.dll了，直接ret

	// f00439420
	// hook到 IGC.dll
	func(buf []uint8, len int, needEnc, isC2 bool) {
		f00DE8A70chkstk() // 0x3124

		if !needEnc {
			v012E4034conn.f004397E3write(buf, len)
			return
		}

		// var ebp3124 int
		var ebp3120 int
		var ebp311C int
		var ebp3118 [1000]uint8 // 原始缓存副本
		var ebp1914 int
		// var ebp1910 [1000]uint8 // C4编码缓存
		var ebp108 [1000]uint8 // C3编码缓存

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
		ebp3118[ebp311C-1] = v08C88F60
		v08C88F60++
		ebp311C--

		ebp1914 = v08C88FB4.f00B98ED0(v012E4034conn, nil, ebp3118[ebp311C:], len-ebp311C) // 得到len为0x11

		if ebp1914 < 0x100 && !isC2 {
			ebp3120 = ebp1914 + 2 // 0x13
			ebp108[0] = 0xC3
			ebp108[1] = uint8(ebp3120)
			v08C88FB4.f00B98ED0(v012E4034conn, ebp108[2:], ebp3118[ebp311C:], len-ebp311C) // 编码数据
			v012E4034conn.f004397E3write(ebp108[:], ebp3120)
			return
		}
		// len >= 0x100或者C2 编码方案
		// ...

	}(t.buf[2:], int(*(*uint16)(unsafe.Pointer(&t.buf[0]))), needEnc, isC2)
}

func (t *pb) f0043968F() {
	var ebp4len uint16
	var ebp8buf []uint8
	var ebp9code uint8
	if t.buf[2] == 0xC1 {
		ebp4len = uint16(t.buf[3])
		ebp9code = t.buf[4]
		ebp8buf = t.buf[5:]
		ebp4len -= 3

	}
	if t.buf[2] == 0xC2 {
		ebp4len = uint16(t.buf[3]<<8 + t.buf[4])
		ebp9code = t.buf[5]
		ebp8buf = t.buf[6:]
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

func (t *pb) f0043EDF5writeUint32(data uint32) {
	var buf [4]uint8
	binary.LittleEndian.PutUint32(buf[:], data)
	t.f00439298writeBuf(buf[:], 4, true)
}
func (t *pb) f004C65EFwriteUint16(data uint16) *pb {
	var buf [2]uint8
	binary.LittleEndian.PutUint16(buf[:], data)
	t.f00439298writeBuf(buf[:], 2, true)
	return t
}

// t6000
var v01310798 t6000

type t6000 struct {
	fs []func()
}

func (t *t6000) f00446D6D() {
	// 带SEH处理
	f00DE8A70chkstk() // 0x2994

	// ...

	var ebp149C pb
	ebp149C.f004393EAsend(false, false) // 发送协议报文

	// ...
}

func (t *t6000) f004CCC07() {
	// ...
	_ = t.fs[11]
	_ = t.fs[12] // f00446D6D
}

// t08C88FB4
var v08C88FB4 t08C88FB4

type t08C88FB4 struct {
	f18 int
}

// 被花了
func (t *t08C88FB4) f00B98ED0(p *conn, dst []uint8, buf []uint8, len int) int {

	return 0
}

func (t *t08C88FB4) f00B98D90() {}
