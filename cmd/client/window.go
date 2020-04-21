package main

import (
	"encoding/binary"
	"sync/atomic"
)

type i011737C4 interface {
	f00BAE450()
	f00BAE460()
}

type window011737C4 struct {
	m00vtabptr uintptr
	m04        int32
}

func (t *window011737C4) f00A3BB4Cconstruct() {
	// t.m00vtabptr = &v011737C4
	t.m04 = 1
}

func (t *window011737C4) f00BAE450() {
	atomic.AddInt32(&t.m04, 1)
}

func (t *window011737C4) f00BAE460() {
	atomic.AddInt32(&t.m04, -1)
}

type window011927E8 struct {
	window011737C4
}

type window01187D5C struct {
	window011737C4
	m08 int
	m0C *window011927E8
}

func (t *window01187D5C) f00BBBA00construct() {
	// inline
	t.window011737C4.f00A3BB4Cconstruct()
	t.m08 = 0x1D
	// t.m00vtabptr = &v01187D5C
	// f00C87090()
	t.m0C = func() *window011927E8 {
		// v09D9BD74mm.do11malloc(8, 0)
		w := new(window011927E8) // 0x0EB41EC0
		// inline
		w.window011737C4.f00A3BB4Cconstruct()
		// w.m00vtabptr = &v011927E8
		return w
	}()
}

type window011737BC struct {
	window011737C4
}

func (t *window011737BC) f00A3BB2Fconstruct() {
	t.window011737C4.f00A3BB4Cconstruct()
	// t.m00vtabptr = &v011737BC
}

type window011737B4 struct {
	window011737BC
}

func (t *window011737B4) f00A3BAEBconstruct() {
	t.window011737BC.f00A3BB2Fconstruct()
	// t.m00vtabptr = &v011737B4
}

type window011737AC struct {
	window011737B4
}

func (t *window011737AC) f00A3BA4Dconstruct() {
	t.window011737B4.f00A3BAEBconstruct()
	// t.m00vtabptr = &v011737AC
}

type window011737A4 struct {
	window011737AC
	m08 int
}

func (t *window011737A4) f00A3B2ABconstruct(x int) {
	t.window011737AC.f00A3BA4Dconstruct()
	// t.m00vtabptr = &v011737A4
	t.m08 = x
}

type window0117438C struct {
	window011737A4
}

func (t *window0117438C) f00A50070construct() {
	t.window011737A4.f00A3B2ABconstruct(0x1E)
	// t.m00vtabptr = &v0117438C
}

type window0117437C struct {
	window0117438C
}

func (t *window0117437C) f00A50053construct() {
	t.window0117438C.f00A50070construct()
	// t.m00vtabptr = &v0117437C
}

type window0117439C struct {
	window011737A4
}

func (t *window0117439C) f00A50148construct() {
	t.window011737A4.f00A3B2ABconstruct(0xA)
	// t.m00vtabptr = &v0117439C
	// v09D9BD70 = 0
}

type window01174368 struct {
	window0117439C
}

func (t *window01174368) f00A50036construct() {
	t.window0117439C.f00A50148construct()
	// t.m00vtabptr = &v01174368
}

type window01187A98 struct {
	window011737C4
	m08 *window01187A80
	// m0C spintex // 占用0x18字节
	m24 int
	m28 uintptr
}

func (t *window01187A98) f00BB44B0construct(parent *window01187A80) {
	// inline
	t.window011737C4.f00A3BB4Cconstruct()
	// t.m00vtabptr = v01187A98[:]
	// t.m0C.f00BFA480(0)
	t.m08 = parent
	t.m24 = 0
	t.m28 = 0
	// var ebp20 struct{}
	// t.m28 = v09D9BD74.do5("_ResourceLib_Images". &ebp20)
}

type window01187A80 struct {
	window011737C4
	m08 *window01187A98
	m0C int
	m10 bool
}

func (t *window01187A80) f00BB4930construct(x bool) {
	// inline
	t.window011737C4.f00A3BB4Cconstruct()
	// t.m00vtabptr = v01187A80[:]
	t.m0C = 0
	t.m10 = x
	// v09D9BD74mm.do11malloc(0x2C, 0)
	w := new(window01187A98)
	w.f00BB44B0construct(t)
	t.m08 = w
}

type window0118CB54 struct {
	// window01187934
}
type window0118CB4C struct {
	// window0118CA34
}

type window0118CC10 struct {
	// window01187934
}
type window0118CC08 struct {
	// window0118CA2C
}

type window0118CC24 struct {
	window011737C4
	window0118CC10
	window0118CC08
	m14 uintptr
	// m1C spintex
}

type window0118CBD0 struct {
	// window0118CA24
}
type window0118CBD8 struct {
	window011737C4
	m08 int
	window0118CBD0
}

type window0118CBC8 struct {
	window011737C4
	m08 int
	m0C uint8
}

type window011934EC struct {
	window011737C4
	m08 int
	m0C int
	m10 int
	m14 int
	m18 int
	m1C int
	m20 int
	m24 int
	m28 int
	m2C float64
	m34 int
	m38 float64
	m3C float64
	m48 int
}

func (t *window011934EC) f00CA0580construct(x, y int) {

}

type window01188DFC struct {
	window011737C4
	m08 int
	m0C int
	m10 int
	m14 int
	m18 int
	m1C int
}

func (t *window01188DFC) f00BDCB40construct() {
	t.f00A3BB4Cconstruct()
	// t.m00vtabptr = v01188DFC[:]
}

type window01188FE8 struct {
	window011737C4
	m08 int
	m0C int
	m10 int
	m14 int
}

func (t *window01188FE8) f00BE0250construct() {

}
func (t *window01188FE8) f00BE0320() {

}

type window0118CB68 struct {
	window011737C4
	window0118CB54
	window0118CB4C
	m14 *window0118CC24
	m18 *window01187A98
	m1C struct {
		m00 uintptr
		m04 uintptr
	}
	// m24 spintex
	m3C uint8
}

func (t *window0118CB68) f00BF5350construct(x *window01187A80, y *int) {
	// 构造虚表指针
	// t.m24.f00BFA480(0)
	// t.m3C = uint8(y & 0xFF)
	if x != nil {
		if x.m08 != nil {
			x.m08.f00BAE450()
		}
		if t.m18 != nil {
			t.m18.f00BAE460()
		} else {
			t.m18 = x.m08
		}
	}
	{
		// v09D9BD74mm.do11malloc(0x34, 0)
		w := new(window0118CC24)
		// inline构造虚表
		w.m14 = 0
		// w.m1C.f00BFA480(0)
		t.m14 = w
	}
	{
		// v09D9BD74mm.do11malloc(0x10, 0)
		w := new(window0118CBD8)
		// inline构造虚表
		// t.m14.window0118CC10.do3(4, w)
		w.f00BAE460()
	}
	{
		// v09D9BD74mm.do11malloc(0x10, 0)
		w := new(window0118CBC8)
		// t.m14.window0118CC10.do3(0xC, w)
		w.f00BAE460()
		// t.m14.window0118CC10.do3(0x11, w)
	}
	{
		// v09D9BD74mm.do11malloc(0x4C, 0)
		w := new(window011934EC)
		w.f00CA0580construct(1, int(t.m3C))
		// t.m14.window0118CC10.do3(0x12, w)
		w.f00BAE460()
	}
	{
		// v09D9BD74mm.do11malloc(0x20, 0)
		w := new(window01188DFC)
		w.f00BDCB40construct()
		// t.m14.window0118CC10.do3(0x18, w)
		w.f00BAE460()
	}
	{
		// v09D9BD74mm.do11malloc(0x18, 0)
		w := new(window01188FE8)
		w.f00BE0250construct()
		w.f00BE0320()
		// t.m14.window0118CC10.do3(0x19, w)
		w.f00BAE460()
	}
}

type window01187A44 struct {
	window011737C4
	m08 int
	m0C int
}

type window01187920 struct {
	window011737C4
	m08 int
	m0C int
	m10 int
	m14 int
	m18 int
}

func (t *window01187920) f00BAFBA0construct(x int) {

}

type window01174340 struct {
	window01188DFC
}

func (t *window01174340) f00A4FAD6construct() {
	t.window01188DFC.f00BDCB40construct()
	// t.m00vtabptr = v01174340[:]
}

type windowManager1 struct {
	m00vtabptr []uintptr
	m04        *window0118CB68
	m08        *window01187A80
	m0C        int
}

func (t *windowManager1) do3(x int, window interface{}) {
	// t.m04.do3()
}

func (t *windowManager1) f00BB0920(x int, ws ...i011737C4) {
	{
		// v09D9BD74mm.do11malloc(0x14, 0)
		w := new(window01187A80)
		w.f00BB4930construct(false)
		t.m08 = w
		t.m0C = x
	}
	{
		// v09D9BD74mm.do11malloc(0x40, 0)
		w := new(window0118CB68)
		w.f00BF5350construct(t.m08, &x)
		t.m04 = w
	}
	t.do3(0x10, ws[0])
	{
		// v09D9BD74mm.do11malloc(0x10, 0)
		w := new(window01187A44)
		// t.do3(0x0D, w)
		w.f00BAE460()
	}
	t.do3(0x1E, ws[1])
	t.do3(0x1D, ws[2])
	{
		// v09D9BD74mm.do11malloc(0x1C, 0)
		w := new(window01187920)
		w.f00BAFBA0construct(x)
	}
}

func (t *windowManager1) f00BB0D20construct(ws ...i011737C4) {
	// t.m00vtabptr = v01187948[:]
	for _, w := range ws {
		if w != nil {
			w.f00BAE450()
		}
	}
	t.f00BB0920(0, ws...)
	for _, w := range ws {
		if w != nil {
			w.f00BAE460()
		}
	}
}

// -------------------------------------------------------------------------

// size:0x34
type windowgame0116A864 struct {
	m00vtabptr []uintptr
	m04        struct{}
	m20        int
	m24        int
	m28        struct{}
}

func (t *windowgame0116A864) f008E20BBconstruct() {
	// ebp10 := t
	// t.m00vtabptr = v0116A864[:]
	// t.m04.f008E294F()
	// t.m28.f00518DB4()
	t.m20 = 0
	t.m24 = 0x80000000
	// t.m28.f00518DDD(1000)
}

func (t *windowgame0116A864) f008E27CFfresh() {

}

type infoTooltip struct{}

func (t *infoTooltip) f00A3C455construct() {}

type infoTooltipText struct{}

func (t *infoTooltipText) f00A3C774construct() {}

// InfoTooltip size:0x40
type windowgameInfoTooltip struct {
	m00 infoTooltip
	m20 infoTooltipText
}

func (t *windowgameInfoTooltip) f00A3BDC3load(s *stdstring) bool {
	// 0x58局部变量
	// ebp60 := t
	ebp18 := f00DE909Efopen(string(s.f004073E0cstr()), "rb")
	if ebp18 == nil {
		// dll.user32.MessageBox(0, "InfoTooltip.bmd File Open Failed", "TxtConvertor", 0)
		return false
	}

	// read count
	var ebp24 [4]uint8
	f00DE8FBDfread(ebp24[:], 4, 1, ebp18)
	ebp24count := uint(binary.LittleEndian.Uint32(ebp24[:])) // 0x16
	if ebp24count <= 0 {
		// dll.user32.MessageBox(0, "Data Count 0", "TxtConvertor", 0)
		return false
	}

	// read elements
	var ebp28 uint = 0x2A // element size
	ebp20 := ebp24count * ebp28
	ebp34 := f00DE64BCnew(ebp20)
	ebp10 := ebp34 // type convert
	f00DE8FBDfread(ebp10, ebp20, 1, ebp18)

	// read crc and check
	var ebp14 [4]uint8
	f00DE8FBDfread(ebp14[:], 4, 1, ebp18)
	ebp14crc := uint(binary.LittleEndian.Uint32(ebp14[:]))
	if ebp14crc != 0xAC6B944A { // if ebp14crc != f004EE17Fcalc(ebp10, ebp20, 0xA4C6) {
		// dll.user32.MessageBox(0, "InfoToolTip.bmd File corrupted", "TxtConvertor", 0)
		return false
	}

	ebp1C := ebp10
	var ebp2C uint
	for {
		if ebp2C >= ebp24count {
			break
		}
		f0043B750xor(ebp1C, int(ebp28))
		// ebp30 := &struct {
		// 	index uint16
		// 	data [40]uint8
		// }{} // ebp30 := f00DE852Fnew(0x2A)
		// f00DE7C90memcpy(ebp30, ebp1C, int(ebp28))
		// ebp44.f00A3FD45(f00A3FD2D(&ebp4C, binary.LittleEndian.Uint16(ebp30), ebp30))
		// t.m00.f00A3C4F9(&ebp58, &ebp44)
		ebp1C = ebp1C[ebp28:]
		ebp2C++
	}
	f00DE8C84close(ebp18)
	if ebp10 != nil {
		// ebp5C := ebp10 // type convert
		// f00DE7538free(ebp5C)
	}
	return true
}

func (t *windowgameInfoTooltip) f00A3BFBCload(s *stdstring) bool {
	// 0x58局部变量
	// ebp60 := t
	ebp18 := f00DE909Efopen(string(s.f004073E0cstr()), "rb")
	if ebp18 == nil {
		// dll.user32.MessageBox(0, "InfoToolTipText.bmd File Open Failed", "TxtConvertor", 0)
		return false
	}

	// read count
	var ebp24 [4]uint8
	f00DE8FBDfread(ebp24[:], 4, 1, ebp18)
	ebp24count := uint(binary.LittleEndian.Uint32(ebp24[:])) // 8
	if ebp24count <= 0 {
		// dll.user32.MessageBox(0, "Data Count 0", "TxtConvertor", 0)
		return false
	}

	// read elements
	var ebp28 uint = 0x102 // element size
	ebp20 := ebp24count * ebp28
	ebp34 := f00DE64BCnew(ebp20)
	ebp10 := ebp34 // type convert
	f00DE8FBDfread(ebp10, ebp20, 1, ebp18)

	// read crc and check
	var ebp14 [4]uint8
	f00DE8FBDfread(ebp14[:], 4, 1, ebp18)
	ebp14crc := uint(binary.LittleEndian.Uint32(ebp14[:]))
	if ebp14crc != 0x7D2AEC3C { // if ebp14crc != f004EE17Fcalc(ebp10, ebp20, 0xA4C6) {
		// dll.user32.MessageBox(0, "InfoToolTip.bmd File corrupted", "TxtConvertor", 0)
		return false
	}

	ebp1C := ebp10
	var ebp2C uint
	for {
		if ebp2C >= ebp24count {
			break
		}
		f0043B750xor(ebp1C, int(ebp28))
		// ebp30 := f00DE852Fnew(0x102)
		// ebp30 := &struct {
		// 	index uint16
		// 	data [256]uint8
		// }{}
		// f00DE7C90memcpy(ebp30, ebp1C, int(ebp28))
		// ebp44.f00A3FD83(f00A3FD6B(&ebp4C, binary.LittleEndian.Uint16(ebp30), ebp30))
		// t.m20.f00A3C818(&ebp58, &ebp44)
		ebp1C = ebp1C[ebp28:]
		ebp2C++
	}
	f00DE8C84close(ebp18)
	if ebp10 != nil {
		// ebp5C := ebp10 // type convert
		// f00DE7538free(ebp5C)
	}
	return true
}

func (t *windowgameInfoTooltip) f00A3BBB9construct() {
	// 0x3C局部变量
	// ebp48 := t
	t.m00.f00A3C455construct()
	t.m20.f00A3C774construct()

	var ebp28 stdstring
	ebp28.f00406FC0stdstring([]byte("Data/Local/InfoTooltip.bmd"))
	t.f00A3BDC3load(&ebp28)
	ebp28.f00407B10free()

	var ebp44 stdstring
	ebp44.f00406FC0stdstring([]byte("Data/Local/InfoTooltipText.bmd")) // 有新加入战盟申请 有新加入队伍的申请
	t.f00A3BFBCload(&ebp44)
	ebp44.f00407B10free()
}

//-------------------------------------------------------------

// barer implements by derived class
type barer interface {
	do1(bool)
	do2()
	do3()
	do4() bool
	do6() bool
	do7() bool
	do16()
}

// fooer implements by base class
type iwindowgame0117373C interface {
	barer
	f00A3A717(bool)
	f00A4977A() int
	f00A3A41F() bool
	f00A3A4C1() bool
}

// base
type windowgame0117373C struct {
	barer // m00vtabptr []uintptr
	m08   struct {
		m00 *mainFrame
	}
	m34 int
	m44 struct{}
	m4B bool
	m49 bool
}

func (t *windowgame0117373C) f00A392A6construct() {
	// t.m00vtabptr = v0117373C[:]
}

func (t *windowgame0117373C) f00A3A717(x bool) {
	t.do16()
	t.m49 = x
}

func (t *windowgame0117373C) f00A3A41F() bool {
	return false
}

func (t *windowgame0117373C) f00A3A4C1() bool {
	b := t.m4B
	t.m4B = false
	return b
}

func (t *windowgame0117373C) f00A4977A() int {
	return t.m34
}

// Caution size:0xA4
type windowgameCaution struct {
	windowgame0117373C
}

func (t *windowgameCaution) f00A99F11construct() {
	t.f00A392A6construct()
	// t.m00vtabptr = v01179074[:]
}

func (t *windowgameCaution) do1(x bool) {

}

// 0x00A39479
func (t *windowgameCaution) do2() {

}
func (t *windowgameCaution) do4() bool { return false }
func (t *windowgameCaution) do6() bool { return false }
func (t *windowgameCaution) do7() bool { return false }

// 0x00A9A05C
func (t *windowgameCaution) do3() {

}

func (t *windowgameCaution) do16() {

}

// mainFrame or dash size:0x340
type windowgameMainFrame struct {
	windowgame0117373C //base
	m88hp              uint16
	m8ChpMax           uint16
	m90                bool
	m94mp              uint16
	m98mpMax           uint16
	m9Csd              uint16
	mA0sdMax           uint16
	mA4ag              uint16
	mA8agMax           uint16
	mB6                bool
	m332               bool
}

func (t *windowgameMainFrame) f00AA9021construct() {
	t.windowgame0117373C.f00A392A6construct()
}

func (t *windowgameMainFrame) do1(x bool) {

}

func (t *windowgameMainFrame) do2() {

}
func (t *windowgameMainFrame) do3() {

}

// f00AA9458
func (t *windowgameMainFrame) do4() bool {
	return false
}

func (t *windowgameMainFrame) f00AAA14Ehpmp() {
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
		f00A3A4F2(t, "SetMP", "%d %d", ebp10mp, ebpCmpMax)
	}
}

func (t *windowgameMainFrame) f00AAA387sd() {
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
func (t *windowgameMainFrame) do6() bool {
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
func (t *windowgameMainFrame) do7() bool {
	return false
}

func (t *windowgameMainFrame) do16() {

}

// GuildPosition size:0x78
type windowgameGuildPosition struct {
}

// ...

// size:0x28
type windowManager01173C54 struct {
	m00vtabptr []uintptr
	m04        struct{}
	m20        int
	m24        int
}

func (t *windowManager01173C54) f00A472A7construct() {
	// t.m00vtabptr = v01173C54[:]
	// t.m04.f00A48308()
	// t.m04.f00A48522()
	t.m20 = 0
	t.m24 = 0
}

func (t *windowManager01173C54) f00A47FAB(w iwindowgame0117373C) {
	// t.m04.f00A4843C(&w)
}

func (t *windowManager01173C54) f00A47461fresh() {
	// 0x24局部变量
	ebp1 := false
	// t.m04.f00A483A4(&ebpC)
	for {
		// if false == ebpC.f00A486CA(t.m04.f00A483CB(&ebp20)) {
		// 	return
		// }
		var ebp10 iwindowgame0117373C // ebp10 := *ebpC.f00A4863A()
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

func f00A3A4F2(w iwindowgame0117373C, name string, format string, param ...interface{}) {
	// 0x408局部变量
	// ebp4 := param
	var ebp408 [1024]uint8
	f00DF30EFstrcpysafe(ebp408[:], 1024, "_root.")
	f00DECB2Estrcatsafe(ebp408[:], 1024, "g_mcMainFrame") // b.m44.f00A3AF9Ecstr()
	f00DECB2Estrcatsafe(ebp408[:], 1024, ".")
	f00DECB2Estrcatsafe(ebp408[:], 1024, name) // "SetSD"
	// if w.m08.m00 != nil {                      // if b.m08.f00A3BAB6() != nil {
	// 	w.m08.m00.f00BB1E00(ebp408[:], format, ebp4...) // b.m08.f00A3AF08().f00BB1E00(ebp408[:], format, ebp4)
	// }
}
