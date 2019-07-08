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

// 0x004D,7CE5
func run() {
	checkupdate()

	var pathenc1 = "Data/Enc1.dat"
	keyenc1 := new([54]uint8)     // 0x08C8,D050
	getenc1(pathenc1, keyenc1[:]) // call 0x00B6,2CF0

	var pathdec2 = "Data/Dec2.dat"
	keydec2 := new([54]uint8)     // 0x08C8,D098
	getdec2(pathdec2, keydec2[:]) // call 0x00B6,2D30

	var pathconf = "> To read config.ini.\r\n"
	bufconf := new([20]uint8)      // 0x0131,9E08
	readconf(pathconf, bufconf[:]) // call 0x00B3,BAE4

	_004D6D64()
}
