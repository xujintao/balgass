package main

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

	_004D6D64()
}
