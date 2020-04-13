package main

func main() {

}

// 0x00C50000 ~ 0x00C51FFF section PE   8k字节
// 0x00C52000 ~ 0x032CFFFF section "unknown1"
// 0x032D0000 ~ 0x0332BFFF section "rsrc"
// 0x0332C000 ~ 0x0332DFFF section "idata"
// 0x0332E000 ~ 0x03759FFF section "unknown2" 壳代码
// 0x0375A000 ~ 0x039D1FFF section "dealuagk" 壳代码
// 0x039D2000 ~ 0x039D3FFF section "crzdzqvw" 8K字节 壳代码
// 0x039D2000 ~ 0x039D2200 实际1k
// 0x039D4000 ~ 0x039D7FFF section "taggant" 16k字节
// 0x039D4000 ~ 0x039D5B37 实际8k不到

// OEP 第一次从系统模块到用户模块的入口
func f039D4000() {
	// ...
	//
}

// 系统断点 0x778A6C51 ntdll.dll 0x77801000 ~ 0x77923FFF
func f778A6C51() {
	// 0x77834731
	// f77841330()
	// f77837A70()
	// ntdll.LdrQueryImageFileKeyOption
	// f77856E00()
	// f77899780()
	// 0x7784369F 0x7784371E
	// ntdll.ZwContinue()
	// ntdll.RtlRaiseStatus()
}
