package main

import (
	"os"
	"unsafe"
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
var _01319A50_ip string

var _01319D68 uint32
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

func _004D6F82(hModule uint32, x uint32) uint32 {
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
			nRet = func(x string, y []uint8) int32 {
				return -1
			}(_01319A50_ip, partition.name[:])
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
func _004D7A1F(haystack []uint8, ip string, z *uint32) uint16 {
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

func _004D7CE5_run(hModule uint32, ebpC uint32, haystack []uint8, ebp14 uint32) {
	// 1B60h，近两页也就是8k字节的局部变量

	// _004D7CE5， 设置一下栈

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

	var ebp_8 uint32 = 10
	var port uint16 = _004D7A1F(haystack, _01319A50_ip, &ebp_8) // 这里重新设置ip地址和端口失败了
	if port > 0 {
		_012E2338_ip = _01319A50_ip
		_012E233C_port = port
	}

	// open main.exe
	if !_004D5368() {
		_004D9F88()
		return
	}

	keyenc1 := new([54]uint8) // 0x08C8,D050
	_00B62CF0("Data/Enc1.dat", keyenc1[:])

	keydec2 := new([54]uint8) // 0x08C8,D098
	_00B62D30("Data/Dec2.dat", keydec2[:])

	_01319E08_log._00B38AE4("> To read config.ini.\r\n")
	bRet := _004D7281()
	if !bRet {
		// _00B38AE4
		// _004D9F88
		return
	}

	_00DE852F(1)

	_00B4C1B8()

	// gg
	// 已经被禁用了
	_01319E08_log._00B38AE4("> gg init success.\r\n")
	_01319E08_log._00B38D19_cut()

	// window
	_01319E08_log._00B38AE4("> Screen size = %d x %d.\r\n", _012E3F08.height, _012E3F08.width)
	_01319D70_hModule = hModule
	_01319D6C_hWnd = _004D6F82(hModule, ebp14) // 创建hWnd
	_01319E08_log._00B38AE4("> Start window success.\r\n")

	// opengl
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

	// ShowWindow(_01319D6C_hWnd, 10)
	// UpdateWindow(_01319D6C_hWnd)

	_00B4C1FF(_01319D6C_hWnd)
	_01319E08_log._00B38AE4("> gg connect window Handle.\r\n")
	_01319E08_log._00B38D19_cut()

	// ime information
	_01319E08_log._00B38EF4(_01319D6C_hWnd)
	_01319E08_log._00B38D19_cut()

	_00DEE871(0, _01319E00)
	_0043BF3F(_01319D6C_hWnd, _012E3F08.height, _012E3F08.width)
	_0043BF9C()
	_0090E94C()
	_0090B256()
	_00A49798(_01319D6C_hWnd, _012E3F08.height, _012E3F08.width)
	_00A49E40() // r6602 floating point support not loaded

	// ...

	_01319E08_log._00B38AE4("> Loading ok.\r\n")
}

func _004D9F88() {

}
