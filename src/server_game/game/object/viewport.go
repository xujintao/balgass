package object

import "github.com/xujintao/balgass/src/server_game/game/math2"

var (
	FrustrumX [MaxArrayFrustrum]int
	FrustrumY [MaxArrayFrustrum]int
)

func InitFrustrum() {
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

func (obj *object) addViewport(tobj *object) {
	if tobj.Class == 523 ||
		tobj.Class == 603 {
		return
	}
	for _, vp := range obj.viewports {
		if vp.number == tobj.index {
			return
		}
	}
	for _, vp := range obj.viewports {
		if vp.state == 0 {
			vp.state = 1
			vp.number = tobj.index
			vp.type_ = int(tobj.Type)
			obj.viewportNum++
			break
		}
	}
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
	start := 0
	// create viewport
	switch obj.Type {
	case ObjectTypePlayer:
		start = 0 // 玩家能看到所有对象
	case ObjectTypeMonster, ObjectTypeNPC:
		start = obj.objectManager.maxMonsterCount // 怪物看不见怪物
	}
	for _, tobj := range obj.objectManager.objects[start:] {
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
		obj.addViewport(tobj)
	}
}

func (obj *object) destoryViewport() {
	if obj.ConnectState != ConnectStatePlaying {
		return
	}
	// remove viewport
	for i, vp := range obj.viewports {
		if vp.state == 0 {
			continue
		}
		tobj := obj.objectManager.objects[vp.number]
		remove := false
		if tobj == nil {
			remove = true
		} else {
			if tobj.ConnectState < ConnectStatePlaying ||
				tobj.index == obj.index ||
				(tobj.State != 1 && tobj.State != 2) ||
				tobj.MapNumber != obj.MapNumber {
				remove = true
			}
			if !obj.checkViewport(tobj.X, tobj.Y) {
				remove = true
			}
		}
		if remove {
			obj.viewports[i].state = 0
			obj.viewports[i].number = -1
			obj.viewportNum--
		}
	}
}

func (obj *object) processViewport() {
	obj.destoryViewport()
	obj.createViewport()
	if obj.State == 1 {
		obj.State = 2
	}
}
