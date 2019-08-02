package main

import (
	"os"
	"unsafe"

	"github.com/xujintao/balgass/client/dll"
)

type conf struct {
	a uint32
	b uint32
}

const (
	_0114D912 uint8  = 0
	_0114D913 uint8  = 0
	_0114DD64 uint16 = 0x0061
)

// 这里使用数组，确保在data段，固定在某个地址上
// 如果使用string，那么读文件后对buf进行string类型转换，会memmove到堆上
var _012E2210 uint32 = 1
var _012E2214 string = "MUCN"
var _012E2224 uint32 = 1
var _012E2228 uint32 = 16
var _012E222C *os.File
var _012E2338_ip string = "192.168.0.100"
var _012E233C_port uint16 = 44405
var _012E3F08 struct {
	height uint32 // 0x258, 600
	width  uint32 // 0x320, 800
}

var _012F7910 [100]uint8
var _01319A38 [8]uint8 // "1.04R+"	// 可执行程序版本号
var _01319A44 [8]uint8 // "1.04.44" // 配置文件版本号
var _01319A50_ip [16]uint8

var _01319D1C_oWndProc uint32
var _01319D68 *uint32
var _01319D6C_hWnd uint32 // hWnd, 0x0073,1366
var _01319D70_hModule uint32
var _01319D74_hDC uint32   // hDC, 0x1401,11A6
var _01319D78_hGLRC uint32 // hGLRC, OpenGL rendering context, 0x0001,0000
var _01319DB8 [64]uint8
var _01319DFC string = "Chs"
var _01319E00 string = "chinese"

type version struct {
	major  uint16 // ebp_244
	minor  uint16 // ebp_242
	patch  uint16 // ebp_240
	fourth uint16 // ebp_23E
}

func _004D52AB(fileName []uint8, cmd string) bool {
	// 从命令行字符串中提取出应用程序名，也就是main.exe，存到fileName数组中
	return true
}

func _004D5368() bool {
	var fileName [260]uint8 // ebp-110
	cmd := os.Args[0]       // ebp-4
	_004D52AB(fileName[:], cmd)

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
	_012E222C = file
	return true
}

func _004D55C6(fileName []uint8, ver *version) bool {
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
func _004D5F98_WndProc(hWnd uint32, message uint32, wParam uint32, lParam uint32) uint32 {
	switch message {
	case dll.WM_CREATE_1:
	case dll.WM_USER_400:
		switch lParam {
		case 1:
			_08C88FF0._006BDA03_Read()
		case 2:
			_08C88FF0._006BD945_Write()
		case 20:
		}
	default:
		return DefWindowProc(hWnd, message, wParam, lParam)
	}
	return 0
}

// new callback
func _004D6C2B_WndProc(hWnd uint32, message uint32, wParam uint32, lParam uint32) uint32 {
	switch message {
	case dll.WM_CHAR_102:
		_00A49798(wParam, lParam)
	default:
		CallWindowProcw(_01319D1C_oWndProc, hWnd, message, wParam, lParam)
	}
	return 0
}

func _004D6D64() bool {
	// 40+4个字节局部变量

	// typedef struct tagPIXELFORMATDESCRIPTOR {
	// 	WORD  nSize; //28h
	// 	WORD  nVersion; // 1
	// 	DWORD dwFlags; // 25h
	// 	BYTE  iPixelType; // 0
	// 	BYTE  cColorBits;	//10h
	// 	BYTE  cRedBits;
	// 	BYTE  cRedShift;
	// 	BYTE  cGreenBits;
	// ...
	// 	BYTE  bReserved;
	// 	DWORD dwLayerMask;
	// 	DWORD dwVisibleMask;
	// 	DWORD dwDamageMask;
	//   } PIXELFORMATDESCRIPTOR, *PPIXELFORMATDESCRIPTOR, *LPPIXELFORMATDESCRIPTOR;
	// 等效
	var pixelfd struct {
		data [40]uint8
	}
	var pfdindex uint32 // ebp-4

	_00DE8100(pixelfd.data[:], 0, 0x28) // 清零

	_01319D74_hDC = GetDC(_01319D6C_hWnd) // 从hWnd中拿到hDC
	if _01319D74_hDC == 0 {
		// var nRet int32 = GetLastError()
		// _01319E08_log._00B38AE4("OpenGL Get DC Error - ErrorCode : %d\n\r", nRet)
		// _004D51C8()
		// nRet = _00436DA8(4, "OpenGL Get DC Error.", 30)
		// MessageBoxA(0, nRet)
		return false
	}

	pfdindex = ChoosePixelFormat(_01319D74_hDC, &pixelfd) // return 4
	if pfdindex == 0 {
		// var nRet int32 = GetLastError()
		// _01319E08_log._00B38AE4("OpenGL Choose Pixel Format Error - ErrorCode : %d\n\r", nRet)
		// _004D51C8()
		// nRet = _00436DA8(4, "OpenGL Choose Pixel Format Error.", 30)
		// MessageBoxA(0, nRet)
		return false
	}

	if !SetPixelFormat(_01319D74_hDC, pfdindex, &pixelfd) { // return true
		// var nRet int32 = GetLastError()
		// _01319E08_log._00B38AE4("OpenGL Set Pixel Format Error - ErrorCode : %d\n\r", nRet)
		// _004D51C8()
		// nRet = _00436DA8(4, "OpenGL Set Pixel Format Error.", 30)
		// MessageBoxA(0, nRet)
		return false
	}

	_01319D78_hGLRC = wglCreateContext(_01319D74_hDC) // return 0x0001,0000
	if _01319D78_hGLRC == 0 {
		// OpenGL Create Context Error
		return false
	}

	if !wglMakeCurrent(_01319D74_hDC, _01319D78_hGLRC) { // return true
		// OpenGL Make Current Error
		return false
	}

	// ShowWindow(_01319D6C_hWnd, 5)
	// SetForegroundWindow(_01319D6C_hWnd)
	// SetFocus(_01319D6C_hWnd)

	return true
}

func _004D6F82_initWindow(hModule uint32, iCmdShow uint32) uint32 {

	// 这里windows那边使用的是数组，编译器会使用movsw和movsb来给字符数组赋值
	var ebp_78 string = "MU"

	// ...

	// 140字节局部变量
	// typedef struct tagWNDCLASSEXA {
	// 	UINT      cbSize;
	// 	UINT      style;
	// 	WNDPROC   lpfnWndProc;
	// 	int       cbClsExtra;
	// 	int       cbWndExtra;
	// 	HINSTANCE hInstance;
	// 	HICON     hIcon;
	// 	HCURSOR   hCursor;
	// 	HBRUSH    hbrBackground;
	// 	LPCSTR    lpszMenuName;
	// 	LPCSTR    lpszClassName;
	// 	HICON     hIconSm;
	// } WNDCLASSEXA, *PWNDCLASSEXA, *NPWNDCLASSEXA, *LPWNDCLASSEXA;

	ebp_34 := struct {
		cbSize        uint32
		style         uint32
		lpfnWndProc   func()
		cbClsExtra    int32
		cbWndExtra    int32
		hInstance     uint32
		hIcon         uint32
		hCursor       uint32
		hbrBackground uint32
		lpszMenuName  string
		lpszClassName string
		hIconSm       uint32
	}{
		cbSize:        0x30, // ebp-34
		style:         0x2B, // CS_OWNDC | CS_DBLCLKS | CS_HREDRAW | CS_VREDRAW
		lpfnWndProc:   _004D5F98_WndProc,
		cbClsExtra:    0,
		cbWndExtra:    0,
		hInstance:     GetModuleHandle(0),
		hIcon:         LoadIcon(0, 0x65),     // IDI_APPLICATION)
		hCursor:       LoadCursor(0, 0x7F00), // IDC_ARROW
		hbrBackground: SetStockObject(4),     // WHITE_BRUSH
		lpszMenuName:  "",
		lpszClassName: ebp_78,     // "MU"
		hIconSm:       0x141F0439, // ebp-8
	}
	RegisterClassEx(&ebp_34)

	// ...
	var ebp_4 uint32 = CreateWindowEx( // CreateWindow会发大概3条消息给WndProc
		0,
		ebp_78,     // "MU"
		ebp_78,     // "MU"
		0x02CA0000, // WS_CLIPCHILDREN | WS_CAPTION | WS_MINIMIZEBOX | WS_SYSMENU
		0x118,
		0x46,
		0x326,
		0x274,
		0,
		0,
		0x00400000,
		0,
	)
	_01319D1C_oWndProc = SetWindowLong(ebp_4, -4, _004D6C2B_WndProc) // GWL_WNDPROC
	return 0
}

// 这个函数被花的厉害，难道是核心业务？
// 追踪只能追踪一个函数
func _004D7281() bool {
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
	buf[0] = _0114D912           // 0
	_00DE8100(buf[1:], 0, 0x113) // 0x113=275

	// ebp-108 ~ ebp-9
	var bufDir [256]uint8
	GetCurrentDirectory(0x100, bufDir[:])

	_00DE8000_strcpy(buf[:], bufDir[:])

	// _00DE7C00
	nameLen := func(bufDir []uint8) uint32 {
		return uint32(len(bufDir))
	}(bufDir[:])
	if bufDir[nameLen-1] == 0x5C {
		_00DE8010_strcat(buf[:], "config.ini")
	}
	_00DE8010_strcat(buf[:], "\\config.ini") // cat拼凑字符串

	GetPrivateProfileString("LOGIN", "Version", &_0114D913, _01319A44[:], 8, buf[:]) // 1.04.44写到全局变量中
	// ebp_8的底层数组在堆上
	// var ebp_8 string = GetCommandLine()
	var ebp_8 string = os.Args[0]

	// _004D52AB
	bRet := _004D52AB(fileName[:], ebp_8)
	if bRet {
		bRet := _004D55C6(fileName[:], &ebp_334)
		if bRet { // je 0x004D,741B
			_00DE817A_sprintf(_01319A38[:], "%d.%02d", ebp_334.major, ebp_334.minor)

			if ebp_334.patch > 0 { // jle 0x004D,7419
				*(*uint16)(unsafe.Pointer(&ebp_338)) = _0114DD64
				ebp_338.patch[2] = 0

				// 我的猜测是从patch=26开始使用字母和+标记标识patch，比如
				// 27就是A+
				// 28就是B+
				// 44就是R+
				if ebp_334.patch > 0x1A { // jle 0x004D,73EE
					ebp_338.patch[0] = 'A' - 27 + uint8(ebp_334.patch) // 65-27+44=82='R'
					ebp_338.patch[1] = '+'
				}
				_00DE8010_strcat(_01319A38[:], string(ebp_338.patch[:]))
			}
		}
	}

	// _004D7174
	// 这个函数在搞什么？
	func(fileName string) bool {
		// 204个字节的局部变量
		_01319DB8[0] = 0
		var ebp_CC uint32 = 0

		var partition struct {
			num  [64]uint8
			name [64]uint8
			ip   [64]uint8
		}
		var ebp_8 uint32
		var ebp_4 []uint8

		// _00DE909E
		ebp_4 = func(fileName string, x string) []uint8 {
			// _00DE8FDA
			return func(fileName string, x string, y uint32) []uint8 {
				// 这个函数会注册一个SEH
				return _012F7910[:]
			}(fileName, x, 0x40)
		}(fileName, "rt")
		if ebp_4 == nil {
			return false
		}

		// 不能超过50条记录
		for ebp_CC < 50 {
			// _00DECD20
			nRet := _00DECD20(ebp_4, "%s", partition.num[:])
			if nRet == -1 {
				break
			}
			nRet = _00DECD20(ebp_4, "%s", partition.name[:])
			if nRet == -1 {
				break
			}
			nRet = _00DECD20(ebp_4, "%s", partition.ip[:])
			if nRet == -1 {
				break
			}

			// _00DE94F0
			nRet = func(x []uint8, y []uint8) int32 {
				return -1
			}(_01319A50_ip[:], partition.name[:])
			if nRet == 0 {
				// _00DECBD1
				nRet := func(x []uint8) uint32 {
					return 0
				}(partition.ip[:])
				if nRet == uint32(_012E233C_port) {
					_00DE8000_strcpy(_01319DB8[:], partition.num[:])
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

func _004D755F(haystack []uint8, y int, buf []uint8) []uint8 {
	return nil
}

// 该函数被花掉了
func _004D7A1F(haystack []uint8, ip []uint8, z *uint32) uint16 {
	// x=0x0B37,3D25，存放的是[0,0,0,'K']

	// ebp-100
	var buf [256]uint8
	buf[0] = 0
	_00DE8100(buf[1:], 0, 0xFF)

	// _004472C9
	sRet := func(haystack []uint8, needle string) []uint8 {
		return _00DE92E0_strstr(haystack, needle)
	}(haystack, "battle")
	if len(sRet) == 0 {
		sRet = _004D755F(haystack, 0x79, buf[:])
		if len(sRet) == 0 {
			sRet = _004D755F(haystack, 0x75, buf[:])
			if len(sRet) == 0 {
				return 0
			}
		}
	}

	// ...
	return 1000
}

func _004D7CE5_WinMain(hModule uint32, hPrevInstance uint32, szCmdLine []uint8, iCmdShow uint32) {
	// 1B60h，近两页也就是8k字节的局部变量

	// _004D7CE5， 设置一下栈

	// ebp-4 和 ebp-8是什么？
	var ebp_4 uint32 = 0x004D7CF2 // 一个函数指针
	var ebp_8 uint32 = 0x0A

	// ebp-28
	var msg struct {
		hWnd     uint32    // ebp-28
		message  uint32    // ebp-24
		wParam   uint32    // ebp-20
		lParam   uint32    // ebp-1C
		time     uint32    // ebp-18
		pt       [2]uint32 // ebp-14, ebp-10
		lPrivate uint32    // ebp-C
	}

	// if _00DE7C00() < 1 { // jae -> jmp
	// 启动mu.exe然后退出
	// }

	// vc编译器把先定义的变量地址在高位置？
	var ebp_230 struct {
		data [248]uint8
	}
	var ebp_238_ver [8]uint8
	copy(ebp_238_ver[:], "unknown")
	_00DE8100(ebp_230.data[:], 0, 0xF8)

	// var cmd string = GetCommandLine()
	ebp_23C_cmd := os.Args[0]

	var ebp_244_ver version // ebp-244 ~ ebp-23E，这里会有清零操作，应该是程序员自己清的
	var ebp_350_fileName [260]uint8
	if _004D52AB(ebp_350_fileName[:], ebp_23C_cmd) {
		if _004D55C6(ebp_350_fileName[:], &ebp_244_ver) {
			_00DE817A_sprintf(ebp_238_ver[:], "%d.%02d", ebp_244_ver.major, ebp_244_ver.minor)
			if ebp_244_ver.patch > 0 {
				var ebp_4DC [2]uint8
				ebp_4DC[0] = uint8('a' + ebp_244_ver.patch - 1) // 0x8C，这里是有问题的，超过ASCII编码范围了
				ebp_4DC[1] = 0                                  // 我猜测作者是想表达R
				_00DE8010_strcat(ebp_238_ver[:], string(ebp_4DC[:]))
			}
		}
	}

	_01319E08_log._00B38AE4("\r\n")
	_01319E08_log._00B38D31_begin()
	_01319E08_log._00B38D19_cut()

	_01319E08_log._00B38AE4("Mu online %s (%s) executed. (%d.%d.%d.%d)\r\n",
		ebp_238_ver[:], _01319DFC,
		ebp_244_ver.major, ebp_244_ver.minor, ebp_244_ver.patch, ebp_244_ver.fourth)

	_01319E08_log._00B38D49(1)

	// system information

	// direct-x information

	// 这里重新设置ip地址和端口失败了，为什么？
	// 从命令行中提取指定的ip地址和端口，因为没有通过mu.exe启动，所有没有命令行
	var port uint16 = _004D7A1F(szCmdLine, _01319A50_ip[:], &ebp_8)
	if port > 0 {
		_012E2338_ip = string(_01319A50_ip[:])
		_012E233C_port = port
	}

	// open main.exe
	if !_004D5368() {
		_004D9F88()
		return
	}

	// enc and dec init
	_08C8D050_enc._00B62CF0_init("Data/Enc1.dat")
	_08C8D098_dec._00B62D30_init("Data/Dec2.dat")

	// read config.ini
	_01319E08_log._00B38AE4("> To read config.ini.\r\n")
	bRet := _004D7281()
	if !bRet {
		_01319E08_log._00B38AE4("config.ini read error\r\n")
		_004D9F88()
		return
	}

	// gg init
	// ebp_1A38也可能是个实例指针变量，值是0x0D9F,EA10
	var ebp_1A38 *uint32 = _00DE852F(1)
	var ebp_1B0C *uint32
	if ebp_1A38 == nil { // if true，disable GameGurad
		ebp_1B0C = nil
	} else {
		// ebp_1B0C = _004D936A(_012E2214)
	}

	_01319D68 = ebp_1B0C
	bRet = _00B4C1B8()
	if !bRet { // if false, disable GameGurad
		// gg init error
		// _004D51C8()
		// _004D9335()
		// _004D9f88()
		return
	}
	_01319E08_log._00B38AE4("> gg init success.\r\n")
	_01319E08_log._00B38D19_cut()

	if _012E2210 == 1 {
		// ShowCursor(0)
	}
	if _012E2224 == 0 {
		// DisplaySetting
	}

	// window init
	_01319E08_log._00B38AE4("> Screen size = %d x %d.\r\n", _012E3F08.height, _012E3F08.width)
	_01319D70_hModule = hModule
	_01319D6C_hWnd = _004D6F82_initWindow(hModule, iCmdShow) // 创建hWnd
	_01319E08_log._00B38AE4("> Start window success.\r\n")

	// opengl init
	bRet = _004D6D64() // 设置hWnd资源
	if !bRet {
		// 失败处理
		return
	}
	_01319E08_log._00B38AE4("> OpenGL init success.\r\n")
	_01319E08_log._00B38D19_cut()

	// opengl information
	_01319E08_log._00B38E3C()
	_01319E08_log._00B38D19_cut()

	// sound card information
	_01319E08_log._00B3902D() // 自带cut

	// ShowWindow(_01319D6C_hWnd, 10) // 发很多消息给WndProc
	// UpdateWindow(_01319D6C_hWnd) // 发很多消息给WndProc

	_00B4C1FF(_01319D6C_hWnd)
	_01319E08_log._00B38AE4("> gg connect window Handle.\r\n")
	_01319E08_log._00B38D19_cut()

	// ime information
	_01319E08_log._00B38EF4(_01319D6C_hWnd)
	_01319E08_log._00B38D19_cut()

	//
	_00DEE871_setlocale(0, _01319E00) // 参数0是哪个LC？
	// _0043BF3F(_01319D6C_hWnd, _012E3F08.height, _012E3F08.width)
	// _0043BF9C()
	// _0090E94C()
	// _0090B256()
	// _00A49798(_01319D6C_hWnd, _012E3F08.height, _012E3F08.width)
	// _00A49E40() // r6602 floating point support not loaded

	// ...

	_01319E08_log._00B38AE4("> Loading ok.\r\n")

	// ...

	var ebp_595 uint8 // ImmIsUIMessage返回值

	// 消息循环
	for {
		if msg.message == 0x10 || // WM_CLOSE
			msg.message == 0x12 { // WM_QUIT
			break
		}

		if PeekMessage(&msg, 0, 0, 0) {
			ImmIsUIMessage()
			if msg.message == 0x100 || // WM_KEYFIRST, WM_KEYDOWN
				msg.message == 0x101 || // WM_KEYUP
				ebp_595 != 0 ||
				msg.message == 0x201 || // WM_LBUTTONDOWN
				msg.message == 0x202 { // WM_LBUTTONUP
				_00A49798(msg.hWnd, msg.message, msg.wParam, msg.lParam, 1)
				_00A4DFB7()
			}

			if !GetMessage(&msg, 0, 0, 0) { // WM_QUIT返回0，出错返回-1
				break
			}
			TranslateMessage(&msg)
			DispatchMessage(&msg) // 到WndProc
		} else {
			if _012E2224 == 0 {

			} else if _012E2224 == 1 {

			}
		}
		_00760292(&_08C88FF0, 0, 0) // 比较重要
		_01319D2C._004B9492()
	}
}

func _004D9F88() {

}

// 也可能是map
var _0075FF6A_cmd []struct{}
var _0075FCB2_cmd []struct{}

// 这个函数被花了
func _00760292(t *t2000, x, y uint32) {
	var ebp18 []uint8 // t.t2002
	var ebp14 uint32
	var ebp10 uint32
	switch ebp18[0] {
	case 0xC1:
		ebp14 = ebp18[2]
	case 0xC2:
		ebp14 = ebp18[3]
		ebp10 = ebp18[1]<<8 + ebp18[2]
	case 0xC3:
		ebp10 = ebp18[1]
		_08C8D014._00B98FC0(t, &ebp181A)
	case 0xC4:
	}
	_086A3C28 += ebp10 // 累计接收字节数
	if y == 1 {
	}

	ebp1820 := 0

	// _0075C3B2 被花到 0x009FE6362
	func(typ uint32, buf []uint8, len uint32, x uint32) {
		ebp5D8 := typ
		if ebp5D8 > 0xFD {
		}

		// 寻找处理器
		_0075FCB2_cmd[4*_0075FF6A_cmd[ebp5D8]]

		// 0x00 处理器
		// _006FA9EB
		func(buf []uint8) {
			// 业务逻辑很复杂
			// send c1 04 f4 06
		}(buf)

		// 0xF1 处理器
		ebp4 = buf
		ebp5DC := ebp4[3]
		if ebp5DC > 4 {
		}
		switch ebp5DC {
		case 4:
		case 0: // version match, buf: 01 2E 87 "10525"
			// _006C75A7 被花到 0x0A4E7A49
			func(buf []uint8) {
				ebp4 := buf
				ebp10 := buf[4]
				if ebp10 == 1 {
					ebp8 := &_0130F728
					ebp8._004A9123(&ebp8.f4880)
				}
				_08C88E0C = ebp4[5]<<8 + ebp4[6]
				_08C88E08 = 2
				ebpC := 0
				if ebpC >= 5 {
				}
				for {
					if _012E4018[ebpC]-ebpC-1 != ebp4[7+ebpC] { // 改为jmp
						break
					}
					ebpC++
					if ebpC >= 5 {
						break
					}
				}

				// _004A9146
				func() {}()

				// _004A9EEB
				func() {}()
				_01319E08_log._00B38AE4("Version dismatch - Join server.\r\n")
				// _0051FB61
				func() {}()
			}(buf)
		}

		// 0xFA 处理器(没有)

		// 0xF4 处理器
		if buf[0] == 0xC1 {
		}
		ebp1C := buf
		ebp14 := ebp1C[4]
		ebp5EC := ebp14
		switch ebp5EC {
		case 3:
		case 5:
		case 6:
			// _006BFA3E
			func(buf []uint8) {
				ebp4 := 6
				ebp9 := buf[ebp4]
				ebp4++
				if ebp14 >= buf[5]<<8|1 {
				}
				ebp18 := buf[ebp4:]
			}(buf)
		}
	}(ebp14, ebp18, ebp10, ebp1820)

	// _006BDC33
	func() {}()
}

// 很复杂 0x004E4F1C ~ 0x004E6232，不到5000行并且几乎都是函数调用
func _004E4F1C() {

	// _004E1CEE
	func() {
		// _006BF89A 连接
		func() {

			_01319E08_log._00B38AE4()
		}()
	}()

}

// 进度加载完毕handle
func _004E6233() {

	_004E4F1C()

}
