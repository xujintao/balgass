package main

import (
	"encoding/binary"
	"sync"
)

var v0131A240 uint32
var v0131A26C bool
var v0131A26D bool
var v0131A270 int

var v08C88AB8 [100]uint8
var v08C88C69 bool
var v08C88C6A bool
var v08C88C74 uint32
var v086A3B94 func()
var v086A3BEC uint32

var v0638C6C4 uint32 = 0x4EB0E616

var v0114EE50 = [...]uint64{
	0x3FE3333340000000,
	0x3FC3333340000000,
	0x3F70624DE0000000,
}

var v01319F28 uint32 = 0x3F2DE762

// f004E46B3handleState2
var v09FB8736 uint32 = 0x2B2F4A3D
var v0A4E24C6 uint32 = 0x4849459E
var v09E035DB uint32 = 0x5433BC5D
var v0A43DD91 uint32 = 0x653CAB4B
var v0B06F40C uint32 = 0x08244423 // v0B06F40C = 0x310F1716 = v09E035DB ^ v0A43DD91
var v09FDFB22 uint32 = 0xFC5D207C
var v0AFD3C52 uint32 = 0xD577
var v09FB6D69 uint32 = 0xDEA7
var v0A441BD1 uint32 = 0x38
var v0AD2DD3A uint32 = 0xE4E2E401
var v0AF890C3 uint32 = 9
var v0B287022label1 uint32 = 0x0AD98B31
var v09E8DF92 uint32 = 8
var v0AAB7324label2 uint32 = 0x0AD2CAED
var v0A933705 uint32 = 7
var v0A935B85label3 uint32 = 0x0A43B47B
var v0ABFAB88blocks = [...]block{
	{0x000393EA, 0x0192},
	{0x00039581, 0x0100},
	{0x0003968F, 0x05},
	// ...
}
var v0A4E85ABimageBase uintptr = 0x00400000
var v0AA2E05D uint32 = 0x5AA7D865
var v09EB9D65 uint32 = 1

var v0A9360CB uint32 = 0x4B0FAFAC
var v0AFD365B uint32 = 0x5EC66E14 // v0AFD365B = v0A9360CB

type t7 struct {
	data [0x2A4]byte
}

var v08C86C50 = [10]t7{}

var v0114EC40 float64

var v012E234C = "data\\music\\mutheme.mp3"
var v012E239C = "data\\music\\main_theme.mp3"

// f004E17B9handleState4
var v0A443F74 uint32 = 0xE9497BA1
var v0AD3B896 uint32 = 0x13EFEDCE // v0AD3B896 = v0A443F74
var v0A890E43 uint32 = 0xC491
var v0A56EB4C uint32 = 0xE521
var v0A84A81D uint32 = 0x3E92EB75 // 0xEB
var v09FE1805 uint32 = 0x125      // v0A84A81D 与 v09FE1805比较

var v0ABE324Flabel1 uint32 = 0x09EBC65D
var v0A952A97 uint32 = 8
var v0A92FC07label2 uint32 = 0x09FDF5DF
var v0A920E49 uint32 = 7
var v0A32B0C2label3 uint32 = 0x0AFDEEEE
var v0A6039A7 uint32 = 6

var v0AFD3C89 uint32 = 0x9B20B37E
var v0AA30425 uint32 = 0x105DA588

var v0A88819B uint32 = 0xD266CC5A // key
var v0AF96824backupCode = [...]uint32{
	0x37664E4F,
	0x36FE23B9,
}

var v0A9F69B2blocks = [...]block{
	{0x000393EA, 0x192}, // f004393EA
	{0x00039581, 0x100},
	{0xFFFFFFFF, 0xFFFFFFFF},
}
var v09FE37F3imageBase uintptr = 0x00400000

var v0ABD75F3 uint32 = 0xE55020DE
var v0A3268A5 uint32 = 0xC01C390C

var v0A0C2FD2 uint32 = 0x05BE6BE3
var v0AF8FD45 uint32 = 0x34B17CF5
var v0AA09787 uint32 = 0x0824448B // v0AA09787 = v0A0C2FD2 ^ v0AF8FD45

var v0AD6FFBE uint32 = 0x90B9
var v0ABB4649 uint32 = 0x975D
var v0AAB36E2 uint32 = 0x531BDC38 // v0AD6FFBE * v0AAB36E2 % v0ABB4649
var v09EAE2F8 uint32 = 0x26

var v0AF781CF uint32 = 0x35C75981
var v0A33BBC9blocks = [...]block{}
var v0A5FCB8Bcodes = [...]uint32{}
var v0A84B964 uint32 = 1

var v09FFF8E0imageBase uintptr = 0x00400000

var v0A9FB214 uint32 = 7
var v0A443580label1 uint32 = 0x09EB9FEC
var v0AFDFD32 uint32 = 6
var v09FB8B32label2 uint32 = 0x0AC9CB58
var v0A849070 uint32 = 5
var v0A32C145label3 uint32 = 0x0A338F26

// f004E0E03handleState5
// 确保第一次不覆盖
var v0A56C1BA uint32 = 0x865EB72E
var v0AF8AF60 uint32 = 0xDCD68D02
var v0AD73D05 uintptr = 0xD2435CB5
var v0ABD3FEA uintptr = 0xD36C2725
var v0B07C274 uintptr = 0x09FDEC49 // v0AD73D05 ^ v0ABD3FEA - v0B07C274 + 0x09FDEC49

// 后面在一个确定的时机进行覆盖
var v0A339687 uint32 = 0xB2F9
var v09F87E15 uint32 = 0x8E25
var v09DEA85C uint32 = 0x76237DA5 //
var v0AC3A339 uint32 = 0xE5

var v0A9FD787backupCode = [...]uint32{}
var v0A56AF60blocks = [...]block{
	{0x000393EA, 0x192},
	{0x00039581, 0x100},
	{0xFFFFFFFF, 0xFFFFFFFF},
}

var v09E72FAA uint32 = 0x03B1C032
var v0A933991 uint32 = 0x922C3447 // v0A933991 = v09E72FAA

var v0B07545C uint32 = 8
var v0ABB631Flabel1 uint32 = 0x0A84CD1A
var v0A38E67E uint32 = 7
var v0AC36B4Dlabel2 uint32 = 0x09FC4127
var v0AF87A7B uint32 = 6
var v0AA2D054label3 uint32 = 0x0AF7BD1B

// f00A0A5E1 shell
var v0AA08598 uint32 = 0xE2209428
var v0AD7A978 uint32 = 0xABF8FFC9
var v0A9D361D uint32 = 0xD2B99D3B
var v0A391C78 uint32 = 0x52B79D4F // 0x7FFE0014 = v0A391C78 - v0A9D361D
var v0AD82FDD uint32 = 0x9EA1
var v0A338B93 uint32 = 0xE003
var v0A562D19 uint32 = 0x5AFBC57D
var v09F876D2 uint32 = 0x4334

var v0AD93A11 uint32 = 0x1A65E9C0
var v0AD93A0D = [...]uint32{0x078202AD}

var v0AD574B8 uint32 = 8
var v0A0519E0label1 uint32 = 0x09E2729C

type stdstring struct {
	m04data []uint8
	m14len  int
	m18cap  int
}

func (t *stdstring) f00406A20init() {
	t.m04data = nil
	t.m14len = 0
	t.m18cap = 15
}

func (t *stdstring) f00406EB0(buf []uint8, len int) {
	// s长度小于16就放在栈上(m04~m13)，否则存在堆上
	t.m04data = buf
	t.m14len = len
	t.m18cap = 0x1F
}

func (t *stdstring) f0043D7E2stdstring(s string) {
	t.f00406EB0([]uint8(s), len(s))
}

func (t *stdstring) f00406FC0stdstring(buf []uint8) {
	t.f00406EB0(buf, len(buf))
}

func (t *stdstring) f004079A0stdstring(s *stdstring) {

}

func (t *stdstring) f004073E0cstr() []uint8 {
	return t.m04data
}

func (t *stdstring) f00407B10free() {
	f00DE7538free(nil)
}

// load serverlist.bmd
var v09D965B0serverListManager serverListManager
var v09D96728 sync.Once

func f00AF7DC3getServerListManager() *serverListManager {
	v09D96728.Do(func() {
		v09D965B0serverListManager.f00AF7CC6()
	})
	return &v09D965B0serverListManager
}

type server struct { // 0x54
	// head
	m00 [32]uint8 // "电信1区"
	m20 uint8     // 0
	m21 uint8     // 1,2
	m22 bool      // 0
	m23 [20]uint8 // ?
	// body
	m38 stdstring
}

func (s *server) f00AF88D8server() {

}

func (s *server) f00AF8903server(s1 *server) {}

type serverListManager struct {
	m2Ccount int

	// 结构
	m13C *treeNode
	size int
}

// 构造函数
func (t *serverListManager) f00AF7CC6() {}

func (t *serverListManager) f00AF7EB6dec(buf []uint8, size uint) {

}

func (t *serverListManager) f00AF7F07load() {
	// 0x66C局部变量
	// ebp670 := t
	ebp54 := f00DE909Efopen("Data/Local/ServerList.bmd", "rb")
	var ebp10 uint = 0x3B // 59
	var ebp4AC server
	ebp4AC.f00AF88D8server()

	// unmarshal
	var ebp50head [60]uint8  // head 59字节
	var ebp454body [10]uint8 // body
	for {
		if f00DE8FBDfread(ebp50head[:], ebp10, 1, ebp54) == 0 {
			break
		}
		t.f00AF7EB6dec(ebp50head[:], ebp10)
		size := uint(binary.LittleEndian.Uint16(ebp50head[58:]))

		f00DE8FBDfread(ebp454body[:], size, 1, ebp54)
		t.f00AF7EB6dec(ebp454body[:], size)

		x := binary.LittleEndian.Uint16(ebp50head[:])
		if x >= 500 {

		}
		f00DE9370strncpy(ebp4AC.m00[:], ebp50head[2:], 32)
		ebp4AC.m20 = ebp50head[34]
		ebp4AC.m21 = ebp50head[35]
		if ebp50head[36] != 0 {
			ebp4AC.m22 = true
		} else {
			ebp4AC.m22 = false
		}
		ebp14 := 0
		for {
			if ebp14 >= 20 {
				break
			}
			ebp4AC.m23[ebp14] = ebp50head[37+ebp14]
			ebp14++
		}
		ebp4AC.m38.f0043D7E2stdstring(string(ebp454body[:]))

		var tmp server
		tmp.f00AF8903server(&ebp4AC)
		// f00AFC1D3(&ebp660, binary.LittleEndian.Uint32(ebp50[:]))
		// ebp604.f00AFC22B(ebp678)
		// t.m13C.f00AF8D6E(&ebp66C, &ebp604)
		// ebp604.f00AF897B()
		// ebp660.f00AF8967()
	}
}

func (t *serverListManager) f00AF883FsetCount(count int) {
	t.m2Ccount = count
}

func (t *serverListManager) f00AF8853getCount() int {
	return t.m2Ccount
}

func (t *serverListManager) f00AF81FF(code, percent int) {
	// 0x58局部变量
	// ebp5C := t
	// t.m04.f00AF89F1(&ebp18)
	// ebp19 := false
	// ebp18.f00AF8CCF(t.m04.f00AF8A14(&ebp24))
	// ebp34 := f00DE852Fnew()
	// 	ebp60 := ebp34.f00AF6EB3()
	// 	ebp30 := ebp60
	// 	ebp10 := ebp30
	// 	t.f00AF8387(code/20, ebp10)
	// 	t.m04.f00AF8A42(ebp50, ebp44.f00AF8F97(ebp10, ebp10.m04))
	// 	t.f00AF846C(ebp10, code, percent)
	// 	t.m04.f00AF89F1(&ebp58)
}

// objectManager
var v08C88E58hpMax uint16
var v08C88E5AmpMax uint16
var v08C88E5CsdMax uint16
var v08C88E5EsdMax uint16
var v086105ECobject = &object{}
var v0805BBACself = &object{}

// object sizeof=0x6BC
type object struct {
	m8        uint32
	m13       uint8
	m24       bool
	m2B       uint8
	m38name   [32]uint8
	m5Eid     uint16
	m10C      uint16
	m122hp    uint16 // hp
	m124mp    uint16 // mp
	m126hpMax uint16 // max hp
	m128mpMax uint16 // max mp
	m12Asd    uint16 // sd
	m12CsdMax uint16 // max sd
	m140ag    uint16 // ag
	m142agMax uint16 // max ag
	m154      uint16
	m160      uint16
	m166      uint16
	m178      [100]uint8
	m410      struct {
		m04 bool
		m0E bool
		m28 uint8
	}
	m438 uint8
}

var v012E31B0 int = 0
var v012E3200 int = 5
var v01308D04objectManager = &objectManager{}

type objectManager struct {
	m08objects []object
}

func (t *objectManager) f00A38D5BgetObject(index int) *object {
	if index < 0 || index >= 0x190 {
		return nil
	}
	return &t.m08objects[index]
}

// draw
var v09D96438 *t09D96438

func f00A49798get() *t09D96438 {
	if v09D96438 == nil {
		// ebp14 := f00DE852Fnew(0x228)
		var ebp18 *t09D96438
		ebp14 := &t09D96438{}
		if ebp14 != nil {
			ebp18 = ebp14.f00A49808construct()
		} else {
			ebp18 = nil
		}
		v09D96438 = ebp18
	}
	return v09D96438
}

type baser interface {
	do4() bool
	do6()
	do7() bool
}

type mainFrame struct {
	m30   uintptr
	m2844 []uint8
}

func (t *mainFrame) f00BB1E00(buf []uint8, format string, param ...interface{}) {
	if t.m30 == 0 {
		return
	}
	var ebp4 struct{}
	ebp4.f00C33740(t.m2844, 1, 0x1B, 0)
	var ebp2C struct{}
	f00DF75E5(&ebp2C, 0, 0)
}

type base struct {
	baser
	m08 struct {
		m00 *mainFrame
	}
	m34 int
	m44 struct{}
	m4B bool
}

func (t *base) f00A3A41F() bool {
	return false
}

func (t *base) f00A3A4C1() bool {
	b := t.m4B
	t.m4B = false
	return b
}

func (t *base) f00A4977A() int {
	return t.m34
}

type dash struct {
	base
	m88hp    uint16
	m8ChpMax uint16
	m90      bool
	m94mp    uint16
	m98mpMax uint16
	m9Csd    uint16
	mA0sdMax uint16
	mA4ag    uint16
	mA8agMax uint16
	mB6      bool
	m332     bool
}

// f00AA9458
func (t *dash) do4() bool {
	return false
}

func (t *dash) f00AAA14Ehpmp() {
	var ebp8hpMax uint16
	var ebpCmpMax uint16
	var ebp14hp uint16
	var ebp10mp uint16
	if f0059D4F6bit4(v0805BBACself.m13) {
		var ebp1C, ebp20 uint16
		ebp8hpMax = v08C88E58hpMax
		if v086105ECobject.m122hp >= 0 {
			ebp1C = v086105ECobject.m122hp
		} else {
			ebp1C = 0
		}
		ebp14hp = ebp1C
		ebpCmpMax = v08C88E5AmpMax
		if v086105ECobject.m124mp >= 0 {
			ebp20 = v086105ECobject.m124mp
		} else {
			ebp20 = 0
		}
		ebp10mp = ebp20
	} else {
		var ebp24, ebp28 uint16
		ebp8hpMax = v086105ECobject.m126hpMax // 0x1B6
		if v086105ECobject.m122hp >= 0 {
			ebp24 = v086105ECobject.m122hp
		} else {
			ebp24 = 0
		}
		ebp14hp = ebp24                       // 0x144
		ebpCmpMax = v086105ECobject.m128mpMax // 0x67
		if v086105ECobject.m124mp >= 0 {
			ebp28 = v086105ECobject.m124mp
		} else {
			ebp28 = 0
		}
		ebp10mp = ebp28 // 0x6B
	}
	// if ebp8hpMax > 0 && ebp14hp > 0 {
	// 	ebp1 = v0805BBACself.m60C.f004CEC95(0x37, ebp8hpMax, ebp14hp)
	// }
	// ebp1 := v0805BBACself.m60C.f004CEC95(0x37, ebp8hpMax, ebp14hp)
	ebp1 := false
	if t.m90 != ebp1 {
		t.m90 = ebp1
		// f00A3A4F2(t, "SetChangeIntoxication", "%d", t.m90)
	}
	if t.m88hp != ebp14hp || t.m8ChpMax != ebp8hpMax {
		t.m88hp = ebp14hp
		t.m8ChpMax = ebp8hpMax
		// f00A3A4F2(t, "SetHP", "%d %d", ebp14hp, ebp8hpMax)
	}
	if t.m94mp != ebp10mp || t.m98mpMax != ebpCmpMax {
		t.m94mp = ebp10mp
		t.m98mpMax = ebpCmpMax
		f00A3A4F2(&t.base, "SetMP", "%d %d", ebp10mp, ebpCmpMax)
	}
}

func (t *dash) f00AAA387sd() {
	var ebp8sdMax uint16
	var ebp4sd uint16
	if f0059D4F6bit4(v0805BBACself.m13) {
		var ebp10, ebp12 uint16
		if v08C88E5CsdMax >= 1 {
			ebp10 = v08C88E5CsdMax
		} else {
			ebp10 = 1
		}
		ebp8sdMax = ebp10
		if ebp8sdMax >= v086105ECobject.m12Asd {
			ebp12 = v086105ECobject.m12Asd
		} else {
			ebp12 = ebp8sdMax
		}
		ebp4sd = ebp12
	} else {
		var ebp18, ebp1A uint16
		if v086105ECobject.m12CsdMax >= 1 {
			ebp18 = v086105ECobject.m12CsdMax
		} else {
			ebp18 = 1
		}
		ebp8sdMax = ebp18 // 0x541
		if ebp8sdMax >= v086105ECobject.m12Asd {
			ebp1A = v086105ECobject.m12Asd
		} else {
			ebp1A = ebp8sdMax
		}
		ebp4sd = ebp1A // 0x541
	}
	if t.m9Csd != ebp4sd || t.mA0sdMax != ebp8sdMax {
		t.m9Csd = ebp4sd
		t.mA0sdMax = ebp8sdMax
		// f00A3A4F2(t, "SetSD", "%d %d", ebp4sd, ebp8sdMax)
	}
}

// f00AA94C1
func (t *dash) do6() bool {
	// if v01353C08.f00537917() == false && t.mB6 == false {
	// 	t.mB6 = true
	// 	f00A3A4F2(t, "SetQuestAlarm", "%d", t.mB6)
	// } else if v01353C08.f00537917() == true && t.mB6 == true {
	// 	t.mB6 = false
	// 	f00A3A4F2(t, "SetQuestAlarm", "%d", t.mB6)
	// }
	t.f00AAA14Ehpmp() // hp mp
	t.f00AAA387sd()   // sd
	// t.f00AAA4C9ag() // ag
	// t.f00AAA60Bexp() // exp
	// t.f00AAAA77()
	// t.f00AAAAAD()
	// if t.m332 == false {
	// 	t.m332 = true
	// 	t.f00AAB4D1(0, 0, 0, 0)
	// }
	return true
}

// f00AA95BA
func (t *dash) do7() bool {
	return false
}

func f00A3A4F2(b *base, name string, format string, param ...interface{}) {
	// 0x408局部变量
	ebp4 := param
	var ebp408 [1024]uint8
	f00DF30EFstrcpysafe(ebp408[:], 1024, "_root.")
	f00DECB2Estrcatsafe(ebp408[:], 1024, "g_mcMainFrame") // b.m44.f00A3AF9E()
	f00DECB2Estrcatsafe(ebp408[:], 1024, ".")
	f00DECB2Estrcatsafe(ebp408[:], 1024, name) // "SetSD"
	if b.m08.m00 != nil {                      // if b.m08.f00A3BAB6() != nil {
		b.m08.m00.f00BB1E00(ebp408[:], format, ebp4...) // b.m08.f00A3AF08().f00BB1E00(ebp408[:], format, ebp4)
	}
}

// main window
type windowManager struct {
	m04 uintptr
}

func (t *windowManager) f00A47461fresh() {
	// 0x24局部变量
	ebp1 := false
	// t.m04.f00A483A4(&ebpC)
	for {
		// if false == ebpC.f00A486CA(t.m04.f00A483CB(&ebp20)) {
		// 	return
		// }
		var ebp10 *base // ebp10 := *ebpC.f00A4863A()
		if ebp10 == nil {
			return
		}
		if ebp10.do4() {
			ebp10.do6()
		}
		if ebp1 == false { // if ebp1 == false && v01319D6ChWnd == dll.user32.GetActiveWindow() && v01319D6ChWnd == dll.user32.GetFocus() {
			ebp1 = ebp10.do7()
		}
		if ebp10.f00A4977A() == 5 && ebp10.f00A3A41F() == true {
			ebp1 = true
		}
		if ebp10.f00A3A4C1() {
			// t.f00A47B22(ebp10)
			break
		}
		// ebpC.f00A4864B(&ebp18, 0)
	}
}

// size:0x228
type t09D96438 struct {
	m70 *windowManager
	m80 uintptr
}

func (t *t09D96438) f00A49808construct() *t09D96438 {
	return nil
}

func (t *t09D96438) f00A4DC94(state int) {
	if t.m80 != 0 {
		// t.m80.f008E27CF()
	}
	if t.m70 != nil {
		t.m70.f00A47461fresh()
	}
}
