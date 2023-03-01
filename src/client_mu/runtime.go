package main

import (
	"os"
	"strings"
	"unsafe"

	"github.com/xujintao/balgass/win"
)

func f00414B5Cnew(size int) unsafe.Pointer {
	return nil
}

func f00433360memset(dst []uint8, c uint8, size int) {

}

func f00433461memcmp(dst []uint8, src []uint8, size int) int {
	return 0
}

func f00417837AfxGetApp() *muApp {
	return &theApp
}

func f00449EEAmfcAfxWinMain(hInstance win.HINSTANCE, hPrevInstance win.HINSTANCE, lpCmdLine string, nCmdShow int) int {
	retCode := -1

	// f004209C7mfcAfxGetThread
	thread := func() winThread {
		// 	// check for current thread in module thread state
		// 	AFX_MODULE_THREAD_STATE* pState = AfxGetModuleThreadState();
		// 	CWinThread* pThread = pState->m_pCurrentWinThread;
		// // if no CWinThread for the module, then use the global app
		// if (pThread == NULL)
		// 	pThread = AfxGetApp();
		// 	return pThread;
		return &theApp
	}()

	app := f00417837AfxGetApp()

	// f0042778AmfcAfxWinInit(hInstance, hPrevInstance, lpCmdLine, nCmdShow)

	if app != nil && !app.InitApplication() { // f90()
		goto InitFailure
	}

	if !thread.InitInstance() { // f50()
		// if thread.mainWnd != 0 {
		// 	thread.OnIdel() // f60
		// }
		// thread.ExitInstance() // f68
		goto InitFailure
	}
	// thread.Run()             // f54

InitFailure:
	// f00426598mfcAfxWinTerm()
	return retCode
}

// f00434EB8, runtime_main, __winMainCRTStartup
func main() {
	// ...
	// 0x00434F22: f00442369heapInit(true)
	// 0x00434F38: f0043C60FflsInit()

	// ...
	szcmdLine := strings.Join(os.Args[1:], " ")
	show := 10
	// 0x00434FBB:
	f00449EEAmfcAfxWinMain(0x00400000, 0, szcmdLine, show)
}

// oep
func f00435036winMainCRTStartup() {
	// f00442399securityInitCookie()
	// jmp f00434EB8
}
