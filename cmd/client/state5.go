package main

import (
	"encoding/binary"
	"unsafe"
)

var v012E3EC8 int

// object sizeof=0x6BC
type object struct {
	m5Eid uint16
	m166  uint16
	m438  uint8
}

// objectManager
var v012E31B0 int = 0
var v012E3200 int = 5
var v01308D04objectManager = &objectManager{}

type objectManager struct {
	m08objects []object
}

func (t *objectManager) f00A38D5BgetObject(index int) *object {
	if index < 0 || index >= 0x190 {
		return nil
	}
	return &t.m08objects[index]
}

// s9 f005A774E
func f005B0120() {
	// ...
	// 0x005B0208 inline send, s9 f005B67D9
	func() {
		var ebp2918 pb
		ebp2918.f00439178init()
		ebp2918.f0043922CwritePrefix(0xC1, 0xD7) // send position, hook D7->D4
		// ...
	}()
	// ...
}

// s9 f005A531A
func f00611C16() {
	var c int // c := x.m59
	switch c {
	case 1:
		// 0x0061385D
	case 2:
		// 0x00618C52
		// v012E3140 = 12
		ebp69D8 := v01308D04objectManager.f00A38D5BgetObject(v012E31B0) // 0x12DABC9C
		if ebp69D8.m438 != 4 {
			v012E31B0 = -1
			return
		}
		ebp69DC := v086105EC.m10C                 // 0x190
		if ebp69D8.m166 == 0xEE && ebp69DC < 10 { // 0xF0
			// ...
			return
		}
		if ebp69D8.m166 == 0xF3 ||
			ebp69D8.m166 == 0xF6 ||
			ebp69D8.m166 == 0xFB ||
			ebp69D8.m166 == 0x1A0 ||
			ebp69D8.m166 == 0x242 ||
			ebp69D8.m166 == 0x2AE ||
			ebp69D8.m166 == 0x2AF {
			// f00A49798().mF0.f00A700C1(1)
		} else {
			// f00A49798().mF0.f00A700C1(0)
		}
		// if f0090E94C().f0090D682(4) {
		// 	f0090E94C().f0090DC7E(4)
		// }
		// if v01351ABC.f005025E9() {
		// ... inline
		// }
		// if f00989D02() {
		// 	... inline
		// }
		// if f00983A58() {
		// 	... inline
		// }
		// if f0093F498().f0093F4F5() {
		// 	... inline
		// }
		// if ebp69D8.m166 >= 0x1D4 && ebp69D8.m166 <= 0x1DB {
		// 	... inline
		// }
		// if ebp69D8.m166 == 0x1DE{
		// 	... inline
		// }
		// if ebp69D8.m166 == 0x21C{
		// 	... inline
		// }
		// if ebp69D8.m166 == 0x223{
		// 	... inline
		// }
		// if ebp69D8.m166 == 0x243{
		// 	... inline
		// }
		// if ebp69D8.m166 == 0x283{
		// 	... inline
		// }
		// if ebp69D8.m166 == 0x28B{
		// 	... inline
		// }
		// if ebp69D8.m166 == 0x181{
		// 	... inline
		// }
		// if ebp69D8.m166 == 0x2AA{
		// 	... inline
		// }
		// if ebp69D8.m166 == 0x29D{
		// 	... inline
		// }
		// if f0094E791() {
		// 	... inline
		// }
		// ebp59538 := f0098C3F9()
		if v012E3EC8 >= 0x62 && v012E3EC8 <= 0x63 {
			// ... inline
		} else if v012E3EC8 < 0x2D || v012E3EC8 > 0x32 {
			// 0x00628420 inline send
			var ebp1896C pb // 0x0014F870
			ebp1896C.f00439178init()
			// ebp4 := 51
			ebp1896C.f0043922CwritePrefix(0xC1, 0x30) // talk
			var ids [2]uint8
			binary.BigEndian.PutUint16(ids[:], ebp69D8.m5Eid) // bigendian 0x0001
			ebp1896C.f00439298writeBuf(ids[:], 2, true)
			ebp1896C.f004393EAsend(true, false)
			// 0x00629345
		}
	case 3:
		// 0x00611C9B, s9 0x005A539E
		// ...
		ebp20 := v01308D04objectManager.f00A38D5BgetObject(v012E3200) // 0x12DC9B08
		if ebp20 == nil {
			return
		}
		// ...
		var ebp14 uint8 // f00DE76F6()
		// ...
		// 0x00611FF5 inline send, s9
		var ebp14A4 pb // 0x00166D38
		ebp14A4.f00439178init()
		// ebp4 := 2
		ebp14A4.f0043922CwritePrefix(0xC1, 0xD9) // normal attack, hook D9->11
		var ids [2]uint8
		binary.BigEndian.PutUint16(ids[:], ebp20.m5Eid) // bigendian 0x166F
		ebp14A4.f00439298writeBuf(ids[:], 2, true)
		ebp14A4.f004397B1writeUint8(0x78)  // AttackAction
		ebp14A4.f004397B1writeUint8(ebp14) // DirDis
		ebp14A4.f004393EAsend(false, false)
	case 4:
		// 0x0062C00B
	case 5:
		// 0x0061362F
	case 6:
		// 0x006314F8
	}
}

// s9 f004E077A
func f004DF0D5handleState5() {
	// 很复杂
	// ...
	// f005B5D6C
	func() {
		// 0x0A3343F6 hook to hide f09FD38F9, disable anti-temper with backup code
		// push 0x009E5BF7
		// push 0x0AAB4C21
		// ret

		// 0x0AAB4C21
		// push 0x0052B20E
		// push 0x0B072E7E
		// ret

		// 0x0B072E7E
		// push 0x00D36E4F
		// push 0x0AAB34A8
		// ret

		// 0x0AAB34A8
		// push edi fd ebx ecx esi

		// 0x09E0464A
		// push 0x0AF7F54A
		// push 0x0A000809
		// ret

		// f0A000809
		func() {
			// 0x0AB4DEB8 0x09FD8038 0x0A934040
			// push esi
			if v0ABD75F3 != 0x0A3268A5 {
				// 0x0AFDFE5A 0x0A888616 0x0A55768D 0x0AA09787
				// push eax
				// push edx
				v0AA09787 = v0A0C2FD2 ^ v0AF8FD45
				// rdtsc
				var tscLeax uint32
				v0AAB36E2 = tscLeax
				// pop edx
				// pop eax
				v0ABD75F3 = v0A3268A5
			}
			// 0x0AF11EA3
			if v0AAB36E2 == 0 {
				// 0x09FCB7AE
				v0AAB36E2++
			}
			// 0x0AD7178C
			v0AAB36E2 = v0AD6FFBE * v0AAB36E2 % v0ABB4649 // 0x2426
			if v0AAB36E2 > v09EAE2F8 {
				// 0x0A3905DC 0x0AB4C513 0x0A84228B 0x0A564287 0x0AC3A1F9
				// [ebp+v0A9FB214*4+8] = v0A443580
				// [ebp+v0AFDFD32*4+8] = v09FB8B32
				// [ebp+v0A849070*4+8] = v0A32C145
				return // 0x0A7473D5
			}
			// 0x0AF83916
			// ebp10 := v0AF781CF
			// ebpC := 0
			ebp18 := v0A33BBC9blocks[:]
			// [ebp+v0A9FB214*4+8] = v0A443580
			for {
				// 0x09FD69FF
				if ebp18[0].addr == ^uintptr(0) {
					// 0x09F839BF
					// [ebp+v0AFDFD32*4+8] = v09FB8B32
					// ebp20 := &ebp10
					// ebp4 := v0A5FCB8Bcodes[:]
					ebp8 := v0A84B964 - 1
					// 0x0A88ED64
					if ebp8 < 0 {
						// 0x0A555686
					} else {
						// 0x0AD754E2
					}
				} else {
					// 0x0AA0D619 0x0AD8105D 0x0A84624D 0x0A5FA9AA 0x09FDEEAB 0x0AC35E44
					if *(*int32)(unsafe.Pointer(ebp18[0].addr)) == -2 {
						// 0x0B074A8B
						// ebp18 := ebp18[0].size + v09FFF8E0imageBase
						// 0x09FD69FF
					} else {
						// 0x0A55C059 0x0A0019A8 0x0AD3B5FE
						// ebp28 := ebp18[0].addr + v09FFF8E0imageBase
						// ...
					}
				}
			} // for loop 0x09FD69FF

			// 0x0A7473D5
			// pop esi
		}()
		// 0x0AF7F54A 0x0AF89454 0x0A4ECB31 0x0AF85A8A 0x0A56D0F9
		// pop esi ecx ebx fd edi
		// label3(0x0A338F26)
		// label2(0x0AC9CB58)
		// label1(0x09EB9FEC)

		// f09FD38F9 隐藏函数
	}()
	// ...
	// f00632F90, s9 f005AE611
	func() {
		// 很复杂
		// ...
		// 0x00634C73
		f005B0120()
		// ...
		// 0x006369DC, s9 0x005AF758
		// if f00552B60(...) {
		f005B0120() // f005B0120(ebp24, ebp2C)
		// } else {
		f00611C16() // f00611C16(ebp24, ebp2C, 1)
		// }
	}()
}

func f004E0E03handleState5() {
	// 0x0A32D9DB hook to hide f0A84CD1A, disable anti-temper with backup code
	var label1 uint32 = 0x0A5FBC25
	// push label1
	// push 0x0A8FCF09
	// ret

	// 0x0A8FCF09
	var label2 uint32 = 0x0B10C62E
	// push label2
	// push 0x0ABF8D43
	// ret

	// 0x0ABF8D43
	var label3 uint32 = 0x00E28EF7
	// push label3
	// push 0x09EADA48
	// ret

	// 0x09EADA48 0x0B072E16 0x0A0C2FAE 0x0A4DF4C8 0x0A559B77
	// push esi fd ecx ebx edx edi

	// 0x0A9FD42D
	// push 0x0A54F706
	// push 0x09FDEC49
	// ret

	// f09FDEC49
	func() {
		// push esi
		// push edi
		var ecx uint32
		if v0A56C1BA != v0AF8AF60 {
			// 0x0AA0736A
			v0A56C1BA = v0AF8AF60
			ecx = *(*uint32)(unsafe.Pointer(v0AD73D05 ^ v0ABD3FEA - v0B07C274 + 0x09FDEC49))
		} else {
			// 0x0A83C16F
			ecx = v09DEA85C
		}
		// 0x0A335B92
		if ecx == 0 {
			// 0x0A38D3F5
			ecx = 1
		}
		// 0x09F846EB
		v09DEA85C = v09F87E15 * ecx % v0A339687
		if v09DEA85C > v0AC3A339 {
			// 0x0AD8083E
			label1 = v0ABB631Flabel1 // [esp+v0B07545C*4+12] = v0ABB631Flabel1
			label2 = v0AC36B4Dlabel2 // [esp+v0A38E67E*4+12] = v0AC36B4Dlabel2
			// pop edi
			label3 = v0AA2D054label3 // [esp+v0AF87A7B*4+8]  = v0AA2D054label3
			// pop esi
			return // 0x0A54F706
		}
		// 0x09FC4087
		label1 = v0ABB631Flabel1 // [esp+v0B07545C*4+12] = v0ABB631Flabel1
		// edx := v0AD324BA
		// push ebx
		// esicodes := v0A9FD787backupCode[:]
		ebxblocks := v0A56AF60blocks[:]
		if v0A56AF60blocks[0].addr == ^uintptr(0) {
			// 0x09FE78E9
			label2 = v0AC36B4Dlabel2 // [esp+v0A38E67E*4+16] = v0AC36B4Dlabel2
			// pop ebx edi
			v0A933991 = v09E72FAA    // 这个有什么意义?
			label3 = v0AA2D054label3 // [esp+v0AF87A7B*4+8]  = v0AA2D054label3
			// pop esi
			return // 0x0A54F706
		}
		// 0x0AFDA2E0 0x0AC396D1
		// push ebp
		if ebxblocks[0].addr == ^uintptr(1) {
			// 0x0A5FCFF7
		} else {
			// 0x0AA30007F
		}
	}()

	// 0x0A54F706
	// pop edi edx ebx ecx fd esi
	// label3(0x0AF7BD1B)
	// label2(0x09FC4127)
	// label1(0x0A84CD1A)

	// f0A84CD1A 隐藏函数
}
