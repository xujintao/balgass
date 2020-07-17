package main

import (
	"syscall"

	"github.com/xujintao/balgass/win"
)

var fs = &tib{}

// the Win32 Thread Information Block
// https://en.wikipedia.org/wiki/Win32_Thread_Information_Block
type tib struct {
	m00seh         uintptr // seh
	m04stackHi     uintptr
	m08stackLo     uintptr
	m18            *tib
	m20pid         uintptr
	m24tid         uintptr
	m2Ctls         uintptr
	m30peb         uintptr
	m34lastErr     uint
	mC0fastSysCall uintptr
	mFCA           uint8
	mFDCoffset     int
}

type CB func(hWnd win.HWND, message uint32, wParam, lParam uintptr) uintptr

// user32
func f77164BE0(cb CB, unk1 int, msg uint32, wParam, lParam uintptr, unk2, unk3 int) uintptr {
	// ebp4 := 0
	// ...
	var ebp1C win.HWND
	var ebxOrigin CB = cb
	// f77184360
	ebp20 := func(cb CB, hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
		// fs.mFCA |= 1
		return cb(hWnd, msg, wParam, lParam)
		// fs.mFCA &= 0xFE
	}(ebxOrigin, ebp1C, msg, wParam, lParam)
	// esi
	// ebp28
	// ebp4 = -2
	// f77165208()
	return ebp20
}

func CallWindowProcW(cbp uintptr, hWnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	if cbp == 0 {
		return 0
	}
	// convert cbp
	var cb CB
	// get tls
	var ecx = 0
	return f77164BE0(cb, 0, msg, wParam, lParam, ecx, 0)
}

// opengl32
func f68E8F770wndProcGL(hWnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	// ...
	ebp4 := syscall.NewCallback(f004D6C2BmainWndProc)
	// ...
	ebp8 := lParam
	return CallWindowProcW(ebp4, hWnd, msg, wParam, ebp8)
}

// ntdll.f776B8370wndProcIME
// user32.f77161D10wndProcIME
func f77161D10wndProcIME(hWnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	// user32.f77163305tls(hWnd)
	eax, edx := func(hWnd win.HWND) (uintptr, uintptr) { return 0x0BFCA540 /*[fs.m18.mF70+0x800+0x48]*/, 0 }(hWnd)
	if eax|edx == 0 {
		return 0
	}
	// user32.f77161D40
	return func(pWnd uintptr, unk2 uintptr, msg uint32, wParam uintptr, lParam uintptr, unk3 uintptr) uintptr {
		// Input Method Manager (IMM) is a technology used by an application to communicate with
		// an input method editor (IME), which runs as a service.
		// The IME allows computer users to enter complex characters and symbols,
		// such as Japanese kanji characters, by using a standard keyboard.
		// ...
		// if ret := imm32.CtfImmDispatchDefImeMessage(pWnd.m00, msg, wParam, lParam); ret != 0 {
		// 	return ret
		// }
		// ...
		switch msg {
		case 0x287:
			// user32.f771610FC
			return func(pWnd uintptr, msg uint32, wParam uintptr, lParam uintptr) uintptr {
				// ...
				switch wParam {
				case 0x18:
					// win10.user32.f76480C2C() win7.user32.f76BFD9DB()
					func(hWnd win.HWND, flag uintptr) {
						// if user32.IsWindow(hWnd) == false {
						// 	return
						// }
						// immCtx := imm32.ImmGetContext(hWnd)
						// imm32.ImmSetActiveContext(hWnd, immCtx, flag)
						func() {
							// ...
							// msctf.CtfImeAssociateFocus(immCtx, flag, hWnd)
							func() {
								// tls33 := kernel32.TlsGetValue(0x33)
								// msctf.f7503F4D0.win10(tls33, immCtx, flag, hWnd, 0, 0) // ecx = tls33.m00
								func() {
									if flag == 0 {
										// win10.0x7503F650:
										// msctf.f75042DB0.win10(tls33, hWnd, immCtx, tls33.m08, 0, 0, 0) // ecx = tls33.m00
										func() {
											// 很复杂
											// ...
											// msctf.f7503F820.win10(0, 1) // before GetFocus，正确的逻辑是跳过这一步
											func() {
												// msctf.f7505E605.win10(2, 0, 0x0B57A728) 带ecx:tls33.m08
												func() {
													// ...
													// msctf.do6win10f750BE300(0x129A5248, 0, 0x0B57A728)
													func() {
														// ...
														// msctf.do4win10f74EEF5F0(0x129498C4, 0x0019D750, 0x0019D714, 0x129A532C)
														func() {
															// ...
															index := 0
															for {
																// obj := unk1.m4F0 + unk1.m4F8*index
																// obj.do4()
																// t0119F638.do4() is f00D86830
																func() {
																	// ...
																	// f00D8692F:
																	eax := func() int { return 0 }()
																	eax--
																	if eax == 0 {
																		// v09D9BD74mm.do13(edi)
																	}
																	// 0x00D86943:
																	// eax = [esi+4] // eax = esi.m04, eax = 0x0F720D00 esi=0x0F720EF4
																	// ecx = [eax+4] // ecx = eax.m04, ecx=0
																	// edx = ecx.m08 // 非法内存访问
																}()
																index++
																// if index >= unk1.m4F4 {
																// 	break
																// }
															}
														}()
													}()
												}()
											}()
										}()
									}
								}()
							}()
							// ...
							win.SendMessage(hWnd, 0x281, flag, 0xC000000F /*ISC_SHOWUIALL*/) // 落到ntdll.kiUserCallbackDispatcher(2, &msg)
						}()
						// imm32.ImmReleaseContext(hWnd, immCtx)
					}(win.HWND(lParam), 0)
				}
				// ...
				return 0
			}(pWnd /*pWnd.m128*/, 0x287, wParam, lParam)
		}
		return 0
	}(eax, edx, msg, wParam, lParam, 0)
}

// ntdll.f776B8430handle
// user32.f77164A40handle
func f77164A40handle(pWnd uintptr, unk2 int, msg uint32, wParam uintptr, lParam int, cb uintptr) {
	// ecxtib := 0
	// ebp1C := f77164BE0(CB(unsafe.Pointer(cb)), 0, msg, wParam, lParam, ecxtib, 1)
	// tls
}

type message struct {
	m00pWnd   uintptr
	m04unk2   int
	m08msg    int // 0x8/WM_KILLFOCUS, 0x287/WM_IME_SYSTEM
	m0CwParam uintptr
	m10lParam uintptr
	m14cb     int
	m18handle uintptr
}

// user32
func f7716E0A0handle(msg *message) {
	// msg.m18handle(msg.m00pWnd, msg.m04unk2, msg.m08msg, msg.m0CwParam, msg.m10lParam, msg.m14cb)
	// 8/WM_KILLFOCUS, hWndEdit, 0
	// msg: f776B83A0handle(0x0C0749F0, 0, 0x008, hWndEdit, 0, vboxgl.f685B28D0) // empty, client_vc has none
	// msg: f776B8430handle(0x0C0749F0, 0, 0x008, hWndEdit, 0, f68E8F770wndProcGL)
	// msg: f776B83D0handle(0x0C0749F0, 0, 0x008, hWndEdit, 0, apphelp.f6E423910) // empty

	// 287/WM_IME_SYSTEM, 0x18/*WM_IME_SETCONTEXT*/, hWndMain
	// msg: f776B83A0handle(0x0C06BD40, 0, 0x287, 0x18, hWndMain, vboxgl.f685B28D0) // empty, client_vc has none
	// msg: f776B8430handle(0x0C06BD40, 0, 0x287, 0x18, hWndMain, f776B8370wndProcIME)

	// var ebp18 [6]int
	// user32.NtCallbackReturn(&ebp18, 0x16, 0) // return to f7716E0A0handle with a new msg
}

// ntdll
func kiUserCallbackDispatcher(index int, msg *message) {
	// seh
	// fs.m30peb.m2CkernelCallbackTable[index](msg) // user32.f7716E0A0handle
	// ntdll.NtCallbackReturn()
	// ntdll.RtlRaiseStatus()
}

// user32
func setFocus(hWnd win.HWND) {
	// fs.mC0fastSysCall() // 落到ntdll.kiUserCallbackDispatcher(2, &msg)
}
