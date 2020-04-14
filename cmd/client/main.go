package main

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/xujintao/balgass/win"
)

type conf struct {
	a uint32
	b uint32
}

const (
	v0114D912 uint8  = 0
	v0114D913 string = ""
	v0114DB18 string = "서버와의 접속종료." // 终止与服务器的连接
	v0114DD64 uint16 = 0x0061
)

// 这里使用数组，确保在data段，固定在某个地址上
// 如果使用string，那么读文件后对buf进行string类型转换，会memmove到堆上
var v012E10F4 [3]uint8 // FC CF AB
var v012E2210 uint32 = 1
var v012E2214 = "MUCN"
var v012E2224 uint32 = 1
var v012E2228 uint32 = 16
var v012E222C *os.File
var v012E2338ip = "192.168.0.102"
var v012E233Cport uint16 = 44405
var v012E2340 int

var v012E3F08 struct {
	height uint32 // 0x258, 600
	width  uint32 // 0x320, 800
}
var v012E4018version = [8]uint8{'2', '2', '7', '8', '9', 0, 0, 0}                                               // "22789"
var v012E4020serial = [16]uint8{'M', '7', 'B', '4', 'V', 'M', '4', 'C', '5', 'i', '8', 'B', 'C', '4', '9', 'b'} // "M7B4VM4C5i8BC49b"
var v012F7910 os.File

var v01319A38version [8]uint8 // "1.04R+"	// 可执行程序版本号
var v01319A44version [8]uint8 // "1.04.44" // 配置文件版本号
var v01319A50ip [16]uint8

var v01319D1CwndProc uintptr
var v01319D65 uint32
var v01319D68 *t1319D68
var v01319D6ChWnd win.HWND // hWnd, 0x0073,1366
var v01319D70hModule win.HMODULE
var v01319D74hDC win.HDC     // hDC, 0x1401,11A6
var v01319D78hGLRC win.HGLRC // hGLRC, OpenGL rendering context, 0x0001,0000
var v01319D90 int
var v01319D95 uint8
var v01319D98 int
var v01319D9C int
var v01319DA0 int
var v01319DB4 int
var v01319DB8 [64]uint8
var v01319DFC = "Chs"
var v01319E00 = "chinese"

type version struct {
	major  uint16 // ebp_244
	minor  uint16 // ebp_242
	patch  uint16 // ebp_240
	fourth uint16 // ebp_23E
}

func f0043B750xor(buf []uint8, len int) {
	ebp4 := 0
	for ebp4 < len {
		// mov eax, dowrd ptr ss:[ebp-4]
		// cdq ;Convert Double to Quad 把EAX的第31bit复制到EDX的每一个bit上
		// push 3
		// pop ecx
		// idiv ecx ;idiv为有符号除法, div为无符号除法
		// 如果除数为8位，则被除数为16位，则结果的商存放与al中，余数存放于ah中。
		// 如果除数为16位，则被除数为32位，则结果的商存放与ax中，余数存放于dx中。
		// 如果除数为32位，则被除数为64位，则结果的商存放与eax中，余数存放于edx中
		buf[ebp4] ^= v012E10F4[ebp4%3]
		ebp4++
	}
}

func f004D52AB(fileName []uint8, cmd string) bool {
	// 从命令行字符串中提取出应用程序名，也就是main.exe，存到fileName数组中
	return true
}

func f004D5368() bool {
	var fileName [260]uint8 // ebp-110
	cmd := os.Args[0]       // ebp-4
	f004D52AB(fileName[:], cmd)

	// dwDesiredAccess = 1<<31, GENERIC_READ
	// dwShareMode = 3, FILE_SHARE_READ | FILE_SHARE_WRITE
	// psa = 0, default security and cannot be inherited by any child
	// dwCreationDisposition = 3, OPEN_EXISTING
	// dwFlagsAndAttributes = 0x80, FILE_ATTRIBUTE_NORMAL
	// dwFileTemplate = 0, When opening an existing file, CreateFile ignores this parameter
	// hFile := CreateFile(fileName[:], 1<<31, 3, 0, 3, 0x80, 0)
	file, err := os.Open(string(fileName[:]))
	if err != nil {
		return false
	}
	v012E222C = file
	return true
}

func f004D55C6(fileName []uint8, ver *version) bool {
	// 获取可执行文件版本
	// GetFileVersionInfoSize()
	// GetFileVersionInfo()
	// VerQueryValue()

	// fileVersion: 1.4.44.0
	ver.major = 0x01
	ver.minor = 0x04
	ver.patch = 0x2C
	ver.fourth = 0
	return true
}

// old callback
func f004D5F98wndProc(hWnd win.HWND, message uint32, wParam, lParam uintptr) uintptr {
	switch message {
	case win.WM_CREATE: // 1
	case win.WM_TIMER: // 0x113
		switch wParam {
		case 0x3E8:
			// f00B4C239
			func() {
				if v01319D68 == nil {
					return
				}
				v01319D68.f00B4CC0D()
				// ...
			}()
			f004D8FF5() // 心跳
		}
	case win.WM_USER: // 0x400
		switch lParam {
		case 1:
			v08C88FF0conn.f006BDA03read()
		case 2:
			v08C88FF0conn.f006BD945write()
		case 20:
		}
	default:
		// v08C88BC0 = 0
		// if v08C88BC5 == 1 && (v01319DA8 != v08C88C80 || v01319DAC != v08C88C80) {
		// 	v08C88BC5 = 0
		// }
		return win.DefWindowProc(hWnd, message, wParam, lParam) // 原始是NtdllDefWindowProc，这是什么？
	}
	return 0
}

// new callback
func f004D6C2BwndProc(hWnd win.HWND, message uint32, wParam, lParam uintptr) int {
	switch message {
	case win.WM_CHAR: // 0x102
		// f00A49798(wParam, lParam)
	default:
		win.CallWindowProc(v01319D1CwndProc, hWnd, message, wParam, lParam)
	}
	return 0
}

func f004D6D64() bool {
	// 40+4个字节局部变量

	pixelfd := win.PIXELFORMATDESCRIPTOR{} // 已经清零了
	// f00DE8100memset(&pixelfd, 0, 0x28) // 清零

	v01319D74hDC = win.GetDC(v01319D6ChWnd) // 从hWnd中拿到hDC
	if v01319D74hDC == 0 {
		// var nRet int32 = GetLastError()
		// v01319E08log.f00B38AE4printf("OpenGL Get DC Error - ErrorCode : %d\n\r", nRet)
		// _004D51C8()
		// nRet = _00436DA8(4, "OpenGL Get DC Error.", 30)
		// MessageBoxA(0, nRet)
		return false
	}

	pfdindex := win.ChoosePixelFormat(v01319D74hDC, &pixelfd) // return 4
	if pfdindex == 0 {
		// var nRet int32 = GetLastError()
		// v01319E08log.f00B38AE4printf("OpenGL Choose Pixel Format Error - ErrorCode : %d\n\r", nRet)
		// _004D51C8()
		// nRet = _00436DA8(4, "OpenGL Choose Pixel Format Error.", 30)
		// MessageBoxA(0, nRet)
		return false
	}

	if !win.SetPixelFormat(v01319D74hDC, pfdindex, &pixelfd) { // return true
		// var nRet int32 = GetLastError()
		// v01319E08log.f00B38AE4printf("OpenGL Set Pixel Format Error - ErrorCode : %d\n\r", nRet)
		// _004D51C8()
		// nRet = _00436DA8(4, "OpenGL Set Pixel Format Error.", 30)
		// MessageBoxA(0, nRet)
		return false
	}

	v01319D78hGLRC = win.WglCreateContext(v01319D74hDC) // return 0x0001,0000
	if v01319D78hGLRC == 0 {
		// OpenGL Create Context Error
		return false
	}

	if !win.WglMakeCurrent(v01319D74hDC, v01319D78hGLRC) { // return true
		// OpenGL Make Current Error
		return false
	}

	// ShowWindow(v01319D6ChWnd, 5)
	// SetForegroundWindow(v01319D6ChWnd)
	// SetFocus(v01319D6ChWnd)

	return true
}

func f004D6F82initWindow(hModule win.HMODULE, iCmdShow int) win.HWND {

	// 这里windows那边使用的是数组，编译器会使用movsw和movsb来给字符数组赋值
	// var ebp_78 string = "MU"
	ebp78, _ := syscall.UTF16FromString("MU")

	// ...

	// 140字节局部变量
	ebp34 := win.WNDCLASSEX{
		CbSize:        0x30, // ebp-34
		Style:         0x2B, // CS_OWNDC | CS_DBLCLKS | CS_HREDRAW | CS_VREDRAW
		LpfnWndProc:   syscall.NewCallback(f004D5F98wndProc),
		CbClsExtra:    0,
		CbWndExtra:    0,
		HInstance:     win.GetModuleHandle(nil),
		HIcon:         win.LoadIcon(0, nil),              // IDI_APPLICATION)
		HCursor:       win.LoadCursor(0, nil),            // IDC_ARROW
		HbrBackground: win.HBRUSH(win.GetStockObject(4)), // WHITE_BRUSH
		LpszMenuName:  nil,
		LpszClassName: &ebp78[0],  // "MU"
		HIconSm:       0x141F0439, // ebp-8
	}
	win.RegisterClassEx(&ebp34)

	// ...

	// CreateWindow会发大概3条消息给WndProc
	ebp4 := win.CreateWindowEx(
		0,
		&ebp78[0],  // "MU"
		&ebp78[0],  // "MU"
		0x02CA0000, // WS_CLIPCHILDREN | WS_CAPTION | WS_MINIMIZEBOX | WS_SYSMENU
		0x118,
		0x46,
		0x326,
		0x274,
		0,
		0,
		0x00400000,
		nil,
	)
	v01319D1CwndProc = win.SetWindowLongPtr(ebp4, win.GWL_WNDPROC, syscall.NewCallback(f004D6C2BwndProc))
	return 0
}

func f004D755F(haystack []uint8, y int, buf []uint8) bool {
	return false
}

// f004D7CE5winMain, WinMain
func f004D7CE5winMain(hModule win.HMODULE, hPrevInstance uint32, szCmdLine string, iCmdShow int) {
	// ----------------------winmain反调式-------------------
	// 0x0A05E61B
	var label1 uint32 = 0x009DCA19
	// push label1
	// 0x0ABDA0EA
	var label2 uint32 = 0x0A3349B0
	// push label2
	// 0x0A4420A9
	var label3 uint32 = 0x0A0518C4
	// push label3
	// 0x0A8485E3
	// push fd ebx edi esi edx ecx

	// 0x09EBC742
	var label4 uint32 = 0x0AD2B081
	// push label4
	// push 0x0AD91CED
	// ret

	// 0x0AD91CED, f0AD91CED
	func() {
		// 0xFC局部变量
		// push ebx esi edi
		var ebp8 uint32
		var ebp10 uint32
		var ebp14 uint32
		//var ebp18 uint32
		var ebp20 uint32
		var ebp24 uint32
		//var ebp28 uint32
		var ebp30 uint32
		var ebp34 uint32
		var ebp38 uint32
		var ebp40 uint32
		var ebp44 uint32

		label1 = v0A56E4E2label1 // [ebp+v0A7483B4*4+8] = v0A56E4E2label1
		if v0A74573D == 0 {      // 0x3F
			// 0x0AA09510
			ebp34 = 8
			ebp14 = 0xFF
			// 0x0A7AC650 0x0A83CE6F
		} else if v0A74573D > 0xFF { // 0x09FC6431 0x0AC37855 没必要再判断v0A74573D==0
			// 0x0A9F5A38 0x0A6B0AEB ...
		} else {
			// 0x0A9F5A3E 0x0AA2F062 0x0A32AFCE
			ebp14 = v0A74573D
			ebp44 = v0A74573D
			for {
				// 0x0A058220 0x0AAB85DB 0x0AF77F6A 0x0AD80BB2 0x0A391FF8 0x0AD3B8F3 0x0A3927B1 0x09EBC1FB
				if ebp44 <= 0 { // 迭代结束, ebp44=0, ebp34=6
					// 0x0AF814F7
					break // 0x0A83CE6F
				}
				// 0x0A442AC6 0x0AF7F68B 0x0A84320A 0x0A8902CD 0x0B071558 0x0AD98514 0x0ABF87EB
				if ebp44&1 != 0 {
					// 0x09E91EC8
					ebp34++ // 1,2,3,4,5,6
				}
				// 0x09FDAD42 0x0AC3847E
				ebp44 >>= 1 // 31,15,7,3,1,0
				// 0x0A058220
			}
		}
		// 0x0A83CE6F 0x0A38F2A7 0x0AF11D82 0x0A565A51
		if ebp34 == 0 {
			// 0x0A8FA42A 0x0A6B0AEB
		}
		// 0x0A84D90D
		ebp38 = 0
		// 0x0A32AF58
		v0A38DB63 = v0AF77CE1 ^ v0B1090BA
		// 0x0A38DB63
		// rdtsc
		var tscLeax uint32 = 0x5274DC2C
		//var tscHedx uint32 = 0x00010865
		// 0x0A32F0F4 0x0AF8EBF1
		ebp20 = tscLeax | (ebp20 >> 16) // 0x5274DE7C
		// 0x0A56E64F 0x0ABE0B07
		ebp20 %= 0x1F // 构造一个[1,30]之间的随机数, 30
		ebp8 = 1
		for {
			// 0x0A9FBE5C 0x0AD83B40 0x09E240E7 0x09FE52D1 0x09FCAF92 0x0B10FA31 0x09E03B4A
			if ebp8&ebp14 != 0 { // 1&0x3F
				// 0x0AD98F0D, 计算ebp8
				ebp40 = ebp20 % ebp34 // 30%6
				for {
					// 0x0AD75B27
					if ebp40 <= 0 {
						break // 0x0ABB3DF5
					}
					// 0x0AF82F72
					ebp8 <<= 1 // 2,4,8,16,32
					// 0x0A3A7943
					if ebp8&ebp14 == 0 { // 2&0x3F, 不超过6次
						// 0x0A049688
					}
					// 0x0AC39B0A 0x09E8EB25
					ebp40-- // 4,3,2,1,0
				} // for loop 0x0AD75B27

				// 0x0ABB3DF5 0x0A742C83 0x0A3360C0 0x0A90096C 0x0A4DF380 0x0A8892C0
				if ebp8&1 != 0 { // 情形1
					// 0x0ABFA5B3 0x0AA0932E
					// push eax
					// pushf
					// eax = [esp]
					// ebp10 = 0x00010202
					// popf
					// 0x0AA3115D
					// pop eax
					// 0x0A55B75A 0x09E0403A 0x0A3A90E4 0x09E28B74 0x0AA2C54F 0x0AD3151D 0x0A39070B
					if ebp10&0x100 != 0 { // TF位(单步跟踪)
						// 0x0A12E17A 0x0A6B0AEB
					}
					// 0x0A7ADFFA 0x0AF872E5
					label2 = v09FC0B5Alabel2 // [ebp+v0ABF8B07*4+8]=v09FC0B5A
					// 0x0A4DE5C1
					ebp38++
					// 0x0AF81189
					ebp8 &= ^uint32(1) // 清零
					// 0x0A4E80E3
				}
				// 0x0A4E80E3 0x0A4DDA43 0x0A05AA37 0x09FD9C5C 0x09E30794 0x0A4DDA22 0x09F8E057
				// 0x0A7441CB 0x0A84A77E
				if ebp8&2 != 0 { // 情形2
					// 0x0A742EF1 0x0A445521
					// pushad
					// push 0x0AFD709D
					// push fs[0]
					// fs[0] = esp
					// push 0
					// push 0
					// int3 ; 报0x80000003, EXCEPTION_BREAKPOINT
					// pop ebp30
					// pop ebp18
					// pop fs[0]
					// esp+=4
					// popad

					if ebp30 == 0 {
						// 0x0A5D66DF 0x0A6B0AEB
					}
					// 0x0A5625C2
					label2 = v09FC0B5Alabel2 // [ebp+v0ABF8B07*4+8]=v09FC0B5Alabel2
					ebp38++
					ebp8 &= ^uint32(2) // 清零
				}
				// 0x0AD31978
				if ebp8&4 != 0 { // 情形3
					// 0x0AFDC124
					ebp30 = 1
					// pushad
					// push 0x0AF94419 ;SEH
					// push fs[0]
					// fs[0]=esp
					// push 0
					// int1 ;单步直接stepover了，没有进入SEH异常处理，而Run是ok的报0x80000004, EXCEPTION_SINGLE_STEP
					// pop ebp30
					// pop fs[0]
					// esp+=4
					// popad

					if ebp30 == 0 {
						// 0x09F84AFD 0x0A6B0AEB
					}
					// 0x0AA30C97
					label2 = v09FC0B5Alabel2 // [ebp+v0ABF8B07*4+8]=v09FC0B5Alabel2
					ebp38++
					ebp8 &= ^uint32(4) // 清零
				}
				// 0x0A92FBD9
				if ebp8&8 != 0 { // 情形4
					// 0x09E2BB81 0x0ABFA7E5 0x0A847A61 0x0A7AD6B6 0x0A904AEE 0x0A4E7E87 0x09E25911
					ebp30 = 1
					// pushad
					// push 0x0A32BD69 ;SEH
					// push fs[0]
					// fs[0]=esp
					// push 0
					// push 0
					// int 1 ;指令是CD 01，报0xC0000005, EXCEPTION_ACCESS_VIOLATION
					// pop ebp30
					// pop ebp18
					// pop fs[0]
					// esp+=4
					// popad

					// 0x0A327A81 0x0A05C5B5 0x09FC28CA 0x0B10F776 0x0ABB564B
					if ebp30 == 0 {
						// 0x09F8768E 0x0A6B0AEB
					}
					// 0x09FC1A2B
					label2 = v09FC0B5Alabel2 // [ebp+v0ABF8B07*4+8]=v09FC0B5Alabel2
					ebp38++
					ebp8 &= ^uint32(8) // 清零
				}
				// 0x0A9353C2 0x0AF8FAB1 0x0B074337 0x0AD9201D 0x0A895B6A 0x09FD5603 0x0ABB3AC2
				// 0x0ABD65F8 0x0AAB22A0
				if ebp8&0x10 != 0 { // 情形5
					// 0x0AFD7373 0x0A892E27
					type ntTIB struct {
						SEH                  uintptr
						StackBase            uintptr
						StackLimit           uintptr
						SubSystemTIB         uintptr
						FiberData            uintptr // or Version
						ArbitraryUserPointer uintptr
						Self                 *ntTIB
					}
					type ntPEB struct {
						InheritedAddressSpace    uint8
						ReadImageFileExecOptions uint8
						BeingDebugger            uint8
						// ...
					}
					type ntTEB struct {
						tib ntTIB
						// ...
						peb *ntPEB // ProcessEnvironmentBlock
					}
					// push eax
					tebp := (*ntTEB)(unsafe.Pointer(&ntTIB{}))
					debuger := tebp.peb.BeingDebugger // IsDebuggerPresent
					// 0x0A844844
					if debuger != 0 {
						// 0x0A6B0AEB
					}
					// 0x0AA0DDAE
					label2 = v09FC0B5Alabel2 // [ebp+v0ABF8B07*4+8]=v09FC0B5Alabel2
					ebp38++
					ebp8 &= ^uint32(0x10) // 清零
				}
				// 0x0A444B3D 0x09E2F73A 0x0AD813DC 0x0AF899B3 0x0A05DA3A
				if ebp8&0x20 != 0 { // 情形6
					var ebp70Name [28]uint8
					// 0x09E6AFDD，加载ntdll.dll--------------------------------------
					ebp4C := 0
					ebp74 := ebp70Name[:]
					for {
						// 0x0AD80EA0 0x09E9120A 0x0ABFAB5E 0x0AD8044B
						ebp74[0] = v09FBF62E[ebp4C] ^ v09F8D538[ebp4C]
						ebp4C++
						// 0x0A5D39E1 0x0AD91C0D
						if ebp74[0] == 0 {
							break // 0x09E6F520，此时ebp70Name="ntdll.dll"
						}
						ebp74 = ebp74[1:]
					}
					// 0x09E6F520 0x0AF79937
					// push ebp70Name[:]
					// push 0x09FB452D
					// push &LoadLibraryA
					// ret
					var ebp28 uint32
					// ebp28 = win.LoadLibraryA(string(ebp70Name[:]))
					// 0x09FB452D
					if ebp28 == 0 {
						// 0x0A44B037 0x0A6B0AEB
					}

					// 0x09EB89D9，得到NtQueryInformationProcess函数------------------
					ebp4C = 0
					ebp78 := ebp70Name[:]
					for {
						// 0x0A4E3C47
						ebp78[0] = v0AD72DB9[ebp4C] ^ v0A55FF3A[ebp4C]
						ebp4C++
						if ebp78[0] == 0 {
							break // 0x0A847E6C
						}
						ebp78 = ebp78[1:]
					}
					// 0x0A847E6C
					// push ebp70Name[:]
					// push ebp28
					// push 0x0A8FEA70
					// push &GetProcAddress ;0x0A8FEA70
					// ret
					var ebp1CNtQueryInformationProcess uint32
					// ebp1C = win.GetProcAddress(ebp28, ebp70Name[:])
					// 0x0A8FEA70
					if ebp1CNtQueryInformationProcess == 0 {
						// 0x0AD2DAF1
						// win.FreeLibrary(ebp28)
						// 0x0A4E108E 0x0A6B0AEB
					}

					// 0x0ABD793E，使用NtQueryInformationProcess函数检测ProcessDebugPort------
					ebp7C := ebp70Name[:]
					// 0x0AF78158
					for {
						// 0x09F8EA54 0x0A04D754 0x0A5D2F76 0x0A0C414C 0x09FB9026 0x0A848CF2 0x0AD311C1
						if ebp7C[0] == 0 {
							break // 0x0A9038FE
						}
						// 0x0AB4D2EB 0x0A562C24
						ebp7C[0] = 0
						ebp7C = ebp7C[1:]
					}
					// 0x0A9038FE 0x0A56FA61 0x0A745C77
					ebp54ProcInfo := 0
					// ebp48RetLen := 0
					// ebp50ProcInfoClass := 7 // ProcessDebugPort
					// ebp1CNtQueryInformationProcess(win.GetCurrentProcess(), ebp50ProcInfoClass, &ebp54ProcInfo, 4, &ebp48RetLen)
					// 0x0A9FF9C9
					// win.FreeLibrary(ebp28)
					// 0x0AD31229
					if ebp54ProcInfo != 0 {
						// 0x0A3A8E2B 0x0A6B0AEB
					}
					// 0x0AD8099B 0x0AF86233 0x0A0C3BDF
					label2 = v09FC0B5Alabel2 // [ebp+v0ABF8B07*4+8]=v09FC0B5Alabel2
					ebp38++
					ebp8 &= ^uint32(20) // 清零
				}
				// 0x0A887F5B
				if ebp8&0x40 != 0 {
					// 0x09FC67DB
				}
				// 0x0A05A2B7
				if ebp8&0x80 != 0 {
					// 0x0A889EB6
				}
				// 0x0A56391F，判等
				if ebp8 != 0 {
					// 0x0B07BF42 0x0A6B0AEB
				}
				// 0x0AD83D30
				if ebp38 != 1 {
					// 0x0A9331C5 0x0A6B0AEB
				}
				// 0x0AF77BA8
				v0A325B9F = v09FBB49C
				label3 = v0A557660label3 // [ebp+v0AA0D71B*4+8] = v0A557660label3
				break                    // 0x09FBBB94
			}
			// 0x0A4E4AF6
			ebp8 <<= 1
		} // for loop 0x0A9FBE5C

		// 0x0A6B0AEB 异常进入0x004D7CED姿势不对，会异常退出
		ebp24 = v0A88E351
		label4 = 0x0AD91CED - v09E2318F + ebp24 // 0x004D7CED

		// 0x09FBBB94
		// pop edi
		// pop esi
		// pop ebx
		// leave
		// ret
	}()

	// 0x0AD2B081 0x09FC860D 0x0AB4CCE0 0x0AF7D9C1 0x09DEB0D9 0x0A5FFBD4
	// pop ecx
	// pop edx
	// pop esi
	// pop edi
	// pop ebx
	// popfd
	// 0x0A9FB977 0x0AA0931B
	// label3(0x0B10933D)
	// label2(0x0A8FE9A6)
	// label1(0x0AD56E8A) f0AD56E8A 隐藏函数

	// push ebp
	// ebp = esp
	// eax = 0x1B60
	// 0x004D7CED
	// 1B60h，近两页也就是8k字节的局部变量
	// ebp-4 和 ebp-8是什么？
	// var ebp_4 uint32 = 0x004D7CF2 // 一个函数指针
	var ebp8port uint16 = 10

	// ebpC := win.CreateMutex
	var ebpC win.HANDLE

	// 0x004D7E1E start mu.exe
	if len(szCmdLine) < 1 { // f00DE7C00strlen(szCmdLine)
		// 启动mu.exe然后退出
	}
	// 0x004D7E99
	// vc编译器把先定义的变量地址在高位置？
	var ebp230 struct {
		data [248]uint8
	}
	var ebp238ver [8]uint8
	copy(ebp238ver[:], "unknown")
	f00DE8100memset(ebp230.data[:], 0, 0xF8)

	// var cmd string = GetCommandLine()
	ebp23Ccmd := os.Args[0]

	var ebp244ver version // ebp-244 ~ ebp-23E，这里会有清零操作，应该是程序员自己清的
	var ebp350fileName [260]uint8
	if f004D52AB(ebp350fileName[:], ebp23Ccmd) {
		if f004D55C6(ebp350fileName[:], &ebp244ver) {
			f00DE817Asprintf(ebp238ver[:], "%d.%02d", ebp244ver.major, ebp244ver.minor)
			if ebp244ver.patch > 0 {
				var ebp4DC [2]uint8
				ebp4DC[0] = uint8('a' + ebp244ver.patch - 1) // 0x8C，这里是有问题的，超过ASCII编码范围了
				ebp4DC[1] = 0                                // 我猜测作者是想表达R
				f00DE8010strcat(ebp238ver[:], string(ebp4DC[:]))
			}
		}
	}

	v01319E08log.f00B38AE4printf("\r\n")
	v01319E08log.f00B38D31begin()
	v01319E08log.f00B38D19cut()

	v01319E08log.f00B38AE4printf("Mu online %s (%s) executed. (%d.%d.%d.%d)\r\n",
		ebp238ver[:], v01319DFC,
		ebp244ver.major, ebp244ver.minor, ebp244ver.patch, ebp244ver.fourth)

	v01319E08log.f00B38D49(1)

	// system information

	// direct-x information

	// 0x004D803B
	// 这里重新设置ip地址和端口失败了，为什么？
	// 从命令行中提取指定的ip地址和端口，因为没有通过mu.exe启动，所有没有命令行
	f004D7A1F := func(haystack []uint8, ip []uint8, port *uint16) bool { // ([]uint8(szCmdLine), v01319A50ip[:], &ebp8)
		// x=0x0B37,3D25，存放的是[0,0,0,'K']
		// 0x0A4411D8 0x004D7A3D
		// 10C局部变量
		var ebp100buf [256]uint8
		ebp100buf[0] = 0
		f00DE8100memset(ebp100buf[1:], 0, 0xFF)
		// 0x012E0FF0 0x004D7A4D
		sRet := func(haystack []uint8, needle string) []uint8 { // f004472C9
			return f00DE92E0strstr(haystack, needle)
		}(haystack, "battle")
		// 0x0AF8FB96
		if len(sRet) != 0 {
			// 0x0A9F95D7 0x004D7A5A
		}
		// 0x004D7A66 0x0A601B5B 0x004D7A72 0x0A916BE0
		if f004D755F(haystack, 0x79, ebp100buf[:]) { // 解析第一个参数
			// 0x0A848F36
		}
		// 0x004D7BC3 0x0A7AC798 0x004D7BCF 0x0A38F870
		if !f004D755F(haystack, 0x75, ebp100buf[:]) { // 解析第二个参数
			// 0x0A12F49A
		}
		// 0x004D7BDF 0x0ABB5DC7 0x004D7BE9
		f00DE8000strcpy(ip, ebp100buf[:])
		// 0x0A7AE561 0x004D7BFC 0x09FC6DA2
		if !f004D755F(haystack, 0x70, ebp100buf[:]) { // 解析第三个参数
			// 0x09E91BB9
		}
		// 0x004D7C0C 0x09E90EF9 0x004D7C13 0x0B284D06
		*port = uint16(f00DECBD1atoi(ebp100buf[:]))
		return true
	}
	if f004D7A1F([]uint8(szCmdLine), v01319A50ip[:], &ebp8port) {
		v012E2338ip = string(v01319A50ip[:])
		v012E233Cport = ebp8port
	}
	// 0x004D8067 // open main.exe
	if !f004D5368() {
		f004D9F88()
		return
	}

	// enc and dec init
	v08C8D050enc.f00B62CF0init("Data/Enc1.dat")
	v08C8D098dec.f00B62D30init("Data/Dec2.dat")

	// 0x004D80A8 read config.ini
	v01319E08log.f00B38AE4printf("> To read config.ini.\r\n")
	// 这个函数被花的厉害，难道是核心业务？
	// 追踪只能追踪一个函数
	f004D7281 := func() bool {
		// CURRENT_USER
		// SOFTWARE\\Webzen\\Mu\\Config
		// 348h个字节的局部变量

		// ebp-338
		var ebp_338 struct {
			patch [3]uint8
		}

		// ebp-334 ~ ebp-330
		var ebp_334 version

		// ebp-328 ~ ebp-225
		var fileName [260]uint8

		// ebp-220 ~ ebp-10D
		var buf [276]uint8
		buf[0] = v0114D912                 // 0
		f00DE8100memset(buf[1:], 0, 0x113) // 0x113=275

		// ebp-108 ~ ebp-9
		var bufDir [256]uint8
		win.GetCurrentDirectory(0x100, (*uint16)(unsafe.Pointer(&bufDir[0])))

		f00DE8000strcpy(buf[:], bufDir[:])

		// _00DE7C00
		nameLen := func(bufDir []uint8) uint32 {
			return uint32(len(bufDir))
		}(bufDir[:])
		if bufDir[nameLen-1] == 0x5C {
			f00DE8010strcat(buf[:], "config.ini")
		}
		f00DE8010strcat(buf[:], "\\config.ini") // cat拼凑字符串

		win.GetPrivateProfileStringA("LOGIN", "Version", v0114D913, v01319A44version[:], 8, string(buf[:])) // 1.04.44写到全局变量中
		// ebp_8的底层数组在堆上
		// var ebp_8 string = GetCommandLine()
		var ebp_8 string = os.Args[0]

		// f004D52AB
		bRet := f004D52AB(fileName[:], ebp_8)
		if bRet {
			bRet := f004D55C6(fileName[:], &ebp_334)
			if bRet { // je 0x004D,741B
				f00DE817Asprintf(v01319A38version[:], "%d.%02d", ebp_334.major, ebp_334.minor)

				if ebp_334.patch > 0 { // jle 0x004D,7419
					*(*uint16)(unsafe.Pointer(&ebp_338)) = v0114DD64
					ebp_338.patch[2] = 0

					// 我的猜测是从patch=26开始使用字母和+标记标识patch，比如
					// 27就是A+
					// 28就是B+
					// 44就是R+
					if ebp_334.patch > 0x1A { // jle 0x004D,73EE
						ebp_338.patch[0] = 'A' - 27 + uint8(ebp_334.patch) // 65-27+44=82='R'
						ebp_338.patch[1] = '+'
					}
					f00DE8010strcat(v01319A38version[:], string(ebp_338.patch[:]))
				}
			}
		}

		// f004D7174 这个函数在搞什么？
		func(fileName string) bool {
			// 204个字节的局部变量
			v01319DB8[0] = 0
			var ebpCC uint32 = 0

			var partition struct {
				num  [64]uint8
				name [64]uint8
				ip   [64]uint8
			}
			ebp4 := f00DE909Efopen(fileName, "rt")
			if ebp4 == nil {
				return false
			}
			// 不能超过50条记录
			for ebpCC < 50 {
				// f00DECD20
				nRet := f00DECD20(ebp4, "%s", partition.num[:])
				if nRet == -1 {
					break
				}
				nRet = f00DECD20(ebp4, "%s", partition.name[:])
				if nRet == -1 {
					break
				}
				nRet = f00DECD20(ebp4, "%s", partition.ip[:])
				if nRet == -1 {
					break
				}

				// _00DE94F0
				nRet = func(x []uint8, y []uint8) int32 {
					return -1
				}(v01319A50ip[:], partition.name[:])
				if nRet == 0 {
					// _00DECBD1
					nRet := func(x []uint8) uint32 {
						return 0
					}(partition.ip[:])
					if nRet == uint32(v012E233Cport) {
						f00DE8000strcpy(v01319DB8[:], partition.num[:])
						// _00DE8C84
						func() {

						}()
						return true
					}
				}
				ebpCC++
			}
			f00DE8C84close(ebp4)
			return false
		}("partition.inf")
		// 接下来还有一段处理，我没看
		return true
	}
	if !f004D7281() {
		v01319E08log.f00B38AE4printf("config.ini read error\r\n")
		f004D9F88()
		return
	}

	// hook: v012E4018version = v01319A44version

	// gg init
	// ebp_1A38也可能是个实例指针变量，值是0x0D9F,EA10
	var ebp1A38 *uint32 = (*uint32)(f00DE852Fnew(1))
	var ebp1B0C *t1319D68
	if ebp1A38 != nil { // 0x004D8102 hook always nil，disable GameGurad
		// ebp1B0C = ebp1A38.f004D936A(v012E2214)
	} else {
		ebp1B0C = nil
	}

	v01319D68 = ebp1B0C
	bRet := f00B4C1B8()
	if !bRet { // 0x004D8145 hook always true, disable GameGurad
		// gg init error
		// f004D51C8()
		// f004D9335()
		// f004D9F88()
		return
	}
	v01319E08log.f00B38AE4printf("> gg init success.\r\n")
	v01319E08log.f00B38D19cut()

	if v012E2210 == 1 {
		// ShowCursor(0)
	}
	if v012E2224 == 0 {
		// DisplaySetting
	}

	// window init
	v01319E08log.f00B38AE4printf("> Screen size = %d x %d.\r\n", v012E3F08.height, v012E3F08.width)
	v01319D70hModule = hModule
	v01319D6ChWnd = f004D6F82initWindow(hModule, iCmdShow) // 创建hWnd
	v01319E08log.f00B38AE4printf("> Start window success.\r\n")

	// opengl init
	bRet = f004D6D64() // 设置hWnd资源
	if !bRet {
		// 失败处理
		return
	}
	v01319E08log.f00B38AE4printf("> OpenGL init success.\r\n")
	v01319E08log.f00B38D19cut()

	// opengl information
	v01319E08log.f00B38E3C()
	v01319E08log.f00B38D19cut()

	// sound card information
	v01319E08log.f00B3902D() // 自带cut

	// ShowWindow(v01319D6ChWnd, 10) // 发很多消息给WndProc
	// UpdateWindow(v01319D6ChWnd) // 发很多消息给WndProc

	f00B4C1FF(v01319D6ChWnd)
	v01319E08log.f00B38AE4printf("> gg connect window Handle.\r\n")
	v01319E08log.f00B38D19cut()

	// ime information
	v01319E08log.f00B38EF4(v01319D6ChWnd)
	v01319E08log.f00B38D19cut()

	//
	f00DEE871setlocale(0, v01319E00) // 参数0是哪个LC？
	// f0043BF3FgetT4003(v01319D6ChWnd, v012E3F08.height, v012E3F08.width)
	// f0043BF9C()
	// f0090E94C()
	// f0090B256()
	// f00A49798(v01319D6ChWnd, v012E3F08.height, v012E3F08.width)
	// f00A49E40() // r6602 floating point support not loaded

	// ...

	// ...

	var ebp595 uint8 // ImmIsUIMessage返回值

	// ebp-28
	var msg win.MSG
	// 消息循环
	for {
		if msg.Message == win.WM_CLOSE ||
			msg.Message == win.WM_QUIT {
			break
		}
		if win.PeekMessage(&msg, 0, 0, 0, 0) {
			// win.ImmIsUIMessage()
			if msg.Message == 0x100 || // WM_KEYFIRST, WM_KEYDOWN
				msg.Message == 0x101 || // WM_KEYUP
				ebp595 != 0 ||
				msg.Message == 0x201 || // WM_LBUTTONDOWN
				msg.Message == 0x202 { // WM_LBUTTONUP
				// f00A49798(msg.HWnd, msg.Message, msg.WParam, msg.LParam, 1)
				// f00A4DFB7()
			}

			if 0 == win.GetMessage(&msg, 0, 0, 0) { // WM_QUIT返回0，出错返回-1
				break
			}
			win.TranslateMessage(&msg)
			win.DispatchMessage(&msg) // 到WndProc
		} else {
			if v012E2224 == 0 {
				v01319D9C++
				if v01319D9C > 30 {
					v01319D98++
					if v01319D98 > 30 {
						v01319D9C = 0
						v01319D98 = 0
						if v01319DB4 == 0 {
							var packet pb
							packet.f00439178init()
							packet.f0043922CwritePrefix(0xC1, 0xF1)
							packet.f004397B1writeUint8(3)
							packet.f004397B1writeUint8(1)
							packet.f004397B1writeUint8(0)
							packet.f004393EAsend(true, false)
							packet.f004391CF()
							v01319DB4 = 1
						}
						win.SetTimer(v01319D6ChWnd, 1001, 1000, 0)
						// v086A3B94 = f007452FD
					}

					// 窗口非最小化
					if !win.IsIconic(v01319D6ChWnd) {
						v01319D98 = 0
						v01319D9C = 0
					}
				}
			}
			if v012E2224 == 1 || v01319D95 != 0 {
				f004E6233handleState(v01319D74hDC) // 状态机 重要
			} else if v012E2224 == 0 {
				win.SetForegroundWindow(v01319D6ChWnd)
				win.SetFocus(v01319D6ChWnd)
				if v01319DA0 > 1 {
					win.SetTimer(v01319D6ChWnd, 1001, 1000, 0) // 定时器id=1001，1秒触发
					win.PostMessage(v01319D6ChWnd, 10, 0, 0)
				} else {
					v01319DA0++
					v01319D90 = 1
					win.ShowWindow(v01319D6ChWnd, win.SW_MINIMIZE)
					v01319D90 = 0
					win.ShowWindow(v01319D6ChWnd, win.SW_MAXIMIZE)
				}
			}
		}
		f00760292parse(&v08C88FF0conn, 0, 0) // 从连接中获取报文并解析
		// v01319D2C.f004B9492()
	}

	// 退出消息循环
	// f00535963()
	if v01319D68 != nil {
		//v01319D68.f004D9335(1) // 释放
		v01319D68 = nil
	}
	win.CloseHandle(ebpC)
	//f004D56B2()
	// ebp1.f004D9F88()
	//return ebp20
}
func f00A49798(hWnd win.HWND, message uint32, wParam, lParam uintptr, x int) {

}

func f00A49E40() {

	// ...

	// _00A4F3F2
	func() {
		// _00AECB6C
		func() {
			// _00AECB6C
			func() {
				// _00BB0020
				func() {
					// _00BF41C0
					func() {
						// _00BF36A0 好复杂
						func() {
							// _00D05350
							func() {
								// _00CEB740
								func() {
									// _00CBFAD0
									func() {
										// _00BCC780
										func() {
											// _00DF02FC
											func() {
												// _00DF01F5
												func() {
													// _00DF0031
													func() {
														// _00E08316 对xmin = %g, ymin = %g, xmax = %g...
														func() {
															// 00E1A4AC
															func() {
																// _00DF0407
																func() {
																	// _00DFEA1B
																	func() {
																		// Runtime Error!\n
																		// program name unknown
																		// win.GetModuleFileName(0, v09D9DBB1, 260)
																		// _00E1AF3C
																		func() {
																			// MessageBox弹出错误
																		}()
																	}()
																}()
															}()
														}()
													}()
												}()
											}()
										}()
									}()
								}()
							}()
						}()
					}()
				}()
			}()
		}()
	}()
}

var v08C88F60 uint8
var v08C88F61 uint8
var v08C88F62 int
var v08C88F64 int
var v08C88F74 int
var v08C88F78username [11]uint8 // username
var v08C88F84 uint32
var v08C88F88 uint32

var v086105ECcharacter = &struct {
	name [100]uint8
	m10C uint16
	m154 uint16 // 0
	m160 uint16 // 0
	m178 [14]uint16
}{}

func f004D8FF5() {
	// SEH
	f00DE8A70chkstk() // 0x1488
	if v08C88F64 == 0 {
		return
	}

	// 构造心跳报文
	var alive pb
	alive.f00439178init()
	alive.f0043922CwritePrefix(0xC1, 0x0E) // 写前缀
	ebp10 := win.GetTickCount()
	alive.f0043974FwriteZero(1)                         // 写0
	alive.f0043EDF5writeUint32(ebp10)                   // 写time
	alive.f004C65EFwriteUint16(v086105ECcharacter.m154) // 什么
	alive.f004C65EFwriteUint16(v086105ECcharacter.m160) // 什么
	alive.f004393EAsend(true, false)                    // 发送心跳报文
	if v08C88F62 == 0 {
		v08C88F62 = 1
		v08C88F84 = ebp10
	}
	alive.f004391CF()
}

func f004D9F88() {}

// 也可能是map
var v0075FF6Acmd []struct{}
var v0075FCB2cmd []struct{}

var v086A3C28 int
