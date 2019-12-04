package user

import (
	"fmt"
	"sync"
	"time"

	"github.com/xujintao/balgass/cmd/server_game/guild"

	"github.com/xujintao/balgass/cmd/server_game/monster/ai"

	"github.com/xujintao/balgass/cmd/server_game/conf"
	"github.com/xujintao/balgass/cmd/server_game/item"
	"github.com/xujintao/balgass/cmd/server_game/monster/kalima/herd"
	mskill "github.com/xujintao/balgass/cmd/server_game/monster/skill"
	"github.com/xujintao/balgass/cmd/server_game/network"
	"github.com/xujintao/balgass/cmd/server_game/skill"
)

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
	MaxResistanceType       = 7
	MaxViewPort             = 75
	MaxArrayFrustrum        = 4
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
	ObjectUser
	ObjectMonster
	ObjectNPC
)

type PlayerType int

const (
	PlayerEmpty PlayerType = iota
	PlayerConnected
	PlayerLogged
	PlayerPlaying
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

type UserData struct {
	ip                            string
	hwid                          string
	experience                    uint
	nextExp                       uint
	masterExperience              uint
	nextMasterExp                 uint
	masterLevel                   uint16
	levelUpPoint                  int
	masterPoint                   int
	masterPointUsed               int
	fruitPoint                    int
	money                         int
	strength                      uint16
	dexterity                     uint16
	vitality                      uint16
	energy                        uint16
	dbClass                       uint8
	changeUP                      uint8
	guild                         *guild.GuildInfo
	guildName                     string
	guildStatus                   int
	guildUnionTimeStamp           int
	guildNumber                   int
	lastMoveTime                  time.Time
	resets                        int
	vipType                       uint8
	vipEffect                     uint8
	santaCount                    uint8
	goblinTime                    time.Time
	securityCheck                 bool
	securityCode                  int
	RegisterLMS                   uint8
	registerLMSRoom               uint8
	jewelHarmonyEffect            item.JewelHarmonyItemEffect
	itemOptionExFor380            item.ItemOptionFor380ItemEffect
	kanturuEntranceByNPC          bool
	gensInfoLoad                  bool
	questInfoLoad                 bool
	wCoinP                        int
	wCoinC                        int
	goblinPoint                   int
	periodItemEffectIndex         int
	seedOptionList                [35]item.SocketOptionList
	bonusOptionList               [7]item.SocketOptionList
	setOptionList                 [2]item.SocketOptionList
	refillHPSocketOption          uint16
	refillMPSocketOption          uint16
	socketOptionMonsterDieGetHP   uint16
	socketOptionMonsterDieGetMana uint16
	agReduceRate                  uint8
	muBotEnable                   bool
	muBotTotalTime                time.Duration
}

type Object struct {
	index                       int
	Connected                   PlayerType
	LoginMsgSend                bool
	LoginMsgCount               byte
	CloseCount                  byte
	CloseTYpe                   byte
	EnableCharacterDel          bool
	conn                        *network.Conn
	UserNumber                  int
	DBNumber                    int
	EnableCharacterCreate       bool
	AutoSaveTime                time.Time
	ConnectCheckTime            time.Time
	CheckTick                   uint
	CheckSpeedHack              bool
	CheckTick2                  uint
	CheckTickCount              byte
	PintTime                    int
	TimeCount                   byte
	PKTimer                     *time.Timer
	CheckSumTableNum            uint16
	CheckSumTime                uint
	Type                        ObjectType
	Live                        byte
	AccountID                   string
	Name                        string
	Class                       uint16
	Level                       uint16
	Life                        float32
	MaxLife                     float32
	ScriptMaxLife               int
	FillLife                    float32
	FillLifeMax                 float32
	Mana                        float32
	MaxMana                     float32
	Leadership                  uint16
	AddLeadership               uint16
	ChatLimitTime               uint16
	ChatLimitTimeSec            byte
	FillLifeCount               byte
	AddStrength                 int
	AddDexterity                int
	AddVitality                 int
	AddEnergy                   int
	BP                          int
	MaxBP                       int
	VitalityToLife              float32
	EnergyToMana                float32
	PKCount                     int
	PKLevel                     byte
	PKTime                      int
	PKTotalCount                int
	X                           uint16
	Y                           uint16
	Dir                         byte
	MapNumber                   byte
	XSave                       uint16
	YSave                       uint16
	MapNumberSave               byte
	XDie                        uint16
	YDie                        uint16
	MapNumberDie                byte
	AddLife                     int
	AddMana                     int
	IShield                     int
	IShieldMax                  int
	IAddShield                  int
	IFillShieldMax              int
	IFillShield                 int
	IFillShieldCount            int
	ShieldAutoRefillTimer       *time.Timer
	DamageMinus                 byte   // 伤害减少
	DamageReflect               byte   // 伤害反射
	MonsterDieGetMoney          uint16 // 杀怪加钱
	MonsterDieGetLife           byte   // 杀怪回生
	MonsterDieGetMana           byte   // 杀怪回蓝
	AutoHPRecovery              byte   // 自动生命恢复
	XStart                      byte
	YStart                      byte
	XOld                        uint16
	YOld                        uint16
	TX                          uint16
	TY                          uint16
	MTX                         uint16
	MTY                         uint16
	PathCount                   int
	PathCur                     int
	PathStartEnd                byte
	PathOri                     [15]uint16
	PathX                       [15]uint16
	PathY                       [15]uint16
	PathDir                     [15]byte
	PathTime                    uint
	Authority                   uint
	AuthorityCode               uint
	Penalty                     uint
	GameMaster                  uint
	PenaltyMask                 uint
	ChatBlockTime               time.Time
	AccountItemBlock            byte
	ActState                    ActionState
	ActionNumber                byte
	ActionTime                  uint
	ActionCount                 byte
	ChatFloodTime               uint
	ChatFloodCount              byte
	State                       uint
	Rest                        byte
	viewState                   byte
	buffEffectCount             byte
	buffEffectList              [MaxBuffEffect]EffectList
	lastMoveTime                uint
	lastAttackTime              uint
	friendServerOnline          byte
	detectSpeedHackTime         time.Duration
	sumLastAttackTime           time.Duration
	detectCount                 uint
	detectHackKickCount         uint
	speedHackPenalty            int
	attackSpeedHackDetectCount  uint
	packetCheckTime             time.Duration
	ShopTime                    time.Duration
	totalAttackTime             time.Duration
	totalAttackCount            uint
	TeleportTIme                time.Duration
	Teleport                    byte
	KillerType                  byte
	DieRegen                    byte
	RegenOK                     byte
	MapNumberRegen              byte
	XRegen                      byte
	YRegen                      byte
	RegenTime                   time.Duration
	RegenTimeMax                time.Duration
	posNum                      uint16
	LifeRefillTimer             *time.Timer
	ActionTimeCur               time.Time
	ActionTimeNext              time.Time
	ActionTimeDelay             time.Duration
	DelayLevel                  byte
	monsterBattleDelay          byte
	kalimaGateExist             byte
	kalimaGateIndex             int
	kalimaGateEnterCount        byte
	AttackObj                   *Object
	skillNumber                 uint16
	skillTime                   time.Duration
	attackerKilled              bool
	manaFillCount               byte
	lifeFillCount               byte
	SelfDefense                 [MaxSelfDefense]int
	SelfDefenseTime             time.Duration
	PartyNumber                 int
	PartyTargetUser             int
	Married                     byte
	MarryName                   string
	MarryRequested              byte
	WinDuels                    int
	LoseDuels                   int
	MarryRequestIndex           uint16
	MarryRequestTime            time.Duration
	recallmon                   int
	change                      int
	TargetNumber                uint16
	TargetNpcNumber             uint16
	LastAttackerID              uint16
	attackDamageMin             int
	attackDamageMax             int
	magicDamageMin              int
	magicDamageMax              int
	curseDamageMin              int
	curseDamageMax              int
	attackDamageLeft            int
	attackDamageRight           int
	attackDamageminLeft         int
	attackDamageMaxLeft         int
	attackDamageMinRight        int
	attackDamageMaxRight        int
	attackRating                int
	attackSpeed                 int
	magicSpeed                  int
	defense                     int
	magicDefense                int
	successfulBlocking          int
	curseSpell                  int
	moveSpeed                   uint16
	moveRange                   uint16
	attackRange                 uint16
	attackType                  uint16
	viewRange                   uint16
	attribute                   uint16
	itemRate                    uint16
	moneyRate                   uint16
	criticalDamage              int
	excellentDamage             int
	magicBack                   *skill.MagicInfo
	Magic                       *skill.MagicInfo
	MagicCount                  byte
	UseMagicNumber              byte
	UseMagicTime                time.Duration
	UseMagicCount               byte
	OSAttackSerial              uint16
	SASCount                    byte
	SkillAttackTime             time.Duration
	CharSet                     string
	resistance                  [MaxResistanceType]byte
	addResistance               [MaxResistanceType]byte
	FrustrumX                   [MaxArrayFrustrum]int
	FrustrumY                   [MaxArrayFrustrum]int
	VPPlayer                    *ViewPort
	VPPlayer2                   *ViewPort
	VPCount                     int
	VPCount2                    int
	HD                          *HitDamage
	HDCount                     uint16
	ifState                     InterfaceState
	iterfaceTime                time.Duration
	Inventory                   *item.Item
	InventoryMap                *uint8
	InventoryCount              *uint8
	Transaction                 uint8
	Inventory1                  *item.Item
	InventoryMap1               *uint8
	InventoryCount1             uint8
	Inventory2                  *item.Item
	InventoryMap2               *uint8
	InventoryCount2             uint8
	Trade                       *item.Item
	TradeMap                    *uint8
	TradeMoney                  int
	TradeOK                     bool
	Warehouse                   *item.Item
	WarehouseID                 uint8
	WarehouseTick               time.Time
	WarehouseMap                *uint8
	WarehouseCount              uint8
	WarehousePW                 uint16
	WarehouseLock               uint8
	WarehouseUnfailLock         uint8
	WarehouseMoney              int
	ChaosBox                    *item.Item
	ChaosBoxMap                 *uint8
	ChaosMoney                  int
	ChaosSuccessRate            int
	ChaosMassMixCurItem         uint8
	ChaosMassMixSuccessCount    uint8
	ChaosLock                   bool
	Option                      uint
	eventScore                  int
	eventExp                    int
	eventMoney                  int
	devilSquareIndex            uint8
	devilSquareAuth             bool
	BloodCastlIndex             uint8
	BloodCastleSubIndex         uint8
	BloodCastleExp              int
	BloodCastleComplete         bool
	ChaosCastleIndex            uint8
	ChaosCastleSubIndex         uint8
	ChaosCastleBlowTime         time.Duration
	isCCFUIUsing                bool
	CCFCanEnter                 uint8
	CCFCertiTick                time.Time
	CCFUserIndex                int
	CCFBlowTime                 time.Time
	killUserCount               uint8
	killMobCount                uint8
	isCCFQuitMsg                bool
	illusionTempleIndex         uint8
	zoneIndex                   uint8
	ckillUserCount              uint8
	cKillMonsterCount           uint8
	duelUserReserved            int
	duelUserRequested           int
	duelUser                    int
	duelRoom                    int
	duelScore                   uint8
	duelTickCount               time.Duration
	IsInBattleGround            bool
	HaveWeaponInHand            bool
	EventChipCount              uint16
	LuckyCoinCount              int
	MutoNumber                  int
	UseEventServer              bool
	LoadWarehouseInfo           bool
	StoneCount                  int
	maxLifePower                int
	checkLifeTime               int
	moveToOtherServer           uint8
	BackName                    string
	isPShopOpen                 bool
	isPShopTransaction          bool
	isPShopItemChange           bool
	isPShopRedrawABS            bool
	PShopText                   string
	isPShopWantDeal             bool
	PShopDealerIndex            int
	PShopDealerName             string
	muPShopTrade                sync.Mutex
	VPPShopPlayer               [MaxViewPort]int
	VPPShopPlayerCount          uint16
	BossGoldDerconMapNumber     uint8
	lastTeleportTime            time.Time
	clientHackLogCount          uint8
	isInMonsterHerd             bool
	isMonsterAttackFirst        bool
	monsterHerd                 *herd.MonsterHerd
	fsKillFrustrumX             [MaxArrayFrustrum]int
	fsKillFrustrumY             [MaxArrayFrustrum]int
	durMagicKeyChecker          *skill.DurMagicKeyChecker
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
	monsterSkillElementInfo     mskill.MonsterSkillElementInfo
	basicAI                     int
	currentAI                   int
	currentAIState              int
	lastAIRunTime               time.Duration
	groupNumber                 int
	subGroupNumber              int
	groupMemberGUID             int
	regenType                   int
	argo                        *ai.MonsterAIAgro
	lastAutoRunTime             time.Time
	lastAutoDelay               time.Duration
	crywolfMVPScore             int
	lastCheckTick               time.Time
	autoRecuperationTime        time.Time
	skillDistanceErrorCount     int
	skillDistanceErrorTick      time.Time
	skillInfo                   skillInfo
	bufferIndex                 int
	buffID                      int
	buffPlayerIndex             int
	agiCheckTime                time.Time
	warehouseSaveLock           uint8
	crcCheckTime                time.Time
	off                         bool
	offLevel                    bool
	offLevelTime                int
}

func (o *Object) Reset() {

}

type bill struct {
	certify uint8
	payCode uint8
	endTime time.Time
	endDays [13]uint8
}

func (*bill) init() {

}

var (
	mu                        sync.Mutex
	Disconnect                bool
	maxObjectCount            int
	objects                   []Object
	objectUserCountStartIndex int
	objectCount               int
	objectUserCount           int
	objectMonsterCount        int
	objectSummonMonsterCount  int
	objectBills               []bill
)

func init() {
	maxObjectCount = conf.Server.MaxObjectMonsterCount + conf.Server.MaxObjectSummonMonsterCount + conf.Server.MaxObjectUserCount
	objects = make([]Object, maxObjectCount)
	objectBills = make([]bill, conf.Server.MaxObjectUserCount)
	// 先有怪后有玩家
	objectUserCountStartIndex = maxObjectCount - conf.Server.MaxObjectUserCount
	objectCount = objectUserCountStartIndex

}

func objectMaxRange(index int) bool {
	if index < 0 || index >= maxObjectCount {
		return false
	}
	return true
}

// ObjectAdd search a avaliable object from pool
// and return the object index
func ObjectAdd(conn *network.Conn) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	if Disconnect {
		return -1, fmt.Errorf("disconnect")
	}

	if objectUserCount > conf.Server.MaxObjectUserCount {
		// 响应
		res := &network.Response{}
		body := []byte{0x04}
		res.WriteHead2(0xC1, 0xF1, 0x01).Write(body)
		conn.Write(res)
		return -1, fmt.Errorf("current user number: [%d], over maximum number of users: [%d]", objectUserCount, conf.Server.MaxObjectUserCount)
	}

	index := objectCount
	cnt := conf.Server.MaxObjectUserCount
	for cnt > 0 {
		if objects[index].Connected == PlayerEmpty {
			break
		}
		index++
		if index >= maxObjectCount {
			index = objectUserCountStartIndex
		}
		cnt--
	}
	if cnt == 0 {
		return 0, fmt.Errorf("have no free index")
	}

	o := &objects[index]
	o.Reset()
	o.LoginMsgSend = false
	o.LoginMsgCount = 0
	o.index = index
	o.conn = conn
	o.ConnectCheckTime = time.Now()
	o.AutoSaveTime = o.ConnectCheckTime
	o.Connected = PlayerConnected
	o.CheckSpeedHack = false
	o.EnableCharacterCreate = false
	o.Type = ObjectUser
	objectBills[index-objectUserCountStartIndex].init()
}
