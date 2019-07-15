package main

import (
	"os"
	"unsafe"
)

var _012EC83C = [16]uint8{
	0x9B, 0xA7, 0x08, 0x3F,
	0x87, 0xC2, 0x5C, 0xE2,
	0xb9, 0x7A, 0xD2, 0x93,
	0xBF, 0xA7, 0xDE, 0x20,
}

func _00B385D6(buf []uint8, size uint32, flag uint32) {
	// 20个局部变量

	var index uint32 = 0 //ebp-14h
	var key = [16]uint8{ //ebp-10h
		0x7C,
		0xBD,
		0x81,
		0x9F,

		0x3D,
		0x93,
		0xE2,
		0x56,

		0x2A,
		0x73,
		0xD2,
		0x3E,

		0xF2,
		0x83,
		0x95,
		0xBF,
	}

	for index < size {
		// buf = buf[index:]
		// el := uint32(buf[0])
		// el |= uint32(key[index+flag])
		// buf[0] = uint8(el)

		buf[index] ^= key[index+flag]
		index++
	}
}

type t1000 struct {
	f00h  uint32
	f04h  int32
	f08h  [260]uint8
	f10Ch uint32
}

// hError, 0x0131,9E08
var _01319E08_log = t1000{}

// var _01319E08_log = t1000{
// 	f00h: 0x01180F50,
// 	f04h: 0xE0, //hFile
// 	// f08h: "MuError.log"
// 	f10Ch: 0x0F,
// }
func init() {
	file := &_01319E08_log
	// OpenFile
	file.f04h = 0xE0
	copy(file.f08h[:], path)
	file.f10Ch = 0x0F
}

func (p *t1000) _00B38AE4(format string, a ...interface{}) []uint8 {
	// 1024+4字节的局部变量

	// ebp_404 := a // go里面可变切片a...貌似没法落在局部变量中，只能直接作为参数
	var buf [1024]uint8
	buf[0] = 0

	_00DE8100(buf[1:], 0, 0x3FF)    // buf清零操作
	_00DF0805(buf[:], format, a...) // 把logconf字符串写到切片里，这个在golang里面很简单

	// 将buf编码后再写进文件
	// _00B38A8D
	// 它只要使用ecx传参数，我就认为作者使用的是匿名调用
	func(buf []uint8) {
		var ebp_4 uint32 // dwNumBytes
		var ebp_8 *t1000 = p
		if p.f04h == -1 {
			return
		}

		// 验证
		// _00DE7C00
		len := func(buf []uint8) uint32 {
			// 验证buf底层数组地址
			// 验证buf底层数组元素
			return uint32(len(buf))
		}(buf)

		// 写文件
		// _00B38A4D
		func(hFile int32, buf []uint8, len uint32, pdwNumBytes *uint32, pOverlapped uint32) {
			// 这个变量是push ecx得到的
			// 局部变量<=8个字节，c编译器使用push，指令数比sub指令少，但性能不行
			var ebp_4 *t1000 = p

			// _00B38653 疑似编码函数
			func(buf []uint8, len uint32, flag uint32) uint32 {
				// 16个字节的局部变量
				var ebp_4 uint32
				var ebp_8 []uint8 = buf
				var ebp_C uint32 = len
				var ebp_10 uint32
				var eax uint32 = (0x10 - flag) & 0x8000000F
				if eax < 0 { // jns, Jump if not sign
					eax--
					eax |= 0xFFFFFFF0
					eax++
				}
				if eax < len {
					eax = (0x10 - flag) & 0x8000000F
					if eax < 0 { // jns, Jump if not sign
						eax--
						eax |= 0xFFFFFFF0
						eax++
					}
					ebp_10 = eax
				} else {
					ebp_10 = ebp_C
				}

				// 编码1个字节？
				ebp_4 = ebp_10                // 1
				_00B385D6(ebp_8, ebp_4, flag) // (0x0018DF84, 1, 15)
				ebp_8 = ebp_8[ebp_4:]
				ebp_C -= ebp_4
				if ebp_C <= 0 { // jg, Jump if greater
					return ebp_4 + ebp_10
				}

				// 编码写16个字节
				for ebp_C >= 0x10 { // jl, Jump if less
					ebp_4 = 0x10
					_00B385D6(ebp_8, ebp_4, 0)
					ebp_8 = ebp_8[ebp_4:]
					ebp_C -= ebp_4
				}

				// 编码剩下的
				ebp_4 = ebp_C
				_00B385D6(ebp_8, ebp_4, 0)
				return ebp_4
			}(buf, len, ebp_4.f10Ch)

			WriteFile(hFile, buf, len, pdwNumBytes, pOverlapped)
		}(ebp_8.f04h, buf, len, &ebp_4, 0)

		//
		if ebp_4 == 0 {
			nRet := CloseHandle(ebp_8.f04h)
			// _00B38781
			func(path string) {
				// 重新设置ebp_8的hFile
			}(string(ebp_8.f08h[:]))
		}
		return
	}(buf[:])

	return buf[:]
}

func (p *t1000) _00B38D19_cut() {
	p._00B38AE4("-----------------\r\n")
}

func (p *t1000) _00B38D31_begin() {
	p._00B38AE4("###### Log Begin ######\r\n")
}

func (p *t1000) _00B38D49(x uint32) {
	// 系统时间
}

func (p *t1000) _00B38E3C() {
	p._00B38AE4("<OpenGL information>\r\n>")
	p._00B38AE4("Vendor\t\t: %s\r\n", glGetString(0x1F00))
	p._00B38AE4("Render\t\t: %s\r\n", glGetString(0x1F01))
	p._00B38AE4("OpenGL version\t: %s\r\n", glGetString(0x1F02))
	var ebp_8 struct {
		data [2]uint32 // ebp-8 and ebp-4
	}
	glGetIntegerv(0xD33, &ebp_8)
	p._00B38AE4("Max Texture size\t: %d x %d\r\n", ebp_8.data[0], ebp_8.data[0])
	glGetIntegerv(0xD3A, &ebp_8)
	p._00B38AE4("Max Viewport size\t: %d x %d\r\n", ebp_8.data[0], ebp_8.data[1])
}

func (p *t1000) _00B3902D() {
	// <Sound card information>\r\n

	// C:\\Windows\\system32\drivers\{...}

	// cut
}

func (p *t1000) _00B38EF4(hWnd uint32) {
	// <IME information>\r\n
}

func _00B4C1B8() bool {
	return false
}

func _00B4C1FF(hWnd uint32) {
	if _01319D68 != 0 {
		// ...
	}
}

type t2000 struct {
	data [64]uint8
}

var _08C8D050_enc t2000
var _08C8D098_dec t2000

func (p *t2000) _00B62CF0_init(path string) {
	return p._00B62EC0(path, 0x1112, 1, 1, 0, 1)
}

func (p *t2000) _00B62D30_init(path string) {
	return p._00B62EC0(path, 0x1112, 1, 0, 1, 1)
}

func (p *t2000) _00B62EC0(path string, begin uint32, x1 uint32, x2 uint32, x3 uint32, x4 uint32) bool {

	var buf struct {
		head [8]uint8
		data [16]uint8
	}

	// var nLen uint32
	// dwDesiredAccess = 1<<31, GENERIC_READ
	// dwShareMode = 1, FILE_SHARE_READ
	// psa = 0, default security and cannot be inherited by any child
	// dwCreationDisposition = 3, OPEN_EXISTING
	// dwFlagsAndAttributes = 0x80, FILE_ATTRIBUTE_NORMAL
	// dwFileTemplate = 0, When opening an existing file, CreateFile ignores this parameter
	// hFile := CreateFile(path, 1<<31, 1, 0, 3, 0x80, 0)
	// ReadFile(hFile, buf[:], 6, &nLen, 0)
	file, error := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	nLen, err := file.Read(buf.head[:6])
	if err != nil {
		return false
	}

	// compare begin
	if uint16(begin) != *(*uint16)(unsafe.Pointer(&buf.head[0])) {
		return false
	}

	// compare len
	if *(*uint32)(unsafe.Pointer(&buf.head[2])) != (x1+x2+x3+x4)<<4+6 {
		return false
	}

	if x1 != 0 {
		file.Read(buf.data[:])
		index := 0
		enc := _08C8D050_enc.data[:]
		for {
			enc[index] = _012EC83C[index] ^ buf.data[index]
			index++
			if index >= 16 {
				break
			}
		}
	}

	if x2 != 0 {
		file.Read(buf.data[:])
		index := 0
		enc := _08C8D050_enc.data[16:]
		for {
			enc[index] = _012EC83C[index] ^ buf.data[index]
			index++
			if index >= 16 {
				break
			}
		}
	}

	if x3 != 0 {
		file.Read(buf.data[:])
		index := 0
		enc := _08C8D050_enc.data[32:]
		for {
			enc[index] = _012EC83C[index] ^ buf.data[index]
			index++
			if index >= 16 {
				break
			}
		}
	}

	if x4 != 0 {
		file.Read(buf.data[:])
		index := 0
		enc := _08C8D050_enc.data[48:]
		for {
			enc[index] = _012EC83C[index] ^ buf.data[index]
			index++
			if index >= 16 {
				break
			}
		}
	}

	return true
}
