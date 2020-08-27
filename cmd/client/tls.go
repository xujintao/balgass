package main

// fls
var v09D9DB84flsAlloc uintptr
var v09D9DB88flsGetValue uintptr
var v09D9DB8CflsSetValue uintptr
var v09D9DB90flsFree uintptr

var v012F7DBCfls int
var v012F7DC0tls int

func f00DFC0E5flsEnc(fls uintptr) uintptr {
	var v3 interface{} // v3 := kernel32.TlsGetValue(v012F7DC0tls)
	var v4 interface{} // v4 := kernel32.TlsGetValue(v012F7DC0tls).(v012F7DBCfls)
	var encPointer func(uintptr) uintptr
	if v3 != nil && v012F7DBCfls != -1 && v4 != nil {
		// encPointer = v3.m01F8
	} else {
		// encPointer = kernel32.GetProcAddress(kernel32.GetModuleHandle("kernel32.dll"), "EncodePointer")
	}
	if encPointer != nil {
		fls = encPointer(fls)
	}
	return fls
}

func f00DFC160flsDec(fls uintptr) uintptr {
	var v3 interface{} // v3 := kernel32.TlsGetValue(v012F7DC0tls)
	var v4 interface{} // v4 := kernel32.TlsGetValue(v012F7DC0tls).(v012F7DBCfls)
	var decPointer func(uintptr) uintptr
	if v3 != nil && v012F7DBCfls != -1 && v4 != nil {
		// decPointer = v3.m01FC
	} else {
		// decPointer = kernel32.GetProcAddress(kernel32.GetModuleHandle("kernel32.dll"), "DecodePointer")
	}
	if decPointer != nil {
		fls = decPointer(fls)
	}
	return fls
}

func f00DFC5ACflsInit() bool {
	// fls初始化
	// v09D9DB84flsAlloc = kernel32.GetProcAddress(kernel32.GetModuleHanlde("kernel32.dll"), "FlsAlloc")
	// v09D9DB88flsGetValue = kernel32.GetProcAddress(kernel32.GetModuleHanlde("kernel32.dll"), "FlsGetValue")
	// v09D9DB8CflsSetValue = kernel32.GetProcAddress(kernel32.GetModuleHanlde("kernel32.dll"), "FlsSetValue")
	// v09D9DB90flsFree = kernel32.GetProcAddress(kernel32.GetModuleHanlde("kernel32.dll"), "FlsFree")

	// use fls
	// v012F7DC0tls = kernel32.TlsAlloc()
	// kernel32.TlsSetValue(v012F7DC0tls, v09D9DB88flsGetValue)

	// 加密fls函数地址
	v09D9DB84flsAlloc = f00DFC0E5flsEnc(v09D9DB84flsAlloc)
	v09D9DB88flsGetValue = f00DFC0E5flsEnc(v09D9DB88flsGetValue)
	v09D9DB8CflsSetValue = f00DFC0E5flsEnc(v09D9DB8CflsSetValue)
	v09D9DB90flsFree = f00DFC0E5flsEnc(v09D9DB90flsFree)

	// v012F7DBCfls = f00DFC160flsDec(v09D9DB84flsAlloc).(f00DFC403)
	return true
}
