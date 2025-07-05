package object

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/math2"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

var (
	FrustumX [MaxArrayFrustum]int
	FrustumY [MaxArrayFrustum]int
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
	var frustum [4][3]float32
	for i := 0; i < 4; i++ {
		frustum[i] = math2.VectorRotate(p[i], matrix)
		FrustumX[i] = int(frustum[i][0] * 0.01)
		FrustumY[i] = int(frustum[i][1] * 0.01)
	}
}

func (obj *Object) CreateFrustum() {
	for i := 0; i < MaxArrayFrustum; i++ {
		obj.FrustumX[i] = FrustumX[i] + obj.X
		obj.FrustumY[i] = FrustumY[i] + obj.Y
	}
}

func (obj *Object) initViewport() {
	for i := range obj.Viewports {
		obj.Viewports[i] = &Viewport{Number: -1}
	}
}

func (obj *Object) checkViewport(x, y int) bool {
	if x < obj.X-15 ||
		x > obj.X+15 ||
		y < obj.Y-15 ||
		y > obj.Y+15 {
		return false
	}
	for i, j := 0, 3; i < MaxArrayFrustum; j, i = i, i+1 {
		frustum := (obj.FrustumX[i]-x)*(obj.FrustumY[j]-y) -
			(obj.FrustumX[j]-x)*(obj.FrustumY[i]-y)
		if frustum < 0 {
			return false
		}
	}
	return true
}

func (obj *Object) addViewport(index, type_ int) bool {
	// if tobj.Class == 523 ||
	// 	tobj.Class == 603 {
	// 	return false
	// }
	for _, vp := range obj.Viewports {
		if vp.Number == index {
			return false
		}
	}
	for _, vp := range obj.Viewports {
		if vp.State == 0 {
			vp.State = 1
			vp.Number = index
			vp.Type = type_
			obj.ViewportNum++
			return true
		}
	}
	return false
}

func (obj *Object) clearViewport() {
	for i := range obj.Viewports {
		obj.Viewports[i].State = 0
		obj.Viewports[i].Number = -1
	}
	obj.ViewportNum = 0
}

func (obj *Object) createViewport() {
	if obj.ConnectState != ConnectStatePlaying {
		return
	}

	// create viewport item
	var viewportItemReply model.MsgCreateViewportItemReply
	if obj.Type == ObjectTypePlayer {
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
			obj.Push(&viewportItemReply)
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
			tobj.Index == obj.Index ||
			(tobj.State != 1 && tobj.State != 2) ||
			tobj.MapNumber != obj.MapNumber {
			continue
		}
		if !obj.checkViewport(tobj.X, tobj.Y) {
			continue
		}
		ok := obj.addViewport(tobj.Index, int(tobj.Type))
		if ok && obj.Type == ObjectTypePlayer {
			switch tobj.Type {
			case ObjectTypePlayer:
				p := model.CreateViewportPlayer{
					Index:                  tobj.Index,
					X:                      tobj.X,
					Y:                      tobj.Y,
					Class:                  tobj.Class,
					ChangeUp:               tobj.GetChangeUp(),
					Inventory:              [9]*item.Item(tobj.GetInventory().Items[:9]),
					Name:                   tobj.Name,
					TX:                     tobj.TX,
					TY:                     tobj.TY,
					Dir:                    tobj.Dir,
					PKLevel:                tobj.GetPKLevel(),
					PentagramMainAttribute: tobj.PentagramAttributePattern,
					Level:                  tobj.Level,
					MaxHP:                  tobj.MaxHP,
					HP:                     tobj.HP,
					ServerCode:             0,
				}
				viewportPlayerReply.Players = append(viewportPlayerReply.Players, &p)
			case ObjectTypeMonster, ObjectTypeNPC:
				m := model.CreateViewportMonster{
					Index:                  tobj.Index,
					Class:                  tobj.Class,
					X:                      tobj.X,
					Y:                      tobj.Y,
					TX:                     tobj.TX,
					TY:                     tobj.TY,
					Dir:                    tobj.Dir,
					PentagramMainAttribute: tobj.PentagramMainAttribute,
					Level:                  tobj.Level,
					MaxHP:                  tobj.MaxHP,
					HP:                     tobj.HP,
				}
				viewportMonsterReply.Monsters = append(viewportMonsterReply.Monsters, &m)
			}
		}
	}
	if len(viewportPlayerReply.Players) > 0 {
		obj.Push(&viewportPlayerReply)
	}
	if len(viewportMonsterReply.Monsters) > 0 {
		obj.Push(&viewportMonsterReply)
	}

	// debug create player viewport
	if conf.ServerEnv.Debug && obj.Type == ObjectTypePlayer {
		if len(viewportItemReply.Items) > 0 {
			var itemLine strings.Builder
			for _, it := range viewportItemReply.Items {
				itemLine.WriteString(fmt.Sprintf("(%d,%s)", it.Index, it.Item.Annotation))
			}
			s := itemLine.String()
			slog.Debug("viewport add", "items", s)
			obj.PushSystemMsg(fmt.Sprintf("viewport add [items]%s", s))
		}
		if len(viewportPlayerReply.Players) > 0 {
			var playerLine strings.Builder
			for _, p := range viewportPlayerReply.Players {
				playerLine.WriteString(fmt.Sprintf("(%d,%s)", p.Index, p.Name))
			}
			s := playerLine.String()
			slog.Debug("viewport add", "players", s)
			obj.PushSystemMsg(fmt.Sprintf("viewport add [players]%s", s))
		}
		if len(viewportMonsterReply.Monsters) > 0 {
			var monsterLine strings.Builder
			for _, m := range viewportMonsterReply.Monsters {
				monsterLine.WriteString(fmt.Sprintf("(%d)", m.Index))
			}
			s := monsterLine.String()
			slog.Debug("viewport add", "monsters", s)
			obj.PushSystemMsg(fmt.Sprintf("viewport add [monsters]%s", s))
		}
		if len(viewportItemReply.Items) > 0 ||
			len(viewportPlayerReply.Players) > 0 ||
			len(viewportMonsterReply.Monsters) > 0 {
			var line strings.Builder
			for _, vp := range obj.Viewports {
				if vp.State != 0 {
					line.WriteString(fmt.Sprintf("(%d)", vp.Number))
				}
			}
			s := line.String()
			slog.Debug("viewport", "all", s)
			obj.PushSystemMsg(fmt.Sprintf("viewport [all]%s", s))
		}
	}
}

func (obj *Object) destroyViewport() {
	if obj.ConnectState != ConnectStatePlaying {
		return
	}
	// remove viewport
	var viewportObjectReply model.MsgDestroyViewportObjectReply
	var viewportItemReply model.MsgDestroyViewportItemReply
	for i, vp := range obj.Viewports {
		if vp.State == 0 {
			continue
		}
		remove := false
		// notify := false
		switch vp.Type {
		case 5:
			maps.MapManager.MapItem(obj.MapNumber, vp.Number, func(item *item.Item, index, x, y int) {
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
			tobj := ObjectManager.objects[vp.Number]
			if tobj == nil {
				remove = true
			} else {
				if tobj.ConnectState < ConnectStatePlaying ||
					tobj.Index == obj.Index ||
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
				d := model.DestroyViewport{Index: vp.Number}
				switch vp.Type {
				case 5:
					viewportItemReply.Items = append(viewportItemReply.Items, &d)
				default:
					viewportObjectReply.Objects = append(viewportObjectReply.Objects, &d)
				}

			}
			obj.Viewports[i].State = 0
			obj.Viewports[i].Number = -1
			obj.ViewportNum--
		}
	}
	if len(viewportItemReply.Items) > 0 {
		obj.Push(&viewportItemReply)
	}
	if len(viewportObjectReply.Objects) > 0 {
		obj.Push(&viewportObjectReply)
	}

	// debug
	if conf.ServerEnv.Debug && obj.Type == ObjectTypePlayer {
		if len(viewportItemReply.Items) > 0 {
			var itemLine strings.Builder
			for _, it := range viewportItemReply.Items {
				itemLine.WriteString(fmt.Sprintf("(%d)", it.Index))
			}
			s := itemLine.String()
			slog.Debug("viewport remove", "items", s)
			obj.PushSystemMsg(fmt.Sprintf("viewport remove [items]%s", s))
		}
		if len(viewportObjectReply.Objects) > 0 {
			var objectLine strings.Builder
			for _, obj := range viewportObjectReply.Objects {
				objectLine.WriteString(fmt.Sprintf("(%d)", obj.Index))
			}
			s := objectLine.String()
			slog.Debug("viewport remove", "objects", s)
			obj.PushSystemMsg(fmt.Sprintf("viewport remove [objects]%s", s))
		}
		if len(viewportItemReply.Items) > 0 ||
			len(viewportObjectReply.Objects) > 0 {
			var line strings.Builder
			for _, vp := range obj.Viewports {
				if vp.State != 0 {
					line.WriteString(fmt.Sprintf("(%d)", vp.Number))
				}
			}
			s := line.String()
			slog.Debug("viewport", "all", s)
			obj.PushSystemMsg(fmt.Sprintf("viewport [all]%s", s))
		}
	}
}

func (obj *Object) ForEachViewportObject(f func(tobj *Object)) {
	for _, vp := range obj.Viewports {
		if vp.State == 0 {
			continue
		}
		if vp.Number < 0 {
			continue
		}
		if vp.Type == 5 {
			continue
		}
		tobj := ObjectManager.objects[vp.Number]
		if tobj == nil {
			continue
		}
		f(tobj)
	}
}

func (obj *Object) processViewport() {
	obj.destroyViewport()
	obj.createViewport()
	if obj.State == 1 {
		obj.State = 2
	}
}

// PushViewport Push msg to self and viewport
func (obj *Object) PushViewport(msg any) {
	obj.Push(msg)
	for _, vp := range obj.Viewports {
		if vp.State == 0 {
			continue
		}
		if vp.Number < 0 {
			continue
		}
		if vp.Type == 5 {
			continue
		}
		tobj := ObjectManager.objects[vp.Number]
		if tobj == nil {
			continue
		}
		tobj.Push(msg)
	}
}
