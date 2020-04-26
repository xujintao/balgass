package main

import (
	"encoding/binary"
	"io"
	"os"
	"sync/atomic"
)

type attr struct {
	m00 int     // 9
	m04 int     // 0x10
	m08 int     // 0x400
	m0C int     // 0x300
	m10 float32 // 0x41F00000
	m14 int     // 1
	m18 int     // 0
	m1C int16   // 0x36F
	m20 int     // 0
}

type i011737C4 interface {
	f00BAE450()
	f00BAE460()
}

// 但是shared_ptr的引用计数器和托管对象不在一起
type t011737C4base struct {
	m00vtabptr []uintptr
	m04        int32
}

func (t *t011737C4base) f00A3BB4Cconstruct() {
	// t.m00vtabptr = &v011737C4
	t.m04 = 1
}

func (t *t011737C4base) f00BAE450() {
	atomic.AddInt32(&t.m04, 1)
}

func (t *t011737C4base) f00BAE460() {
	atomic.AddInt32(&t.m04, -1)
}

type window011927E8 struct {
	t011737C4base
}

type t01187D5C struct {
	t011737C4base
	m08 int
	m0C *window011927E8
}

func (t *t01187D5C) f00BBBA00construct() {
	// inline
	t.t011737C4base.f00A3BB4Cconstruct()
	t.m08 = 0x1D
	// t.m00vtabptr = &v01187D5C
	// f00C87090()
	t.m0C = func() *window011927E8 {
		// v09D9BD74mm.do11malloc(8, 0)
		w := new(window011927E8) // 0x0EB41EC0
		// inline
		w.t011737C4base.f00A3BB4Cconstruct()
		// w.m00vtabptr = &v011927E8
		return w
	}()
}

type t011737BC struct {
	t011737C4base
}

func (t *t011737BC) f00A3BB2Fconstruct() {
	t.t011737C4base.f00A3BB4Cconstruct()
	// t.m00vtabptr = &v011737BC
}

type t011737B4 struct {
	t011737BC
}

func (t *t011737B4) f00A3BAEBconstruct() {
	t.t011737BC.f00A3BB2Fconstruct()
	// t.m00vtabptr = &v011737B4
}

type window011737AC struct {
	t011737B4
}

func (t *window011737AC) f00A3BA4Dconstruct() {
	t.t011737B4.f00A3BAEBconstruct()
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

type t0117437C struct {
	window0117438C
}

func (t *t0117437C) f00A50053construct() {
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
	t011737C4base
	m08 *t01187A80
	// m0C spintex // 占用0x18字节
	m24 int
	m28 uintptr
}

func (t *window01187A98) f00BB44B0construct(parent *t01187A80) {
	// inline
	t.t011737C4base.f00A3BB4Cconstruct()
	// t.m00vtabptr = v01187A98[:]
	// t.m0C.f00BFA480(0)
	t.m08 = parent
	t.m24 = 0
	t.m28 = 0
	// var ebp20 struct{}
	// t.m28 = v09D9BD74.do5("_ResourceLib_Images". &ebp20)
}

type t01187A80 struct {
	t011737C4base
	m08 *window01187A98
	m0C int
	m10 bool
}

func (t *t01187A80) f00BB4930construct(x bool) {
	// inline
	t.t011737C4base.f00A3BB4Cconstruct()
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
	t011737C4base
	window0118CC10
	window0118CC08
	m14 uintptr
	// m1C spintex
}

type window0118CBD0 struct {
	// window0118CA24
}
type t0118CBD8 struct {
	t011737C4base
	m08 int
	window0118CBD0
}

type window0118CBC8 struct {
	t011737C4base
	m08 int
	m0C uint8
}

type window011934EC struct {
	t011737C4base
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
	t011737C4base
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
	t011737C4base
	m08 int
	m0C int
	m10 int
	m14 int
}

func (t *window01188FE8) f00BE0250construct() {

}
func (t *window01188FE8) f00BE0320() {

}

type ifile interface {
	i011737C4
	do3valid() bool
	do6getFileOffset() int64
	do7getSize() int
	do11read(buf []uint8, size int) int
	do19close() bool
}

// size:0x20 stdcall且无ebp帧栈
type t01189700file struct {
	t011737C4base
	m08path   *xstring
	m0Cvalid  bool
	m10fd     *os.File
	m14mode   int
	m18err    int
	m1Cwhence int
}

func (t *t01189700file) f00BEED40open() {
	// 0x2C局部变量
	mode := "rb" // t.m14mode
	f00DF6F11fopen(&t.m10fd, []rune(t.m08path.f00A3AF9Ecstr()), []rune(mode))
	// f00DF777C(t.m10fd)
	t.m0Cvalid = true
	if t.m0Cvalid == false {

	}
	t.m18err = 0
	t.m1Cwhence = io.SeekStart
}

// f00BEEED0
func (t *t01189700file) do3valid() bool {
	return t.m0Cvalid
}

// f00BEEF00
func (t *t01189700file) do5ftell() int64 {
	pos := f00DEFCD4ftell(t.m10fd)
	if pos < 0 {
		// ...
	}
	return pos
}

// 0x00BEEF70
func (t *t01189700file) do6getFileOffset() int64 {
	pos := f00DEFCD4ftell(t.m10fd)
	if pos < 0 {
		// ...
	}
	return pos
}

// f00BEEFE0
func (t *t01189700file) do7getSize() int {
	l := int(t.do5ftell())
	if l < 0 {
		return -1
	}
	defer t.do15seek(l, io.SeekStart)
	return int(t.do5ftell())
}

// f00BEF0F0
func (t *t01189700file) do11read(buf []uint8, size int) int {
	if t.m1Cwhence != io.SeekStart && t.m1Cwhence != io.SeekCurrent {
		// f00DF4A9Eclearerr(t.m10fd)
	}
	t.m1Cwhence = io.SeekCurrent
	n := f00DE8FBDfread(buf, 1, uint(size), t.m10fd)
	if n >= uint(size) {
		return int(n)
	}
	f00DF7693errno()
	return 0
}

// f00BEF290
func (t *t01189700file) do15seek(offset int, whence int) int {
	f00DEFA34fseek(t.m10fd, 0, whence)
	return int(t.do5ftell())
}

// f00BEF3E0
func (t *t01189700file) do19close() bool {
	if f00DE8C84close(t.m10fd) != 0 {
		t.m18err = f00DF7693errno()
	}
	t.m0Cvalid = false
	t.m10fd = nil
	t.m18err = 0
	return true
}

type t0118C870 struct {
	t011737C4base
}

func (t *t0118C870) construct() {
	t.t011737C4base.f00A3BB4Cconstruct()
	// t.m00vtabptr = v0118C870[:]
}

// size:0x28 stdcall且无ebp帧栈
type t01195FE0filebuf struct {
	t0118C870
	m08ifile  ifile
	m0Cbuf    []uint8 // size:0x1FF8 不到8k
	m10whence int
	m14r      int
	m18w      int
	m20offset int64
}

func (t *t01195FE0filebuf) f00D02B60construct(f ifile) {
	t.t0118C870.construct()
	if f != nil {
		f.f00BAE450()
	}
	t.m08ifile = f
	// t.m00vtabptr = v01195FE0[:]
	t.m0Cbuf = make([]uint8, 0x1FF8) // t.m0Cbuf = v09D9BD74mm.do10malloc(0x1FF8,0x20,0) 且置零
	t.m10whence = io.SeekStart
	t.m20offset = f.do6getFileOffset()
	t.m14r = 0
	t.m18w = 0
}

func (t *t01195FE0filebuf) do3valid() bool {
	return t.m08ifile.do3valid()
}

func (t *t01195FE0filebuf) do6getFileOffset() int64 {
	return t.m08ifile.do6getFileOffset()
}

// f00D02500
func (t *t01195FE0filebuf) do7getSize() int {
	return t.m08ifile.do7getSize()
}

func (t *t01195FE0filebuf) f00D023D0() {
	if t.m10whence == io.SeekCurrent {
		// ...
	} else if t.m10whence == io.SeekEnd {
		// ...
	}
	return
}

// f00D029A0
func (t *t01195FE0filebuf) do11read(buf []uint8, size int) int {
	if t.m10whence != io.SeekCurrent {
		t.f00D023D0()
		t.m10whence = io.SeekCurrent
		t.m14r = 0
		t.m18w = 0
	}
	eaxr := t.m14r
	ebxw := t.m18w
	ecxbuf := t.m0Cbuf
	edi := size
	ebxw -= eaxr
	ecxbuf = ecxbuf[eaxr:]
	if ebxw >= edi {
		f00DE7C90memcpy(buf, ecxbuf, edi)
		t.m14r = edi
		return 0
	}
	// 读一次
	ebp := buf
	f00DE7C90memcpy(ebp, ecxbuf, ebxw)
	t.m14r = t.m18w
	edi -= ebxw
	ebp = ebp[ebxw:]

	if edi > 0x1000 {
		// 直接io
		n := t.m08ifile.do11read(ebp, edi)
		if n > 0 {
			t.m20offset += int64(n)
			t.m14r = 0
			t.m18w = 0
		}
		return n
	}
	// t.f00D02440()
	// ecx = t.m14
	// eax = t.m18
	// eax -= ecx
	// if eax < edi {
	// 	edi = eax
	// }
	// eax = t.m0C
	// eax += ecx
	// f00DE7C90memcpy(ebp, eax, edi)
	// t.m14 = edi
	// return // ebx+edi
	return 0
}

// 0x00D02810
func (t *t01195FE0filebuf) do19close() bool {
	switch t.m10whence {
	case io.SeekCurrent:
		t.m10whence = io.SeekStart
	case io.SeekEnd:
		t.f00D023D0()
	}
	return t.m08ifile.do19close()
}

// size:0x08
type t0118C8C0fileempty struct {
	t011737C4base
}

func (t *t0118C8C0fileempty) construct() {
	t.t011737C4base.f00A3BB4Cconstruct()
	// t.m00vtabptr = v0118C8C0[:]
}

// size:0xC stdcall有ebp帧栈，应该是业务代码
type t0118C910fileio struct {
	t011737C4base
	m08ifile ifile
}

func (t *t0118C910fileio) f00BEFF80construct(path *xstring, mode, y int) {
	// inline
	// t.t011737C4base.f00A3BB4Cconstruct()
	// t.m00vtabptr = v0118C910[:]

	// 构造file
	// t.f00BEFE50(path, mode, y)
	func() bool {
		// fd
		// f00BEF630(path, mode, y)
		t.m08ifile = func() ifile {
			f := new(t01189700file)
			// inline
			f.t011737C4base.f00A3BB4Cconstruct()
			// w.m00vtabptr = v01189700[:]
			f.m08path.f00BADF50assign(path)
			f.m14mode = mode
			f.f00BEED40open()
			return f
		}()
		if t.m08ifile.do3valid() == false {
			// ...
			return false
		}
		// bufio.Reader
		w := new(t01195FE0filebuf)
		w.f00D02B60construct(t.m08ifile)

		t.m08ifile.f00BAE460()
		t.m08ifile = w
		return true
	}()
}

// f00BEFA80
func (t *t0118C910fileio) do3valid() bool {
	return t.m08ifile.do3valid()
}

func (t *t0118C910fileio) do6getFileOffset() int64 {
	return 0
}

// f00AEEE51
func (t *t0118C910fileio) do7getSize() int {
	return t.m08ifile.do7getSize()
}

// f00AEEE74
func (t *t0118C910fileio) do11read(buf []uint8, size int) int {
	return t.m08ifile.do11read(buf, size)
}

// f00BEFF20
func (t *t0118C910fileio) do19close() bool {
	if t.do3valid() == false {
		return false
	}
	t.m08ifile.do19close()
	e := new(t0118C8C0fileempty) // v09D9BD74mm.do11malloc(0x8, 0)
	// inline
	e.construct()
	// t.m08ifile = e
	return true
}

type t0117EE84 struct {
	t011737B4
}

func (t *t0117EE84) f00AEC019construct() {
	t.t011737B4.f00A3BAEBconstruct()
}

type t0117EE34 struct {
	t0117EE84
}

func (t *t0117EE34) f00AEBD83construct() {
	t.t0117EE84.f00AEC019construct()
}

// size:0x1C stdcall带ebp帧栈
type t0117EDE4 struct {
	t0117EE34
	m08ozg    *xstring
	m0Cbuf    []uint8
	m10size   int
	m14offset int
	m18valid  bool
}

func (t *t0117EDE4) f00AEBCF0construct(ozg string, buf []uint8, size int) {
	t.t0117EE34.f00AEBD83construct()
	t.m08ozg.f00BADDD0xstring(ozg)
	t.m0Cbuf = buf
	t.m10size = size
	t.m14offset = 0
	t.m18valid = true
}

// f00AEBE00
func (t *t0117EDE4) do9() bool {
	return false
}

// f00AEBF25
func (t *t0117EDE4) do15seek(offset int, whence int) int {
	switch whence {
	case io.SeekStart:
		t.m14offset = offset
	case io.SeekCurrent:
		t.m14offset += offset
	case io.SeekEnd:
		t.m14offset = t.m10size - offset
	}
	return t.m14offset
}

func (t *t0117EDE4) f00AEBCD8seek() {
	t.do15seek(0, io.SeekStart)
}

func (t *t0117EDE4) f00AEB3BBdestruct() {

}

type listNode struct {
	p1    *listNode
	p2    *listNode
	value interface{}
}
type list2 struct {
	unk  uintptr
	head *listNode
	size int
}

func (t *list2) f00AEBBC2listMaxSize() int {
	return 0x3FFFFFFF
}

// stdcall带ebp帧栈
type t22 struct {
	m00     *t23
	m04head *listNode
}

func (t *t22) f00AEB64B(x *listNode, y *gfxcss) *t22 {
	// t.f00AEBA3B(x, y)
	func(x *listNode, y *gfxcss) {
		// t.f00AEBA88()
		func() {
			// t.f00AEBA9C
			func() {
				// t.f0042DBD0()
				func() {
					t.m00 = nil
				}()
			}()
		}()
		t.m04head = x
		if y == nil {
			f00DE84C8()
		}
		// t.f0042DBE0(y)
		func(x *gfxcss) {
			t.m00 = x.m00
		}(y)
	}(x, y)
	return t
}

func (t *t22) f00AEBA79() *listNode {
	return t.m04head
}

type t23 struct{}

func f00AEB9F8newnode(num int) *listNode {
	// f00AEBBFA(size, 0)
	return func(num, zero int) *listNode {
		if num > 0 {
			if ^uint(0)/uint(num) < 12 {
				// ebpC.f00407570(0)
				// f00DE84E3(&ebpC, v012DA28C)
			}
		}
		return new(listNode) // return f00DE852Fnew(num * 12)
	}(num, 0)
}

func f00AEBA0Fnodeassign(nodepdst **listNode, nodepsrc **listNode) {
	// f00AEBC40(x, root)
	func(nodepdst **listNode, nodepsrc **listNode) {
		if nodepdst == nil {
			return
		}
		*nodepdst = *nodepsrc
	}(nodepdst, nodepsrc)
}

func f00AEBAF9nodevalue(dst *interface{}, src *interface{}) {
	*dst = *src
}

func f00AEB425getp1(node *listNode) **listNode {
	return &node.p1
}
func f00AEB763getp2(node *listNode) **listNode {
	return &node.p2
}
func f00AEB76Egetvalue(node *listNode) *interface{} {
	return &node.value
}

// stdcall带ebp帧栈
type gfxcss struct {
	m00     *t23
	m08     *t01187D5C
	m0C     *t0117437C
	m10list list2 // m14head *listNode
}

func (t *gfxcss) f00AEB31Fappend(buf []uint8) {
	var ebp10 t22
	// t.f00AEB300(&ebp10)
	func(x *t22) {
		x.f00AEB64B(t.m10list.head, t)
	}(&ebp10)

	// t.f00AEB42Dappend(ebp10, buf)
	func(x t22, buf []uint8) {
		ebp8head := x.f00AEBA79()
		// ebp4 := t.f00AEB805nodeappend(ebp8head, *f00AEB763getp2(ebp8head), buf)
		ebp4node := func(node *listNode, p2 *listNode, buf []uint8) *listNode {
			// ebp1C := t
			ebp18node := f00AEB9F8newnode(1)
			ebp14 := 0
			// ebp20nodep := f00AEB425getp1(ebp18node)
			// f00AEBA0Fnodeassign(ebp20nodep, &node)
			ebp18node.p1 = node
			ebp14++
			// ebp24nodep := f00AEB763getp2(ebp18node)
			// f00AEBA0Fnodeassign(ebp24nodep, &p2)
			ebp18node.p2 = p2
			ebp14++
			// ebp28valuep := f00AEB76Egetvalue(ebp18node)
			// f00AEBAF9nodevalue(ebp28, buf)
			ebp18node.value = buf
			return ebp18node
		}(ebp8head, ebp8head.p2, buf)

		// t.f00AEB91ClistSizeInc(1)
		func(delta int) {
			// 0x48局部变量
			// ebp54 := t
			// t.f00AEBAB0listMaxSize()
			if func() int { return t.m10list.f00AEBBC2listMaxSize() }()-t.m10list.size < delta {
				// list<T> too long
			}
			t.m10list.size += delta
		}(1)

		// ebp4node被指向
		*f00AEB763getp2(ebp8head) = ebp4node
		*f00AEB425getp1(*f00AEB763getp2(ebp4node)) = ebp4node
	}(ebp10, buf)
}

type t0117EDD0gfxFileOpener struct {
	m0C gfxcss
}

func (t *t0117EDD0gfxFileOpener) f00BEFFC0open(file string, mode, y int) ifile {
	w := new(t0118C910fileio) // v09D9BD74.do11malloc(0x0C, 0)
	var ebp4 *xstring
	ebp4.f00BADDD0xstring(file)
	w.f00BEFF80construct(ebp4, mode, y)
	ebp4.destruct()
	return w
}

// f00AEAD21 stdcall且有ebp帧栈，业务代码
func (t *t0117EDD0gfxFileOpener) do2load(file string, mode, y int) *t0117EDE4 {
	// 0xA0局部变量
	// ebp6C := t
	var ebp28 *xstring
	ebp28.f00BADDD0xstring(file)
	switch {
	case f00DE92E0strstr(ebp28.f00A3AF9Ecstr(), ".dds") != nil: // f009235EC(ebp28.f00A3AF9Ecstr(), ".dds")
		// var ebp44 *xstring
		// ebp70 := ebp28.f00BAE080(ebp44, 0, ebp28.f00BAC0B0()-4)
		// ebp74 := ebp70
		// ebp28.f00BACE40(ebp74)
		// ebp44.f00A3AF16destruct()
		// ebp28.f00A4FF11("ozd")
	case f00DE92E0strstr(ebp28.f00A3AF9Ecstr(), ".gfx") != nil:
		// var ebp44 *xstring
		// ebp70 := ebp28.f00BAE080(ebp44, 0, ebp28.f00BAC0B0()-4)
		// ebp74 := ebp70
		// ebp28.f00BACE40(ebp74)
		// ebp44.f00A3AF16destruct()
		// ebp28.f00A4FF11(".ozg")
	case f00DE92E0strstr(ebp28.f00A3AF9Ecstr(), ".ozp") != nil:
		ebp38file := t.f00BEFFC0open(ebp28.f00A3AF9Ecstr(), mode, y) // ebp38file := f00AEC036assign(f)
		if ebp38file.do3valid() == false {
			ebp38file.f00BAE460()
			// ebp28.f00A3AF16destruct()
			return nil
		}
		ebp2Cbuf := make([]uint8, ebp38file.do7getSize()) // ebp2Cbuf := f00A3BA24malloc(ebp38file.do7getSize())
		ebp3Csize := ebp38file.do7getSize()
		ebp38file.do11read(ebp2Cbuf, ebp3Csize)

		// c++搞个切片好难
		ebp40size := ebp3Csize - 4
		ebp30buf := make([]uint8, ebp40size) // ebp30buf := f00A3BA24malloc(ebp40size)
		f00DE7C90memcpy(ebp30buf, ebp2Cbuf, ebp40size)

		// unmarshal
		ebp54 := new(t0117EDE4) // f00A3BA10newobject(0x1C)
		ebp54.f00AEBCF0construct(ebp28.f00A3AF9Ecstr(), ebp30buf, ebp40size)
		ebp90 := ebp54
		ebp50 := ebp90
		ebp34 := ebp50 // ebp34.f00AEB397assign(ebp50)
		ebp34.f00AEBCD8seek()

		ebp38file.do19close()
		ebp38file.f00BAE460()
		// f00A3AF52free(ebp2Cbuf)
		t.m0C.f00AEB31Fappend(ebp30buf)
		ebp34.f00BAE460()
		ebp38file.f00BAE460()
		// ebp28.f00A3AF16destruct()
		return ebp34
	}

	// 0x00AEB01F
	ebp1Cfile := t.f00BEFFC0open(ebp28.f00A3AF9Ecstr(), mode, y) // ebp1Cfile.f00AEC036assign(f)
	if ebp1Cfile.do3valid() == false {
		// ...
		return nil
	}
	ebp10bufsrc := make([]uint8, ebp1Cfile.do7getSize()) // ebp10 := f00A3BA24malloc(ebp1Cfile.do7getSize())
	ebp20size := ebp1Cfile.do7getSize()
	ebp1Cfile.do11read(ebp10bufsrc, ebp20size)
	// decode
	ebp24size := f00658C4Ddec(nil, ebp10bufsrc, ebp20size) // 2A2DF->2A2BD
	ebp14bufdst := make([]uint8, ebp24size)                // ebp14 := f00A3BA24malloc(ebp24size)
	f00658C4Ddec(ebp14bufdst, ebp10bufsrc, ebp20size)
	// unmarshal
	ebp64 := new(t0117EDE4) // f00A3BA10newobject(0x1C)
	ebp64.f00AEBCF0construct(ebp28.f00A3AF9Ecstr(), ebp14bufdst, ebp24size)
	ebpA8 := ebp64
	ebp60 := ebpA8
	ebp18 := ebp60 // ebp18.f00AEB397assign(ebp60)
	ebp18.f00AEBCD8seek()
	ebp1Cfile.do19close()
	ebp1Cfile.f00BAE460() // 引用计数-1，析构函数里面再做一次-1然后释放内存
	// f00A3AF52(ebp10bufsrc) // free
	t.m0C.f00AEB31Fappend(ebp14bufdst)
	// ebp18.f00AEB3BBdestruct()
	// ebp1Cfile.f00A504C8destruct()
	// ebp28.f00A3AF16destruct()
	return ebp18
}

// f00BF00C0
func (t *t0117EDD0gfxFileOpener) do4load(file string, i i011737C4, mode, unk int) *t0117EDE4 {
	f := t.do2load(file, mode, unk)
	if f != nil && f.do9() == false {
		return f
	}
	// ...
	return nil
}

type t0118CE40 struct {
	m08gfxFileOpener *t0117EDD0gfxFileOpener
}

// size:0x54 0x0E9A25C0
type t01196128gfxLoader struct {
	t011737C4base
	m08 *t0118CE40
	m0C *t0118CBD8
}

func (t *t01196128gfxLoader) f00D05DA0construct() {}

func (t *t01196128gfxLoader) f00D03CA0load(file string, unk int) *t0117EDE4 {
	if t.m08.m08gfxFileOpener == nil {
		// t.m0C.f00BB0310(t, "GFxLoader failed to open %s, GFxFileOpener not installed\n", file)
		return nil
	}
	return t.m08.m08gfxFileOpener.do4load(file, t.m0C, 0x21, 0x1B6)
}

type t0118CB68 struct {
	t011737C4base
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

func (t *t0118CB68) f00BF5350construct(x *t01187A80, y *int) {
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
		w := new(t0118CBD8)
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

// "./Data/Interface/GFx/MainFrame.ozg", x, 0, 0x80
func (t *t0118CB68) f00BF55E0(ozg string, x *attr, y int, z int) bool {
	// 0x2B8局部变量
	w := new(t01196128gfxLoader)
	w.f00D05DA0construct()
	// f00BADDD0xstring(ozg)
	// ...
	// 0x00BF5796
	if w.f00D03CA0load(ozg, 0) != nil {

	}

	return false
}

type window01187A44 struct {
	t011737C4base
	m08 int
	m0C int
}

type window01187920 struct {
	t011737C4base
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

// stdcall无ebp帧栈
type windowManager1 struct {
	m00vtabptr []uintptr
	m04        *t0118CB68
	m08        *t01187A80
	m0C        int
}

func (t *windowManager1) do3(x int, window interface{}) {
	// t.m04.do3()
}

func (t *windowManager1) f00BB0920(x int, ws ...i011737C4) {
	{
		// v09D9BD74mm.do11malloc(0x14, 0)
		w := new(t01187A80)
		w.f00BB4930construct(false)
		t.m08 = w
		t.m0C = x
	}
	{
		// v09D9BD74mm.do11malloc(0x40, 0)
		w := new(t0118CB68)
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

func (t *windowManager1) f00BAFFD0(ozg string, x *attr, y int, z int) {
	t.m04.f00BF55E0(ozg, x, y, z)
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
	do2(wm *windowManager1, x int, ozgfile *xstring, lang *xstring, unk1, unk2, unk3, unk4 int) bool
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
	barer   // m00vtabptr []uintptr
	m04     struct{}
	m08     struct{}
	m0C     attr
	m34     int
	m38ozg  *xstring // "/Data/Interface/GFx/MainFrame.ozg"
	m3Cozg  *xstring // "/Data/Interface/GFx/MainFrame.ozg"
	m44name *xstring // "g_mcMainFrame"
	m49     bool
	m4B     bool
}

func (t *windowgame0117373C) f00A392A6construct() {
	// t.m00vtabptr = v0117373C[:]
	// t.m04.f00A3AE0E(nil)
	// t.m08.f00A3AE91(nil)
}

// f00A39479
func (t *windowgame0117373C) do2(wm *windowManager1, x int, ozgfile *xstring, lang *xstring, unk1, unk2, unk3, unk4 int) bool {
	// 0x4C4局部变量
	// ebp484 := t
	t.m38ozg = ozgfile // t.m38ozg.f00BACE40(&ozgfile)
	var ebp414 [1024]uint8
	f00DF30EFstrcpysafe(ebp414[:], len(ebp414), ozgfile.f00A3AF9Ecstr())
	ebp10 := ebp414[:] // ebp10 := f00DF3001(ebp414[:], "\\", &ebp14)
	for ebp10 != nil {
		t.m3Cozg.f00BACD30xstring(string(ebp10))
		ebp10 = nil // ebp10 := f00DF3001(nil, "\\", &ebp14)
	}
	t.m34 = x
	wm.f00BAFFD0(ozgfile.f00A3AF8Acstr(), &t.m0C, 0, 0x80)
	return false
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

func (t *windowgameCaution) do1(x bool) {}
func (t *windowgameCaution) do3()       {}
func (t *windowgameCaution) do4() bool  { return false }
func (t *windowgameCaution) do6() bool  { return false }
func (t *windowgameCaution) do7() bool  { return false }
func (t *windowgameCaution) do16()      {}

type windowgame0117A544 struct {
	windowgame0117373C
}

// mainFrame or dash size:0x340
type windowgameMainFrame struct {
	windowgame0117A544
	m88hp    int
	m8ChpMax int
	m90      bool
	m94mp    int
	m98mpMax int
	m9Csd    int
	mA0sdMax int
	mA4ag    int
	mA8agMax int
	mAC      int
	mB0      int
	mB4      bool
	mB5      bool
	mB6      bool
	mB8      int
	mBC      struct{ data [24]uint8 }
	mD4      struct{ data [8]uint8 }
	mDC      struct{ data [8]uint8 }
	mE4      struct{ data [40]uint8 }
	m10C     struct{ data [40]uint8 }
	m134     struct{ data [160]uint8 }
	m1D4     struct{ data [160]uint8 }
	m280     int
	m284     struct{}
	m2A0     struct{}
	m328     int
	m330     bool
	m331     bool
	m332     bool
	m333     bool
	m334     bool
}

func (t *windowgameMainFrame) f00AA9021construct() {
	t.windowgame0117373C.f00A392A6construct()
	// t.mBC.f00938EFD()
	// t.mD4.f00A3AFFF()
	// t.mDC.f00A3AFFF()
	// f0043D7C1(&t.mE4, 8, 5, f00A3AFFF)   // 8*5
	// f0043D7C1(&t.m10C, 8, 5, f00A3AFFF)  // 8*5
	// f0043D7C1(&t.m134, 8, 20, f00A3AFFF) // 8*20
	// f0043D7C1(&t.m1D4, 8, 20, f00A3AFFF) // 8*20
	// t.m284.f0052AA7D()
	// t.m2A0.f0049FA13(3, 1)
	// ebp14 := new(windowgame0117A590) // 0x14
	// ebp14.f00AAD8C4construct(t)
	// ebp1C := ebp14
	// ebp10 := ebp1C
	// t.m30.f00A3ADBC(ebp10)
}

func (t *windowgameMainFrame) do1(x bool) {}

func (t *windowgameMainFrame) do3() {}

// f00AA9458
func (t *windowgameMainFrame) do4() bool {
	return false
}

func (t *windowgameMainFrame) f00AAA14Ehpmp() {
	var ebp8hpMax int
	var ebpCmpMax int
	var ebp14hp int
	var ebp10mp int
	if f0059D4F6bit4(v0805BBACself.m13) {
		var ebp1C, ebp20 int
		ebp8hpMax = int(v08C88E58hpMax)
		if v086105ECobject.m122hp >= 0 {
			ebp1C = int(v086105ECobject.m122hp)
		} else {
			ebp1C = 0
		}
		ebp14hp = ebp1C
		ebpCmpMax = int(v08C88E5AmpMax)
		if v086105ECobject.m124mp >= 0 {
			ebp20 = int(v086105ECobject.m124mp)
		} else {
			ebp20 = 0
		}
		ebp10mp = ebp20
	} else {
		var ebp24, ebp28 int
		ebp8hpMax = int(v086105ECobject.m126hpMax) // 0x1B6
		if v086105ECobject.m122hp >= 0 {
			ebp24 = int(v086105ECobject.m122hp)
		} else {
			ebp24 = 0
		}
		ebp14hp = ebp24                            // 0x144
		ebpCmpMax = int(v086105ECobject.m128mpMax) // 0x67
		if v086105ECobject.m124mp >= 0 {
			ebp28 = int(v086105ECobject.m124mp)
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
	var ebp8sdMax int
	var ebp4sd int
	if f0059D4F6bit4(v0805BBACself.m13) {
		var ebp10, ebp12 uint16
		if v08C88E5CsdMax >= 1 {
			ebp10 = v08C88E5CsdMax
		} else {
			ebp10 = 1
		}
		ebp8sdMax = int(ebp10)
		if ebp8sdMax >= int(v086105ECobject.m12Asd) {
			ebp12 = v086105ECobject.m12Asd
		} else {
			ebp12 = uint16(ebp8sdMax)
		}
		ebp4sd = int(ebp12)
	} else {
		var ebp18, ebp1A uint16
		if v086105ECobject.m12CsdMax >= 1 {
			ebp18 = v086105ECobject.m12CsdMax
		} else {
			ebp18 = 1
		}
		ebp8sdMax = int(ebp18) // 0x541
		if ebp8sdMax >= int(v086105ECobject.m12Asd) {
			ebp1A = v086105ECobject.m12Asd
		} else {
			ebp1A = uint16(ebp8sdMax)
		}
		ebp4sd = int(ebp1A) // 0x541
	}
	if t.m9Csd != ebp4sd || t.mA0sdMax != ebp8sdMax {
		t.m9Csd = ebp4sd
		t.mA0sdMax = ebp8sdMax
		f00A3A4F2(t, "SetSD", "%d %d", ebp4sd, ebp8sdMax)
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

func f00A3A4F2(iw iwindowgame0117373C, attr string, format string, param ...interface{}) {
	// 0x408局部变量
	// reflect
	w := iw.(*windowgameMainFrame)
	// ebp4 := param
	var ebp408 [1024]uint8 // "_root.g_mcMainFrame.SetSD"
	f00DF30EFstrcpysafe(ebp408[:], 1024, "_root.")
	f00DECB2Estrcatsafe(ebp408[:], 1024, w.m44name.f00A3AF9Ecstr()) // "g_mcMainFrame"
	f00DECB2Estrcatsafe(ebp408[:], 1024, ".")
	f00DECB2Estrcatsafe(ebp408[:], 1024, attr) // "SetSD"
	// if w.m08.m00 != nil {                      // if b.m08.f00A3BAB6() != nil {
	// 	w.m08.m00.f00BB1E00(ebp408[:], format, ebp4...) // b.m08.f00A3AF08().f00BB1E00(ebp408[:], format, ebp4)
	// }

}

type t01190C50 struct {
}
