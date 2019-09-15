package main

import (
	"debug/pe"
	"fmt"
	"os"
	"unsafe"

	"github.com/xujintao/balgass/client/win"
)

const (
	v01149E44 uintptr = 0x004075A0
	v012DA28C int     = 0
)

var v012F7B90 uintptr = 0x5F709C77
var v012F8738 *t15    // [v012F8738]=0x012F8660
var v012F8558 *uint32 // [v012F8558]=0x0D921A80

var v09D9D9FC int
var v09D9DA08 uint8
var v09D9DA0C uint32 = 1
var v09D9DB88 uintptr = 0xBD6A970F

func f00DE7C00strlen(str []uint8) {

}

func f00DE7C90memcpy(dst, src []uint8, len int) {

}

func f00DE8000strcpy(dst []uint8, src []uint8) {
	copy(dst, src)
}

func f00DE8010strcat(dst []uint8, src string) {
	dst = append(dst, src...)
}

func f00DE8100memset(buf []uint8, value uint8, size int) uint32 {
	if size == 0 {
		return 0
	}
	if value == 0 {
		if size >= 100 {
			// 这里对一个全局变量进行判断
			if true {
				// 编译器把这个函数调用优化为jmp
				return f00DFC9DD(buf, value, size)
			}
		}
	}

	// ...
	return 0
}

func f00DE817Asprintf(buf []uint8, strfmt string, a ...interface{}) int {
	str := fmt.Sprint(strfmt, a)
	buf = []uint8(str) // 从堆上复制到data段
	return len(str)
}

func f00DE852F(x uint32) unsafe.Pointer {

	// var ebp_C uint32
	for {
		// _00DF0F2F
		entry := func(x uint32) string {
			return ":\r"
		}(x)
		if len(entry) != 0 {
			return nil
		}

		// // _00DFB084
		// bRet = func(x uint32) bool {
		// 	return false
		// }(x)
		// if !bRet {
		// 	if v09D9DA08&1 == 0 {
		// 		v09D9DA08 |= 1

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
		// 		}(f01148F33)

		// 		// _004075D0
		// 		func(x *uint32) {
		// 			// 能使用到ebp_C
		// 			// _00DE8615
		// 			func(x *uint32) {

		// 			}(x)
		// 			ebp_C = _01149E44
		// 		}(&v09D9D9FC)

		// 		// _00DE84E3
		// 		// func() {

		// 		// }(&v012DA28C)

		// 		// int3
		// 	}
		// }
	}
}

// 拿着eax做栈分配，这个骚套路是什么API？目的是什么？
func f00DE8A70chkstk() {}

// 猜测
func f00DE8AADrand() int { return 0 }

func f00DEE871setlocale(category uint32, locale string) {

}

func f00DE92E0strstr(haystack []uint8, needle string) []uint8 {
	return nil
}

func f00DECD20(x []uint8, strfmt string, y []uint8) int32 {
	return -1
}

// setlocale?
func f00DEE8171(x uint32, lang string) {
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

// f00DECDA2
func f00DECDA2(prevt16p *t16, t16p *t16) {
	t16p.f0Ch[0] = 0
	if prevt16p == nil {
		// _00DFC3E9
		v := func() unsafe.Pointer {
			// _00DFC370
			v := func() unsafe.Pointer {
				errno := win.GetLastError()
				// _00DFC1FB
				cb := func() func(uint32) unsafe.Pointer {
					const v012F7DC0 uint32 = 25
					// v := unsafe.Pointer(win.TlsGetValue(v012F7DC0))
					var v interface{}
					if v == nil {
						// _00DFC160
						func(x uintptr) {
							// ...
						}(v09D9DB88)
					}
					return v.(func(uint32) unsafe.Pointer) // kernel32.FlsGetValue，这里unsafe.Pointer值能转换成函数类型吗？
				}()
				const v012F7DBC uint32 = 4
				v := cb(v012F7DBC) // FlsGetValue
				if v == nil {
					// ...
				}
				win.SetLastError(errno)
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
		if t16p.f00h != v012F8738 {
			// ...
		}
		if t16p.f04h != v012F8558 {
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
func f00DFCCB0(infop *info, format string, t16p *t16, a ...interface{}) int {
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

	// var stack uintptr = v012F7B90 ^ 0x0018DF14 // ebp-4 这个是什么意思？

	f00DECDA2(t16p, &ebp_25c)

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
			f00DECDA2(&ebp_25c, &ebp_10)
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

func f00DF0787(buf []uint8, format string, x *t16, a ...interface{}) int {
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
	cnt := f00DFCCB0(&i, format, x, a...) // 其实就是把logconf字符串copy到buf切片
	// i.f00h = append(i.f00h,0) // golang不需要追加0
	return cnt
}

func f00DF0805(buf []uint8, format string, a ...interface{}) {
	f00DF0787(buf, format, nil, a...)
}

var v00DF490A [0x5B]uint8
var v0B2BE8D6 [0x19]uint8
var v0B2BE90B [0x58]uint8

// shell entry point
func f00DF490A() {
	// vp := win.GetProcAddress(win.LoadLibrary("kernel32.dll"), "virtualProtect")
	// var oldProtect uint32
	// vp(0x00DF490A, 0x5B, 0x40, &oldProtect) // PAGE_EXECUTE_READWRITE
	// vp(0x0B2BE8D6, 0x19, 0x40, &oldProtect) // PAGE_EXECUTE_READWRITE
	// jmp 0x0B2BE8D6
}

func f0B2BE8D6() {
	// f00DE7C90memcpy(v00DF490A[:], v0B2BE90B[:], len(v00DF490A))
	// jmp 0x00DF490A
}

func f0B2BE90B() {
	// E8 4A89BDFF ;call 0x006CD259, f006CD259
	// E9 78FEFFFF ;jmp 0x00DF478C
}

var v0082505D uint32           // 0x6A7B9E4D
var v00825070imageBase uintptr // 0x00400000
var v0A327567imageBase uint32 = 0x00400000
var v0A0032DB uint32 = 0x1CA9625E
var v09E2BC68 uint32 = 0x2DBE21DC
var v0A8FE0E0 uint32 = 0x12A225EC
var v0A746588 uint32 = 0xF7598681
var v0AAB950C uint32 = 0x1000

type crcdata struct {
	offset uint32
	size   uint32
}

// 0x0AAB950C ~ 0x0AB4C31B
var v0AAB950CcrcdataSet = [...]crcdata{
	{0x1000, 0x46},
	{0x1050, 0x4A},
	{0x10A0, 0x1D},
	// text段全部
	{0x00D48F00, 0x47},
	{0x00D48F50, 0x0F},

	{0x00EDDA10, 0x1B4},
	// rdata段高地址处
	{0x00EE0ED2, 0x13},
	{0x00EE0EE9, 0x117},

	{0x099E8490, 0x16E},
	// rsc段几乎全部
	{0x0AE8839E, 0x11B},
	{0x0AE884BD, 0x25},
	{0xFFFFFFFF, 0xFFFFFFFF},
}

// v0AF77FE4sum是0x00401000~0x00401040共17个32位指令值的累加
var v0AF77FE4sum uint32 = 0x839A83C7

func f0082508F(imageBase uintptr, size uint32) int {
	if size == 0 {
		// 0x00A10652
		return 0
	}

	// 0x00825049
	ebp1C := struct {
		baseAddress       uintptr
		allocationBase    uintptr
		allocationProtect uint32
		regionSize        uint32
		state             uint32
		protect           uint32
		dwtype            uint32
	}{
		baseAddress:       0x00400000,
		allocationBase:    0x00400000,
		allocationProtect: 0x80, // PAGE_EXECUTE_WRITECOPY
		regionSize:        0x1000,
		state:             0x1000,     // MEM_COMMIT
		protect:           2,          // PAGE_READONLY
		dwtype:            0x01000000, //MEM_IMAGE
	}
	// win.VirtualQuery(imageBase, ebp1C[:], 28) // 内存布局中region的概念

	// 0x00A0B204
	ebp1C.protect &= 0x80
	if int8(ebp1C.protect) < 0 {
		// 0x00A10657
		return 1
	}
	// 0x0075ACD9
	// win.VirtualProtect(imageBase, size, win.PAGE_EXECUTE_READWRITE, &ebp1C.protect)

	// 0x00825022
	return 1
}

// 0x0A05D00F
func f0A05D00Fmemcpy(dstAddr *uint32, srcAddr uint32, size uint32) {
	var index uint32
	if index < size {
		// 0x0AF88180
		for {
			*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(dstAddr)) + uintptr(index))) = *(*uint8)(unsafe.Pointer(uintptr(srcAddr) + uintptr(index))) // 6->0x4F, 0x0018FF58的值由0x2A070006->0xFFC0C94F

			// 0x0AF8F817
			index++
		}
	}

	// 0x0A8FBAC6
	// lea esp, dword ptr ss:[esp+4]
	// jmp dword ptr ss:[esp-4]
	// these two instructions is equivalent to ret
}

// 壳逻辑
func f006CD259() {
	// push eax
	// push ecx
	// push esi
	// pushfd
	// push edx
	// push ebx
	// push edi

	// 0x00705113
	// 0x007269EA, 设置所有段属性为ERWC
	func() int {
		ebp8imageBase := v00825070imageBase
		if v0082505D == 0 {
			return 0
		}

		// 0x00825061
		ebp2C := struct {
			dummyName  uint32
			dwPageSize uint32
			data       [24]uint8
		}{}
		// win.GetSystemInfo(&ebp2C)
		ebp2C.dwPageSize = 0x1000

		// 修改PE段访问权限
		f0082508F(ebp8imageBase, ebp2C.dwPageSize-1) // (0x00400000, 0x00000FFF)

		// 0x005C751F
		var ebp4numberofSections uint32
		// 0x0082500A
		imageSectionHeader := func(imageBase uintptr, numberofSectionsp *uint32) []pe.SectionHeader32 {
			fileHeader := (*pe.FileHeader)(unsafe.Pointer(imageBase + uintptr(*(*uint32)(unsafe.Pointer(imageBase + 0x3C))) + 4)) // IMAGE_NT_HEADERS.IMAGE_FILE_HEADER, 0x00400144
			if fileHeader == nil {
				return (*[100]pe.SectionHeader32)(unsafe.Pointer(fileHeader))[:0]
			}

			// 0x00824FF8
			*numberofSectionsp = uint32(fileHeader.NumberOfSections)                                                                                                             // IMAGE_NT_HEADERS.IMAGE_FILE_HEADER.NumberOfSections, 0x00400146结果是0x4
			return (*[100]pe.SectionHeader32)(unsafe.Pointer((uintptr(unsafe.Pointer(fileHeader)) + unsafe.Sizeof(pe.FileHeader{}) + unsafe.Sizeof(pe.OptionalHeader32{}))))[:0] // IMAGE_SECTION_HEADER

		}(ebp8imageBase, &ebp4numberofSections)

		// 0x0082502A
		if ebp4numberofSections <= 0 {
			v0082505D = 0
			return 0 // 0x00400238
		}

		for ebp4numberofSections > 0 {
			// 0x00A1063F
			sectionBase := imageSectionHeader[0].VirtualAddress  // &textSection.VirtualAddress, 0x0040244值是0x1000
			sectionSize := imageSectionHeader[0].VirtualSize     // &textSection.VirtualSize, 0x0040240值是0x00D48000
			sectionChar := imageSectionHeader[0].Characteristics // &textSection.Characteristics, 0x004025C值是0x60000020, code, Execute, Read
			// 0x00825080
			if *(*uint8)(unsafe.Pointer((uintptr(unsafe.Pointer(&sectionChar)) + 3)))&0x10 == 0 { // IMAGE_SCN_MEM_SHARED
				// 0x0079ACB5
				// 修改text段访问权限
				f0082508F(ebp8imageBase+uintptr(sectionBase), sectionSize) // (0x00401000, 0x00D48000), (0x01149000, 0x00198000), (0x012E1000, 0x08AC37E4), (0x09DA5000, 0x01519966)
			}

			// 0x006F4DE8
			imageSectionHeader = imageSectionHeader[1:]
			ebp4numberofSections--
		}

		// 0x0082503D
		v0082505D = 0
		return 1
	}()

	// 0x00825074
	// pop edi
	// pop ebx
	// pop edx
	// popfd
	// pop esi
	// pop ecx
	// pop eax

	// 0x0A048F52
	// push eax

	// 0x0AFDE4A1
	// push ecx

	// 0x0A43CFF3
	// push esi
	// pushfd
	// push edx

	// 0x0A930692
	// push ebx

	// 0x09F84632
	// push edi

	// 0x0AD9893F
	// push 0x0A56ED4A
	// push 0x0AA31853
	// ret

	// 0x0AA31853
	func() {
		// 20个字节的局部变量

		ebp4 := v0A327567imageBase
		if v0A0032DB == v09E2BC68 {
			return
		}
		// 0x09FD9EB3
		// push ebx
		// push esi
		v0A8FE0E0 ^= 0 // 0x12A225EC
		v0A746588 ^= 0 // 0xF7598681
		// ebx = 0
		// ecx = 0
		// mov esi, 0x0AAB950C
		// mov eax,esi
		var ebp8len uint32
		v0AF77FE4sum = 0
		var ebpCvalue uint32

		for {
			if v0AAB950CcrcdataSet[0].offset == ^uint32(0) {
				// 0x0A55F42F
			} else {
				// 0x0A8466A8
				// push ebx

				// 0x0A4E17E5
				// push &ebp8len
				// push 4

				// 0x0ABF738E
				// push &ebpCvalue

				// 0x0A889E57
				// push eax // 相当于push &v0AAB950C
				// push ebp4

				// 0x0AD36283
				// push 0x0A0C4D0F
				// push 0x0AAB88A1
				// ret

				// 0x0AAB88A1
				func(imageBase uint32, crcdatas []crcdata, valuep *uint32, offset uint32, lenp *uint32) {
					// 0x09DE88F1
					// push ecx
					//eax := valuep

					// 0x0A7475A4
					// ecx := offset

					// 0x09E6AAB7
					// push ebx
					// cmp eax, eax+ecx

					// 0x09E91A45
					// push esi
					// push edi

					// 0x0A84BA1D
					// edi := crcdatas
					// valuep = eax

					valuep = valuep

					// 0x0A44AD5E
					// ebp4 := eax + ecx
					ebp4 := uintptr(unsafe.Pointer(valuep)) + uintptr(offset)
					// offset = ecx

					offset = offset

					if valuep == (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(valuep))+uintptr(offset))) {
						return // 0x0AF11712
					}

					// 0x0A563538
					//ebx := lenp

					// 0x0A745BAE
					for {
						// ecx = *ebx ;ecx = *lenp

						// 0x0AD362DE
						// eax := *edi
						eax := crcdatas[0].offset
						esi := crcdatas[0].size // v0AAB9510=0x46
						// eax += ecx
						eax += *lenp // 0x1000, 0x1004

						// 0x0A970D30
						eax += imageBase // 0x00401000, 0x00401004
						esi -= *lenp     // 0x46, 0x42
						// cmp offset, esi
						// push 0x0A4EA4E2

						// 0x0AD98EF4
						// push edx
						// push ebx
						// edx = [esp+8] // 0x0A4EA4E2

						// 0x0AF77208
						// ebx = 0x0AFD7E1C
						// cmovae edx, ebx
						// [esp+8] = edx

						// 0x0AF7875F
						// pop ebx
						// pop edx
						// ret

						if offset >= esi {
							// 0x0AFD7E1C
							// push esi

							// 0x0AF12435
							// push eax
							// push valuep

							// 0x0A931011
							// push 0x0B06F488
							// push ecx ;push *lenp

							// 0x0AF975AE
							// push edx ;edx哪里来的？push lenp
							// ecx := 0x0B06F488

							// 0x0A38CEAE
							// edx = 0x0ABE03B6

							// 0x0B071867
							// cmovne ecx, edx
							// [esp+8] = ecx
							// pop edx

							// 0x0A558DE4
							// pop ecx

							if offset == esi {
								// 0x0B06F488
							} else {
								// 0x0ABE03B6
								// push 0x0A440234
								// push 0x0A05D00F
								// ret

								f0A05D00Fmemcpy(valuep, eax, esi)

								// 0x0A440234
								valuep = (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(valuep)) + uintptr(esi))) // 0x0018FF58 + 2
								*lenp = 0
								offset -= esi
								crcdatas := crcdatas[1:]
								if crcdatas[0].offset == ^uint32(0) { // 0x1050
									// 0x0AFDFFD3
								}
								// else 0x0ABE3F38
							}
						} else {
							// 0x0A4EA4E2
							// push offset
							// push eax
							// push valuep
							// push 0x0A32E5EE
							// push 0x0A05D00F
							// ret

							f0A05D00Fmemcpy(valuep, eax, offset)

							// 0x0A32E5EE
							valuep = (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(valuep)) + uintptr(offset))) // 0x0018FF58 + 4
							*lenp += offset                                                                       // 4,8
							// push 0x0ABE3F38
							// ret
						}

						// 0x0ABE3F38
						if uintptr(unsafe.Pointer(valuep)) == ebp4 { // 正常都是0x0018FF5C，校验结束后就不等了
							// 0x0A931F7A
							// push 0x0AF11712
							// ret
							break
						}
					} // for loop 0x0A745BAE

				}(ebp4, v0AAB950CcrcdataSet[:], &ebpCvalue, 4, &ebp8len)

				// 0x0AF11712
				// push edi
				// pop eax
				// pop edi
				// pop esi
				// pop ebx

				// 0x0A0C4D0F
				v0AF77FE4sum += ebpCvalue
			}
		} // for loop
	}()
}

// OEP 0x00DF478C
func main() {
	// check pe

	checkupdate()

	var szcmdLine string = os.Args[1]
	// hInstance
	// hPrevInstance
	// szcmdline
	// SW_SHOWDEFAULT
	f004D7CE5winMain(0x00400000, 0, []uint8(szcmdLine), 10) // call 0x004D,7CE5
}

// 大范围清零操作
func f00DFC986(buf []uint8, size int) {
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
func f00DFC9DD(buf []uint8, value uint8, size int) uint32 {
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
	offset := 16 - (*[3]int)(unsafe.Pointer(&buf))[0]%16
	if offset == 0 {
		s := size & 0x7F
		if s != size {
			// ...
		}
		size := size - s
		f00DFC986(buf, size)
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

	size -= offset
	f00DFC9DD(buf[offset:], 0, size)
	return 0
}
