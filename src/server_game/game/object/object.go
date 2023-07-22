package object

import (
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"sync"
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
	maxObjectCount int
	objects        []iobject
	// monsterStartIndex int
	// lastMonsterIndex  int
	monsterCount int
	// summonMonsterCount int
	playerStartIndex int
	lastPlayerIndex  int
	playerCount      int
}

type iobject interface {
	reset()
	addSkill(int, int) bool
	processRegen()
}

func (m *objectManager) init() {
	m.maxObjectCount = conf.Server.GameServerInfo.MaxMonsterCount +
		conf.Server.GameServerInfo.MaxSummonMonsterCount +
		conf.Server.GameServerInfo.MaxPlayerCount
	m.objects = make([]iobject, m.maxObjectCount)
	// objectBills = make([]bill, conf.Server.MaxPlayerCount)
	m.playerStartIndex = conf.Server.GameServerInfo.MaxMonsterCount +
		conf.Server.GameServerInfo.MaxSummonMonsterCount
	m.lastPlayerIndex = m.playerStartIndex
	// 先有怪后有玩家
	m.spawnMonster()
}

func (m *objectManager) objectMaxRange(index int) bool {
	if index < 0 || index >= m.maxObjectCount {
		return false
	}
	return true
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
				Spawn       []struct {
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

	randPosition := func(number, x1, y1, x2, y2 int) (int, int) {
		w := x2 - x1
		if w <= 0 {
			w = 1
		}
		h := y2 - y1
		if h <= 0 {
			h = 1
		}
		for i := 0; i < 100; i++ {
			x := x1 + rand.Intn(w)
			y := y1 + rand.Intn(h)
			attr := maps.MapManager.GetMapAttr(number, x, y)
			if attr&1 == 0 && attr&4 == 0 && attr&8 == 0 {
				return x, y
			}
		}
		panic(fmt.Sprintf("randPosition failed [number]%d", number))
	}
	for _, _map := range monsterSpawn.Map {
		for _, spot := range _map.Spot {
			for _, spawn := range spot.Spawn {
				cnt := spawn.Count
				if cnt == 0 {
					cnt = 1
				}
				for i := 0; i < cnt; i++ {
					// wrong
					if _map.Number == maps.Atlans && spawn.StartX == 251 && spawn.StartY == 51 ||
						_map.Number == maps.Atlans && spawn.StartX == 7 && spawn.StartY == 52 ||
						_map.Number == maps.LandOfTrial && spawn.StartX == 14 && spawn.StartY == 43 ||
						_map.Number == maps.KanturuBoss && spawn.Index == 106 {
						continue
					}
					monster, err := m.AddMonster(spawn.Index)
					if err != nil {
						log.Fatalf("AddMonster failed err[%v]", err)
					}
					monster.MapNumber = _map.Number
					switch spot.Type {
					case 0: // npc
						monster.StartX, monster.StartY = spawn.StartX, spawn.StartY
					case 1, 3: // multiple
						monster.StartX, monster.StartY = randPosition(_map.Number, spawn.StartX, spawn.StartY, spawn.EndX, spawn.EndY)
					case 2: // single
						monster.StartX, monster.StartY = randPosition(_map.Number, spawn.StartX-3, spawn.StartY-3, spawn.StartX+3, spawn.StartY+3)
					}
					monster.X, monster.Y = monster.StartX, monster.StartY
					maps.MapManager.SetMapAttrStand(monster.MapNumber, monster.X, monster.Y)
					monster.OldX = monster.X
					monster.OldY = monster.Y
					monster.Dir = spawn.Dir
					if monster.Dir < 0 {
						monster.Dir = rand.Intn(8)
					}
					if spot.Type == 3 {
						monster.pentagramMainAttribute = spawn.Element
					}
				}
			}
		}
	}
}

func (m *objectManager) AddMonster(kind int) (*Monster, error) {
	if m.monsterCount >= conf.Server.GameServerInfo.MaxMonsterCount {
		return nil, fmt.Errorf("over max monster count")
	}
	index := m.monsterCount
	cnt := conf.Server.GameServerInfo.MaxMonsterCount
	for cnt > 0 {
		if m.objects[index] == nil {
			break
		}
		index++
		if index >= conf.Server.GameServerInfo.MaxMonsterCount {
			index = 0
		}
		cnt--
	}
	if cnt == 0 {
		panic(fmt.Errorf("have no free monster index"))
	}
	m.monsterCount++
	monster := NewMonster(kind)
	monster.index = index
	m.objects[index] = monster
	return monster, nil
}

func (m *objectManager) AddPlayer(conn Conn) (int, error) {
	// limit max player count
	if m.playerCount >= conf.Server.GameServerInfo.MaxPlayerCount {
		// reply
		msg := model.MsgConnectFailed{Result: 4}
		conn.Write(&msg)
		return -1, fmt.Errorf("over max player count")
	}

	// get unified object index
	index := m.lastPlayerIndex
	cnt := conf.Server.GameServerInfo.MaxPlayerCount
	for cnt > 0 {
		if m.objects[index] == nil {
			break
		}
		index++
		if index >= m.maxObjectCount {
			index = m.playerStartIndex
		}
		cnt--
	}
	if cnt == 0 {
		panic(fmt.Errorf("have no free player index"))
	}
	m.lastPlayerIndex = index
	m.playerCount++
	player := NewPlayer(conn)
	player.index = index
	// register the new player to object manager
	m.objects[index] = player

	// reply
	msg := model.MsgConnectSuccess{
		Result:  1,
		ID:      index,
		Version: conf.MapServers.ServerInfo.Version,
	}
	player.Push(&msg)
	log.Printf("player online [id]%d [addr]%s", player.index, player.conn.Addr())
	return index, nil
}

func (m *objectManager) DeletePlayer(id int) {
	player := m.objects[id].(*Player)
	// unregister player from object manager
	m.objects[id] = nil
	m.playerCount--
	player.delete()
	log.Printf("player offline [id]%d [addr]%s", player.index, player.conn.Addr())
}

func (m *objectManager) GetPlayer(id int) *Player {
	return m.objects[id].(*Player)
}

func (m *objectManager) object(v iobject) *object {
	var obj *object
	if monster, ok := v.(*Monster); ok {
		obj = &monster.object
	} else if player, ok := v.(*Player); ok {
		obj = &player.object
	} else {
		panic(fmt.Sprintf("object failed [err]%v", v))
	}
	return obj
}

func (m *objectManager) ProcessRegen() {
	for _, v := range m.objects {
		if v == nil {
			continue
		}
		v.processRegen()
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

type ActionState struct {
	Rest         byte
	Attack       byte
	Move         byte
	Escape       byte
	Emotion      byte
	EmotionCount byte
}

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

type ViewPort struct {
	State  uint8
	Number uint16
	Type   uint8
	Index  uint16
	Dis    int
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
	ObjectEmpty ObjectType = iota - 1
	ObjectPlayer
	ObjectMonster
	ObjectNPC
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

type object struct {
	index        int
	ConnectState ConnectState
	Live         bool
	// 4:等待重生
	State                     int
	StartX                    int
	StartY                    int
	X                         int // x坐标
	Y                         int // y坐标
	OldX                      int
	OldY                      int
	Dir                       int        // 方向
	MapNumber                 int        // 地图号
	Type                      ObjectType // 对象种类：玩家，怪物，NPC
	Class                     int        // 对象类别。怪物和玩家都有类别
	Name                      string     // 对象名称
	Annotation                string     // 对象备注
	Level                     int
	HP                        int // HP
	MaxHP                     int // MaxHP
	AddHP                     int
	ScriptMaxHP               int
	FillHP                    int
	FillHPMax                 int
	MP                        int // MP
	MaxMP                     int // MaxMP
	AddMP                     int
	SD                        int // SD
	MaxSD                     int
	AddSD                     int
	AG                        int // AG
	MaxAG                     int
	AddAG                     int
	attackDamageMin           int // 物攻min
	attackDamageMax           int // 物攻max
	attackSpeed               int // 物攻速度
	attackRate                int // 攻击率
	defense                   int // 防御力
	defenseRate               int // 防御率
	successfulBlocking        int // 防御率
	magicDefense              int // 魔法防御率
	moveSpeed                 int // 移动速度
	moveRange                 int // 移动范围
	attackRange               int // 攻击范围
	attackType                int // 攻击类型
	viewRange                 int // 视野范围
	itemDropRate              int // 道具掉落率
	moneyDropRate             int // 金钱掉落率
	attribute                 int
	dieRegen                  bool
	regenOK                   byte
	regenTime                 time.Duration // 重生时间
	maxRegenTime              time.Duration // 最大重生时间
	pentagramMainAttribute    int
	pentagramAttributePattern int
	pentagramAttackMin        int
	pentagramAttackMax        int
	pentagramAttackRate       int
	pentagramDefense          int

	basicAI         int
	currentAI       int
	currentAIState  int
	lastAIRunTime   time.Duration
	groupNumber     int
	subGroupNumber  int
	groupMemberGUID int
	regenType       int
	// argo                    *monster.MonsterAIAgro
	lastAutoRunTime time.Time
	lastAutoDelay   time.Duration
	expType         int

	attackDamageLeft           int // 物攻左
	attackDamageRight          int // 物攻右
	attackDamageLeftMin        int // 物攻左min
	attackDamageLeftMax        int // 物攻左min
	attackDamageRightMin       int // 物攻右min
	attackDamageRightMax       int // 物攻右max
	magicDamageMin             int // 魔攻min
	magicDamageMax             int // 魔攻max
	magicSpeed                 int // 魔攻速度
	curseDamageMin             int // 诅咒min
	curseDamageMax             int // 诅咒max
	curseSpell                 int
	LoginMsgSend               bool
	LoginMsgCount              byte
	CloseCount                 byte
	CloseType                  byte
	EnableCharacterDel         bool
	UserNumber                 int
	DBNumber                   int
	EnableCharacterCreate      bool
	AutoSaveTime               time.Time
	ConnectCheckTime           time.Time
	CheckTick                  uint
	CheckSpeedHack             bool
	CheckTick2                 uint
	CheckTickCount             byte
	PintTime                   int
	TimeCount                  byte
	PKTimer                    *time.Timer
	CheckSumTableNum           uint16
	CheckSumTime               uint
	Leadership                 int
	AddLeadership              int
	ChatLimitTime              uint16
	ChatLimitTimeSec           byte
	FillLifeCount              byte
	AddStrength                int
	AddDexterity               int
	AddVitality                int
	AddEnergy                  int
	VitalityToLife             float32
	EnergyToMana               float32
	PKCount                    int
	PKLevel                    byte
	PKTime                     int
	PKTotalCount               int
	XSave                      uint16
	YSave                      uint16
	MapNumberSave              byte
	XDie                       uint16
	YDie                       uint16
	MapNumberDie               byte
	IFillShieldMax             int
	IFillShield                int
	IFillShieldCount           int
	ShieldAutoRefillTimer      *time.Timer
	DamageMinus                int  // 伤害减少
	DamageReflect              int  // 伤害反射
	MonsterDieGetMoney         int  // 杀怪加钱
	MonsterDieGetLife          byte // 杀怪回生
	MonsterDieGetMana          byte // 杀怪回蓝
	AutoHPRecovery             byte // 自动生命恢复
	TX                         int
	TY                         int
	MTX                        uint16
	MTY                        uint16
	PathCount                  int
	PathCur                    int
	PathStartEnd               byte
	PathOri                    [15]uint16
	PathX                      [15]uint16
	PathY                      [15]uint16
	PathDir                    [15]byte
	PathTime                   uint
	Authority                  uint
	AuthorityCode              uint
	Penalty                    uint
	GameMaster                 uint
	PenaltyMask                uint
	ChatBlockTime              time.Time
	AccountItemBlock           byte
	ActState                   ActionState
	ActionNumber               byte
	ActionTime                 uint
	ActionCount                byte
	ChatFloodTime              uint
	ChatFloodCount             byte
	Rest                       byte
	viewState                  byte
	buffEffectCount            byte
	buffEffectList             [MaxBuffEffect]EffectList
	lastMoveTime               uint
	lastAttackTime             uint
	friendServerOnline         byte
	detectSpeedHackTime        time.Duration
	sumLastAttackTime          time.Duration
	detectCount                uint
	detectHackKickCount        uint
	speedHackPenalty           int
	attackSpeedHackDetectCount uint
	packetCheckTime            time.Duration
	ShopTime                   time.Duration
	totalAttackTime            time.Duration
	totalAttackCount           uint
	TeleportTIme               time.Duration
	Teleport                   byte
	KillerType                 byte
	MapNumberRegen             byte
	XRegen                     byte
	YRegen                     byte
	posNum                     uint16
	LifeRefillTimer            *time.Timer
	ActionTimeCur              time.Time
	ActionTimeNext             time.Time
	ActionTimeDelay            time.Duration
	DelayLevel                 byte
	monsterBattleDelay         byte
	kalimaGateExist            byte
	kalimaGateIndex            int
	kalimaGateEnterCount       byte
	AttackObj                  *object
	skillNumber                uint16
	skillTime                  time.Duration
	attackerKilled             bool
	manaFillCount              byte
	lifeFillCount              byte
	SelfDefense                [MaxSelfDefense]int
	SelfDefenseTime            time.Duration
	PartyNumber                int
	PartyTargetUser            int
	Married                    byte
	MarryName                  string
	MarryRequested             byte
	WinDuels                   int
	LoseDuels                  int
	MarryRequestIndex          uint16
	MarryRequestTime           time.Duration
	RecallMon                  int
	change                     int
	TargetNumber               int
	TargetNpcNumber            int
	LastAttackerID             int
	criticalDamage             int
	excellentDamage            int // 卓越一击概率
	// magicBack                   *skill.MagicInfo
	// Magic                       *skill.MagicInfo
	Skills          map[int]*skill.Skill
	UseMagicNumber  byte
	UseMagicTime    time.Duration
	UseMagicCount   byte
	OSAttackSerial  uint16
	SASCount        byte
	SkillAttackTime time.Duration
	CharSet         string
	// resistance               [MaxResistanceType]byte
	// addResistance            [MaxResistanceType]byte
	FrustrumX                [MaxArrayFrustrum]int
	FrustrumY                [MaxArrayFrustrum]int
	VPPlayer                 *ViewPort
	VPPlayer2                *ViewPort
	VPCount                  int
	VPCount2                 int
	HD                       *HitDamage
	HDCount                  uint16
	ifState                  InterfaceState
	iterfaceTime             time.Duration
	Inventory                []item.Item
	InventoryMap             *uint8
	InventoryCount           *uint8
	Transaction              uint8
	Inventory1               *item.Item
	InventoryMap1            *uint8
	InventoryCount1          uint8
	Inventory2               *item.Item
	InventoryMap2            *uint8
	InventoryCount2          uint8
	Trade                    *item.Item
	TradeMap                 *uint8
	TradeMoney               int
	TradeOK                  bool
	Warehouse                *item.Item
	WarehouseID              uint8
	WarehouseTick            time.Time
	WarehouseMap             *uint8
	WarehouseCount           uint8
	WarehousePW              uint16
	WarehouseLock            uint8
	WarehouseUnfailLock      uint8
	WarehouseMoney           int
	ChaosBox                 *item.Item
	ChaosBoxMap              *uint8
	ChaosMoney               int
	ChaosSuccessRate         int
	ChaosMassMixCurItem      uint8
	ChaosMassMixSuccessCount uint8
	ChaosLock                bool
	Option                   uint
	eventScore               int
	eventExp                 int
	eventMoney               int
	devilSquareIndex         uint8
	devilSquareAuth          bool
	BloodCastlIndex          uint8
	BloodCastleSubIndex      uint8
	BloodCastleExp           int
	BloodCastleComplete      bool
	ChaosCastleIndex         uint8
	ChaosCastleSubIndex      uint8
	ChaosCastleBlowTime      time.Duration
	isCCFUIUsing             bool
	CCFCanEnter              uint8
	CCFCertiTick             time.Time
	CCFUserIndex             int
	CCFBlowTime              time.Time
	killUserCount            uint8
	killMobCount             uint8
	isCCFQuitMsg             bool
	illusionTempleIndex      uint8
	zoneIndex                uint8
	ckillUserCount           uint8
	cKillMonsterCount        uint8
	duelUserReserved         int
	duelUserRequested        int
	duelUser                 int
	duelRoom                 int
	duelScore                uint8
	duelTickCount            time.Duration
	IsInBattleGround         bool
	HaveWeaponInHand         bool
	EventChipCount           uint16
	LuckyCoinCount           int
	MutoNumber               int
	UseEventServer           bool
	LoadWarehouseInfo        bool
	StoneCount               int
	maxLifePower             int
	checkLifeTime            int
	moveToOtherServer        uint8
	BackName                 string
	isPShopOpen              bool
	isPShopTransaction       bool
	isPShopItemChange        bool
	isPShopRedrawABS         bool
	PShopText                string
	isPShopWantDeal          bool
	PShopDealerIndex         int
	PShopDealerName          string
	muPShopTrade             sync.Mutex
	VPPShopPlayer            [MaxViewPort]int
	VPPShopPlayerCount       uint16
	BossGoldDerconMapNumber  uint8
	lastTeleportTime         time.Time
	clientHackLogCount       uint8
	isInMonsterHerd          bool
	isMonsterAttackFirst     bool
	// monsterHerd              *monster.MonsterHerd
	fsKillFrustrumX [MaxArrayFrustrum]int
	fsKillFrustrumY [MaxArrayFrustrum]int
	// durMagicKeyChecker          *skill.DurMagicKeyChecker
	IsChaosMixCompleted         bool
	SkillLongSpearChange        bool
	objectSecTimer              time.Timer
	mapSvrMoveQuit              bool
	mapSvrMoveReq               bool
	mapSvrMoveReq2              bool
	mapSvrQuitTick              time.Time
	prevMapSvrCode              uint16
	destMapNumber               uint16
	destX                       uint8
	destY                       uint8
	csNpcExistVal               int
	csNpcType                   uint8
	csGateOPen                  uint8
	csGateLeverLinkIndex        int
	csNpcDfLevel                uint8
	csNpcRgLevel                uint8
	csJoinSide                  uint8
	csGuildInvolved             bool
	IsCastleNPCUpgradeCompleted bool
	weaponState                 uint8
	weaponUser                  int
	killCount                   uint8
	accumulatedDamage           int
	lifeStoneCount              uint8
	creationState               uint8
	createdActiviationTime      int
	accumulatedCrownAccessTime  int
	// monsterSkillElementInfo     monster.MonsterSkillElementInfo
	crywolfMVPScore         int
	lastCheckTick           time.Time
	autoRecuperationTime    time.Time
	skillDistanceErrorCount int
	skillDistanceErrorTick  time.Time
	skillInfo               skillInfo
	bufferIndex             int
	buffID                  int
	buffPlayerIndex         int
	agiCheckTime            time.Time
	warehouseSaveLock       uint8
	crcCheckTime            time.Time
	off                     bool
	offLevel                bool
	offLevelTime            int
}

func (obj *object) init() {
	obj.Skills = make(map[int]*skill.Skill)
}

func (obj *object) reset() {
	obj.Skills = nil
}

// AddSkill  object add skill
func (obj *object) addSkill(skillIndex, level int) bool {
	if _, ok := obj.Skills[skillIndex]; ok {
		log.Printf("object[%s] skill[%s] already exists", obj.Name, skill.SkillTable[skillIndex].Name)
		return false
	}
	obj.Skills[skillIndex] = skill.Get(skillIndex, level)
	return true
}
