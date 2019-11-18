package main

import (
	"unsafe"

	"github.com/xujintao/balgass/win"
)

type point struct {
	x int
	y int
}
type rect struct {
	left   int
	top    int
	right  int
	bottom int
}
type servicer interface {
	do4(int) bool
	do6(bool)
	do10() bool
	do11(float64)
	do12()
	do13(float64)
}

type subservice struct {
	mA9done bool
	mAA     bool
	mAB     bool
}

func (s *subservice) f00448976() bool { return false }
func (s *subservice) f00436088() bool {
	if s.mA9done {
		return s.f00448976()
	}
	return false
}
func (s *subservice) f004CCE44(x bool)  { s.mA9done = x }
func (s *subservice) f004360AB(float64) {} // 比较复杂
func (s *subservice) f0043912C() bool   { return s.mAA }
func (s *subservice) f00438F7B(x bool)  { s.mAB = x }

type windower interface {
	f004CCC07(float64)
	f004CCB8A(bool)
	f004AA018isActive() bool
	f004AA068() int
	f004AA027() bool
}

// window
// window需要实现windower接口
// window需要继承servicer接口同时各派生类需要实现servicer接口
type window struct {
	servicer
	m04x      int
	m08y      int
	m0Cactive bool
	m0D       bool
	m0E       bool
	m10       int
	m18left   int // rect left
	m1Ctop    int // rect top
	m20width  int // rect width
	m24heigth int // rect height
	m28left   int // = m18left
	m2Ctop    int // = m1Ctop
	subs      list
}

func (b *window) f004CCA35(x int) bool {
	if b.m0Cactive == false {
		return false
	}
	ebp18 := f0043BF3FgetT4003()
	// ebp10rect := rect{}
	switch x {
	case 0:
		// SetRect(&ebp10rect, b.m18left, b.m1Ctop, b.m18left+b.m20width, b.m1Ctop+b.m24heigth)
		ebp24pt := point{}
		ebp18.f0043BE81getPoint(&ebp24pt)
		// return PtInRect(&ebp10rect, ebp24pt.x, ebp24pt.y)
	case 2:
		// SetRect(&ebp10rect, b.m18left, b.m1Ctop, b.m18left+b.m20width, b.m1Ctop+0x1A)
		ebp2Cpt := point{}
		ebp18.f0043BE81getPoint(&ebp2Cpt)
		// return PtInRect(&ebp10rect, ebp2Cpt.x, ebp2Cpt.y)
	case 3:
		ebp14 := b.subs.f004409AAgetList()
		for ebp14 != nil {
			ebp1C := b.subs.f00445530getNodeValue(&ebp14)
			if ebp1C.(*subservice).f00436088() == true {
				return true
			}
		}
		b.do10()
	}
	return false
}

// do4 虽然每个service实现了do4，但是实现细节是一样的
// 所以由window实现供各个service调用就行
func (b *window) do4(x int) bool {
	if b.m0Cactive == false {
		return false
	}
	if x == 2 {
		return false
	}
	return b.f004CCA35(x)
}

// f004CCB8A 设置所有子业务done标识
func (b *window) f004CCB8A(x bool) {
	// init sub list
	ebpC := b
	ebp4 := ebpC.subs.f004409AAgetList()
	for ebp4 != nil {
		ebp8 := ebpC.subs.f00445530getNodeValue(&ebp4)
		ebp8.(*subservice).f004CCE44(x)
	}
}

func (b *window) f004CCC07(unk float64) {
	if b.m0Cactive == false {
		return
	}
	ebp4 := f0043BF3FgetT4003()
	if ebp4.f004366A5() {
		b.m10 = 0
	}
	if b.m10 == 0 {
		ebp8 := b.subs.f004409AAgetList()
		for ebp8 != nil {
			ebpC := b.subs.f00445530getNodeValue(&ebp8)
			ebpC.(*subservice).f004360AB(unk)
		}
	}
	b.do11(unk)
	if b.m0D == false {
		return
	}
	if ebp4.f00436696() {
		if b.do4(2) {
			ebp14pt := point{}
			ebp4.f0043BE81getPoint(&ebp14pt)
			b.m04x = ebp14pt.x
			b.m08y = ebp14pt.y
			b.m28left = b.m18left
			b.m2Ctop = b.m1Ctop
			b.m10 = 2
		}
	}
	if b.m10 == 2 {
		b.m28left = ebp4.f00448D8BgetX() - b.m04x + b.m28left
		b.m2Ctop = ebp4.f00448191getY() - b.m08y + b.m2Ctop
		if b.m0E == false {
			// b.do3(b.m2Ctop, b.m28left)
		}
		ebp1C := point{}
		ebp4.f0043BE81getPoint(&ebp1C)
		b.m04x = ebp1C.x
		b.m08y = ebp1C.y
	}
	b.do12()
	b.do13(unk)
}
func (b *window) f004AA018isActive() bool { return b.m0Cactive }
func (b *window) f004AA068() int          { return b.m10 }
func (b *window) f004AA027() bool         { return b.m0D }

// service1 sizeof(service{})=0x408
type service1 struct {
	window
}

// do4->f0043F608
func (s *service1) do4(x int) bool {
	return s.window.do4(x)
}
func (s *service1) do6(bool)         {}
func (s *service1) do10() bool       { return false }
func (s *service1) do11(unk float64) {}
func (s *service1) do12()            {}

// do13->f0043F640
func (s *service1) do13(unk float64) {}

// service2 sizeof(service{})=0x408
type service2 struct {
	window
}

// do4->f0043AA4E
func (s *service2) do4(x int) bool {
	return s.window.do4(x)
}
func (s *service2) do6(bool)         {}
func (s *service2) do10() bool       { return false }
func (s *service2) do11(unk float64) {}
func (s *service2) do12()            {}
func (s *service2) do13(unk float64) {}

// service3 sizeof(service{})=0x408
type service3 struct {
	window
}

// do4->f0043E5D4
func (s *service3) do4(x int) bool {
	return s.window.do4(x)
}
func (s *service3) do6(bool)         {}
func (s *service3) do10() bool       { return false }
func (s *service3) do11(unk float64) {}
func (s *service3) do12()            {}

// do->f0043E60C 发送login报文
func (s *service3) do13(unk float64) {}

// service4 sizeof(service{})=0x408
type service4 struct {
	window
	subs  [62]subservice
	m3760 int
	m3764 int
}

// do4->f00446D35
func (s *service4) do4(x int) bool {
	return s.window.do4(x)
}
func (s *service4) do6(bool)         {}
func (s *service4) do10() bool       { return false }
func (s *service4) do11(unk float64) {}
func (s *service4) do12()            {}

// do13->f00446D6D
func (s *service4) do13(unk float64) {
	// 带SEH处理
	f00DE8A70chkstk() // 0x2994
	// ...
	ebp10 := 0
	for ebp10 < 0x15 {
		if s.subs[ebp10].f0043912C() {
			if s.m3760 != 0 {
				s.subs[s.m3760].f00438F7B(false)
			}
			s.subs[ebp10].f00438F7B(true)
			s.m3760 = ebp10

			var ebp149C pb
			ebp149C.f00439178init()
			ebp149C.f0043922CwritePrefix(0xC1, 0xF4)
			var ebp15 [1]uint8
			ebp15[0] = 6
			ebp149C.f00439298writeBuf(ebp15[:], 1, false)
			ebp149C.f004393EAsend(false, false) // 发送协议报文
			ebp149C.f004391CF()
		}
		ebp10++
	}
	if s.m3764 == 0 {
		return
	}
}

// service5 v01310598
type service5 struct {
	window
	s1 subservice // offset: 0x040, v013105D8
	s2 subservice // offset: 0x120, v013106B8
}

// do4->f0043DCD3
func (s *service5) do4(x int) bool {
	return s.window.do4(x)
}

// do6->f00435D53
func (s *service5) do6(x bool) {
	s.window.m0D = x
}
func (s *service5) do10() bool { return false }

// do11->f00435D86
func (s *service5) do11(unk float64) {
	// return 0.0 - unk
}

// do12->f00435D96
func (s *service5) do12() {}

// do13->f0043DD0B
func (s *service5) do13(unk float64) {
	// seh
	f00DE8A70chkstk() // 0x1494
	if 0.0-unk > 0.41 {
		return
	}
	if s.s1.f0043912C() {
		ebp10 := f004A7D34getServiceManager()
		ebp10.f004A9123(&ebp10.s7)
		ebp10.f00439161(1)
	}
	if s.s2.f0043912C() {
		// 二次请求服务器列表
		var ebp149C pb
		ebp149C.f00439178init()
		// ebp4 := 0
		ebp149C.f0043922CwritePrefix(0xC1, 0xF4)
		var ebp15 [1]uint8
		ebp15[0] = 6
		ebp149C.f00439298writeBuf(ebp15[:], 1, false)
		ebp149C.f004393EAsend(false, false)
		ebp149C.f004391CF()
		ebp14 := f004A7D34getServiceManager()
		ebp14.f004A9123(&ebp14.s2)
		// f004D4F77(v012E239C, 0) // main_theme.mp3
		// f004D4FB9(v012E234C, 0) // mutheme.mp3
	}
}

// service6 sizeof(service{})=0x408
type service6 struct {
	window
}

// do4->f00441096
func (s *service6) do4(x int) bool {
	return s.window.do4(x)
}
func (s *service6) do6(bool)         {}
func (s *service6) do10() bool       { return false }
func (s *service6) do11(unk float64) {}
func (s *service6) do12()            {}
func (s *service6) do13(unk float64) {}

// service7 sizeof(service{})=0x408
type service7 struct {
	window
}

// do4->f00449300
func (s *service7) do4(x int) bool {
	return s.window.do4(x)
}
func (s *service7) do6(bool)         {}
func (s *service7) do10() bool       { return false }
func (s *service7) do11(unk float64) {}
func (s *service7) do12()            {}
func (s *service7) do13(unk float64) {}

type node struct {
	next  *node
	prev  *node
	value interface{}
}
type list struct {
	head *node
	tail *node
	num  uint32 // service num
}

func (t *list) f004AA077isNumZero() bool {
	if t.num > 0 {
		return false
	}
	return true
}

func (t *list) f00445530getNodeValue(nodep **node) interface{} {
	// ebp8 := t
	ebp4 := *nodep
	// if IsBadReadPtr(ebp4, ebp4) {
	// 	return nil
	// }
	*nodep = ebp4.next
	return ebp4.value
}

func (t *list) f004409AAgetList() *node {
	return t.head
}

func (t *list) f004AA009getNodeNum() uint32 {
	return t.num
}

func (t *list) f004452A7getFirstNodeValue() interface{} {
	ebp4 := t
	if ebp4.head == nil {
		return nil
	}
	return ebp4.head.value
}

type t4002 struct {
	f04 [10]uint8
	f18 int
}

func (t *t4002) f00406EB0(buf []uint8, len int) {
	if t.f18 < 16 {
		// if buf < t.f04[:] {

		// }
	}
	if t.f18 < 16 {
		// t.f04[]
	}
}

func (t *t4002) f00406F90(buf []uint8) {
	len := 0
	for _, v := range buf {
		if v == 0 {
			break
		}
		len++
	}

	t.f00406EB0(buf, len)
}

func (t *t4002) f0043D7E2(buf []uint8) {
	t.f00406F90(buf)
}

// serviceManager
var v0130F728 serviceManager

type serviceManager struct {
	s001 service1 // v0130F730
	// v0130F8C0 v0130F9A0

	s7 service7 // v0130FB38
	// v0130FBC0 v0130FCA0 v0130FD80 v0130FE60

	s6 service6 // v0130FF40
	// v0130FFC8 v01310268

	s5 service5 // v01310598, 激活频率高
	// v013105D8
	// v013106B8

	s4 service4 // v01310798
	// v013107D8 v013108B8 v01310998 ... 64个

	s3 service3 // v01313FA8
	// v01313FE8
	// v013140C8

	s2 service2 // offset:0x4BD8, v01314300
	// v013147D8

	s010 service1 // v01317CD8
	// v01317E68
	// v01317F48

	s100 service1 // v013180E0
	// v01318270
	// v01318350

	// ... // 100个service
	f9FD4activeServices list // v013196FC
	f9FE0               bool // v01319708
	f9FE1               bool
	f9FE4               int  // v0131970C
	f9FE8               bool // v01319710
	f9FE9               bool
	f9FEC               *t4002 // v01319714
}

var v01319730 uint32 // sync.Once

func f004A7D34getServiceManager() *serviceManager {
	if v01319730&1 == 0 {
		v01319730 |= 1
		v0130F728.f004A7A82()
		// f00DE8BF6(0x0114817A)
	}
	return &v0130F728
}
func (t *serviceManager) f004A7A82() {}
func (t *serviceManager) f004A9083(s interface{}) interface{} {
	ebp4 := t.f9FD4activeServices.f004452A7getFirstNodeValue()
	if ebp4 == nil {
		return nil
	}
	if ebp4.(windower).f004AA027() {
		ebp4.(servicer).do6(false)
	} else {
		ebp4 = nil
	}
	if s.(windower).f004AA018isActive() {
		// if t.f9FD4activeServices.f004455FB(t.f9FD4activeServices.f004457B6(s, 0)) == false {
		// 	return nil
		// }
		t.f9FE8 = true
		// t.f9FD4activeServices.f004453C6(s)
	}
	return ebp4
}
func (t *serviceManager) f004A9123(s interface{}) {
	// s.(servicer).do5(1)
	t.f004A9083(s)
}
func (t *serviceManager) f00439161(x int) {}
func (t *serviceManager) f004A91CE() {
	ebp98 := t
	ebp8 := ebp98.f9FD4activeServices.f004409AAgetList()
	for ebp8 != nil {
		ebp34 := ebp98.f9FD4activeServices.f00445530getNodeValue(&ebp8)
		if ebp34.(windower).f004AA068() != 2 {
			return
		}
		// ...
	}
}
func (t *serviceManager) f004A9B5B(unk float64) {
	ebp30 := t
	if ebp30.f9FE4 == 0 { // 0, 2
		return
	}
	if ebp30.f9FD4activeServices.f004AA077isNumZero() == true {
		return
	}
	if ebp30.f9FE8 { // 0, 1
		ebp18 := ebp30.f9FD4activeServices.f004452A7getFirstNodeValue()
		if ebp18.(windower).f004AA018isActive() {
			ebp18.(servicer).do6(true)
			ebp30.f9FE8 = false
		}
	}

	ebp10 := f0043BF3FgetT4003() // v01308D18
	if ebp10.f0043913E(0x1B) {   // always be false
		// ...
	}
	ebp30.f9FE0 = false
	if ebp10.f00436696() {
		// ...
	}
	if ebp10.f004366A5() {
		ebp30.f9FE1 = false
	}
	ebp14 := ebp30.f9FD4activeServices.f004AA009getNodeNum() // 9
	ebp8 := make([]interface{}, 9)                           // 0x37000C68 new?
	ebp4 := ebp30.f9FD4activeServices.f004409AAgetList()     // 0x36C8E8E8
	var ebp20 uint32
	for ebp20 < ebp14 {
		ebp8[ebp20] = ebp30.f9FD4activeServices.f00445530getNodeValue(&ebp4)
		ebp8[ebp20].(windower).f004CCB8A(false)
		ebp20++
	}

	ebp4 = ebp30.f9FD4activeServices.f004409AAgetList()
	for ebp4 != nil {
		ebpC := ebp30.f9FD4activeServices.f00445530getNodeValue(&ebp4)
		if ebpC.(servicer).do4(0) {
			ebpC.(windower).f004CCB8A(true)
			break
		}
	}

	var ebp24 uint32
	for ebp24 < ebp14 {
		ebp8[ebp24].(windower).f004CCC07(unk) // 重要
		ebp24++
	}

	if len(ebp8) > 0 {
		// f00DE7BEA(ebp8) // delete?
		ebp8 = ebp8[:0]
	}
	ebp30.f004A91CE()
	ebp4 = ebp30.f9FD4activeServices.f004409AAgetList()
	for ebp4 != nil {
		ebpC := ebp30.f9FD4activeServices.f00445530getNodeValue(&ebp4)
		ebp34 := ebpC.(windower).f004AA068()
		switch ebp34 {
		case 1, 2, 3, 4:
			ebp30.f9FE0 = true
		}

		if ebp30.f9FE0 {
			break
		}
		if ebpC.(servicer).do4(0) == true {
			ebp30.f9FE0 = true
			break
		}
	}
}

func (t *serviceManager) f004A9F3B(buf []uint8) {
	if len(buf) == 0 {
		return
	}
	if t.f9FE4 == 4 {
		// t.f9DD8.f00445A2A(buf)
	} else {
		t.f9FEC.f0043D7E2(buf)
	}
}

var v01308D18 t4003
var v01308D80 uint32 // sync.Once

type t4003 struct {
	m0Cx int
	m10y int
	m1C  bool
	m1E  bool
	m31  bool
}

func f0043BF3FgetT4003() *t4003 {
	if v01308D80&1 != 0 {
		return &v01308D18
	}
	v01308D80 |= 1
	v01308D18.f0043BF18init()
	// f00DE8BF6(&v01148111)
	return &v01308D18
}
func (t *t4003) f0043BF18init() {}
func (t *t4003) f0043913E(x uint32) bool {
	if t.m31 == false {
		// f008AEFAD(x)
		return false
	}
	return false
}
func (t *t4003) f00436696() bool { return t.m1C }
func (t *t4003) f004366A5() bool { return t.m1E }
func (t *t4003) f0043BE81getPoint(p *point) *point {
	p.x = t.m0Cx
	p.y = t.m10y
	return p
}
func (t *t4003) f00448D8BgetX() int { return t.m0Cx }
func (t *t4003) f00448191getY() int { return t.m10y }

func f004DD578handleState1(hDC win.HDC) {
	// ...
	v01319E08log.f00B38AE4printf("> Loading ok.\r\n")
	// f004DAACA(v01319D6ChWnd)
	v012E2340 = 2
}

// s9 0x004DEDAD
func f004E46B3(hDC win.HDC) {
	if v0131A26C == false {
		return
	}
	// ...
	// 1.04R location: 0x004E4819
	// s9    location: 0x004DEED2
	// f0051B429 hook-> f0AD98B31
	func() {
		// 0x0A8FD064
		var label1 uint32 = 0x0B07473F
		// push label1
		// push 0x0ABF9BB5
		// ret

		// 0x0ABF9BB5
		var label2 uint32 = 0x00955C2C
		// push label2
		// push 0x0AAB3582
		// ret

		// 0x0AAB3582
		var label3 uint32 = 0x0AD33A1B
		// push label3
		// push 0x0A4EAF02
		// ret

		// 0x0A4EAF02
		// push eax
		// push edx
		// pushfd
		// push esi
		// push edi
		// push ebx
		// push ecx

		// 0x0AF837C2
		// push 0x0A4E2E5E
		// push 0x0A556D38
		// ret

		// 0x0A556D38
		func() {
			// 0x0AD3DEEC
			// push esi
			// 0x09FDD040 0x0A902E30 0x0A04CEA5 0x0A9F7A5E 0x0ABF73ED
			if v09FB8736 != v0A4E24C6 {
				// 0x0AF940DF
				// push eax
				// 0x0AD84F85
				// push edx
				// 0x0A334194
				v0B06F40C = v09E035DB ^ v0A43DD91
				// 0x0B06F40C
				// rdtsc
				var tscLeax uint32 = 0x80AA05B4
				// 0x0A4E4C04
				v09FDFB22 = tscLeax
				// pop edx
				// pop eax
				v09FB8736 = v0A4E24C6
			}
			// 0x0AD3F139
			if v09FDFB22 == 0 {
				// 0x0ABF91B8
				v09FDFB22++
			}
			// 0x0ABB5E49
			v09FDFB22 = (v0AFD3C52 * v09FDFB22) % v09FB6D69 // v09FDFB22 = 0xC8D6
			if v09FDFB22 <= v0A441BD1 {
				// 引导 trap message
				// 0x0A904725
				ebpC := v0AD2DD3A
				ebp14 := v0ABFAB88blocks[:]
				// 0x0A88F009
				label1 = v0B287022label1 // *(ebp + v0AF890C3*4 + 8) = v0B287022
				// 0x09FE38AE
				for {
					if ebp14[0].addr == ^uintptr(0) {
						break // 0x0AC33DD2
					}
					// 0x0AAB2B15
					if ebp14[0].addr == ^uintptr(1) {
						// 0x0A88C1C8
						// ebp14 = ebp14[0].size + v0A4E85ABimageBase
						continue
					}
					// 0x0AF86C65
					ebp20 := ebp14[0].addr + v0A4E85ABimageBase
					// 0x09FC5349 0x09FC3291
					ebp10 := ebp14[0].size
					// ebp1C := ebp14[0].addr
					// 0x0AF949B9 0x0A84B475 0x0AA092FB 0x0A94F4DC 0x0A846F88
					// 0x0AF7D073
					for ebp10 >= 4 {
						ebpC *= (*(*uint32)(unsafe.Pointer(ebp20))) // 4 byte block
						ebp20 += 4
						ebp10 -= 4
					}
					// 0x0AD75422
					for ebp10 > 0 {
						// 0x0ABD786A
						ebpC *= (uint32)(*(*uint8)(unsafe.Pointer(ebp20))) // 1 byte block
						ebp10--
						ebp20++
					}
					// 0x09FDF2EA
					ebp14 = ebp14[1:] // next block
				} // for loop 0x09FE38AE

				// 0x0AC33DD2
				label2 = v0AAB7324label2
				ebp18 := &ebpC
				ebp4 := &v0AA2E05D
				ebp8 := v09EB9D65 - 1
				for {
					// 0x09FB976E 0x09EB8681 0x0A5561A8 0x0A9F8156 0x0AB4D04C 0x0AA05397 0x0AD3ED19
					if ebp8 >= 0 {
						// 0x0A9D4F03
						// *(ebp8*4 + &v0A84EAA4) = edx
						// 0x0A55E6BE
						ebp8--
						// 0x09FB976E
					} else {
						break // 0x09F8EDD8
					}
				}
				// 0x09F8EDD8
				ebp8 = v09EB9D65 - 1
				// 0x0AD2DAD0
				for {
					// 0x0A057AD3 0x0ABB5BA4 0x0A6045A4 0x0AC38DB6 0x0AAB75A9 0x0A4EBF51
					if ebp8 >= 0 {
						// 0x0AFDC80E
						if *ebp18 == *ebp4 {
							// 0x0AD8EC9E 0x0A4E2236
							ebp8--
							// 0x0A057AD3
						} else {
							// 0x0A33AA09
							// push v0A936E16 // 0x2715
							// push 0x0A849753
							// push v0A4E002C // 0x007C6515
							// ret

							// 0x007C6515 0x0B10BD49
							// var label2 uint32 = 0x0069A764
							// push label2
							// push 0x0A05E2CB
							// ret

							// 0x0A05E2CB
							// push 0x008A2556
							// push 0x09F8A379
							// ret

							// 0x09F8A379
							// push 0x0AF87992
							// push 0x0A9FD0D8
							// ret

							// 0x0A9FD0D8 0x09E7064B 0x0A4ED56A 0x0A9035D7 0x0B10DDBC 0x0A43C3DC
							// push edx
							// push edi
							// push ebx
							// push ecx
							// push esi
							// pushfd
							// 0x0A8919E3
							// push 0x0ABDA08C
							// push 0x09FC49DC
							// ret

							// 0x09FC49DC 0x0B10B4CC 0x0A0585C1 0x0B2861D1 0x0A5619F2
							// func() {
							// 	// 0x09E8E8CE 0x09E6B4AB
							// 	ecx := v09FFD056
							// 	eax := v09FFD052
							// 	edx := v0AF947E0
							// 	// push ebx
							// 	// push esi
							// 	// push edi
							// 	// 0x0ABB6C9F
							// 	ebpC := ecx
							// 	ebp10 := eax
							// 	eax = v0A8FB710label2
							// 	// 0x0AF13A88
							// 	ecx = &ebp10
							// 	// push 0x0A84F076
							// 	// 0x0B286D02 0x0AF95B7F
							// 	// push ecx
							// 	ebx := 0x09F94DAA
							// 	ebp10 = 0x09FFD052
							// 	label2 = v0A8FB710label2
							// 	// 0x0ABD639F
							// 	// push 0x0A891A5A
							// 	// push 0x0AD9788A
							// 	// ret

							// 	// 0x0AD9788A
							// 	// 一个迭代
							// }()

						}
					}
					// 0x0A9D43C2 发完trap message后会走到这里
					v0AFD365B = v0A9360CB
					label3 = v0A935B85label3 // [ebp+v0A933705*4+8] = v0A935B85
					break                    // 0x0A888DC9
				}
			} else {
				// 0x0A83F3CE
				label1 = v0B287022label1 // *(ebp + v0AF890C3*4 + 8) = v0B287022
				// 0x0A131864
				label2 = v0AAB7324label2 // *(ebp + v09E8DF92*4 + 8) = v0AAB7324
				// 0x0AD97228
				label3 = v0A935B85label3 // *(ebp + v0A933705*4 + 8) = v0A935B85
				// 0x0A888DC9
			}
			// 0x0A888DC9
			// pop esi
		}()

		// 0x0A4E2E5E
		// pop ecx
		// pop ebx
		// pop edi
		// pop esi
		// popfd
		// pop edx
		// pop eax
		// 0x0A9D45B8 0x0AF0F2C8
		// label3(0x0A43B47B)
		// label2(0x0AD2CAED)
		// label1(0x0AD98B31)
		// f0AD98B31 被掩藏的函数
		func() {
			ebp4 := 0
			for {
				// 0x0051B43B
				if ebp4 >= 10 {
					break
				}
				ebp8 := &v08C86C50[ebp4]
				// f0051B25B
				eax := func(x *t7, y uint32) bool {
					return true
				}(ebp8, 0)
				if eax == false {
					break
				}
				ebp4++
			}
		}()
	}()
	// ...
}

func f004E4F1ChandleState245(hDC win.HDC) {
	// SEH
	// f00552D0D()
	ebp178 := v0131A270 // 0, 0x28, 0x44
	for ebp178 >= 0x28 {
		// f008AF00D().f008AF06A() // v09D24B80.f008AF06A()
		if v012E2340 == 2 || v012E2340 == 4 {
			// v01319D8C.f00A08B5D()
			var ebp184, ebp408 float64
			if v0114EC40-ebp184 == 0x41 {
				ebp408 = ebp184
			} else {
				ebp408 = v0114EC40
			}
			ebp184 = ebp408 // 0x4069 << 20
			// v01319D8C.f00A08BF0()
			// f0043BF3FgetT4003().f0043C06B() // v01308D18.f0043C06B()
			f004A7D34getServiceManager().f004A9B5B(ebp184)
		}
		// v01308ED4 = 0
		ebp40C := v012E2340
		switch ebp40C {
		case 2:
			// f004E1E1E()
			func() {
				// 带SEH
				f00DE8A70chkstk() // 0x46E8
				if !v0131A26C {
					v0131A26C = true
					// f004E1CEE()
					func() {
						// ...
						// 0x004E1D36 hook到 dll
						// f006BF89A 拨号
						f006BF89ADial(v012E2338ip, int(v012E233Cport))
						// ...
					}()
				}
				ebp1498 := &f004A7D34getServiceManager().s2
				ebp1499 := ebp1498.window.m0Cactive
				if ebp1499 == false {
					// f00657C13() f00670FFE() f0051B219() f0084EBF9() f00576F03() f0084B501() f0086BA70()
					// f00884C77() f0051CFAA() v0131A294.f009D8054() v0131A2A0.f00B2136D() f004DB0B1()
				}
				var ebp14A1 bool
				ebp14A0 := f0043BF3FgetT4003()
				if ebp14A0.m31 {
					ebp14A1 = false
				} else {
					// ebp14A1 = f008AEFAD(0x1B)
				}
				if ebp14A1 {
					// ebp10 := f004A7D34getServiceManager()
					// if !ebp10.f14 && !ebp10.f488C && !ebp10.f41C && !ebp10.f824 && !ebp10.f4BE4 && ebp10.fE7C && ebp10.f107C && ebp10.f9FE9 {
					// 	// f007DAFE0(0x19, 0, 0)
					// 	// ebp10.f410
					// 	ebp10.f004A9123(ebp10.f410)
					// }
				}
				if v08C88E08 != 0x14 {
					return
				}
				v01319E08log.f00B38AE4printf("> Request Character list\r\n")
				// f004E9975(0, 0, 0).f004E99D2()
				v012E2340 = 4
				v08C88E08 = 0x32
				// 0x004E2050 压缩
				var reqCharList pb // [c1 04 f3 00]
				reqCharList.f00439178init()
				reqCharList.buf[0] = 0xC1
				reqCharList.buf[2] = 0xF3
				reqCharList.buf[3] = 0
				reqCharList.len = 4
				reqCharList.buf[1] = uint8(reqCharList.len)
				reqCharList.f004393EAsend(false, false)
			}()
		case 4:
			// f004DDD4F()
		case 5:
			// f004DF0D5()
		}
		// 0x04E502F
		ebp188 := 0
		for ebp188 < 5 {
			// v0131A27C.f00534AFA(v0114EE48) // float32
			ebp188++
		}
		// f005AC5A0()
		// f005A4BC5(0x2C)
		if v08C88F88 > 0 {
			v08C88F88--
		}
		if v086A3BEC > 0 {
			v086A3BEC--
		}
		// v08C7CC18++
		// v08C7CC18 %= 32
		v0131A240++
		ebp178 -= 0x28
	}
	// 0x004E50FB
	if v01319D65 != 0 {
		return
	}
	// v09D24A20.f00514F8F()
	// f007DB28F()
	systime := struct {
		wYear         uint16
		wMonth        uint16
		wDayOfWeek    uint16
		wDay          uint16
		wHour         uint16
		wMinute       uint16
		wSecond       uint16
		wMilliseconds uint16
	}{}
	// GetLocalTime(&ebp64)
	f00DE817Asprintf(v08C88AB8[:], "Screen(%02d_%02d-%02d_%02d)-%04d.jpg", systime.wMonth, systime.wDay, systime.wHour, systime.wMinute, v08C88C74)
	// ebp220 := v08610600.f00436DF1(0x1CB)
	// var ebp410 uint32
	// if ebp220.m18 >= 0x10 {
	// 	ebp410 = ebp220.m04
	// } else {
	// 	ebp410 = &ebp220.m04 // [25 73 3A BB AD C3 E6 D2 D1 B4 A2 B4 E6 00] -> gbk "%s:画面已储存" https://r12a.github.io/app-encodings/
	// }
	// var ebp174 [100]uint8
	// f00DE817Asprintf(ebp174[:], ebp410, v08C88AB8[:])
	// ...
	// 0x004E5529
	switch v012E2340 {
	case 2:
		f004E46B3(hDC)
	case 4:
		// f004E17B9(hDC)
	case 5:
		// f004E0E03(hDC)
	}
	// ...
	ebp6C := 0x28
	v0131A270 = ebp6C
}

func f004DB77ChandleState3(hDC win.HDC) {
	// SEH
	// ebp10 := f004A7D34getServiceManager()
	if v0131A26D == false {
		// v086A3B94 = f0098967F
		v0131A26D = true
		// ...
	}
	v08C88C69 = false
	// f006BA133(0, 0, 0x280, 0x1E0)
	// glClearColor()
	// glClear()
	// f006BB37E()
	// ebp10.f9FC8.f0043DA73()
	// f006BB45B()
	// f006BA38B()
	// SwapBuffers(hDC)
	// if ebp10.f9FC8 != 0 {
	// 	// 这是什么调用？
	// 	// ebp10.9FC8()
	// 	ebp10.f9FC8 = 0
	// }
	v012E2340 = 5
	ebp14 := 0
	for ebp14 < 4 {
		// f0089DB5A(ebp14+0x798D, 0)
		ebp14++
	}
	// f005A4C09(1)
}

// 状态机
func f004E6233handleState(hDC win.HDC) {
	// SEH
	f00DE8A70chkstk() //0x4734
	// st0 = v0638C6C4 // st0~st7,浮点寄存器,80bit
	// f00DE7C90memcpy()
	ebp4738 := v012E2340
	switch ebp4738 {
	case 1:
		f004DD578handleState1(hDC)
	case 2, 4, 5:
		f004E4F1ChandleState245(hDC)
	case 3:
		f004DB77ChandleState3(hDC)
	}

	// 0x004E62F0
	// f00A49798(v012E2340).f00A4E1BF()
	func(state int) {
		// if v09D96438 == nil{
		// ...
		// }
		// return v09D96438
	}(v012E2340)

	if v01319D9C <= 0x1F {
		return
	}
	// var codef103 pb
	// codef103.f00439178init()
	// // 构造报文
	// codef103.f0043968Fxor()
	// v08C88FB4decrypt
}
