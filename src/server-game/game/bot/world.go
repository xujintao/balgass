package bot

import (
	"time"

	"github.com/xujintao/balgass/src/server-game/game/model"
)

type Actor struct {
	Index      int
	Class      int
	MapNumber  int
	X          int
	Y          int
	TX         int
	TY         int
	Dir        int
	HP         int
	MaxHP      int
	Alive      bool
	Attackable bool
}

func (a Actor) position() Position {
	return Position{X: a.X, Y: a.Y}
}

type Snapshot struct {
	Now     time.Time
	Self    Actor
	Players []Actor
	Objects []Actor
	Moving  bool
}

func (s Snapshot) blockers(except int) map[Position]struct{} {
	blockers := make(map[Position]struct{}, len(s.Players)+len(s.Objects))
	for _, actor := range s.Players {
		if actor.Alive && actor.Index != except {
			blockers[actor.position()] = struct{}{}
		}
	}
	for _, actor := range s.Objects {
		if actor.Alive && actor.Index != except {
			blockers[actor.position()] = struct{}{}
		}
	}
	return blockers
}

type movement struct {
	path     []Position
	next     int
	nextStep time.Time
}

type world struct {
	resources *resources
	self      Actor
	players   map[int]Actor
	objects   map[int]Actor
	movement  movement
	playing   bool
}

func newWorld(resources *resources) *world {
	return &world{
		resources: resources,
		self:      Actor{Index: -1},
		players:   make(map[int]Actor),
		objects:   make(map[int]Actor),
	}
}

func (w *world) handle(msg any) {
	switch msg := msg.(type) {
	case *model.MsgCreateViewportPlayerReply:
		for _, player := range msg.Players {
			actor := Actor{
				Index:     player.Index,
				Class:     player.Class,
				MapNumber: w.self.MapNumber,
				X:         player.X,
				Y:         player.Y,
				TX:        player.TX,
				TY:        player.TY,
				Dir:       player.Dir,
				HP:        player.HP,
				MaxHP:     player.MaxHP,
				Alive:     player.HP > 0,
			}
			if actor.Index == w.self.Index {
				w.mergeSelf(actor)
				continue
			}
			w.players[actor.Index] = actor
		}
	case *model.MsgCreateViewportMonsterReply:
		for _, monster := range msg.Monsters {
			w.objects[monster.Index] = Actor{
				Index:      monster.Index,
				Class:      monster.Class,
				MapNumber:  w.self.MapNumber,
				X:          monster.X,
				Y:          monster.Y,
				TX:         monster.TX,
				TY:         monster.TY,
				Dir:        monster.Dir,
				HP:         monster.HP,
				MaxHP:      monster.MaxHP,
				Alive:      monster.HP > 0,
				Attackable: w.resources.attackable(monster.Class),
			}
		}
	case *model.MsgDestroyViewportObjectReply:
		for _, actor := range msg.Objects {
			delete(w.players, actor.Index)
			delete(w.objects, actor.Index)
		}
	case *model.MsgMoveReply:
		w.handleMoveReply(msg)
	case *model.MsgSetPositionReply:
		w.handleSetPositionReply(msg)
	case *model.MsgAttackHPReply:
		w.handleAttackHPReply(msg)
	case *model.MsgAttackDieReply:
		w.markDead(msg.Target)
	case *model.MsgTeleportReply:
		w.clearSight()
		w.setSelfPosition(msg.MapNumber, msg.X, msg.Y, msg.Dir)
	case *model.MsgReloadCharacterReply:
		w.clearSight()
		w.setSelfPosition(msg.MapNumber, msg.X, msg.Y, msg.Dir)
		w.self.HP = msg.HP
		w.self.Alive = msg.HP > 0
	case *model.MsgHPReply:
		if msg.Position == -1 {
			w.self.HP = msg.HP
			if msg.HP <= 0 {
				w.markDead(w.self.Index)
			}
		}
	}
}

func (w *world) setCharacter(msg *model.MsgLoadCharacterReply) {
	w.setSelfPosition(msg.MapNumber, msg.X, msg.Y, msg.Dir)
	w.self.HP = msg.HP
	w.self.MaxHP = msg.MaxHP
	w.self.Alive = true
	w.playing = true
}

func (w *world) setSelfIndex(index int) {
	w.self.Index = index
}

func (w *world) mergeSelf(actor Actor) {
	w.self.Class = actor.Class
	w.self.HP = actor.HP
	w.self.MaxHP = actor.MaxHP
	if actor.HP > 0 {
		w.self.Alive = true
	}
}

func (w *world) setSelfPosition(mapNumber, x, y, dir int) {
	if w.self.MapNumber != mapNumber {
		w.clearSight()
	}
	w.self.MapNumber = mapNumber
	w.self.X = x
	w.self.Y = y
	w.self.TX = x
	w.self.TY = y
	w.self.Dir = dir
	w.movement = movement{}
}

func (w *world) clearSight() {
	clear(w.players)
	clear(w.objects)
	w.movement = movement{}
}

func (w *world) handleMoveReply(msg *model.MsgMoveReply) {
	if msg.Number == w.self.Index {
		w.self.TX = msg.X
		w.self.TY = msg.Y
		w.self.Dir = msg.Dir >> 4
		return
	}
	actor, ok := w.players[msg.Number]
	if ok {
		actor.X, actor.Y = msg.X, msg.Y
		actor.TX, actor.TY = msg.X, msg.Y
		actor.Dir = msg.Dir >> 4
		w.players[msg.Number] = actor
		return
	}
	actor, ok = w.objects[msg.Number]
	if ok {
		actor.X, actor.Y = msg.X, msg.Y
		actor.TX, actor.TY = msg.X, msg.Y
		actor.Dir = msg.Dir >> 4
		w.objects[msg.Number] = actor
	}
}

func (w *world) handleSetPositionReply(msg *model.MsgSetPositionReply) {
	if msg.Number == w.self.Index {
		w.setSelfPosition(w.self.MapNumber, msg.X, msg.Y, w.self.Dir)
		return
	}
	actor, ok := w.players[msg.Number]
	if ok {
		actor.X, actor.Y = msg.X, msg.Y
		actor.TX, actor.TY = msg.X, msg.Y
		w.players[msg.Number] = actor
		return
	}
	actor, ok = w.objects[msg.Number]
	if ok {
		actor.X, actor.Y = msg.X, msg.Y
		actor.TX, actor.TY = msg.X, msg.Y
		w.objects[msg.Number] = actor
	}
}

func (w *world) handleAttackHPReply(msg *model.MsgAttackHPReply) {
	if msg.Target == w.self.Index {
		w.self.HP = msg.HP
		w.self.MaxHP = msg.MaxHP
		if msg.HP <= 0 {
			w.markDead(msg.Target)
		}
		return
	}
	actor, ok := w.objects[msg.Target]
	if !ok {
		return
	}
	actor.HP = msg.HP
	actor.MaxHP = msg.MaxHP
	actor.Alive = msg.HP > 0
	w.objects[msg.Target] = actor
}

func (w *world) markDead(index int) {
	if index == w.self.Index {
		w.self.Alive = false
		w.movement = movement{}
		return
	}
	actor, ok := w.objects[index]
	if ok {
		actor.Alive = false
		actor.HP = 0
		w.objects[index] = actor
	}
}

func (w *world) startMove(path []Position, now time.Time) {
	if len(path) == 0 {
		return
	}
	w.movement = movement{
		path:     append([]Position(nil), path...),
		nextStep: now.Add(stepDuration(w.self.position(), path[0])),
	}
	end := path[len(path)-1]
	w.self.TX, w.self.TY = end.X, end.Y
}

func (w *world) advance(now time.Time) {
	for len(w.movement.path) > 0 &&
		w.movement.next < len(w.movement.path) &&
		!now.Before(w.movement.nextStep) {
		next := w.movement.path[w.movement.next]
		w.self.Dir = calcDir(w.self.position(), next)
		w.self.X, w.self.Y = next.X, next.Y
		w.movement.next++
		if w.movement.next >= len(w.movement.path) {
			w.movement = movement{}
			w.self.TX, w.self.TY = w.self.X, w.self.Y
			return
		}
		w.movement.nextStep = w.movement.nextStep.Add(stepDuration(next, w.movement.path[w.movement.next]))
	}
}

func (w *world) moving() bool {
	return len(w.movement.path) > 0
}

func (w *world) snapshot(now time.Time) Snapshot {
	snapshot := Snapshot{
		Now:    now,
		Self:   w.self,
		Moving: w.moving(),
	}
	for _, actor := range w.players {
		snapshot.Players = append(snapshot.Players, actor)
	}
	for _, actor := range w.objects {
		snapshot.Objects = append(snapshot.Objects, actor)
	}
	return snapshot
}

func stepDuration(from, to Position) time.Duration {
	if from.X != to.X && from.Y != to.Y {
		return 520 * time.Millisecond
	}
	return 400 * time.Millisecond
}
