package main

import "sync"

// key recorder
func f008AF00D() *t09D24B80 {
	v09D24C80once.Do(func() {
		v09D24B80.f008AEFD5init()
		// f00DE8BF6atexit(f011485E8)
	})
	return &v09D24B80
}

var v09D24B80 t09D24B80
var v09D24C80once sync.Once

type t09D24B80 struct {
	vkeys [256]uint8
}

func (t *t09D24B80) f008AEFD5init() {

}

func (t *t09D24B80) f008AF06ArecordKey() {
	for i := 0; i < 256; i++ {
		s := 0 // dll.user32.GetAsyncKeyState(i)
		if s&0x8000 != 0 {
			switch t.vkeys[i] {
			case 0, 1:
				t.vkeys[i] = 2
			case 2:
				t.vkeys[i] = 3
			}
		} else {
			switch t.vkeys[i] {
			case 1:
				t.vkeys[i] = 0
			case 2, 3:
				t.vkeys[i] = 1
			}
		}
	}
}

func (t *t09D24B80) f008AF121(vkey uint8) uint8 {
	v := t.vkeys[vkey]
	// return v>0?0,1;
	if v > 0 {
		v = 1
	}
	return v
}

func (t *t09D24B80) f008AF156(vkey uint8) uint8 {
	v := t.vkeys[vkey]
	if v == 2 {
		v = 1
	}
	return v
}

func (t *t09D24B80) f008AF170(vkey uint8) uint8 {
	v := t.vkeys[vkey]
	if v == 3 {
		v = 1
	}
	return v
}

func f008AEFAD(vkey uint8) uint8 {
	return f008AF00D().f008AF156(vkey)
}

func f008AEFC1(vkey uint8) uint8 {
	return f008AF00D().f008AF170(vkey)
}

var v086A3778hotkeys [256]int

func f005A4BC5queryHotKey(vkey uint8) bool {
	if f008AEFAD(vkey) == 1 {
		if v086A3778hotkeys[vkey] == 0 {
			v086A3778hotkeys[vkey] = 1
			return true
		}
	} else {
		v086A3778hotkeys[vkey] = 0
	}
	return false
}
