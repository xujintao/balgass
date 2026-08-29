package formula

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/xujintao/balgass/src/server-game/conf"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	f.init()
}

var f formula

type formula struct {
	CalcCharacter    *lua.LState
	StatSpec         *lua.LState
	ItemCalc         *lua.LState
	RegularSkillCalc *lua.LState
	MasterSkillPoint *lua.LState
	ExpCalc          *lua.LState
}

func (f *formula) init() {
	load := func(file string) *lua.LState {
		file = path.Join(conf.PathCommon, "Scripts", file)
		slog.Info("load lua", "file", file)
		l := lua.NewState()
		if err := l.DoFile(file); err != nil {
			slog.Error("l.DoFile", "err", err)
			os.Exit(1)
		}
		return l
	}
	f.CalcCharacter = load("Character/CalcCharacter.lua")
	f.StatSpec = load("Specialization/StatSpec.lua")
	f.ItemCalc = load("Misc/ItemCalc.lua")
	f.RegularSkillCalc = load("Skills/RegularSkillCalc.lua")
	f.MasterSkillPoint = load("Skills/MasterSkillPoint.lua")
	f.ExpCalc = load("Misc/ExpCalc.lua")
}

func call(ls *lua.LState, method string, sig string, args ...any) error {
	sigs := strings.Split(sig, ">")
	in, out := "", ""
	in = sigs[0]
	if len(sigs) > 1 {
		out = sigs[1]
	}
	nIn := len(in)
	nRet := len(out)

	// prepare lua state args
	lvArgs := make([]lua.LValue, nIn)
	for i, v := range in {
		switch v {
		case 'i':
			lvArgs[i] = lua.LNumber(args[i].(int))
		case 'd':
			lvArgs[i] = lua.LNumber(args[i].(float64))
		}
	}

	// call
	err := ls.CallByParam(lua.P{
		Fn:      ls.GetGlobal(method),
		NRet:    nRet,
		Protect: true,
	}, lvArgs...)
	if err != nil {
		slog.Error("formula CallByParam", "err", err, "method", method)
		return err
	}

	// returned value
	defer ls.Pop(nRet)
	for i, v := range out {
		lv := ls.Get(0 - nRet + i)
		switch v {
		case 'i':
			ln, ok := lv.(lua.LNumber)
			if !ok {
				err := fmt.Errorf("formula %s returned %s, want number", method, lv.Type())
				slog.Error("formula CallByParam returned value i", "err", err)
				return err
			}
			r := args[nIn+i].(*int)
			*r = int(ln)
		case 'd':
			ln, ok := lv.(lua.LNumber)
			if !ok {
				err := fmt.Errorf("formula %s returned %s, want number", method, lv.Type())
				slog.Error("formula CallByParam returned value d", "err", err)
				return err
			}
			r := args[nIn+i].(*float64)
			*r = float64(ln)
		}
	}
	return nil
}
