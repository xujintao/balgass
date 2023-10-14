package win

import (
	"syscall"
	"unsafe"
)

const socket_error = uintptr(^uint32(0))

var (
	libws2_32 = syscall.NewLazyDLL("ws2_32.dll")

	procWSAStartup        = libws2_32.NewProc("WSAStartup")
	procWSACleanup        = libws2_32.NewProc("WSACleanup")
	procWSAIoctl          = libws2_32.NewProc("WSAIoctl")
	procsocket            = libws2_32.NewProc("socket")
	procsetsockopt        = libws2_32.NewProc("setsockopt")
	procgetsockopt        = libws2_32.NewProc("getsockopt")
	procbind              = libws2_32.NewProc("bind")
	procconnect           = libws2_32.NewProc("connect")
	procgetsockname       = libws2_32.NewProc("getsockname")
	procgetpeername       = libws2_32.NewProc("getpeername")
	proclisten            = libws2_32.NewProc("listen")
	procshutdown          = libws2_32.NewProc("shutdown")
	procclosesocket       = libws2_32.NewProc("closesocket")
	procWSARecv           = libws2_32.NewProc("WSARecv")
	procWSASend           = libws2_32.NewProc("WSASend")
	procWSARecvFrom       = libws2_32.NewProc("WSARecvFrom")
	procWSASendTo         = libws2_32.NewProc("WSASendTo")
	procgethostbyname     = libws2_32.NewProc("gethostbyname")
	procgetservbyname     = libws2_32.NewProc("getservbyname")
	procntohs             = libws2_32.NewProc("ntohs")
	procgetprotobyname    = libws2_32.NewProc("getprotobyname")
	procGetAddrInfoW      = libws2_32.NewProc("GetAddrInfoW")
	procFreeAddrInfoW     = libws2_32.NewProc("FreeAddrInfoW")
	procWSAEnumProtocolsW = libws2_32.NewProc("WSAEnumProtocolsW")
)

func WSAStartup(verreq uint32, data *syscall.WSAData) (sockerr error) {
	r0, _, _ := syscall.Syscall(procWSAStartup.Addr(), 2, uintptr(verreq), uintptr(unsafe.Pointer(data)), 0)
	if r0 != 0 {
		sockerr = syscall.Errno(r0)
	}
	return
}

func Socket(af int32, typ int32, protocol int32) (handle syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procsocket.Addr(), 3, uintptr(af), uintptr(typ), uintptr(protocol))
	handle = syscall.Handle(r0)
	if handle == syscall.InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WSARecv(s syscall.Handle, bufs *syscall.WSABuf, bufcnt uint32, recvd *uint32, flags *uint32, overlapped *syscall.Overlapped, croutine *byte) (err error) {
	r1, _, e1 := syscall.Syscall9(procWSARecv.Addr(), 7, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(recvd)), uintptr(unsafe.Pointer(flags)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0, 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Recv(s syscall.Handle, buf []uint8, size int, flags int) (int, error) {
	return 0, nil
}

func WSASend(s syscall.Handle, bufs *syscall.WSABuf, bufcnt uint32, sent *uint32, flags uint32, overlapped *syscall.Overlapped, croutine *byte) (err error) {
	r1, _, e1 := syscall.Syscall9(procWSASend.Addr(), 7, uintptr(s), uintptr(unsafe.Pointer(bufs)), uintptr(bufcnt), uintptr(unsafe.Pointer(sent)), uintptr(flags), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0, 0)
	if r1 == socket_error {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Send(s syscall.Handle, buf []uint8, size int, flags int) (int, error) {
	return 0, nil
}

func WSAGetLastError() uint32 {
	return GetLastError()
}
