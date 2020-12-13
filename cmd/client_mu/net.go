package main

import (
	"encoding/binary"
	"syscall"

	"github.com/xujintao/balgass/win"
)

var v004633D0conn conn

// conn或者叫netFD
type conn struct {
	hWnd    win.HWND       // m0000
	m04     int            // m0004
	fd      syscall.Handle // m0008
	bufw    [4096]uint8    // m000C
	w       int            // m100C
	bufr    [4096]uint8    // m1010
	r       int            // m2010
	m2014   int
	packets [10]struct { // packets
		flag int
		unk  int
		data [4096]uint8
	}
	// once  sync.Once      // m40EC, v08C8D0DC
}

func (c *conn) f0040CB70() {
	// win.WSAStartup()
}

func (c *conn) f0040CC30socket(hWnd win.HWND) {

}

func (c *conn) f0040CD00dial(ip string, port uint16, msg uint) int {
	/*
		if c.hWnd == 0 {
			dll.user32.MessageBoxA(0, v0046F448msg.Get(144), "Error", d.hWnd) // 错误的窗口句柄
			return 0
		}
		if c.fd == -1 {
			v004633D0conn.f0040CC30socket(d.hWnd)
		}
		dll.ws2_32.htons()
		if dll.ws2_32.inet_addr() == -1 {
			dll.ws2_32.gethostbyname()
		}
		dll.ws2_32.WSAAsyncSelect(c.fd, c.hWnd, msg, FD_CONNECT)
		if SOCKET_ERROR == dll.ws2_32.connect(c.fd) {
			if dll.ws2_32.WSAGetLastError() != 10035 {
				dll.ws2_32.closesocket(d.fd)
				return 0
			}
		}
	*/
	return 1
}

func (c *conn) f0040D010write() {
	win.Send(c.fd, c.bufw[:], c.w, 0)
}

func (c *conn) f0040D090read() int {
	n, err := win.Recv(c.fd, c.bufr[c.r:], 4096-c.r, 0)
	if c.m2014 != 0 {
		// n = c.f0040CE40(c.bufr[:], n)
	}
	if n == 0 {
		return 1
	}
	if err != nil {
		errcode := win.WSAGetLastError()
		if errcode == 10035 { // WSAEWOULDBLOCK
			return 1
		}
		// c.f0040CE50(c, "recv() 받기 에러 %d", errcode)
		return 1
	}
	c.r += n
	if c.r < 3 {
		return 3
	}
	size := 0
	if c.bufr[0] == 0xC1 {
		size = int(c.bufr[1])
	} else if c.bufr[0] == 0xC2 {
		size = int(binary.BigEndian.Uint16(c.bufr[1:]))
	} else {
		c.r = 0
		return 0
	}
	if size <= 0 {
		// f0040CE50(c, "size 가 %d이다.", size)
		return 0
	}

	if size > c.r {
		// f004337A0(c.bufr[:], c.bufr[:], c.r)
		return 0
	}
	// c.f0040CE70()
	func(buf []uint8, size int) int {
		if size > 0x1000 {
			return 2
		}
		if c.w != 0 {

		}
		// f004337A0()
		return 0
	}(c.bufr[:], size)
	return 1
}

func (c *conn) f0040CF30write(buf []uint8, l int) bool {
	if c.fd == syscall.Handle(^uintptr(0)) {
		return false
	}
	for l == 0 {
		n := 0 // n := dll.ws2_32.send(c.fd, buf, l, 0)
		if n == -1 {
			if errcode := win.WSAGetLastError(); errcode != 10035 {
				// f0040CE50(c, "send 받기 에러 %d", errcode)
				c.f0040CCE0close()
				return false
			}
			if c.w+l > 4096 {
				c.f0040CCE0close()
				return false
			}
			// f004337A0(c.bufw[c.w:], buf, l)
			c.w += l
			return false
		}
		l -= n
	}
	return true
}

func (c *conn) f0040CCE0close() {

}

func (c *conn) f0040CEE0getPacket() []uint8 {
	for i := range c.packets {
		if c.packets[i].flag == 1 {
			c.packets[i].flag = 0
			return c.packets[i].data[:]
		}
	}
	return nil
}
