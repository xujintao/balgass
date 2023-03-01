package main

import (
	"os"
	"sync"

	"github.com/xujintao/balgass/cmd/client_mu/mumsg"
)

// msg
var v0046F448msg mumsg.Msg

// reg
var v00463220once sync.Once
var v00463180reg *reg

func f00401870reg() *reg {
	v00463220once.Do(func() {
		v00463180reg = new(reg)
	})
	return v00463180reg
}

type reg struct {
	UserID                  [11]uint8  // 0x00
	DisplayDeviceName       [129]uint8 // 0x0B
	DisplayDeviceModelIndex int        // 0x8C
	FullScreenMode          int        // 0x90
	DisplayColorBit         int        // 0x94
	MusicOn                 int        // 0x98
	SoundOn                 int        // 0x9C
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

func (r *reg) f004015D0get() bool {
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
	// advapi32.RegQueryValueEx(hReg, "UserID", 0, 0, r.UserID[:], &size)

	// size = 128
	// advapi32.RegQueryValueEx(hReg, "DisplayDeviceName", 0, 0, r.DisplayDeviceName[:], &size)

	// size = 4
	// advapi32.RegQueryValueEx(hReg, "DisplayDeviceModelIndex", 0, 0, &r.DisplayDeviceModelIndex, &size)
	// advapi32.RegQueryValueEx(hReg, "FullScreenMode", 0, 0, &r.FullScreenMode, &size)
	// advapi32.RegQueryValueEx(hReg, "DisplayColorBit", 0, 0, &r.DisplayColorBit, &size)
	// advapi32.RegQueryValueEx(hReg, "MusicOn", 0, 0, &r.MusicOn, &size)
	// advapi32.RegQueryValueEx(hReg, "SoundOn", 0, 0, &r.SoundOn, &size)
	// RegCloseKey(hReg)
	// r.f004014B0Set()
	return true
}

func (r *reg) f004017F0getDeviceName() []uint8 {
	return r.DisplayDeviceName[:]
}

// -----------------------------------------------------------------------------
// card
var v00463254once sync.Once
var v00463224card card

func f00403EE0card() *card {
	v00463220once.Do(func() {
		v00463224card.f00403E70construct()
		// f00432C21(f0044C3C0)
	})
	return &v00463224card
}

type DEVMODEA struct {
	dmBitsPerPel       uint32 // 0x68
	dmPelsWidth        uint32 // 0x6C
	dmPelsHeight       uint32 // 0x70
	dmDisplayFrequency uint32 // 0x78
}
type DEVMODEx struct {
	index   int
	devMode DEVMODEA
	index2  int
	id      int // EnumDisplaySettings
	w       int
	h       int
	isWide  bool
}

// size:0xC4
type mode struct {
	m00      *mode
	m04      *mode
	m08      *mode
	devModex DEVMODEx
	mC0      bool
	mC1      bool
}

type card struct {
	m00      **card
	m18modes []*mode // tree
	m1Csize  int
	m20      **card
	m24      *mode
	m28      int
	m2C      int
}

func (t *card) f00403E70construct() {
	// t.f00403780()
	func() {
		c := new(*card) // f00414B5Cnew(4)
		*c = t
		t.m00 = c
		// t.f00402FD0()
		t.m18modes = func() []*mode {
			return nil
		}()
		t.m1Csize = 0
	}()
	t.m20 = nil
	t.m24 = nil
	// t.f00403020()
}

func (t *card) f00403080filter(devMode *DEVMODEA) bool {
	if devMode.dmDisplayFrequency != 60 { // hook to <48
		return false
	}
	if devMode.dmPelsWidth < 800 || devMode.dmPelsHeight < 600 || devMode.dmBitsPerPel != 32 {
		return false
	}
	for _, node := range t.m18modes {
		m := node.devModex.devMode
		if devMode.dmPelsWidth == m.dmPelsWidth &&
			devMode.dmPelsHeight == m.dmPelsHeight {
			// && devMode.dmDisplayFrequency == m.dmDisplayFrequency { // hook discard the equality operator of dmDisplayFrequency, fixed resolution mismatching
			return false
		}
	}
	return true
}

// 640*480(4:3), remove
// 800*600(4:3)
// 1024*768(4:3)
// 1152*864(4:3)
// 1280*960(4:3)
// 1280*1024(5:4), remove
// 1400*1050(4:3)
// 1920*1080(16:9)
func (t *card) f00403C20() bool {
	// t.f00403020()
	t.m28 = 0
	t.m2C = -1
	var devMode DEVMODEA
	index := 0
	id := 0
	// var buf [3]uint
	modex := DEVMODEx{}
	for {
		// if user32.EnumDisplaySettings(0, id, &devMode) == false {
		// 	break
		// }
		modex.index = index
		modex.id = id

		if t.f00403080filter(&devMode) == true {
			switch {
			case devMode.dmPelsWidth/4*3 == devMode.dmPelsHeight:
				modex.w = 4
				modex.h = 3
			case devMode.dmPelsWidth/16*9 == devMode.dmPelsHeight:
				modex.w = 16
				modex.h = 9
				modex.isWide = true
			case devMode.dmPelsWidth/16*10 == devMode.dmPelsHeight:
				modex.w = 16
				modex.h = 10
				modex.isWide = true
			default:
				goto label1
			}
			modex.index2 = index
			modex.devMode = devMode
			// t.f00403B20addMode(buf[:], &modex)
			func() {
				// t.f00403810()
				func() {
					// t.f00402F80()
					func() {

					}()
				}()
			}()
			index++
		}
	label1:
		id++
	}
	return false
}

var v004633CCfd *os.File
