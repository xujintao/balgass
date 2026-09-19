package object

import (
	"testing"
	"time"

	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
)

type eventDeathActor struct {
	*skillTestActor
	drops, regenerated int
}

func (a *eventDeathActor) DieDropItem(*Object) { a.drops++ }
func (a *eventDeathActor) Regen()              { a.regenerated++ }

func TestEventMonsterDeathSettlementAndRegen(t *testing.T) {
	dead, base := newSkillTestObject(1, ObjectTypeMonster)
	killer, _ := newSkillTestObject(32, ObjectTypePlayer)
	withTestObjectManager(t, dead, killer)
	ObjectManager.monsterCount = 1
	a := &eventDeathActor{skillTestActor: base}
	dead.Objecter = a
	now := time.Now()
	dead.Live = false
	dead.State = 4
	dead.dieRegen = true
	dead.dieTime = now
	dead.NoRegen = true
	dead.MaxRegenTime = 100 * time.Second
	dead.AddDelayMsg(1, 0, 800, killer.Index)
	dead.AddDelayMsg(3, 0, 2000, killer.Index)
	// Even a late cleanup must wait for messages the object loop has not run.
	if ObjectManager.DeleteEventMonster(dead, now.Add(6*time.Second)) {
		t.Fatal("removed pending death")
	}
	for _, msg := range dead.msgs {
		switch msg.code {
		case 1:
			msg.time = time.Now().Add(-time.Hour)
		case 3:
			msg.time = time.Now().Add(time.Hour)
		}
	}
	dead.processDelayMsg()
	if a.drops != 1 {
		t.Fatal("drop not settled")
	}
	if ObjectManager.DeleteEventMonster(dead, now.Add(6*time.Second)) {
		t.Fatal("removed pending recovery")
	}
	for _, msg := range dead.msgs {
		if msg.code == 3 {
			msg.time = time.Now().Add(-time.Hour)
		}
	}
	dead.processDelayMsg()
	dead.processDelayMsg()
	if a.drops != 1 {
		t.Fatal("duplicate drop")
	}
	if ObjectManager.DeleteEventMonster(dead, now.Add(4*time.Second)) {
		t.Fatal("removed before corpse phase ended")
	}
	dead.dieTime = time.Now().Add(-time.Hour)
	dead.processRegen()
	if dead.Live || a.regenerated != 0 {
		t.Fatal("event monster respawned")
	}
	if !ObjectManager.DeleteEventMonster(dead, now.Add(101*time.Second)) || ObjectManager.monsterCount != 0 {
		t.Fatal("settled corpse not retired")
	}

	// The same regen path still revives ordinary monsters.
	ordinary, base := newSkillTestObject(2, ObjectTypeMonster)
	a = &eventDeathActor{skillTestActor: base}
	ordinary.Objecter = a
	ordinary.Live = false
	ordinary.dieRegen = true
	ordinary.dieTime = time.Now()
	ordinary.MaxRegenTime = 100 * time.Second
	ordinary.processRegen()
	if ordinary.Live {
		t.Fatal("ordinary monster regenerated early")
	}
	ordinary.dieTime = time.Now().Add(-time.Hour)
	ordinary.processRegen()
	if !ordinary.Live || a.regenerated != 1 {
		t.Fatal("ordinary regen changed")
	}
}

func TestDeleteEventMonsterCleansReferencesAndGuardsSlot(t *testing.T) {
	mon, base := newSkillTestObject(1, ObjectTypeMonster)
	viewer, viewerActor := newSkillTestObject(32, ObjectTypePlayer)
	withTestObjectManager(t, mon, viewer)
	ObjectManager.monsterCount = 1
	a := &eventDeathActor{skillTestActor: base}
	mon.Objecter = a
	mon.NoRegen = true
	mon.TX, mon.TY = findBasicAttackMapPosition(t, false)
	maps.MapManager.SetMapAttrStand(0, mon.TX, mon.TY)
	if !viewer.addViewportObject(mon) || !mon.addViewportObject(viewer) {
		t.Fatal("viewport setup")
	}
	viewer.TargetNumber = mon.Index
	viewer.AddDelayMsg(4, 0, 500, mon.Index)
	if !ObjectManager.DeleteEventMonster(mon, time.Now()) {
		t.Fatal("live cleanup failed")
	}
	if a.drops != 0 || maps.MapManager.GetMapAttr(0, mon.TX, mon.TY)&2 != 0 {
		t.Fatal("cleanup rewarded or leaked occupied tile")
	}
	if viewer.ViewportsNum != 0 || viewer.ViewportsPassiveNum != 0 || viewer.TargetNumber != -1 || viewer.msgs[0].code != -1 {
		t.Fatal("dangling references")
	}
	if countBasicAttackMessages[*model.MsgDestroyViewportObjectReply](viewerActor.messages) != 1 {
		t.Fatal("missing client despawn")
	}
	replacement, _ := newSkillTestObject(1, ObjectTypeMonster)
	ObjectManager.objects[1] = replacement
	ObjectManager.monsterCount = 1
	if !ObjectManager.DeleteEventMonster(mon, time.Now()) || ObjectManager.objects[1] != replacement || ObjectManager.monsterCount != 1 {
		t.Fatal("stale handle deleted replacement")
	}
	if ObjectManager.DeleteEventMonster(replacement, time.Now()) {
		t.Fatal("ordinary monster removed")
	}
}

func TestAddMonsterCapacityReturnsError(t *testing.T) {
	m := objectManager{maxMonsterCount: 1, objects: make([]*Object, 1), lastMonsterIndex: -1}
	created := 0
	create := func() *Object { created++; return &Object{} }
	if _, err := m.AddMonster(create); err != nil {
		t.Fatal(err)
	}
	if _, err := m.AddMonster(create); err == nil {
		t.Fatal("capacity allowed")
	}
	m.monsterCount = 0 // Exhausted slots must also return an error.
	if _, err := m.AddMonster(create); err == nil || created != 1 {
		t.Fatal("full slots invoked constructor")
	}
}
