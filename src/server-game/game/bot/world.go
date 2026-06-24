package bot

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

type Actor struct {
	Index      int
	Class      int
	ChangeUp   int
	Level      int
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

type WorldSnapshot struct {
	Phase           Phase
	Failure         string
	Self            Actor
	Players         []Actor
	Objects         []Actor
	Skills          []CombatSkill
	LearnedSkills   []int
	LearnSkills     []LearnSkill
	MP              int
	AG              int
	AttackSpeed     int
	MagicSpeed      int
	PositionVersion uint64
}

type CombatSkill struct {
	Index    int
	Damage   int
	Distance int
	MP       int
	AG       int
	Type     int
	Delay    int
}

type LearnSkill struct {
	Index    int
	Position int
}

type Phase int

const (
	PhaseDisconnected Phase = iota
	PhaseConnected
	PhaseLoggedIn
	PhasePlaying
	PhaseFailed
)

func isImplementedDamageSkill(index int) bool {
	switch index {
	case skill.SkillIndexPoison, // 1毒咒
		skill.SkillIndexMeteorite,         // 2陨石
		skill.SkillIndexLightning,         // 3掌心雷
		skill.SkillIndexFireBall,          // 4火球
		skill.SkillIndexFlame,             // 5火龙
		skill.SkillIndexIce,               // 7冰封
		skill.SkillIndexTwister,           // 8龙卷风
		skill.SkillIndexEvilSpirit,        // 9黑龙波
		skill.SkillIndexPowerWave,         // 11真空波
		skill.SkillIndexAquaBeam,          // 12极光
		skill.SkillIndexCometFall,         // 13爆炎
		skill.SkillIndexEnergyBall,        // 17能量球(初始)
		skill.SkillIndexFallingSlash,      // 19地裂斩(武器)
		skill.SkillIndexLunge,             // 20牙突刺(武器)
		skill.SkillIndexUppercut,          // 21升龙击(武器)
		skill.SkillIndexCyclone,           // 22旋风斩(武器)
		skill.SkillIndexSlash,             // 23天地十字剑(武器)
		skill.SkillIndexTripleShot,        // 24多重箭(武器)
		skill.SkillIndexDecay,             // 38单毒炎
		skill.SkillIndexIceStorm,          // 39暴风雪
		skill.SkillIndexTwistingSlash,     // 41霹雳回旋斩
		skill.SkillIndexRagefulBlow,       // 42雷霆裂闪
		skill.SkillIndexDeathStab,         // 43袭风刺
		skill.SkillIndexCrescentMoonSlash, // 44半月斩(攻城)
		skill.SkillIndexImpale,            // 47钻云枪
		skill.SkillIndexFireBreath,        // 49流星焰(彩云兽)
		skill.SkillIndexIceArrow,          // 51冰封箭
		skill.SkillIndexPenetration,       // 52穿透箭
		skill.SkillIndexFireSlash,         // 55玄月斩
		skill.SkillIndexPowerSlash,        // 56天雷闪(武器)
		skill.SkillIndexSpiralSlash,       // 57风舞回旋斩(攻城)
		skill.SkillIndexForce,             // 60冲击(初始)
		skill.SkillIndexFireBurst,         // 61星云火链
		skill.SkillIndexElectricSpike,     // 65圣极光
		skill.SkillIndexForceWave:         // 66冲击波
		return true
	default:
		return false
	}
}

func (s WorldSnapshot) blockers(except int) map[Position]struct{} {
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

const (
	skillTypeSiege    = -1
	skillTypePhysical = 0
	skillTypeMagic    = 1
)

type world struct {
	resources       *resources
	apis            map[reflect.Type]*api
	phase           Phase
	failure         string
	self            Actor
	players         map[int]Actor
	objects         map[int]Actor
	skills          skill.Skills
	inventory       []*item.Item
	mp              int
	ag              int
	attackSpeed     int
	magicSpeed      int
	strength        int
	dexterity       int
	vitality        int
	energy          int
	leadership      int
	skillsSet       bool
	itemsSet        bool
	characterSet    bool
	selfClassSet    bool
	positionVersion uint64
}

func newWorld(resources *resources) *world {
	return &world{
		resources: resources,
		apis:      worldAPIs,
		self:      Actor{Index: -1},
		players:   make(map[int]Actor),
		objects:   make(map[int]Actor),
		skills:    make(skill.Skills),
	}
}

type connectFailed struct {
	Err error
}

type api struct {
	msg    any
	handle string
}

var apis = [...]*api{
	{(*connectFailed)(nil), "HandleConnectFailed"},
	{(*model.MsgConnectReply)(nil), "HandleConnectReply"},
	{(*model.MsgLoginReply)(nil), "HandleLoginReply"},
	{(*model.MsgLoadCharacterReply)(nil), "HandleLoadCharacterReply"},
	{(*model.MsgItemListReply)(nil), "HandleItemListReply"},
	{(*model.MsgSkillListReply)(nil), "HandleSkillListReply"},
	{(*model.MsgSkillOneReply)(nil), "HandleSkillOneReply"},
	{(*model.MsgDeleteInventoryItemReply)(nil), "HandleDeleteInventoryItemReply"},
	{(*model.MsgPickItemReply)(nil), "HandlePickItemReply"},
	{(*model.MsgCreateViewportPlayerReply)(nil), "HandleCreateViewportPlayerReply"},
	{(*model.MsgCreateViewportMonsterReply)(nil), "HandleCreateViewportMonsterReply"},
	{(*model.MsgDestroyViewportObjectReply)(nil), "HandleDestroyViewportObjectReply"},
	{(*model.MsgMoveReply)(nil), "HandleMoveReply"},
	{(*model.MsgSetPositionReply)(nil), "HandleSetPositionReply"},
	{(*model.MsgAttackHPReply)(nil), "HandleAttackHPReply"},
	{(*model.MsgAttackDieReply)(nil), "HandleAttackDieReply"},
	{(*model.MsgTeleportReply)(nil), "HandleTeleportReply"},
	{(*model.MsgReloadCharacterReply)(nil), "HandleReloadCharacterReply"},
	{(*model.MsgHPReply)(nil), "HandleHPReply"},
	{(*model.MsgMPReply)(nil), "HandleMPReply"},
	{(*model.MsgAttackSpeedReply)(nil), "HandleAttackSpeedReply"},
	{(*model.MsgLevelUpReply)(nil), "HandleLevelUpReply"},
	{(*model.MsgAddLevelPointReply)(nil), "HandleAddLevelPointReply"},
}

var worldAPIs = newAPIs()

func newAPIs() map[reflect.Type]*api {
	table := make(map[reflect.Type]*api, len(apis))
	worldType := reflect.TypeOf((*world)(nil))
	for _, api := range apis {
		t := reflect.TypeOf(api.msg)
		if t.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("bot api %s message must be a pointer", api.handle))
		}
		if _, ok := table[t]; ok {
			panic(fmt.Sprintf("bot api message %s is duplicated", t))
		}
		if _, ok := worldType.MethodByName(api.handle); !ok {
			panic(fmt.Sprintf("bot world has no handler %s", api.handle))
		}
		table[t] = api
	}
	return table
}

func (w *world) Handle(msg any) bool {
	if msg == nil {
		return false
	}
	api, ok := w.apis[reflect.TypeOf(msg)]
	if !ok {
		return false
	}
	handler := reflect.ValueOf(w).MethodByName(api.handle)
	if !handler.IsValid() {
		slog.Error("bot world has no handler", "handle", api.handle)
		return false
	}
	handler.Call([]reflect.Value{reflect.ValueOf(msg)})
	return true
}

func (w *world) HandleConnectFailed(msg *connectFailed) {
	w.fail(fmt.Sprintf("connect failed: %v", msg.Err))
}

func (w *world) HandleConnectReply(msg *model.MsgConnectReply) {
	if w.phase != PhaseDisconnected {
		return
	}
	if msg.Result != 1 {
		w.fail(fmt.Sprintf("connect failed: result %d", msg.Result))
		return
	}
	w.self.Index = msg.ID
	w.phase = PhaseConnected
}

func (w *world) HandleLoginReply(msg *model.MsgLoginReply) {
	if w.phase != PhaseConnected {
		return
	}
	if msg.Result != 1 {
		w.fail(fmt.Sprintf("login failed: result %d", msg.Result))
		return
	}
	w.phase = PhaseLoggedIn
}

func (w *world) HandleLoadCharacterReply(msg *model.MsgLoadCharacterReply) {
	if w.phase != PhaseLoggedIn {
		return
	}
	w.setCharacter(msg)
	w.phase = PhasePlaying
}

func (w *world) HandleItemListReply(msg *model.MsgItemListReply) {
	w.inventory = append([]*item.Item(nil), msg.Items...)
	w.itemsSet = true
}

func (w *world) HandleSkillListReply(msg *model.MsgSkillListReply) {
	w.skills = make(skill.Skills, len(msg.Skills))
	for _, s := range msg.Skills {
		if s != nil {
			w.skills[s.Index] = s
		}
	}
	w.skillsSet = true
}

func (w *world) HandleSkillOneReply(msg *model.MsgSkillOneReply) {
	if msg.Skill == nil {
		return
	}
	if msg.Flag == -1 {
		delete(w.skills, msg.Skill.Index)
		return
	}
	w.skills[msg.Skill.Index] = msg.Skill
}

func (w *world) HandleDeleteInventoryItemReply(msg *model.MsgDeleteInventoryItemReply) {
	if msg.Position >= 0 && msg.Position < len(w.inventory) {
		w.inventory[msg.Position] = nil
	}
}

func (w *world) HandlePickItemReply(msg *model.MsgPickItemReply) {
	if msg.Result < 0 || msg.Item == nil {
		return
	}
	if msg.Result >= len(w.inventory) {
		items := make([]*item.Item, msg.Result+1)
		copy(items, w.inventory)
		w.inventory = items
	}
	w.inventory[msg.Result] = msg.Item
}

func (w *world) HandleMPReply(msg *model.MsgMPReply) {
	if msg.Position != -1 {
		return
	}
	w.mp = msg.MP
	w.ag = msg.AG
}

func (w *world) HandleAttackSpeedReply(msg *model.MsgAttackSpeedReply) {
	w.attackSpeed = msg.AttackSpeed
	w.magicSpeed = msg.MagicSpeed
}

func (w *world) HandleLevelUpReply(msg *model.MsgLevelUpReply) {
	w.self.Level = msg.Level
}

func (w *world) HandleAddLevelPointReply(msg *model.MsgAddLevelPointReply) {
	switch msg.Type {
	case 0x10:
		w.strength++
	case 0x11:
		w.dexterity++
	case 0x12:
		w.vitality++
	case 0x13:
		w.energy++
	case 0x14:
		w.leadership++
	}
}

func (w *world) HandleCreateViewportPlayerReply(msg *model.MsgCreateViewportPlayerReply) {
	for _, player := range msg.Players {
		actor := Actor{
			Index:     player.Index,
			Class:     player.Class,
			ChangeUp:  player.ChangeUp,
			Level:     player.Level,
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
}

func (w *world) HandleCreateViewportMonsterReply(msg *model.MsgCreateViewportMonsterReply) {
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
}

func (w *world) HandleDestroyViewportObjectReply(msg *model.MsgDestroyViewportObjectReply) {
	for _, actor := range msg.Objects {
		delete(w.players, actor.Index)
		delete(w.objects, actor.Index)
	}
}

func (w *world) HandleMoveReply(msg *model.MsgMoveReply) {
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

func (w *world) HandleSetPositionReply(msg *model.MsgSetPositionReply) {
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

func (w *world) HandleAttackHPReply(msg *model.MsgAttackHPReply) {
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

func (w *world) HandleAttackDieReply(msg *model.MsgAttackDieReply) {
	w.markDead(msg.Target)
}

func (w *world) HandleTeleportReply(msg *model.MsgTeleportReply) {
	w.clearSight()
	w.setSelfPosition(msg.MapNumber, msg.X, msg.Y, msg.Dir)
}

func (w *world) HandleReloadCharacterReply(msg *model.MsgReloadCharacterReply) {
	w.clearSight()
	w.setSelfPosition(msg.MapNumber, msg.X, msg.Y, msg.Dir)
	w.self.HP = msg.HP
	w.self.Alive = msg.HP > 0
	w.mp = msg.MP
	w.ag = msg.AG
}

func (w *world) HandleHPReply(msg *model.MsgHPReply) {
	if msg.Position != -1 {
		return
	}
	w.self.HP = msg.HP
	if msg.HP <= 0 {
		w.markDead(w.self.Index)
	}
}

func (w *world) setCharacter(msg *model.MsgLoadCharacterReply) {
	w.setSelfPosition(msg.MapNumber, msg.X, msg.Y, msg.Dir)
	w.self.HP = msg.HP
	w.self.MaxHP = msg.MaxHP
	w.self.Alive = true
	w.mp = msg.MP
	w.ag = msg.AG
	w.strength = msg.Strength
	w.dexterity = msg.Dexterity
	w.vitality = msg.Vitality
	w.energy = msg.Energy
	w.leadership = msg.Leadership
	w.characterSet = true
}

func (w *world) setSelfIndex(index int) {
	w.self.Index = index
}

func (w *world) mergeSelf(actor Actor) {
	w.self.Class = actor.Class
	w.self.ChangeUp = actor.ChangeUp
	w.self.Level = actor.Level
	w.self.HP = actor.HP
	w.self.MaxHP = actor.MaxHP
	if actor.HP > 0 {
		w.self.Alive = true
	}
	w.selfClassSet = true
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
	w.positionVersion++
}

func (w *world) clearSight() {
	clear(w.players)
	clear(w.objects)
}

func (w *world) markDead(index int) {
	if index == w.self.Index {
		w.self.Alive = false
		return
	}
	actor, ok := w.objects[index]
	if ok {
		actor.Alive = false
		actor.HP = 0
		w.objects[index] = actor
	}
}

func (w *world) Snapshot() WorldSnapshot {
	snapshot := WorldSnapshot{
		Phase:           w.phase,
		Failure:         w.failure,
		Self:            w.self,
		MP:              w.mp,
		AG:              w.ag,
		AttackSpeed:     w.attackSpeed,
		MagicSpeed:      w.magicSpeed,
		PositionVersion: w.positionVersion,
	}
	for _, actor := range w.players {
		snapshot.Players = append(snapshot.Players, actor)
	}
	for _, actor := range w.objects {
		snapshot.Objects = append(snapshot.Objects, actor)
	}
	for _, s := range w.skills {
		if s == nil {
			continue
		}
		snapshot.LearnedSkills = append(snapshot.LearnedSkills, s.Index)
		if s.SkillBase == nil || !isImplementedDamageSkill(s.Index) ||
			!supportedSkillType(s.Type) {
			continue
		}
		distance := s.Distance
		if distance < 1 {
			distance = 1
		}
		snapshot.Skills = append(snapshot.Skills, CombatSkill{
			Index:    s.Index,
			Damage:   s.DamageMax,
			Distance: distance,
			MP:       s.ManaUsage,
			AG:       s.BPUsage,
			Type:     s.Type,
			Delay:    s.Delay,
		})
	}
	if w.itemsSet && w.skillsSet && w.characterSet && w.selfClassSet {
		for position, it := range w.inventory {
			index, ok := skillIndexFromItem(it)
			if !ok {
				continue
			}
			if !w.canUseSkillItem(it) {
				continue
			}
			if _, ok := w.skills[index]; ok {
				continue
			}
			snapshot.LearnSkills = append(snapshot.LearnSkills, LearnSkill{
				Index:    index,
				Position: position,
			})
		}
	}
	return snapshot
}

func (w *world) fail(reason string) {
	w.failure = reason
	w.phase = PhaseFailed
}

func supportedSkillType(skillType int) bool {
	return skillType == skillTypeSiege ||
		skillType == skillTypePhysical ||
		skillType == skillTypeMagic
}

func skillIndexFromItem(it *item.Item) (int, bool) {
	if it == nil || it.ItemBase == nil || it.KindA != item.KindASkill {
		return 0, false
	}
	index := it.SkillIndex
	if it.Code == item.Code(12, 11) {
		index += it.Level
	}
	return index, index > 0
}

func (w *world) canUseSkillItem(it *item.Item) bool {
	if it == nil || it.ItemBase == nil || it.KindA != item.KindASkill ||
		w.self.Class < 0 || w.self.Class >= len(it.ReqClass) ||
		w.self.Level < it.ReqLevel ||
		w.strength < it.ReqStrength ||
		w.dexterity < it.ReqDexterity ||
		w.vitality < it.ReqVitality ||
		w.energy < it.ReqEnergy ||
		w.leadership < it.ReqCommand {
		return false
	}
	reqClass := it.ReqClass[w.self.Class]
	return reqClass != 0 && reqClass <= w.self.ChangeUp+1
}

func (s WorldSnapshot) object(index int) (Actor, bool) {
	for _, actor := range s.Objects {
		if actor.Index == index {
			return actor, true
		}
	}
	return Actor{}, false
}

func (s WorldSnapshot) combatSkill(index int) (CombatSkill, bool) {
	for _, combatSkill := range s.Skills {
		if combatSkill.Index == index {
			return combatSkill, true
		}
	}
	return CombatSkill{}, false
}

func (s WorldSnapshot) learnedSkill(index int) bool {
	for _, learned := range s.LearnedSkills {
		if learned == index {
			return true
		}
	}
	return false
}
