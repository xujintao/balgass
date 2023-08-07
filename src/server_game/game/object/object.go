package object

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/xujintao/balgass/src/server_game/conf"
	"github.com/xujintao/balgass/src/server_game/game/item"
	"github.com/xujintao/balgass/src/server_game/game/maps"
	"github.com/xujintao/balgass/src/server_game/game/math2"
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
	objects        []iobject
}

type iobject interface {
	process300ms()
	process500ms()
	process1000ms()
	processRegen()
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
	m.objects = make([]iobject, m.maxObjectCount)

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
					monster, err := m.AddMonster(spawn.Index)
					if err != nil {
						log.Fatalf("AddMonster failed err[%v]", err)
					}
					monster.MapNumber = _map.Number
					switch spot.Type {
					case 0: // npc
						monster.spawnStartX = spawn.StartX
						monster.spawnStartY = spawn.StartY
						monster.spawnEndX = spawn.StartX
						monster.spawnEndY = spawn.StartY
					case 1, 3: // multiple
						monster.spawnStartX = spawn.StartX
						monster.spawnStartY = spawn.StartY
						monster.spawnEndX = spawn.EndX
						monster.spawnEndY = spawn.EndY
					case 2: // single
						monster.spawnStartX = spawn.StartX - 3
						monster.spawnStartY = spawn.StartY - 3
						monster.spawnEndX = spawn.StartX + 3
						monster.spawnEndY = spawn.StartY + 3
					}
					monster.spawnDir = spawn.Dir
					monster.spawnDis = spawn.Distance
					monster.spawnPosition()
					if spot.Type == 3 {
						monster.pentagramMainAttribute = spawn.Element
					}
				}
			}
		}
	}
}

func (m *objectManager) AddMonster(kind int) (*Monster, error) {
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
	monster := NewMonster(kind)
	monster.objectManager = m
	monster.index = index
	m.objects[index] = monster
	return monster, nil
}

func (m *objectManager) AddCallMonster(kind int) (*Monster, error) {
	if m.callMonsterCount > m.maxCallMonsterCount {
		return nil, fmt.Errorf("over max call monster count")
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
	monster := NewMonster(kind)
	monster.objectManager = m
	monster.index = index
	m.objects[index] = monster
	return monster, nil
}

func (m *objectManager) AddPlayer(conn Conn) (int, error) {
	// limit max player count
	if m.playerCount > m.maxPlayerCount {
		// reply
		msg := model.MsgConnectFailed{Result: 4}
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
	player := NewPlayer(conn)
	player.objectManager = m
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
	}
	return obj
}

func (m *objectManager) Process300ms() {
	for _, v := range m.objects {
		if v == nil {
			continue
		}
		v.process300ms()
	}
}

func (m *objectManager) Process500ms() {
	for _, v := range m.objects {
		if v == nil {
			continue
		}
		v.process500ms()
	}
}

func (m *objectManager) Process1000ms() {
	for _, v := range m.objects {
		if v == nil {
			continue
		}
		v.process1000ms()
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
	state  int // 3消失
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

type object struct {
	*objectManager
	index        int
	ConnectState ConnectState
	Live         bool
	State        int // 1:初始 2:视野 4:死亡

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
	attackDamageMin    int // 物攻min
	attackDamageMax    int // 物攻max
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
	viewports2                [MaxViewportNum]*viewport // 被动视野
	viewportNum               int
	viewportNum2              int
	msgs                      [20]*messageStateMachine

	groupNumber     int
	subGroupNumber  int
	groupMemberGUID int
	regenType       int
	// argo                    *monster.MonsterAIAgro
	lastAutoRunTime            time.Time
	lastAutoDelay              time.Duration
	expType                    int
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
	Authority                  uint
	AuthorityCode              uint
	Penalty                    uint
	GameMaster                 uint
	PenaltyMask                uint
	ChatBlockTime              time.Time
	AccountItemBlock           byte
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
	targetNumber               int
	TargetNpcNumber            int
	LastAttackerID             int
	criticalDamage             int
	excellentDamage            int // 卓越一击概率
	// magicBack                   *skill.MagicInfo
	// Magic                       *skill.MagicInfo
	UseMagicNumber  byte
	UseMagicTime    time.Duration
	UseMagicCount   byte
	OSAttackSerial  uint16
	SASCount        byte
	SkillAttackTime time.Duration
	CharSet         string
	// resistance               [MaxResistanceType]byte
	// addResistance            [MaxResistanceType]byte
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
	obj.targetNumber = -1
	obj.initSkill()
	obj.initViewport()
	obj.initMessage()
}

func (obj *object) reset() {
	obj.clearSkill()
	obj.clearViewport()
}

func (obj *object) initSkill() {
	obj.skills = make(map[int]*skill.Skill)
}

// AddSkill  object add skill
func (obj *object) addSkill(index, level int) bool {
	if _, ok := obj.skills[index]; ok {
		log.Printf("[object]%s [skill]%d already exists", obj.Name, index)
		return false
	}
	obj.skills[index] = skill.SkillManager.Get(index, level, obj.skills)
	return true
}

func (obj *object) clearSkill() {
	obj.skills = nil
}

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

func (obj *object) initMessage() {
	for i := range obj.msgs {
		obj.msgs[i] = &messageStateMachine{
			code: -1,
		}
	}
}

func (obj *object) calcDistance(tobj *object) int {
	x := obj.X - tobj.X
	y := obj.Y - tobj.Y
	if x == 0 && y == 0 {
		return 0
	}
	return int(math.Sqrt(float64(x*x + y*y)))
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
		tobj := obj.objectManager.object(v)
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
			tobj := obj.objectManager.object(obj.objectManager.objects[tnum])
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
		tobj := obj.objectManager.object(obj.objectManager.objects[vp.number])
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

func (obj *object) process300ms() {
	if obj.ConnectState < ConnectStatePlaying ||
		!obj.Live ||
		obj.State != 2 ||
		obj.pathCount == 0 {
		return
	}
	moveTime := obj.moveSpeed
	if obj.delayLevel != 0 {
		moveTime += 300
	}
	pathTime := time.Now().UnixMilli()
	if pathTime-obj.pathTime+1 < int64(moveTime) {
		return
	}
	obj.pathTime = pathTime
	x := obj.pathX[obj.pathCur]
	y := obj.pathY[obj.pathCur]
	dir := obj.pathDir[obj.pathCur]
	attr := maps.MapManager.GetMapAttr(obj.MapNumber, x, y)
	if attr&4 != 0 && attr&8 != 0 {
		log.Printf("process300ms object move check [index]%d [class]%d [map]%d [position](%d,%d)",
			obj.index, obj.Class, obj.MapNumber, x, y)
		for i := 0; i < len(obj.pathDir); i++ {
			obj.pathX[i] = 0
			obj.pathY[i] = 0
			obj.pathDir[i] = 0
		}
		obj.pathCount = 0
		obj.pathCur = 0
		obj.pathMoving = false
		return
	}
	obj.X = x
	obj.Y = y
	// if obj.index == 6 && obj.pathMoving {
	// 	fmt.Println(obj.X, obj.Y)
	// }
	obj.Dir = dir
	obj.createFrustrum()
	obj.pathCur++
	if obj.pathCur >= obj.pathCount {
		for i := 0; i < len(obj.pathDir); i++ {
			obj.pathX[i] = 0
			obj.pathY[i] = 0
			obj.pathDir[i] = 0
		}
		obj.pathCount = 0
		obj.pathCur = 0
		obj.pathMoving = false
	}
}

func (obj *object) Move(msg *model.MsgMove) {
	n := len(msg.Path)
	if n < 1 || n > 15 {
		return
	}
	for i := range msg.Path {
		obj.pathX[i] = msg.Path[i].X
		obj.pathY[i] = msg.Path[i].Y
		obj.pathDir[i] = msg.Dirs[i]
	}
	obj.pathCount = n
	obj.pathCur = 0
	obj.pathMoving = true
	maps.MapManager.ClearMapAttrStand(obj.MapNumber, obj.TX, obj.TX)
	obj.TX = msg.Path[n-1].X
	obj.TY = msg.Path[n-1].Y
	maps.MapManager.SetMapAttrStand(obj.MapNumber, obj.TX, obj.TY)
	// if obj.index == 6 {
	// 	fmt.Printf("(%d,%d)->(%d,%d)\n", obj.X, obj.Y, obj.TX, obj.TY)
	// }

	msgRelpy := model.MsgMoveReply{
		Number: obj.index,
		X:      obj.TX,
		Y:      obj.TY,
		Dir:    msg.Dirs[0] << 4,
	}
	if obj.Type == ObjectTypePlayer {
		om := obj.objectManager
		tobj := om.objects[obj.index]
		p := tobj.(*Player)
		p.Push(&msgRelpy)
	}
	for _, vp := range obj.viewports2 {
		if vp.state != 1 && vp.state != 2 {
			continue
		}
		tnum := vp.number
		if tnum < 0 {
			continue
		}
		om := obj.objectManager
		tobj := om.objects[tnum]
		p, ok := tobj.(*Player)
		if !ok {
			continue
		}
		if p.ConnectState == ConnectStatePlaying && p.Live {
			p.Push(&msgRelpy)
		}
	}
}

func (obj *object) Attack(msg *model.MsgAttack) {

}

func (obj *object) SkillAttack(msg *model.MsgSkillAttack) {

}
