package object

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server-game/conf"
	"github.com/xujintao/balgass/src/server-game/game/item"
	"github.com/xujintao/balgass/src/server-game/game/maps"
	"github.com/xujintao/balgass/src/server-game/game/model"
	"github.com/xujintao/balgass/src/server-game/game/skill"
)

func init() {
	ObjectManager.init()
}

var ObjectManager objectManager

type objectManager struct {
	// monster
	maxMonsterCount   int
	monsterStartIndex int
	lastMonsterIndex  int
	monsterCount      int

	// call monster
	maxCallMonsterCount   int
	callMonsterStartIndex int
	lastCallMonsterIndex  int
	callMonsterCount      int

	// player
	maxPlayerCount   int
	playerStartIndex int
	lastPlayerIndex  int
	playerCount      int

	// objects
	maxObjectCount int
	objects        []*Object

	// users
	maxUserCount      int
	userStartIndex    int
	lastUserIndex     int
	userCount         int
	users             []*user
	mapSubscribeTable map[int]map[*user]struct{}
}

func (m *objectManager) init() {
	// monster
	m.maxMonsterCount = conf.Server.GameServerInfo.MaxMonsterCount
	m.monsterStartIndex = 0
	m.lastMonsterIndex = m.monsterStartIndex - 1

	// call monster
	m.maxCallMonsterCount = conf.Server.GameServerInfo.MaxSummonMonsterCount
	m.callMonsterStartIndex = m.maxMonsterCount
	m.lastCallMonsterIndex = m.callMonsterStartIndex - 1

	// player
	m.maxPlayerCount = conf.Server.GameServerInfo.MaxPlayerCount
	m.playerStartIndex = m.maxMonsterCount + m.maxCallMonsterCount
	m.lastPlayerIndex = m.playerStartIndex - 1

	// objects
	m.maxObjectCount = m.maxMonsterCount + m.maxCallMonsterCount + m.maxPlayerCount
	m.objects = make([]*Object, m.maxObjectCount)

	// users
	m.maxUserCount = 10
	m.userStartIndex = 0
	m.lastUserIndex = m.userStartIndex - 1
	m.users = make([]*user, m.maxUserCount)
	m.mapSubscribeTable = make(map[int]map[*user]struct{})
}

func (om *objectManager) AddMonster(newMonster func() *Object) (*Object, error) {
	if om.monsterCount > om.maxMonsterCount {
		return nil, fmt.Errorf("over max monster count")
	}
	index := om.lastMonsterIndex
	cnt := om.maxMonsterCount
	for cnt > 0 {
		index++
		if index >= om.maxMonsterCount {
			index = om.monsterStartIndex
		}
		if om.objects[index] == nil {
			break
		}
		cnt--
	}
	if cnt == 0 {
		panic(fmt.Errorf("have no free monster index"))
	}
	om.lastMonsterIndex = index
	om.monsterCount++
	m := newMonster()
	m.Index = index
	om.objects[index] = m
	return m, nil
}

// func (m *objectManager) AddCallMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element int) (int, error) {
// 	if m.callMonsterCount > m.maxCallMonsterCount {
// 		return -1, fmt.Errorf("over max call monster count")
// 	}
// 	index := m.lastCallMonsterIndex
// 	cnt := m.maxCallMonsterCount
// 	for cnt > 0 {
// 		index++
// 		if index >= m.playerStartIndex {
// 			index = m.callMonsterStartIndex
// 		}
// 		if m.objects[index] == nil {
// 			break
// 		}
// 		cnt--
// 	}
// 	if cnt == 0 {
// 		panic(fmt.Errorf("have no free call monster index"))
// 	}
// 	m.lastCallMonsterIndex = index
// 	m.callMonsterCount++
// 	monster := NewMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element)
// 	monster.index = index
// 	m.objects[index] = &monster.object
// 	return index, nil
// }

type Conn interface {
	Addr() string
	Write(any) error
	Close() error
}
type Actioner interface {
	PlayerAction(int, string, any)
}

func (m *objectManager) AddPlayer(conn Conn, newPlayer func() *Object) (int, error) {
	// limit max player count
	if m.playerCount >= m.maxPlayerCount {
		// reply
		msg := model.MsgLoginReply{Result: 4}
		conn.Write(&msg)
		return -1, fmt.Errorf("over max player count")
	}

	// get unified object index
	index := m.lastPlayerIndex
	cnt := m.maxPlayerCount
	for cnt > 0 {
		index++
		if index >= m.maxObjectCount {
			index = m.playerStartIndex
		}
		if m.objects[index] == nil {
			break
		}
		cnt--
	}
	if cnt == 0 {
		panic(fmt.Errorf("have no free player index"))
	}
	m.lastPlayerIndex = index
	m.playerCount++
	p := newPlayer()
	p.Index = index
	m.objects[index] = p

	// reply
	msg := model.MsgConnectReply{
		Result:  1,
		ID:      index,
		Version: conf.MapServers.ServerInfo.Version,
	}
	p.Push(&msg)
	slog.Info("player online", "id", p.Index, "addr", p.Addr())
	return index, nil
}

func (m *objectManager) DeletePlayer(id int) {
	if id < m.playerStartIndex {
		slog.Error("DeletePlayer id < m.playerStartIndex", "id", id)
		return
	}
	p := m.objects[id]
	if p == nil {
		return
	}
	p.Offline()
	slog.Info("player offline", "id", p.Index, "addr", p.Addr())

	// unregister player from object manager
	p.Reset()
	m.objects[id] = nil
	m.playerCount--
}

func (m *objectManager) GetObject(id int) *Object {
	if id >= m.maxObjectCount {
		return nil
	}
	return m.objects[id]
}

func (m *objectManager) GetPlayerByName(name string) *Object {
	for _, tobj := range m.objects[m.playerStartIndex:] {
		if tobj == nil {
			continue
		}
		if tobj.Name == name {
			return tobj
		}
	}
	return nil
}

func (m *objectManager) GetPlayerPercent() int {
	return m.playerCount / m.maxPlayerCount * 100
}

func (m *objectManager) OfflineAllObjects() {
	for _, obj := range m.objects[m.playerStartIndex:] {
		if obj == nil {
			continue
		}
		obj.Offline()
	}
	for _, u := range m.users[m.userStartIndex:] {
		if u == nil {
			continue
		}
		u.Offline()
	}
}

func (m *objectManager) GetOnlineObjectsNumber() *model.MsgGetOnlineObjectNumberReply {
	msg := model.MsgGetOnlineObjectNumberReply{
		PlayerNumber: m.playerCount,
		UserNumber:   m.userCount,
	}
	return &msg
}

func (m *objectManager) AddUser(conn Conn) (int, error) {
	if m.userCount >= m.maxUserCount {
		return -1, fmt.Errorf("over max user count")
	}
	index := m.lastUserIndex
	cnt := m.maxUserCount
	for cnt > 0 {
		cnt--
		index++
		if index >= m.maxUserCount {
			index = m.userStartIndex
		}
		if m.users[index] == nil {
			break
		}
	}
	if cnt == 0 {
		panic(fmt.Errorf("have no free user index"))
	}
	m.lastUserIndex = index
	m.userCount++
	u := NewUser(conn)
	u.index = index
	m.users[index] = u
	slog.Info("user online", "id", u.index, "addr", u.conn.Addr())
	return index, nil
}

func (m *objectManager) DeleteUser(id int) {
	u := m.users[id]
	if u == nil {
		return
	}
	u.Offline()
	slog.Info("user offline", "id", u.index, "addr", u.conn.Addr())

	// unregister user from object manager
	m.users[id] = nil
	m.userCount--
}

func (m *objectManager) GetUser(id int) *user {
	return m.users[id]
}

func (m *objectManager) Process100ms() {
	for _, obj := range m.objects {
		if obj == nil {
			continue
		}
		obj.processMove()
		obj.processDelayMsg()
		obj.ProcessAction()
	}
}

func (m *objectManager) Process1000ms() {
	type objects struct {
		Players  []*maps.Pot `json:"players"`
		Monsters []*maps.Pot `json:"monsters"`
		Npcs     []*maps.Pot `json:"npcs"`
		Stands   []*maps.Pot `json:"stands"`
	}
	table := make(map[int]*objects)
	for i := range m.mapSubscribeTable {
		table[i] = &objects{Stands: maps.MapManager.GetMapStands(i)}
	}

	for _, obj := range m.objects {
		if obj == nil {
			continue
		}
		obj.Process1000ms()
		obj.processViewport() // 1->2
		obj.processRegen()    // 4->1

		// process map subscripion
		if _, ok := m.mapSubscribeTable[obj.MapNumber]; ok {
			if !obj.Live {
				continue
			}
			p := maps.Pot{X: obj.X, Y: obj.Y}
			switch obj.Type {
			case ObjectTypePlayer:
				table[obj.MapNumber].Players = append(table[obj.MapNumber].Players, &p)
			case ObjectTypeMonster:
				table[obj.MapNumber].Monsters = append(table[obj.MapNumber].Monsters, &p)
			case ObjectTypeNPC:
				table[obj.MapNumber].Npcs = append(table[obj.MapNumber].Npcs, &p)
			default:
			}
		}
	}
	for i, users := range m.mapSubscribeTable {
		for u := range users {
			u.publishMap(table[i])
		}
	}
}

const (
	MaxMonsterSendMsg       = 20
	MaxMonsterSendAttackMsg = 100
	MaxMonsterType          = 1024
	MaxGuildLen             = 8
	MaxAccountIDLen         = 10
	MaxCharacterNameLen     = 10
	TradeBoxSize            = 32
	MagicSize               = 150
	InventoryWearSize       = 12
	InventoryBagStart       = InventoryWearSize
	TradeBoxMapSize         = 32
	PShopSize               = 32
	PShopMapSize            = 32
	PShopRangeStart         = 204
	PShopRangeEnd           = 236
	MaxSelfDefense          = 5
	MaxBuffEffect           = 32
	// MaxResistanceType = 7
	MaxViewport     = 75
	MaxArrayFrustum = 4
	MaxZen          = 2000000000
)

type EffectList struct {
	BuffIndex      byte
	EffectCategory byte
	EffectType1    byte
	EffectType2    byte
	EffectValue1   int
	EffectValue2   int
	EffectSetTime  uint
	EffectDuration int
}

const (
	MaxViewportNum = 75
)

type Viewport struct {
	State  int
	Number int
	Type   int
	// index  int
	Dis int
}

type HitDamage struct {
	Number      uint16
	HitDamage   int
	LastHitTime time.Time
}

type InterfaceState struct {
	Use   uint8
	State uint8
	Type  uint16
}

type ObjectType int

const (
	ObjectTypeEmpty ObjectType = iota
	ObjectTypePlayer
	ObjectTypeMonster
	ObjectTypeNPC
)

type NpcType int

const (
	NpcTypeNone NpcType = iota
	NpcTypeShop
	NpcTypeWarehouse
	NpcTypeChaosMix
	NpcTypeGoldarcher
	NpcTypePentagramMix
)

type ConnectState int

const (
	ConnectStateEmpty ConnectState = iota
	ConnectStateConnected
	ConnectStateLogged
	ConnectStatePlaying
)

type skillInfo struct {
	ghostPhantomX        uint8
	ghostPhantomY        uint8
	remedyOfLoveEffect   uint16
	remedyOfLoveTime     uint16
	lordSummonTime       uint16
	lordSummonMapNumber  uint8
	lordSummonX          uint8
	lordSummonY          uint8
	soulBarrierDefence   int
	soulBarrierManaRate  int
	posionType           uint8
	iceType              uint8
	infinityArrowIncRate float32
	circleShieldRate     float32
}

type DelayMsg struct {
	code    int
	subcode int
	sender  int
	time    time.Time
}

type Objecter interface {
	// object actions:
	// 1. connection safety
	KeepLive(*model.MsgKeepLive)
	Hack(*model.MsgHack)
	// 2. account and character
	Login(*model.MsgLogin)
	Logout(*model.MsgLogout)
	GetCharacterList(*model.MsgGetCharacterList)
	CreateCharacter(*model.MsgCreateCharacter)
	DeleteCharacter(*model.MsgDeleteCharacter)
	CheckCharacter(*model.MsgCheckCharacter)
	LoadCharacter(*model.MsgLoadCharacter)
	// 3. item management
	PickItem(*model.MsgPickItem) // implemented by Object
	DropItem(*model.MsgDropItem) // implemented by Object
	BuyItem(*model.MsgBuyItem)   // implemented by Object
	SellItem(*model.MsgSellItem) // implemented by Object
	MoveItem(*model.MsgMoveItem) // implemented by Object
	UseItem(*model.MsgUseItem)   // implemented by Object
	// 4. behavior management
	Chat(*model.MsgChat)
	Whisper(*model.MsgWhisper)
	Talk(*model.MsgTalk)
	CloseTalkWindow(*model.MsgCloseTalkWindow)
	CloseWarehouseWindow(*model.MsgCloseWarehouseWindow)
	Move(*model.MsgMove)               // implemented by Object
	Attack(*model.MsgAttack)           // implemented by Object
	UseSkill(*model.MsgUseSkill)       // implemented by Object
	SetPosition(*model.MsgSetPosition) // implemented by Object
	MapMove(*model.MsgMapMove)
	MapDataLoadingOK(*model.MsgMapDataLoadingOK)
	Teleport(*model.MsgTeleport)
	Action(*model.MsgAction)
	BattleCoreNotice(*model.MsgBattleCoreNotice)
	AddLevelPoint(*model.MsgAddLevelPoint)
	LearnMasterSkill(*model.MsgLearnMasterSkill)
	DefineMuKey(*model.MsgDefineMuKey)
	DefineMuBot(*model.MsgDefineMuBot)
	EnableMuBot(*model.MsgEnableMuBot)
	UsePet(*model.MsgUsePet)
	// 5. Others
	MuunSystem(*model.MsgMuunSystem)

	// object actions implemented by derived object:
	Addr() string
	Offline()
	Push(any)
	GetPKLevel() int
	GetMasterLevel() int
	IsMasterLevel() bool
	GetSkillMPAG(s *skill.Skill) (int, int)
	ProcessAction()
	Process1000ms()
	SpawnPosition()
	Die(*Object, int)
	LevelUp(int) bool
	DieDropItem(*Object)
	Regen()
	GetChangeUp() int
	GetAttackRatePVP() int
	GetDefenseRatePVP() int
	GetIgnoreDefenseRate() int
	GetCriticalAttackRate() int
	GetCriticalAttackDamage() int
	GetExcellentAttackRate() int
	GetExcellentAttackDamage() int
	GetMonsterDieGetHP() float64
	GetMonsterDieGetMP() float64
	GetAddDamage() int
	GetArmorReduceDamage() int
	GetWingIncreaseDamage() int
	GetWingReduceDamage() int
	GetHelperReduceDamage() int
	GetPetIncreaseDamage() int
	GetPetReduceDamage() int
	GetDoubleDamageRate() int
	GetMonsterDieGetMoney() float64
	GetKnightGladiatorCalcSkillBonus() float64
	GetImpaleSkillCalc() float64
	SetMoney(int)
	GetMoney() int
	GetInventory() *item.Inventory
	GetWarehouse() *item.Warehouse
	EquipmentChanged()
	SetDelayRecoverHP(int, int)
	SetDelayRecoverSD(int, int)
	decreaseItemDurability(int)
}

type Object struct {
	Objecter
	Index                     int
	ConnectState              ConnectState
	Live                      bool
	State                     int // 1:初始 2:视野 4:死亡 8:清理
	StartX                    int
	StartY                    int
	X                         int // x坐标
	Y                         int // y坐标
	Dir                       int // 方向
	TX                        int // 目标x坐标
	TY                        int // 目标y坐标
	pathX                     [15]int
	pathY                     [15]int
	pathDir                   [15]int
	pathCount                 int
	pathCur                   int
	pathTime                  time.Time
	PathMoving                bool
	delayLevel                int
	MapNumber                 int        // 地图号
	Type                      ObjectType // 对象种类：玩家，怪物，NPC
	NpcType                   NpcType
	Class                     int    // 对象类别。怪物和玩家都有类别
	Name                      string // 对象名称
	Annotation                string // 对象备注
	Level                     int
	HP                        int // HP
	MaxHP                     int // MaxHP
	ScriptMaxHP               int
	MP                        int // MP
	MaxMP                     int // MaxMP
	SD                        int // SD
	MaxSD                     int
	AG                        int // AG
	MaxAG                     int
	TargetNumber              int
	AttackMin                 int // 物攻min
	AttackMax                 int // 物攻max
	AttackSpeed               int // 物攻速度
	AttackRate                int // 攻击率
	Defense                   int // 防御力
	DefenseRate               int // 防御率
	MagicDefense              int // 魔法防御率
	MoveSpeed                 int // 移动速度
	AttackRange               int // 攻击范围
	AttackType                int // 攻击类型
	ViewRange                 int // 视野范围
	ItemDropRate              int // 道具掉落率
	MoneyDropRate             int // 金钱掉落率
	MoneyDrop                 int
	Attribute                 int // 0:passive monster 1: invisible monster 2: normal monster
	dieTime                   time.Time
	dieRegen                  bool
	MaxRegenTime              time.Duration // 最大重生时间
	PentagramMainAttribute    int
	PentagramAttributePattern int
	PentagramAttackMin        int
	PentagramAttackMax        int
	PentagramAttackRate       int
	PentagramDefense          int
	Skills                    skill.Skills
	FrustumX                  [MaxArrayFrustum]int
	FrustumY                  [MaxArrayFrustum]int
	SkillFrustumX             [MaxArrayFrustum]int
	SkillFrustumY             [MaxArrayFrustum]int
	Viewports                 [MaxViewportNum]*Viewport // for attack
	ViewportsNum              int
	ViewportsPassive          [MaxViewportNum]*Viewport // for push
	ViewportsPassiveNum       int
	msgs                      [20]*DelayMsg

	// groupNumber     int
	// subGroupNumber  int
	// groupMemberGUID int
	// regenType       int
	// // argo                    *monster.MonsterAIAgro
	// lastAutoRunTime            time.Time
	// lastAutoDelay              time.Duration
	// expType                    int
	// LoginMsgSend               bool
	// LoginMsgCount              byte
	// CloseCount                 byte
	// CloseType                  byte
	// EnableCharacterDel         bool
	// UserNumber                 int
	// DBNumber                   int
	// EnableCharacterCreate      bool
	// AutoSaveTime               time.Time
	// ConnectCheckTime           time.Time
	// CheckTick                  uint
	// CheckSpeedHack             bool
	// CheckTick2                 uint
	// CheckTickCount             byte
	// PintTime                   int
	// TimeCount                  byte
	// PKTimer                    *time.Timer
	// CheckSumTableNum           uint16
	// CheckSumTime               uint
	// ChatLimitTime              uint16
	// ChatLimitTimeSec           byte
	// FillLifeCount              byte
	// VitalityToLife             float32
	// EnergyToMana               float32
	// XSave                      uint16
	// YSave                      uint16
	// MapNumberSave              byte
	// XDie                       uint16
	// YDie                       uint16
	// MapNumberDie               byte
	// IFillShieldMax             int
	// IFillShield                int
	// IFillShieldCount           int
	// ShieldAutoRefillTimer      *time.Timer
	// AutoHPRecovery             byte // 自动生命恢复
	// Authority                  uint
	// AuthorityCode              uint
	// Penalty                    uint
	// GameMaster                 uint
	// PenaltyMask                uint
	// ChatBlockTime              time.Time
	// AccountItemBlock           byte
	// ActionNumber               byte
	// ActionTime                 uint
	// ActionCount                byte
	// ChatFloodTime              uint
	// ChatFloodCount             byte
	// Rest                       byte
	// viewState                  byte
	// buffEffectCount            byte
	// buffEffectList             [MaxBuffEffect]EffectList
	// lastMoveTime               uint
	// lastAttackTime             uint
	// friendServerOnline         byte
	// detectSpeedHackTime        time.Duration
	// sumLastAttackTime          time.Duration
	// detectCount                uint
	// detectHackKickCount        uint
	// speedHackPenalty           int
	// attackSpeedHackDetectCount uint
	// packetCheckTime            time.Duration
	// ShopTime                   time.Duration
	// totalAttackTime            time.Duration
	// totalAttackCount           uint
	// TeleportTIme               time.Duration
	// Teleport                   byte
	// KillerType                 byte
	// MapNumberRegen             byte
	// XRegen                     byte
	// YRegen                     byte
	// posNum                     uint16
	// LifeRefillTimer            *time.Timer
	// ActionTimeCur              time.Time
	// ActionTimeNext             time.Time
	// ActionTimeDelay            time.Duration
	// monsterBattleDelay         byte
	// kalimaGateExist            byte
	// kalimaGateIndex            int
	// kalimaGateEnterCount       byte
	// AttackObj                  *object
	// skillNumber                uint16
	// skillTime                  time.Duration
	// attackerKilled             bool
	// manaFillCount              byte
	// lifeFillCount              byte
	// SelfDefense                [MaxSelfDefense]int
	// SelfDefenseTime            time.Duration
	// PartyNumber                int
	// PartyTargetUser            int
	// Married                    byte
	// MarryName                  string
	// MarryRequested             byte
	// WinDuels                   int
	// LoseDuels                  int
	// MarryRequestIndex          uint16
	// MarryRequestTime           time.Duration
	// RecallMon                  int
	// change                     int
	// TargetNpcNumber            int
	// LastAttackerID             int
	// magicBack                   *skill.MagicInfo
	// Magic                       *skill.MagicInfo
	// UseMagicNumber  byte
	// UseMagicTime    time.Duration
	// UseMagicCount   byte
	// OSAttackSerial  uint16
	// SASCount        byte
	// UseSkillTime time.Duration
	// CharSet         string
	// resistance               [MaxResistanceType]byte
	// addResistance            [MaxResistanceType]byte
	// HD                       *HitDamage
	// HDCount                  uint16
	// ifState                  InterfaceState
	// iterfaceTime             time.Duration
	// InventoryMap             *uint8
	// InventoryCount           *uint8
	// Transaction              uint8
	// Inventory1               *item.Item
	// InventoryMap1            *uint8
	// InventoryCount1          uint8
	// Inventory2               *item.Item
	// InventoryMap2            *uint8
	// InventoryCount2          uint8
	// Trade                    *item.Item
	// TradeMap                 *uint8
	// TradeMoney               int
	// TradeOK                  bool
	// Warehouse                *item.Item
	// WarehouseID              uint8
	// WarehouseTick            time.Time
	// WarehouseMap             *uint8
	// WarehouseCount           uint8
	// WarehousePW              uint16
	// WarehouseLock            uint8
	// WarehouseUnfailLock      uint8
	// WarehouseMoney           int
	// ChaosBox                 *item.Item
	// ChaosBoxMap              *uint8
	// ChaosMoney               int
	// ChaosSuccessRate         int
	// ChaosMassMixCurItem      uint8
	// ChaosMassMixSuccessCount uint8
	// ChaosLock                bool
	// Option                   uint
	// eventScore               int
	// eventExp                 int
	// eventMoney               int
	// devilSquareIndex         uint8
	// devilSquareAuth          bool
	// BloodCastlIndex          uint8
	// BloodCastleSubIndex      uint8
	// BloodCastleExp           int
	// BloodCastleComplete      bool
	// ChaosCastleIndex         uint8
	// ChaosCastleSubIndex      uint8
	// ChaosCastleBlowTime      time.Duration
	// isCCFUIUsing             bool
	// CCFCanEnter              uint8
	// CCFCertiTick             time.Time
	// CCFUserIndex             int
	// CCFBlowTime              time.Time
	// killUserCount            uint8
	// killMobCount             uint8
	// isCCFQuitMsg             bool
	// illusionTempleIndex      uint8
	// zoneIndex                uint8
	// ckillUserCount           uint8
	// cKillMonsterCount        uint8
	// duelUserReserved         int
	// duelUserRequested        int
	// duelUser                 int
	// duelRoom                 int
	// duelScore                uint8
	// duelTickCount            time.Duration
	// IsInBattleGround         bool
	// HaveWeaponInHand         bool
	// EventChipCount           uint16
	// LuckyCoinCount           int
	// MutoNumber               int
	// UseEventServer           bool
	// LoadWarehouseInfo        bool
	// StoneCount               int
	// maxLifePower             int
	// checkLifeTime            int
	// moveToOtherServer        uint8
	// BackName                 string
	// isPShopOpen              bool
	// isPShopTransaction       bool
	// isPShopItemChange        bool
	// isPShopRedrawABS         bool
	// PShopText                string
	// isPShopWantDeal          bool
	// PShopDealerIndex         int
	// PShopDealerName          string
	// muPShopTrade             sync.Mutex
	// VPPShopPlayer            [MaxViewport]int
	// VPPShopPlayerCount       uint16
	// BossGoldDerconMapNumber  uint8
	// lastTeleportTime         time.Time
	// clientHackLogCount       uint8
	// isInMonsterHerd      bool
	// isMonsterAttackFirst bool
	// monsterHerd          *monster.MonsterHerd
	// skillFrustumX      [MaxArrayFrustum]int
	// skillFrustumY      [MaxArrayFrustum]int
	// // durMagicKeyChecker          *skill.DurMagicKeyChecker
	// IsChaosMixCompleted         bool
	// SkillLongSpearChange        bool
	// objectSecTimer              time.Timer
	// mapSvrMoveQuit              bool
	// mapSvrMoveReq               bool
	// mapSvrMoveReq2              bool
	// mapSvrQuitTick              time.Time
	// prevMapSvrCode              uint16
	// destMapNumber               uint16
	// destX                       uint8
	// destY                       uint8
	// csNpcExistVal               int
	// csNpcType                   uint8
	// csGateOPen                  uint8
	// csGateLeverLinkIndex        int
	// csNpcDfLevel                uint8
	// csNpcRgLevel                uint8
	// csJoinSide                  uint8
	// csGuildInvolved             bool
	// IsCastleNPCUpgradeCompleted bool
	// weaponState                 uint8
	// weaponUser                  int
	// killCount                   uint8
	// accumulatedDamage           int
	// lifeStoneCount              uint8
	// creationState               uint8
	// createdActiviationTime      int
	// accumulatedCrownAccessTime  int
	// monsterSkillElementInfo     monster.MonsterSkillElementInfo
	// crywolfMVPScore             int
	// lastCheckTick               time.Time
	// autoRecuperationTime        time.Time
	// skillDistanceErrorCount     int
	// skillDistanceErrorTick      time.Time
	// skillInfo                   skillInfo
	// bufferIndex                 int
	// buffID                      int
	// buffPlayerIndex             int
	// agiCheckTime                time.Time
	// warehouseSaveLock           uint8
	// crcCheckTime                time.Time
	// off                         bool
	// offLevel                    bool
	// offLevelTime                int
}

func (obj *Object) Init() {
	obj.TargetNumber = -1
	obj.initSkill()
	obj.initViewport()
	obj.initMessage()
}

func (obj *Object) Reset() {
	obj.Name = ""
	obj.TargetNumber = -1
	obj.Live = false
	obj.clearSkill()
	obj.clearViewport()
}

func (obj *Object) initMessage() {
	for i := range obj.msgs {
		obj.msgs[i] = &DelayMsg{
			code: -1,
		}
	}
}

func (obj *Object) Test(msg *model.MsgTest) {
	obj.Push(msg)
}

func (obj *Object) RandPosition(number, x1, y1, x2, y2 int) (int, int) {
	w := x2 - x1
	if w <= 0 {
		w = 1
	}
	h := y2 - y1
	if h <= 0 {
		h = 1
	}
	if w == 1 && h == 1 {
		return x1, y1
	}
	for i := 0; i < 100; i++ {
		x := x1 + rand.Intn(w)
		y := y1 + rand.Intn(h)
		attr := maps.MapManager.GetMapAttr(number, x, y)
		if attr&1 == 0 && attr&4 == 0 && attr&8 == 0 {
			return x, y
		}
	}
	slog.Error("RandPosition",
		"index", obj.Index,
		"name", obj.Name,
		"map", number,
		"start", fmt.Sprintf("(%d,%d)", x1, y1),
		"end", fmt.Sprintf("(%d,%d)", x2, y2),
	)
	return x1, y1
}

func (obj *Object) processRegen() {
	if !obj.dieRegen {
		return
	}
	if obj.ConnectState < ConnectStatePlaying {
		return
	}
	now := time.Now()
	if now.After(obj.dieTime.Add(5 * time.Second)) {
		if obj.State == 4 {
			obj.State = 8
		}
	}
	if now.Before(obj.dieTime.Add(obj.MaxRegenTime)) {
		return
	}
	obj.SpawnPosition()
	obj.Regen()
	obj.dieRegen = false
	obj.State = 1
	obj.Live = true
}

func (obj *Object) AddDelayMsg(code, subcode, delay, sender int) {
	for _, msg := range obj.msgs {
		if msg.code == -1 {
			msg.code = code
			msg.subcode = subcode
			msg.time = time.Now().Add(time.Duration(delay) * time.Millisecond)
			msg.sender = sender
			break
		}
	}
}

func (obj *Object) processDelayMsg() {
	now := time.Now()
	for _, msg := range obj.msgs {
		if msg.code == -1 {
			continue
		}
		if now.Before(msg.time) {
			continue
		}
		switch msg.code {
		case 0: // die give target experience
		case 1: // die drop item
			tobj := ObjectManager.objects[msg.sender]
			if tobj == nil || !tobj.Live {
				break
			}
			obj.DieDropItem(tobj)
		case 2: // skill knockback
			tobj := ObjectManager.objects[msg.sender]
			if tobj == nil || !tobj.Live {
				break
			}
			obj.Knockback(tobj)
		case 3: // die recover target hp/mp
			tobj := ObjectManager.objects[msg.sender]
			if tobj == nil || !tobj.Live {
				break
			}
			obj.DieRecoverHPMP(tobj)
		}
		msg.code = -1
	}
}

func (obj *Object) PushHPSD(hp, sd int) {
	obj.Push(&model.MsgHPReply{Position: -1, HP: hp, SD: sd})
}

func (obj *Object) PushMPAG(mp, ag int) {
	obj.Push(&model.MsgMPReply{Position: -1, MP: mp, AG: ag})
}

func (obj *Object) PushMaxHPSD(hp, sd int) {
	obj.Push(&model.MsgHPReply{Position: -2, HP: hp, SD: sd})
}

func (obj *Object) PushMaxMPAG(mp, ag int) {
	obj.Push(&model.MsgMPReply{Position: -2, MP: mp, AG: ag})
}

func (obj *Object) PushSystemMsg(msg string) {
	reply := model.MsgWhisperReply{}
	reply.Name = "system"
	reply.Msg = msg
	obj.Push(&reply)
}
