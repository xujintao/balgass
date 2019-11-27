package main

var v0131A240 uint32
var v0131A26C bool
var v0131A26D bool
var v0131A270 int

var v08C88AB8 [100]uint8
var v08C88C69 bool
var v08C88C6A bool
var v08C88C74 uint32
var v086A3B94 func()
var v086A3BEC uint32

var v0638C6C4 uint32 = 0x4EB0E616

var v0114EE50 = [...]uint64{
	0x3FE3333340000000,
	0x3FC3333340000000,
	0x3F70624DE0000000,
}

var v01319F28 uint32 = 0x3F2DE762

// f004E46B3handleState2
var v09FB8736 uint32 = 0x2B2F4A3D
var v0A4E24C6 uint32 = 0x4849459E
var v09E035DB uint32 = 0x5433BC5D
var v0A43DD91 uint32 = 0x653CAB4B
var v0B06F40C uint32 = 0x08244423 // v0B06F40C = 0x310F1716 = v09E035DB ^ v0A43DD91
var v09FDFB22 uint32 = 0xFC5D207C
var v0AFD3C52 uint32 = 0xD577
var v09FB6D69 uint32 = 0xDEA7
var v0A441BD1 uint32 = 0x38
var v0AD2DD3A uint32 = 0xE4E2E401
var v0AF890C3 uint32 = 9
var v0B287022label1 uint32 = 0x0AD98B31
var v09E8DF92 uint32 = 8
var v0AAB7324label2 uint32 = 0x0AD2CAED
var v0A933705 uint32 = 7
var v0A935B85label3 uint32 = 0x0A43B47B
var v0ABFAB88blocks = [...]block{
	{0x000393EA, 0x0192},
	{0x00039581, 0x0100},
	{0x0003968F, 0x05},
	// ...
}
var v0A4E85ABimageBase uintptr = 0x00400000
var v0AA2E05D uint32 = 0x5AA7D865
var v09EB9D65 uint32 = 1

var v0A9360CB uint32 = 0x4B0FAFAC
var v0AFD365B uint32 = 0x5EC66E14 // v0AFD365B = v0A9360CB

type t7 struct {
	data [0x2A4]byte
}

var v08C86C50 = [10]t7{}

var v0114EC40 float64

var v012E234C = "data\\music\\mutheme.mp3"
var v012E239C = "data\\music\\main_theme.mp3"

// f004E17B9handleState4
var v0A443F74 uint32 = 0xE9497BA1
var v0AD3B896 uint32 = 0x13EFEDCE // v0AD3B896 = v0A443F74
var v0A890E43 uint32 = 0xC491
var v0A56EB4C uint32 = 0xE521
var v0A84A81D uint32 = 0x3E92EB75 // 0xEB
var v09FE1805 uint32 = 0x125      // v0A84A81D 与 v09FE1805比较

var v0ABE324Flabel1 uint32 = 0x09EBC65D
var v0A952A97 uint32 = 8
var v0A92FC07label2 uint32 = 0x09FDF5DF
var v0A920E49 uint32 = 7
var v0A32B0C2label3 uint32 = 0x0AFDEEEE
var v0A6039A7 uint32 = 6

var v0AFD3C89 uint32 = 0x9B20B37E
var v0AA30425 uint32 = 0x105DA588

var v0A88819B uint32 = 0xD266CC5A // key
var v0AF96824backupCode = [...]uint32{
	0x37664E4F,
	0x36FE23B9,
}

var v0A9F69B2blocks = [...]block{
	{0x000393EA, 0x192}, // f004393EA
	{0x00039581, 0x100},
	{0xFFFFFFFF, 0xFFFFFFFF},
}
var v09FE37F3imageBase uintptr = 0x00400000

// f004DF0D5handleState5
var v0ABD75F3 uint32 = 0xE55020DE
var v0A3268A5 uint32 = 0xC01C390C

var v0A0C2FD2 uint32 = 0x05BE6BE3
var v0AF8FD45 uint32 = 0x34B17CF5
var v0AA09787 uint32 = 0x0824448B // v0AA09787 = v0A0C2FD2 ^ v0AF8FD45

var v0AD6FFBE uint32 = 0x90B9
var v0ABB4649 uint32 = 0x975D
var v0AAB36E2 uint32 = 0x531BDC38 // v0AD6FFBE * v0AAB36E2 % v0ABB4649
var v09EAE2F8 uint32 = 0x26

var v0AF781CF uint32 = 0x35C75981
var v0A33BBC9blocks = [...]block{}
var v0A5FCB8Bcodes = [...]uint32{}
var v0A84B964 uint32 = 1

var v09FFF8E0imageBase uintptr = 0x00400000

var v0A9FB214 uint32 = 7
var v0A443580label1 uint32 = 0x09EB9FEC
var v0AFDFD32 uint32 = 6
var v09FB8B32label2 uint32 = 0x0AC9CB58
var v0A849070 uint32 = 5
var v0A32C145label3 uint32 = 0x0A338F26

// f004E0E03handleState5
// 确保第一次不覆盖
var v0A56C1BA uint32 = 0x865EB72E
var v0AF8AF60 uint32 = 0xDCD68D02
var v0AD73D05 uintptr = 0xD2435CB5
var v0ABD3FEA uintptr = 0xD36C2725
var v0B07C274 uintptr = 0x09FDEC49 // v0AD73D05 ^ v0ABD3FEA - v0B07C274 + 0x09FDEC49

// 后面在一个确定的时机进行覆盖
var v0A339687 uint32 = 0xB2F9
var v09F87E15 uint32 = 0x8E25
var v09DEA85C uint32 = 0x76237DA5 //
var v0AC3A339 uint32 = 0xE5

var v0A9FD787backupCode = [...]uint32{}
var v0A56AF60blocks = [...]block{
	{0x000393EA, 0x192},
	{0x00039581, 0x100},
	{0xFFFFFFFF, 0xFFFFFFFF},
}

var v09E72FAA uint32 = 0x03B1C032
var v0A933991 uint32 = 0x922C3447 // v0A933991 = v09E72FAA

var v0B07545C uint32 = 8
var v0ABB631Flabel1 uint32 = 0x0A84CD1A
var v0A38E67E uint32 = 7
var v0AC36B4Dlabel2 uint32 = 0x09FC4127
var v0AF87A7B uint32 = 6
var v0AA2D054label3 uint32 = 0x0AF7BD1B

// f00A0A5E1 shell
var v0AA08598 uint32 = 0xE2209428
var v0AD7A978 uint32 = 0xABF8FFC9
var v0A9D361D uint32 = 0xD2B99D3B
var v0A391C78 uint32 = 0x52B79D4F // 0x7FFE0014 = v0A391C78 - v0A9D361D
var v0AD82FDD uint32 = 0x9EA1
var v0A338B93 uint32 = 0xE003
var v0A562D19 uint32 = 0x5AFBC57D
var v09F876D2 uint32 = 0x4334

var v0AD93A11 uint32 = 0x1A65E9C0
var v0AD93A0D = [...]uint32{0x078202AD}

var v0AD574B8 uint32 = 8
var v0A0519E0label1 uint32 = 0x09E2729C
