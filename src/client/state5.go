package main

import (
	"encoding/binary"
	"unsafe"
)

var v012E3EC8mapNumber int

// s9 f005A774E
func f005B0120() {
	// ...
	// 0x005B0208 inline send, s9 f005B67D9
	func() {
		var ebp2918 pb
		ebp2918.f00439178init()
		ebp2918.f0043922CwriteHead(0xC1, 0xD7) // send position, hook D7->D4
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
	case 2: // code:30, talk
		// 0x00618C52
		// v012E3140 = 12
		ebp69D8 := v01308D04objectPool.f00A38D5BgetObject(v012E31B0) // 0x12DABC9C
		if ebp69D8.m438 != 4 {
			v012E31B0 = -1
			return
		}
		ebp69DC := v086105ECpanel.m10Clevel       // 0x190
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
			// f00A49798ui().mF0.f00A700C1(1)
		} else {
			// f00A49798ui().mF0.f00A700C1(0)
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
		if v012E3EC8mapNumber >= 0x62 && v012E3EC8mapNumber <= 0x63 {
			// ... inline
		} else if v012E3EC8mapNumber < 0x2D || v012E3EC8mapNumber > 0x32 {
			// 0x00628420 inline send
			var ebp1896C pb // 0x0014F870
			ebp1896C.f00439178init()
			// ebp4 := 51
			ebp1896C.f0043922CwriteHead(0xC1, 0x30) // talk
			var ids [2]uint8
			binary.BigEndian.PutUint16(ids[:], ebp69D8.m5Eid) // bigendian 0x0001
			ebp1896C.f00439298writeBuf(ids[:], 2, true)
			ebp1896C.f004393EAsend(true, false)
			// 0x00629345
		}
	case 3: // code:D9, attack
		// 0x00611C9B, s9 0x005A539E
		// ...
		ebp20 := v01308D04objectPool.f00A38D5BgetObject(v012E3200) // 0x12DC9B08
		if ebp20 == nil {
			return
		}
		// ...
		var ebp14 uint8 // f00DE76F6round()
		// ...
		// 0x00611FF5 inline send, s9
		var ebp14A4 pb // 0x00166D38
		ebp14A4.f00439178init()
		// ebp4 := 2
		ebp14A4.f0043922CwriteHead(0xC1, 0xD9) // normal attack, hook D9->11
		var ids [2]uint8
		binary.BigEndian.PutUint16(ids[:], ebp20.m5Eid) // bigendian 0x166F
		ebp14A4.f00439298writeBuf(ids[:], 2, true)
		ebp14A4.f004397B1writeUint8(0x78)  // AttackAction
		ebp14A4.f004397B1writeUint8(ebp14) // DirDis
		ebp14A4.f004393EAsend(false, false)
	case 4:
		// 0x0062C00B
	case 5: // code:DA, position get
		// 0x0061362F
		// ...
		// 0x00613847
		// f00A0A4BE()
		func() {
			// 0x00A0A547
			// f00A0A4E1()
			func() {
				// 0x00A0B3C5
				// f00A15E91(), s9 f00936DF7
				func() {
					// 0x00A16082
					var ebp14A4 pb
					ebp14A4.f00439178init()
					// ebp4 := 0
					ebp14A4.f0043922CwriteHead(0xC1, 0xDA) // positionGet DA, s9 0x15
					// ebp14A4.f004397B1writeUint8(ebp12)       // 0x23
					// ebp14A4.f004397B1writeUint8(ebp11)       // 0x2E
					ebp14A4.f004393EAsend(false, false)
					// ebp4 := -1
					ebp14A4.f004391CF()
				}()
			}()
		}()
	case 6:
		// 0x006314F8
	}
}

// s9 what?
func f007E4FC4(x1, x2, x3, x4 int, x5 uint16, x6 bool) bool {
	// 0x09FE4C03 0x007E4FE1 0x09FE4C03 0x0A4400D7, s9 0x0AA9C932
	var ebp24 []uint32
	// ebp2C := 5
	ebp4 := 0
	ebp34 := 0
	ebp28skillID := v086105ECpanel.m178skillIDs[x1]
	var ebp30 int // ebp30 := f0090E94C().f0090EA51().f008F6A00(ebp28)
	ebpD := 0
	if v012E3EC8mapNumber == 0x1E {
		// 0x007E5035
	}
	// 0x007E5076 0x09EB65A0
	var ebp3C int
	for ebp3C < 0x190 {
		// 0x007E5083 0x09FD45F8
		ebp48obj := v01308D04objectPool.f00A38D5BgetObject(ebp3C)
		// 陨石: 13040D54(053B),13041410(053C),13041ACC(053D),13042188(053E),13042844(053F),13042F00(0540),130435BC(0541)
		// 死神骑士: 13043C78(068B),13044334(0691),130449F0(0692),130450AC(0694),13045E24(0695),13047FD0(068D)
		// 死神戈登: 13045768(06A2),13046B9C(06CE),13047258(06D1)
		// 恶魔: 130464E0(06CD),13047914(06CB)
		// 福袋: 1304868C(1545)
		// 月兔: 13048D48(160D),13049404(161E)
		ebp50 := &ebp48obj.m410
		if ebp50.m04 && ebp50.m0E && !ebp48obj.m24 && ebp48obj != v0805BBACobjectself { // 1304A838(2E82)
			// 渲染？
			// 0x007E50EC 0x0AC3A648
			// ...
			// 0x007E512E
			// f00DEF7A4()
			// 0x09F8DFA8
			// ...
			var ah uint8
			if ah == 0 {
				// 0x0AC34BE5
				if ebp50.m28 != 2 {
					// 0x0A9FF67F
					if ebp50.m28 == 1 {
						// 0x007E51BB
					} else {
						// 0x0AF11C9D 0x0A43A6DA
						if ebp48obj.m5Eid != x5 {
							// 0x09EAEC52
							if ebpD == 0 {
								// 0x007E51BB
							} else {
								// 0x007E5182
							}
						} else {
							// 0x007E5182
						}
					}
				} else {
					// 0x007E5182 0x0AD88F1A
					if ebp30 == 8 {
						// 0x007E518E 0x0A8899E7
						if ebp48obj.m8 != 0 {
							// 0x0AD3ED3C
							ebp48obj.m2B = 10
						} else {
							// 0x007E519E
						}
					} else {
						// 0x0AD9813E
						if ebp30 == 9 {
							// 0x007E518E
						} else {
							// 0x007E519E 0x0AA0E2E6
							ebp24[ebp34] = uint32(ebp48obj.m5Eid)
							ebp34++
							if ebp34 >= 5 {
								break // 0x0A916873 0x007E51C0
							} else {
								// 0x007E51BB
							}
						}
					}
				}
			} else {
				// 0x007E51BB 0x09FCD104 0x007E507C
			}
		}
		// 0x007E50EA 0x007E506F 0x0A32430C 0x007E507C 0x0A83F696
		ebp3C++
	}
	// 0x007E51C0 0x0A84040B
	if ebp34 <= 0 {
		return false // 0x007E9327
	}
	// 0x0A38F1B3
	if ebp30 == 0x4E {
		// 0x09FC6B67 // 0x007E51E5
		// ebp51 := x6
		// ebp2990.f00406FC0stdstring(&v0114ECB0)
		ebp4 = 0
		// v0805BBACobjectself.m38name
	}
	// 0x007E7271
	if ebp28skillID != 0 && ebp28skillID != 234 {
		// 0x09E27707, s9 0x0AA33867
		// 0x007E7293 0x0A5D51ED 0x007E72AE 0x0AD57EA6 0x0AD57EA6 0x09FB45A3 0x007E72C5 0x0A05ED3A 0x09FBE557
		// ebp29CC.f00406FC0stdstring(&v0114ECB0)
		ebp4 = 4
		// ebp29E8.f00406FC0stdstring(v0805BBACobjectself.m38name[:])
		ebp4 = 5
		var ebp29AD bool // ebp29AD := f004CE0BDstrstr(&ebp29E8, &ebp29CC)
		// 0x09FE37E6 0x0A55B22D 0x0A889689
		ebp4 = 4
		// 0x0A8FB3C7
		ebp4 = -1
		// 0x007E72FB
		// ebp29CC.f00407AC0(1, 0)
		// 0x0A4ECFA3
		if ebp29AD == false {
			// 0x007E9323 0x007E932E
		}
		// 0x0A436F59, s9 0x0AFE744C
		var ebp2964 pb
		ebp2964.f00439178init()
		// 0x0A0C438C
		ebp4 = 6
		ebp2964.f0043922CwriteHead(0xC1, 0x1D)
		println(ebp4)

	}
	// 0x007E9323 0x007E932E 0x0A9FEC0D
	al := true
	// 0x007E9329 0x0AD971F6
	return al
}

// s9 f004E077A
func f004DF0D5handleState5() {
	// 很复杂
	// ...
	// 0x004E0A71
	f00A49798ui().f00A4DC94fresh(v012E2340state)
	// ...
	//...
	// 0x004E0B99: f005B5D6C
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

		// 0x09FD38F9: true logic
	}()
	// ...
	// 0x004E0DA2: f00632F90
	// s9 0x?: s9 f005AE611
	func() {
		// 很复杂
		// ...
		// 0x00634C73, s9 0x005AF044
		f005B0120()
		// ...
		// 0x006369DC, s9 0x005AF758
		// if f00552B60(...) {
		f005B0120() // f005B0120(ebp24, ebp2C)
		// } else {
		f00611C16() // f00611C16(ebp24, ebp2C, 1)
		// }
		// ...
		// 0x0063DCD8: f00A09A07
		// s9 0x005B0D93: s9 f0092D446
		f00A09635skill().f00A09A07()
		/*
			ebp10 := f00DE76C0roundf(v0805BBACobjectself.m528)/100<<8 | f00DE76C0roundf(v0805BBACobjectself.m524)/100
			if ebp10 < 0 {
				ebp10 = 0
			} else if ebp10 > 0xFFFF {
				ebp10 = 0xFFFF
			}
			v086A3B8C = v088EBC10[ebp10]
		*/
	}()
	// ...
	// 0x004E0DD1: f0086BA70
	// s9 0x004E0DD2: s9 f007879A0
	func() {
		// f0085A04E, s9 f00775F5E
		func() {
			// 0x0AC33981:
			// s9 0x0A3F0B8:
			// ...
			f007E4FC4(1, 1, 1, 1, 1, true)
			// ...
		}()
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
