package object

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/skill"
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

func (m *objectManager) AddMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element int,
	newMonster func(class, mapNumber, startX, startY, endX, endY, dir, dis, element int) *Object) (*Object, error) {
	if m.monsterCount > m.maxMonsterCount {
		return nil, fmt.Errorf("over max monster count")
	}
	index := m.lastMonsterIndex
	cnt := m.maxMonsterCount
	for cnt > 0 {
		index++
		if index >= m.maxMonsterCount {
			index = m.monsterStartIndex
		}
		if m.objects[index] == nil {
			break
		}
		cnt--
	}
	if cnt == 0 {
		panic(fmt.Errorf("have no free monster index"))
	}
	m.lastMonsterIndex = index
	m.monsterCount++
	monster := newMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element)
	monster.Index = index
	m.objects[index] = monster
	return monster, nil
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

func (m *objectManager) AddPlayer(conn Conn, actioner Actioner, newPlayer func(Conn, Actioner) *Object) (int, error) {
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
	player := newPlayer(conn, actioner)
	player.Index = index
	// register the new player to object manager
	m.objects[index] = player

	// reply
	msg := model.MsgConnectReply{
		Result:  1,
		ID:      index,
		Version: conf.MapServers.ServerInfo.Version,
	}
	player.Push(&msg)
	log.Printf("player online [id]%d [addr]%s", player.Index, player.Addr())
	return index, nil
}

func (m *objectManager) DeletePlayer(id int) {
	if id < m.playerStartIndex {
		log.Printf("delete player failed [id]%d\n", id)
		return
	}
	player := m.objects[id]
	if player == nil {
		return
	}
	player.Offline()
	log.Printf("player offline [id]%d [addr]%s", player.Index, player.Addr())

	// unregister player from object manager
	player.Reset()
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
	log.Printf("user online [id]%d [addr]%s", u.index, u.conn.Addr())
	return index, nil
}

func (m *objectManager) DeleteUser(id int) {
	u := m.users[id]
	if u == nil {
		return
	}
	u.Offline()
	log.Printf("user offline [id]%d [addr]%s", u.index, u.conn.Addr())

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
		obj.ProcessAction()
	}
}

func (m *objectManager) Process1000ms() {
	type objects struct {
		Players  []*maps.Pot `json:"players"`
		Monsters []*maps.Pot `json:"monsters"`
		Npcs     []*maps.Pot `json:"npcs"`
	}
	table := make(map[int]*objects)
	for i := range m.mapSubscribeTable {
		table[i] = &objects{}
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
	MaxViewPort      = 75
	MaxArrayFrustrum = 4
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

type messageStateMachine struct {
	code    int
	subcode int
	sender  int
	time    time.Duration
}

type Objecter interface {
	Addr() string
	Offline()
	Push(any)
	PushMPAG(int, int)
	Chat(*model.MsgChat)
	Whisper(*model.MsgWhisper)
	Login(*model.MsgLogin)
	Logout(*model.MsgLogout)
	Hack(*model.MsgHack)
	GetCharacterList(*model.MsgGetCharacterList)
	CreateCharacter(*model.MsgCreateCharacter)
	DeleteCharacter(*model.MsgDeleteCharacter)
	CheckCharacter(*model.MsgCheckCharacter)
	LoadCharacter(*model.MsgLoadCharacter)
	MapDataLoadingOK(*model.MsgMapDataLoadingOK)
	DefineMuKey(*model.MsgDefineMuKey)
	DefineMuBot(*model.MsgDefineMuBot)
	EnableMuBot(*model.MsgEnableMuBot)
	LearnMasterSkill(*model.MsgLearnMasterSkill)
	GetPKLevel() int
	GetMasterLevel() int
	GetSkillMPAG(s *skill.Skill) (int, int)
	ProcessAction()
	Action(*model.MsgAction)
	Process1000ms()
	SpawnPosition()
	Die(*Object)
	Regen()
	GetChangeUp() int
	GetInventory() [9]*item.Item
	Teleport(*model.MsgTeleport)
	GetItem(*model.MsgGetItem)
	DropInventoryItem(*model.MsgDropInventoryItem)
	MoveItem(*model.MsgMoveItem)
	UseItem(*model.MsgUseItem)
	Talk(*model.MsgTalk)
	CloseTalkWindow(*model.MsgCloseTalkWindow)
	BuyItem(*model.MsgBuyItem)
	SellItem(*model.MsgSellItem)
	CloseWarehouseWindow(*model.MsgCloseWarehouseWindow)
	MapMove(*model.MsgMapMove)
	GetAttackRatePVP() int
	GetDefenseRatePVP() int
	GetIgnoreDefenseRate() int
	GetCriticalAttackRate() int
	GetCriticalAttackDamage() int
	GetExcellentAttackRate() int
	GetExcellentAttackDamage() int
	GetAddDamage() int
	GetHelperReduceDamage() int
	GetArmorReduceDamage() int
	GetWingIncreaseDamage() int
	GetWingReduceDamage() int
	GetDoubleDamageRate() int
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
	AddHP                     int
	ScriptMaxHP               int
	MP                        int // MP
	MaxMP                     int // MaxMP
	AddMP                     int
	SD                        int // SD
	MaxSD                     int
	AddSD                     int
	AG                        int // AG
	MaxAG                     int
	AddAG                     int
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
	Attribute                 int
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
	FrustrumX                 [MaxArrayFrustrum]int
	FrustrumY                 [MaxArrayFrustrum]int
	SkillFrustrumX            [MaxArrayFrustrum]int
	SkillFrustrumY            [MaxArrayFrustrum]int
	Viewports                 [MaxViewportNum]*Viewport // 主动视野
	ViewportNum               int
	msgs                      [20]*messageStateMachine

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
	// VPPShopPlayer            [MaxViewPort]int
	// VPPShopPlayerCount       uint16
	// BossGoldDerconMapNumber  uint8
	// lastTeleportTime         time.Time
	// clientHackLogCount       uint8
	// isInMonsterHerd      bool
	// isMonsterAttackFirst bool
	// monsterHerd          *monster.MonsterHerd
	// fsKillFrustrumX      [MaxArrayFrustrum]int
	// fsKillFrustrumY      [MaxArrayFrustrum]int
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
		obj.msgs[i] = &messageStateMachine{
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
	// panic(fmt.Sprintf("RandPosition failed [number]%d", number))
	log.Printf("RandPosition failed [map]%d [start](%d,%d) [end](%d,%d)\n", number, x1, y1, x2, y2)
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
