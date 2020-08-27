package mumsg

import "os"

type msg struct {
	m00next  *msg
	m04id    int
	m08value string
}

type Msg struct {
	head      *msg
	tail      *msg
	msgs      [0x8000]*msg
	frameHead struct {
		m00flag uint8
		m01ver  uint8
		m18size int
	}
	err [32]uint8 // "Msg error"
}

func (t *Msg) msgListInit() bool {
	m := &msg{} // new
	if t.head = m; m == nil {
		// user32.MessageBox("Error")
		return false
	}
	t.head.m00next = nil
	t.head.m04id = 0
	t.head.m08value = ""
	t.tail = t.head
	return true
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
	m.m08value = string(text)
	t.head.m00next = m
	t.tail = m
	if index < 0 || index >= 0x7FFF {
		// user32.MessageBox(...)
		return
	}
	t.msgs[index] = m // index in [0,0x7FFF] is rapid
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
	if t.msgListInit() == false {
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
	if id < 0 || id >= 0x7FFF || t.msgs[id] == nil {
		for m := t.head.m00next; m != nil && m.m04id == id; m = m.m00next {
			return m.m08value
		}
		return t.head.m08value
	}
	return t.msgs[id].m08value
}
