package main

import (
	"fmt"
	"unsafe"
)

const (
	_01149E44 uint32 = 0x004075A0
	_012DA28C uint32 = 0
)

var _012F7B90 uint32 = 0x5F709C77
var _012F8738 *t15    // [_012F8738]=0x012F8660
var _012F8558 *uint32 // [_012F8558]=0x0D921A80

var _09D9D9FC uint32
var _09D9DA08 uint8 = 0
var _09D9DA0C uint32 = 1
var _09D9DB88 uint32 = 0xBD6A970F

func _00DE8000_strcpy(dst []uint8, src []uint8) {
	copy(dst, src)
}

func _00DE8010_strcat(dst []uint8, src string) {
	dst = append(dst, src...)
}

func _00DE8100(buf []uint8, value uint8, size uint32) uint32 {
	if size == 0 {
		return 0
	}
	if value == 0 {
		if size >= 100 {
			// 这里对一个全局变量进行判断
			if true {
				// 编译器把这个函数调用优化为jmp
				return _00DFC9DD(buf, value, size)
			}
		}
	}

	// ...
	return 0
}

// _00DE817A_sprintf
func _00DE817A_sprintf(buf []uint8, strfmt string, a ...interface{}) int {
	str := fmt.Sprint(strfmt, a...)
	buf = []uint8(str) // 从堆上复制到data段
	return len(str)
}

func _00DE852F(x uint32) {

	var ebp_C uint32
	for {
		// _00DF0F2F
		entry := func(x uint32) string {
			return ":\r"
		}(x)
		if len(entry) != 0 {
			return
		}

		// // _00DFB084
		// bRet = func(x uint32) bool {
		// 	return false
		// }(x)
		// if !bRet {
		// 	if _09D9DA08&1 == 0 {
		// 		_09D9DA08 |= 1

		// 		// _00415D00
		// 		func() {

		// 			// _00DE85F8
		// 			func() {

		// 			}()
		// 		}()

		// 		// _00DE8BF6
		// 		func(fval func()) {
		// 			// _00DE8BBA
		// 			nRet := func(fval func()) uint32 {
		// 				return 1
		// 			}(fval)
		// 		}(_01148F33)

		// 		// _004075D0
		// 		func(x *uint32) {
		// 			// 能使用到ebp_C
		// 			// _00DE8615
		// 			func(x *uint32) {

		// 			}(x)
		// 			ebp_C = _01149E44
		// 		}(&_09D9D9FC)

		// 		// _00DE84E3
		// 		// func() {

		// 		// }(&_012DA28C)

		// 		// int3
		// 	}
		// }
	}
}

func _00DECD20(x []uint8, strfmt string, y []uint8) int32 {
	return -1
}

// setlocale?
func _00DEE8171(x uint32, lang string) {
	// _00DFD850
	// _00DFC3E9
	// _00DEDA19
	// _00DF9D56
	// _00DFF2B2
	// _00DED9B5
	// _00DEE99F
	// _00DEE556
	// _00DE94F0
	// _00DFF2B2
	// _00DED9D8
	// _00DED91C
	// ...
}

type t14 struct {
	data [124]uint8
	f7Ch uint16
}

type t15 struct {
	data [200]uint8
	fC8h *t14
}

type tlsvalue struct {
	data [68]uint8
	f68h *uint32
	f6Ch *t15
	f70h uint8
}

type t16 struct {
	f00h *t15    // tlsvalue.f6Ch
	f04h *uint32 // tlsvalue.f68h
	f08h *tlsvalue
	f0Ch [4]uint8
}

// _00DECDA2
func _00DECDA2(prevt16p *t16, t16p *t16) {
	t16p.f0Ch[0] = 0
	if prevt16p == nil {
		// _00DFC3E9
		v := func() unsafe.Pointer {
			// _00DFC370
			v := func() unsafe.Pointer {
				errno := GetLastError()
				// _00DFC1FB
				cb := func() func(uint32) unsafe.Pointer {
					const _012F7DC0 uint32 = 25
					v := unsafe.Pointer(TlsGetValue(_012F7DC0))
					if v == nil {
						// _00DFC160
						func(x uint32) {
							// ...
						}(_09D9DB88)
					}
					return v // kernel32.FlsGetValue，这里unsafe.Pointer值能转换成函数类型吗？
				}()
				const _012F7DBC uint32 = 4
				v := cb(_012F7DBC) // FlsGetValue
				if v == nil {
					// ...
				}
				SetLastError(errno)
				return v
			}()
			if v == nil {
				// ...
			}
			return v
		}()
		tlsv := (*tlsvalue)(v)
		t16p.f08h = tlsv      // [0x0D92,07D0]=01BC
		t16p.f00h = tlsv.f6Ch // [0x0D92,083C]=0x012F,8660
		t16p.f04h = tlsv.f68h // [0x0D92,0838]=0x0D92,1A80
		if t16p.f00h != _012F8738 {
			// ...
		}
		if t16p.f04h != _012F8558 {
			// ...
		}
		if b := tlsv.f70h; b == 2 { // [0x0D92,0840]=1 // 这里有问题，应该跳
			tlsv.f70h |= 1 // [0x0D92,0840]=3
			t16p.f0Ch[0] = 1
			return
		}
	}

	t16p.f00h = prevt16p.f00h
	t16p.f04h = prevt16p.f04h
	return
}

// 好复杂的函数
func _00DFCCB0(infop *info, format string, t16p *t16, a ...interface{}) int {
	// 278h字节的局部变量

	// 功能就是 format = fmt.Sprintf(format, a...)
	// 再copy(buf, format)

	var ebp_25c t16 // 0x0018,DCB8
	// ebp-25c
	// ebp-258
	// ebp-254
	// ebp-250

	var ebp_24c *info = infop
	// ebp-248
	// ebp-240
	// ebp-238
	// ebp-234
	// ebp-230
	// var ebp_224 *conf = c
	// ebp-218
	// ebp-210

	var stack uint32 = _012F7B90 ^ 0x0018DF14 // ebp-4 这个是什么意思？

	_00DECDA2(t16p, &ebp_25c)

	if infop == nil || len(format) == 0 {
		// ...
	}
	if infop.f0Ch == 40 {
		// ...
	}

	var ebp_23C string = format[1:] // 拿掉 '>' 字符
	// ...

	var ebp_228 uint32 // 0x0018,DCEC

	// var ebp_211 uint8 = ebp_23C[0] // c语言风格
	for len(ebp_23C) != 0 {
		ebp_23C = ebp_23C[1:]

		if ebp_228 < 0 { // jl(jump less) 有符号跳转
			break
		}

		// _00E036C2
		nRet := func(c uint8) uint32 {
			var ebp_10 t16
			_00DECDA2(&ebp_25c, &ebp_10)
			nRet := uint32(ebp_10.f00h.fC8h.f7Ch & 8000)
			if ebp_10.f0Ch[0]&0xFF != 0 {
				// ...
			}
			return nRet
		}(ebp_23C[0])
		if nRet != 0 {
			// ...
		}

		// _00DFCBD0
		func() {
			// 初始ecx为ebp_24C, esi=&ebp_228
			if ebp_24c.f0Ch == 0x40 || len(ebp_24c.f08h) != 0 {
				ebp_24c.f04h--
				if ebp_24c.f04h >= 0 {
					ebp_24c.f00h[0] = ebp_23C[0]
					ebp_24c.f00h = ebp_24c.f00h[1:]
				} else {
					// _00DFCA6C
				}
				if ebp_23C[0] == ^uint8(0) {
					ebp_228 |= uint32(ebp_23C[0])
					return
				}
			}
			ebp_228++
			return
		}()
	}
	return int(ebp_228)
}

type info struct {
	f00h []uint8
	f04h uint32
	f08h []uint8
	f0Ch uint32
}

func _00DF0787(buf []uint8, format string, x *t16, a ...interface{}) int {
	// 如果接下来会用到ebx，那么先把ebx压栈

	// c里面判断字符串指针变量是否为空
	// 等效为go里面判断string类型变量长度是否为0
	// 或者切片长度是否为0
	if len(format) == 0 || len(buf) == 0 {
		// ...
		return -1 // ?
	}

	i := info{
		f00h: buf,       // ebp-20
		f04h: 0x7FFFFFF, // ebp-1C
		f08h: buf,       // ebp-18
		f0Ch: 0x42,
	}
	cnt := _00DFCCB0(&i, format, x, a...) // 其实就是把logconf字符串copy到buf切片
	// i.f00h = append(i.f00h,0) // golang不需要追加0
	return cnt
}

func _00DF0805(buf []uint8, format string, a ...interface{}) {
	_00DF0787(buf, format, nil, a...)
}

// OEP: 0x00DF,478C
func main() {
	// check pe

	_004D7CE5_run(0x00400000, 0, 0x0B433D25, 10) // call 0x004D,7CE5
}

// 大范围清零操作
func _00DFC986(buf []uint8, size uint32) {
	size >>= 7
	for size != 0 {
		// 它使用了8次 movdqa xmm指令
		for i := range buf[:128] {
			buf[i] = 0
		}
		buf = buf[128:]
		size--
	}
}

// setzero， 16字节对齐
// 某种意义上说，把迭代改为递归也可以有效防止饥饿调度问题
func _00DFC9DD(buf []uint8, value uint8, size uint32) uint32 {
	// 假设eax是0x0018,DF85，然后edx是0x3FF
	// cdq          ;edx:eax组合成为64位
	// mov edi,eax
	// xor edi,edx
	// sub edi,edx
	// and edi,F
	// xor edi,edx
	// sub edi,edx
	// test edi,edi ;edi值为5

	// neg edi
	// add edi,10

	// offset := 16 - uint32(&buf[0])%16
	// offset := 16 - uint32(unsafe.Pointer(&buf))%16
	offset := 16 - (*[3]uint)(unsafe.Pointer(&buf))[0]%16
	if offset == 0 {
		s := size & 0x7F
		if s != size {
			// ...
		}
		size := size - s
		_00DFC986(buf, size)
		if size == 0 {
			return 0
		}
		buf = buf[int(size):]
		for i := range buf {
			buf[i] = 0
		}
	}
	for i := range buf[:offset] {
		buf[i] = 0
	}

	size -= uint32(offset)
	_00DFC9DD(buf[offset:], 0, size)
	return 0
}
