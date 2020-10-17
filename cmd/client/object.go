package main

var v0860F290 [8]classDefault

// sizeof=0x29C
type classDefault struct {
	m00strength  uint16
	m02dexterity uint16
	m04vitality  uint16
	m06energy    uint16
	m08          uint16
	m0A          uint16
	m0E          uint8
	m0F          uint8
	m10          uint8
	m11          uint8
}

func f005A2F91(id int, strength, dexterity, vitality, energy uint16, unk1 uint16, unk2 uint16, unk3, unk4, unk5, unk6 uint8) {
	v0860F290[id].m00strength = strength
	v0860F290[id].m02dexterity = dexterity
	v0860F290[id].m04vitality = vitality
	v0860F290[id].m06energy = energy
	v0860F290[id].m08 = unk1
	v0860F290[id].m0A = unk2
	v0860F290[id].m0E = unk3
	v0860F290[id].m0F = unk4
	v0860F290[id].m10 = unk5
	v0860F290[id].m11 = unk6
}

func f005A3016() {
	f005A2F91(0, 18, 18, 15, 30, 80, 60, 1, 2, 1, 2)
	f005A2F91(1, 28, 20, 25, 10, 110, 20, 2, 1, 2, 1)
	f005A2F91(2, 50, 50, 50, 30, 110, 30, 110, 30, 6, 3)
	f005A2F91(3, 30, 30, 30, 30, 120, 80, 1, 1, 2, 2)
	f005A2F91(4, 30, 30, 30, 30, 120, 80, 1, 1, 2, 2)
	f005A2F91(5, 50, 50, 50, 30, 110, 30, 110, 30, 6, 3)
	f005A2F91(6, 32, 27, 25, 20, 100, 40, 1, 3, 1, 1)
}

var v08C88CACid int
var v08C88E0CconnID uint16

// master level
var v08C88E40exp uint64
var v08C88E48expNext uint64
var v08C88E58hpMax uint16
var v08C88E5AmpMax uint16
var v08C88E5CsdMax uint16
var v08C88E5EsdMax uint16

var v0805BBACobjectself *object

// sizeof=0x6BC
type object struct {
	m8            uint32
	m13class      uint8
	m14           uint8
	m15ctlCode    uint8
	m18guildTitle uint8
	m20pkLevel    uint8
	m1E           uint8
	m24           bool
	m2B           uint8
	m32           uint8
	m38name       [10]uint8
	m42           uint8
	m58           uint8
	m5C           int16
	m5Eid         uint16
	m64           [255]uint8
	m10C          uint16
	m164level     uint16
	m166          uint16
	m410          struct {
		m04  bool
		m0E  bool
		m28  uint8
		mB0  float32
		mB4  float32
		mB8  float32
		mBC  float32
		m114 float32
		m118 float32
	}
	m438 uint8
}

func (t *object) f004D332C(x uint16) {
	// t.m3D2 = 0xFFFF
}

// player
var v086105E0 []uint8 // 7k
var v086105E4 []uint8 // 91*800
var v086105E8player *player
var v086105ECpanel = &v086105E8player.m04panel

// sizeof=0x84
type bag struct {
	id    int
	slots [64]uint16
}

func (t *bag) f00A2678Ainit() {}

func (t *bag) f00A26765construct() {
	t.f00A2678Ainit()
}

// sizeof=0x1C30
type player struct {
	m00 int
	// normal level
	m04panel struct {
		m00name          [11]uint8
		m11              [255]uint8
		m10Aclass        uint8
		m10B             uint8
		m10Clevel        uint16
		m110exp          uint32
		m114expNext      uint32
		m118strength     uint16
		m11Aagility      uint16
		m11Cvitality     uint16
		m11Eenergy       uint16
		m120leadship     uint16
		m122hp           uint16 // hp
		m124mp           uint16 // mp
		m126hpMax        uint16 // max hp
		m128mpMax        uint16 // max mp
		m12Asd           uint16 // sd
		m12CsdMax        uint16 // max sd
		m140ag           uint16 // ag
		m142agMax        uint16 // max ag
		m144             uint16
		m146             uint16
		m148             uint16
		m14A             uint16
		m14CpointsAdd    uint16
		m14EpointsAddMax uint16
		m150pointsDec    uint16
		m152pointsDecMax uint16
		m154             uint16
		m160             uint16
		m174points       uint16
		m178             [100]uint8
	}
	m1344items          [12]item
	m19D4item           item // 可能是卷轴位置
	m1A60bag1           bag
	m1AE4items          [2]item
	m1BFC               [2]int
	m1C04money          uint
	m1C08moneyWarehouse uint
	m1C2C               uint8
	m1C2D               uint8
}

func (t *player) f005A31C1construct() {
	// f00DE8931(t.m1344[:], 0x8C, 12, f00A270D6construct, f00A270EA)
	t.m19D4item.f00A270D6construct()
	// t.m1A60.f00A26765() // t.m1A60.f00A2678A()
	// f00DE8931(t.m1AE4[:], 0x8C, 2, f00A270D6construct, f00A270EA)
}

func (t *player) f005A3337init() {
	f005A3016()
	for ebp4 := 0; ebp4 < 12; ebp4++ {
		t.m1344items[ebp4].m04code = -1
	}
	t.m19D4item.m04code = -1
	for ebp8 := 0; ebp8 < 2; ebp8++ {
		t.m1AE4items[ebp8].m04code = -1
		t.m1BFC[ebp8] = 0
	}
	t.m1C04money = 0
	t.m1C08moneyWarehouse = 0
	t.m1C2C = 0
	t.m1C2D = 0
}

func (t *player) f005A3727(changeup0 int) {

}

// objectPool
func f00592888getIndex(id int) int {
	return 0
}

func f004373C5objectPool() *objectPool {
	return v01308D04objectPool
}

var v01308D04objectPool *objectPool

type objectPool struct {
	m00vtabptr []uintptr
	m04objects []object
	m08object  *object
	m0C        int
}

func (t *objectPool) f00A38D5BgetObject(index int) *object {
	if index < 0 || index >= 400 {
		return nil
	}
	return &t.m04objects[index]
}

type t01173638 struct {
	objectPool
}

func (t *t01173638) f00A39237construct() *t01173638 {
	// t.m00vtabptr = v01173638[:]
	if v01308D04objectPool == nil {
		v01308D04objectPool = &t.objectPool
	}
	return t
}

var v01319D44objectPool *t01173630objectPool

// sizeof=0x10
type t01173630objectPool struct {
	t01173638
}

func (t *t01173630objectPool) f00A38C2Cinit() {
	// ebp14 := f00DE64BCnew(529*1724 + 4) // 0xDEA80: 912000
	// var ebp1C []uint8
	// if ebp14 != nil {
	// 	binary.LittleEndian.PutUint32(ebp14[:], 529)
	// 	f00DE8931(ebp14[4:], 1724, 529, f004D2B85, f004D2C3D)
	// 	ebp1C = ebp14[4:]
	// }
	ebp14 := &struct {
		number  int
		objects [529]object
	}{}
	ebp14.number = 529
	t.m04objects = ebp14.objects[:]
	t.m08object = &t.m04objects[f00DE8AADrand()%128]
	v0805BBACobjectself = t.m08object
	t.m0C = 0
}

func (t *t01173630objectPool) f00A38B97construct() *t01173630objectPool {
	t.t01173638.f00A39237construct()
	// t.m00vtabptr = v01173630[:]
	t.f00A38C2Cinit()
	return t
}

//
var v08C88FAC bool
var v012E4080 int16
var v08C88DB8 int

func f006C64A1init() {
	v08C88FAC = false
	v012E4080 = -1
	v08C88DB8 = 0
	for i := 0; i < 400; i++ {
		obj := v01308D04objectPool.f00A38D5BgetObject(i)
		// f006BF9A3(obj)
		func(obj *object) {
			ebp4 := &obj.m410
			if ebp4.m04 == false {
				return
			}
			obj.m1E = 0
			if v0805BBACobjectself.m5C == -1 || obj.m5C == v0805BBACobjectself.m5C {
				obj.m1E = 1
			}
			if v08C88FAC == false {
				return
			}
			if v012E4080 == -1 && v08C88DB8 != 0 {
				// v012E4080 = f006BF959(&v08C88DB8)
			}
			if v012E4080 >= 0 && obj.m5C == v012E4080 {
				obj.m1E = 2
			}

		}(obj)
	}
}

func f004DAAB3setMaxClass(n uint8) {
	v0131A233 = n
}

func f0059CFA6class(class uint8) uint8 {
	// packet class: bit567-class, bit4-changeup1, bit3-changeup2
	// client class: bit012-class, bit3-changeup1, bit4-changeup2
	return class>>4&1<<3 | class>>3&1<<4 | class>>5
}

func f004398F6changeup0(class uint8) uint8 {
	return class & 7
}

var v0114AAE4 float32 = 0.3

func f0059CA40newObject(id int, class uint8, unk1 uint8, unk2, unk3, dir float32) *object {
	ebp4obj := f004373C5objectPool().f00A38D5BgetObject(id)
	if ebp4obj == nil {
		return nil
	}
	// f00592B90objSet(ebp4obj, 0x4C4, 0, 0, dir)
	ebp8 := &ebp4obj.m410
	ebp8.mB4 = v0114AAE4
	ebp8.mB8 = v0114AAE4
	ebp8.mBC = v0114AAE4
	ebp4obj.m5Eid = uint16(id)
	ebp8.m114 = unk2
	ebp8.m118 = unk3
	ebp4obj.m13class = uint8(class)
	ebp4obj.m14 = unk1
	// ebp8.m1FC.f004CF282()
	if v012E2340state == 2 {
		// ebp4obj.m1FC = 0x1314
		// ebp4obj.m1FE = 7
		// ebp4obj.m220 = 0x1514
		// ebp4obj.m222 = 7
		// ebp4obj.m244 = 0x1714
		// ebp4obj.m246 = 7
		// ebp4obj.m268 = 0x1914
		// ebp4obj.m26A = 7
		// ebp4obj.m28C = 0x1B14
		// ebp4obj.m28E = 7
		// ebp4obj.m2B0 = 0x906
		// ebp4obj.m2B2 = 13
		// ebp4obj.m2D4 = 0x1EFD
		// ebp4obj.m2F8 = 0x1F16
		// ebp4obj.m31C = 0x1EFC
	} else {
		// ebp4obj.m1FC = f004398F6changeup0(class) + 0x2512
		// ebp4obj.m220 = f004398F6changeup0(class) + 0x252A
		// ebp4obj.m244 = f004398F6changeup0(class) + 0x2542
		// ebp4obj.m268 = f004398F6changeup0(class) + 0x255A
		// ebp4obj.m28C = f004398F6changeup0(class) + 0x2572
		// ebp4obj.m2B0 = 0xFFFF
		// ebp4obj.m2D4 = 0xFFFF
		// ebp4obj.m2F8 = 0xFFFF
		// ebp4obj.m31C = 0xFFFF
		// ebp4obj.m2D0 = 0xFFFF
		ebp4obj.f004D332C(0xFFFF)
		// ebp4obj.m3D0 = 0xFF
	}
	// f00736F42(28, ebp8)
	ebp4obj.m15ctlCode = 0
	// f00593E76(ebp4obj)
	// f0055E188(ebp4obj)
	return ebp4obj
}

func f00594B19objSet(id int, inventory []uint8, unk1, unk2 int) {

}
