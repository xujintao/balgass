package main

type blockCipher interface {
	do2(key []uint8, l int)
}

type t01186D00 struct {
	m00vtabptr []uintptr
}

func (t *t01186D00) construct() {
	// t.m00vtabptr = v01186D00[:]
}

// size:0x20
type t01186D2C struct {
	m00vtabptr []uintptr
	m04        t01186D00
	m18        int
	m1C        int
}

func (t *t01186D2C) f00BA03C0construct() {
	// t.m04.f00401820(1)
	t.m18 = 0
	t.m1C = 0
	// t.m00vtabptr = v01186D2C[:]
	// inline
	t.m04.construct()
}

// f00401890
func (t *t01186D2C) do7(buf []uint8, l int, x interface{}) {
	// f00401A10()
	func() {

	}()
}

type t01186D70 struct {
	m00vtabptr []uintptr
}

func (t *t01186D70) construct() {
	// t.m00vtabptr = v01186D70[:]
}

// size:0x24
type t01186D9C struct {
	m00vtabptr []uintptr
	m04        t01186D70
	m18        int
	m1C        int
}

func (t *t01186D9C) f00BA0750construct() {
	// t.m04.f00401820(1)
	t.m18 = 0
	t.m1C = 0
	// t.m00vtabptr = v01186D9C[:]
	// inline
	t.m04.construct()
}

func (t *t01186D9C) do7(buf []uint8, l int, x interface{}) {

}

// size:0x4C
type t01186EBC struct {
	m00vtabptr []uintptr
	m08        t01186D2C
	m28        t01186D9C
}

func (t *t01186EBC) f00BA0C80construct() {
	// t.m00vtabptr = v01186EBC[:]
	t.m08.f00BA03C0construct()
	t.m28.f00BA0750construct()
}

var v0130714C int

// f00BA0CE0
func (t *t01186EBC) do2(buf []uint8, l int) {
	t.m08.do7(buf, 0x10, &v0130714C)
	t.m28.do7(buf, 0x10, &v0130714C)
}

// ------------------------------------------------------------
type bmdCipher struct {
	m00bc blockCipher
}

func (t *bmdCipher) f00B99D60construct() {
	t.m00bc = nil
}

func (t *bmdCipher) f00B99D70destruct() {

}

func (t *bmdCipher) f00BA1120expandKey(flag int, key []uint8, l int) {
	var bc blockCipher
	if t.m00bc != nil {
		// ...
		t.m00bc = nil
	}
	switch flag & 7 {
	case 0: // ecb
	case 1: // cbc
	case 2: // cfb
	case 3: // ctr
		c := new(t01186EBC) // f00DE852Fnew(t01186EBC)
		c.f00BA0C80construct()
		bc = c
	case 4: // ofb
	case 5:
	case 6:
	case 7:
	}
	t.m00bc = bc
	bc.do2(key, l)
}

func (t *bmdCipher) f00B99E20(size int) int {
	return 0
}

func (t *bmdCipher) f00B99D90enc(dst []uint8, src []uint8, len int) {

}

func (t *bmdCipher) f00B99DC0dec(dst []uint8, src []uint8, len int) {

}

func f00658A9Cenc(dst []uint8, src []uint8, size int) int {
	// 0x68局部变量
	if dst == nil {
		return size + 0x22
	}
	var ebp6C [34]uint8
	ebp6C[0] = uint8(f00DE8AADrand() % 11)
	ebp6C[1] = uint8(f00DE8AADrand() % 7)
	// c++的ebp6C[2:34]可以保证是随机的
	f00DE7C90memcpy(dst, ebp6C[:], 34)
	f00DE7C90memcpy(dst[34:], src, size)

	// 使用key1加密
	var ebp14 bmdCipher
	ebp14.f00B99D60construct()
	ebp14.f00BA1120expandKey(int(ebp6C[0]), ebp6C[2:34], 32)
	ebp14.f00B99D90enc(dst[34:], dst[34:], ebp14.f00B99E20(size))
	ebp14.f00B99E20(size)

	// 使用key2强化加密
	var ebp44key [33]uint8
	copy(ebp44key[:], "webzen#@!01webzen#@!01webzen#@!0")
	ebp44key[32] = 0
	var ebp1C bmdCipher
	ebp1C.f00B99D60construct()
	ebp1C.f00BA1120expandKey(int(ebp6C[1]), ebp44key[:], f00DE7C00strlen(ebp44key[:]))
	ebp20 := ebp1C.f00B99E20(1024)
	if size > ebp20 {
		// 使用key2对头部的1024字节加密
		ebp48 := dst[2:]
		ebp1C.f00B99D90enc(ebp48, ebp48, ebp20)
		// 使用key2对尾部的1024字节加密
		ebp48 = dst[34+size-ebp20:]
		ebp1C.f00B99D90enc(ebp48, ebp48, ebp20)
	}
	if size > ebp20*4 {
		// 使用key2对中间的1024字节加密
		ebp48 := dst[2+size/2:]
		ebp1C.f00B99D90enc(ebp48, ebp48, ebp20)
	}
	ebp1C.f00B99D70destruct()
	ebp14.f00B99D70destruct()
	return size
}

// cdecl带ebp帧栈
// bmd ozg
func f00658C4Ddec(dst []uint8, src []uint8, size int) int {
	// 0x6C局部变量
	if dst == nil {
		return size - 0x22
	}
	ebp70mod1 := src[0] // 8
	ebp6Fmod2 := src[1] // 3

	// 强化解密
	var ebp48key [33]uint8
	copy(ebp48key[:], "webzen#@!01webzen#@!01webzen#@!0")
	ebp48key[32] = 0
	var ebp20 bmdCipher
	ebp20.f00B99D60construct()
	ebp20.f00BA1120expandKey(int(ebp6Fmod2), ebp48key[:], f00DE7C00strlen(ebp48key[:]))
	ebp24 := ebp20.f00B99E20(1024)
	ebp10size := size - 0x22
	if ebp10size > ebp24*4 {
		// 使用key2对中间的1024字节解密
		ebp4C := src[2+ebp10size/2:]
		ebp20.f00B99DC0dec(ebp4C, ebp4C, ebp24)
	}
	if ebp10size > ebp24 {
		// 使用key2对尾部的1024字节解密
		ebp4C := src[size-ebp24:]
		ebp20.f00B99DC0dec(ebp4C, ebp4C, ebp24)
		// 使用key2对头部的1024字节解密
		ebp4C = src[2:]
		ebp20.f00B99DC0dec(ebp4C, ebp4C, ebp24)
	}

	// 使用key1解密
	var ebp6Ekey [32]uint8
	f00DE7C90memcpy(ebp6Ekey[:], src[2:], 32)
	f00DE7C90memcpy(dst, src[34:], ebp10size)
	var ebp18 bmdCipher
	ebp18.f00B99D60construct()
	ebp18.f00BA1120expandKey(int(ebp70mod1), ebp6Ekey[:], 32)
	ebp18.f00B99DC0dec(dst, dst, ebp18.f00B99E20(ebp10size))
	ebp18.f00B99D70destruct()
	ebp20.f00B99D70destruct()
	return size
}
