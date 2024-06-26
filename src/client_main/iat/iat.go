package iat

// advapi32.dll
// 这里定义所有用到的API指针变量，36个
// 最后一个DWORD为0
type advapi32 interface {
	CryptGetUserKey()
	RegCLoseKey()
	RegSetValueExA()
	// ...
	RegQueryValueExW()
}

// crypt32.dll
// 这里定义所有用到的API指针变量，33个
// 最后一个DWORD为0
type crypt32 interface {
	CertFreeCertificateChainEngine()
	CertFreeCertificateChain()
	CertGetIntendedKeyUsage()
	// ...
	CertCreateCertificateChainEngine()
}

// dsound.dll, /c/windows/syswow64/dsound.dll
// 这里定义所有用到的API指针变量，2个
// 最后一个DWORD为0
type dsound interface {
	DirectSoundEnumerateA()
	DirectSoundCreate()
}

// gdi32.dll, /c/windows/syswow64/gdi32.dll
// 这里定义所有用到的API指针变量，20个
// 最后一个DWORD为0
type gdi32 interface {
	CreateCompatibleDC()
	CreateDIBSection()
	DeleteDC()
	// ...
	GetTextMetricsW()
	CreateFontA()
}

// glu32.dll, /c/windows/syswow64/glu32.dll
// 这里定义所有用到的API指针变量，3个
// 最后一个DWORD为0
type glu32 interface {
	gluPerspective()
	gluLookAt()
	gluOrtho2D()
}

// imm32.dll, /c/windows/syswow64/imm32.dll
// 这里定义所有用到的API指针变量，22个
// 最后一个DWORD为0
type imm32 interface {
	ImmSetConversionStatus()
	ImmGetCompositionWindow()
	ImmGetCompositionStringA()
	ImmIsUIMessageA()
	// ...
	ImmReleaseContext()
}

// iphlpapi.dll, /c/windows/syswow64/iphlpapi.dll
// 这里定义所有用到的API指针变量，1个
// 最后一个DWORD为0
type iphlpapi interface {
	GetAdaptersInfo()
}

// kernel32.dll, /c/windows/syswow64/kernel32.dll
// 这里定义所有用到的API指针变量，219个
// 最后一个DWORD为0
// IAT多了第一个VirtualQueryEx
type kernel32 interface {
	GlobalLock()
	CreateFileA()
	GetCommandLineA()
	CloseHandle()
	// ...
	GetTickCount()
}

// oleaut32.dll, /c/windows/syswow64/oleaut32.dll
// 这里定义所有用到的API指针变量，6个
// 最后一个DWORD为0
type oleaut32 interface {
	SysStringLen()
	VariantInit()
	VariantClear()
	SysAllocStringByteLen()
	SysAllocString()
	SysFreeString()
}

// opengl32.dll, /c/windows/syswow64/opengl32.dll
// 这里定义所有用到的API指针变量，75个
// 最后一个DWORD为0
type opengl32 interface {
	glScissor()
	glStencilMask()
	glTexCoordPointer()
	// ...
	wglGetProcAddress()
}

// shell32.dll, /c/windows/syswow64/shell32.dll
// 这里定义所有用到的API指针变量，1个
// 最后一个DWORD为0
type shell32 interface {
	ShellExecuteA()
}

// sh1wapi.dll, /c/windows/syswow64/sh1wapi.dll
// 这里定义所有用到的API指针变量，1个
// 最后一个DWORD为0
type sh1wapi interface {
	PathFileExistsA()
}

// user32.dll, /c/windows/syswow64/user32.dll
// 这里定义所有用到的API指针变量，100个
// 最后一个DWORD为0
type user32 interface {
	FindWindowA()
	IntersetRect()
	wsprintfA()
	// ...
	SetWindowLongW()
	ReleaseDC()
}

// version.dll, /c/windows/syswow64/version.dll
// 这里定义所有用到的API指针变量，3个
// 最后一个DWORD为0
type version interface {
	VerQueryValueA()
	GetFileVersionInfoSizeA()
	GetFileVersionInfoA()
}

// wininet.dll, /c/windows/syswow64/wininet.dll
// 这里定义所有用到的API指针变量，15个
// 最后一个DWORD为0
type wininet interface {
	InternetOpenA()
	InternetConnectA()
	// ...
	InternetOpenUrlA()
	InternetConnectW()
}

// winmm.dll, /c/windows/syswow64/winmm.dll
// 这里定义所有用到的API指针变量，12个
// 最后一个DWORD为0
type winmm interface {
	timeGetDevCaps()
	timeBeginPeriod()
	timeKillEvent()
	// ...
	timeEndPeriod()
}

// ws2_32.dll, /c/windows/syswow64/ws2_32.dll
// 这里定义所有用到的API指针变量，35个
// 最后一个DWORD为0
type ws2_32 interface {
	ntohl()
	getnameinfo()
	gethostname()
	ioctlsocket()
	// ...
	WSAStartup()
}

// dbghelp.dll, /c/windows/syswow64/dbghelp.dll
// debghelp.dll, /c/mu/dbghelp.dll 这个和系统的有什么区别？
// 这里定义所有用到的API指针变量，7个
// 最后一个DWORD为0
type dbghelp interface {
	SymGetLineFromAddr64()
	SymFromAddr()
	StackWalk64()
	SymInitialize()
	SymSetOptions()
	MiniDumpWriteDump()
	SymCleanup()
}

// ole32.dll, /c/windows/syswow64/ole32.dll
// 这里定义所有用到的API指针变量，4个
// 最后一个DWORD为0
type ole32 interface {
	CoTaskMemFree()
	CoCreateInstance()
	CoUninitialize()
	CoInitialize()
}

// urlmon.dll, /c/windows/syswow64/urlmon.dll
// 这里定义所有用到的API指针变量，2个
// 最后一个DWORD为0
type urlmon interface {
	URLDownloadToFileA()
	URLDownloadToFileW()
}

// wzaudio.dll, /c/mu/wzAudio.dll
// 这里定义所有用到的API指针变量，6个
// 最后一个DWORD为0
type wzaudio interface {
	wzAudioPlay()
	wzAudioGetStreamOffsetRange()
	wzAudioDestroy()
	wzAudioOption()
	wzAudioCreate()
	wzAudioStop()
}
