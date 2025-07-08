package object

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

func (obj *Object) CalcDistance(tobj *Object) int {
	x := obj.X - tobj.X
	y := obj.Y - tobj.Y
	if x == 0 && y == 0 {
		return 0
	}
	return int(math.Sqrt(float64(x*x + y*y)))
}

func (obj *Object) processMove() {
	if obj.ConnectState < ConnectStatePlaying ||
		!obj.Live ||
		obj.State != 2 ||
		!obj.PathMoving {
		return
	}
	moveTime := 400
	if obj.delayLevel != 0 {
		moveTime += 300
	}
	if obj.pathDir[obj.pathCur]%2 == 0 {
		moveTime = int(float64(moveTime) * 1.3)
	}
	now := time.Now()
	if now.Before(obj.pathTime.Add(time.Duration(moveTime) * time.Millisecond)) {
		return
	}
	obj.pathTime = now
	x := obj.pathX[obj.pathCur]
	y := obj.pathY[obj.pathCur]
	dir := obj.pathDir[obj.pathCur]
	attr := maps.MapManager.GetMapAttr(obj.MapNumber, x, y)
	if attr&4 != 0 && attr&8 != 0 {
		slog.Warn("processMove maps.MapManager.GetMapAttr",
			"index", obj.Index, "map", obj.MapNumber, "position", fmt.Sprintf("(%d,%d)", x, y))
		for i := 0; i < len(obj.pathDir); i++ {
			obj.pathX[i] = 0
			obj.pathY[i] = 0
			obj.pathDir[i] = 0
		}
		obj.pathCount = 0
		obj.pathCur = 0
		obj.PathMoving = false
		return
	}
	obj.X = x
	obj.Y = y
	obj.Dir = dir
	obj.CreateFrustum()
	obj.pathCur++
	if obj.pathCur >= obj.pathCount {
		for i := 0; i < len(obj.pathDir); i++ {
			obj.pathX[i] = 0
			obj.pathY[i] = 0
			obj.pathDir[i] = 0
		}
		obj.pathCount = 0
		obj.pathCur = 0
		obj.PathMoving = false
	}
}

func (obj *Object) Move(msg *model.MsgMove) {
	n := len(msg.Path)
	if n > 15 {
		slog.Warn("Move object check",
			"index", obj.Index, "name", obj.Name, "map", obj.MapNumber, "path_count", n)
		return
	}
	// debug
	if conf.ServerEnv.Debug {
		if msg.X != obj.X || msg.Y != obj.Y {
			slog.Debug("Move object check",
				"index", obj.Index, "name", obj.Name, "map", obj.MapNumber,
				"client_position", fmt.Sprintf("(%d,%d)", msg.X, msg.Y),
				"server_position", fmt.Sprintf("(%d,%d)", obj.X, obj.Y),
			)
		}
	}
	// set move state machine
	obj.X = msg.X
	obj.Y = msg.Y
	obj.Dir = msg.Dir
	for i := 0; i < n; i++ {
		obj.pathX[i] = msg.Path[i].X
		obj.pathY[i] = msg.Path[i].Y
		obj.pathDir[i] = msg.Dirs[i]
	}
	obj.pathCount = n
	obj.pathCur = 0
	if n == 0 {
		obj.PathMoving = false
		obj.TX = msg.X
		obj.TY = msg.Y
	} else {
		obj.PathMoving = true
		maps.MapManager.ClearMapAttrStand(obj.MapNumber, obj.TX, obj.TY)
		obj.TX = msg.Path[n-1].X
		obj.TY = msg.Path[n-1].Y
		maps.MapManager.ClearMapAttrStand(obj.MapNumber, obj.X, obj.Y)
		maps.MapManager.SetMapAttrStand(obj.MapNumber, obj.TX, obj.TY)
	}
	// reply
	msgRelpy := model.MsgMoveReply{
		Number: obj.Index,
		X:      obj.TX,
		Y:      obj.TY,
		Dir:    obj.Dir << 4,
	}
	obj.PushViewport(&msgRelpy)
}

func (obj *Object) LoadMiniMap() {
	maps.MiniManager.ForEachMapNpc(obj.MapNumber, func(id, display, x, y int, name string) {
		reply := model.MsgMiniMapReply{
			ID:          id,
			IsNpc:       1,
			DisplayType: display,
			Type:        0,
			X:           x,
			Y:           y,
			Name:        name,
		}
		obj.Push(&reply)
	})
	maps.MiniManager.ForEachMapEntrance(obj.MapNumber, func(id, display, x, y int, name string) {
		reply := model.MsgMiniMapReply{
			ID:          id,
			IsNpc:       0,
			DisplayType: 1,
			Type:        0,
			X:           x,
			Y:           y,
			Name:        name,
		}
		obj.Push(&reply)
	})
}

func (obj *Object) gateMove(gateNumber int) bool {
	success := false
	maps.GateMoveManager.Move(gateNumber, func(mapNumber, x, y, dir int) {
		slog.Debug("teleport",
			"object", obj.Name,
			"from", fmt.Sprintf("[%d](%d,%d)", obj.MapNumber, obj.X, obj.Y),
			"to", fmt.Sprintf("[%d](%d,%d)", mapNumber, x, y))
		reply := model.MsgTeleportReply{
			GateNumber: gateNumber,
			MapNumber:  mapNumber,
			X:          x,
			Y:          y,
			Dir:        dir,
		}
		obj.Push(&reply)
		maps.MapManager.ClearMapAttrStand(obj.MapNumber, obj.TX, obj.TY)
		if obj.MapNumber != mapNumber {
			obj.MapNumber = mapNumber
			obj.LoadMiniMap()
		}
		obj.X, obj.Y = x, y
		obj.TX, obj.TY = x, y
		obj.Dir = dir
		maps.MapManager.SetMapAttrStand(obj.MapNumber, obj.TX, obj.TY)
		obj.CreateFrustum()
		success = true

	})
	return success
}

func (obj *Object) Teleport(msg *model.MsgTeleport) {
	obj.PathMoving = false
	obj.gateMove(msg.GateNumber)
}

func (obj *Object) MapMove(msg *model.MsgMapMove) {
	obj.PathMoving = false
	reply := model.MsgMapMoveReply{
		Result: 0,
	}
	defer obj.Push(&reply)
	maps.MapMoveManager.Move(msg.MoveIndex, func(gateNumber, level, money int) {
		objMoney := obj.GetMoney()
		if obj.Level < level || objMoney < money {
			return
		}
		ok := obj.gateMove(gateNumber)
		if ok {
			objMoney -= money
			obj.SetMoney(objMoney)
			reply := model.MsgMoneyReply{
				Result: -2,
				Money:  objMoney,
			}
			obj.Push(&reply)
		}
	})
}

func (obj *Object) SetPosition(msg *model.MsgSetPosition) {
	// debug
	if conf.ServerEnv.Debug {
		slog.Debug("SetPosition",
			"index", obj.Index, "name", obj.Name, "map", obj.MapNumber,
			"from", fmt.Sprintf("(%d,%d)", obj.X, obj.Y),
			"to", fmt.Sprintf("(%d,%d)", msg.X, msg.Y),
		)
	}
	obj.PathMoving = false
	maps.MapManager.ClearMapAttrStand(obj.MapNumber, obj.TX, obj.TY)
	obj.X, obj.Y = msg.X, msg.Y
	obj.TX, obj.TY = msg.X, msg.Y
	maps.MapManager.SetMapAttrStand(obj.MapNumber, obj.TX, obj.TY)
	reply := model.MsgSetPositionReply{
		Number: obj.Index,
		X:      msg.X,
		Y:      msg.Y,
	}
	obj.CreateFrustum()
	obj.PushViewport(&reply)
}

func (obj *Object) Knockback(tobj *Object) {
	dir := tobj.Dir
	if rand.Intn(3) == 0 {
		// reverse dir
		dir = (tobj.Dir + 4) % 8
	}
	dirPot := maps.Dirs[dir]
	x := tobj.X + dirPot.X
	y := tobj.Y + dirPot.Y
	attr := maps.MapManager.GetMapAttr(tobj.MapNumber, x, y)
	if attr&1 != 0 &&
		attr&2 != 0 &&
		attr&4 != 0 &&
		attr&8 != 0 &&
		attr&16 != 0 {
		return
	}
	tobj.SetPosition(&model.MsgSetPosition{
		X: x,
		Y: y,
	})
}
