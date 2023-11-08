package object

import (
	"log"
	"math"
	"time"

	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

func (obj *object) calcDistance(tobj *object) int {
	x := obj.X - tobj.X
	y := obj.Y - tobj.Y
	if x == 0 && y == 0 {
		return 0
	}
	return int(math.Sqrt(float64(x*x + y*y)))
}

func (obj *object) processMove() {
	if obj.ConnectState < ConnectStatePlaying ||
		!obj.Live ||
		obj.State != 2 ||
		obj.pathCount == 0 {
		return
	}
	moveTime := obj.moveSpeed
	if obj.delayLevel != 0 {
		moveTime += 300
	}
	pathTime := time.Now().UnixMilli()
	if pathTime-obj.pathTime+1 < int64(moveTime) {
		return
	}
	obj.pathTime = pathTime
	x := obj.pathX[obj.pathCur]
	y := obj.pathY[obj.pathCur]
	dir := obj.pathDir[obj.pathCur]
	attr := maps.MapManager.GetMapAttr(obj.MapNumber, x, y)
	if attr&4 != 0 && attr&8 != 0 {
		log.Printf("process300ms object move check [index]%d [class]%d [map]%d [position](%d,%d)",
			obj.index, obj.Class, obj.MapNumber, x, y)
		for i := 0; i < len(obj.pathDir); i++ {
			obj.pathX[i] = 0
			obj.pathY[i] = 0
			obj.pathDir[i] = 0
		}
		obj.pathCount = 0
		obj.pathCur = 0
		obj.pathMoving = false
		return
	}
	obj.X = x
	obj.Y = y
	// if obj.index == 6 && obj.pathMoving {
	// 	fmt.Println(obj.X, obj.Y)
	// }
	obj.Dir = dir
	obj.createFrustrum()
	obj.pathCur++
	if obj.pathCur >= obj.pathCount {
		for i := 0; i < len(obj.pathDir); i++ {
			obj.pathX[i] = 0
			obj.pathY[i] = 0
			obj.pathDir[i] = 0
		}
		obj.pathCount = 0
		obj.pathCur = 0
		obj.pathMoving = false
	}
}

func (obj *object) Move(msg *model.MsgMove) {
	n := len(msg.Path)
	if n < 1 || n > 15 {
		return
	}
	for i := range msg.Path {
		obj.pathX[i] = msg.Path[i].X
		obj.pathY[i] = msg.Path[i].Y
		obj.pathDir[i] = msg.Dirs[i]
	}
	obj.pathCount = n
	obj.pathCur = 0
	obj.pathMoving = true
	maps.MapManager.ClearMapAttrStand(obj.MapNumber, obj.X, obj.Y)
	obj.TX = msg.Path[n-1].X
	obj.TY = msg.Path[n-1].Y
	maps.MapManager.SetMapAttrStand(obj.MapNumber, obj.TX, obj.TY)
	// if obj.Name == "asdf" {
	// 	fmt.Printf("(%d,%d)->(%d,%d)\n", obj.X, obj.Y, obj.TX, obj.TY)
	// }

	msgRelpy := model.MsgMoveReply{
		Number: obj.index,
		X:      obj.TX,
		Y:      obj.TY,
		Dir:    msg.Dirs[0] << 4,
	}
	obj.push(&msgRelpy)
	obj.pushViewport(&msgRelpy)
}
