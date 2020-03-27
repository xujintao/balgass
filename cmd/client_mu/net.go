package main

import (
	"encoding/binary"
	"syscall"

	"github.com/xujintao/balgass/win"
)

var v004633D0conn conn

// conn或者叫netFD
type conn struct {
	hWnd  win.HWND       // m0000
	m04   int            // m0004
	fd    syscall.Handle // m0008
	bufw  [1000]uint8    // m000C
	w     int            // m100C
	bufr  [1000]uint8    // m1010
	r     int            // m2010, 0x004653E0
	data  [4]int
	frame [2000]uint8 // m2024, 0x004653F4
	// w    int         // m3014
	// m4020 *t3001         // m4020
	// once  sync.Once      // m40EC, v08C8D0DC
}

func (c *conn) f0040CB70() {

}

func (c *conn) f0040CC30socket(hWnd win.HWND) {

}

func (c *conn) f0040CD00dial(ip string, port uint16, msg uint) {

}

func (c *conn) f0040D010write() {
	win.Send(c.fd, c.bufw[:], c.w, 0)
}

func (c *conn) f0040D090read() {
	n, err := win.Recv(c.fd, c.bufr[c.r:], 0x1000-c.r, 0)
	if err != nil {
		win.WSAGetLastError()
	}
	c.r = n
	if n < 3 {
		return
	}
	size := 0
	if c.bufr[0] == 0xC1 {
		size = int(c.bufr[1])
	} else if c.bufr[0] == 0xC2 {
		size = int(binary.BigEndian.Uint16(c.bufr[1:]))
	} else {
		return
	}
	// c.f0040CE70()
	func(buf []uint8, size int) int {
		if size > 0x1000 {
			return 2
		}
		if c.w != 0 {

		}
		// f004337A0handle()
		func() {

		}()
		return 0
	}(c.bufr[:], size)
}

func (c *conn) f0040CCE0close() {

}
