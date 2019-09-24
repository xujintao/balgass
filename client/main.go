package main

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/xujintao/balgass/client/win"
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
var v012E2338ip = "192.168.0.100"
var v012E233Cport uint16 = 44405
var v012E2340 int

var v012E3F08 struct {
	height uint32 // 0x258, 600
	width  uint32 // 0x320, 800
}

var v012F7910 [100]uint8
var v01319A38 [8]uint8 // "1.04R+"	// 可执行程序版本号
var v01319A44 [8]uint8 // "1.04.44" // 配置文件版本号
var v01319A50ip [16]uint8

var v01319D1CwndProc uintptr
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

type t4 struct {
	f50h func([]uint8, int)
}

type t3 struct {
	f200h *t4
	f204h *t4
}

func (t *t3) f0043E60C() {
	// ...
	t.f0043E9B6()
	// ...
}

// s9 006E855E
func (t *t3) f0043E9B6() {
	// ...
	f00DE8A70chkstk() // 0x152C
	if v08C88E08 == 0 {
		return
	}

	f004A7D34()
	f004A9146()

	var ebp18, ebp24 [11]uint8 // username, pwd
	f00DE8100memset(ebp18[:], 0, 11)
	t.f200h.f50h(ebp18[:], 11)       // f00342929, fill username
	f00DE8100memset(ebp24[:], 0, 11) // f00452929, fill pwd
	t.f204h.f50h(ebp24[:], 11)

	// f0043EE13(ebp18[:]) // 判断字符串长度
	// f0043EE13(ebp24[:])
	if v08C88E08 != 2 {
		return
	}

	v01319E08log.f00B38AE4printf("> Login Request.\r\n")
	v01319E08log.f00B38AE4printf("> Try to Login \"%s\"\r\n", ebp18[:])

	v08C88F74 = 1
	f00DE8000strcpy(v08C88F78username[:], ebp18[:])
	v08C88E08 = 13

	// 构造登录报文
	var ebp14BClogin pb
	ebp14BClogin.f00439178init()
	ebp14BClogin.f0043922CwritePrefix(0xC1, 0xF1) // 前缀
	ebp14BClogin.f004397B1writeUint8(1)           // 可能是subcode

	var ebp30, ebp3C [11]uint8
	f00DE8100memset(ebp30[:], 0, 11)
	f00DE8100memset(ebp3C[:], 0, 11)
	f00DE7C90memcpy(ebp30[:], ebp18[:], 10)
	f00DE7C90memcpy(ebp3C[:], ebp24[:], 10)
	f0043B750xor(ebp30[:], 10) // 与igc.dll的TOOLTIP_FIX_XOR_BUFF一样了，可能有问题
	f0043B750xor(ebp3C[:], 10) // 与igc.dll的TOOLTIP_FIX_XOR_BUFF一样了，可能有问题

	ebp14BClogin.f00439298writeBuf(ebp30[:], 10, true)    // 写username
	ebp14BClogin.f00439298writeBuf(ebp3C[:], 10, true)    // 写pwd
	ebp14BClogin.f0043EDF5writeUint32(win.GetTickCount()) // 写时间戳
	ebp14C0 := 0
	for ebp14C0 < 5 {
		ebp14BClogin.f004397B1writeUint8(v012E4018[ebp14C0] - byte(ebp14C0+1)) // 写版本 VERSION_HOOK1
		ebp14C0++
	}
	ebp14C0 = 0
	for ebp14C0 < 16 {
		ebp14BClogin.f004397B1writeUint8(v012E4020[ebp14C0]) // 写序列号 SERIAL_HOOK1
		ebp14C0++
	}

	// 发送登录报文
	ebp14BClogin.f004393EAsend(true, false)

	// ...
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
func f004A7D34() *t4000 { return nil }
func f004A9146()        {}

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
			v08C88FF0.f006BDA03read()
		case 2:
			v08C88FF0.f006BD945write()
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

// 这个函数被花的厉害，难道是核心业务？
// 追踪只能追踪一个函数
func f004D7281() bool {
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

	win.GetPrivateProfileStringA("LOGIN", "Version", v0114D913, v01319A44[:], 8, string(buf[:])) // 1.04.44写到全局变量中
	// ebp_8的底层数组在堆上
	// var ebp_8 string = GetCommandLine()
	var ebp_8 string = os.Args[0]

	// f004D52AB
	bRet := f004D52AB(fileName[:], ebp_8)
	if bRet {
		bRet := f004D55C6(fileName[:], &ebp_334)
		if bRet { // je 0x004D,741B
			f00DE817Asprintf(v01319A38[:], "%d.%02d", ebp_334.major, ebp_334.minor)

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
				f00DE8010strcat(v01319A38[:], string(ebp_338.patch[:]))
			}
		}
	}

	// _004D7174
	// 这个函数在搞什么？
	func(fileName string) bool {
		// 204个字节的局部变量
		v01319DB8[0] = 0
		var ebp_CC uint32 = 0

		var partition struct {
			num  [64]uint8
			name [64]uint8
			ip   [64]uint8
		}
		// var ebp_8 uint32
		var ebp_4 []uint8

		// _00DE909E
		ebp_4 = func(fileName string, x string) []uint8 {
			// _00DE8FDA
			return func(fileName string, x string, y uint32) []uint8 {
				// 这个函数会注册一个SEH
				return v012F7910[:]
			}(fileName, x, 0x40)
		}(fileName, "rt")
		if ebp_4 == nil {
			return false
		}

		// 不能超过50条记录
		for ebp_CC < 50 {
			// f00DECD20
			nRet := f00DECD20(ebp_4, "%s", partition.num[:])
			if nRet == -1 {
				break
			}
			nRet = f00DECD20(ebp_4, "%s", partition.name[:])
			if nRet == -1 {
				break
			}
			nRet = f00DECD20(ebp_4, "%s", partition.ip[:])
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
			ebp_CC++
		}

		// _00DE8C84
		func(x []uint8) {
			// 清零
		}(ebp_4)

		return false
	}("partition.inf")

	// 接下来还有一段处理，我没看

	return true
}

func f004D755F(haystack []uint8, y int, buf []uint8) []uint8 {
	return nil
}

// 该函数被花掉了
func f004D7A1F(haystack []uint8, ip []uint8, z *uint32) uint16 {
	// x=0x0B37,3D25，存放的是[0,0,0,'K']

	// ebp-100
	var buf [256]uint8
	buf[0] = 0
	f00DE8100memset(buf[1:], 0, 0xFF)

	// _004472C9
	sRet := func(haystack []uint8, needle string) []uint8 {
		return f00DE92E0strstr(haystack, needle)
	}(haystack, "battle")
	if len(sRet) == 0 {
		sRet = f004D755F(haystack, 0x79, buf[:])
		if len(sRet) == 0 {
			sRet = f004D755F(haystack, 0x75, buf[:])
			if len(sRet) == 0 {
				return 0
			}
		}
	}

	// ...
	return 1000
}

// f004D7CE5winMain, WinMain
func f004D7CE5winMain(hModule win.HMODULE, hPrevInstance uint32, szCmdLine []uint8, iCmdShow int) {
	// 0x0A05E61B
	var label1 uint32 = 0x009DCA19
	// push label1
	// push 0x0ABDA0EA
	// ret

	// 0x0ABDA0EA
	var label2 uint32 = 0x0A3349B0
	// push label2
	// push 0x0A4420A9
	// ret

	// 0x0A4420A9
	var label3 uint32 = 0x0A0518C4
	// push 0x0A051BC4
	// push 0x0A8485E3
	// ret

	// 0x0A8485E3
	// pushfd
	// push ebx
	// push edi
	// push esi
	// push edx
	// push ecx

	// 0x09EBC742
	// push 0x0AD2B081
	// push 0x0AD91CED
	// ret

	// 0x0AD91CED
	func() {
		// 0xFC局部变量
		// push ebx
		// push esi
		// push edi
		var ebp8 uint32
		var ebp14 uint32
		//var ebp18 uint32
		var ebp20 uint32
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
			// 0x0A9F5A38
			// 0x0A6B0AEB
			// ...
		} else {
			// 0x0A9F5A3E

			ebp14 = v0A74573D

			// 0x0AA2F062
			ebp44 = v0A74573D
			// 0x0A32AFCE
			for {
				// 0x0A058220
				// 0x0AAB85DB
				// 0x0AF77F6A
				// 0x0AD80BB2
				// 0x0A391FF8
				// 0x0AD3B8F3
				// 0x0A3927B1
				// 0x09EBC1FB
				if ebp44 <= 0 { // 迭代结束, ebp44=0, ebp34=6
					// 0x0AF814F7
					break // 0x0A83CE6F
				} else {
					// 0x0A442AC6
					// 0x0AF7F68B
					// 0x0A84320A
					// 0x0A8902CD
					// 0x0B071558
					// 0x0AD98514
					// 0x0ABF87EB
					if ebp44&1 != 0 {
						// 0x09E91EC8
						ebp34++ // 1,2,3,4,5,6
					}
					// 0x09FDAD42
					// 0x0AC3847E
					ebp44 >>= 1
					// 0x0A058220
				}
			}
		}

		// 0x0A83CE6F
		// 0x0A38F2A7
		// 0x0AF11D82
		// 0x0A565A51
		if ebp34 == 0 {
			// 0x0A8FA42A
			// 0x0A6B0AEB
		}
		// 0x0A84D90D
		ebp38 = 0
		// 0x0A32AF58
		v0A38DB63 = v0AF77CE1 ^ v0B1090BA
		// 0x0A38DB63
		var tscLeax uint32 = 0x5274DC2C
		//var tscHedx uint32 = 0x00010865
		// 0x0A32F0F4
		ebp20 = tscLeax
		// 0x0AF8EBF1
		ebp20 |= ebp20 >> 16 // 0x5274DE7C
		// 0x0A56E64F
		// eax = ebp20
		// edx = 0
		// ecx = 0x1F
		// 0x0ABE0B07
		ebp20 %= 0x1F // 3
		ebp8 = 1
		for {
			// 0x0A9FBE5C
			// 0x0AD83B40
			// 0x09E240E7
			// 0x09FE52D1
			// 0x09FCAF92
			// 0x0B10FA31
			// 0x09E03B4A
			if ebp8&ebp14 != 0 { // 1&0x3F
				// 0x0AD98F0D
				ebp40 = ebp20 % ebp34 // 3%6
				for {
					// 0x0AD75B27
					if ebp40 <= 0 {
						break // 0x0ABB3DF5
					}
					// 0x0AF82F72
					ebp8 <<= 1 // 1,2,4,8
					// 0x0A3A7943
					if ebp8&ebp14 == 0 { // 2&0x3F, 不超过6次
						// 0x0A049688
					}
					// 0x0AC39B0A
					// 0x09E8EB25
					ebp40-- // 3, 2, 1, 0
				} // for loop 0x0AD75B27

				// 0x0ABB3DF5
				// 0x0A742C83
				// 0x0A3360C0
				// 0x0A90096C
				// 0x0A4DF380
				// 0x0A8892C0
				if ebp8&1 != 0 {
					// 0x0ABFA5B3
				}
				// 0x0A4E80E3
				// 0x0A4DDA43
				// 0x0A05AA37
				// 0x09FD9C5C
				// 0x09E30794
				// 0x0A4DDA22
				// 0x09F8E057
				// 0x0A7441CB
				// 0x0A84A77E
				if ebp8&2 != 0 {
					// 0x0A742EF1
				}
				// 0x0AD31978
				if ebp8&4 != 0 {
					// 0x0AFDC124
				}
				// 0x0A92FBD9
				if ebp8&8 != 0 {
					// 0x09E2BB81
					ebp30 = 1
					// pushad
					// push 0x0A32BD69
					// 0x0ABFA7E5
					// SEH
					// 0x0A847A61
					// push 0
					// push 0
					// int 1
					ebp30 = 1
					//ebp18 = 1
					// 0x0A904AEE
					// 0x0A4E7E87
					// esp+=4
					// 0x09E25911
					// popad
					// 0x0A327A81
					// 0x0A05C5B5
					// 0x09FC28CA
					// 0x0B10F776
					// 0x0ABB564B
					if ebp30 == 0 {
						// 0x09F8768E
						// 0x0A6B0AEB
					}
					// 0x09FC1A2B
					label2 = v09FC0B5Alabel2 // [ebp+v0ABF8B07*4+8] = v09FC0B5Alabel2
					ebp38++
					ebp8 &= ^uint32(8) // 8&(^8)
				}
				// 0x0A9353C2
				// 0x0AF8FAB1
				// 0x0B074337
				// 0x0AD9201D
				// 0x0A895B6A
				// 0x09FD5603
				// 0x0ABB3AC2
				// 0x0ABD65F8
				// 0x0AAB22A0
				if ebp8&0x10 != 0 {
					// 0x0AFD7373
				}
				// 0x0A444B3D
				// 0x09E2F73A
				// 0x0AD813DC
				// 0x0AF899B3
				// 0x0A05DA3A
				if ebp8&0x20 != 0 {
					// 0x09E6AFDD
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
					// 0x0B07BF42
					// 0x0A6B0AEB
				}
				// 0x0AD83D30
				if ebp38 != 1 {
					// 0x0A9331C5
					// 0x0A6B0AEB
				}
				// 0x0AF77BA8
				v0A325B9F = v09FBB49C
				label3 = v0A557660label3 // [ebp+v0AA0D71B*4+8] = v0A557660label3
				break                    // 0x09FBBB94
			}
			// 0x0A4E4AF6
			ebp8 <<= 1
		} // for loop 0x0A9FBE5C

		// 0x09FBBB94
		// pop edi
		// pop esi
		// pop ebx
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
	// label1(0x0AD56E8A)
	// push ebp
	// ebp = esp
	// eax = 0x1B60
	// 0x004D7CED
	// -------------------------肉区---------------------------------
	// 1B60h，近两页也就是8k字节的局部变量

	// _004D7CE5， 设置一下栈

	// ebp-4 和 ebp-8是什么？
	// var ebp_4 uint32 = 0x004D7CF2 // 一个函数指针
	var ebp_8 uint32 = 0x0A

	// ebpC := win.CreateMutex
	var ebpC win.HANDLE

	// if _00DE7C00() < 1 { // jae -> jmp
	// 启动mu.exe然后退出
	// }

	// vc编译器把先定义的变量地址在高位置？
	var ebp_230 struct {
		data [248]uint8
	}
	var ebp_238_ver [8]uint8
	copy(ebp_238_ver[:], "unknown")
	f00DE8100memset(ebp_230.data[:], 0, 0xF8)

	// var cmd string = GetCommandLine()
	ebp_23C_cmd := os.Args[0]

	var ebp_244_ver version // ebp-244 ~ ebp-23E，这里会有清零操作，应该是程序员自己清的
	var ebp_350_fileName [260]uint8
	if f004D52AB(ebp_350_fileName[:], ebp_23C_cmd) {
		if f004D55C6(ebp_350_fileName[:], &ebp_244_ver) {
			f00DE817Asprintf(ebp_238_ver[:], "%d.%02d", ebp_244_ver.major, ebp_244_ver.minor)
			if ebp_244_ver.patch > 0 {
				var ebp_4DC [2]uint8
				ebp_4DC[0] = uint8('a' + ebp_244_ver.patch - 1) // 0x8C，这里是有问题的，超过ASCII编码范围了
				ebp_4DC[1] = 0                                  // 我猜测作者是想表达R
				f00DE8010strcat(ebp_238_ver[:], string(ebp_4DC[:]))
			}
		}
	}

	v01319E08log.f00B38AE4printf("\r\n")
	v01319E08log.f00B38D31begin()
	v01319E08log.f00B38D19cut()

	v01319E08log.f00B38AE4printf("Mu online %s (%s) executed. (%d.%d.%d.%d)\r\n",
		ebp_238_ver[:], v01319DFC,
		ebp_244_ver.major, ebp_244_ver.minor, ebp_244_ver.patch, ebp_244_ver.fourth)

	v01319E08log.f00B38D49(1)

	// system information

	// direct-x information

	// 这里重新设置ip地址和端口失败了，为什么？
	// 从命令行中提取指定的ip地址和端口，因为没有通过mu.exe启动，所有没有命令行
	var port uint16 = f004D7A1F(szCmdLine, v01319A50ip[:], &ebp_8)
	if port > 0 {
		v012E2338ip = string(v01319A50ip[:])
		v012E233Cport = port
	}

	// open main.exe
	if !f004D5368() {
		f004D9F88()
		return
	}

	// enc and dec init
	v08C8D050enc.f00B62CF0init("Data/Enc1.dat")
	v08C8D098dec.f00B62D30init("Data/Dec2.dat")

	// read config.ini
	v01319E08log.f00B38AE4printf("> To read config.ini.\r\n")
	bRet := f004D7281()
	if !bRet {
		v01319E08log.f00B38AE4printf("config.ini read error\r\n")
		f004D9F88()
		return
	}

	// gg init
	// ebp_1A38也可能是个实例指针变量，值是0x0D9F,EA10
	var ebp_1A38 *uint32 = (*uint32)(f00DE852F(1))
	var ebp_1B0C *t1319D68
	if ebp_1A38 == nil { // if true，disable GameGurad
		ebp_1B0C = nil
	} else {
		// ebp_1B0C = _004D936A(v012E2214)
	}

	v01319D68 = ebp_1B0C
	bRet = f00B4C1B8()
	if !bRet { // if false, disable GameGurad
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
	// _0043BF3F(v01319D6ChWnd, v012E3F08.height, v012E3F08.width)
	// _0043BF9C()
	// _0090E94C()
	// _0090B256()
	// f00A49798(v01319D6ChWnd, v012E3F08.height, v012E3F08.width)
	// f00A49E40() // r6602 floating point support not loaded

	// ...

	v01319E08log.f00B38AE4printf("> Loading ok.\r\n")

	// ...

	var ebp_595 uint8 // ImmIsUIMessage返回值

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
				ebp_595 != 0 ||
				msg.Message == 0x201 || // WM_LBUTTONDOWN
				msg.Message == 0x202 { // WM_LBUTTONUP
				f00A49798(msg.HWnd, msg.Message, msg.WParam, msg.LParam, 1)
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
				f004E6233(v01319D74hDC) // 网络处理
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
		f00760292parse(&v08C88FF0, 0, 0) // 从连接中获取报文并解析
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

var v08C88F60 uint8
var v08C88F61 uint8
var v08C88F62 int
var v08C88F64 int
var v08C88F74 int
var v08C88F78username [11]uint8 // username
var v08C88F84 uint32

type t2 struct {
	f154 uint16 // 0
	f160 uint16 // 0
}

var v086105EC struct {
	f00 *t2 // 指向0x0D94B0AC
}

func f004D8FF5() {
	// SEH
	var ebp1494heartbeat pb
	f00DE8A70chkstk() // 0x1488

	if v08C88F64 == 0 {
		return
	}

	// 构造心跳报文
	ebp1494heartbeat.f00439178init()
	ebp1494heartbeat.f0043922CwritePrefix(0xC1, 0x0E) // 写前缀
	ebp10 := win.GetTickCount()
	ebp1494heartbeat.f0043974FwriteZero(1)                    // 写0
	ebp1494heartbeat.f0043EDF5writeUint32(ebp10)              // 写time
	ebp1494heartbeat.f004C65EFwriteUint16(v086105EC.f00.f154) // 什么
	ebp1494heartbeat.f004C65EFwriteUint16(v086105EC.f00.f160) // 什么
	ebp1494heartbeat.f004393EAsend(true, false)               // 发送心跳报文
	if v08C88F62 == 0 {
		v08C88F62 = 1
		v08C88F84 = ebp10
	}
	ebp1494heartbeat.f004391CF()
}

func f004D9F88() {}

// 也可能是map
var v0075FF6Acmd []struct{}
var v0075FCB2cmd []struct{}

var v086A3B94 func()
var v086A3C28 int

// 网络
func f004E6233(hDC win.HDC) {

	// f004E4F1C
	func() {
		// 很复杂 0x004E4F1C ~ 0x004E6232，不到5000行并且几乎都是函数调用
		// ...

		v0130F728.f004A9B5B() // 请求服务器列表

		ebp40C := v012E2340
		switch ebp40C {
		case 2:
			// f004E1E1E
			func() {
				// 带SEH
				f00DE8A70chkstk()

				// f004E1CEE
				func() {
					// _006BF89A 拨号
					func(ip string, port int) {
						v08C88FF0.f006BD3A7init()
						v01319E08log.f00B38AE4printf("[Connect to Server] ip address = %s, port = %d\r\n")
						v08C88FF0.f006BD509socket(v01319D6ChWnd, 1)
						v08C88FF0.f006BD708dial(ip, port, 400)
					}(v012E2338ip, int(v012E233Cport))
				}()
			}()
		case 4:
			// f004DDD4F
			func() {}()
		case 5:
			// f004DF0D5
			func() {}()
		}

		// ...
	}()
}
