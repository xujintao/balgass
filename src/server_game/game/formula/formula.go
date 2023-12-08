package formula

import (
	"log"
	"path"
	"strings"

	"github.com/xujintao/balgass/src/server_game/conf"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	f.init()
}

var f formula

type formula struct {
	CalcCharacter *lua.LState
	StatSpec      *lua.LState
}

func (f *formula) init() {
	load := func(file string) *lua.LState {
		file = path.Join(conf.PathCommon, "Scripts", file)
		log.Printf("Load %s", file)
		l := lua.NewState()
		if err := l.DoFile(file); err != nil {
			log.Fatalln(err)
		}
		return l
	}
	f.CalcCharacter = load("Character/CalcCharacter.lua")
	f.StatSpec = load("Specialization/StatSpec.lua")
}

func call(ls *lua.LState, method string, sig string, args ...any) {
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
			lvArgs[i] = lua.LNumber(args[i].(float32))
		}
	}

	// call
	err := ls.CallByParam(lua.P{
		Fn:      ls.GetGlobal(method),
		NRet:    nRet,
		Protect: true,
	}, lvArgs...)
	if err != nil {
		log.Printf("formula CallByParam failed [method]%s\n", method)
		return
	}

	// returned value
	defer ls.Pop(nRet)
	for i, v := range out {
		lv := ls.Get(0 - nRet + i)
		switch v {
		case 'i':
			ln, ok := lv.(lua.LNumber)
			if !ok {
				log.Printf("formula CallByParam returned value i invalid [method]%s\n", method)
				return
			}
			r := args[nIn+i].(*int)
			*r = int(ln)
		case 'd':
			ln, ok := lv.(lua.LNumber)
			if !ok {
				log.Printf("formula CallByParam returned value d invalid [method]%s\n", method)
				return
			}
			r := args[nIn+i].(*float32)
			*r = float32(ln)
		}
	}
}
