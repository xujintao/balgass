package main

import (
	"sync"

	"github.com/xujintao/balgass/client/dll"
)

// t3004
type t3004 struct {
	f14 uint32 // 0x0DA6D830
	f18 uint32 // 1
}

// 比较复杂，可能是解析
func (t *t3004) _006BDE09() {
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

				}
			}()
		}

	}()
}

// t3003
type t3003 struct{}

func (t *t3003) _006BEB1E(x uint32) {

}

func (t *t3003) _006BEB35(p **t3002, pp **t3002) {
	t._006BF322(p, pp)
}

func (t *t3003) _006BF322(p **t3002, pp **t3002) {
	ebp4 := p

	// _00405BF0
	ebp8 := func(x uint32, p **t3002) {
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
	len uint32
}

func (t *t3002) _006BDFED(buf []uint8, len uint32) *t3002 {
	_00DE7C90_memcpy(t.buf[:], buf, len)
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

func (t *t3001) _006BDDF5() {
	t.f20._006BDE09()
}

func (t *t3001) _006BDF76(buf []uint8, len uint32) {
	// 带SEH了
	var ebp18 *t3002 = _00DE852F(2004) // 可能是new
	var ebp20 *t3002
	if ebp18 == nil {
		ebp20 = nil
	} else {
		ebp20 = ebp18._006BDFED(buf, len) // 可能是初始化
	}
	ebp14 := ebp20
	ebp10 := ebp14
	t._006BE241(&ebp10) // 填表
}

func (t *t3001) _006BE241(pp **t3002) {
	t._006BE6D9(pp)
}

func (t *t3001) _006BE6D9(pp **t3002) {
	if (t.f18+t.f1C)%4 == 0 && t.f14 <= (t.f1C+4)>>2 {
		t._006BEC73(1)
	}
	ebp4 := t.f18 + t.f1C
	ebp8 := ebp4 >> 2
	if t.f14 <= ebp8 {
		ebp8 -= t.f14
	}
	if t.f10[4*ebp8] == 0 {
		t.f10[4*ebp8] = t.f0C._006BEB1E()
	}
	t.f0C._006BEB35(t.f10[4*ebp8][ebp4%4*4], pp)
	t.f1C++
}

func (t *t3001) _006BEC73(x uint32) {

}

// t3000
type t3000 struct {
	hWnd  int         // f0000
	f0004 int         // f0004
	f0008 int         // f0008
	fd    uint32      // f000C
	bufw  [2000]uint8 // f0010
	w     uint32      // f2010
	bufr  [2000]uint8 // f2014
	r     uint32      // f4014
	f4020 *t3001      // f4020
	once  sync.Once   // f40EC, _08C8D0DC
}

var _08C88FF0 t3000

func (t *t3000) _006BD3A7_Init() {
	dll.Ws2_32.WSAStartup()
}

func (t *t3000) _006BD509_Socket(hWnd uintptr, x int) int {
	t.fd = dll.Ws2_32.Socket(2, 1, 0) // AF_INT, SOCK_STREAM, 0
	t.f0004 = x
	if t.fd == -1 { // INVALID_SOCKET
		dll.Ws2_32.WSAGetLastError()
		// log
		// MessageBoxA
		return 0
	}
	t.hWnd = hWnd
	return 1
}

func (t *t3000) _006BD708_Dial(ip string, port int, x int) {

}

func (t *t3000) _006BDA03_Read() uint32 {
	if t.r >= 0x2000 {
		_01319E08_log._00B38AE4("Receive Packet Buffer Overflow")
		return 1
	}
	ebpC := dll.Ws2_32.Recv(t.s, t.bufr[t.r:], 0x2000-t.r, 0)
	if ebpC == 0 {
		// 服务器关闭连接
		return 1
	}
	if ebpC == -1 { // SOCKET_ERROR
		WSAGetLastError()
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
			ebp10 = bufr[ebp8:]
			ebp4_len = ebp10[1]
		case 0xC2, 0xC4:
			ebp14 = bufr[ebp8:]
			ebp4_len = ebp14[1]<<8 + ebp14[2] // big endian
		default:
			t.r = 0
			return 0
		}

		if ebp4_len <= 0 {
			return 0
		}
		if ebp4_len <= t.r { // 得到一个完整的协议报文了
			t.f4020._006BDF76(t.bufr[ebp8:], ebp4_len) // 不是解析，仅仅是创建对象
			ebp8 += ebp4_len
			t.r -= ebp4_len
			if t.r > 0 {
				continue
			}
			break
		}
		if ebp8 > 0 && t.r >= 1 {
			_00DE7C90(t.bufr[:], t.bufr[ebp8:], t.r)
		}
		break
	}

	t.f4020._006BDDF5()
	return 0
}

func (t *t3000) _006BD945_Write() uint32 {
	if t.w <= 0 {
		return 1
	}
}

// 发送协议报文
func (t *t3000) f004397E3_Write(buf []uint8, len int) int {
	// ebp10 := t // ecx也要落到栈上
	ebp4 := len
	ebp8_sum := 0
	if t.fd == -1 {
		return 0
	}
	for {
		ebpC := dll.Ws2_32.Send(t.fd, buf[ebp8_sum:], len-ebp8_sum, 0)
		if ebpC == -1 {
			if 0x2733 != dll.Ws2_32.WSAGetLastError() {
				_01319E08_log._00B38AE4("[Send Packet Error] WSAGetLastError() != WSAEWOULDBLOCK\r\n")
				//_006BD5BB()
				return 0
			}
			if t.w+len > 0x2000 {
				_01319E08_log._00B38AE4("[Send Packet Error] SendBuffer Overflow\r\n")
				//_006BD5BB()
				return 0
			}

			// 发送缓存满了，留着下一次发送
			_00DE7C90_memcpy(t.bufw[t.w:], buf, ebp4) // buf不切一下？
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
var _0130F728 t4000

type t4000 struct {
	data  [4880]uint8
	f4880 t4001 // 01313FA8
}

func (t *t4000) _004A9083(p *t4000) {}

func (t *t4000) _004A9123(p *t4001) {
	p.fs[10]
	t._004A9083(p)
}

func (t *t4000) _004A9B5B() {
	// ...

	// ebp8数组在0x377E3138地址处，应该是new出来的
	// ebp8数组有9个元素，应该是个表
	// 分别存放的是t6000的派生类型，t6000可能是接口类型
	// {_01310056, _01310598, _013180E0, _01317CD8,_01314300, _01313FA8, _0130FF40, _0130FB38, _0130F730}
	ebp8[ebp24*4]._004CCC07()
}

var _08C88E08 uint32
var _08C88E0C uint32
var _012E4018 string = "22789" // 版本怎么会是这个？

// heartbeat
type heartbeat struct {
	len uint16
	buf []uint8
}

var v012E4034 *t3000
var v08C88F60 int

func (t *heartbeat) f00439178(x, y int) {}
func (t *heartbeat) f004391CF()         {}
func (t *heartbeat) f0043922C()         {}

func (t *heartbeat) f004393EA(needEnc, C2 int) {
	t._00439612()
	t._0043968F()

	// _00439420
	func(buf []uint8, len int, needEnc, C2 int) {
		// 0x3124字节的局部变量，还是很复杂的
		_00DE8A70()

		if needEnc == 0 {
			v012E4034.f004397E3_Write(buf, len)
			return
		}

		var ebp3124 int
		var ebp3120 int
		var ebp311C int
		var ebp3118 [1000]uint8 // 原始缓存
		var ebp1914 int
		var ebp1910 [1000]uint8 // C4编码缓存
		var ebp108 [1000]uint8  // C3编码缓存

		// 编码数据
		_00DE7C90_memcpy(ebp3118, buf, len)
		ebp3118[len] = uint8(_00DE8AAD_rand() & 0xFF) // 源码带绝对值
		if C2 == 1 {
			ebp3118[0] = 0xC2
		}
		if ebp3118[0] != 0xC1 {
			ebp311C = 1
		}
		ebp311C += 2
		ebp3118[ebp311C-1] = v08C88F60
		v08C88F60++
		ebp311C--

		ebp1914 = v08C88FB4.f00B98ED0(v012E4034, 0, ebp3118[ebp311C:], len-ebp311C) // 得到len为0x11

		if ebp1914 < 0x100 {
			if C2 == 0 {
				ebp3120 = ebp1914 + 2 // 0x13
				ebp108[0] = 0xC3
				ebp108[1] = ebp3120
				v08C88FB4.f00B98ED0(v012E4034, ebp108[2:], ebp3118[ebp311C:], len-ebp311C) // 编码数据
				v012E4034.f004397E3_Write(ebp108[:], ebp3120)
				return
			}
		}
		// len >= 0x100或者C2 编码方案
		// ...

	}(t.buf, int(t.len), needEnc, C2)
}

func (t *heartbeat) f00439612()         {}
func (t *heartbeat) f0043968F()         {}
func (t *heartbeat) f0043974F(x int)    {}
func (t *heartbeat) f004C65EF(x, y int) {}

// t6000
var v01310798 t6000

type t6000 struct {
	fs []func()
}

func (t *t6000) _00446D6D() {
	// 带SEH处理
	// 0x2994字节的全局变量
	_00DE8A70()

	// ...

	var ebp149C heartbeat
	ebp149C.f004393EA(0, 0) // 发送协议报文

	// ...
}

func (t *t6000) _004CCC07() {
	// ...
	t.fs[11]
	t.fs[12] // _00446D6D
}

// t08C88FB4
var v08C88FB4 t08C88FB4

type t08C88FB4 struct {
	f18 int
}

func (t *t08C88FB4) f00B98ED0(p *t3000, x int, buf []uint8, len int) int {
	// 被花了
}

func (t *t08C88FB4) f00B98D90() {}
