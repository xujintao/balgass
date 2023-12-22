package object

import (
	"log"

	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/math2"
	"github.com/xujintao/balgass/src/server_game/game/model"
)

var (
	FrustrumX [MaxArrayFrustrum]int
	FrustrumY [MaxArrayFrustrum]int
)

func init() {
	var cameraViewFar float32 = 3200.0
	var cameraviewNear float32 = cameraViewFar * 0.19
	var cameraViewTarget float32 = cameraViewFar * 0.53
	var widthFar float32 = 1390.0
	var widthNear float32 = 750.0

	p := [4][3]float32{
		{-widthFar, cameraViewFar - cameraViewTarget, 0.0},
		{widthFar, cameraViewFar - cameraViewTarget, 0.0},
		{widthNear, cameraviewNear - cameraViewTarget, 0.0},
		{-widthNear, cameraviewNear - cameraViewTarget, 0.0},
	}
	angle := [3]float32{0.0, 0.0, 45.0}
	matrix := math2.Angle2Matrix(angle)
	var frustrum [4][3]float32
	for i := 0; i < 4; i++ {
		frustrum[i] = math2.VectorRotate(p[i], matrix)
		FrustrumX[i] = int(frustrum[i][0] * 0.01)
		FrustrumY[i] = int(frustrum[i][1] * 0.01)
	}
}

func (obj *object) createFrustrum() {
	for i := 0; i < MaxArrayFrustrum; i++ {
		obj.FrustrumX[i] = FrustrumX[i] + obj.X
		obj.FrustrumY[i] = FrustrumY[i] + obj.Y
	}
}

func (obj *object) initViewport() {
	for i := range obj.viewports {
		obj.viewports[i] = &viewport{number: -1}
	}
}

func (obj *object) checkViewport(x, y int) bool {
	if x < obj.X-15 ||
		x > obj.X+15 ||
		y < obj.Y-15 ||
		y > obj.Y+15 {
		return false
	}
	for i, j := 0, 3; i < MaxArrayFrustrum; j, i = i, i+1 {
		frustrum := (obj.FrustrumX[i]-x)*(obj.FrustrumY[j]-y) -
			(obj.FrustrumX[j]-x)*(obj.FrustrumY[i]-y)
		if frustrum < 0 {
			return false
		}
	}
	return true
}

func (obj *object) addViewport(index, type_ int) bool {
	// if tobj.Class == 523 ||
	// 	tobj.Class == 603 {
	// 	return false
	// }
	for _, vp := range obj.viewports {
		if vp.number == index {
			return false
		}
	}
	for _, vp := range obj.viewports {
		if vp.state == 0 {
			vp.state = 1
			vp.number = index
			vp.type_ = type_
			obj.viewportNum++
			return true
		}
	}
	return false
}

func (obj *object) clearViewport() {
	for i := range obj.viewports {
		obj.viewports[i].state = 0
		obj.viewports[i].number = -1
	}
	obj.viewportNum = 0
}

func (obj *object) createViewport() {
	if obj.ConnectState != ConnectStatePlaying {
		return
	}

	if obj.Type == ObjectTypePlayer {
		// create viewport item
		var viewportItemReply model.MsgCreateViewportItemReply
		maps.MapManager.MapEachItem(obj.MapNumber, func(item *item.Item, index, x, y int) {
			if !obj.checkViewport(x, y) {
				return
			}
			ok := obj.addViewport(index, 5)
			if ok {
				i := model.CreateViewportItem{
					Index: index,
					X:     x,
					Y:     y,
					Item:  item,
				}
				viewportItemReply.Items = append(viewportItemReply.Items, &i)
			}
		})
		if len(viewportItemReply.Items) > 0 {
			obj.push(&viewportItemReply)
		}
	}

	// create viewport object
	start := 0
	switch obj.Type {
	case ObjectTypePlayer:
		start = 0 // 玩家能看到所有对象
	case ObjectTypeMonster, ObjectTypeNPC:
		start = ObjectManager.maxMonsterCount // 怪物看不见怪物
	}
	var viewportPlayerReply model.MsgCreateViewportPlayerReply
	var viewportMonsterReply model.MsgCreateViewportMonsterReply
	for _, tobj := range ObjectManager.objects[start:] {
		if tobj == nil {
			continue
		}
		if tobj.ConnectState < ConnectStatePlaying ||
			tobj.index == obj.index ||
			(tobj.State != 1 && tobj.State != 2) ||
			tobj.MapNumber != obj.MapNumber {
			continue
		}
		if !obj.checkViewport(tobj.X, tobj.Y) {
			continue
		}
		ok := obj.addViewport(tobj.index, int(tobj.Type))
		if ok && obj.Type == ObjectTypePlayer {
			switch tobj.Type {
			case ObjectTypePlayer:
				p := model.CreateViewportPlayer{
					Index:                  tobj.index,
					X:                      tobj.X,
					Y:                      tobj.Y,
					Class:                  tobj.Class,
					ChangeUp:               tobj.GetChangeUp(),
					Inventory:              tobj.GetInventory(),
					Name:                   tobj.Name,
					TX:                     tobj.TX,
					TY:                     tobj.TY,
					Dir:                    tobj.Dir,
					PKLevel:                tobj.getPKLevel(),
					PentagramMainAttribute: tobj.pentagramAttributePattern,
					Level:                  tobj.Level,
					MaxHP:                  tobj.MaxHP,
					HP:                     tobj.HP,
					ServerCode:             0,
				}
				viewportPlayerReply.Players = append(viewportPlayerReply.Players, &p)
			case ObjectTypeMonster, ObjectTypeNPC:
				m := model.CreateViewportMonster{
					Index:                  tobj.index,
					Class:                  tobj.Class,
					X:                      tobj.X,
					Y:                      tobj.Y,
					TX:                     tobj.TX,
					TY:                     tobj.TY,
					Dir:                    tobj.Dir,
					PentagramMainAttribute: tobj.pentagramMainAttribute,
					Level:                  tobj.Level,
					MaxHP:                  tobj.MaxHP,
					HP:                     tobj.HP,
				}
				viewportMonsterReply.Monsters = append(viewportMonsterReply.Monsters, &m)
			}
		}
	}
	if len(viewportPlayerReply.Players) > 0 {
		obj.push(&viewportPlayerReply)
	}
	if len(viewportMonsterReply.Monsters) > 0 {
		obj.push(&viewportMonsterReply)
	}
}

func (obj *object) destroyViewport() {
	if obj.ConnectState != ConnectStatePlaying {
		return
	}
	// remove viewport
	var viewportObjectReply model.MsgDestroyViewportObjectReply
	var viewportItemReply model.MsgDestroyViewportItemReply
	for i, vp := range obj.viewports {
		if vp.state == 0 {
			continue
		}
		remove := false
		// notify := false
		switch vp.type_ {
		case 5:
			maps.MapManager.MapItem(obj.MapNumber, vp.number, func(item *item.Item, index, x, y int) {
				if item == nil {
					remove = true
					// notify = true
					return
				}
				if !obj.checkViewport(x, y) {
					remove = true
				}
			})
		default:
			tobj := ObjectManager.objects[vp.number]
			if tobj == nil {
				remove = true
			} else {
				if tobj.ConnectState < ConnectStatePlaying ||
					tobj.index == obj.index ||
					(tobj.State == 8) ||
					tobj.MapNumber != obj.MapNumber {
					remove = true
					// notify = true
				}
				if !obj.checkViewport(tobj.X, tobj.Y) {
					remove = true
				}
			}
		}

		if remove {
			// if obj.Type == ObjectTypePlayer && notify {
			if obj.Type == ObjectTypePlayer {
				d := model.DestroyViewport{Index: vp.number}
				switch vp.type_ {
				case 5:
					viewportItemReply.Items = append(viewportItemReply.Items, &d)
				default:
					viewportObjectReply.Objects = append(viewportObjectReply.Objects, &d)
				}

			}
			obj.viewports[i].state = 0
			obj.viewports[i].number = -1
			obj.viewportNum--
		}
	}
	if len(viewportObjectReply.Objects) > 0 {
		obj.push(&viewportObjectReply)
	}
	if len(viewportItemReply.Items) > 0 {
		obj.push(&viewportItemReply)
	}
}

func (obj *object) processViewport() {
	obj.destroyViewport()
	obj.createViewport()
	if obj.State == 1 {
		obj.State = 2
	}
}

// pushViewport push msg to self and viewport
func (obj *object) pushViewport(msg any) {
	obj.push(msg)
	for _, vp := range obj.viewports {
		if vp.state == 0 {
			continue
		}
		if vp.number < 0 {
			log.Printf("pushViewport warning [obj]%d [msg]%v -> [tobj]%d\n",
				obj.index, msg, vp.number)
			continue
		}
		tobj := ObjectManager.objects[vp.number]
		if tobj == nil {
			continue
		}
		tobj.push(msg)
	}
}
