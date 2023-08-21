package object

func (obj *object) initViewport() {
	for i := range obj.viewports {
		obj.viewports[i] = &viewport{number: -1}
	}
	for i := range obj.viewports2 {
		obj.viewports2[i] = &viewport{number: -1}
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
		frustrum := (obj.FrustrumX[i]-x)*(obj.FrustrumY[i]-y) -
			(obj.FrustrumX[j]-x)*(obj.FrustrumY[j]-y)
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
	// type_ := tobj.Type
	// index := tobj.index
	// k := int(type_)<<16 + index
	// if _, ok := obj.viewports[k]; !ok {
	// 	v := &viewport{
	// 		state:  1,
	// 		number: index,
	// 		type_:  int(type_),
	// 	}
	// 	obj.viewports[k] = v
	// }
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

func (obj *object) addViewport2(tobj *object) {
	if tobj.Class == 523 ||
		tobj.Class == 603 {
		return
	}
	for _, vp := range obj.viewports2 {
		if vp.state == 0 {
			vp.state = 1
			vp.number = tobj.index
			vp.type_ = int(tobj.Type)
			obj.viewportNum2++
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

	for i := range obj.viewports2 {
		obj.viewports2[i].state = 0
		obj.viewports2[i].number = -1
	}
	obj.viewportNum2 = 0
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
	for _, v := range obj.objectManager.objects[start:] {
		if v == nil {
			continue
		}
		tobj := v.getObject()
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
		tobj.addViewport2(obj)
	}
}

func (obj *object) destoryViewport() {
	if obj.ConnectState != ConnectStatePlaying {
		return
	}
	// remove viewport
	for i, vp := range obj.viewports {
		if vp.state != 1 && vp.state != 2 {
			continue
		}
		tnum := vp.number
		switch vp.type_ {
		case 5: // items
		default: // objects
			tobj := obj.objectManager.objects[tnum].getObject()
			if tobj == nil {
				obj.viewports[i].state = 3
			} else {
				if tobj.ConnectState < ConnectStatePlaying ||
					tobj.index == obj.index ||
					(tobj.State != 1 && tobj.State != 2) ||
					tobj.MapNumber != obj.MapNumber {
					obj.viewports[i].state = 3
				}
				if !obj.checkViewport(tobj.X, tobj.Y) {
					obj.viewports[i].state = 3
				}
			}
		}
	}
	for i, vp := range obj.viewports2 {
		if vp.state != 1 && vp.state != 2 {
			continue
		}
		tobj := obj.objectManager.objects[vp.number].getObject()
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
			obj.viewports2[i].state = 0
			obj.viewports2[i].number = -1
			obj.viewportNum2--
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
