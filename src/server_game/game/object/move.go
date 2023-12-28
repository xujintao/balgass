package object

import (
	"log"
	"math"
	"time"

	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
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
		obj.pathCount == 0 {
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
		log.Printf("process300ms object move check [index]%d [class]%d [map]%d [position](%d,%d)",
			obj.Index, obj.Class, obj.MapNumber, x, y)
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
	// if obj.index == 6 && obj.PathMoving {
	// 	fmt.Println(obj.X, obj.Y)
	// }
	obj.Dir = dir
	obj.CreateFrustrum()
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
	if n < 1 || n > 14 {
		return
	}
	for i := range msg.Path {
		obj.pathX[i] = msg.Path[i].X
		obj.pathY[i] = msg.Path[i].Y
		obj.pathDir[i] = msg.Dirs[i]
	}
	obj.pathCount = n
	obj.pathCur = 0
	obj.PathMoving = true
	maps.MapManager.ClearMapAttrStand(obj.MapNumber, obj.X, obj.Y)
	obj.TX = msg.Path[n-1].X
	obj.TY = msg.Path[n-1].Y
	maps.MapManager.SetMapAttrStand(obj.MapNumber, obj.TX, obj.TY)
	// if obj.Name == "asdf" {
	// 	fmt.Printf("(%d,%d)->(%d,%d)\n", obj.X, obj.Y, obj.TX, obj.TY)
	// }

	msgRelpy := model.MsgMoveReply{
		Number: obj.Index,
		X:      obj.TX,
		Y:      obj.TY,
		Dir:    obj.Dir << 4,
	}
	obj.destroyViewport()
	obj.createViewport()
	obj.PushViewport(&msgRelpy)
}

func (obj *Object) SetPosition(msg *model.MsgSetPosition) {
	obj.X = msg.X
	obj.Y = msg.Y
	obj.TX = msg.X
	obj.TY = msg.Y
	reply := model.MsgSetPositionReply{
		Number: obj.Index,
		X:      msg.X,
		Y:      msg.Y,
	}
	obj.CreateFrustrum()
	obj.destroyViewport()
	obj.createViewport()
	obj.PushViewport(&reply)
}
