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
	v0114DB18 string = "서버와의 접속종료." // 终止与服务器的连接
	v0114DD64 uint16 = 0x0061
)

// 这里使用数组，确保在data段，固定在某个地址上
// 如果使用string，那么读文件后对buf进行string类型转换，会memmove到堆上
var v012E10F4 [3]uint8 // FC CF AB
var v012E2210 uint32 = 1
var v012E2220 uint32 = 900

type displayMode int

const (
	displayModeFullScreen displayMode = 0
	displayModeWindow     displayMode = 1
)

var v012E2224displayMode displayMode

var v012E2228displayColorBit int = 16
var v012E222C *os.File
var v012E2338ip = "192.168.0.102"
var v012E233Cport uint16 = 44405

// 1: loading
// 2: server select
// 4: character select
// 5: game
var v012E2340state = 1

var v012E3F08resolution struct {
	width  int // 0x320, 800
	height int // 0x258, 600
}
var v012E4018versionDLL = [8]uint8{'2', '2', '7', '8', '9', 0, 0, 0}                                            // "22789"
var v012E4020serial = [16]uint8{'M', '7', 'B', '4', 'V', 'M', '4', 'C', '5', 'i', '8', 'B', 'C', '4', '9', 'b'} // "M7B4VM4C5i8BC49b"
var v012F7910 os.File

var v01319A38version [8]uint8     // "1.04R+"	// 可执行程序版本号
var v01319A44versionFile [8]uint8 // "1.04.44" // 配置文件版本号
var v01319A50ip [16]uint8

var v01319D18 []uint8
var v01319D1CwndProc uintptr
var v01319D54 *t0114EEF8
var v01319D65 uint32
var v01319D68gg *t1319D68
var v01319D6ChWnd win.HWND
var v01319D70hInstance win.HINSTANCE
var v01319D74hDC win.HDC     // hDC, 0x1401,11A6
var v01319D78hGLRC win.HGLRC // hGLRC, OpenGL rendering context, 0x0001,0000
var v01319D90 int
var v01319D95 bool
var v01319D98 int
var v01319D9C int
var v01319DA0 int
var v01319DB4 int
var v01319DB8name [64]uint8
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

func f004D52ABtruncateCmd(cmdshort []string, cmd []string) bool {
	// 从命令行字符串中提取出应用程序名，也就是main.exe，存到fileName数组中
	return true
}

func f004D5368() bool {
	var ebp110cmd []string // ebp-110
	cmd := os.Args         // ebp-4
	f004D52ABtruncateCmd(ebp110cmd, cmd)

	// dwDesiredAccess = 1<<31, GENERIC_READ
	// dwShareMode = 3, FILE_SHARE_READ | FILE_SHARE_WRITE
	// psa = 0, default security and cannot be inherited by any child
	// dwCreationDisposition = 3, OPEN_EXISTING
	// dwFlagsAndAttributes = 0x80, FILE_ATTRIBUTE_NORMAL
	// dwFileTemplate = 0, When opening an existing file, CreateFile ignores this parameter
	// hFile := CreateFile(fileName[:], 1<<31, 3, 0, 3, 0x80, 0)
	file, err := os.Open(ebp110cmd[0])
	if err != nil {
		return false
	}
	v012E222C = file
	return true
}

func f004D55C6queryVersion(fileName string, ver *version) bool {
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

func f004D5F98mainWndProcOrigin(hWnd win.HWND, message uint32, wParam, lParam uintptr) uintptr {
	ebp29B4message := message
	switch ebp29B4message {
	case 0x1: //win.WM_CREATE
	case 0x6: //win.WM_ACTIVATE, wparam:1, lparam:0
		// 0x004D6620
		active := uint(wParam) & 0xFFFF // WA_ACTIVE:1, WA_CLICKACTIVE:2, WA_INACTIVE:0
		if active == 0 {                // WA_INACTIVE
			// 0x004D6633
			if v012E2224displayMode == displayModeFullScreen {
				v01319D95 = false
			}
			if v012E2224displayMode == displayModeWindow {
				// v08C88BC0 = 0
				// v08C88BC1 = 0
				// v08C88BC2 = 0
				// v08C88BC3 = 0
				// v08C88BC4 = 0
				// v08C88BC5 = 0
				// v08C88BC6 = 0
				// v08C88A99 = 0
				// v08C88A9A = 0
				// v08C88A9B = 0
				// v08C88BBC = 0
			}
			// 0x004D6692 0x004D669F 0x004D6BD7
		} else {
			// 0x004D6694 0x004D6698
			// 防窗口非激活作弊
			v01319D95 = true
			// 0x004D6BD7
		}
	case 0x113: // win.WM_TIMER
		switch wParam {
		case 1000: // 0x3E8
			// f00B4C239
			func() {
				if v01319D68gg == nil {
					return
				}
				v01319D68gg.f00B4CC0D()
				// ...
			}()
			// f004D8FF5() // 发送心跳0E
			func() {
				f00DE8A70chkstk() // 0x1488
				if v08C88F64 == 0 {
					return
				}
				// 构造心跳报文
				var alive pb
				alive.f00439178init()
				alive.f0043922CwriteHead(0xC1, 0x0E) // 写前缀
				ebp10 := win.GetTickCount()
				alive.f0043974FwriteZero(1)                     // 写0
				alive.f0043EDF5writeUint32(ebp10)               // 写time
				alive.f004C65EFwriteUint16(v086105ECpanel.m154) // 什么
				alive.f004C65EFwriteUint16(v086105ECpanel.m160) // 什么
				alive.f004393EAsend(true, false)                // 发送心跳报文
				if v08C88F62 == 0 {
					v08C88F62 = 1
					v08C88F84 = ebp10
				}
				alive.f004391CF()
			}()
			// 0x004D6823 0x004D6BD7
		}
	case 0x400: // win.WM_USER
		switch lParam {
		case 1:
			v08C88FF0conn.f006BDA03read()
		case 2:
			v08C88FF0conn.f006BD945write()
		case 20:
		}
	default: // win.WM_KILLFOCUS/8
	}
	// 0x004D6BD7:
	// v08C88BC0 = 0
	// if v08C88BC5 == 1 && (v01319DA8 != v08C88C80 || v01319DAC != v08C88C80) {
	// 	v08C88BC5 = 0
	// }
	return win.DefWindowProc(hWnd, message, wParam, lParam) // 原始是NtdllDefWindowProc，这是什么？
}

func f004D6C2BmainWndProc(hWnd win.HWND, message uint32, wParam, lParam uintptr) uintptr {
	switch message {
	case 0x102: // win.WM_CHAR
		// f00A49798ui().f00A4DF93handleKeyPress(wParam, lParam)
	}
	return win.CallWindowProc(v01319D1CwndProc, hWnd, message, wParam, lParam) // 会经过 dll.user32.xxx 再回调到 f004D5F98mainWndProcOrigin
}

func f004D6F82initWindow(hInstance win.HINSTANCE, iCmdShow int) win.HWND {

	// 这里windows那边使用的是数组，编译器会使用movsw和movsb来给字符数组赋值
	// var ebp_78 string = "MU"
	ebp78, _ := syscall.UTF16FromString("MU")

	// ...

	// 140字节局部变量
	ebp34wndClass := win.WNDCLASSEX{
		CbSize:        0x30, // ebp-34
		Style:         0x2B, // CS_OWNDC | CS_DBLCLKS | CS_HREDRAW | CS_VREDRAW
		LpfnWndProc:   syscall.NewCallback(f004D5F98mainWndProcOrigin),
		CbClsExtra:    0,
		CbWndExtra:    0,
		HInstance:     win.GetModuleHandle(nil),
		HIcon:         win.LoadIcon(0, nil),              // IDI_APPLICATION)
		HCursor:       win.LoadCursor(0, nil),            // IDC_ARROW
		HbrBackground: win.HBRUSH(win.GetStockObject(4)), // WHITE_BRUSH
		LpszMenuName:  nil,
		LpszClassName: &ebp78[0],  // "MU"
		HIconSm:       0x141F0439, // ebp8
	}
	win.RegisterClassEx(&ebp34wndClass)

	// ...

	// CreateWindow会发大概3条消息给WndProc且create IME window
	ebp4hWnd := win.CreateWindowEx(
		0,          // ExStyle
		&ebp78[0],  // ClassName, "MU"
		&ebp78[0],  // WindowName, "MU"
		0x02CA0000, // Style, WS_CLIPCHILDREN | WS_CAPTION | WS_MINIMIZEBOX | WS_SYSMENU
		0x118,      // x
		0x46,       // y
		0x326,      // width
		0x274,      // height
		0,          // parent window
		0,          // menu
		hInstance,  // Instance
		nil,
	)
	v01319D1CwndProc = win.SetWindowLongPtr(ebp4hWnd, win.GWL_WNDPROC, syscall.NewCallback(f004D6C2BmainWndProc))
	return ebp4hWnd
}

func f004D6D64glInit() bool {
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

	win.ShowWindow(v01319D6ChWnd, 5 /*SW_SHOW*/) // also create "MSCTFIME UI" window
	// SetForegroundWindow: Brings the thread that created the specified window
	// into the foreground and activates the window.
	// Keyboard input is directed to the window, and various visual cues are changed for the user.
	// The system assigns a slightly higher priority to the thread that
	// created the foreground window than it does to other threads.
	win.SetForegroundWindow(v01319D6ChWnd)
	win.SetFocus(v01319D6ChWnd)
	return true
}

func f004D755F(haystack string, y int, buf []uint8) bool {
	return false
}

// f004D7CE5winMain, WinMain
func f004D7CE5winMain(hInstance win.HINSTANCE, hPrevInstance win.HINSTANCE, szCmdLine string, iCmdShow int) int {
	func() {
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
						// int3 ; 频繁，报0x80000003, EXCEPTION_BREAKPOINT
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
						// int1 ;频发，单步直接stepover了，没有进入SEH异常处理，而Run是ok的报0x80000004, EXCEPTION_SINGLE_STEP
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
	}()
	// 0x0AD56E8A: winMain logic
	// 0x1B60局部变量 近两页也就是8k字节的局部变量
	// 0x004D7CED:
	f00DE8A70chkstk()
	// f00DE8100memset(ebp28msg[:], 0, 0x1C)
	// ebp1.f004D9EBB(f004D9460get())
	func() {
		func() { // f004D9ED4()
			func() { // f004D9EEB()
				// var ebp20 struct {
				// 	m00 int // 0
				// 	m04 int // 0x10
				// 	m08 int // 0x4000
				// 	m0C int // 0x4000
				// 	m10 int // 0x40000
				// 	m14 int // 0
				// 	m18 int // 1
				// 	m1C int // 0
				// }
				// f00BABDE0(ebp20.f004D9F06())
				func() {
					// if v09D9BD6C != nil {
					// 	return
					// }
					// v01319F18.do2(&ebp20)
					// v09D9BD6C = &v01319F18
				}()
			}()
		}()
	}()
	// f00B43232(f004D7C24, 0)
	// 0x004D7D2D: check named mutex
	var ebpC win.HANDLE
	if f0043A2DFgetBattleCoreSpec().f0043A33Cenable() {
		// ebpC = dll.kernel32.CreateMutex(0, 1, "MuBattleCore")
		// if dll.kernel32.GetLastError() == 183 { // ERROR_ALREADY_EXISTS
		// 	dll.kernel32.CloseHandle(ebpC)
		// 	ebp1A1C := 0
		// 	ebp1.f004D9F88()
		// 	return ebp1A1C
		// }
	} else if v012E2224displayMode == displayModeWindow {
		// ebpC = dll.kernel32.CreateMutex(0, 1, "MuOnline")
		// if dll.kernel32.GetLastError() == 183 {
		// 	dll.kernel32.CloseHandle(ebpC)
		// 	ebp1A20 := 0
		// 	ebp1.f004D9F88()
		// 	return ebp1A20
		// }
	}
	// 0x004D7DD3: check whether auto update window is working
	var ebp2C win.HWND // ebp2C := dll.user32.FindWindow("#32770", "MU Auto Update")
	if ebp2C != 0 {
		win.SendMessage(ebp2C, 0x10, 0, 0)
	}
	// if f004D77E0() == false {
	// 	ebp1A24 := 0
	// 	ebp1.f004D9F88()
	// 	return ebp1A24
	// }
	// 0x004D7E1E: check command parameter number, or start mu.exe
	if len(szCmdLine) < 1 { // f00DE7C00strlen(szCmdLine)
		// launch mu.exe and exit
	}
	// 0x004D7E99: output version to log
	// var ebp244 struct {
	// 	versem version    // ebp244
	// 	cmd    string     // ebp23C
	// 	ver    [8]uint8   // ebp238
	// 	data   [248]uint8 // ebp230
	// }
	var ebp230data [248]uint8
	var ebp238ver [8]uint8
	copy(ebp238ver[:], "unknown")
	f00DE8100memset(ebp230data[:], 0, 0xF8)
	// 命令行和环境变量参数已经在进程地址空间了，kernel32.GetCommandLine系统调用只是返回地址
	ebp23Ccmd := os.Args
	var ebp244ver version // ebp-244 ~ ebp-23E，这里会有清零操作，应该是程序员自己清的
	var ebp350cmd []string
	if f004D52ABtruncateCmd(ebp350cmd, ebp23Ccmd) && f004D55C6queryVersion(ebp350cmd[0], &ebp244ver) {
		f00DE817Asprintf(ebp238ver[:], "%d.%02d", ebp244ver.major, ebp244ver.minor)
		if ebp244ver.patch > 0 {
			var ebp4DC [2]uint8
			ebp4DC[0] = uint8('a' + ebp244ver.patch - 1) // 0x8C，这里是有问题的，超过ASCII编码范围了 需要MemAssign(0x004D7F64+3, byte(0xC5))
			ebp4DC[1] = 0                                // 我猜测作者是想表达R
			f00DE8010strcat(ebp238ver[:], string(ebp4DC[:]))
		}
	}
	v01319E08log.f00B38AE4printf("\r\n")
	v01319E08log.f00B38D31begin()
	v01319E08log.f00B38D19cut()
	v01319E08log.f00B38AE4printf("Mu online %s (%s) executed. (%d.%d.%d.%d)\r\n",
		ebp238ver[:], // version
		v01319DFC,    // local
		ebp244ver.major, ebp244ver.minor, ebp244ver.patch, ebp244ver.fourth)
	v01319E08log.f00B38D49(1)

	// system information

	// direct-x information

	// 0x004D803B: 从命令行中提取指定的ip地址和端口，因为没有通过mu.exe启动，所有没有命令行
	f004D7A1F := func(haystack string, ip []uint8, port *uint16) bool { // ([]uint8(szCmdLine), v01319A50ip[:], &ebp8)
		// x=0x0B37,3D25，存放的是[0,0,0,'K']
		// 0x0A4411D8 0x004D7A3D
		// 10C局部变量
		var ebp100buf [256]uint8
		ebp100buf[0] = 0
		f00DE8100memset(ebp100buf[1:], 0, 0xFF)
		// 0x012E0FF0 0x004D7A4D
		sRet := func(haystack string, needle string) []uint8 { // f004472C9
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
	var ebp8port uint16 = 10
	if f004D7A1F(szCmdLine, v01319A50ip[:], &ebp8port) {
		v012E2338ip = string(v01319A50ip[:])
		v012E233Cport = ebp8port
	}

	// 0x004D8067: open main.exe
	if !f004D5368() {
		f004D9F88()
		ebp1A2C := 0
		return ebp1A2C
	}

	// 0x004D808A: enc and dec init
	v08C8D050enc.f00B62CF0init("Data/Enc1.dat")
	v08C8D098dec.f00B62D30init("Data/Dec2.dat")

	// 0x004D80A8: read config.ini
	v01319E08log.f00B38AE4printf("> To read config.ini.\r\n")
	// f004D7281函数被花了
	f004D7281 := func() bool {
		// 348h个字节的局部变量

		// config file version
		var buf [276]uint8
		buf[0] = v0114D912                 // 0
		f00DE8100memset(buf[1:], 0, 0x113) // 0x113=275
		var bufDir [256]uint8
		win.GetCurrentDirectory(0x100, (*uint16)(unsafe.Pointer(&bufDir[0])))
		f00DE8000strcpy(buf[:], bufDir[:])
		nameLen := f00DE7C00strlen(bufDir[:])
		if bufDir[nameLen-1] == 0x5C {
			f00DE8010strcat(buf[:], "config.ini")
		}
		f00DE8010strcat(buf[:], "\\config.ini")                                                          // cat拼凑字符串
		win.GetPrivateProfileStringA("LOGIN", "Version", "", v01319A44versionFile[:], 8, string(buf[:])) // 1.04.44写到全局变量中

		// 0x004D7338: exe file version
		var ebp338 struct {
			patch [3]uint8
		}
		var ebp334 version
		var ebp328cmd []string
		ebp8cmd := os.Args // GetCommandLine()
		if f004D52ABtruncateCmd(ebp328cmd, ebp8cmd) == false {
			// 0x004D742E
		}
		if f004D55C6queryVersion(ebp328cmd[0], &ebp334) == false {
			// 0x004D741B
		}
		f00DE817Asprintf(v01319A38version[:], "%d.%02d", ebp334.major, ebp334.minor)
		if ebp334.patch > 0 { // jle 0x004D7419
			*(*uint16)(unsafe.Pointer(&ebp338)) = v0114DD64
			ebp338.patch[2] = 0

			// 我的猜测是从patch=26开始使用字母和+标记标识patch，比如
			// 27就是A+
			// 28就是B+
			// 44就是R+
			if ebp334.patch > 0x1A { // jle 0x004D,73EE
				ebp338.patch[0] = 'A' - 27 + uint8(ebp334.patch) // 65-27+44=82='R'
				ebp338.patch[1] = '+'
			}
			f00DE8010strcat(v01319A38version[:], string(ebp338.patch[:]))
		}

		// 0x004D743F: partition.inf
		func(fileName string) bool { // f004D7174("partition.inf")
			// 204个字节的局部变量
			v01319DB8name[0] = 0
			ebp4file := f00DE909Efopen(fileName, "rt")
			if ebp4file == nil {
				return false
			}
			// 不能超过50条记录
			ebpCC := 0
			var partition struct {
				name [64]uint8
				ip   [64]uint8
				port [64]uint8
			}
			for ebpCC < 50 {
				if f00DECD20fscanf(ebp4file, "%s", partition.name[:]) == -1 {
					break
				}
				if f00DECD20fscanf(ebp4file, "%s", partition.ip[:]) == -1 {
					break
				}
				if f00DECD20fscanf(ebp4file, "%s", partition.port[:]) == -1 {
					break
				}
				if f00DE94F0strcmp(v01319A50ip[:], partition.ip[:]) == 0 && f00DECBD1atoi(partition.port[:]) == int(v012E233Cport) {
					f00DE8000strcpy(v01319DB8name[:], partition.name[:])
					f00DE8C84close(ebp4file)
					return true
				}
				ebpCC++
			}
			f00DE8C84close(ebp4file)
			return false
		}("partition.inf")

		// 0x004D744A: EnumDisplaySettings
		f00B0EF1EmodeManager().f00B0EF7BenumMode()
		// f00B0E9BFreg().f00B0EA1Cconstruct("SOFTWARE\\Webzen\\Mu\\Config")
		// f00B0E9BFreg().f00B0EBC4get()

		// 0x004D7473: FullScreenMode
		// if f00B0E9BFreg().f00B0EE48getFullScreen() == 0 {
		// 	v012E2224displayMode = displayModeWindow
		// } else {
		// 	v012E2224displayMode = displayModeFullScreen
		// }

		// 0x004D74A1: DisplayDeviceModelIndex/resolution
		// m := f00B0EF1EmodeManager().f00B0F1F6getMode(f00B0E9BFreg().f00B0EE36getModeIndex())
		// if m == nil {
		// 	m = f00B0EF1EmodeManager().f00B0F1F6getMode(0)
		// }
		// v012E3F08resolution.width = m.dmPelsWidth
		// v012E3F08resolution.height = m.dmPelsHeight
		// v01308EC4resolution.width = v012E3F08resolution.width / v0114AE18 // widht/640.0
		// v01308EC4resolution.height = v012E3F08resolution.height / v0114AE10 // height/480.0

		// 0x004D7532: DisplayColorBit
		// if f00B0E9BFreg().f00B0EE5AgetColorBit() == 1 {
		// 	v012E2228displayColorBit = 32
		// } else {
		// 	v012E2228displayColorBit = 16
		// }
		return true
	}
	if f004D7281() == true {
		v01319E08log.f00B38AE4printf("config.ini read error\r\n")
		ebp1A30 := 0
		f004D9F88()
		return ebp1A30 // 0x004D8FE2
	}

	// 0x004D80ED: gg init
	var ebp1A38 *uint32 = (*uint32)(f00DE852Fnew(1)) // hook: v012E4018versionDLL = v01319A44versionFile
	var ebp1B0C *t1319D68
	if ebp1A38 != nil { // 0x004D8102 gameguard 1, hook always nil，disable GameGurad
		// ebp1B0C = ebp1A38.f004D936A("MUCN")
	} else {
		ebp1B0C = nil
	}
	v01319D68gg = ebp1B0C
	if false { // if f00B4C1B8() == true { // 0x004D8145 gameguard 2, hook always true, disable GameGurad
		// gg init error
		// f004D51C8()
		// f004D9335()
		// f004D9F88()
		return 0
	}
	v01319E08log.f00B38AE4printf("> gg init success.\r\n")
	v01319E08log.f00B38D19cut()

	if v012E2210 == 1 {
		// user32.ShowCursor(0)
	}
	if v012E2224displayMode == displayModeFullScreen {
		// user32.DisplaySetting
	}

	// 0x004D828B: window init
	v01319E08log.f00B38AE4printf("> Screen size = %d x %d.\r\n", v012E3F08resolution.height, v012E3F08resolution.width)
	v01319D70hInstance = hInstance
	v01319D6ChWnd = f004D6F82initWindow(hInstance, iCmdShow) // 创建hWnd
	v01319E08log.f00B38AE4printf("> Start window success.\r\n")

	// 0x004D82D4: opengl init
	if f004D6D64glInit() == false {
		// ...
		return 0
	}
	v01319E08log.f00B38AE4printf("> OpenGL init success.\r\n")
	v01319E08log.f00B38D19cut()
	v01319E08log.f00B38E3C() // opengl information
	v01319E08log.f00B38D19cut()

	// 0x004D8375: sound card information
	v01319E08log.f00B3902D()                             // 自带cut
	win.ShowWindow(v01319D6ChWnd, 10 /*SW_SHOWDEFAULT*/) // SW_SHOWDEFAULT 发很多消息给WndProc
	win.UpdateWindow(v01319D6ChWnd)                      // 发很多消息给WndProc

	// 0x004D839A: gg connect
	func(hWnd win.HWND) { // f00B4C1FF
		if v01319D68gg != nil {
			// ...
		}
	}(v01319D6ChWnd)
	v01319E08log.f00B38AE4printf("> gg connect window Handle.\r\n")
	v01319E08log.f00B38D19cut()

	// 0x004D83C1: ime information
	v01319E08log.f00B38EF4(v01319D6ChWnd)
	v01319E08log.f00B38D19cut()

	// 0x004D83E3:
	f00DEE871setlocale(0 /*LC_ALL*/, v01319E00)
	f0043BF3FgetT4003().f0043BF9C(v01319D6ChWnd, v012E3F08resolution.width, v012E3F08resolution.height) // 800*600
	// f0090E94C().f0090B256()
	f00A49798ui().f00A49E40init1(v01319D6ChWnd, v012E3F08resolution.height, v012E3F08resolution.width) // r6602 floating point support not loaded
	/*
		if f00B0E9BF().f00B0EE6C() {
			v01319D67 = true
		} else {
			v01319D67 = false
		}
		if dll.wzAudio.wzAudioCreate(v01319D6ChWnd) {
			dll.user32.MessageBox(v01319D6ChWnd, "No audio devices in this system", 0, 0)
			v01319E08log.f00B38AE4printf("> wzAudio Create Fail.\r\n")
			v01319E08log.f00B38D19cut()
		} else {
			dll.wzAudio.wzAudioOption(0, 1)
		}
		if f00B0E9BF().f00B0EE7E() && f007DA2B1(v01319D6ChWnd) < 0 {
			dll.user32.MessageBox(v01319D6ChWnd, "Direct Sound Init Fail")
			v01319E08log.f00B38AE4printf("> DirectSound Init Fail.\r\n")
			v01319E08log.f00B38D19cut()
			f007DA4B7(0)
		}
	*/
	// 0x004D84F9:
	// user32.SetTimer(v01319D6ChWnd, 1000, 20*1000, nil)

	// 0x004D8511: v01319B88 init
	f00DE8A9Bsrand(int64(f00DEE9E1time(nil))) // f004D8FE7time
	for ebp58C := 0; ebp58C < 100; ebp58C++ {
		v01319B88[ebp58C] = f00DE8AADrand() % 360
	}
	// 0x004D855B: array size <= 100
	ebp1A54 := f00DE64BCnew(uint(f00DE8AADrand()%100) + 1)
	v01319D18 = ebp1A54

	// 0x004D857F: array size = 7k
	ebp1A58 := f00DE64BCnew(7 * 1024)
	v086105E0 = ebp1A58

	// 0x004D859B:
	ebp1A5C := make([]t086105E4, 650)
	v086105E4 = ebp1A5C

	// 0x004D85B7: object pool init
	var ebp1B18 *t01173630objectPool
	ebp1A64 := &t01173630objectPool{}
	if ebp1A64 != nil {
		ebp1B18 = ebp1A64.f00A38B97construct()
	}
	ebp1A60 := ebp1B18
	v01319D44objectPool = ebp1A60

	// 0x004D85FF:
	v086105E8player = &player{}
	v086105E8player.f005A31C1construct()
	// 0x004D8689:
	v086105ECpanel = &v086105E8player.m04panel
	v086105E8player.f005A3337init()

	// 0x004D86A1:
	if v012E2210 == 1 {
		// 0x004D86AE: new v01319D20
		var ebp1B20 *t0114E1ACwindow
		ebp1A74 := &t0114E1ACwindow{} // f00DE852Fnew(0x300)
		if ebp1A74 != nil {
			ebp1B20 = ebp1A74.f004D9AFCconstruct()
		}
		ebp1A70 := ebp1B20
		v01319D20 = ebp1A70

		// 0x004D86F9: new v01319D24
		var ebp1B24 *t0114AE34edit
		ebp1A7C := &t0114AE34edit{}
		if ebp1A7C != nil {
			ebp1B24 = ebp1A7C.f004521F5construct()
		}
		ebp1A78 := ebp1B24
		v01319D24 = ebp1A78

		// 0x004D8744: v01319D28
		var ebp1B28 *t0114AE34edit
		ebp1A84 := &t0114AE34edit{}
		if ebp1A84 != nil {
			ebp1B28 = ebp1A84.f004521F5construct()
		}
		ebp1A80 := ebp1B28
		v01319D28 = ebp1A80
	}
	// 0x004D878F: v01319D2C
	// 0x004D87DA: v01319D30
	// 0x004D8822: v01319D34
	// 0x004D886A: v01319D38
	// 0x004D88B5: v01319D3C
	// 0x004D8906: v01319D40
	// 0x004D894E: v01319D48
	// 0x004D8996: v01319D4C
	// 0x004D89DE: v01319D50
	// 0x004D8A26:
	var ebp1B50 *t0114EEF8
	ebp1AD4 := &t0114EEF8{}
	if ebp1AD4 != nil {
		ebp1B50 = ebp1AD4.f004EB62Bconstruct()
	}
	ebp1AD0 := ebp1B50
	v01319D54 = ebp1AD0
	// 0x004D8A6E: v01319D8C
	// 0x004D8AB6: v01319D58
	// 0x004D8AFE: v01319D5C
	// 0x004D8B46:
	// v013199F4.f004D9395(f00AFEE9F(&ebp1AF4))
	// ebp1AF4.f004CE748()
	// v09D8D0AC.f004D9503(f0093AF95(&ebp1AFC))
	// ebp1AFC.f004D9381()
	f004A7D34getWindowManager().f004A89FFconstruct()
	// 0x004D8B98:
	if v012E2210 == 1 {
		// 0x004D8BA5: win10非法内存访问
		v01319D20.do2(v01319D6ChWnd)
		v01319D24.do12initWindow(v01319D6ChWnd, 200, 20, 50, false)
		v01319D28.do12initWindow(v01319D6ChWnd, 140, 20, 9, true)
		v01319D24.do2showWindow(3) // hide
		v01319D28.do2showWindow(3) // hide
		/*
			// 0x004D8C21: imm
			v01308EF8 = 0
			ebp590ctx := imm32.ImmGetContext(v01319D6ChWnd)
			imm32.ImmSetConversionStatus(ebp590ctx, 0, 0)
			imm32.ImmReleaseContext(v01319D6ChWnd, ebp590ctx)
			// f00451FB9()
			func() {
				ebp8ctx := imm32.ImmGetContext(v01319D6ChWnd)
				imm32.ImmGetConversionStatus(ebp8ctx, &v01308EE8, &v01308EEC)
				imm32.ImmSetConversionStatus(ebp8ctx, 0, 0)
				imm32.ImmReleaseContext(v01319D6ChWnd, ebp8ctx)
			}()
			v01308EF8 = 1
		*/
	}
	// 0x004D8C68: SystemparametersInfo
	if v012E2224displayMode == displayModeFullScreen {
		/*
			user32.SystemParametersInfo(0x61, 1, &ebp594, 0)
			user32.SystemParametersInfo(0x0E, 0, &v012E2220, 0)
			user32.SystemParametersInfo(0x0F, 0x4650, nil, 0)
		*/
	}
	// 0x004D8CA6:
	// f0053594C(v01319D70hInstance, v01319D6ChWnd)
	// 0x004D8CB9: 消息循环
	var ebp28msg win.MSG
	for {
		if ebp28msg.Message == win.WM_CLOSE || ebp28msg.Message == win.WM_QUIT {
			break // 0x004D8F6C
		}
		if win.PeekMessage(&ebp28msg, 0, 0, 0, 0) {
			// 0x004D8CE7
			// Checks for messages intended for the IME window and sends those messages to the window.
			// Returns a nonzero value if the message is processed by the IME window, or 0 otherwise.
			// ebp595 := bool(dll.imm32.ImmIsUIMessage(0, ebp28msg.Message, ebp28msg.WParam, ebp28msg.LParam)) // 输入法
			if ebp28msg.Message == 0x100 || // WM_KEYFIRST, WM_KEYDOWN
				ebp28msg.Message == 0x101 || // WM_KEYUP
				// ebp595 == true ||
				ebp28msg.Message == 0x201 || // WM_LBUTTONDOWN
				ebp28msg.Message == 0x202 { // WM_LBUTTONUP
				f00A49798ui().f00A4DFB7handleClick(ebp28msg.HWnd, ebp28msg.Message, ebp28msg.WParam, ebp28msg.LParam, true)
			}
			if 0 == win.GetMessage(&ebp28msg, 0, 0, 0) { // WM_QUIT返回0，出错返回-1
				break // 0x004D8F6C
			}
			win.TranslateMessage(&ebp28msg)
			win.DispatchMessage(&ebp28msg) // 到WndProc
		} else {
			// 0x004D8D80
			if v012E2224displayMode == displayModeFullScreen { // 0 means hack?
				// 0x004D8D91 0x004D8D95
				v01319D9C++
				if v01319D9C > 30 {
					v01319D98++
					if v01319D98 > 30 {
						v01319D9C = 0
						v01319D98 = 0
						if v01319DB4 == 0 {
							// 发消息
							var inhack pb
							inhack.f00439178init()
							inhack.f0043922CwriteHead(0xC1, 0xF1)
							inhack.f004397B1writeUint8(3) // subcode
							inhack.f004397B1writeUint8(1) // flag
							inhack.f004397B1writeUint8(0) // subflag
							inhack.f004393EAsend(true, false)
							inhack.f004391CF()
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
			// 0x004D8E7F 0x004D8E83
			if v012E2224displayMode == displayModeWindow || v01319D95 == true {
				f004E6233handleState(v01319D74hDC) // 状态机 重要
			} else if v012E2224displayMode == displayModeFullScreen {
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
		// 0x004D8F4B
		f00760292parse(&v08C88FF0conn, 0, 0) // 从连接中获取报文并解析
		// v01319D2C.f004B9492()
	}
	// 0x004D8F6C
	// 退出消息循环
	// f00535963()
	if v01319D68gg != nil {
		//v01319D68gg.f004D9335(1) // 释放
		v01319D68gg = nil
	}
	win.CloseHandle(ebpC)
	//f004D56B2()
	// ebp1.f004D9F88()
	return int(ebp28msg.WParam) // return ebp20
}

func _f00A49E40init1() {
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

func f004D9F88() {}
