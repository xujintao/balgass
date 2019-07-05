package dll

// advapi32.dll
// 这里定义所有用到的API指针变量，36个
// 最后一个DWORD为0
type advapi32 struct {
	CryptGetUserKey func()
	RegCLoseKey     func()
	RegSetValueExA  func()
	// ...
	RegQueryValueExW func()
	end              uintptr
}

// crypt32.dll
// 这里定义所有用到的API指针变量，33个
// 最后一个DWORD为0
type crypt32 struct {
	CertFreeCertificateChainEngine func()
	CertFreeCertificateChain       func()
	CertGetIntendedKeyUsage        func()
	// ...
	CertCreateCertificateChainEngine func()
	end                              uintptr
}

// dsound.dll, /c/windows/syswow64/dsound.dll
// 这里定义所有用到的API指针变量，2个
// 最后一个DWORD为0
type dsound struct {
	DirectSoundEnumerateA func()
	DirectSoundCreate     func()
	end                   uintptr
}

// gdi32.dll, /c/windows/syswow64/gdi32.dll
// 这里定义所有用到的API指针变量，24个
// 最后一个DWORD为0
type gdi32 struct {
	CreateCompatibleDC func()
	CreateDIBSection   func()
	DeleteDC           func()
	// ...
	GetTextMetricsW func()
	CreateFontA     func()
	end             uintptr
}

// glu32.dll, /c/windows/syswow64/glu32.dll
// 这里定义所有用到的API指针变量，3个
// 最后一个DWORD为0
type glu32 struct {
	gluPerspective func()
	gluLookAt      func()
	gluOrtho2D     func()
	end            uintptr
}

// imm32.dll, /c/windows/syswow64/imm32.dll
// 这里定义所有用到的API指针变量，22个
// 最后一个DWORD为0
type imm32 struct {
	ImmSetConversionStatus   func()
	ImmGetCompositionWindow  func()
	ImmGetCompositionStringA func()
	ImmIsUIMessageA          func()
	// ...
	ImmReleaseContext func()
	end               uintptr
}

// iphlpapi.dll, /c/windows/syswow64/iphlpapi.dll
// 这里定义所有用到的API指针变量，1个
// 最后一个DWORD为0
type iphlpapi struct {
	GetAdaptersInfo func()
	end             uintptr
}

// kernel32.dll, /c/windows/syswow64/kernel32.dll
// 这里定义所有用到的API指针变量，219个
// 最后一个DWORD为0
type kernel32 struct {
	GlobalLock      func()
	CreateFileA     func()
	GetCommandLineA func()
	CloseHandle     func()
	// ...
	GetTickCOunt func()
	end          uintptr
}

// oleaut32.dll, /c/windows/syswow64/oleaut32.dll
// 这里定义所有用到的API指针变量，6个
// 最后一个DWORD为0
type oleaut32 struct {
	SysStringLen          func()
	VariantInit           func()
	VariantClear          func()
	SysAllocStringByteLen func()
	SysAllocString        func()
	SysFreeString         func()
	end                   uintptr
}

// opengl32.dll, /c/windows/syswow64/opengl32.dll
// 这里定义所有用到的API指针变量，75个
// 最后一个DWORD为0
type opengl32 struct {
	glScissor         func()
	glStencilMask     func()
	glTexCoordPointer func()
	// ...
	wglGetProcAddress func()
	end               uintptr
}

// shell32.dll, /c/windows/syswow64/shell32.dll
// 这里定义所有用到的API指针变量，1个
// 最后一个DWORD为0
type shell32 struct {
	ShellExecuteA func()
	end           uintptr
}

// sh1wapi.dll, /c/windows/syswow64/sh1wapi.dll
// 这里定义所有用到的API指针变量，1个
// 最后一个DWORD为0
type sh1wapi struct {
	PathFileExistsA func()
	end             uintptr
}

// user32.dll, /c/windows/syswow64/user32.dll
// 这里定义所有用到的API指针变量，100个
// 最后一个DWORD为0
type user32 struct {
	FindWindowA  func()
	IntersetRect func()
	wsprintfA    func()
	// ...
	SetWindowLongW func()
	ReleaseDC      func()
	end            uintptr
}

// version.dll, /c/windows/syswow64/version.dll
// 这里定义所有用到的API指针变量，3个
// 最后一个DWORD为0
type version struct {
	VerQueryValueA          func()
	GetFileVersionInfoSizeA func()
	GetFileVersionInfoA     func()
}

// wininet.dll, /c/windows/syswow64/wininet.dll
// 这里定义所有用到的API指针变量，15个
// 最后一个DWORD为0
type wininet struct {
	InternetOpenA    func()
	InternetConnectA func()
	// ...
	InternetOpenUrlA func()
	InternetConnectW func()
	end              uintptr
}

// winmm.dll, /c/windows/syswow64/winmm.dll
// 这里定义所有用到的API指针变量，12个
// 最后一个DWORD为0
type winmm struct {
	timeGetDevCaps  func()
	timeBeginPeriod func()
	timeKillEvent   func()
	// ...
	timeEndPeriod func()
	end           uintptr
}

// ws2_32.dll, /c/windows/syswow64/ws2_32.dll
// 这里定义所有用到的API指针变量，35个
// 最后一个DWORD为0
type ws2_32 struct {
	ntohl       func()
	getnameinfo func()
	gethostname func()
	ioctlsocket func()
	// ...
	WSAStartup func()
	end        uintptr
}

// dbghelp.dll, /c/windows/syswow64/dbghelp.dll
// debghelp.dll, /c/mu/dbghelp.dll 这个和系统的有什么区别？
// 这里定义所有用到的API指针变量，35个
// 最后一个DWORD为0
type dbghelp struct {
	SymGetLineFromAddr64 func()
	SymFromAddr          func()
	StackWalk64          func()
	SymInitialize        func()
	SymSetOptions        func()
	MiniDumpWriteDump    func()
	SymCleanup           func()
	end                  uintptr
}

// ole32.dll, /c/windows/syswow64/ole32.dll
// 这里定义所有用到的API指针变量，4个
// 最后一个DWORD为0
type ole32 struct {
	CoTaskMemFree    func()
	CoCreateInstance func()
	CoUninitialize   func()
	CoInitialize     func()
	end              uintptr
}

// urlmon.dll, /c/windows/syswow64/urlmon.dll
// 这里定义所有用到的API指针变量，2个
// 最后一个DWORD为0
type urlmon struct {
	URLDownloadToFileA func()
	URLDownloadToFileW func()
	end                uintptr
}

// wzaudio.dll, /c/mu/wzAudio.dll
// 这里定义所有用到的API指针变量，6个
// 最后一个DWORD为0
type wzaudio struct {
	wzAudioPlay                 func()
	wzAudioGetStreamOffsetRange func()
	wzAudioDestroy              func()
	wzAudioOption               func()
	wzAudioCreate               func()
	wzAudioStop                 func()
	end                         uintptr
}
