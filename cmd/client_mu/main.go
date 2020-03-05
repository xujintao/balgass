package main

import (
	"sync"
)

var v00461DF0stack int = 0xC58267AF
var v004632B4port uint16
var v0048F4A0 func()
var v0046327C int = 0

type reg struct {
	UserID                  [11]uint8  // 0x00
	DisplayDeviceName       [129]uint8 // 0x0B
	DisplayDeviceModelIndex int        // 0x8C
	FullScreenMode          int        // 0x90
	DisplayColorBit         int        // 0x94
	MusicOn                 int        // 0x98
	SoundOn                 int        // 0x9C
}

var regOnce sync.Once
var v00463180Reg *reg

func f00401870once() *reg {
	regOnce.Do(func() {
		v00463180Reg = new(reg)
	})
	return v00463180Reg
}

func (r *reg) f004014B0Set() bool {
	// hReg := 0
	// dwDisposition := 0
	// if ERROR_SUCCESS != RegCreateKeyEx(HKEY_CURRENT_USER, // hKey
	// 	"SOFTWARE\\Webzen\\Mu\\Config", // lpSubKey
	// 	0,                              // Reserved
	// 	0,                              // lpClass
	// 	0,                              // dwOptions
	// 	KEY_ALL_ACCESS,                 // samDesired 0xF003F
	// 	0,                              // lpSecurityAttributes
	// 	&hReg,                          // phkResult
	// 	&dwDisposition,                 // lpdwDisposition
	// ) {
	// 	return false
	// }

	// RegSetValueEx(hReg, "UserID", 0, 1, r.UserID[:], len(r.UserID))
	// RegSetValueEx(hReg, "DisplayDeviceName", 0, 1, r.DisplayDeviceName[:], len(r.DisplayDeviceName))
	// RegSetValueEx(hReg, "DisplayDeviceModelIndex", 0, 4, &r.DisplayDeviceModelIndex, 4)
	// RegSetValueEx(hReg, "FullScreenMode", 0, 4, &r.FullScreenMode, 4)
	// RegSetValueEx(hReg, "DisplayColorBit", 0, 4, &r.DisplayColorBit, 4)
	// RegSetValueEx(hReg, "MusicOn", 0, 4, &r.MusicOn, 4)
	// RegSetValueEx(hReg, "SoundOn", 0, 4, &r.SoundOn, 4)
	// RegCloseKey(hReg)
	return true
}

func (r *reg) f004015D0Get() bool {
	// hReg := 0
	// dwDisposition := 0
	// if ERROR_SUCCESS != RegCreateKeyEx(HKEY_CURRENT_USER, // hKey
	// 	"SOFTWARE\\Webzen\\Mu\\Config", // lpSubKey
	// 	0,                              // Reserved
	// 	0,                              // lpClass
	// 	0,                              // dwOptions
	// 	KEY_ALL_ACCESS,                 // samDesired 0xF003F
	// 	0,                              // lpSecurityAttributes
	// 	&hReg,                          // phkResult
	// 	&dwDisposition,                 // lpdwDisposition
	// ) {
	// 	return false
	// }
	// size := 10
	// RegQueryValueEx(hReg, "UserID", 0, 0, r.UserID[:], &size)

	// size = 128
	// RegQueryValueEx(hReg, "DisplayDeviceName", 0, 0, r.DisplayDeviceName[:], &size)

	// size = 4
	// RegQueryValueEx(hReg, "DisplayDeviceModelIndex", 0, 0, &r.DisplayDeviceModelIndex, &size)
	// RegQueryValueEx(hReg, "FullScreenMode", 0, 0, &r.FullScreenMode, &size)
	// RegQueryValueEx(hReg, "DisplayColorBit", 0, 0, &r.DisplayColorBit, &size)
	// RegQueryValueEx(hReg, "MusicOn", 0, 0, &r.MusicOn, &size)
	// RegQueryValueEx(hReg, "SoundOn", 0, 0, &r.SoundOn, &size)
	// RegCloseKey(hReg)
	// r.f004014B0Set()
	return true
}

func (r *reg) f004017F0(x int) []uint8 {
	return r.DisplayDeviceName[:]
}

type block struct {
	m00 uint
	m04 uint
	m08 uint
	m0C uint
	m10 uint
	m14 uint
	m18 uint
	m1C uint
	m20 uint
	m24 uint8
	m28 uint
	m2C uint
	m30 uint
	m34 uint
	m38 uint
	m3C uint
	m40 uint
	m44 uint
	m48 uint
	m4C uint
	m50 uint
	m54 uint
	m58 uint
	m68 uint
	m74 uint
	m78 uint
	m84 uint
	m88 uint
}

func (b *block) f00417837() uint {
	// v00A6DE08 := v0048F8C8.f0042566A(f0041547C)
	// if v00A6DE08 == 0 {
	// 	// 0x00415460
	// }
	// if v00A6DE08.f04 != 0 {
	// 	return 1
	// }
	var ret uint
	// ret = v0048F8C4.f00425146(f00417808)
	// if ret == 0 {
	// 	// 0x00415460
	// }
	return ret // 0x00A6DF18
}

func (b *block) f0041987A() {
	// b.f0041DE41()
	func() {
		// b.f00417837
		b.m1C = b.f00417837()
		b.m04 = 1
		b.m08 = 0
		b.m0C = 0
		b.m10 = 0
		b.m14 = 1
		b.m18 = 0
	}()
	b.m00 = 0x004510AC
	b.m20 = 0
	b.m24 = 0
	b.m28 = 0
	b.m2C = 0
	b.m30 = 0x0045101C
	b.m34 = 0x00451090
	b.m38 = 0
	b.m3C = 0
	b.m40 = 0
	b.m44 = 0
	b.m48 = 0
	b.m4C = 0
	b.m50 = 0
}

func (b *block) f0040DA60() {
	b.f0041987A()
	b.m00 = 0x0044ECF4
}

func (b *block) f0040F670() {
	b.m00 = 0x0044F568
}

type t0014AC88 struct {
	blocks [10]block // 0x00 0x98 0xEC
	// ...
	pathName [256]uint8 // 0x3F96C
	major    uint8      // 0x3FAFC
	minor    uint8      // 0x3FAFD
	patch    uint8      // 0x3FAFE
	name     [256]uint8 // 0x44928, "奇迹("
	dir      [260]uint8 // 0x44A28
	block2   [100]block // 0x44B40 0x44B94 0x44BEC 0x44C60 0x44CD4
}

func (t *t0014AC88) f0040AEC0(x uint) {
	// t.blocks[0].f0040E510(x, 0x66)
	func(x, y uint) {
		// t.f00415984(x, y)
		func(x, y uint) {
			t.blocks[0].f0041987A()
			t.blocks[0].m00 = 0x0045061C
			// f00433360memset(t.f54[:], 0, 32)
			t.blocks[0].m68 = y
			t.blocks[0].m54 = x
			t.blocks[0].m58 = x & 0x0000FFFF
		}(x, y)
		t.blocks[0].m00 = 0x0044F28C
		// t.f0040E430()
		func() {
			t.blocks[0].m74 = 1
			t.blocks[0].m78 = 0
			t.blocks[0].m84 = 0
			// v0048F4A0 = dll.user32.SetLayeredWindowsAttributes
		}()
	}(x, 0x66)
	t.blocks[0].m00 = 0x0044E654
	t.blocks[0].m88 = 0x0044E620
	t.blocks[1].f0041987A() // 0x98
	t.blocks[2].f0041987A() // 0xEC
	v0014AAC4state := 0
	v0014AAC4state = 2
	// v0018A75C.f00422CDD()
	v0014AAC4state = 3
	// v0014A75C.f00 = 0x0044E1FC
	// v0018A75C.f00422D3E(1)
	v0014AAC4state = 4
	t.block2[0].f0041987A() // 0x44B40
	t.block2[0].m00 = 0x0044E4AC
	v0014AAC4state = 5
	t.block2[1].f0041987A() // 0x44B94
	t.block2[1].m00 = 0x0044E4AC
	v0014AAC4state = 6
	// f0041E335().f0C()
	v0014AAC4state = 7
	t.block2[2].f0041987A() // 0x44BEC
	t.block2[2].m00 = 0x0044DCAC
	v0014AAC4state = 8
	t.block2[3].f0041987A() // 0x44C60
	t.block2[3].m00 = 0x0044DCAC
	v0014AAC4state = 9
	t.block2[4].f0041987A() // 0x44CD4
	t.block2[4].m00 = 0x0044DCAC
	v0014AAC4state = 10
	t.block2[5].f0041987A() // 0x44D48
	t.block2[5].m00 = 0x0044DCAC
	v0014AAC4state = 11
	t.block2[6].f0041987A() // 0x44D48
	t.block2[6].m00 = 0x0044DCAC
	v0014AAC4state = 12
	t.block2[7].f0041987A() // 0x44E30
	t.block2[7].m00 = 0x0044DCAC
	v0014AAC4state = 13
	t.block2[8].f0041987A() // 0x44EA4
	t.block2[8].m00 = 0x0044DCAC
	v0014AAC4state = 14
	t.block2[9].f0040DA60() // 0x44F18
	v0014AAC4state = 15
	t.block2[10].f0040F670() // 0x44F84
	v0014AAC4state = 16
	t.block2[11].f0040F670() // 0x44FAC
	v0014AAC4state = 17
	t.block2[12].f0040F670() // 0x44FEC
	v0014AAC4state = 18
	t.block2[13].f0040F670() // 0x45004
	v0014AAC4state = 19
	t.block2[14].f0040F670() // 0x45110
	v0014AAC4state = 20
	t.block2[15].f0040F670() // 0x45128
	v0014AAC4state = 21
	t.block2[15].f00417837()
	t.block2[15].f00417837()
	// LoadIcon()
	// t.f00406920()
	func() {
		// read config.ini
		// fileName := path.Join(filepath.Dir(os.Args[0]), "config.ini")
		// var version [10]uint8
		// dll.kernel32.GetPrivateProfileString("LOGIN", "Version", "0.0.0", version[:], 10, fileName[:])
		var minor, major, patch uint8
		minor = 4
		major = 1
		patch = 44

		// dll.kernel32.GetPrivateProfileString("LOGIN", "TestVersion", "0.0.0", version[:], 10, fileName[:])
		// var testmajor uint8
		// testmajor = 1
		if v0046327C == 1 {
			// ...
		}
		t.major = major
		t.minor = minor
		t.patch = patch
	}()
	println(v0014AAC4state)
	// memcpy(t.name[:], textcode(81)) // "奇迹(MU)"
	// f00433360memset(t.pathName, 0, 256)
}

func (t *t0014AC88) f00416045() {

}

type base interface {
	f50()
	f54()
	f68()
	f90()
}

// v004632B8
type service struct {
	m20 *t0014AC88
	m48 []uint8
}

// f00405AB0
func (s *service) f50() {
	// 栈初始化

	// CreateMutexA(NULL, TRUE, "MU AutoUpdate")

	// v0046F448.LoadWTF("message.wtf")

	// if 0 == f00403EE0().f00403C20() { // v00453224.f00403C20()
	// f0041E6D3("设备未初始化 请安装最新的图形驱动程序", 10, 0) // EUC-KR  --> translate.google.com
	// return
	// }

	// 注册表
	f00401870once().f004015D0Get() // 单例v00463180.f004015D0()

	// dd := struct {
	// 	displayDeviceName [128]uint8
	// }{}
	// f00433360memset(dd.buf[:], 0, 0x1A8)
	// EnumDisplayDevicesA(0, 0, &dd, 0)
	// if 0 != f00433461memcmp(f00401870once().f004017F0(), dd.displayDeviceName[:], 0x80) {
	// 	// not equal
	// }

	// if ERROR_ALREADY_EXISTS == GetLastError() {
	// 	return
	// }

	// f004223A2(0)
	if s.m48[0] != 0 {
		// ...
	}
	v004632B4port = 44405

	v0014AC88 := t0014AC88{}
	v0014AC88.f0040AEC0(0)
	s.m20 = &v0014AC88
	v0014AC88.f00416045()
}

func (s *service) f54() {}
func (s *service) f68() {}

// f0041FF0E
func (s *service) f90() {}

func f00449EEAwinMain() {
	// f004209C7/f00417837
	v004632B8 := func() base {
		return &service{}
	}()
	v004632B8.f90()
	v004632B8.f50()
	// ...
	// v006
}

// entry point, f00435036
func main() {
	// f00442399()
	// runtime function
	f00449EEAwinMain()
}
