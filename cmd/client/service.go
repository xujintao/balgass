package main

/*
class base
{
	m0C bool
	m10 int
public:
	virtual do3(int) bool = 0
	f004CCC07(uint64) {}
	f004CCB8A(int) {}
	f004AA068() int {}
}
*/

type base1 interface {
	do3(int) bool
}

type base2 interface {
	f004CCC07(uint64)
	f004CCB8A(int)
	f004AA068() int
}

type base2imp struct {
	m0C bool
	m10 int
}

func (b *base2imp) f004CCB8A(x int) {
	// init logic
}
func (b *base2imp) f004CCC07(unk uint64) {
	// var t *t3
	// t.f0043E60C() // v01313FA8.f0043E60C 发送login报文
	ebp20 := b
	if s.m0C == false {
		return
	}
}
func (b *base2imp) f004AA068() int {
	return b.m10
}

// service1 sizeof(service{})=0x408
type service1 struct {
	base2imp
}

// do3->f0043F608
func (s *service1) do3(x int) bool {}

func (s *service1) f00446D6DreqServerList() {
	// 带SEH处理
	f00DE8A70chkstk() // 0x2994
	// ...
	var ebp149C pb
	ebp149C.f004393EAsend(false, false) // 发送协议报文
	// ...
}

// service2 sizeof(service{})=0x408
type service2 struct {
	base2imp
}

// do3->f0043AA4E
func (s *service2) do3(x int) {}

// service3 sizeof(service{})=0x408
type service3 struct {
	base2imp
}

// do3->f0043E5D4
func (s *service3) do3(x int) {}

// service4 sizeof(service{})=0x408
type service4 struct {
	base2imp
}

// do3->f00446D35
func (s *service4) do3(x int) {}

// service5 sizeof(service{})=0x408
type service5 struct {
	base2imp
}

// do3->f0043DCD3
func (s *service5) do3(x int) {}

// service6 sizeof(service{})=0x408
type service6 struct {
	base2imp
}

// do3->f00441096
func (s *service6) do3(x int) {}

// service7 sizeof(service{})=0x408
type service7 struct {
	base2imp
}

// do3->f00449300
func (s *service7) do3(x int) {}



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
type serviceNode struct {
	next *serviceNode
	prev *serviceNode
	service  interface{}
}
type services struct {
	head *serviceNode
	tail *serviceNode
	num uint32  // service num
}

func (t *services) f004AA077() bool {
	return false
}

func (t *services) f00445530nextServiceNode(nodep **serviceNode) interface{} {
	ebp8 := t
	ebp4 := *nodep
	// if IsBadReadPtr(ebp4, ebp4) {
	// 	return nil
	// }
	*nodep := ebp4.next
	return ebp4.service
}

func (t *services) f004409AAfirstServiceNode() *serviceNode {
	return t.head
}

func (t *services) f004AA009getServiceNum() uint32 {
	return t.num
}

// serviceManager
var v0130F728 serviceManager

type serviceManager struct {
	s001 service1 // v0130F730
	s7   service7 // v0130FB38
	s6   service6 // v0130FF40
	s5   service5 // v01310598
	s4   service4 // v01310798
	s3   service3 // v01313FA8
	s2   service2 // v01314300
	s010 service1 // v01317CD8
	s100 service1 // v013180E0
	// ... // 100个service
	f9FD4activeServices services // v013196FC
	f9FE0 bool
	f9FE1 bool
	f9FE4 int
	f9FE8 bool
	f9FE9 bool
	f9FEC *t4002 // v01319714
}

func f004A7D34() *serviceManager                  { return nil }
func (t *serviceManager) f004A9083(p interface{}) {}
func (t *serviceManager) f004A9123(p interface{}) {}
func (t *serviceManager) f004A9B5B(unk uint64) {
	ebp30 := t
	if ebp30.f9FE4 == 0 { // 0, 2
		return
	}
	if ebp30.f9FD4activeServices.f004AA077() == false {
		return
	}
	if ebp30.f9FE8 {
		// ...
	}

	// ebp10 := f0043BF3F() // v01308D18
	// if ebp10.f0043913E(0x1B) {
	// 	// ...
	// }
	ebp30.f9FE0 = false
	// if ebp10.f00436696() {
	// 	// ...
	// }
	if ebp10.f004366A5() {
		ebp30.f9FE1 = false
	}
	ebp14 := ebp30.f9FD4activeServices.f004AA009getServiceNum() // 9
	ebp8 := make([]interface{}, 9)   // ebp8 := f00DE64BC(ebp14 * 4) // f00DE64BC(4*9) new?
	ebp4 := ebp30.f9FD4activeServices.f004409AAfirstServiceNode()  // 一个虚拟内存地址0x36C8E8E8
	ebp20 := 0
	for ebp20 < ebp14 {
		ebp8[ebp20] = ebp30.f9FD4activeServices.f00445530nextServiceNode(&ebp4)
		ebp8[ebp20].(base2).f004CCB8A(0)
		ebp20++
	}

	ebp4 = ebp30.f9FD4activeServices.f004409AAfirstServiceNode()
	for ebp4 != nil {
		ebpC := ebp30.f9FD4activeServices.f00445530nextServiceNode(&ebp4)
		ebpC.(base1).do3(0)
		ebpC.(base2).f004CCB8A(1)
		break
	}

	ebp24 := 0
	for ebp24 < ebp14 {
		ebp8[ebp24].(base2).f004CCC07(unk)
		ebp24++
	}

	if len(ebp8) > 0 {
		// f00DE7BEA(ebp8) // delete?
		ebp8 = ebp8[:0]
	}
	// ebp30.f004A91CE()
	ebp4 = ebp30.f9FD4activeServices.f004409AAfirstServiceNode()
	for ebp4 != nil {
		ebpC := ebp30.f9FD4activeServices.f00445530nextServiceNode(&ebp4)
		ebp34 := ebpC.(base2).f004AA068()
		switch ebp34 {
		case 1, 2, 3, 4:
			ebp30.f9FE0 = true
		}
		if ebp30.f9FE0 {
			break
		}
		if ebpC.(base1).do3(0) == false {
			break
		}
		ebp30.f9FE0 = true
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
