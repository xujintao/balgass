package main

// windows ABI
type findData struct {
	dwFileAttributes   int32
	ftCreationTime     [2]int32
	ftLastAccessTime   [2]int32
	ftLastWriteTIme    [2]int32
	nFileSizeHigh      int32
	nFileSizeLow       int32
	dwReserved0        int32
	dwReserved1        int32
	cFileName          [260]int8
	cAlternateFileName [14]int8
}

// DLL表
func FindFirstFile(name string, fd *findData) int32 { return -1 }
func FindClose(h int32)                             {}

func checkupdate() bool {
	var h int32      // [ebp-184]
	fd := findData{} // [epb-180]
	// iterater [epb -40]
	fileNames := [...]string{ // [ebp-3C]~[ebp-18]
		"mu._xe",
		"mumsg._ll",
		"wz_zp._ll",
		"message._tf",
		"mu.exe",
		"mumsg.dll",
		"wz_zp.dll",
		"message.wtf",
	}
	// [ebp-20]
	updates := [4]bool{false} // [ebp-14]~[ebp-8]
	update := false           // [ebp-4]

	for i, fileName := range fileNames[:4] {
		h = FindFirstFile(fileName, &fd)
		if -1 != h {
			update = true
			updates[i] = true
		}
		FindClose(h)
	}
	if !update {
		return true
	}

	// update
	// launch mu.exe
	return false
}
