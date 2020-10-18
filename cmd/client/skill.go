package main

import (
	"sync"
	"unsafe"

	"github.com/xujintao/balgass/win"
)

// skill
func f005A477A(id int) bool {
	// switch f0090E94C().f0090EA51().f008F6A00(id) {
	// // Nova, Beast Uppercut, 光速拳
	// case 40, 261, 263:
	// 	return true
	// }
	return false
}

func f005A457A(id int) {
	ebp4 := v086105E4[id].m2C
	if ebp4 <= 0 {
		return
	}
	for ebp8 := 0; ebp8 <= int(v086105ECpanel.m176skillNum); ebp8++ {
		if int(v086105ECpanel.m178skillIDs[ebp8]) == id && v086105ECpanel.m918[ebp8] == 0 {
			v086105ECpanel.m918[ebp8] = ebp4
		}
	}
}

// sizeof=0x70
type t086105E4 struct {
	m2C int
}

func f00A09635skill() *t09D91F4Cskill {
	v09D91F74once.Do(func() {
		v09D91F4C.f00A095DAinit()
		// f00DE8BF6atexit(f011486FF)
	})
	return &v09D91F4C
}

var v09D91F4C t09D91F4Cskill
var v09D91F74once sync.Once
var v08C88F68skillUseTimestamp uint32

type t09D91F4Cskill struct {
	m00obj      *object
	m04         *object410
	m08         int
	m0C         int
	m10id       int
	m14id       int
	m1CtargetID int
	m20         float32
	m24         bool
}

func (t *t09D91F4Cskill) f00A095DAinit() {
	t.m00obj = v0805BBACobjectself
	t.m04 = &v0805BBACobjectself.m410
	t.m08 = 0
	t.m0C = 0
	t.m14id = 0
	t.m20 = 0.0
	t.m24 = false
}

func (t *t09D91F4Cskill) f00A14596() bool {
	return false
}

// f00A0A5E1->0x09E25AB0
// s9 f0092E01C->0x0A37F316
func (t *t09D91F4Cskill) f00A0A5E1() bool {
	// 0x0AD31885 hook to hide 0x09E25AB0, disable what?
	// push 0x09FC14E5
	// push 0x0A0084AA
	// ret

	// 0x0A0084AA
	// push 0x00D0C808
	// push 0x0A0C4915
	// ret

	// 0x0A0C4915
	// push 0x00C04000
	// push 0x0A4455E0
	// ret

	// 0x0A4455E0
	// push ecx edi ebx esi edx fd

	// 0x0A9F5744
	// push 0x0A95171B
	// push 0x0AD98CB3
	// ret

	// 0x0AD98CB3
	// push ebx esi edi
	var ecx uint32
	if v0AA08598 != v0AD7A978 {
		// 0x0A9F74BC
		ecx = *(*uint32)(unsafe.Pointer(uintptr(v0A391C78 - v0A9D361D))) // 0x7FFE0014地址处的值一直在变化，体现了随机性
		v0AA08598 = v0AD7A978
	} else {
		// 0x0ABB5B0F
		ecx = v0A562D19
	}
	// 0x0ABD5366
	if ecx == 0 {
		// 0x09FCBB36
		ecx = 1
	}
	// 0x0A9F74CD
	v0A562D19 = v0AD82FDD * ecx % v0A338B93
	if v0A562D19 > v09F876D2 {
		// 0x0AD92B1F
		// [ebp+v0AD574B8*4+8] = v0A0519E0label1
		// [ebp+v0A059559*4+8] = v0AD72EA4label2
		// [ebp+v0A048EEB*4+8] = v0A3A8521label3
		// pop edi esi ebx
		// leave
		// 0x0A95171B
	} else {
		// 0x0A05EB51
		// 壳业务
		// ebpC := v0AD93A11
		// ebp10 := v0AD93A0D[0]
		// // ebp18 := v0A33BBC9blocks[:]
		// // eax = &ebp10
		// // push 0x0A7434E5
		// // push eax
		// // ebx = 0x0A392992
		// ebp18 := v0AD93A0D[:]
		// // [ebp+v0AD574B8*4+8] = v0A0519E0label1
		// // push 0x0A4E10F4
		// // push 0x012DDE52
		// // ret
		// // 0x012DDE52
		// edx := v0AD93A0D[0]
		// ecx := v0AD93A11
		// // push ebx
		// ebx := v0AFD8339
		// // push ebp
		// // push esi
		// esi := v0A74649D
		// esi ^= ebx
		// // push edi
	}
	// 0x0A95171B
	// pop fd edx esi ebx edi ecx
	// label3(0x0AFD4238)
	// label2(0x0A057FB2)
	// label1(0x09E2729C)

	// true logic
	// 0x09E25AB0 0x0A950A8D 0x09EB974E 0x00A0A5FE
	f00DE8A70chkstk() // 0x23100局部变量
	// 0x0A9FAFD5 0x00A0A65A 0x0A83CE92
	ebp230C8 := 0
	/*
		t.m04.m128 = f00550BED(t.m04.m114, t.m04.m118, t.m00obj.m1C0, t.m00obj.m1C4) // 296.245
		// 0x00A0A68C 0x0A4418EE 0x00A0A6B5 0x09FF745D
		ebp1A := uint8(f00DE76F6round(float64(t.m04.m128) / 360 * 255))         // 0xD1
		ebp1B := uint8(f00DE76F6round((float64(t.m04.m128) + 180) / 360 * 255)) // 0x51
		// 0x00A0A6CB 0x09E219C5 0x00A0A6DF 0x0A60144E
		// 0x00A0A70B 0x09E8AD99
		ebp19 := f005513BC(t.m04.m114, t.m04.m118, f00DE76C0roundf(t.m00obj.m1C0), f00DE76C0roundf(t.m00obj.m1C4)) // 0x77
	*/
	ebp18 := 1.0
	ebp14 := 1.0
	ebp10 := 1.0
	ebp4 := 0
	switch t.m14id {
	case 18:
		// 0x00A0A91E
	case 19:
		// 0x00A0AB00
	case 20:
		// 0x00A0ACFE
	case 21:
		// 0x00A0AF5D
	case 22:
		// 0x00A0B1C3
	case 23:
		// 0x00A0B430
	case 41: // 霹雳回旋斩
		// 0x00A0BE11
	case 42: // 雷霆裂闪
		// 0x00A0BFD1
	case 43: // 袭风刺
		// 0x00A0C191 0x0AF80590 0x00A0C19C
		var ebp22D24 stdstring
		ebp22D24.f00406FC0stdstring([]byte("webzen"))
		// 0x0A55742C
		ebp4 = 33
		// 0x0A331BF5
		ebp230C8 |= 0x40000
		// 0x09E037FD 0x00A0C1C8
		var ebp22D40 stdstring
		ebp22D40.f00406FC0stdstring(v0805BBACobjectself.m38name[:])
		// 0x0A052CE2
		ebp4 = 34
		ebp230C8 |= 0x80000
		var ebp230F8 bool
		// 0x09FD746A 0x00A0C1F3 0x0A4E942D 0x0A050889
		if f004CE0BDstrstr(&ebp22D40, &ebp22D24) ||
			(f005A477A(t.m10id) == false && f00DF08E8abs(int(win.GetTickCount()-v08C88F68skillUseTimestamp)) <= 300) {
			// 0x00A0C23E 0x0AF79988
			ebp230F8 = false
		} else {
			// 0x00A0C232 0x09FB8382
			ebp230F8 = true
		}
		// 0x00A0C245 0x0A048F6F
		ebp22D05 := ebp230F8
		ebp4 = 33
		if ebp230C8&0x80000 != 0 {
			// 0x00A0C265 0x0AF80B75
			ebp230C8 &= 0xFFF7FFFF
			ebp22D40.f00407B10free()
		}
		// 0x0AAB2E1A
		ebp4 = -1
		if ebp230C8&0x40000 != 0 {
			// 0x0A8FB1F2 0x00A0C29B
			ebp230C8 &= 0xFFFBFFFF
			ebp22D24.f00407B10free()
		}
		// 0x0AD2EDC7
		if ebp22D05 {
			// 0x00A0C2AF 0x09F8B7E6
			v08C88F68skillUseTimestamp = win.GetTickCount()
			// 0x00A0C2C0
			var ebpF674skillAttack pb
			ebpF674skillAttack.f00439178init()
			// 0x0AF0DCF0
			ebp4 = 35
			ebpF678skillID := t.m10id
			// 0x00A0C2EA
			ebpF674skillAttack.f0043922CwriteHead(0xC1, 0x19)
			// 0x0A43F254 0x00A0C32E 0x00A0C2BA 0x0A93630F
			// 0x00A0C335 0x00A0C2CA 0x09EBAE2F 0x00A0C33C
			// 0x00A0C2CF 0x0AFD402F 0x00A0C343
			ebpF674skillAttack.f004397B1writeUint8(uint8(t.m1CtargetID >> 8))
			ebpF674skillAttack.f004397B1writeUint8(uint8(ebpF678skillID >> 8))
			ebpF674skillAttack.f004397B1writeUint8(uint8(t.m1CtargetID))
			ebpF674skillAttack.f004397B1writeUint8(uint8(ebpF678skillID))
			// 0x0A566077 0x00A0C352
			ebpF674skillAttack.f004393EAsend(true, false)
			// 0x09FCA7C6 0x00A0C360 0x0A32997F
			f005A457A(t.m10id)
			ebp4 = -1
			// 0x00A0C370
			ebpF674skillAttack.f004391CF()
		}
		/*
			// 0x00A0C375 0x0AD2E3E4 0x00A0C380 0x09FD53A8
			f005521C7(t.m04, 71)
			// 0x0ABB47CD 0x0A38EA9A 0x0A55D67E 0x0A935814
			// 0x0A338DDC 0x0A4ED4F7 0x00A0C3BF 0x0AD35A50
			f008707A9(0x7D73, &t.m04.m114, &t.m04.m120, &ebp18, 0, 0.0, t.m04)
			// 0x00A0C3CB 0x09FC352D 0x00A0C3DC 0x00A0C37A
			// 0x0A05B713 0x00A0C3E0 0x09FB5250
			f007DAFE0(f00DE8AADrand()%2+60, 0, 0)
		*/
	case 44: // Crescent Moon Slash
		// 0x00A0C3ED
	case 47: // 钻云枪
		// 0x00A0B69D
	case 48: // 生命之光
		// 0x00A0B984
	case 49: // 流星焰
		// 0x00A0BB77
	case 55: // 玄月斩
		// 0x00A0CC11
	case 56: // 天雷闪
		// 0x00A0CDD1
	case 57: // Spiral Slash
		// 0x00A0C64C 0x0AFDD4FF 0x00A0C657
		var ebp22D9C stdstring
		ebp22D9C.f00406FC0stdstring([]byte("webzen"))
		// 0x09E27F89
		ebp4 = 39
		ebp230C8 |= 0x40000
		// 0x00A0C683
		var ebp22DB8 stdstring
		ebp22DB8.f00406FC0stdstring(v0805BBACobjectself.m38name[:])
		// 0x0AF938ED
		ebp4 = 40
		// 0x0AD91B0F 0x0AA05F47
		ebp230C8 |= 0x800000
		// 0x0AFD8DD8 0x00A0C6AE 0x0A9004F9 0x0A92FF5C
		if f004CE0BDstrstr(&ebp22DB8, &ebp22D9C) {
			// 0x00A0C6F9
		} else {
			// 0x0AFE2129
		}
	case 60, 61:
		// 0x00A0D3A7
	case 62:
		// 0x00A0DB9E
	case 63:
		// 0x00A0E068
	case 64:
		// 0x00A0E2B6
	case 78:
		// 0x00A0D90E
	case 232:
		// 0x00A0C8AB
	case 236:
		// 0x00A0D020
	case 237:
		// 0x00A0D1E4
	case 238:
		// 0x00A0D672
	case 344:
		// 0x00A0CA5E
	case 521:
		// 0x00A0E2B6
	}
	// 0x00A0E492 0x00A0E4C7 0x0ABD65E3 0x00A0E494 0x0AD816B9
	print(ebp4, ebp10, ebp14, ebp18)
	return true
}

func (t *t09D91F4Cskill) f00A0E4A0() bool {
	return false
}

func (t *t09D91F4Cskill) f00A0FE1A() bool {
	return false
}

func (t *t09D91F4Cskill) f00A11A00() bool {
	return false
}

func (t *t09D91F4Cskill) f00A133DD() bool {
	return false
}

func (t *t09D91F4Cskill) f00A14890() bool {
	return false
}

// s9 f0092DEF9
func (t *t09D91F4Cskill) f00A0A4BE() bool {
	// f00560030(t.m04)
	if f0050E6D7().f0050E893isIllusionTemple(v012E3EC8mapNumber) && f0050E6D7().f0050F876(t.m14id) {
		return t.f00A14596()
	}

	ebp4 := f004398F6changeup0(t.m00obj.m13class)
	// 战士 魔剑 圣 格斗
	if (ebp4 == 1 || ebp4 == 3 || ebp4 == 4 || ebp4 == 6) && t.f00A0A5E1() == false {
		return false
	}
	// 法师 魔剑 召唤
	if (ebp4 == 0 || ebp4 == 3 || ebp4 == 5) && t.f00A0E4A0() == false {
		return false
	}
	// 弓箭手
	if ebp4 == 2 && t.f00A0FE1A() == false {
		return false
	}
	// 召唤
	if ebp4 == 5 && t.f00A11A00() == false {
		return false
	}
	// 格斗
	if ebp4 == 6 && t.f00A133DD() == false {
		return false
	}
	if t.f00A14890() == false {
		return false
	}
	return true
}

// s9 f0092D873
func (t *t09D91F4Cskill) f00A09E34() {
	// ...
	// 0x00A0A472, s9 0x0092DEAD:
	t.f00A0A4BE()
	// ...
}

func (t *t09D91F4Cskill) f00A09A07() {
	// ...
	// 0x00A09E11, s9 0x0092D84D:
	t.f00A09E34()
	// ...
}
