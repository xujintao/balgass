package main

import (
	"debug/pe"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/xujintao/balgass/win"
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

// libc cdecl无ebp帧栈

func f00DE84C8() {

}

func f00DF7693errno() int {
	return 0
}

// libc.stdlib--------------------------------
// atof atoi strtod strtol strtoul
// calloc malloc free
// abort atexit exit getenv system
// bsearch qsort abs div abls ldiv rand srand
// mblen mbstowcs mbtowc wcstombs wctomb
func f00DECBD1atoi(buf []uint8) int {
	n, err := strconv.Atoi(string(buf))
	if err != nil {
		n = 0
	}
	return n
}

func f00DE7538free(p unsafe.Pointer) {}

func f00DE84E3abort() {}

func f00DE8A9Bsrand(seed int64) {
	rand.Seed(seed)
}

func f00DE8AADrand() int {
	return rand.Int()
}

func f00DE852Fnew(x int) unsafe.Pointer {
	return nil
}
func f00DE64BCnew(n uint) []uint8 {
	return make([]uint8, n)
}
func f00DE7BEAdelete(buf []uint8) {
	buf = buf[:0]
}

func f00DF08E8abs(v int) int {
	return 0
}

// libc.string
// memchr memcmp memcpy memmove memset
// strcat strncat strcmp strncmp strcoll
// strcpy strncpy strcspn strerror strlen
// strpbrk strrchr strspn strstr strtok strxfrm
func f00DE8857memcpy_s(dst unsafe.Pointer, dstSize int, src unsafe.Pointer, size int) int {
	return 0
}
func f00DE7C00strlen(str []uint8) int {
	return 0
}

func f00DE7C90memcpy(dst, src []uint8, len int) {
}

func f00DE8000strcpy(dst []uint8, src []uint8) {
	copy(dst, src)
}

func f00DE9370strncpy(dst []uint8, src []uint8, size int) {
	copy(dst, src[:size])
}

func f00DEC9B0memchr(s []uint8, c uint8, n int) []uint8 {
	return nil
}

func f00DE94F0strcmp(str1 []uint8, str2 []uint8) int {
	return 0
}

func f00DF30EFstrcpysafe(dst []uint8, size int, src string) {
	copy(dst[:size], src[:])
}

func f00DE8010strcat(dst []uint8, src string) {
	dst = append(dst, src...)
}

func f00DECB2Estrcatsafe(dst []uint8, size int, src string) {
	dst = append(dst[:size], src...)
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

func f00DE817Asprintf(buf []uint8, fmt string, a ...interface{}) int {
	// equal
	// str := fmt.Sprint(strfmt, a)
	// buf = []uint8(str) // 从堆上复制到data段
	// return len(str)

	if fmt == "" || len(buf) == 0 {
		return 0
	}
	// ebp20
	i := info{
		m00: buf,       // ebp20
		m04: 0x7FFFFFF, // ebp1C
		m08: buf,       // ebp18
		m0C: 0x42,      // ebp14
	}
	n := f00DFCCB0vsprintf(&i, fmt, nil, a)
	if n < 0 {
		// ...
	}
	i.m00[0] = 0
	return n
}

// 拿着eax做栈分配，这个骚套路是什么API？目的是什么？
func f00DE8A70chkstk() {}

func f00DEE871setlocale(category uint32, locale string) {

}

// libc.io------------------------------
// fclose clearerr feof ferror fflush fgetpos
// fopen fread freopen fseek fsetpos ftell fwrite
// remove rename rewind setbuf setvbuf tmpfile tmpnam
// fprintf printf sprintf vfprintf vprintf vsprintf
// fscanf scanf sscanf fgetc fgets fputc fputs
// getc getchar gets putc putchar puts ungetc
// perror snprintf
func f00DF6F11fopen(f **os.File, filename []rune, mode []rune) {
	// *f = f00DF6E34(filename, mode, 0x80)
}

func f00DE909Efopen(fileName, mode string) *os.File {
	f, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	return f
}

func f00DE8C84close(f *os.File) int {
	err := f.Close()
	if err != nil {
		return -1
	}
	return 0
}

func f00DE8FBDfread(buf []uint8, size uint, num uint, f *os.File) uint {
	n, err := f.Read(buf[:size*num])
	if err != nil {
		return 0
	}
	return uint(n)
}

// ftell: Get current position in stream
func f00DEFCD4ftell(f *os.File) int64 {
	pos, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return -1
	}
	return pos
}

func f00DEFA34fseek(f *os.File, offset int, whence int) int {
	pos, err := f.Seek(int64(offset), whence)
	if err != nil {
		return -1
	}
	return int(pos)
}

func f00DE92E0strstr(haystack string, needle string) []uint8 {
	return nil
}

func f00DECD20fscanf(f *os.File, strfmt string, y []uint8) int {
	return -1
}

func f00DE8BF6atexit(func()) {}

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
func f00DECDA2(t16p *t16, prevt16p *t16) {
	t16p.f0Ch[0] = 0
	if prevt16p == nil {
		// f00DFC3E9()
		v := func() interface{} {
			// f00DFC370
			v := func() interface{} {
				const v012F7DBC int = 4 // -1
				errno := win.GetLastError()
				// f00DFC1FB()
				flsGetValue := func() func(int) interface{} {
					const v012F7DC0tls int = 25 // -1
					var v interface{}           // v := kernel32.TlsGetValue(v012F7DC0tls)
					if v == nil {
						v2 := f00DFC160flsDec(v09D9DB88flsGetValue)
						// v = kernel32.TlsSetValue(v012F7DC0tls, v2)
						v = v2
					}
					return v.(func(int) interface{}) // kernel32.FlsGetValue
				}()
				v := flsGetValue(v012F7DBC)
				if v == nil {
					// ...
				}
				win.SetLastError(errno)
				return v
			}()
			if v == nil {
				// f00DF0407(0x10)
			}
			return v
		}()
		tlsv := v.(*tlsvalue)
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

// n := func(bufp *[]uint8, fmt string, locale int, a []interface{}) int {
// 好复杂的函数
func f00DFCCB0vsprintf(infop *info, format string, t16p *t16, a []interface{}) int {
	// 278h字节的局部变量
	// 功能就是copy(buf, fmt.Sprintf(format, a...))

	var ebp25C t16 // 0x0018,DCB8
	ebp24C := infop
	// ebp224 := a
	// ebp4 = v012F7B90 ^ ebp

	f00DECDA2(&ebp25C, t16p) // ebp25C.f00DECDA2(t16p)
	if infop == nil {
		// ...
	}
	if infop.m0C == 40 {
		// ...
	}

	var ebp23C string = format[1:] // 拿掉 '>' 字符
	// ...

	var ebp228 uint32 // 0x0018DCEC

	// var ebp211 uint8 = ebp23C[0] // c语言风格
	for len(ebp23C) != 0 {
		ebp23C = ebp23C[1:]
		if ebp228 < 0 { // jl(jump less) 有符号跳转
			break
		}

		// f00E036C2()
		nRet := func(c uint8) uint32 {
			var ebp10 t16
			f00DECDA2(&ebp25C, &ebp10)
			nRet := uint32(ebp10.f00h.fC8h.f7Ch & 8000)
			if ebp10.f0Ch[0]&0xFF != 0 {
				// ...
			}
			return nRet
		}(ebp23C[0])
		if nRet != 0 {
			// ...
		}

		// f00DFCBD0()
		func() {
			// 初始ecx为ebp24C, esi=&ebp228
			if ebp24C.m0C == 0x40 || len(ebp24C.m08) != 0 {
				ebp24C.m04--
				if ebp24C.m04 >= 0 {
					ebp24C.m00[0] = ebp23C[0]
					ebp24C.m00 = ebp24C.m00[1:]
				} else {
					// f00DFCA6C
				}
				if ebp23C[0] == ^uint8(0) {
					ebp228 |= uint32(ebp23C[0])
					return
				}
			}
			ebp228++
			return
		}()
	}
	return int(ebp228)
}

type info struct {
	m00 []uint8
	m04 int
	m08 []uint8
	m0C int
}

func f00DF0805(buf []uint8, format string, a ...interface{}) {
	// f00DF0787
	func(buf []uint8, format string, x *t16, a ...interface{}) int {
		// 如果接下来会用到ebx，那么先把ebx压栈

		// c里面判断字符串指针变量是否为空
		// 等效为go里面判断string类型变量长度是否为0
		// 或者切片长度是否为0
		if len(format) == 0 || len(buf) == 0 {
			// ...
			return -1 // ?
		}

		i := info{
			m00: buf,       // ebp20
			m04: 0x7FFFFFF, // ebp1C
			m08: buf,       // ebp18
			m0C: 0x42,
		}
		cnt := f00DFCCB0vsprintf(&i, format, x, a) // 其实就是把logconf字符串copy到buf切片
		// i.f00h = append(i.f00h,0) // golang不需要追加0
		return cnt
	}(buf, format, nil, a...)
}

// libc.time
func f00DEE9E1time(t *time.Duration) time.Duration {
	if t == nil {
		return time.Duration(time.Now().Unix())
	}
	*t = time.Duration(time.Now().Unix())
	return *t
}

// libc.math
func f00DE76F6round(v float64) float64 {
	return math.Round(v)
}

func f00DE76C0roundf(v float32) float32 {
	return float32(math.Round(float64(v)))
}

// -----------------section ERWC-----------------
var v0082505D uint32 = 0x6A7B9E4D // 0表示done
var v00825070imageBase uintptr    // 0x00400000

// -----------------crc validate--------------
var v0A327567imageBase uintptr = 0x00400000
var v0A8FE0E0unused uint32 = 0x12A225EC  // 未使用
var v0A746588sum uint32 = 0xF7598681     // 密文校验和
var v0AF77FE4sumcalc uint32 = 0x839A83C7 // 密文完整性校验和
var v09E2BC68 uint32 = 0x2DBE21DC
var v0A0032DB uint32 = 0x1CA9625E // v0A0032DB = v09E2BC68 只解密一次

type block struct {
	addr uintptr
	size uint32
}

// 0x0AAB950C ~ 0x0AB4C31C
var v0AAB950CcrcDataSet = [...]block{
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
	{0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF}, // end
}

var v0A330A7ELo uint32 = 0x55E8636A
var v0A330A82Hi uint32 = 0x278AC0C9
var v0A5578FA = [4]uint32{
	0x16A621E8, // ^= v0B287362 解密成功会做异或运算
	0x24719CBE,
	0x6A91B35E,
	0xD0BADE22,
}

// v0B287362 解密成功会做异或运算
var v0B287362 uint32 = 0xCD6ABE04

// v0A88F1F8 完整性验证失败入口
var v0A88F1F8 uint32 = 0x0A9F55F3

// --------------validate read and write----------
var v09E8FF96 uint32 = 0xFFFFFFFF // 0表示done
var v0A049D16 uint32 = 0x0AD2D162
var v0ABF9FFCimageBase uint32 = 0x00400000

// ---------------iat move------------
var v0A563A89 uint32 = 0xCCCCCCCC // 0表示done

// 没必要分成2段
var v09F8682AiatSet = [...]block{
	{0x099E917C, 0x1F0},
	{0x099E9370, 0x7D0},
	{0xFFFFFFFF, 0xFFFFFFFF},
}

var v0A3A1688rdataSet = [...]block{
	{0x00D49000, 0x9C0},
	{0xFFFFFFFF, 0xFFFFFFFF},
}

// -------------text move-----------------
var v0A32F9C5textSet = [...]block{
	{0x000D7CE5, 5},
	{0x000D7CED, 0x1078},
	{0x000D8D67, 0x28},
	{0x000D8D91, 2},
	{0x000D8D95, 0xEC},
	{0x000D8E83, 0x164},
	{0x0A956E8A, 0xD},
	{0xFFFFFFFF, 0xFFFFFFFF},
}
var v0AFE0B1FimageBase uintptr = 0x00400000
var v0A84CEC1 uint32 = 0x7B853693
var v0AD2F33DopcodeSet = [...]uint8{
	0xAA,
	0x04,
	0x1C,
	// ...
}

// done label
var v0ABDABC7 uint32 = 8
var v0ABE3A4F uint32 = 0x0A39221A // label11
var v0A5D6CCF uint32 = 7
var v0A5557A4 uint32 = 0x0AD3B521 // label2
var v0A0526C0 uint32 = 6
var v0AC33EA1 uint32 = 0x0AF1023E // label3

// done
var v0A0C61AD uint32 = 0x5420110F
var v0AF0E89A uint32 = 0x546CE7AB // v0AF0E89A = v0A0C61AD

// ------------------winmain反调式-------------------
var v0A7483B4 uint32 = 8
var v0A56E4E2label1 uint32 = 0x0AD56E8A
var v0ABF8B07 uint32 = 7
var v09FC0B5Alabel2 uint32 = 0x0A8FE9A6
var v0AA0D71B uint32 = 6
var v0A557660label3 uint32 = 0x0B10933D

var v0A74573D uint32 = 0x3F // 0x00111111，只做6种反调试场景检测

// 异或计算的意义是什么
var v0AF77CE1 uint32 = 0x89B8FCE2
var v0B1090BA uint32 = 0xB8B76C72
var v0A38DB63 uint32 = 0xFC24548D // v0A38DB63 = 0x310F9090 = v0AF77CE1 ^ v0B1090BA

var v09FBB49C uint32 = 0x10648B16
var v0A325B9F uint32 = 0x721655D3 // v0A325B9F = v09FBB49C

// 异常
var v0A88E351 uint32 = 0x004D7CED
var v09E2318F uint32 = 0x0AD91CED

// 异或得到"ntdll.dll"
var v09FBF62E = [12]uint8{0x3E, 0x6C, 0x0C, 0x55, 0x17, 0xF9, 0x70, 0xB8, 0x18, 0x7B, 0xFF, 0x75}
var v09F8D538 = [12]uint8{0x50, 0x18, 0x68, 0x39, 0x7B, 0xD7, 0x14, 0xD4, 0x74, 0x7B, 0x8B, 0x45}

// 异或得到"NtQueryInformationProcess"
var v0AD72DB9 = [28]uint8{}
var v0A55FF3A = [28]uint8{}

//
//
// ------------------end-----------------------

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
func f0A05D00Fmemcpy(dstAddr unsafe.Pointer, srcAddr unsafe.Pointer, size uint32) uint32 {
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
	return index
}

func f0AAB88A1getValue(imageBase uintptr, crcDatas []block, valuep unsafe.Pointer, step uint32, lenp *uint32) []block {
	// 0x09DEB8F1 0x0A7475A4 0x09E6AAB7 0x09E91A45 0x0A84BA1D 0x0A44AD5E
	// push ecx // 产生ebp4
	// push ebx
	// push esi
	// push edi
	// edi = crcDatas
	ebp4 := uintptr(valuep) + uintptr(step) // valuep + 8
	if valuep == unsafe.Pointer(uintptr(valuep)+uintptr(step)) {
		// 0x0AF11712
		return crcDatas // eax = edi
	}
	// 0x0A563538
	// ebx = lenp
	for {
		// 0x0A745BAE
		// ecx = [ebx] // ecx = *lenp
		// 0x0AD362DE
		addr := crcDatas[0].addr // 0x1000
		size := crcDatas[0].size // 0x46
		addr += uintptr(*lenp)   // 0x1000, 0x1004

		// 0x0A970D30
		addr += imageBase // 0x00401000, 0x00401004
		size -= *lenp     // 0x46, 0x42

		// 0x0AD98EF4 0x0AF77208 0x0AF7875F
		if step == size {
			// 0x0AFD7E1C 0x0AF12435 0x0A931011 0x0AF975AE 0x0A38CEAE 0x0B071867 0x0A558DE4
			// 0x0B06F488
			// 0x0A05D00F
			f0A05D00Fmemcpy(valuep, unsafe.Pointer(addr), size)

			// 0x0AF8F334
			valuep = unsafe.Pointer(uintptr(valuep) + uintptr(size)) // 0x0018FF58 + 4
			crcDatas := crcDatas[1:]
			if crcDatas[0].addr == 0xFFFFFFFF {
				// 0x0AF11712
				return crcDatas
			}
			// 0x0A43CE90
			*lenp = 0

			// 0x0ABE3F38 // 不需要再判断凑整
		} else if step > size {
			// 0x0ABE03B6
			// push 0x0A440234
			// push 0x0A05D00F
			// ret

			// 0x0A05D00F
			f0A05D00Fmemcpy(valuep, unsafe.Pointer(addr), size)

			// 0x0A440234
			valuep = unsafe.Pointer(uintptr(valuep) + uintptr(size)) // 0x0018FF58 + 2
			*lenp = 0
			step -= size
			crcDatas := crcDatas[1:]
			if crcDatas[0].size == 0xFFFFFFFF { // 最后不足部分用0凑
				// 0x0AFDFFD3 0x0AF7CDF7
				// 0x0A83C5D1 memset
				func(p unsafe.Pointer, value uint32, size uint32) {
					// 0x0A55A765 0x09FC7C50 0x09EB9524 0x09F8436A
					if step > 0 {
						// 0x0A38CF80 0x0A9F8454 0x0A88BF90 0x0AA05578 0x0A437662 0x0A38C4F4 0x09FD1F88
						// 0x0A842745 0x0A6035A3 0x09E270C3 0x0AF7A605
					}
					// 0x0AD70014
				}(valuep, 0, step)

				// 0x09FDCDCA
				return crcDatas // 0x0AF11712
			}
			// 0x0ABE3F38
		} else {
			// 0x0A4EA4E2
			// push 0x0A32E5EE
			// push 0x0A05D00F
			// ret

			// 0x0A05D00F
			nSize := f0A05D00Fmemcpy(valuep, unsafe.Pointer(addr), step)

			// 0x0A32E5EE
			valuep = unsafe.Pointer(uintptr(valuep) + uintptr(nSize)) // 0x0018FF58 + 4
			*lenp += nSize                                            // 4, 8, ... or 8, 16, ...
			// 0x0ABE3F38
		}

		// 0x0ABE3F38
		if uintptr(unsafe.Pointer(valuep)) == ebp4 { // 正常都是0x0018FF5C，否则就要循环读取下一组校验数据来凑满4个字节
			// 0x0A931F7A
			// 0x0AF11712
			return crcDatas // eax = edi
		}
	} // for loop 0x0A745BAE
}

// 壳逻辑
func f006CD259securityInitCookie() {
	// shell logic
	// ---------------------section ERWC--------------------
	// push eax
	// push ecx
	// push esi
	// pushfd
	// push edx
	// push ebx
	// push edi

	// 0x00705113
	// 0x007269EA, f007269EA, 设置所有段属性为ERWC
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

		// 修改text段，rdata段，data段，rsrc段访问权限
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

	// -----------------crc validate--------------
	// 0x006CD253 0x0AFDF647 0x0A048F52 0x0AFDE4A1 0x0A43CFF3 0x0A930692 0x09F84632
	// push eax
	// push ecx
	// push esi
	// pushfd
	// push edx
	// push ebx
	// push edi

	// 0x0AD9893F
	var label1 uint32 = 0x0A56ED4A
	// push 0x0A56ED4A
	// push 0x0AA31853
	// ret

	// 0x0AA31853 crc校验
	func() {
		// 20个字节的局部变量
		ebp4imageBase := v0A327567imageBase
		if v0A0032DB == v09E2BC68 {
			// 0x0AFDC089
			return // 0x0A56ED4A
		}
		// 0x09FD9EB3
		// push ebx
		// push esi
		v0A8FE0E0unused ^= 0               // 0x12A225EC
		v0A746588sum ^= 0                  // 0xF7598681
		crcDatas := v0AAB950CcrcDataSet[:] // esi
		var ebp8len uint32
		v0AF77FE4sumcalc = 0
		var ebpCvalue uint32
		var ebp14value uint64

		// 计算密文完整性校验和
		for {
			if crcDatas[0].addr == 0xFFFFFFFF {
				break // 0x0A55F42F
			}
			// 0x0A8466A8 0x0A4E17E5 0x0ABF738E 0x0A889E57
			// push ebx 是参数???
			// push &ebp8len
			// push 4
			// push &ebpCvalue
			// push crcDatas
			// push ebp4imageBase

			// 0x0AD36283
			// push 0x0A0C4D0F
			// push 0x0AAB88A1
			// ret

			// 0x0AAB88A1
			crcDatas = f0AAB88A1getValue(ebp4imageBase, crcDatas, unsafe.Pointer(&ebpCvalue), 4, &ebp8len)

			// 0x0A0C4D0F
			v0AF77FE4sumcalc += ebpCvalue
		} // for loop

		// 0x0A55F42F，验证密文完整性
		// 0x0AFDEC04 0x0A440CC2 0x0ABD4564 0x0AD33BA2 0x0A84AE2E 0x0A43E793
		if v0AF77FE4sumcalc == v0A746588sum {
			// 0x09FC934F
			// push 0x09EBB9CC
			// push ebp4imageBase

			// 0x0A5508D5
			// push 0x09EAC485
			// push 0x0B286D7E
			// ret

			// 0x0B286D7E 验证段长度
			nRet := func(imageBase uintptr, sectionInfos []sectionInfo) int {
				fileHeader := (*pe.FileHeader)(unsafe.Pointer(imageBase + uintptr(*(*uint32)(unsafe.Pointer(imageBase + 0x3C))) + 4)) // IMAGE_NT_HEADERS.IMAGE_FILE_HEADER, 0x00400144
				if fileHeader == nil {
					// 0x0A8FC9CD
					// 有问题，或者非法内存访问
					// return (*[100]pe.SectionHeader32)(unsafe.Pointer(fileHeader))[:0]
				}
				// 0x0A83C137
				imageBase = 0
				// 0x0AD73ED5
				eaxSectionHeader := (*[100]pe.SectionHeader32)(unsafe.Pointer((uintptr(unsafe.Pointer(fileHeader)) + unsafe.Sizeof(pe.FileHeader{}) + unsafe.Sizeof(pe.OptionalHeader32{}))))[:]
				if fileHeader.NumberOfSections > 0 {
					// 0x09EAED20
					for {
						// 0x0A83FA73
						if sectionInfos[0].addr&sectionInfos[0].flag1 == 0xFFFFFFFF { // 判断是否需要校验
							// 0x0AA090C3
							// 0x0ABB6AFC
							return 0
						}
						// 0x0A55B8FD
						if eaxSectionHeader[0].VirtualSize != sectionInfos[0].size { // 校验段大小
							// 0x0AA090C3
							// 0x0ABB6AFC
							return 0
						}
						// 0x0AAB49F0
						if sectionInfos[0].flag2 != 0 {
							// 0x0AA090C3
							// 0x0ABB6AFC
							return 0
						}
						// 0x0A558FFD
						eaxSectionHeader = eaxSectionHeader[1:]
						imageBase++
						sectionInfos = sectionInfos[1:]

						// 0x0AF789C1
						if imageBase >= uintptr(fileHeader.NumberOfSections) {
							break // 0x0AA2BBD5
						}
					} // for loop 0x0A83FA73
				}
				// 0x0AA2BBD5
				if sectionInfos[0].addr&sectionInfos[0].flag1 != 0xFFFFFFFF {
					// 0x0AA090C3
					// 0x0ABB6AFC
					return 0
				}
				// 0x0A564B2F
				// 0x0ABB6AFC
				return 1
			}(ebp4imageBase, v09EBB9CCsectionInfoSet[:])

			// 0x09EAC485, 解密过程
			if nRet != 0 {
				// 0x0A84B9A6
				// 0x0A56CFD4
				crcDatas = v0AAB950CcrcDataSet[:] // esi
				ebp8len = 0
				v0A0032DB = v09E2BC68
				for {
					// 0x0AD80743 0x0A55799E 0x0A84064A 0x0A12E533 0x0A335558 0x0AA07FD0
					if crcDatas[0].addr == 0xFFFFFFFF {
						break // 0x09FB8D2B
					}
					// 0x09FD2C99
					// push 0 // ?
					ebpCvalue = ebp8len
					// ebx := crcDatas // 为后续写明文做准备
					// push 0x0A9527D0
					// push 0x0AAB88A1
					// ret

					// 0x0AAB88A1, 每次读取8字节密文
					crcDatas = f0AAB88A1getValue(ebp4imageBase, crcDatas, unsafe.Pointer(&ebp14value), 8, &ebp8len)

					// 0x0A9527D0
					// push 0x0AD70D9F
					// push 0x0A9003B1
					// ret

					// 0x0A9003B1, 计算密钥
					func() {
						eax := v0A330A7ELo
						ecx := v0A330A82Hi
						var edx uint32
						var esi uint32 = 0x40
						for {
							// 0x0AD80CB9
							eax += ((ecx>>5 ^ ecx<<4) + ecx) ^ (v0A5578FA[edx%4] + edx)
							edx -= 0x61C88647
							ecx += ((eax>>5 ^ eax<<4) + eax) ^ (v0A5578FA[(edx>>11)%4] + edx)
							esi--
							if esi == 0 {
								break // 0x0A05177A
							}
						}

						// 0x0A05177A
						v0A330A7ELo = eax
						v0A330A82Hi = ecx
					}()

					// 0x0AD70D9F, 解密
					ebp14value = uint64(uint32(ebp14value) ^ v0A330A7ELo + (uint32(ebp14value>>32)^v0A330A82Hi)<<32)
					// push &ebpCvalue，当长度用了
					// push ebx
					// push 8
					// push &ebp14value
					// push ebp4imageBase
					// push 0x0AFD8B4F
					// push 0x09FE64D2
					// ret

					// 0x09FE64D2, 写明文
					func(imageBase uintptr, valuep unsafe.Pointer, step uint32, crcDatas []block, lenp *uint32) {
						var ebp4 uint32
						var ebp8 uintptr
						ebp8 = uintptr(valuep) + uintptr(step)
						if valuep != unsafe.Pointer(uintptr(valuep)+uintptr(step)) {
							return // 0x0AD97243
						}
						// 0x0AD7BE4D 0x0AF7C50F 0x0AD9274B
						// push ebx
						// ebx := crcDatas
						// push esi
						// push edi
						for {
							// 0x0AF97BBD 0x0A9D32BC 0x0A604C8D 0x0AFDAA70
							addr := crcDatas[0].addr
							size := crcDatas[0].size
							addr += uintptr(*lenp)
							addr += imageBase
							size -= *lenp
							// 0x0A057A29 0x0AF8A955 0x0A84973E 0x09FDF556
							if size <= step {
								// 0x0A0023F3
								// push size
								// push valuep
								// push addr
								// push 0x0AAB911B
								// push 0x0A05D00F
								// ret

								// 0x0A05D00F
								f0A05D00Fmemcpy(unsafe.Pointer(uintptr(addr)), valuep, size)

								// 0x0AAB911B
								ebp4 += size
								crcDatas = crcDatas[1:]
								if crcDatas[0].addr == 0xFFFFFFFF {
									// 0x0A0485AD 0x0AF7ED9E
									return // 0x0AD97243
								}
								// 0x0A4ED146
								*lenp = 0
								// 0x09E2D4D0
							} else {
								// 0x09FC481F
								// push step
								// push valuep
								// push addr
								// push 0x09FE2A75
								// push 0x0A05D00F
								// ret

								// 0x0A05D00F
								f0A05D00Fmemcpy(unsafe.Pointer(uintptr(addr)), valuep, step)

								// 0x09FE2A75
								ebp4 += step
								// 0x0A95081C
								*lenp += step
								// 0x0AD310AF 0x09E2D4D0
							}
							// 0x09E2D4D0 0x0A131034
							if uintptr(ebp4)+uintptr(valuep) == ebp8 {
								// 0x0A0485AD 0x0AF7ED9E
								return // 0x0AD97243
							}
							// 0x0A5FC621
							// 0x0AF97BBD
						} // for loop 0x0AF97BBD
					}(ebp4imageBase, unsafe.Pointer(&ebp14value), 8, crcDatas, &ebp8len)

					// 0x0AFD8B4F
					// 0x0AFD8B52 等效于 0x0AD80743
				} // for loop 0x0AD80743

				// 0x09FB8D2B
				v0A5578FA[0] ^= v0B287362
				// 0x09FC0180
				// 0x0ABDAB3B
			}
		} // if v0AF77FE4sumcalc == v0A746588sum
		// 0x0AF7E275，完整性验证失败
		ebpCvalue = uint32(v0A327567imageBase) + v0A88F1F8 // 0x00400000 + 0x0A9F55F3 = 0x0ADF55F3
		label1 = ebpCvalue
		// 0x0ABDAB3B
		// 0x0ABDAB3B 与 0x09FD9EB3 对应
		// pop esi
		// pop ebx
		return // 0x0AFDC089
	}()

	// 0x0ADF55F3，完整性验证失败，异常

	// 0x0A56ED4A，解密成功
	// pop edi
	// pop ebx
	// pop edx
	// popfd
	// pop esi
	// pop ecx
	// pop eax

	// --------------validate read and write----------
	// 0x0A33348C
	// 0x09F8A8BD
	// 0x0B10B013
	// push eax
	// push ecx
	// push esi
	// pushfd
	// push edx
	// push ebx
	// push edi

	// 0x0A9F55F3
	// 0x0AD2D162, f0AD2D162, 验证读写权限？
	func() {
		// 0x24字节局部变量
		if v09E8FF96 == 0 {
			//0x0AF7FFD5
			return
		}
		// 0x0A896CF3
		ebp24 := struct {
			dummyName  uint32
			dwPageSize uint32
			data       [24]uint8
		}{}
		ebp24.dwPageSize = 0x1000 // win.GetSystemInfo(&ebp24)

		// 0x0A9F5118
		esi := ebp24.dwPageSize
		var eaxImageBase uint32 = 0x0AD2D162
		eaxImageBase -= v0A049D16 // 0
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
			ecxCodeBase += esi // 一页为一个单位
		}

		// 0x0AF98284
		v09E8FF96 = 0
		// 0x0AF7FFD5
		return
	}()

	// 0x0A9F55F8 0x0ABD82B9
	// pop edi
	// pop ebx
	// pop edx
	// popfd
	// pop esi
	// pop ecx
	// pop eax

	// ---------------iat move------------
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
	// 0x0A44A7FD, f0A44A7FD, IAT move
	func() {
		// 0xC字节的布局变量
		if v0A563A89 == 0 {
			// 0x0A8429DF
			return
		}
		// 0x0A4E8FC1
		ebp4iats := v09F8682AiatSet[:]
		ecxIATaddr := ebp4iats[0].addr + 0x00400000 // 0x09DE917C advapi32.&CryptGetUserKey
		ebxIATsize := ebp4iats[0].size
		ebp8rdatas := v0A3A1688rdataSet[:]
		esiRDATAaddr := ebp8rdatas[0].addr + 0x00400000 // 0x01149000
		ediRDATAsize := ebp8rdatas[0].size
		if v09F8682AiatSet[0].addr != 0xFFFFFFFF {
			// 0x09EAE30B
			ebpC := v0A3A1688rdataSet[0].addr
			for {
				// 0x0ABE28CC
				if ebpC == 0x0FFFFFFFF {
					break // 0x0A56CD8C
				}
				// 0x0A9F96DB
				*(*uint32)(unsafe.Pointer(esiRDATAaddr)) = *(*uint32)(unsafe.Pointer(ecxIATaddr)) // 0x75AF350E advapi32.CryptGetUserKey
				var eax uint32 = 4
				esiRDATAaddr += uintptr(eax) // 0x01149000, 0x01149004
				ecxIATaddr += uintptr(eax)   // 0x09DE917C, 0x09DE9180
				ebxIATsize -= eax            // 0x1F0, 0x1EC
				ediRDATAsize -= eax          // 0x09C0, 0x09BC
				if ebxIATsize == 0 {         // kernel32.dll
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
					break // 0x0A56CD8C
				}
			}
		}
		// 0x0A56CD8C
		v0A563A89 = 0
		// 0x0A8429DF
		return
	}()

	// 0x0AFD8512
	// 0x0AAB829D
	// popfd
	// pop esi
	// pop ebx
	// pop edi
	// pop edx
	// pop ecx

	// -------------text move-----------------
	// 0x0AF8BBBA 0x0A60483B 0x0A56C231
	var label11 uint32 = 0x09FB8C65
	// 0x0AD323B3, text move
	var label2 uint32 = 0x0A56DEAE
	// 0x09FFDE30
	var label3 uint32 = 0x00DDD650
	// 0x0A0509F4
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

	// 0x0AF7CDBD
	func() {
		// push esi
		label11 = v0ABE3A4F // [esp+v0ABDABC7*4+8] = v0ABE3A4F
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
				// 0x0A43E504 0x0AD3B4F7 0x09E2F54F 0x09EBC1C0 0x0AF126B7
				if eaxAddr == 0xFFFFFFFE {
					// 0x0A8FC82C
					// editexts = editexts[0].size + v0AFE0B1FimageBase，无意义
					// 0x09E91166
					// 0x09FE278E
				} else {
					// 0x0A4416F2 0x0A9F9C72
					eaxAddr += v0AFE0B1FimageBase
					// 0x0AD738B4
					edxSize := editexts[0].size
					// 0x0A9F6726 0x0A936A10 0x0AFDD768
					if edxSize > 3 {
						// 0x09FB7EC3
						ecx := edxSize / 4
						edxSize %= 4
						// 商迭代
						for {
							// 0x0AD84E83
							*(*uint32)(unsafe.Pointer(eaxAddr)) = v0A84CEC1 - *(*uint32)(unsafe.Pointer((&esiops[0])))
							// 0x0AA2F6ED
							eaxAddr += 4
							esiops = esiops[4:]
							ecx--
							// 0x0A8417D3 0x0ABE42C0 0x0A041758 0x0ABFA56D 0x0B284F49 0x0A55F76C 0x0AB4FADB 0x0A900853
							if ecx == 0 {
								break // // 0x0AF0FE48
							}
						} // for loop 0x0AD84E83
					}
					// 0x0AF0FE48
					// 0x0A603BC5, 余数迭代
					for edxSize > 0 {
						// 0x0A4EAF80
						*(*uint8)(unsafe.Pointer(eaxAddr)) = uint8(v0A84CEC1) - *(*uint8)(unsafe.Pointer((&esiops[0])))
						eaxAddr++
						esiops = esiops[1:]
						edxSize--
					} // for loop 0x0A4EAF80

					// 0x09F88775, next
					editexts = editexts[1:]
				}
				// 0x09FE278E 0x0A7462B5 0x0B10D023 0x09EB7BBB
				if editexts[0].addr == 0xFFFFFFFF {
					break // 0x09FB5670
				}
			} // for loop 0x0AAB316A
			// 0x09FB5670
			// pop ebp
			// pop ebx
		}
		// 0x0A89291D
		label2 = v0A5557A4 // [esp+v0A5D6CCF*4+12] = v0A5557A4
		// pop edi
		v0AF0E89A = v0A0C61AD //
		label3 = v0AC33EA1    // [esp+v0A0526C0*4+8] = v0AC33EA1
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
	// label3 (0x0AF1023E 0x09E2FFC0)
	// label2 (0x0AD3B521)
	// label11 (0x0A39221A)

	// 0x00E10691, 计算 v012F7B90 和 v012F7B94
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

var v09D9DB10env []string
var v09D9E18CfileName string
var v09DA37A0heapNum int
var v09DA37ACcmdline string

// --------------------------------------------------------------------
// f00D16821, runtime_main, __winMainCRTStartup, en
// f00DF478C, runtime_main, __winMainCRTStartup, zh
func main() {
	// f00DFD850(&v012AC0A8, 0x58)
	// ebp4 := 0

	// 0x00DF479D: startup info
	var ebp68si struct {
		// DWORD  cb;
		// DWORD  cb;
		// LPSTR  lpReserved;
		// LPSTR  lpDesktop;
		// LPSTR  lpTitle;
		// DWORD  dwX;
		// DWORD  dwY;
		// DWORD  dwXSize;
		// DWORD  dwYSize;
		// DWORD  dwXCountChars;
		// DWORD  dwYCountChars;
		// DWORD  dwFillAttribute;
		dwFlags     int   // DWORD  dwFlags;
		wShowWindow int16 // WORD   wShowWindow;
		// WORD   cbReserved2;
		// LPBYTE lpReserved2;
		// HANDLE hStdInput;
		// HANDLE hStdOutput;
		// HANDLE hStdError;
	}

	// dll.kernel32.GetStartupInfo(&ebp68)
	// ebp4 = -2

	// 0x00DF47AD: check pe
	// ...

	// 0x00DF47FA: heap init
	f00DFF074heapInit := func(safe bool) int {
		// flag := !safe
		// v09D9DEB8heap1 = kernel32.HeapCreate(flag, 4096, 0)
		v09DA37A0heapNum = 1
		return v09DA37A0heapNum
	}
	if f00DFF074heapInit(true) == 0 {
		// f00DF471C(0x1C)
	}

	// 0x00DF480C: fls init
	if f00DFC5ACflsInit() == false {
		// f00DF471C(0x10)
	}

	// 0x00DF481D:
	// f00E0A38C()
	// ebp4 = 1
	// if f00E00A78() < 0 {
	// 	f00DF0407(0x1B)
	// }

	// 0x00DF4836: command line
	// v09DA37ACcmdline = dll.kernel32.GetCommandLine()
	v09DA37ACcmdline = `"\\path\\main.exe" connect /u192.168.0.102 /p444405`

	// 0x00DF4841: environment varible
	v09D9DB10env = []string(nil) // f00E1055AgetEnv()

	// 0x00DF484B:
	// if f00E1049F() < 0 {
	// 	f00DF0407(0x8)
	// }

	// 0x00DF485C:
	// if f00E10218() < 0 {
	// 	f00DF0407(0x9)
	// }

	// 0x00DF486D: log init
	// if err := f00DF053ElogInit(); err != 0 {
	// 	f00DF0407(err)
	// }

	// 0x00DF487F:
	szcmdLine := strings.Join(os.Args[1:], " ") // f00E101B9
	show := 10                                  /*win.SW_SHOWDEFAULT*/
	if ebp68si.dwFlags&1 /*STARTF_USESHOWWINDOW*/ != 0 {
		show = int(ebp68si.wShowWindow)
	}
	f004D7CE5winMain(0x00400000, 0, szcmdLine, show)
}

// --------------------------------------------------------------------
var v00DF490A [0x5B]uint8
var v0B2BE8D6 [0x19]uint8
var v0B2BE90B [0x58]uint8

// OEP.shell logic
func f00DF490AwinMainCRTStartup() {
	// vp := win.GetProcAddress(win.LoadLibrary("kernel32.dll"), "virtualProtect")
	// var oldProtect uint32
	// vp(0x00DF490A, 0x5B, 0x40, &oldProtect) // PAGE_EXECUTE_READWRITE
	// vp(0x0B2BE8D6, 0x19, 0x40, &oldProtect) // PAGE_EXECUTE_READWRITE
	// jmp f0B2BE8D6copy
}

// OEP.shell copy
func f0B2BE8D6copy() {
	// f00DE7C90memcpy(v00DF490A[:], v0B2BE90B[:], len(v00DF490A))
	// jmp f00DF490AwinMainCRTStartup
}

// OEP.real logic, __winMainCRTStartup
func f0B2BE90BwinMainCRTStartup() {
	f006CD259securityInitCookie() // -GS选项提供
	// jmp f00DF478C ;0x00DF478C, hard hook as E9 6CA04C0A/jmp 0x0B2BE980 to load main.dll
	// 0x0B2BE980: f0B2BE980
	// if false == kernel32.LoadLibrary("main.dll") {
	// 	code := kernel32.GetLastError()
	// 	func(code int) {
	// 		var buf [32]uint8
	// 		heapInit(true)
	// 		flsInit()
	// 		sprintf(buf[:], "getlasterror %d", code)
	// 	}(code)
	// 	dll.kernel32.ExitProcess(0)
	// }
	// jmp f00DF478C ;main/runtime_main
}

// os创建主线程并执行一段逻辑再跳转到sysep
func oslogic() {
	func() {
		func() {
			func() {
				func() {
					func() {
						// ntdll system entrypoint
					}()
				}()
			}()
		}()
		f00DF490AwinMainCRTStartup()
	}()
}
