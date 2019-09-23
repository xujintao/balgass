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

// 要确保v012F7B90 + v012F7B94 = 0
var v012F7B90 uint32 = 0x5F709C77
var v012F7B94 uint32 = 0x44BF19B1
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

// shell entry point, 0x00DF490A
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
	f006CD259()
	// 0x00DF478C
}

var v0082505D uint32           // 0x6A7B9E4D
var v00825070imageBase uintptr // 0x00400000
var v0A327567imageBase uint32 = 0x00400000
var v0A0032DB uint32 = 0x1CA9625E
var v0B287362 uint32
var v0A5578FA uint32
var v09E2BC68 uint32 = 0x2DBE21DC
var v0A8FE0E0 uint32 = 0x12A225EC
var v0A746588sum uint32 = 0xF7598681
var v0AF77FE4sumcalc uint32 = 0x839A83C7 // v0AF77FE4sumcalc是0x00401000~0x00401040共17个32位指令值的累加，需要和v0A746588

type sectionInfo struct {
	addr  uint32
	flag1 uint32
	size  uint32
	flag2 uint32
}

var v09EBB9CCsectionInfoSet = [...]sectionInfo{
	{0x00401000, 0, 0x00D48000, 0},                   // .text
	{0x01149000, 0, 0x00198000, 0},                   // .radata
	{0x012E1000, 0, 0x08AC37E4, 0},                   // .data
	{0x09DA5000, 0, 0x01519966, 0},                   // .rsrc
	{0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF}, // .rsrc
}

type crcData struct {
	addr uint32
	size uint32
}

// 0x0AAB950C ~ 0x0AB4C31B
var v0AAB950CcrcDataSet = [...]crcData{
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

var v0A330A7E uint32
var v0A330A82 uint32

//
//
// --------------validate read and write----------
var v09E8FF96 uint32 = 0xFFFFFFFF
var v0A049D16 uint32 = 0x0AD2D162
var v0ABF9FFCimageBase uint32 = 0x00400000

//
//
// ---------------iat move------------
var v0A563A89 uint32 = 0xCCCCCCCC

type iat struct {
	addr uint32
	size uint32
}

// 没必要
var v09F8682AiatSet = [...]iat{
	{0x099E917C, 0x1F0},
	{0x099E9370, 0x7D0},
	{0xFFFFFFFF, 0xFFFFFFFF},
}

type rdata struct {
	addr uint32
	size uint32
}

var v0A3A1688rdataSet = [...]rdata{
	{0x00D49000, 0x9C0},
	{0xFFFFFFFF, 0xFFFFFFFF},
}

//
//
// -------------text move-----------------
var v0ABDABC7 uint32 = 8
var v0ABE3A4F uint32 = 0x0A39221A
var v0A5D6CCF uint32 = 7
var v0A5557A4 uint32 = 0x0AD3B521
var v0A0C61AD uint32 = 0x5420110F
var v0A0526C0 uint32 = 6
var v0AC33EA1 uint32 = 0x0AF1023E
var v0AF0E89A uint32 = 0x546CE7AB

type block struct {
	addr uint32
	size uint32
}

var v0A32F9C5textSet = [...]block{
	{0x0D7CE5, 5},
	{0x0D7CED, 0x1078},
	{0x0D8D67, 0x28},
	{0x0D8D91, 2},
	{0x0D8D95, 0xEC},
	{0x0D8E83, 0x164},
	{0x0A956E8A, 0xD},
	{0xFFFFFFFF, 0xFFFFFFFF},
}
var v0AFE0B1FimageBase uint32 = 0x00400000
var v0A84CEC1 uint32 = 0x7B853693

var v0AD2F33DopcodeSet = [...]uint8{
	0xAA,
	0x04,
	0x1C,
}

// ------------------winmain-------------------
//
var v0A7483B4 uint32 = 8
var v0A56E4E2label1 uint32 = 0x0AD56E8A
var v0ABF8B07 uint32 = 7
var v09FC0B5Alabel2 uint32 = 0x0A8FE9A6
var v0AA0D71B uint32 = 6
var v0A557660label3 uint32 = 0x0B10933D

var v0A74573D uint32 = 0x3F
var v0A38DB63 uint32 = 0xFC24548D // v0A38DB63 = 0x310F9090 = v0AF77CE1 ^ v0B1090BA
var v0AF77CE1 uint32 = 0x89B8FCE2
var v0B1090BA uint32 = 0xB8B76C72
var v09FBB49C uint32 = 0x10648B16
var v0A325B9F uint32 = 0x721655D3

// ------------------end-----------------------
//
//
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
func f0A05D00Fmemcpy(dstAddr unsafe.Pointer, srcAddr unsafe.Pointer, size uint32) {
	var index uint32
	if index < size {
		// 0x0AF88180
		for {
			*(*uint8)(unsafe.Pointer(uintptr(dstAddr) + uintptr(index))) = *(*uint8)(unsafe.Pointer(uintptr(srcAddr) + uintptr(index))) // 6->0x4F, 0x0018FF58的值由0x2A070006->0xFFC0C94F

			// 0x0AF8F817
			index++
		}
	}

	// 0x0A8FBAC6
	// lea esp, dword ptr ss:[esp+4]
	// jmp dword ptr ss:[esp-4]
	// these two instructions is equivalent to ret
}

func f0AAB88A1(imageBase uint32, crcDatas []crcData, valuep unsafe.Pointer, step uint32, lenp *uint32) []crcData {
	// 0x09DE88F1
	ebp4 := uintptr(valuep) + uintptr(step) // valuep + 8
	if valuep == unsafe.Pointer(uintptr(valuep)+uintptr(step)) {
		return crcDatas // 0x0AF11712
	}
	// 0x0A563538
	// 0x0A745BAE
	for {
		// 0x0AD362DE
		addr := crcDatas[0].addr
		size := crcDatas[0].size // v0AAB9510=0x46
		addr += *lenp            // 0x1000, 0x1004

		// 0x0A970D30
		addr += imageBase // 0x00401000, 0x00401004
		size -= *lenp     // 0x46, 0x42

		// 0x0AD98EF4
		// 0x0AF77208
		// 0x0AF7875F
		if step >= size {
			// 0x0AFD7E1C
			// 0x0AF12435
			// 0x0A931011
			// 0x0AF975AE
			// 0x0A38CEAE
			// 0x0B071867
			// 0x0A558DE4
			if step == size {
				// 0x0B06F488
				f0A05D00Fmemcpy(valuep, unsafe.Pointer(uintptr(addr)), size)

				// 0x0AF8F334
				valuep = unsafe.Pointer(uintptr(valuep) + uintptr(size)) // 0x0018FF58 + 4
				crcDatas := crcDatas[1:]
				if crcDatas[0].addr == ^uint32(0) {
					// 0x0AF11712
					return crcDatas
				}
				// 0x0A43CE90
				*lenp = 0

				// 0x0ABE3F38 // 不需要再判断凑整
			} else {
				// 0x0ABE03B6
				// 0x0A05D00F
				f0A05D00Fmemcpy(valuep, unsafe.Pointer(uintptr(addr)), size)

				// 0x0A440234
				valuep = unsafe.Pointer(uintptr(valuep) + uintptr(size)) // 0x0018FF58 + 2
				*lenp = 0
				step -= size
				crcDatas := crcDatas[1:]
				if crcDatas[0].size == 0xFFFFFFFF { // 0x1050
					// 0x0AFDFFD3
					// 0x0AF7CDF7
					// 0x0A83C5D1 memset
					func(p unsafe.Pointer, value uint32, size uint32) {
						// 0x0A55A765
						// 0x09FC7C50
						// 0x09EB9524
						// 0x09F8436A

						if step > 0 {
							// 0x0A38CF80
							// 0x0A9F8454
							// 0x0A88BF90
							// 0x0AA05578
							// 0x0A437662
							// 0x0A38C4F4
							// 0x09FD1F88
							// 0x0A842745
							// 0x0A6035A3
							// 0x09E270C3
							// 0x0AF7A605
						}
						// 0x0AD70014
					}(valuep, 0, step)

					// 0x09FDCDCA
					return crcDatas // 0x0AF11712
				}
				// else 0x0ABE3F38
			}
		} else {
			// 0x0A4EA4E2
			// 0x0A05D00F
			f0A05D00Fmemcpy(valuep, unsafe.Pointer(uintptr(addr)), step)

			// 0x0A32E5EE
			valuep = unsafe.Pointer(uintptr(valuep) + uintptr(step)) // 0x0018FF58 + 4
			*lenp += step                                            // 4, 8, ... or 8, 16, ...
			// 0x0ABE3F38
		}

		// 0x0ABE3F38
		if uintptr(unsafe.Pointer(valuep)) == ebp4 { // 正常都是0x0018FF5C，否则就要循环读取下一组校验数据来凑满4个字节
			// 0x0A931F7A
			// 0x0AF11712
			return crcDatas
		}
	} // for loop 0x0A745BAE
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

	// 0x0AA31853 crc校验
	func() {
		// 20个字节的局部变量
		ebp4imageBase := v0A327567imageBase
		if v0A0032DB == v09E2BC68 {
			return
		}
		// 0x09FD9EB3
		v0A8FE0E0 ^= 0    // 0x12A225EC
		v0A746588sum ^= 0 // 0xF7598681
		crcDatas := v0AAB950CcrcDataSet[:]
		var ebp8len uint32
		v0AF77FE4sumcalc = 0
		var ebpCvalue uint32
		var ebp14value uint64

		for {
			if crcDatas[0].size == 0xFFFFFFFF {
				// 0x0A55F42F
				if v0AF77FE4sumcalc == v0A746588sum {
					// 0x09FC934F
					// push 0x09EBB9CC
					// push ebp4imageBase

					// 0x0A5508D5
					// push 0x09EAC485
					// push 0x0B286D7E

					// 0x0B286D7E
					nRet := func(imageBase uintptr, sectionInfos []sectionInfo) int {
						fileHeader := (*pe.FileHeader)(unsafe.Pointer(imageBase + uintptr(*(*uint32)(unsafe.Pointer(imageBase + 0x3C))) + 4)) // IMAGE_NT_HEADERS.IMAGE_FILE_HEADER, 0x00400144
						if fileHeader == nil {
							// return (*[100]pe.SectionHeader32)(unsafe.Pointer(fileHeader))[:0]
						}
						imageBase = 0

						// 0x0AD73ED5
						eaxSectionHeader := (*[100]pe.SectionHeader32)(unsafe.Pointer((uintptr(unsafe.Pointer(fileHeader)) + unsafe.Sizeof(pe.FileHeader{}) + unsafe.Sizeof(pe.OptionalHeader32{}))))[:]

						if fileHeader.NumberOfSections > 0 {
							for {
								// 0x0A83FA73
								if sectionInfos[0].addr&sectionInfos[0].flag1 != 0xFFFFFFFF { // 判断是否需要校验
									// 0x0A55B8FD
									if eaxSectionHeader[0].VirtualSize == sectionInfos[0].size { // 校验段大小
										if sectionInfos[0].flag2 == 0 {
											// 0x0A558FFD
											eaxSectionHeader = eaxSectionHeader[1:]
											imageBase++
											sectionInfos = sectionInfos[1:]

											// 0x0AF789C1
											if imageBase >= uintptr(fileHeader.NumberOfSections) {
												break // 0x0AA2BBD5
											}
										}
									}
								}
								return 0
							} // for loop
						}

						// 0x0AA2BBD5
						if sectionInfos[0].addr&sectionInfos[0].flag1 != 0xFFFFFFFF {
							return 0
						}
						// 0x0ABB6AFC
						return 1
					}(uintptr(ebp4imageBase), v09EBB9CCsectionInfoSet[:])

					// 0x09EAC485, 解密过程
					if nRet != 0 {
						// 0x0A84B9A6
						// 0x0A56CFD4

						ebp8len = 0
						v0A0032DB = v09E2BC68

						for {
							// 0x0AD80743
							if crcDatas[0].size == 0xFFFFFFFF {
								// 0x09FB8D2B
								v0A5578FA ^= v0B287362

								// 0x09F0180
								// 0x0ABDAB3B
								// pop esi
								// pop ebx

								// 0x0AFDC089
								return
							} else {
								// 0x09FD2C99
								// push 0
								ebpCvalue = ebp8len
								// ebx := crcDatas // 返回值？
								// 0x0AAB88A1, 密文
								crcDatas = f0AAB88A1(ebp4imageBase, crcDatas, unsafe.Pointer(&ebp14value), 8, &ebp8len)

								// 0x0A9527D0
								// push 0x0AD70D9F
								// push 0x0A9003B1
								// ret

								// 0x0A9003B1, 密钥
								func() {
									eax := v0A330A7E
									ecx := v0A330A82
									// push ebx
									// push esi
									// push edi
									var edi uint32
									esi := 0x40

									for {
										// 0x0AD80CB9
										edi = ecx
										edi >>= 5
										ebx := ecx
										ebx <<= 4
										edi ^= ebx
										edi += ecx
										edx := ebx
										ebx &= 3
										ebx = *(*uint32)(unsafe.Pointer(uintptr(ebx*4) + uintptr(unsafe.Pointer(&v0A5578FA))))
										ebx += edx
										edi ^= ebx
										eax += edi
										edi = eax
										edi >>= 5
										ebx = eax
										ebx <<= 4
										edi ^= ebx
										edx -= 0x61C88647
										ebx = edx
										ebx >>= 0xB
										ebx &= 3
										ebx = *(*uint32)(unsafe.Pointer(uintptr(ebx*4) + uintptr(unsafe.Pointer(&v0A5578FA))))
										edi += eax
										ebx += edx
										edi ^= ebx
										ecx += edi
										esi--
										if esi == 0 {
											break // 0x0A05177A
										}
									}

									// 0x0A05177A
									// pop edi
									// pop esi
									v0A330A7E = eax
									v0A330A82 = ecx
									// pop ebx
								}()

								// 0x0AD70D9F, 明文
								ebp14value = uint64(uint32(ebp14value) ^ v0A330A7E + (uint32(ebp14value>>32)^v0A330A82)<<32)
								// push &ebpCvalue，当长度用了
								// push ebx
								// push 8
								// push &ebp14value
								// push ebp4imageBase
								// push 0x0AFD8B4F
								// push 0x09FE64D2
								// ret

								// 0x09FE64D2, 写明文
								func(imageBase uint32, valuep unsafe.Pointer, step uint32, crcDatas []crcData, lenp *uint32) {
									var ebp4 uint32
									var ebp8 uint32
									// push ecx
									// push ecx
									// eax = valuep
									// edx = step
									// edx += eax
									// ecx = 0
									ebp4 = 0
									ebp8 = uint32(uintptr(valuep) + uintptr(step))
									if valuep != unsafe.Pointer(uintptr(valuep)+uintptr(step)) {
										return // 0x0AD97243
									}

									// 0x0AD7BE4D
									// push ebx
									// ebx := crcDatas
									// push esi

									// 0x0AF7C50F
									// push edi

									// 0x0AD9274B

									for {
										// 0x0AF97BBD
										// eax = lenp

										// 0x0A9D32BC
										// edx = *eax

										// 0x0A604C8D
										// eax = crcDatas[0].addr
										// esi = crcDatas[0].size
										// edi = step
										// eax += edx
										addr := crcDatas[0].addr
										size := crcDatas[0].size
										addr += *lenp

										// 0x0AFDAA70
										// eax += imageBase
										// esi -= edx
										// edx = valuep
										addr += imageBase
										size -= *lenp

										// 0x0A057A29
										// edi -= ecx
										// ecx += edx
										// cmp esi,edi

										// 0x0AF8A955
										// 0x0A84973E
										// 0x09FDF556
										if size <= step {
											// 0x0A0023F3
										} else {
											// 0x09FC481F
											// push edi
											// push ecx
											// push eax
											// push 0x09FE2A75
											// push 0x0A05D00F
											// ret

											// 0x0A05D00F
											f0A05D00Fmemcpy(unsafe.Pointer(uintptr(addr)), valuep, step)

											// 0x09FE2A75
											// eax = lenp
											// ebp4 += edi
											ebp4 += step

											// 0x0A95081C
											// *eax = edi
											*lenp += step

											// 0x0AD310AF
											// 0x09E2D4D0
											// eax = valuep
											// ecx = ebp4
											// ecx += eax

											// 0x0A131034
											if uintptr(ebp4)+uintptr(valuep) == uintptr(ebp8) {
												// 0x0A0485AD
												// 0x0AF7ED9E
												return // 0x0AD97243
											}
											// 0x0A5FC621
											// ecx = ebp4
											// 0x0AF97BBD
										}
									} // for loop 0x0AF97BBD

								}(ebp4imageBase, unsafe.Pointer(&ebp14value), 8, crcDatas, &ebp8len)

								// 0x0AFD8B4F
								// 0x0AFD8B52 等效于 0x0AD80743
							}
						} // for loop
					}
					// 0x0AF7E275
				}
				// 0x0AF7E275

			} else {
				// 0x0AAB88A1
				crcDatas = f0AAB88A1(ebp4imageBase, crcDatas, unsafe.Pointer(&ebpCvalue), 4, &ebp8len)

				// 0x0A0C4D0F
				v0AF77FE4sumcalc += ebpCvalue
			}
		} // for loop

	}()

	// 0x0A56ED4A
	// pop edi
	// pop ebx
	// pop edx
	// popfd
	// pop esi
	// pop ecx
	// pop eax

	// 0x0A33348C
	// 0x09F8A8BD
	// 0x0B10B013
	// push eax
	// push ecx
	// push esi
	// push edx
	// push ebx
	// push edi

	// 0x0A9F55F3
	// 0x0AD2D162, 验证读写权限？
	func() {
		// 0x24字节局部变量
		if v09E8FF96 == 0 {
			return
		}
		// 0x0A896CF3
		ebp24 := struct {
			dummyName  uint32
			dwPageSize uint32
			data       [24]uint8
		}{}
		// win.GetSystemInfo(&ebp24)
		ebp24.dwPageSize = 0x1000

		// 0x0A9F5118
		esi := ebp24.dwPageSize
		var eaxImageBase uint32 = 0x0AD2D162
		eaxImageBase -= v0A049D16
		eaxImageBase += v0ABF9FFCimageBase
		ecxNToffset := *(*uint32)(unsafe.Pointer((uintptr(eaxImageBase) + 0x3C)))                                     // 0x140
		optionHeader := (*pe.OptionalHeader32)(unsafe.Pointer((uintptr(eaxImageBase) + uintptr(ecxNToffset) + 0x18))) // 0x00400158
		ecxCodeBase := optionHeader.BaseOfCode + eaxImageBase                                                         // 0x00401000
		edxCodeSize := optionHeader.SizeOfCode + eaxImageBase                                                         // 0x0B2BE8D6

		for {
			// 0x0AD70CC7
			if ecxCodeBase >= edxCodeSize {
				break // 0x0AF98284
			}
			*(*uint8)(unsafe.Pointer(uintptr(ecxCodeBase))) = *(*uint8)(unsafe.Pointer(uintptr(ecxCodeBase))) // r and w
			// 0x0A950359
			ecxCodeBase += esi
		}

		// 0x0AF98284
		v09E8FF96 &= 0

		// 0x0AF7FFD5
		return
	}()

	// 0x0ABD82B9
	// pop edi
	// pop ebx
	// pop edx
	// popfd
	// pop esi
	// pop ecx
	// pop eax

	// 0x0AD81E07
	// 0x0A333AD2
	// 0x0AA0E7FA
	// push ecx
	// push edx
	// push edi
	// push ebx
	// push esi
	// pushfd

	// 0x0AFD850D
	// 0x0A44A7FD, IAT move
	func() {
		// 0xC字节的布局变量
		if v09F8682AiatSet[0].addr == 0 {
			// 0x0A8429DF
			return
		}
		// 0x0A4E8FC1
		ebp4iats := v09F8682AiatSet[:]
		ebp8rdatas := v0A3A1688rdataSet[:]
		esiRDATAaddr := ebp8rdatas[0].addr + 0x00400000 // 0x01149000
		ediRDATAsize := ebp8rdatas[0].size
		ebxIATsize := ebp4iats[0].size
		ecxIATaddr := ebp4iats[0].addr + 0x00400000 // 0x09DE917C advapi32.&CryptGetUserKey
		if v09F8682AiatSet[0].addr == 0xFFFFFFFF {
			return // 0x0A56CD8C
		}

		// 0x09EAE30B
		ebpC := v0A3A1688rdataSet[0].addr

		for {
			// 0x0ABE28CC
			if ebpC == 0x0FFFFFFFF {
				return // 0x0A56CD8C
			}

			// 0x0A9F96DB
			*(*uint32)(unsafe.Pointer(uintptr(esiRDATAaddr))) = *(*uint32)(unsafe.Pointer(uintptr(ecxIATaddr))) // 0x75AF350E advapi32.CryptGetUserKey
			var eax uint32 = 4
			esiRDATAaddr += eax // 0x01149000, 0x01149004
			ecxIATaddr += eax   // 0x09DE917C, 0x09DE9180
			ebxIATsize -= eax   // 0x1F0, 0x1EC
			ediRDATAsize -= eax // 0x09C0, 0x09BC

			if ebxIATsize == 0 { // kernel32.dll
				// 0x0AC353B9
				ebp4iats = ebp4iats[1:]
				ecxIATaddr = ebp4iats[0].addr + 0x00400000
				ebxIATsize = ebp4iats[0].size
			}

			// 0x0AB4F818
			if ediRDATAsize == 0 {
				// 0x0A4EB197
				ebp8rdatas = ebp8rdatas[1:]
				esiRDATAaddr = ebp8rdatas[0].addr + 0x00400000
				ediRDATAsize = ebp8rdatas[0].size
				ebpC = esiRDATAaddr
			}

			// 0x0A931EAB
			if ebp4iats[0].addr == 0xFFFFFFFF { // 0x099E917C
				return // 0x0A56CD8C
			}
		}
	}()

	// 0x0A56CD8C
	// 0x0A8429DF
	v0A563A89 = 0

	// 0x0AFD8512
	// 0x0AAB829D
	// popfd
	// pop esi
	// pop ebx
	// pop edi
	// pop edx
	// pop ecx

	// 0x0AF8BBBA
	// 0x0A60483B
	// 0x0A56C231
	// push 0x09FB8C65
	// push 0x0AD323B3
	// ret

	// 0x0AD323B3, text move
	func() {
		// push 0x0A56DEAE
		// push 0x09FFDE30
		// ret

		// 0x09FFDE30
		func() {
			// push 0x00DDD650
			// push 0x0A0509F4
			// ret

			// 0x0A0509F4
			func() {
				// push esi
				// push ecx
				// push ebx
				// push edx
				// pushfd
				// push edi

				// 0x0AD707CA
				// push 0x0B07273B
				// push 0x0AF7CDBD
				// ret

				var ebp4 uint32
				var ebp8 uint32
				var ebpC uint32

				// 0x0AF7CDBD
				func() {
					// push esi
					ebp4 = v0ABE3A4F // 0x0A39221A
					// push edi
					esiops := v0AD2F33DopcodeSet[:]
					editexts := v0A32F9C5textSet[:]
					if v0A32F9C5textSet[0].addr != 0xFFFFFFFF {
						// 0x09FFCBF1
						// push ebx
						// push ebp

						for {
							// 0x0AAB316A
							eaxAddr := editexts[0].addr

							// 0x0A43E504
							// 0x0AD3B4F7
							// 0x09E2F54F
							// 0x09EBC1C0
							// 0x0AF126B7
							if eaxAddr == 0xFFFFFFFE {
								// 0x0A8FC82C
							}
							// 0x0A4416F2
							// 0x0A9F9C72
							eaxAddr += v0AFE0B1FimageBase

							// 0x0AD738B4
							edxSize := editexts[0].size

							// 0x0A9F6726
							// 0x0A936A10
							// 0x0AFDD768
							if edxSize > 3 {
								// 0x09FB7EC3
								ecx := edxSize / 4
								edxSize %= 4

								// 商迭代
								for {
									// 0x0AD84E83
									*(*uint32)(unsafe.Pointer(uintptr(eaxAddr))) = v0A84CEC1 - *(*uint32)(unsafe.Pointer((&esiops[0])))

									// 0x0AA2F6ED
									eaxAddr += 4
									esiops = esiops[4:]
									ecx--

									// 0x0A8417D3
									// 0x0ABE42C0
									// 0x0A041758
									// 0x0ABFA56D
									// 0x0B284F49
									// 0x0A55F76C
									// 0x0AB4FADB
									// 0x0A900853
									if ecx == 0 {
										break // // 0x0AF0FE48
									}
								} // for loop 0x0AD84E83
							}
							// 0x0AF0FE48
							// 0x0A603BC5, 余数迭代
							for edxSize > 0 {
								// 0x0A4EAF80
								*(*uint8)(unsafe.Pointer(uintptr(eaxAddr))) = uint8(v0A84CEC1) - *(*uint8)(unsafe.Pointer((&esiops[0])))
								eaxAddr++
								esiops = esiops[1:]
								edxSize--
							} // for loop 0x0A4EAF80

							// 0x09F88775, next
							editexts = editexts[1:]

							// 0x09FE278E
							// 0x0A7462B5
							// 0x0B10D023
							// 0x09EB7BBB
							if editexts[0].addr == 0xFFFFFFFF {
								// 0x09FB5670
								break
							}
						} // for loop 0x0AAB316A
						// pop ebp
						// pop ebx
					}
					// 0x0A89291D
					ebp8 = v0A5557A4 // 0x0AD3B521
					// pop edi
					v0AF0E89A = v0A0C61AD // 0x5420110F
					ebpC = v0AC33EA1      // 0x0AF1023E
					// pop esi
				}()

				// 0x0B07273B
				// pop edi
				// popfd
				// pop edx
				// pop ebx
				// pop ecx
				// pop esi

				// 0x0A55A721
				// replace 0x0B07273B with (ebpC 0x0AF1023E)
				// 0x9E2FFC0
			}()

			// replace 0x00DDD650 with (ebp8 0x0AD3B521)
		}()

		// replace 0x0A56DEAE with (ebp4 0x0A39221A)
	}()

	// 0x00E10691, 计算v012F7B90和v012F7B94
	func() {
		// 0x10字节局部变量
		type FileTime struct {
			dwLowDateTime  uint32
			dwHighDateTime uint32
		}
		ebp8fileTime := FileTime{}
		// push ebx
		// push edi
		var edi uint32 = 0xBB40E64E
		var ebx uint32 = 0xFFFF0000
		// var v012F7B90 uint32 = 0xBB40E64E
		if v012F7B90 != edi || v012F7B90 != ebx {
			v012F7B94 = -v012F7B90
			return
		}

		// push esi
		//win.GetSystemTimeAsFileTime(&ebp8fileTime)
		esi := ebp8fileTime.dwHighDateTime
		esi ^= ebp8fileTime.dwLowDateTime

		var eax uint32
		// win.GetCurrentProcessId()
		esi ^= eax

		eax = win.GetCurrentThreadId()
		esi ^= eax

		eax = win.GetTickCount()
		esi ^= eax

		var ebp10 uint64
		// win.QueryPerfomanceCounter(&ebp10)
		eax = uint32(ebp10 >> 32)
		eax ^= uint32(ebp10)
		esi ^= eax

		if esi != edi && esi == ebx { // esi是动态变化的，不可能等于ebx
			eax = esi
			eax <<= 16
			esi |= eax
		}
		v012F7B90 = esi
		v012F7B94 = -esi
		// pop esi
		// pop edi
		// pop ebx
	}()
}

// --------------------------------------------------------------------
// OEP 0x00DF478C, f00DF478C
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
