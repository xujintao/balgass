package object

import (
	"testing"

	"github.com/xujintao/balgass/src/server-game/game/model"
)

func newViewportTestObject(index int, typ ObjectType) *Object {
	obj := &Object{
		Index: index,
		Type:  typ,
	}
	obj.Init()
	return obj
}

func withTestObjectManager(t *testing.T, objs ...*Object) {
	t.Helper()

	old := ObjectManager
	t.Cleanup(func() {
		ObjectManager = old
	})

	ObjectManager.objects = make([]*Object, 128)
	ObjectManager.maxObjectCount = len(ObjectManager.objects)
	ObjectManager.maxMonsterCount = 16
	ObjectManager.maxCallMonsterCount = 16
	ObjectManager.maxPlayerCount = ObjectManager.maxObjectCount - ObjectManager.maxMonsterCount - ObjectManager.maxCallMonsterCount
	ObjectManager.playerStartIndex = ObjectManager.maxMonsterCount + ObjectManager.maxCallMonsterCount

	for _, obj := range objs {
		if obj == nil {
			continue
		}
		ObjectManager.objects[obj.Index] = obj
	}
}

func findViewport(slots [MaxViewportNum]*Viewport, number int) (*Viewport, bool) {
	for _, vp := range slots {
		if vp != nil && vp.Number == number {
			return vp, true
		}
	}
	return nil, false
}

func assertViewportReset(t *testing.T, vp *Viewport) {
	t.Helper()
	if vp == nil {
		t.Fatal("viewport slot is nil")
	}
	if vp.State != 0 {
		t.Fatalf("State = %d, want 0", vp.State)
	}
	if vp.Number != -1 {
		t.Fatalf("Number = %d, want -1", vp.Number)
	}
	if vp.Type != int(ObjectTypeEmpty) {
		t.Fatalf("Type = %d, want %d", vp.Type, ObjectTypeEmpty)
	}
	if vp.Dis != 0 {
		t.Fatalf("Dis = %d, want 0", vp.Dis)
	}
}

func TestInitViewportAllocatesAndResetsSlots(t *testing.T) {
	obj := newViewportTestObject(1, ObjectTypePlayer)

	for i := 0; i < MaxViewportNum; i++ {
		assertViewportReset(t, obj.Viewports[i])
		assertViewportReset(t, obj.ViewportsPassive[i])
	}
}

func TestAddViewportItemCreatesItemEntry(t *testing.T) {
	obj := newViewportTestObject(1, ObjectTypePlayer)

	if ok := obj.addViewportItem(42); !ok {
		t.Fatal("addViewportItem returned false, want true")
	}

	vp, ok := findViewport(obj.Viewports, 42)
	if !ok {
		t.Fatal("item viewport entry was not created")
	}
	if vp.State != 1 {
		t.Fatalf("State = %d, want 1", vp.State)
	}
	if vp.Type != 5 {
		t.Fatalf("Type = %d, want 5", vp.Type)
	}
	if obj.ViewportsNum != 1 {
		t.Fatalf("ViewportsNum = %d, want 1", obj.ViewportsNum)
	}
	if obj.ViewportsPassiveNum != 0 {
		t.Fatalf("ViewportsPassiveNum = %d, want 0", obj.ViewportsPassiveNum)
	}
	for _, vpp := range obj.ViewportsPassive {
		if vpp.Number == 42 {
			t.Fatal("item viewport should not create passive entry")
		}
	}
}

func TestAddViewportObjectCreatesActiveAndPassiveEntries(t *testing.T) {
	viewer := newViewportTestObject(32, ObjectTypePlayer)
	target := newViewportTestObject(5, ObjectTypeMonster)
	withTestObjectManager(t, viewer, target)

	if ok := viewer.addViewportObject(target); !ok {
		t.Fatal("addViewportObject returned false, want true")
	}

	vp, ok := findViewport(viewer.Viewports, target.Index)
	if !ok {
		t.Fatal("active viewport entry was not created")
	}
	if vp.State != 1 {
		t.Fatalf("active State = %d, want 1", vp.State)
	}
	if vp.Type != int(ObjectTypeMonster) {
		t.Fatalf("active Type = %d, want %d", vp.Type, ObjectTypeMonster)
	}

	vpp, ok := findViewport(target.ViewportsPassive, viewer.Index)
	if !ok {
		t.Fatal("passive viewport entry was not created")
	}
	if vpp.State != 1 {
		t.Fatalf("passive State = %d, want 1", vpp.State)
	}
	if vpp.Type != int(ObjectTypePlayer) {
		t.Fatalf("passive Type = %d, want %d", vpp.Type, ObjectTypePlayer)
	}
	if viewer.ViewportsNum != 1 {
		t.Fatalf("viewer ViewportsNum = %d, want 1", viewer.ViewportsNum)
	}
	if target.ViewportsPassiveNum != 1 {
		t.Fatalf("target ViewportsPassiveNum = %d, want 1", target.ViewportsPassiveNum)
	}
}

func TestAddViewportObjectRollsBackWhenPassiveFull(t *testing.T) {
	viewer := newViewportTestObject(32, ObjectTypePlayer)
	target := newViewportTestObject(5, ObjectTypeMonster)
	withTestObjectManager(t, viewer, target)

	for i, vp := range target.ViewportsPassive {
		vp.State = 1
		vp.Number = i + 100
		vp.Type = int(ObjectTypePlayer)
		target.ViewportsPassiveNum++
	}

	if ok := viewer.addViewportObject(target); ok {
		t.Fatal("addViewportObject returned true, want false")
	}
	if _, ok := findViewport(viewer.Viewports, target.Index); ok {
		t.Fatal("active viewport entry was not rolled back")
	}
	if viewer.ViewportsNum != 0 {
		t.Fatalf("viewer ViewportsNum = %d, want 0", viewer.ViewportsNum)
	}
	if target.ViewportsPassiveNum != MaxViewportNum {
		t.Fatalf("target ViewportsPassiveNum = %d, want %d", target.ViewportsPassiveNum, MaxViewportNum)
	}
}

func TestRemoveViewportClearsReversePassiveEntry(t *testing.T) {
	viewer := newViewportTestObject(32, ObjectTypePlayer)
	target := newViewportTestObject(5, ObjectTypeMonster)
	withTestObjectManager(t, viewer, target)

	if ok := viewer.addViewportObject(target); !ok {
		t.Fatal("addViewportObject returned false, want true")
	}

	viewer.removeViewport(0)

	assertViewportReset(t, viewer.Viewports[0])
	if _, ok := findViewport(target.ViewportsPassive, viewer.Index); ok {
		t.Fatal("reverse passive viewport entry was not cleared")
	}
	if viewer.ViewportsNum != 0 {
		t.Fatalf("viewer ViewportsNum = %d, want 0", viewer.ViewportsNum)
	}
	if target.ViewportsPassiveNum != 0 {
		t.Fatalf("target ViewportsPassiveNum = %d, want 0", target.ViewportsPassiveNum)
	}
}

func TestClearViewportClearsActivePassiveAndReverseLinks(t *testing.T) {
	viewer := newViewportTestObject(32, ObjectTypePlayer)
	target := newViewportTestObject(5, ObjectTypeMonster)
	watcher := newViewportTestObject(40, ObjectTypePlayer)
	withTestObjectManager(t, viewer, target, watcher)

	if ok := viewer.addViewportItem(10); !ok {
		t.Fatal("addViewportItem returned false, want true")
	}
	if ok := viewer.addViewportObject(target); !ok {
		t.Fatal("addViewportObject returned false, want true")
	}
	viewer.Viewports[2].State = 1
	viewer.Viewports[2].Number = 99
	viewer.Viewports[2].Type = int(ObjectTypeMonster)
	viewer.ViewportsNum++

	viewer.ViewportsPassive[0].State = 1
	viewer.ViewportsPassive[0].Number = watcher.Index
	viewer.ViewportsPassive[0].Type = int(watcher.Type)
	viewer.ViewportsPassiveNum++

	viewer.clearViewport()

	for i := 0; i < MaxViewportNum; i++ {
		assertViewportReset(t, viewer.Viewports[i])
		assertViewportReset(t, viewer.ViewportsPassive[i])
	}
	if _, ok := findViewport(target.ViewportsPassive, viewer.Index); ok {
		t.Fatal("reverse passive viewport entry was not cleared")
	}
	if viewer.ViewportsNum != 0 {
		t.Fatalf("viewer ViewportsNum = %d, want 0", viewer.ViewportsNum)
	}
	if viewer.ViewportsPassiveNum != 0 {
		t.Fatalf("viewer ViewportsPassiveNum = %d, want 0", viewer.ViewportsPassiveNum)
	}
	if target.ViewportsPassiveNum != 0 {
		t.Fatalf("target ViewportsPassiveNum = %d, want 0", target.ViewportsPassiveNum)
	}
}

func TestCreateViewportExcludesHiddenMonster(t *testing.T) {
	viewerActor := &skillTestActor{}
	viewer := newViewportTestObject(32, ObjectTypePlayer)
	viewer.Objecter = viewerActor
	viewer.ConnectState = ConnectStatePlaying
	viewer.State = 1
	viewer.MapNumber = 0
	viewer.X, viewer.Y = 130, 130
	viewer.TX, viewer.TY = viewer.X, viewer.Y
	viewer.CreateFrustum()

	hidden := newViewportTestObject(5, ObjectTypeMonster)
	hidden.ConnectState = ConnectStatePlaying
	hidden.Live = true
	hidden.State = 1
	hidden.Class = 105
	hidden.Hidden = true
	hidden.MapNumber = viewer.MapNumber
	hidden.X, hidden.Y = viewer.X, viewer.Y
	hidden.TX, hidden.TY = hidden.X, hidden.Y

	withTestObjectManager(t, viewer, hidden)

	viewer.createViewport()

	if _, ok := findViewport(viewer.Viewports, hidden.Index); ok {
		t.Fatal("hidden monster was added to the player's viewport")
	}
	if _, ok := findViewport(hidden.ViewportsPassive, viewer.Index); ok {
		t.Fatal("hidden monster subscribed the player to its passive viewport")
	}
	for _, message := range viewerActor.messages {
		if _, ok := message.(*model.MsgCreateViewportMonsterReply); ok {
			t.Fatal("MsgCreateViewportMonsterReply was pushed for hidden monster")
		}
	}
}
