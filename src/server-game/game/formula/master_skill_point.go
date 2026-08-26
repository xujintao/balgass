package formula

import lua "github.com/yuin/gopher-lua"

type MasterSkillPointFormula struct {
	STID      int
	Method    string
	ValueType int
}

func MasterSkillPointFormulas() ([]MasterSkillPointFormula, error) {
	ls := f.MasterSkillPoint
	old := ls.GetGlobal("AddToMLSTable")
	defer ls.SetGlobal("AddToMLSTable", old)

	formulas := make([]MasterSkillPointFormula, 0, 39)
	ls.SetGlobal("AddToMLSTable", ls.NewFunction(func(ls *lua.LState) int {
		formulas = append(formulas, MasterSkillPointFormula{
			STID:      ls.CheckInt(1),
			Method:    ls.CheckString(2),
			ValueType: ls.CheckInt(3),
		})
		return 0
	}))
	if err := call(ls, "MLS_ValueInit", ">"); err != nil {
		return nil, err
	}
	return formulas, nil
}

func MasterSkillPointValue(method string, level int) (float64, error) {
	value := 0.0
	if err := call(f.MasterSkillPoint, method, "i>d", level, &value); err != nil {
		return 0, err
	}
	return value, nil
}
