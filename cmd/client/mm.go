package main

import "sync"

type t21 struct{}

func (t *t21) f00D06BF0(size uint, y *int) uintptr {
	return 0
}

type t20 struct {
	m00 *t01319F18
	m08 int // 0x10
	m0C t21
}

func (t *t20) f00BFAFC0(size uint, y *int) uintptr {
	// b := false
	x := 0
	return t.m0C.f00D06BF0(size, &x)
}

func (t *t20) f00BFB480(size uint, y *int) uintptr {
	if t.m08 > 0x10 {
		// return t.f00BFB280(size, t.m08, y)
	}
	if size <= 0x200 {
		size = (size + 0xF) & ^uint(0xF) // 16 byte align
		return t.f00BFAFC0(size, y)
	}
	return 0
}

var v09D9BD74mm *mm // 0x0EB40B84

type mm struct {
	m00      uintptr // 虚表: v01187810
	m4Cmutex sync.Mutex
	m64      bool
	m68      *t20
}

// f00BAAFB0
func (t *mm) do11malloc(size uint, x int) uintptr {
	ebpC := 0
	if t.m64 {
		defer t.m4Cmutex.Unlock()
		t.m4Cmutex.Lock()
		return t.m68.f00BFB480(size, &ebpC)
	}
	return t.m68.f00BFB480(size, &ebpC)
}

func f00A3BA10newobject(size uint) uintptr {
	// return f00A3BA24(id)
	return func(size uint) uintptr {
		return v09D9BD74mm.do11malloc(size, 0)
	}(size)
}
