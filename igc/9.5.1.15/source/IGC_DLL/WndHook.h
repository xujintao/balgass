#ifndef WND_HOOK_H
#define WND_HOOK_H
//004D804D   C745 D4 FF6E4D00 MOV DWORD PTR SS:[EBP-2C],main_f.004D6EF>


extern WNDPROC muWndProc;
extern HWND MuWindow;
extern HWND LauncherWindowHandle;
extern bool gLauncherProxyEnabled;
extern bool gLauncherProxyReconnectInitiated;
LRESULT CALLBACK WndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam);
#endif