package main

import (
	"encoding/binary"
	"sync"
	"unsafe"

	"github.com/xujintao/balgass/win"
)

type t17 struct {
	m00 *treeNode
	m04 uintptr
}

func (t *t17) f0043700F(x *treeNode, y *textManager) {

}

func (t *t17) f0043702B(x *t17) bool {
	return false
}

type treeNode struct {
	p1    *treeNode
	p2    *treeNode
	p3    *treeNode
	index interface{}
	value interface{} // stdstring
	m2D   bool
}

func (t *treeNode) f00436FD5validate(x *t17) bool {
	// t.f0043702B(x)
	return false
}

func (t *treeNode) f004370B2getValue() interface{} {
	return t.value
}

type t19 struct {
	m00 *treeNode
	m04 uintptr
	m08 bool
}

func (t *t19) f004A753B(x *t17, y *bool) *t19 {
	t.m00 = x.m00
	t.m04 = x.m04
	t.m08 = *y
	return t
}

type text3 struct {
	m00index uint
	m04      *stdstring
}

func (t *text3) f004A6078text3(pindex *uint, s *stdstring) *text3 {
	// ebp4 := t
	t.m00index = *pindex
	t.m04.f004079A0stdstring(s)
	func() {

	}()
	return t
}

func (t *text3) f004A584Efree() {
	t.m04.f00407B10free()
}

// 很重要
var v08610600textManager textManager

type textManager struct {
	// 可能是个结构
	m18root *treeNode
	m1Csize int
	m20     stdstring
}

func (t *textManager) f006B8574(x uint, y uint32) bool {
	return false
}

func (t *textManager) f006B851Ddec(buf []uint8, size uint) {
	ebpCi := uint(0)
	ebp8 := buf
	ebp4key := [3]uint8{0xFC, 0xCF, 0xAB}
	for {
		if ebpCi >= size {
			break
		}
		ebp8[ebpCi] = ebp8[ebpCi] ^ ebp4key[ebpCi%3]
		// ...
		ebpCi++
	}
}

func (t *textManager) f00436FA1() **treeNode {
	// f0043712E(t.m18root)
	return func(tree *treeNode) **treeNode {
		return &tree.p2
	}(t.m18root)
}

func (t *textManager) f00436F1E(pindex *uint32, t1 *textManager) *treeNode {
	// ebpC := t
	// ebp8 := t.f00436FA1()
	ebp4 := t.m18root
	// f00437123(ebp8)
	// if *func(p *treeNode) *bool {
	// 	return &p.m2D
	// }(ebp8) == true {
	// 	return ebp4
	// }
	return ebp4
}

func (t *textManager) f00436FB6(x *t17) *t17 {
	return x
}

func (t *textManager) f004A793A(x *t17) *t17 {
	return x
}

func (t *textManager) f004A6371(list1 *t17, num int, list2 *treeNode, v *text3) *t17 {
	return nil
}

func (t *textManager) f00436E48(x *t17, pindex *uint) *t17 {
	// 0x24局部变量
	// ebp20 := t

	var ebp8 t17 // 0x0ED26630, 0x0ED1FF00
	// t.f00436EEA(&ebp8, pindex)
	func(x *t17, pindex *uint) *t17 {
		// x.f0043700F(t.f00436F1E(pindex, t)) // f0043700F有两个参数
		return x
	}(&ebp8, pindex)

	var ebp14 t17 // 0x0ED26630, 0x0ED1FF00
	ebp8.f0043702B(t.f00436FB6(&ebp14))

	var ebp1C t17 // 0x0ED26630, 0x0ED1FF00
	ebp24 := t.f00436FB6(&ebp1C)
	ebpC := ebp24
	x.m00 = ebpC.m00
	x.m04 = ebpC.m04
	return x
}

func (t *textManager) f00436DF1findstdstring(index uint) *stdstring {
	// 1C局部变量
	// ebp1C := t
	var ebp10 t17
	p := t.f00436E48(&ebp10, &index)

	var ebp18 t17
	if p.m00.f00436FD5validate((t.f00436FB6(&ebp18))) {
		v := p.m00.f004370B2getValue().(stdstring)
		return &v
	}
	return &t.m20
}

func (t *textManager) f00436DA8findcstr(index uint) []uint8 {
	// t.f00436DBE(index)
	return func(index uint) []uint8 {
		// t.f00436DD4(index)
		return func(index uint) []uint8 {
			return t.f00436DF1findstdstring(index).f004073E0cstr()
		}(index)
	}(index)
}

func (t *textManager) f006B8D79assign(index uint, buf []uint8) bool {
	// 0x64局部变量
	// ebp68 := t
	var ebp14 t17
	t.f00436E48(&ebp14, &index)

	var ebp1C t17 // 0x0ED26630, 0x0ED1FF00
	if ebp14.f0043702B(t.f00436FB6(&ebp1C)) == false {
		return false
	}

	var ebp58 stdstring
	ebp58.f00406FC0stdstring(buf)
	ebp4 := 0
	var ebp3C text3
	ebp6C := ebp3C.f004A6078text3(&index, &ebp58) // 构造text3对象
	ebp70 := ebp6C
	ebp4 = 1
	var ebp64 t19
	// t.f004A60D1(&ebp64, ebp70) // append text3对象
	func(x *t19, v *text3) *t19 {
		// 4C局部变量
		// ebp48 := t
		ebpC := *t.f00436FA1()
		ebp8 := t.m18root
		ebp1 := true
		for {
			// f00437123
			if ebpC.m2D {
				break
			}
		}
		var ebp14 t17
		ebp14.f0043700F(ebp8, t)
		if ebp1 {
			var ebp28 t17
			if ebp14.f0043702B(t.f004A793A(&ebp28)) {
				ebp29 := true
				var ebp34 t17
				x.f004A753B(t.f004A6371(&ebp34, 1, ebp8, v), &ebp29)
				return x
			}
		}
		return nil
	}(&ebp64, ebp70)
	ebp4 = 0
	ebp3C.f004A584Efree()
	ebp4 = -1
	ebp58.f00407B10free()
	println(ebp4)
	return true
}

func (t *textManager) f006B83FD(fileName *stdstring, x uint32) {
	// 0x28局部变量
	// ebp28 := t
	ebp4file := f00DE909Efopen(string(fileName.f004073E0cstr()), "rb") // v012F7910 为什么返回一个全局变量的地址？
	if ebp4file == nil {
		return
	}
	var ebpCheader [8]uint8
	f00DE8FBDfread(ebpCheader[:], 6, 1, ebp4file)
	if binary.LittleEndian.Uint16(ebpCheader[:]) != 0x5447 { // leadcode: "GT"
		return
	}
	ebpAsize := binary.LittleEndian.Uint32(ebpCheader[2:]) // 0x10C5
	ebp10index := uint32(0)
	var ebp1C [8]uint8
	for {
		if ebp10index >= ebpAsize {
			break
		}
		f00DE8FBDfread(ebp1C[:], 8, 1, ebp4file)
		record := struct {
			index uint    // ebp1C
			size  uint    // ebp18
			buf   []uint8 // ebp14
		}{}
		record.index = uint(binary.LittleEndian.Uint32(ebp1C[:])) // 0
		record.size = uint(binary.LittleEndian.Uint32(ebp1C[4:])) // 4
		record.buf = f00DE64BCnew(uint(record.size + 1))
		f00DE8FBDfread(record.buf, 1, uint(record.size), ebp4file) // 37 01 67 19
		if t.f006B8574(record.index, x) == true || record.index < 0x270F {
			t.f006B851Ddec(record.buf, record.size)
			record.buf[record.size] = 0
			t.f006B8D79assign(record.index, record.buf) // store
		}
		f00DE7BEAdelete(record.buf)
		ebp10index++
	}
	f00DE8C84close(ebp4file)
}

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
	do5(bool)
	do6(bool)
	do10() bool
	do11(float64)
	do12()
	do13(float64)
}

type button struct {
	mA9done bool
	mAA     bool
	mAB     bool
}

func (s *button) f00448976() bool { return false }
func (s *button) f00436088() bool {
	if s.mA9done {
		return s.f00448976()
	}
	return false
}
func (s *button) f004CCE44(x bool)  { s.mA9done = x }
func (s *button) f004360AB(float64) {} // 比较复杂
func (s *button) f0043912C() bool   { return s.mAA }
func (s *button) f00438F7B(x bool)  { s.mAB = x }
func (s *button) f00436060(x bool) {
	// s.f0043672D(x)
	if x == false {
		s.mAA = false
	}
}

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
	m0Cshow   bool
	m0D       bool
	m0E       bool
	m10       int
	m14       *int
	m18left   int // rect left
	m1Ctop    int // rect top
	m20width  int // rect width
	m24heigth int // rect height
	m28left   int // = m18left
	m2Ctop    int // = m1Ctop
	subs      list
}

func (b *window) f004CCA35(x int) bool {
	if b.m0Cshow == false {
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
			if ebp1C.(*button).f00436088() == true {
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
	if b.m0Cshow == false {
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
		ebp8.(*button).f004CCE44(x)
	}
}

func (b *window) f004CCC07(unk float64) {
	if b.m0Cshow == false {
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
			ebpC.(*button).f004360AB(unk)
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
func (b *window) f004AA018isActive() bool { return b.m0Cshow }
func (b *window) f004AA068() int          { return b.m10 }
func (b *window) f004AA027() bool         { return b.m0D }

// service1 v0130F730 v01317CD8 v013180E0
// exit
type service1 struct {
	window
	m3D4 uint
	m3D8 uint
}

// do4->f0043F608
func (s *service1) do4(x int) bool   { return s.window.do4(x) }
func (s *service1) do5(x bool)       {}
func (s *service1) do6(bool)         {}
func (s *service1) do10() bool       { return false }
func (s *service1) do11(unk float64) {}
func (s *service1) do12()            {}

// do13->f0043F640 may exit
func (s *service1) do13(unk float64) {}

func (s *service1) f0043FBCC(x, y uint) { // 渲染文本字符串
	// A4局部变量
	// ebpA0 := s
	// ebpC := f004A7D34getWindowManager()
	// ebp8 := 0
	// ebp94 := 0
	// ebp4 := 2
	s.m3D4 = x
	s.m3D8 = 3
	switch x {
	case 0x36:
		// 0x0044034D
		// ebpC.f004A9123(&ebpC.m92D8) // ebpC.f004A9123(&ebpC.s???)
		// ebpC.m92D8.f0043752C(&ebpC.m9DC8)
		v08610600textManager.f00436DA8findcstr(0x6B4) // 1716
	}

}

// service2 v01314300
type service2 struct {
	window
}

// do4->f0043AA4E
func (s *service2) do4(x int) bool   { return s.window.do4(x) }
func (s *service2) do5(x bool)       {}
func (s *service2) do6(bool)         {}
func (s *service2) do10() bool       { return false }
func (s *service2) do11(unk float64) {}
func (s *service2) do12()            {}
func (s *service2) do13(unk float64) {}

// service3 v01313FA8
// login
type service3 struct {
	window
	subs [2]button // offset: 0E0, v01313FE8
	// offset: 120, v013140C8
	// m200username text
	// m204password text
	// m208 [2]struct{ data [0xA8]uint8 }
}

// do4->f0043E5D4
func (s *service3) do4(x int) bool { return s.window.do4(x) }

// do5->f0043E575
func (s *service3) do5(x bool) {
	// s.f004CCBCA(x)
	func(x bool) {
		if s.m14 != nil {
			// s.m14.f0043672D(x)
		}
		s.m0Cshow = x
		if s.m0Cshow == false {
			s.m0D = false
		}
	}(x)

	ebp4 := 0
	for ebp4 < 2 {
		// s.m208[ebp4].f0043672D(x)
		s.subs[ebp4].f00436060(x)
	}
}

// do6->f00435D53
func (s *service3) do6(x bool)       { s.m0D = x }
func (s *service3) do10() bool       { return false }
func (s *service3) do11(unk float64) {}
func (s *service3) do12()            {}

// 关闭当前连接并重新拨号
func (s *service3) f0043ED78() {
	// s.f0043ED98()
	func() {
		v08C88FF0conn.f006BD5BBclose()
		v08C88F74 = 0
		v08C88E08 = 0
		f006BF89ADial(v012E2338ip, int(v012E233Cport))
	}()
	f004A7D34getWindowManager().f004A9146LRU(s)
}

// s9: 006E855E
func (s *service3) f0043E9B6() {
	// seh
	f00DE8A70chkstk() // 0x152C
	// ebp1538 := s
	if v08C88E08 == 0 {
		return
	}
	f004A7D34getWindowManager().f004A9146LRU(s)

	var ebp18, ebp24 [11]uint8 // username, pwd
	f00DE8100memset(ebp18[:], 0, 11)
	// s.m200username.f50h(ebp18[:], 11) // f00342929, fill username
	f00DE8100memset(ebp24[:], 0, 11)
	// s.m204password.f50h(ebp24[:], 11) // f00452929, fill pwd

	// f0043EE13(ebp18[:]) // 判断字符串长度
	// f0043EE13(ebp24[:])
	if v08C88E08 != 2 {
		return
	}

	v01319E08log.f00B38AE4printf("> Login Request.\r\n")
	v01319E08log.f00B38AE4printf("> Try to Login \"%s\"\r\n", ebp18[:])

	v08C88F74 = 1
	f00DE8000strcpy(v08C88F78username[:], ebp18[:])
	v08C88E08 = 13

	// 构造登录报文
	var ebp14BClogin pb
	ebp14BClogin.f00439178init()
	ebp14BClogin.f0043922CwriteHead(0xC1, 0xF1) // 前缀
	ebp14BClogin.f004397B1writeUint8(1)         // 可能是subcode

	var ebp30, ebp3C [11]uint8
	f00DE8100memset(ebp30[:], 0, 11)
	f00DE8100memset(ebp3C[:], 0, 11)
	f00DE7C90memcpy(ebp30[:], ebp18[:], 10)
	f00DE7C90memcpy(ebp3C[:], ebp24[:], 10)
	f0043B750xor(ebp30[:], 10) // 与igc.dll的TOOLTIP_FIX_XOR_BUFF一样了，可能有问题
	f0043B750xor(ebp3C[:], 10) // 与igc.dll的TOOLTIP_FIX_XOR_BUFF一样了，可能有问题

	ebp14BClogin.f00439298writeBuf(ebp30[:], 10, true)    // 写username
	ebp14BClogin.f00439298writeBuf(ebp3C[:], 10, true)    // 写pwd
	ebp14BClogin.f0043EDF5writeUint32(win.GetTickCount()) // 写时间戳
	ebp14C0 := 0
	for ebp14C0 < 5 {
		ebp14BClogin.f004397B1writeUint8(v012E4018versionDLL[ebp14C0] - byte(ebp14C0+1)) // 写版本 VERSION_HOOK1
		ebp14C0++
	}
	ebp14C0 = 0
	for ebp14C0 < 16 {
		ebp14BClogin.f004397B1writeUint8(v012E4020serial[ebp14C0]) // 写序列号 SERIAL_HOOK1
		ebp14C0++
	}

	// 发送登录报文
	ebp14BClogin.f004393EAsend(true, false)

	// var ebp4 int
	// ebp14E0.f00406FC0stdstring(v08610600textManager.f00436DA8(0x1D8))
	// ebp4 = 1
	// ebp14FC.f00406FC0stdstring(&v0114A327)
	// ebp4 = 2
	// f00A49798ui().f0043EE21().f00A9FB38print(&ebp14FC, &ebp14E0, 3, 0)
	// ebp4 = 1
	// ebp14FC.f00407B10()
	// ebp4 = 0
	// ebp14E0.f00407B10()

	// ebp1518.f00406FC0stdstring(v08610600textManager.f00436DA8(0x1D9))
	// ebp4 = 3
	// ebp1534.f00406FC0stdstring(&v0114A33E)
	// f00A49798ui().f0043EE21().f00A9FB38print(&ebp1534, &ebp1518, 3, 0)
	// ebp4 = 3
	// ebp1534.f00407B10()
	// ebp4 = 0
	// ebp1518.f00407B10()
	// ebp4 = -1

	ebp14BClogin.f004391CF()
}

// do13->f0043E60C 发送login报文
func (s *service3) do13(unk float64) {
	if s.subs[0].f0043912C() {
		s.f0043E9B6() // login
		return
	}
	if s.subs[1].f0043912C() {
		s.f0043ED78() // 关闭当前连接并重新拨号
		return
	}
	if f0043BF3FgetT4003().f0043913E(13) {
		// f007DAFE0(25, 0, 0)
		s.f0043E9B6() // login
		return
	}
	if f0043BF3FgetT4003().f0043913E(27) {
		// f007DAFE0(25, 0, 0)
		s.f0043ED78() // 关闭当前连接并重新拨号
		f004A7D34getWindowManager().f00439161(false)
	}
}

// service4 v01310798
// serverList and serverInfo
type service4 struct {
	window
	serverConns [21]button
	serverGames [40]button // offset: 0x12A0
	m3758       int
	m3760       int
	m3764       *int
}

// do4->f00446D35
func (s *service4) do4(x int) bool   { return s.window.do4(x) }
func (s *service4) do5(x bool)       {}
func (s *service4) do6(bool)         {}
func (s *service4) do10() bool       { return false }
func (s *service4) do11(unk float64) {}
func (s *service4) do12()            {}

// do13->f00446D6D 请求服务器列表，请求服务器信息
func (s *service4) do13(unk float64) {
	// 带SEH处理
	f00DE8A70chkstk() // 0x2994

	// ... 浮点运算
	ebp10 := 0
	for ebp10 < 0x15 {
		if s.serverConns[ebp10].f0043912C() {
			if s.m3760 != 0 {
				s.serverConns[s.m3760].f00438F7B(false)
			}
			s.serverConns[ebp10].f00438F7B(true)
			s.m3760 = ebp10

			// 请求服务器列表
			var ebp149C pb
			ebp149C.f00439178init()
			ebp149C.f0043922CwriteHead(0xC1, 0xF4)
			var ebp15 [1]uint8
			ebp15[0] = 6
			ebp149C.f00439298writeBuf(ebp15[:], 1, false)
			ebp149C.f004393EAsend(false, false) // 发送协议报文
			ebp149C.f004391CF()
		}
		ebp10++
	}
	if s.m3764 == nil {
		return
	}

	ebp10 = 0
	for ebp10 < s.m3758 {
		if s.serverGames[ebp10].f0043912C() {
			// ebp14 := s.m3764.f00AF712E(ebp10)
			// if ebp14 == nil {
			// 	break
			// }
			// if ebp14.m10 < 0x64 {
			{
				f004A7D34getWindowManager().f004A9146LRU(s)
				// if f004472C9(&ebp14.m15, v0114A5F8) {
				// 	f0043A2DF(true).f004472DB()
				// } else {
				// 	f0043A2DF(false).f004472DB()
				// }
				// 请求服务器信息
				var ebp292CserverInfo pb
				ebp292CserverInfo.f00439178init()
				ebp292CserverInfo.f0043922CwriteHead(0xC1, 0xF4) // header
				var ebp14A9 [1]uint8
				ebp14A9[0] = 3
				ebp292CserverInfo.f00439298writeBuf(ebp14A9[:], 1, false) // subcode
				var ebp14A8code uint16                                    // = ebp14.m0C
				var codes [2]uint8
				binary.BigEndian.PutUint16(codes[:], ebp14A8code)
				ebp292CserverInfo.f00439298writeBuf(codes[:], 2, false) // server code
				ebp292CserverInfo.f004393EAsend(false, false)

				// var ebp4 int
				// ebp2948.f00406FC0stdstring(v08610600textManager.f00436DA8(0x1D6))
				// ebp4 = 2
				// ebp2964.f00406FC0stdstring(&v0114A5FF)
				// ebp4 = 3
				// f00A49798ui().f0043EE21().f00A9FB38print(&ebp2964, &ebp2948, 3, 0)
				// ebp4 = 2
				// ebp2964.f00407B10()
				// ebp4 = 1
				// ebp2948.f00407B10()
				// ebp2980.f00406FC0stdstring(v08610600textManager.f00436DA8(0x1D7))
				// ebp4 = 4
				// ebp299C.f00406FC0stdstring(&v0114A600)
				// ebp4 = 5
				// f00A49798ui().f0043EE21().f00A9FB38print(&ebp299C, &ebp2980, 3, 0)
				// ebp4 = 4
				// ebp299C.f00407B10()
				// ebp4 = 1
				// ebp2980.f00407B10()
				// ebp4 = -1
				// ebp292C.f004391CF()

				// ebp14A0 := 1
				// if s.m3764.m1C {
				// 	ebp14A0 = 3
				// } else if ebp14.m14&1 != 0 {
				// 	ebp14A0 = 2
				// }

				// ebp14A1 := 0
				// if s.m3764.m04 == false {
				// 	ebp14A1 = 1
				// }
				// f00AF7DC3(&s.m3764.m31, ebp14.m08, ebp14A0, ebp14.m14, ebp14A1).f00AF876E()
				break
			}
			// if ebp14.m10 < 0x80 {
			// 	f004A7D34getWindowManager().f004A9EEB(0x63, 0)
			// }
		}
		ebp10++
	}
}

// service5 v01310598
type service5 struct {
	window
	s1 button // offset: 0x040, v013105D8
	s2 button // offset: 0x120, v013106B8
}

// do4->f0043DCD3
func (s *service5) do4(x int) bool { return s.window.do4(x) }
func (s *service5) do5(x bool)     {}

// do6->f00435D53
func (s *service5) do6(x bool) {
	s.m0D = x
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
		ebp10 := f004A7D34getWindowManager()
		ebp10.f004A9123(&ebp10.s7)
		ebp10.f00439161(true)
	}
	if s.s2.f0043912C() {
		// 请求服务器列表(没有执行)
		var ebp149C pb
		ebp149C.f00439178init()
		// ebp4 := 0
		ebp149C.f0043922CwriteHead(0xC1, 0xF4)
		var ebp15 [1]uint8
		ebp15[0] = 6
		ebp149C.f00439298writeBuf(ebp15[:], 1, false)
		ebp149C.f004393EAsend(false, false)
		ebp149C.f004391CF()
		ebp14 := f004A7D34getWindowManager()
		ebp14.f004A9123(&ebp14.s2)
		// f004D4F77(v012E239C, 0) // main_theme.mp3
		// f004D4FB9(v012E234C, 0) // mutheme.mp3
	}
}

// service6 v0130FF40
type service6 struct {
	window
}

// do4->f00441096
func (s *service6) do4(x int) bool   { return s.window.do4(x) }
func (s *service6) do5(x bool)       {}
func (s *service6) do6(bool)         {}
func (s *service6) do10() bool       { return false }
func (s *service6) do11(unk float64) {}
func (s *service6) do12()            {}
func (s *service6) do13(unk float64) {}

// service7 v0130FB38
type service7 struct {
	window
}

// do4->f00449300
func (s *service7) do4(x int) bool   { return s.window.do4(x) }
func (s *service7) do5(x bool)       {}
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

func (t *list) f004455FBdeleteNode(n *node) bool {
	ebp4 := n
	// if IsBadReadPtr(ebp4, 12){
	// 	return false
	// }

	if t.head == ebp4 {
		t.head = ebp4.next
	} else {
		// if IsBadReadPtr(ebp4.prev, 12) {
		// 	return false
		// }
		ebp4.prev.next = ebp4.next
	}

	if t.tail == ebp4 {
		t.tail = ebp4.prev
	} else {
		// if IsBadReadPtr(ebp4.next, 12) {
		// 	return false
		// }
		ebp4.next.prev = ebp4.prev
	}
	// t.f004458A1(ebp4)
	func(n *node) {
		// f00DE7538(n)
		t.num--
		if t.num == 0 {
			// t.f004454EA()
			func() {
				for t.head != nil {
					// ebp4 := t.head
					t.head = t.head.next
					// f00DE7538(ebp4)
				}
				t.tail = nil
				t.num = 0
			}()
		}
	}(ebp4)
	return true
}

func (t *list) f004457B6findNode(val interface{}, n *node) *node {
	ebp4 := n
	if ebp4 == nil {
		ebp4 = t.head
	} else {
		// if IsBadReadPtr(ebp4, 12) {
		// 	return nil
		// }
		ebp4 = ebp4.next
	}
	for ebp4 != nil {
		if ebp4.value == val {
			return ebp4
		}
		ebp4 = ebp4.next
	}
	return nil
}
func (t *list) f00445416appendNode(val interface{}) *node {
	// ebp4 := t.f0044585BnewNode(t.tail, nil)
	ebp4 := func(n1, n2 *node) *node {
		ebp8 := (*node)(f00DE852Fnew(12))
		ebp4 := ebp8
		ebp4.prev = n1
		ebp4.next = n2
		t.num++
		ebp4.value = nil
		return ebp4
	}(t.tail, nil)
	ebp4.value = val
	if t.tail != nil {
		t.tail.next = ebp4
	} else {
		t.head = ebp4
	}
	t.tail = ebp4
	return ebp4
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

// serviceManager
func f004A7D34getWindowManager() *serviceManager {
	v01319730once.Do(func() {
		v0130F728.f004A7A82()
		// f00DE8BF6atexit(f0114817A)
	})
	return &v0130F728
}

var v0130F728 serviceManager
var v01319730once sync.Once

type serviceManager struct {
	s001 service1 // v0130F730
	// v0130F8C0 v0130F9A0

	s7 service7 // v0130FB38
	// v0130FBC0 v0130FCA0 v0130FD80 v0130FE60

	s6 service6 // v0130FF40
	// v0130FFC8 v01310268

	s5 service5 // offset:0xE70, v01310598
	// v013105D8
	// v013106B8

	s4 service4 // offset:0x1070, v01310798
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

	// m92D8 v01318A00
	m9DC8 uint8

	// ... // 100个service
	f9FD4activeServices list // v013196FC
	m9FE0               bool // v01319708
	m9FE1               bool
	m9FE4               uint // v0131970C
	m9FE8               bool // v01319710
	m9FE9               bool
	m9FEC               *stdstring // v01319714
}

func (t *serviceManager) f004A89FFconstruct() {
	t.m9FE0 = false
	t.m9FE1 = false
	t.m9FE4 = 0
	t.m9FE8 = false
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
		if t.f9FD4activeServices.f004455FBdeleteNode(t.f9FD4activeServices.f004457B6findNode(s, nil)) == false {
			return nil
		}
		t.m9FE8 = true
		// t.f9FD4activeServices.f004453C6(s)
	}
	return ebp4
}
func (t *serviceManager) f004A9123(s interface{}) {
	// s.(servicer).do5(true)
	t.f004A9083(s)
}
func (t *serviceManager) f004A9146LRU(s interface{}) {
	if t.f9FD4activeServices.f004455FBdeleteNode(t.f9FD4activeServices.f004457B6findNode(s, nil)) == false {
		return
	}
	s.(servicer).do5(false)
	s.(servicer).do6(false)
	t.f9FD4activeServices.f00445416appendNode(s)
	s = t.f9FD4activeServices.f004452A7getFirstNodeValue()
	if s.(windower).f004AA018isActive() {
		t.m9FE8 = true
	}
}
func (t *serviceManager) f00439161(x bool) { t.m9FE9 = x }
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
	if ebp30.m9FE4 == 0 { // 0, 2
		return
	}
	if ebp30.f9FD4activeServices.f004AA077isNumZero() == true {
		return
	}
	if ebp30.m9FE8 { // 0, 1
		ebp18 := ebp30.f9FD4activeServices.f004452A7getFirstNodeValue()
		if ebp18.(windower).f004AA018isActive() {
			ebp18.(servicer).do6(true)
			ebp30.m9FE8 = false
		}
	}

	ebp10 := f0043BF3FgetT4003() // v01308D18
	if ebp10.f0043913E(0x1B) {   // always be false
		// ...
	}
	ebp30.m9FE0 = false
	if ebp10.f00436696() {
		// ...
	}
	if ebp10.f004366A5() {
		ebp30.m9FE1 = false
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
		// f00DE7BEAdelete(ebp8)
		ebp8 = ebp8[:0]
	}
	ebp30.f004A91CE()
	ebp4 = ebp30.f9FD4activeServices.f004409AAgetList()
	for ebp4 != nil {
		ebpC := ebp30.f9FD4activeServices.f00445530getNodeValue(&ebp4)
		ebp34 := ebpC.(windower).f004AA068()
		switch ebp34 {
		case 1, 2, 3, 4:
			ebp30.m9FE0 = true
		}

		if ebp30.m9FE0 {
			break
		}
		if ebpC.(servicer).do4(0) == true {
			ebp30.m9FE0 = true
			break
		}
	}
}

func (t *serviceManager) f004A9EEB(x, y uint) {
	switch t.m9FE4 {
	case 0, 1, 3, 5:
		return
	default:
		t.s001.f0043FBCC(x, y)
	}
}

func (t *serviceManager) f004A9F3B(buf []uint8) {
	if len(buf) == 0 {
		return
	}
	if t.m9FE4 == 4 {
		// t.f9DD8.f00445A2A(buf)
	} else {
		t.m9FEC.f0043D7E2stdstring(string(buf))
	}
}

func f0043BF3FgetT4003() *t4003 {
	v01308D80once.Do(func() {
		v01308D18.f0043BF18init()
		// f00DE8BF6atexit(f01148111)
	})
	return &v01308D18
}

var v01308D18 t4003
var v01308D80once sync.Once

type t4003 struct {
	m08hWnd            win.HWND
	m0Cpoint           point
	m1C                bool
	m1D                bool
	m1E                bool
	m21                bool
	m25                bool
	m28width           int
	m2Cheight          int
	m30                bool
	m31                bool
	m34pointClient     point
	m40doubleClickTime float64
	m48                float64
	m50                float64
	m58                float64
	m60                bool
	m61                bool
	m62                bool
}

func (t *t4003) f0043BF18init() {
	// 虚表初始化？
}

// f00439105
func (t *t4003) do1() {

}
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
	*p = t.m0Cpoint
	return p
}
func (t *t4003) f0043BF9C(hWnd win.HWND, width, height int) bool {
	// ebp4 := t
	if hWnd == 0 {
		return false
	}
	t.m08hWnd = hWnd
	t.m28width = width
	t.m2Cheight = height
	// dll.user32.GetCursorPos(&t.m0Cpoint) // Retrieves the position of the mouse cursor, in screen coordinates.
	// dll.user32.ScreenToClient(hWnd, &t.m0Cpoint) // converts the screen coordinates of a specified point on the screen to client-area coordinates.
	t.m34pointClient = t.m0Cpoint
	t.m30 = false
	var ebpC uint // ebpC := dll.user32.GetDoubleClickTime()
	t.m40doubleClickTime = float64(ebpC)
	return true
}
func (t *t4003) f00448D8BgetX() int { return t.m0Cpoint.x }
func (t *t4003) f00448191getY() int { return t.m0Cpoint.y }

func f004DD578handleState1(hDC win.HDC) {
	// if ebp198.f00407AC0() == true { // 0x004DD639 gameguard 3, hook hook always false, disable GameGuard
	if false {
		v01319E08log.f00B38AE4printf("> ResourceGuard Error!!\r\n")
		win.SendMessage(v01319D6ChWnd, 2, 0, 0)
	}
	// load Interface/...
	// ...
	// 0x004DD963
	// f006B4C78(hdC)
	func(hDC win.HDC) {
		// load NPC/...
		// ...
		// 0x006B7AAE
		// ebp170.f004DBD5C()
		// ebp64.f00B60770(&ebp170)
		var ebp164 [200]uint8
		f00DE817Asprintf(ebp164[:], "data/local/Gameguard.csr")
		var ebp190tmpstr stdstring
		ebp190tmpstr.f00406FC0stdstring(ebp164[:])
		// var ebp171 bool
		// if ebp64.f00B607F0(&ebp190tmpstr) == 0 {
		// 	ebp171 = true
		// } else {
		// 	ebp171 = false
		// }
		ebp190tmpstr.f00407B10free()
		// if ebp171 == true { // 0x006B7B2A gameguard 4, hook always false, disable GameGuard
		if false {
			v01319E08log.f00B38AE4printf("> ResourceGuard Error!!\r\n")
			win.SendMessage(v01319D6ChWnd, 2, 0, 0)
		}
		// ...
		// 0x006B7C18
		var ebp1B0tmpstr stdstring
		ebp1B0tmpstr.f00406FC0stdstring(ebp164[:])
		// var ebp191 bool
		// if ebp64.f00B607F0(&ebp190tmpstr) == 0 {
		// 	ebp191 = true
		// } else {
		// 	ebp191 = false
		// }
		ebp1B0tmpstr.f00407B10free()
		// if ebp191 == true { // 0x006B7C63 gameguard 5, hook always false, disable GameGuard
		if false {
			v01319E08log.f00B38AE4printf("> ResourceGuard Error!!\r\n")
			win.SendMessage(v01319D6ChWnd, 2, 0, 0)
		}
		// 0x006B7C88
		f00A49798ui().f00A4A521loadResource()
		// ...
		// f006B7DDC() // load Data/local/...
		func() {
			// 0xA0局部变量
			ebp4 := -1
			var ebp74fileName [100]uint8
			f006B8367iocopy := func(fileName string, buf []uint8, b bool) bool {
				// 0x3C局部变量
				ebp4 := -1
				if fileName == "" || buf == nil {
					return false
				}
				var ebp44 stdstring
				ebp44.f00406A20init()
				ebp4 = 0
				var ebp28 stdstring
				ebp28.f00406A20init()
				ebp4 = 1
				ebp44.f0043D7E2stdstring(fileName)
				f00DE8000strcpy(buf, ebp44.f004073E0cstr())
				ebp45ret := true
				ebp4 = 0
				ebp28.f00407B10free()
				ebp4 = -1
				ebp44.f00407B10free()
				println(ebp4)
				return ebp45ret
			}
			// "Data/Local/Text.bmd"
			f006B8367iocopy("Data/Local/Text.bmd", ebp74fileName[:], true)
			var ebp90 stdstring
			ebp90.f00406FC0stdstring(ebp74fileName[:])
			ebp4 = 0
			v08610600textManager.f006B83FD(&ebp90, 0x80000010)
			ebp4 = -1
			ebp90.f00407B10free()
			// "Data/Macro.txt"
			// ...
			// 0x006B81F8
			f00AF7DC3serverList().f00AF7F07load() // v09D965B0.f00AF7F07()

			println(ebp4)
		}()
		// ...
	}(hDC)
	// ...
	// 0x004DDA33
	v01319E08log.f00B38AE4printf("> Loading ok.\r\n")
	// f004DAACA(v01319D6ChWnd)
	v012E2340state = 2
	v0131A250 = win.GetTickCount()
}

func f004E1E1EhandleState2() {
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
	ebp1498 := &f004A7D34getWindowManager().s2
	ebp1499 := ebp1498.m0Cshow
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
		// ebp10 := f004A7D34getWindowManager()
		// if !ebp10.f14 && !ebp10.f488C && !ebp10.f41C && !ebp10.f824 && !ebp10.f4BE4 && ebp10.fE7C && ebp10.f107C && ebp10.m9FE9 {
		// 	// f007DAFE0(0x19, 0, 0)
		// 	// ebp10.f410
		// 	ebp10.f004A9123(ebp10.f410)
		// }
	}
	if v08C88E08 != 20 {
		return
	}
	v01319E08log.f00B38AE4printf("> Request Character list\r\n")
	// f004E9975(0, 0, 0).f004E99D2()
	v012E2340state = 4
	v08C88E08 = 50
	// 0x004E2050 压缩
	var reqCharList pb // [c1 04 f3 00]
	reqCharList.f00439178init()
	reqCharList.buf[0] = 0xC1
	reqCharList.buf[2] = 0xF3
	reqCharList.buf[3] = 0
	reqCharList.len = 4
	reqCharList.buf[1] = uint8(reqCharList.len)
	reqCharList.f004393EAsend(false, false)
}

// s9 f004DEDAD
func f004E46B3handleState2(hDC win.HDC) bool {
	if v0131A26C == false {
		return false
	}
	// ...
	// 1.04R location: 0x004E4819
	// s9    location: 0x004DEED2
	// f0051B429
	func() {
		// 0x0A8FD064 hook to hide f0AD98B31, disable anti-temper with backup code
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
	return true
}
func f004DDD4FhandleState4() {}

func f004E17B9handleState4(hDC win.HDC) bool {
	// 0x0ADB3664 hook to hide f09EBC65D, disable anti-temper with backup code
	var label1 uint32 = 0x0091AD10
	// push label1
	// push 0x0A4E45B0
	// ret

	// 0x0A4E45B0
	var label2 uint32 = 0x00BFE32C
	// push label2
	// push 0x0A4E0636
	// ret

	// 0x0A4E0636
	var label3 uint32 = 0x00E26B30
	// push label3
	// push 0x09FE242C
	// ret

	// 0x09FE242C 0x0A84C76E 0x09E72845 0x09E2466D 0x0A84D026
	// push esi
	// push ecx
	// push edi
	// push ebx
	// pushfd
	// push edx

	// 0x0A8FCAAA
	// push 0x0AF7DCB3
	// push 0x09F94AB8
	// ret

	// 0x09F94AB8
	func() {
		// push ebx
		// push esi
		// push edi
		var ecx uint32
		if v0A443F74 == v0AD3B896 {
			// 0x0B10C7F0
			ecx = v0A84A81D
		} else {
			// 0x0AF98AAC
			ecx = v012F7B94 // *(v09F8DB1C ^ v0A437DC7 - v0AFDC2F2 + &v09F94AB8)为v012F7B94
			v0AD3B896 = v0A443F74
		}

		// 0x0A83DF81 0x0AF8A2A5 0x0AF7E149 0x0A8FB91A 0x0A9FE3BC 0x09E2B7D1 0x0A5D4DD2
		if ecx == 0 {
			// 0x0AAB6395
			ecx++
		}
		// 0x0A846783 0x0A88F548
		v0A84A81D = v0A890E43 * ecx % v0A56EB4C
		// ecx = v0ABE324F
		if v0A84A81D <= v09FE1805 {
			// 0x09E2E91B
			ebp8 := v0A88819B
			ebp4 := ebp8 & 0x1F // 0x1A
			ebpC := v0AF96824backupCode[:]
			edi := v0A9F69B2blocks[:]
			label1 = v0ABE324Flabel1 // [ebp+v0A952A97*4+8] = v0ABE324F
			for edi[0].addr != ^uintptr(0) {
				// 0x0AD7AFCB
				addr := edi[0].addr
				if addr == ^uintptr(1) {
					// 0x09FB7DFD
					addr = uintptr(edi[0].size) + v09FE37F3imageBase
					break // 0x0AF8264D 0x0AF7FFE8
				}
				// 0x0A83CFB2 0x0AFD65B3 0x0A193940 0x0AA2DC0B 0x0ABE38DC 0x0A32E2E1
				edx := addr + v09FE37F3imageBase
				eax := edi[0].size
				if eax > 3 { // 商
					// 0x0A057F34 0x09E2E3C1 0x0ABD8DA1
					esi := eax / 4
					eax %= 4
					// num = size%4
					for {
						// 0x0AFDFEC4
						ebp8 <<= ebp4 // 循环左移在c语言里如何表示？
						ebx := ebpC[0] - ebp4*ebp8
						ebp4 ^= ^ebx
						*(*uint32)(unsafe.Pointer(edx)) = ebx // 使用备份代码覆盖当前代码防篡改
						ebpC = ebpC[1:]
						edx += 4
						esi--
						if esi == 0 {
							break // 0x0A04C6C1
						}
					} // for loop 0x0AFDFEC4
				}
				// 0x0A04C6C1
				if eax > 0 { // 余数
					// 0x0AD92EFD
					for {
						// 0x09EABF12 0x0AFE088E 0x09FE279B 0x0A84E240 0x0B07017C 0x0A5FD6D8 0x0A4DFCAE
						esi := eax
						ebp8 <<= ebp4
						ebp8 *= ebp4
						cl := uint8(ebpC[0] - ebp8)
						// ebpC快进一个字节
						ebp4 ^= ^uint32(cl)
						*(*uint8)(unsafe.Pointer(edx)) = cl
						edx++
						esi--
						if esi == 0 {
							break // 0x0A936780
						}
					} // for loop 0x09EABF12
				}

				// 0x0A936780
				edi = edi[1:]
			} // for loop 0x0AD7AFCB

			// 0x0AF7FFE8 0x0A916F27
			label2 = v0A92FC07label2 // [ebp+v0A902E49*4+8] = v0A92FC07
			v0AA30425 = v0AFD3C89
			// 0x0A4DFFBE
		} else {
			// 0x0A4375F9
			label1 = v0ABE324Flabel1 // [ebp+v0A952A97*4+8] = v0ABE324F
			label2 = v0A92FC07label2 // [ebp+v0A902E49*4+8] = v0A92FC07
			// 0x0A4DFFBE
		}
		// 0x0A4DFFBE 0x0A903547 0x09E6F742 0x0A05A39D
		// pop edi
		// pop esi
		label3 = v0A32B0C2label3 // [ebp+v0A6039A7*4+8] = v0A32B0C2
		// pop ebx
	}()

	// 0x0AF7DCB3
	// pop edx
	// popfd
	// pop ebx
	// pop edi
	// pop ecx
	// pop esi
	// 0x0AF76F24 0x0AAB7B16
	// label3(0x0AFDEEEE)
	// label2(0x09FDF5DF)
	// label1(0x09EBC65D)

	// f09EBC65D 隐藏函数
	return true
}

// s9 f004E14D3
func f004E4F1ChandleState245(hDC win.HDC) {
	// SEH
	// f00552D0D()
	ebp178 := v0131A270 // 0, 0x28, 0x44
	for ebp178 >= 40 {
		f008AF00D().f008AF06ArecordKey()
		if v012E2340state == 2 || v012E2340state == 4 {
			var ebp184 float64 // v01319D8C.f00A08B5D()
			var ebp408 float64
			if v0114EC40-ebp184 == 0x41 {
				ebp408 = ebp184
			} else {
				ebp408 = v0114EC40
			}
			ebp184 = ebp408 // 0x4069 << 20
			// v01319D8C.f00A08BF0()
			// f0043BF3FgetT4003().f0043C06B() // v01308D18.f0043C06B()
			f004A7D34getWindowManager().f004A9B5B(ebp184)
		}
		// v01308ED4 = 0
		ebp40C := v012E2340state
		switch ebp40C {
		case 2:
			f004E1E1EhandleState2()
		case 4:
			// f004DDD4FhandleState4()
		case 5:
			f004DF0D5handleState5()
		}
		// 0x04E502F
		ebp188 := 0
		for ebp188 < 5 {
			// v0131A27C.f00534AFA(v0114EE48) // float32
			ebp188++
		}
		// f005AC5A0()
		if f005A4BC5queryHotKey(0x2C) { // VK_SNAPSHOT(0x2C): PRINT SCREEN key
			if !v08C88C6AprintScreen {
				v08C88C6AprintScreen = true
			} else {
				v08C88C6AprintScreen = false
			}
		}
		if v08C88F88 > 0 {
			v08C88F88--
		}
		if v086A3BEC > 0 {
			v086A3BEC--
		}
		// v08C7CC18++
		// v08C7CC18 %= 32
		v0131A240++
		ebp178 -= 40
	}
	// 0x004E50FB
	if v01319D65 != 0 {
		return
	}
	// v09D24A20.f00514F8F()
	// f007DB28F()
	// 0x004E511A:
	systime := struct {
		wYear         uint16
		wMonth        uint16
		wDayOfWeek    uint16
		wDay          uint16
		wHour         uint16
		wMinute       uint16
		wSecond       uint16
		wMilliseconds uint16
	}{} // GetLocalTime(&ebp64)
	f00DE817Asprintf(v08C88AB8snapshot[:], "Screen(%02d_%02d-%02d_%02d)-%04d.jpg", systime.wMonth, systime.wDay, systime.wHour, systime.wMinute, v08C88C74)
	ebp220 := v08610600textManager.f00436DF1findstdstring(459)
	ebp410 := string(ebp220.m04data) // [25 73 3A BB AD C3 E6 D2 D1 B4 A2 B4 E6] -> gbk "%s:画面已储存" https://r12a.github.io/app-encodings/
	var ebp174 [100]uint8
	f00DE817Asprintf(ebp174[:], ebp410, v08C88AB8snapshot[:])
	var ebp54 [100]uint8
	// dll.user32.wsprintfA(ebp54[:], " [%s / %s]", f00AF7DC3serverList().f00AF87BAserverName(), v0805BBACobjectself.m38name[:])
	f00DE8010strcat(ebp174[:], string(ebp54[:]))
	ebp10 := 1
	if f008AEFC1(0x10) == 1 { // VK_SHIFT
		ebp10 = 0
	}
	if v08C88C6AprintScreen && ebp10 == 1 {
		var ebp1AC, ebp1C8 stdstring
		ebp1AC.f00406FC0stdstring(ebp174[:])
		ebp1C8.f00406FC0stdstring(nil)
		// f00A49798ui().m124ChatWindow.f00A9FB38print(ebp1C8, ebp1AC, 3, 0)
		// ebp1C8.f00407AC0(1, 0)
		// ebp1AC.f00407AC0(1, 0)
	}
	// 0x004E52A8:
	switch {
	case v012E3EC8mapNumber == 10: // 天空之城
	case v012E3EC8mapNumber == 30: // 罗兰峡谷
	case v012E3EC8mapNumber >= 45 && v012E3EC8mapNumber < 50: // 幻影寺院1～6
	case v012E3EC8mapNumber == 51: // 幻术园
	case v012E3EC8mapNumber == 73:
	case v012E3EC8mapNumber == 74:
	default:
		// dll.opengl32.glClearColor(...)
	}
	// 0x004E5503:
	// dll.opengl32.glClear(0x4100)
	ebp74 := v0131A250
	v0131A250 = win.GetTickCount()
	ebp6D := false
	// 0x004E5529:
	switch v012E2340state {
	case 2:
		ebp6D = f004E46B3handleState2(hDC)
	case 4:
		ebp6D = f004E17B9handleState4(hDC)
	case 5:
		// ebp6D = f004E0E03handleState5(hDC)
	}
	// v0131A27C.f00534BE8()
	if v08C88C6AprintScreen {
		// f006B93DFsnapshot()
	}
	if v08C88C6AprintScreen && ebp10 == 0 {
		var ebp1E4, ebp200 stdstring
		ebp1E4.f00406FC0stdstring(ebp174[:])
		ebp200.f00406FC0stdstring(nil)
		// f00A49798ui().m124ChatWindow.f00A9FB38print(ebp200, ebp1E4, 3, 0)
		// ebp1C8.f00407AC0(1, 0)
		// ebp1AC.f00407AC0(1, 0)
	}
	v08C88C6AprintScreen = false
	if ebp6D {
		win.SwapBuffers(hDC) // dll.gdi32.SwapBuffers(hDC)
	}
	ebp68 := v0131A250 - ebp74
	if ebp68 < 40 {
		offset := 40 - ebp68
		// dll.kernel32.Sleep(offset)
		v0131A250 += offset
		ebp68 = 40
	}
	ebp6C := ebp178 + ebp68
	/*if v0131A23F &&
		v012E2340state == 5 &&
		v08C88FF0conn.f006BD6F9fd() == syscall.Handle(^uint(0)) &&
		v0131A2B4 == false {
		v0131A2B4 = true
		v01319E08log.f00B38AE4printf("> Connection closed.")
		v01319E08log.f00B38D49(1)
		ebp39C := v08610600textManager.f00436DF1findstdstring(381) // "结束游戏"
		ebp41C := string(ebp39C.m04data)
		ebp3E8 := v08610600textManager.f00436DF1findstdstring(402) // "和服务器的连接中断"
		ebp420 := string(ebp3E8.m04data)
		w := f00A49798ui().m88caution
		w.f00A9B633(w, 48, ebp420, ebp41C, 0, 0, 1)
	}*/
	switch v012E2340state {
	case 2:
		// f004D4FB9(v012E23C8, 0) // "data/music/login_theme.mp3"
	case 4:
		// f004D4F77(v012E23C8, 0) // "data/music/login_theme.mp3"
	}
	if v012E2340state == 5 {
		switch v012E3EC8mapNumber {

		}
	}
	// 0x004E621F:
	v0131A270 = ebp6C
}

func f004DB77ChandleState3(hDC win.HDC) {
	// SEH
	// ebp10 := f004A7D34getWindowManager()
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
	v012E2340state = 5
	ebp14 := 0
	for ebp14 < 4 {
		// f0089DB5A(ebp14+0x798D, 0)
		ebp14++
	}
	// f005A4C09(1)
}

// 状态机 s9 f004E2622
func f004E6233handleState(hDC win.HDC) {
	// SEH
	f00DE8A70chkstk() //0x4734
	// st0 = v0638C6C4 // st0~st7,浮点寄存器,80bit
	// f00DE7C90memcpy()
	ebp4738 := v012E2340state
	switch ebp4738 {
	case 1:
		f004DD578handleState1(hDC) // 会阻塞当前主线程，那么这段时间怎么保证响应消息循环？
	case 2, 4, 5:
		f004E4F1ChandleState245(hDC)
	case 3:
		f004DB77ChandleState3(hDC)
	}
	// 0x004E62F0
	f00A49798ui().f00A4E1BFchangeState(v012E2340state)
	// 0x004E6306
	if v01319D9C <= 31 {
		return
	}
	var inhack pb
	inhack.f00439178init()
	// inline send
	inhack.f0043922CwriteHead(0xC1, 0xF1)
	inhack.f004397B1writeUint8(3) // subcode
	inhack.f004397B1writeUint8(1) // flag
	inhack.f004397B1writeUint8(0) // subflag
	inhack.f004393EAsend(true, false)
	inhack.f004391CF()
}
