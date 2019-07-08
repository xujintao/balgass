package main

// 0x00B3,BAE4
func readconf(path string, buf1 []uint8) {
	// 1024+4字节的局部变量

	var p10 *uint32 // 指向参数
	var buf [1024]uint8
	buf[0] = 0

	_00DE8100(buf[1:], 0, 0x3FF) // call 0x00DE,8100，buf清零操作

}

// 0x00B6,2CF0
func getenc1(path string, key []uint8) {}

// 0x00B6,2D30
func getdec2(path string, key []uint8) {}
