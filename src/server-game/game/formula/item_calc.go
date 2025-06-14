package formula

func Wings_CalcIncAttack(code, level int, rate *int) {
	call(f.ItemCalc, "Wings_CalcIncAttack", "iii>i", 100, code, level, rate)
}

func Wings_CalcAbsorb(code, level int, rate *int) {
	call(f.ItemCalc, "Wings_CalcAbsorb", "iiid>i", 100, code, level, 0.0, rate)
}
