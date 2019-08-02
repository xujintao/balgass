package main

// t2004
type t2004 struct {
	f14 uint32 // 0x0DA6D830
	f18 uint32 // 1
}

// 比较复杂，可能是解析
func (t *t2004) _006BDE09() {
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

// t2003
type t2003 struct{}

func (t *t2003) _006BEB1E(x uint32) {

}

func (t *t2003) _006BEB35(p **t2002, pp **t2002) {
	t._006BF322(p, pp)
}

func (t *t2003) _006BF322(p **t2002, pp **t2002) {
	ebp4 := p

	// _00405BF0
	ebp8 := func(x uint32, p **t2002) {
		return p
	}(4, p)
	if ebp8 == nil {
		return
	}
	*p = *pp
}

// t2002
type t2002 struct {
	buf [2000]uint8
	len uint32
}

func (t *t2002) _006BDFED(buf []uint8, len uint32) *t2002 {
	_00DE7C90_memcpy(t.buf[:], buf, len)
	t.buf[len] = 0xFD
	t.len = len
	return t
}

// t2001
type t2001 struct {
	data [12]uint8
	f0C  *t2003      // 0x0D906A78
	f10  [][]**t2002 // 指针数组, 0x2817EFA0
	f14  uint32
	f18  uint32
	f1C  uint32
	f20  *t2004 // 0x0DA6FFC0
}

func (t *t2001) _006BDDF5() {
	t.f20._006BDE09()
}

func (t *t2001) _006BDF76(buf []uint8, len uint32) {
	// 带SEH了
	var ebp18 *t2002 = _00DE852F(2004) // 可能是new
	var ebp20 *t2002
	if ebp18 == nil {
		ebp20 = nil
	} else {
		ebp20 = ebp18._006BDFED(buf, len) // 可能是初始化
	}
	ebp14 := ebp20
	ebp10 := ebp14
	t._006BE241(&ebp10) // 填表
}

func (t *t2001) _006BE241(pp **t2002) {
	t._006BE6D9(pp)
}

func (t *t2001) _006BE6D9(pp **t2002) {
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

func (t *t2001) _006BEC73(x uint32) {

}

// t2000
type t2000 struct {
	f0000 [12]uint8   // f0000
	fd    uint32      // f000C
	bufw  [2000]uint8 // f0010
	w     uint32      // f2010
	bufr  [2000]uint8 // f2014
	r     uint32      // f4014
	f4020 *t2001      // f4020
}

var _08C88FF0 t2000

func (t *t2000) _006BDA03_Read() uint32 {
	if t.r >= 0x2000 {
		_01319E08_log._00B38AE4("Receive Packet Buffer Overflow")
		return 1
	}
	ebp_C := recv(t.s, t.bufr[t.r:], 0x2000-t.r, 0)
	if ebp_C == 0 {
		// 服务器关闭连接
		return 1
	}
	if ebp_C == -1 { // SOCKET_ERROR
		WSAGetLastError()
		//...
		return 1
	}
	t.r += ebp_C
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

func (t *t2000) _006BD945_Write() uint32 {
	if t.w <= 0 {
		return 1
	}
}

// t3001
type t3001 struct {
	fs []func()
}

// t3000
var _0130F728 t3000

type t3000 struct {
	data  [4880]uint8
	f4880 t3001 // 01313FA8
}

func (t *t3000) _004A9083(p *t3000) {

}

func (t *t3000) _004A9123(p *t3001) {
	p.fs[10]
	t._004A9083(p)
}

var _08C88E08 uint32
var _08C88E0C uint32
var _012E4018 string = "22789" // 版本怎么会是这个？
