package main

import (
	"sync/atomic"
)

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

func (t *windowManager1) f00BB0920(x int, ws ...*window011737C4) {
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

func (t *windowManager1) f00BB0D20construct(ws ...*window011737C4) {
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
