package main

import (
	"sync"
	"syscall"

	"github.com/xujintao/balgass/win"
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
	var ebp18 *t3002 = (*t3002)(f00DE852Fnew(2004)) // 可能是new
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

var v08C88FF0conn conn

// conn或者叫netFD
type conn struct {
	hWnd  win.HWND       // f0000
	f04   int            // f0004
	f08   int            // f0008
	fd    syscall.Handle // f000C
	bufw  [2000]uint8    // f0010
	w     int            // f2010
	bufr  [2000]uint8    // f2014
	r     int            // f4014
	f4020 *t3001         // f4020
	once  sync.Once      // f40EC, v08C8D0DC
}

var v08C8D0DC bool // like once

func (t *conn) f006BD3A7init() {
	t.once.Do(func() {
		win.WSAStartup(0, nil)
	})
	// win.WSAStartup(0, nil)
}

func (t *conn) f006BD509socket(hWnd win.HWND, x int) int {
	var err error
	t.fd, err = win.Socket(2, 1, 0) // AF_INT, SOCK_STREAM, 0
	t.f04 = x
	if err != nil { // INVALID_SOCKET
		win.WSAGetLastError()
		// log
		// MessageBoxA
		return 0
	}
	t.hWnd = hWnd
	return 1
}

func (t *conn) f006BD5BBclose() bool {
	if t.f04 != 0 {
		v08C88F64 = 0
	}
	// ebp4 := struct {
	// 	onoff  uint16
	// 	linger uint16
	// }{}
	// ebp4.onoff = 1
	// ebp4.linger = 0
	// setsocketopt(t.fd, SOL_SOCKET, 0x80, &ebp4, 4)
	f00DE8100memset(t.bufr[:], 0, 4)
	f00DE8100memset(t.bufw[:], 0, 4)
	t.w = 0
	t.r = 0
	// for t.f4020.f006BDF65() == false {
	// 	t.f4020.f006BDF28()
	// }
	v01319E08log.f00B38AE4printf("[Socket CLosed][Clear PacketQueue]\r\n")
	// v08C88FB4decrypt.f00B997E0()
	// v08C8D014encrypt.f00B997E0()
	// closesocket(t.fd)
	// t.fd = syscall.Handle(-1)
	return true
}

func (t *conn) f006BD6F9fd() syscall.Handle {
	return t.fd
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
	ebp8sum := 0
	if uintptr(t.fd) == uintptr(^uint32(0)) {
		return 0
	}
	for {
		ebpC, err := win.Send(t.fd, buf[ebp8sum:], len-ebp8sum, 0)
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

		ebp8sum += ebpC
		ebp4 -= ebpC
		if ebp4 <= 0 {
			return 1
		}
	}
	return 0
}

var v08C88E08 uint32 // 可能是状态

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

// f006BF89A 拨号
func f006BF89ADial(ip string, port int) {
	v08C88FF0conn.f006BD3A7init()
	v01319E08log.f00B38AE4printf("[Connect to Server] ip address = %s, port = %d\r\n", ip, port)
	v08C88FF0conn.f006BD509socket(v01319D6ChWnd, 1)
	v08C88FF0conn.f006BD708dial(ip, port, 400)
	v08C88F60 = 0
	v08C88F61 = 0
}
