package main

import (
	"os"
	"strings"

	"github.com/xujintao/balgass/win"
)

func f00449EEAmfcAfxWinMain(hInstance win.HINSTANCE, hPrevInstance win.HINSTANCE, szCmdLine string, iCmdShow int) int {
	retCode := -1

	// f004209C7mfcAfxGetThread
	v004632B8app := func() winThread {
		// 	// check for current thread in module thread state
		// 	AFX_MODULE_THREAD_STATE* pState = AfxGetModuleThreadState();
		// 	CWinThread* pThread = pState->m_pCurrentWinThread;
		// // if no CWinThread for the module, then use the global app
		// if (pThread == NULL)
		// 	pThread = AfxGetApp();
		// 	return pThread;
		return &theApp
	}()
	// f00417837mfcAfxGetModuleState()
	// f0042778AmfcAfxWinInit()
	v004632B8app.InitApplication()    // f90()
	if !v004632B8app.InitInstance() { // f50()
		// if v004632B8app.mainWnd != 0 {
		// 	v004632B8app.OnIdel() // f60
		// }
		// v004632B8app.ExitInstance() // f68
	} else {
		// v004632B8app.Run()             // f54
	}
	// f00426598mfcAfxWinTerm()
	return retCode
}

// f00434EB8, runtime_main, __winMainCRTStartup
func main() {
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
