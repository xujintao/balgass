package cmd

import (
	"encoding/json"
	"log"
)

type level int

const (
	Player level = iota
	GM
	Admin
)

type cmd struct {
	id     int
	level  level
	name   string
	code   int
	handle func(int, []uint8)
	enc    bool
}

var FindCmd map[int]*cmd

func init() {
	for _, v := range cmds {
		if vv, ok := FindCmd[v.code]; ok {
			log.Printf("duplicated cmd code[%d] name[%s] with code[%d] name[%s]", v.code, v.name, vv.code, vv.name)
		}
		FindCmd[v.code] = v
	}
}

var cmds = [...]*cmd{
	{1, Player, "object_use_item", 0x26, ObjectUseItem, false},
}

type MsgObjectUseItem struct{}

func ObjectUseItem(id int, data []uint8) {
	msg := MsgObjectUseItem{}
	json.Unmarshal(data, &msg)
}
