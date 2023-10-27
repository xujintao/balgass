package object

import (
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/model"
	"github.com/xujintao/balgass/src/server_game/game/skill"
)

func init() {
	InitFrustrum()
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
	objects        []*object

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
	m.objects = make([]*object, m.maxObjectCount)

	// users
	m.maxUserCount = 10
	m.userStartIndex = 0
	m.lastUserIndex = m.userStartIndex - 1
	m.users = make([]*user, m.maxUserCount)
	m.mapSubscribeTable = make(map[int]map[*user]struct{})

	// 先有怪后有玩家
	m.spawnMonster()
}

func (m *objectManager) spawnMonster() {
	// MonsterSpawn was generated 2023-07-17 16:05:41 by https://xml-to-go.github.io/ in Ukraine.
	type MonsterSpawn struct {
		XMLName xml.Name `xml:"MonsterSpawn"`
		Text    string   `xml:",chardata"`
		Map     []*struct {
			Text       string `xml:",chardata"`
			Number     int    `xml:"Number,attr"`
			Name       string `xml:"Name,attr"`
			Annotation string `xml:"annotation,attr"`
			Spot       []*struct {
				Text        string `xml:",chardata"`
				Type        int    `xml:"Type,attr"`
				Description string `xml:"Description,attr"`
				Spawn       []*struct {
					Text     string `xml:",chardata"`
					Index    int    `xml:"Index,attr"`
					Distance int    `xml:"Distance,attr"`
					StartX   int    `xml:"StartX,attr"`
					StartY   int    `xml:"StartY,attr"`
					Dir      int    `xml:"Dir,attr"`
					EndX     int    `xml:"EndX,attr"`
					EndY     int    `xml:"EndY,attr"`
					Count    int    `xml:"Count,attr"`
					Element  int    `xml:"Element,attr"`
				} `xml:"Spawn"`
			} `xml:"Spot"`
		} `xml:"Map"`
	}
	var monsterSpawn MonsterSpawn
	conf.XML(conf.PathCommon, "Monsters/IGC_MonsterSpawn.xml", &monsterSpawn)
	for _, _map := range monsterSpawn.Map {
		for _, spot := range _map.Spot {
			for _, spawn := range spot.Spawn {
				cnt := spawn.Count
				if cnt == 0 {
					cnt = 1
				}
				for i := 0; i < cnt; i++ {
					spawnClass := spawn.Index
					spawnMapNumber := _map.Number
					spawnStartX := 0
					spawnStartY := 0
					spawnEndX := 0
					spawnEndY := 0
					switch spot.Type {
					case 0: // npc
						spawnStartX = spawn.StartX
						spawnStartY = spawn.StartY
						spawnEndX = spawn.StartX
						spawnEndY = spawn.StartY
					case 1, 3: // multiple
						spawnStartX = spawn.StartX
						spawnStartY = spawn.StartY
						spawnEndX = spawn.EndX
						spawnEndY = spawn.EndY
					case 2: // single
						spawnStartX = spawn.StartX - 3
						spawnStartY = spawn.StartY - 3
						spawnEndX = spawn.StartX + 3
						spawnEndY = spawn.StartY + 3
					}
					spawnDir := spawn.Dir
					spawnDis := spawn.Distance
					spawnElement := spawn.Element
					_, err := m.AddMonster(
						spawnClass,
						spawnMapNumber,
						spawnStartX,
						spawnStartY,
						spawnEndX,
						spawnEndY,
						spawnDir,
						spawnDis,
						spawnElement,
					)
					if err != nil {
						log.Fatalf("AddMonster failed err[%v]", err)
					}
				}
			}
		}
	}
}

func (m *objectManager) AddMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element int) (int, error) {
	if m.monsterCount > m.maxMonsterCount {
		return -1, fmt.Errorf("over max monster count")
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
	monster := NewMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element)
	monster.objecter = monster
	monster.objectManager = m
	monster.index = index
	m.objects[index] = &monster.object
	return index, nil
}

func (m *objectManager) AddCallMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element int) (int, error) {
	if m.callMonsterCount > m.maxCallMonsterCount {
		return -1, fmt.Errorf("over max call monster count")
	}
	index := m.lastCallMonsterIndex
	cnt := m.maxCallMonsterCount
	for cnt > 0 {
		index++
		if index >= m.playerStartIndex {
			index = m.callMonsterStartIndex
		}
		if m.objects[index] == nil {
			break
		}
		cnt--
	}
	if cnt == 0 {
		panic(fmt.Errorf("have no free call monster index"))
	}
	m.lastCallMonsterIndex = index
	m.callMonsterCount++
	monster := NewMonster(class, mapNumber, startX, startY, endX, endY, dir, dis, element)
	monster.objectManager = m
	monster.index = index
	m.objects[index] = &monster.object
	return index, nil
}

func (m *objectManager) AddPlayer(conn Conn, actioner actioner) (int, error) {
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
	player := NewPlayer(conn, actioner)
	player.objecter = player
	player.objectManager = m
	player.index = index
	// register the new player to object manager
	m.objects[index] = &player.object

	// reply
	msg := model.MsgConnectReply{
		Result:  1,
		ID:      index,
		Version: conf.MapServers.ServerInfo.Version,
	}
	player.push(&msg)
	log.Printf("player online [id]%d [addr]%s", player.index, player.addr())
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
	log.Printf("player offline [id]%d [addr]%s", player.index, player.addr())

	// unregister player from object manager
	player.reset()
	m.objects[id] = nil
	m.playerCount--
}

func (m *objectManager) GetPlayer(id int) *object {
	return m.objects[id]
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
	u.objectManager = m
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
		obj.processAction()
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
		obj.processViewport() // 1->2
		obj.processRegen()    // 4->1

		// process map subscripion
		if _, ok := m.mapSubscribeTable[obj.MapNumber]; ok {
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

type viewport struct {
	state  int
	number int
	type_  int
	// index  int
	dis int
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

type objecter interface {
	addr() string
	Offline()
	push(any)
	Login(*model.MsgLogin)
	// SetAccount(*model.MsgSetAccount)
	CreateCharacter(*model.MsgCreateCharacter)
	GetCharacterList(*model.MsgGetCharacterList)
	DeleteCharacter(*model.MsgDeleteCharacter)
	PickCharacter(*model.MsgPickCharacter)
	SetCharacter(*model.MsgSetCharacter)
	getPKLevel() int
	processAction()
	// processRegen()
	spawnPosition()
}

type object struct {
	objecter
	*objectManager
	index              int
	ConnectState       ConnectState
	Live               bool
	State              int // 1:初始 2:视野 4:死亡
	StartX             int
	StartY             int
	X                  int // x坐标
	Y                  int // y坐标
	Dir                int // 方向
	TX                 int // 目标x坐标
	TY                 int // 目标y坐标
	pathX              [15]int
	pathY              [15]int
	pathDir            [15]int
	pathCount          int
	pathCur            int
	pathTime           int64
	pathMoving         bool
	delayLevel         int
	MapNumber          int        // 地图号
	Type               ObjectType // 对象种类：玩家，怪物，NPC
	Class              int        // 对象类别。怪物和玩家都有类别
	Name               string     // 对象名称
	Annotation         string     // 对象备注
	Level              int
	HP                 int // HP
	MaxHP              int // MaxHP
	AddHP              int
	ScriptMaxHP        int
	FillHP             int
	FillHPMax          int
	MP                 int // MP
	MaxMP              int // MaxMP
	AddMP              int
	SD                 int // SD
	MaxSD              int
	AddSD              int
	AG                 int // AG
	MaxAG              int
	AddAG              int
	targetNumber       int
	attackPanelMin     int // 物攻min
	attackPanelMax     int // 物攻max
	attackSpeed        int // 物攻速度
	attackRate         int // 攻击率
	defense            int // 防御力
	defenseRate        int // 防御率
	successfulBlocking int // 防御率
	magicDefense       int // 魔法防御率
	moveSpeed          int // 移动速度
	attackRange        int // 攻击范围
	attackType         int // 攻击类型
	viewRange          int // 视野范围
	itemDropRate       int // 道具掉落率
	moneyDropRate      int // 金钱掉落率
	attribute          int
	dieRegen           bool
	// regenOK                   byte
	regenTime                 time.Duration // 重生时间
	maxRegenTime              time.Duration // 最大重生时间
	pentagramMainAttribute    int
	pentagramAttributePattern int
	pentagramAttackMin        int
	pentagramAttackMax        int
	pentagramAttackRate       int
	pentagramDefense          int
	skills                    map[int]*skill.Skill
	FrustrumX                 [MaxArrayFrustrum]int
	FrustrumY                 [MaxArrayFrustrum]int
	viewports                 [MaxViewportNum]*viewport // 主动视野
	viewportNum               int
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
	// SkillAttackTime time.Duration
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

func (obj *object) init() {
	obj.targetNumber = -1
	obj.initSkill()
	obj.initViewport()
	obj.initMessage()
}

func (obj *object) reset() {
	obj.clearSkill()
	obj.clearViewport()
}

func (obj *object) initMessage() {
	for i := range obj.msgs {
		obj.msgs[i] = &messageStateMachine{
			code: -1,
		}
	}
}

func (obj *object) Test(msg *model.MsgTest) {
	obj.push(msg)
}

func (obj *object) randPosition(number, x1, y1, x2, y2 int) (int, int) {
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
	// panic(fmt.Sprintf("randPosition failed [number]%d", number))
	log.Printf("randPosition failed [map]%d [start](%d,%d) [end](%d,%d)\n", number, x1, y1, x2, y2)
	return x1, y1
}

func (obj *object) processRegen() {
	if !obj.dieRegen {
		return
	}
	if obj.ConnectState < ConnectStatePlaying {
		return
	}
	if time.Now().Unix()-int64(obj.regenTime) < int64(obj.maxRegenTime) {
		return
	}
	obj.HP = obj.MaxHP + obj.AddHP
	obj.MP = obj.MaxMP + obj.AddMP
	obj.Live = true
	obj.spawnPosition()
	obj.dieRegen = false
	obj.State = 1
}
