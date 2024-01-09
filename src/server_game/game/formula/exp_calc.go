package formula

func SetExpTable_Normal(level int, exp *int) {
	call(f.ExpCalc, "SetExpTable_Normal", "i>i", level, exp)
}

func SetExpTable_Master(level, base int, exp *int) {
	call(f.ExpCalc, "SetExpTable_Master", "ii>i", level, base, exp)
}
