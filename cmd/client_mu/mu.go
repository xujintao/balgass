package main

import "sync"

// winThread gui thread
type winThread interface {
	InitApplication() bool
	InitInstance() bool
}

// winApp base app
type winApp struct {
	// winThread
	// 注： c++里面winApp继承了winThread，而go的winApp实现了winThread
	// c++只能通过继承方式来实现某接口且不需要构造，而go这边继承和实现做了区分，含义不同
	hInstance uintptr
	cmdLine   string
}

// InitApplication hooks for your initiatization code
// f90
func (app *winApp) InitApplication() bool {
	return true
}

// f50
func (app *winApp) InitInstance() bool {
	return true
}

var theApp muApp

// muApp the app, v004632B8
type muApp struct {
	winApp
}

// f50
func (app *muApp) InitInstance() bool {
	// CreateMutexA(NULL, TRUE, "MU AutoUpdate")
	// v0046F448.LoadWTF("message.wtf")
	// if 0 == f00403EE0().f00403C20() { // v00453224.f00403C20()
	// f0041E6D3("设备未初始化 请安装最新的图形驱动程序", 10, 0) // EUC-KR  --> translate.google.com
	// return
	// }

	// 注册表
	f00401870once().f004015D0Get() // 单例v00463180.f004015D0()
	// dd := struct {
	// 	displayDeviceName [128]uint8
	// }{}
	// f00433360memset(dd.buf[:], 0, 0x1A8)
	// EnumDisplayDevicesA(0, 0, &dd, 0)
	// if 0 != f00433461memcmp(f00401870once().f004017F0(), dd.displayDeviceName[:], 0x80) {
	// 	// not equal
	// }
	// if ERROR_ALREADY_EXISTS == GetLastError() {
	// 	return
	// }

	// f004223A2mfcAfxEnableControlContainer()
	if app.cmdLine != "" {
		// ...
	}
	v004632B4port = 44405
	v0014AC88dlg := muDlg{}
	v0014AC88dlg.f0040AEC0(0)
	// app.m20mainWnd = &v0014AC88dlg
	v0014AC88dlg.DoModal() // f00416045

	return true
}

type reg struct {
	UserID                  [11]uint8  // 0x00
	DisplayDeviceName       [129]uint8 // 0x0B
	DisplayDeviceModelIndex int        // 0x8C
	FullScreenMode          int        // 0x90
	DisplayColorBit         int        // 0x94
	MusicOn                 int        // 0x98
	SoundOn                 int        // 0x9C
}

// reg
var regOnce sync.Once
var v00463180Reg *reg

func f00401870once() *reg {
	regOnce.Do(func() {
		v00463180Reg = new(reg)
	})
	return v00463180Reg
}

func (r *reg) f004014B0Set() bool {
	// hReg := 0
	// dwDisposition := 0
	// if ERROR_SUCCESS != RegCreateKeyEx(HKEY_CURRENT_USER, // hKey
	// 	"SOFTWARE\\Webzen\\Mu\\Config", // lpSubKey
	// 	0,                              // Reserved
	// 	0,                              // lpClass
	// 	0,                              // dwOptions
	// 	KEY_ALL_ACCESS,                 // samDesired 0xF003F
	// 	0,                              // lpSecurityAttributes
	// 	&hReg,                          // phkResult
	// 	&dwDisposition,                 // lpdwDisposition
	// ) {
	// 	return false
	// }

	// RegSetValueEx(hReg, "UserID", 0, 1, r.UserID[:], len(r.UserID))
	// RegSetValueEx(hReg, "DisplayDeviceName", 0, 1, r.DisplayDeviceName[:], len(r.DisplayDeviceName))
	// RegSetValueEx(hReg, "DisplayDeviceModelIndex", 0, 4, &r.DisplayDeviceModelIndex, 4)
	// RegSetValueEx(hReg, "FullScreenMode", 0, 4, &r.FullScreenMode, 4)
	// RegSetValueEx(hReg, "DisplayColorBit", 0, 4, &r.DisplayColorBit, 4)
	// RegSetValueEx(hReg, "MusicOn", 0, 4, &r.MusicOn, 4)
	// RegSetValueEx(hReg, "SoundOn", 0, 4, &r.SoundOn, 4)
	// RegCloseKey(hReg)
	return true
}

func (r *reg) f004015D0Get() bool {
	// hReg := 0
	// dwDisposition := 0
	// if ERROR_SUCCESS != RegCreateKeyEx(HKEY_CURRENT_USER, // hKey
	// 	"SOFTWARE\\Webzen\\Mu\\Config", // lpSubKey
	// 	0,                              // Reserved
	// 	0,                              // lpClass
	// 	0,                              // dwOptions
	// 	KEY_ALL_ACCESS,                 // samDesired 0xF003F
	// 	0,                              // lpSecurityAttributes
	// 	&hReg,                          // phkResult
	// 	&dwDisposition,                 // lpdwDisposition
	// ) {
	// 	return false
	// }
	// size := 10
	// RegQueryValueEx(hReg, "UserID", 0, 0, r.UserID[:], &size)

	// size = 128
	// RegQueryValueEx(hReg, "DisplayDeviceName", 0, 0, r.DisplayDeviceName[:], &size)

	// size = 4
	// RegQueryValueEx(hReg, "DisplayDeviceModelIndex", 0, 0, &r.DisplayDeviceModelIndex, &size)
	// RegQueryValueEx(hReg, "FullScreenMode", 0, 0, &r.FullScreenMode, &size)
	// RegQueryValueEx(hReg, "DisplayColorBit", 0, 0, &r.DisplayColorBit, &size)
	// RegQueryValueEx(hReg, "MusicOn", 0, 0, &r.MusicOn, &size)
	// RegQueryValueEx(hReg, "SoundOn", 0, 0, &r.SoundOn, &size)
	// RegCloseKey(hReg)
	// r.f004014B0Set()
	return true
}

func (r *reg) f004017F0(x int) []uint8 {
	return r.DisplayDeviceName[:]
}
