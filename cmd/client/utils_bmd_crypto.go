package main

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

// f00BA04D0
func (t *t01186D2C) do5(len int) int {
	if len >= 255 {
		return 255
	}
	return len
}

// f00415000
func (t *t01186D2C) do6IsValidKeyLength(len int) bool {
	return len == t.do5(len)
}

// f00401890
func (t *t01186D2C) do7SetKey(key []uint8, len int, x interface{}) {
	// t.f00401A10ThrowIfInvalidKeyLength(len)
	func() {
		if t.do6IsValidKeyLength(len) == false {
			// f004968E0InvalidKeyLength(t.do14GetAlgorithm().do3AlgorithmName(), len)
			// f00DE84E3abort()
		}
	}()
	t.do15UncheckedSetKey(key, len, x)
}

func (t *t01186D2C) do15UncheckedSetKey(key []uint8, len int, x interface{}) {

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

func (t *t01186D9C) do7SetKey(buf []uint8, len int, x interface{}) {

}

type blockCipher interface {
	do2setkey(key []uint8, l int)
	do3encrypt(src []uint8, len int, dst []uint8) int
	do4decrypt(src []uint8, len int, dst []uint8) int
	do5blocksize() int
}

// size:0x4C
type t01186EBCrc5 struct {
	m00vtabptr   []uintptr
	m08encryptor t01186D2C
	m28decryptor t01186D9C
}

func (t *t01186EBCrc5) f00BA0C80construct() {
	// t.m00vtabptr = v01186EBC[:]
	t.m08encryptor.f00BA03C0construct()
	t.m28decryptor.f00BA0750construct()
}

var v0130714C int

// f00BA0CE0
func (t *t01186EBCrc5) do2setkey(key []uint8, len int) {
	// t.m08encryptor.do7SetKey(key, t.KEYLENGTH, &v0130714C)
	// t.m28decryptor.do7SetKey(key, t.KEYLENGTH, &v0130714C)
}

// f00BA0D40
func (t *t01186EBCrc5) do3encrypt(src []uint8, len int, dst []uint8) int {
	if dst == nil || len <= 0 {
		return len
	}
	len /= 8
	for len > 0 {
		// t.m08encryptor.do3ProcessAndXorBlock(src, nil, dst)
		src = src[8:]
		dst = dst[8:]
		len--
	}
	return len
}

// f00BA0D90
func (t *t01186EBCrc5) do4decrypt(src []uint8, len int, dst []uint8) int {
	if dst == nil || len <= 0 {
		return len
	}
	len /= 8
	for len > 0 {
		// t.m28decryptor.do3ProcessAndXorBlock(src, nil, dst)
		src = src[8:]
		dst = dst[8:]
		len--
	}
	return len
}

// f00BA0D30
func (t *t01186EBCrc5) do5blocksize() int {
	return 0
	// return t.BLOCKSIZE
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

func (t *bmdCipher) f00BA1120setkey(alg int, key []uint8, l int) {
	var bc blockCipher
	if t.m00bc != nil {
		// ...
		t.m00bc = nil
	}
	switch alg & 7 {
	case 0: // TEA block cipher, 16, 8
	case 1: // 3-Way block cipher, 12, 12
	case 2: // CAST-128 block cipher, 16, 8
	case 3: // RC5 block cipher, 16, 8
		c := new(t01186EBCrc5) // f00DE852Fnew(t01186EBCrc5)
		c.f00BA0C80construct()
		bc = c
	case 4: // RC6 block cipher, 16, 16
	case 5: // MARS block cipher, 16, 16
	case 6: // IDEA block cipher, 16, 8
	case 7: // GOST block cipher, 32, 8
	}
	t.m00bc = bc
	bc.do2setkey(key, l)
}

func (t *bmdCipher) f00B99E20align(size int) int {
	if size <= 0 {
		return 0
	}
	return size - size%t.m00bc.do5blocksize() // 1024 - 1024%8
}

func (t *bmdCipher) f00B99D90enc(dst []uint8, src []uint8, len int) int {
	if t.m00bc == nil {
		return -1
	}
	return t.m00bc.do3encrypt(src, len, dst)
}

func (t *bmdCipher) f00B99DC0dec(dst []uint8, src []uint8, len int) int {
	if t.m00bc == nil {
		return -1
	}
	return t.m00bc.do4decrypt(src, len, dst)
}

// ------------------------------------------------------
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
	ebp14.f00BA1120setkey(int(ebp6C[0]), ebp6C[2:34], 32)
	ebp14.f00B99D90enc(dst[34:], dst[34:], ebp14.f00B99E20align(size))
	ebp14.f00B99E20align(size)

	// 使用key2强化加密
	var ebp44key [33]uint8
	copy(ebp44key[:], "webzen#@!01webzen#@!01webzen#@!0")
	ebp44key[32] = 0
	var ebp1C bmdCipher
	ebp1C.f00B99D60construct()
	ebp1C.f00BA1120setkey(int(ebp6C[1]), ebp44key[:], f00DE7C00strlen(ebp44key[:]))
	ebp20 := ebp1C.f00B99E20align(1024)
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
	ebp70mode1 := src[0] // 8
	ebp6Fmode2 := src[1] // 3

	// 强化解密
	var ebp48key [33]uint8
	copy(ebp48key[:], "webzen#@!01webzen#@!01webzen#@!0")
	ebp48key[32] = 0
	var ebp20 bmdCipher
	ebp20.f00B99D60construct()
	ebp20.f00BA1120setkey(int(ebp6Fmode2), ebp48key[:], f00DE7C00strlen(ebp48key[:]))
	ebp24 := ebp20.f00B99E20align(1024)
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
	ebp18.f00BA1120setkey(int(ebp70mode1), ebp6Ekey[:], 32)
	ebp18.f00B99DC0dec(dst, dst, ebp18.f00B99E20align(ebp10size))
	ebp18.f00B99D70destruct()
	ebp20.f00B99D70destruct()
	return size
}
