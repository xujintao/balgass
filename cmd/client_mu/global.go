package main

import (
	"sync"

	"github.com/xujintao/balgass/cmd/client_mu/mumsg"
)

// msg
var v0046F448msg mumsg.Msg

// reg
var v00463220once sync.Once
var v00463180Reg *reg

func f00401870reg() *reg {
	v00463220once.Do(func() {
		v00463180Reg = new(reg)
	})
	return v00463180Reg
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
