package main

import (
	"fmt"
	"unsafe"
)

func _004D6D64() {
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

	hWnd := 0x01319D6C // hWnd是全局变量
	// hDC := GetDC(hWnd) //把hDC的值存到0x01319D74
	// if hDC == nil {
	// 	nRet := GetLastError()
	// 	_00B38AE4(0x01319E08, "OpenGL Get DC Error - ErrorCode : %d\n\r", nRet)
	// 	_004D51C8()
	// 	nRet = _00436DA8(4, "OpenGL Get DC Error.", 30)
	// 	MessageBoxA(0, nRet)
	// 	return
	// }

	// pfdindex = ChoosePixelFormat(hDC, &pixelfd) // return 4
	// if pfdindex == 0 {
	// 	nRet := GetLastError()
	// 	_00B38AE4(0x01319E08, "OpenGL Choose Pixel Format Error - ErrorCode : %d\n\r", nRet)
	// 	_004D51C8()
	// 	nRet = _00436DA8(4, "OpenGL Choose Pixel Format Error.", 30)
	// 	MessageBoxA(0, nRet)
	// 	return
	// }

	// if !SetPixelFormat(hDC, pfdindex, &pixelfd) { // return true
	// 	nRet := GetLastError()
	// 	_00B38AE4(0x01319E08, "OpenGL Set Pixel Format Error - ErrorCode : %d\n\r", nRet)
	// 	_004D51C8()
	// 	nRet = _00436DA8(4, "OpenGL Set Pixel Format Error.", 30)
	// 	MessageBoxA(0, nRet)
	// 	return
	// }

}

type conf struct {
	a uint32
	b uint32
}

type t1000 struct {
	f00h  uint32
	f04h  int32
	f08h  [260]uint8
	f10Ch uint32
}

// hError
var _01319E08 = t1000{
	f00h: 0x01180F50,
	f04h: 0xE0, //hFile
	// f08h: "MuError.log"
	f10Ch: 0x0F,
}

// 如果是结构体的话，程序会用偏移
const (
	_0114D912 uint8  = 0
	_0114D913 uint8  = 0
	_0114DD64 uint16 = 0x0061
)

// 这里使用数组，确保在data段，固定在某个地址上
// 如果使用string，那么读文件后对buf进行string类型转换，会memmove到堆上
var _012E233C uint16 // 0xAD75
var _012F7910 [100]uint8
var _01319A38 [8]uint8 // "1.04R+"	// 可执行程序版本号
var _01319A44 [8]uint8 // "1.04.44" // 配置文件版本号
var _01319A50 uint32 = 0
var _01319DB8 [64]uint8

// 这个函数被花的厉害，难道是核心业务？
// 追踪只能追踪一个函数
func _004D7281() {
	// CURRENT_USER
	// SOFTWARE\\Webzen\\Mu\\Config
	// 348h个字节的局部变量

	// ebp-338
	var ebp_338 struct {
		patch [3]uint8
	}

	// ebp-334 ~ ebp-330
	var ebp_334 struct {
		v1 uint16 // ebp_334, major
		v2 uint16 // ebp_332, minor
		v3 uint16 // ebp_330, patch
		v4 uint16 // ebp_32E
	}

	// ebp-328 ~ ebp-225
	var fileName [260]uint8

	// ebp-220 ~ ebp-10D
	var buf [276]uint8
	buf[0] = _0114D912           // 0
	_00DE8100(buf[1:], 0, 0x113) // 0x113=275

	// ebp-108 ~ ebp-9
	var bufDir [256]uint8
	GetCurrentDirectory(0x100, bufDir[:])

	_00DE8000(buf[:], bufDir[:])

	// _00DE7C00
	nameLen := func(bufDir []uint8) uint32 {
		return uint32(len(bufDir))
	}(bufDir[:])
	if bufDir[nameLen-1] == 0x5C {
		_00DE8010(buf[:], "config.ini")
	}
	_00DE8010(buf[:], "\\config.ini") // cat拼凑字符串

	GetPrivateProfileString("LOGIN", "Version", &_0114D913, _01319A44[:], 8, buf[:]) // 1.04.44写到全局变量中
	// ebp_8的底层数组在堆上
	var ebp_8 string = GetCommandLine()

	// _004D52AB
	bRet := func(fileName []uint8, cmd string) bool {
		// 从命令行字符串中提取出应用程序名，也就是main.exe，存到fileName数组中
		return true
	}(fileName[:], ebp_8)
	if bRet {
		// _004D55C6
		bRet := func(fileName []uint8) bool {
			// 获取可执行文件版本
			// fileVersion: 1.4.44.0
			ebp_334.v1 = 0x01
			ebp_334.v2 = 0x04
			ebp_334.v3 = 0x2C
			ebp_334.v4 = 0
			return true
		}(fileName[:]) // }(fileName[:], &ebp_334)

		if bRet { // je 0x004D,741B
			// _00DE817A
			func(verbuf []uint8, strfmt string, major uint16, minor uint16) int {
				verstr := fmt.Sprint(strfmt, major, minor)
				verbuf = []uint8(verstr) // 从堆上复制到rdata
				return len(verstr)
			}(_01319A38[:], "%d.%02d", ebp_334.v1, ebp_334.v2)

			if ebp_334.v3 > 0 { // jle 0x004D,7419
				*(*uint16)(unsafe.Pointer(&ebp_338)) = _0114DD64
				ebp_338.patch[2] = 0

				// 我的猜测是从patch=26开始使用字母和+标记标识patch，比如
				// 27就是A+
				// 28就是B+
				// 44就是R+
				if ebp_334.v3 > 0x1A { // jle 0x004D,73EE
					ebp_338.patch[0] = 'A' - 27 + uint8(ebp_334.v3) // 65-27+44=82='R'
					ebp_338.patch[1] = '+'
				}
				_00DE8010(_01319A38[:], string(ebp_338.patch[:]))
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
			nRet = func(x *uint32, y []uint8) int32 {
				return -1
			}(&_01319A50, partition.name[:])
			if nRet == 0 {
				// _00DECBD1
				nRet := func(x []uint8) uint32 {
					return 0
				}(partition.ip[:])
				if nRet == uint32(_012E233C) {
					_00DE8000(_01319DB8[:], partition.num[:])
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

	var ebp_4 uint32
}

// 0x004D,7CE5
func run() {
	//

	checkupdate()

	keyenc1 := new([54]uint8)            // 0x08C8,D050
	getenc1("Data/Enc1.dat", keyenc1[:]) // call 0x00B6,2CF0

	keydec2 := new([54]uint8)            // 0x08C8,D098
	getdec2("Data/Dec2.dat", keydec2[:]) // call 0x00B6,2D30

	// 似乎返回值并没有体现出来
	s := _00B3BAE4(&_01319E08, "> To read config.ini.\r\n")
	_004D7281()

	_004D6D64()
}
