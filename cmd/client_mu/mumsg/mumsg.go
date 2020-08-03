package mumsg

import "os"

type msg struct {
	m00     uintptr
	m04id   int
	m08text []uint8
}

type Msg struct {
	msgs      [0x8000]*msg
	frameHead struct {
		m00flag uint8
		m01ver  uint8
		m18size int
	}
	err [32]uint8 // "Msg error"
}

func (t *Msg) msgListInit() int {
	m := func(int) *msg { return &msg{} }(12) // new
	if t.msgs[0] = m; m == nil {
		// user32.MessageBox("Error")
		return 0
	}
	t.msgs[0].m00 = 0
	t.msgs[0].m04id = 0
	t.msgs[0].m08text = nil
	t.msgs[1] = t.msgs[0]
	return 1
}

func (t *Msg) msgListNew() *msg {
	return &msg{}
}

func (t *Msg) msgListAdd(index int, text []uint8) {
	l := len(text)
	if l == 0 {
		return
	}
	m := t.msgListNew()
	if m == nil {
		// user32.MessageBox(...)
		return
	}
	m.m04id = index
	m.m08text = make([]uint8, l)
	copy(m.m08text, text)
	t.msgs[1] = m
	t.msgs[2+index] = m
}

func (*Msg) xorBuffer(buf []uint8, size int) {
	if size == 0 {
		return
	}
	for i := 0; i < size; i++ {
		buf[i] = buf[i] ^ 0xCA
	}
}

func (t *Msg) dataFileLoadVer01(f *os.File) {
	if t.frameHead.m18size <= 0 {
		// user32.MessageBox(...)
		return
	}
	for t.frameHead.m18size > 0 {
		index := 0         // fread(s0014A932unk[:], 2, 1, f) // 0x0000
		size := 0xC        // fread(size[:], 2, 1, f) // 0x000C
		var buf [256]uint8 // fread(buf[:], size, 1, f)
		t.xorBuffer(buf[:], size)
		buf[size] = 0
		t.msgListAdd(index, buf[:])
		t.frameHead.m18size--
	}
}

func (t *Msg) LoadWTF(name string) {
	f, err := os.Open(name) // fopen(name, "rb")
	if err != nil {
		// user32.MessageBox(...)
		return
	}
	if t.msgListInit() == 0 {
		return
	}
	// fread(m.frame[:], 28, 1, f)
	if t.frameHead.m00flag != 0xCC {
		// user32.MessageBox(...)
		return
	}
	if t.frameHead.m01ver != 1 {
		// user32.MessageBox(...)
		return
	}
	t.dataFileLoadVer01(f)
	f.Close() // fclose(f)
}

func (t *Msg) Get(id int) string {
	// t := m.texts[id]
	// if t != nil {
	// 	return t.text
	// }
	return "msg error"
}
