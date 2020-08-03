package main

import (
	"github.com/xujintao/balgass/win"
)

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

// f00405AB0
func (app *muApp) InitInstance() bool {
	// 很大的栈

	// CreateMutexA(NULL, TRUE, "MU AutoUpdate")
	v0046F448msg.LoadWTF("message.wtf")
	// if 0 == f00403EE0().f00403C20() { // v00453224.f00403C20()
	// 	f0041E6D3("设备未初始化 请安装最新的图形驱动程序", 10, 0) // EUC-KR  --> translate.google.com
	// 	return
	// }

	// 0x00405B31: 注册表
	f00401870reg().f004015D0Get() // 单例v00463180.f004015D0()
	// dd := struct {
	// 	displayDeviceName [128]uint8
	// }{}
	// f00433360memset(dd.buf[:], 0, 0x1A8)
	// EnumDisplayDevicesA(0, 0, &dd, 0)
	// ...
	// if 0 != f00433461memcmp(f00401870reg().f004017F0(), dd.displayDeviceName[:], 0x80) {
	// 	// 0x00405BCD:
	// 	f0041E6D3("分辨率已重置。请按“选项设置”按钮以指定新的分辨率。", 10, 0) // https://r12a.github.io/app-encodings/
	// }

	// 0x00405BDA:
	if 183 /*ERROR_ALREADY_EXISTS*/ == win.GetLastError() {
		return false
	}

	// f004223A2mfcAfxEnableControlContainer()
	if app.cmdLine != "" {
		// ...
	}
	v004632B4port = 44405
	v0014AC88dlg := muDlg{}
	v0014AC88dlg.f0040AEC0construct(0)
	// app.m20mainWnd = &v0014AC88dlg
	v0014AC88dlg.DoModal() // f00416045

	return true
}
