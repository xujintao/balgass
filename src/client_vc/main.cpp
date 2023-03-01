#include <windows.h>
#include <imm.h>
#pragma comment(lib, "imm32.lib")

WNDPROC wndProcEditOrigin;

int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpszCmdLine, int nCmdShow) {
	// main window
	WNDCLASSEXA mainWndClass = {
		sizeof(WNDCLASSEX),	// cbSize
		0,					// Style
		0,					// wndProc
		0,					// cbClsExt
		0,					// cbWndEXT
		hInstance,			// hInstance
		LoadIcon((HINSTANCE)NULL, IDI_APPLICATION),	// hIcon
		LoadCursor((HINSTANCE)NULL, IDC_ARROW),		// hCursor
		HBRUSH(GetStockObject(GRAY_BRUSH)),		// hBrushBackgroud
		0,					// menuName
		"main",				// className
		0					// hIconSmall
	};
	mainWndClass.lpfnWndProc = [](HWND hWnd, UINT msg, WPARAM wParam, LPARAM lParam)->LRESULT {
		switch (msg) {
		case WM_CLOSE:
			break;
		case WM_DESTROY:
			PostQuitMessage(0);
			break;
		case WM_NCDESTROY:
			break;
		}
		return DefWindowProc(hWnd, msg, wParam, lParam);
	};
	RegisterClassEx(&mainWndClass);
	HWND hWndMain = CreateWindowEx(
		0,				// ExStyle
		"main",			// ClassName
		"main",			// WindowName
		WS_SYSMENU,		// Style
		CW_USEDEFAULT,	// x
		CW_USEDEFAULT,	// y
		800,			// width
		600,			// height
		(HWND)NULL,		// parent window
		(HMENU)NULL,	// menu
		hInstance,		// Instance
		(LPVOID)NULL	// lpParam
	);
	ShowWindow(hWndMain, nCmdShow);
	UpdateWindow(hWndMain);
	
	// child edit window
	HWND hWndEdit =  CreateWindowEx(
		0,			// ExStyle
		"edit",		// ClassName
		"",			// WindowName
		ES_AUTOHSCROLL | WS_VISIBLE | WS_CHILD,	// Style
		0,			// x
		0,			// y
		100,		// width
		20,			// height
		hWndMain,	// parent window
		HMENU(1),	// menu
		hInstance,	// Instance
		nullptr		// ?
	);
	wndProcEditOrigin = (WNDPROC)SetWindowLongA(
		hWndEdit, 
		GWL_WNDPROC, 
		(LONG)(WNDPROC)[](HWND hWnd, UINT msg, WPARAM wParam, LPARAM lParam)->LRESULT { // wndProcEdit
			switch (msg) {}
			return CallWindowProc(wndProcEditOrigin, hWnd, msg, wParam, lParam);
		});
	ShowWindow(hWndEdit, SW_SHOWDEFAULT); // SW_HIDE
	SetFocus(hWndEdit);

	// msg loop
	MSG msg;
	ZeroMemory(&msg, sizeof(MSG));
	DWORD tick = GetTickCount();
	while(1){
		if (PeekMessage(&msg, 0, 0, 0, PM_NOREMOVE)) {
			BOOL bUIMsg = ImmIsUIMessage(NULL, msg.message, msg.wParam, msg.lParam);
			if (bUIMsg // Returns a nonzero value if the message is processed by the IME window
			|| msg.message == WM_KEYDOWN
			|| msg.message == WM_KEYUP
			|| msg.message == WM_LBUTTONDOWN
			|| msg.message == WM_LBUTTONUP) {
				
			}
			BOOL bRet = GetMessage(&msg, 0, 0, 0);
			if (bRet == FALSE) {
				break; // WM_QUIT
			}
			else if (bRet == -1) {
				break; // error
			}
			TranslateMessage(&msg);
			DispatchMessage(&msg);
		}
		else {
			// do others
			DWORD tickLast = tick;
			tick = GetTickCount();
			DWORD delta = tick - tickLast;
			if (delta < 40) {
				DWORD wait = 40 - delta;
				Sleep(wait);
				delta = 40;
			}
			for (int i = 0; i < 256; i++) {
				GetAsyncKeyState(i);
			}
		}
		// realtime handle
	}
}