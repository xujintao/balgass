package main

func _00B385D6(buf []uint8, size uint32, flag uint32) {
	// 20个局部变量

	var index uint32 = 0 //ebp-14h
	var key = [16]uint8{ //ebp-10h
		0x7C,
		0xBD,
		0x81,
		0x9F,

		0x3D,
		0x93,
		0xE2,
		0x56,

		0x2A,
		0x73,
		0xD2,
		0x3E,

		0xF2,
		0x83,
		0x95,
		0xBF,
	}

	for index < size {
		// buf = buf[index:]
		// el := uint32(buf[0])
		// el |= uint32(key[index+flag])
		// buf[0] = uint8(el)

		buf[index] ^= key[index+flag]
		index++
	}
}

func _00B3BAE4(p *t1000, log string) conf {
	// 1024+4字节的局部变量

	// 当返回值是个大结构时，c编译器会使用命名返回值，
	// 所以这里经过优化后在汇编层面其实是var c *conf，指向外面的结构
	// 在go里面就比较简单了，func类型返回值特征标写为*conf就行或者直接使用匿名调用
	var c *conf
	var buf [1024]uint8
	buf[0] = 0

	_00DE8100(buf[1:], 0, 0x3FF) // buf清零操作
	_00DF0805(buf[:], log, c)    // 把logconf字符串写到切片里，这个在golang里面很简单

	// 将buf编码后再写进文件
	// _00B38A8D
	// 它只要使用ecx传参数，我就认为作者使用的是匿名调用
	func(buf []uint8) {
		var ebp_4 uint32 // dwNumBytes
		var ebp_8 *t1000 = p
		if p.f04h == -1 {
			return
		}

		// 验证
		// _00DE7C00
		len := func(buf []uint8) uint32 {
			// 验证buf底层数组地址
			// 验证buf底层数组元素
			return uint32(len(buf))
		}(buf)

		// 写文件
		// _00B38A4D
		func(hFile int32, buf []uint8, len uint32, pdwNumBytes *uint32, pOverlapped uint32) {
			// 这个变量是push ecx得到的
			// 局部变量<=8个字节，c编译器使用push，性能比sub指令高
			var ebp_4 *t1000 = p

			// _00B38653 疑似编码函数
			func(buf []uint8, len uint32, flag uint32) uint32 {
				// 16个字节的局部变量
				var ebp_4 uint32
				var ebp_8 []uint8 = buf
				var ebp_C uint32 = len
				var ebp_10 uint32
				var eax uint32 = (0x10 - flag) & 0x8000000F
				if eax < 0 { // jns, Jump if not sign
					eax--
					eax |= 0xFFFFFFF0
					eax++
				}
				if eax < len {
					eax = (0x10 - flag) & 0x8000000F
					if eax < 0 { // jns, Jump if not sign
						eax--
						eax |= 0xFFFFFFF0
						eax++
					}
					ebp_10 = eax
				} else {
					ebp_10 = ebp_C
				}

				// 编码1个字节？
				ebp_4 = ebp_10                // 1
				_00B385D6(ebp_8, ebp_4, flag) // (0x0018DF84, 1, 15)
				ebp_8 = ebp_8[ebp_4:]
				ebp_C -= ebp_4
				if ebp_C <= 0 { // jg, Jump if greater
					return ebp_4 + ebp_10
				}

				// 编码写16个字节
				for ebp_C >= 0x10 { // jl, Jump if less
					ebp_4 = 0x10
					_00B385D6(ebp_8, ebp_4, 0)
					ebp_8 = ebp_8[ebp_4:]
					ebp_C -= ebp_4
				}

				// 编码剩下的
				ebp_4 = ebp_C
				_00B385D6(ebp_8, ebp_4, 0)
				return ebp_4
			}(buf, len, ebp_4.f10Ch)

			WriteFile(hFile, buf, len, pdwNumBytes, pOverlapped)
		}(ebp_8.f04h, buf, len, &ebp_4, 0)

		//
		if ebp_4 == 0 {
			nRet := CloseHandle(ebp_8.f04h)
			// _00B38781
			func(path string) {
				// 重新设置ebp_8的hFile
			}(string(ebp_8.f08h[:]))
		}
		return
	}(buf[:])

	return *c
}

// 0x00B6,2CF0
func getenc1(path string, key []uint8) {}

// 0x00B6,2D30
func getdec2(path string, key []uint8) {}
