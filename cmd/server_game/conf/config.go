package conf

import (
	"encoding/xml"
	"io/ioutil"
	"log"

	"gopkg.in/ini.v1"
)

var (
	// Server server config
	Server configServer

	// ConnectMember connect memeber config
	ConnectMember configConnectMember

	// VipSystem vip system config
	VipSystem configVipSystem

	// Common represents common config
	Common configCommon

	// ChaosBox represents chaosBox mix rate
	ChaosBox configChaosBox

	// PetRing represents pet and ring
	PetRing configPetRing

	// OffTrade represents personal shop works when player is offline
	OffTrade configOffTrade

	// CalcChar represents calculate percent config
	CalcChar configCalcCharacter
)

func init() {
	mapINISection("GameServer.ini", "GameServerInfo", &Server)
	mapXML("IGC_ConnectMember.xml", &ConnectMember)
	mapXML("IGC_VipSettings.xml", &VipSystem)
	mapINI("../../config/common/IGCData/IGC_Common.ini", &Common)
	mapXML("../../config/common/IGCData/IGC_ChaosBox.xml", &ChaosBox)
	mapXML("../../config/common/IGCData/IGC_PetSettings.xml", &PetRing)
	mapXML("../../config/common/IGCData/IGC_OffTrade.xml", &OffTrade)
	mapXML("../../config/common/IGCData/IGC_CalcCharacter.xml", &CalcChar)
}

func mapINI(file, v interface{}) {
	log.Printf("Load %s", file)
	f, err := ini.Load(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := f.MapTo(v); err != nil {
		log.Fatalln(err)
	}
}

func mapINISection(file, section string, v interface{}) {
	log.Printf("Load %s[%s]", file, section)
	f, err := ini.Load(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := f.Section(section).MapTo(v); err != nil {
		log.Fatalln(err)
	}
}

func mapXML(file string, v interface{}) {
	log.Printf("Load %s", file)
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := xml.Unmarshal(buf, v); err != nil {
		log.Fatalf("Failed to unmarshal %s, %v", file, err)
	}
}

type configServer struct {
	Name                string `ini:"ServerName"`
	Code                int    `ini:"ServerCode"`
	NonPVP              bool   `ini:"NonPK"`
	EnableConnectMember bool   `ini:"ConnectMemberLoad"`
	Type                int    `ini:"ServerType"`
	Port                int    `ini:"GameServerPort"`
	ConnectServerIP     string `ini:"ConnectServerIP"`
	ConnectServerPort   int    `ini:"ConnectServerPort"`
	JoinServerIP        string `ini:"JoinServerIP"`
	JoinServerPort      int    `ini:"JoinServerPort"`
	DataServerIP        string `ini:"DataServerIP"`
	DataServerPort      int    `ini:"DataServerPort"`
	ExDBIP              string `ini:"ExDBIP"`
	ExDBPort            int    `ini:"ExDBPort"`
	MaxConnectCount     int    `ini:"MachineIDConnectionLimitCount"`
	// Log
	MaxObjectUserCount          int `ini:"PlayerCount"`
	MaxObjectMonsterCount       int `ini:"MonsterCount"`
	MaxObjectSummonMonsterCount int `ini:"SummonMonsterCount"`
	MaxObjectItemCount          int `ini:"MapItemCount"`
}

type configConnectMember struct {
	XMLName  xml.Name `xml:"ConnectMember"`
	Accounts []struct {
		Name string `xml:"Name,attr"`
	} `xml:"Account"`
}

type vipBonus struct {
	ExpBonus           float32 `xml:"ExpBonus,attr"`
	DropBonus          int     `xml:"DropBonus,attr"`
	ExcDropBonus       int     `xml:"ExcDropBonus,attr"`
	MasterExpBonus     float32 `xml:"MasterExpBonus,attr"`
	MasterDropBonus    int     `xml:"MasterDropBonus,attr"`
	MasterExcDropBonus int     `xml:"MasterExcDropBonus,attr"`
}

type rateChaosBoxMix struct {
	Level10                     int  `xml:"Level10,attr"`
	Level11                     int  `xml:"Level11,attr"`
	Level12                     int  `xml:"Level12,attr"`
	Level13                     int  `xml:"Level13,attr"`
	Level14                     int  `xml:"Level14,attr"`
	Level15                     int  `xml:"Level15,attr"`
	EnableLevel15Notice         bool `xml:"Level15Notice,attr"`
	AddLuck                     int  `xml:"AddLuck,attr"`
	SocketWeapon                int  `xml:"SocketWeapon,attr"`
	SocketWeaponMixRequireMoney int  `xml:"SocketWeaponMixRequireMoney,attr"`
	Second                      int  `xml:"Second,attr"`
	Monster                     int  `xml:"Monster,attr"`
	Third                       int  `xml:"Third,attr"`
	Cape                        int  `xml:"Cape,attr"`
	FeatherOfCondor             int  `xml:"FeatherOfCondor,attr"`
}

type rateCHaosBoxMixs struct {
	Normal    rateChaosBoxMix `xml:"Normal"`
	Enhanced  rateChaosBoxMix `xml:"Enhanced"`
	Socket    rateChaosBoxMix `xml:"Socket"`
	Pentagram rateChaosBoxMix `xml:"Pentagram"`
	Wing      rateChaosBoxMix `xml:"Wing"`
}

type configVipSystem struct {
	XMLName                xml.Name `xml:"VipSystem"`
	LevelType              int      `xml:"LevelType,attr"`
	SendRatesChangeMessage bool     `xml:"SendRatesChangeMessage,attr"`
	Message                struct {
		Day   string `xml:"Day,attr"`
		Night string `xml:"Night,attr"`
	} `xml:"Message"`
	VipTypes struct {
		Vip []struct {
			Type              int              `xml:"Type,attr"`
			Name              string           `xml:"Name,attr"`
			MLMonsterMinLevel int              `xml:"ML_MonsterMinLevel,attr"`
			PointPerReset     int              `xml:"PointPerReset,attr"`
			NightStartHour    int              `xml:"NightStartHour,attr"`
			NightStartMinute  int              `xml:"NightStartMinute,attr"`
			NightEndHour      int              `xml:"NightEndHour,attr"`
			NightEndMinute    int              `xml:"NightEndMinute,attr"`
			Day               vipBonus         `xml:"Day"`
			Night             vipBonus         `xml:"Night"`
			RateChaosBoxMixs  rateCHaosBoxMixs `xml:"ChaosBoxMixRates"`
		} `xml:"Vip"`
	} `xml:"VipTypes"`
}

type color []int

type postCMD struct {
	Enable   bool  `ini:"Enable"`
	Cost     int   `ini:"Cost"`
	MinLevel int   `ini:"MinLevel"`
	Color    color `ini:"Color"`
	CoolDown int   `ini:"CoolDown"`
}

type configCommon struct {
	General struct {
		MaxLevelNormal                int     `ini:"MaxNormalLevel"`
		MaxLevelMaster                int     `ini:"MaxMasterLevel"`
		MasterPointPerLevel           int     `ini:"MasterPointPerLevel"`
		MinMonsterLevelForMasterExp   int     `ini:"MonsterMinLevelForMLExp"`
		ZenDropMultiplier             float32 `ini:"ZenDropMultipler"`
		EnableGuardSpeak              bool    `ini:"GuardSpeak"`
		GuardSpeakChance              int     `ini:"GuardSpeakChance"`
		GuardSpeakMsg                 string  `ini:"GuardSpeakMsg"`
		WelcomeMessage                string  `ini:"WelcomeMessage"`
		EnableBossMonsterKillNotice   bool    `ini:"BossMonsterKillNotice"`
		EnableEnterGameMessage        bool    `ini:"EnterGameMessageEnable"`
		Enable3thQuestMonsterCountMsg bool    `ini:"ThirdQuestMonsterCountMsg"`
		EnableTrade                   bool    `ini:"TradeBlock"`
		EnableTradeHarmonyItem        bool    `ini:"CanTradeHarmonyItem"`
		EnableTradeFullExcItem        bool    `ini:"CanTradeFullExcItem"`
		EnableUseSocketExcItem        bool    `ini:"CanUseSocketExcItem"`
		EnableUseAnciHarmonyItem      bool    `ini:"CanUseAnciHarmonyItem"`
		EnableTrade0SerialItem        bool    `ini:"CanTradeFFFFFFFFSerialItem"`
		EnableCheckValidItem          bool    `ini:"CheckValidItem"`
		EnableSellAllItem             bool    `ini:"EnableSellAllItems"`
		EnablePickLuckyItem           bool    `ini:"AllowToGetLuckyItemFromGround"`
		EnableSellFullExcItemInPShop  bool    `ini:"CanSellInStoreFullExcItem"`
		EnableSellFullExcItemToShop   bool    `ini:"CanSellToShopFullExcItem"`
		EnableEnhanceLuckyItemByJewel bool    `ini:"AllowEnchantLuckyItemByJewels"`
		LuckyItemDurabilityTime       int     `ini:"LuckyItemDurabilityTime"`
		EnableEnterEventWithPK        bool    `ini:"CanEnterEventWithPK"`
		RateSoul                      int     `ini:"UseSoulRate"`
		RateSoulLucky                 int     `ini:"UseSoulWithLuckRate"`
		RateLife                      int     `ini:"UseLifeRate"`
		EnableLifeOption28            bool    `ini:"Is28Option"`
		EnableSavePrivateChat         bool    `ini:"SavePrivateChat"`
		EnableAutoParty               bool    `ini:"AutoParty"`
		EnableReconnectSystem         bool    `ini:"ReconnectSystem"`
		EnableItem380DropMap          int     `ini:"DropMap380Items"`
	} `ini:"General"`

	PostCMD       postCMD `ini:"PostCMD"`
	PostCMDGlobal postCMD `ini:"GlobalPostCMD"`

	ChatColor struct {
		Info     color `ini:"Info"`
		Error    color `ini:"Error"`
		Chat     color `ini:"Chat"`
		Whisper  color `ini:"Whisper"`
		Party    color `ini:"Party"`
		Guild    color `ini:"Guild"`
		Alliance color `ini:"Alliance"`
		Gens     color `ini:"Gens"`
		GMChat   color `ini:"GMChat"`
	} `ini:"ChatColors"`

	AntiHack struct {
		EnableAgilityCheck            bool   `ini:"EnableAgilityCheck"`
		AgilityCheckTime              int    `ini:"AgilityDelayCheckTime"`
		EnableAntiRefCheck            bool   `ini:"EnableAntiRefCheckTime"`
		AntiRefCheckTime              int    `ini:"AntiRefCheckTimeMSEC"`
		EnableHitHackCheck            bool   `ini:"EnableHitHackDetection"`
		HitHackMaxAgility             int    `ini:"HitHackMaxAgility"`
		CRCMain                       uint   `ini:"MainExeCRC"`
		CRCDLL                        uint   `ini:"DLLCRC"`
		CRCPlayer                     uint   `ini:"PlayerBmdCRC"`
		CRCSkill                      uint   `ini:"SkillCRC"`
		CRCItem                       uint   `ini:"ItemCRC"`
		CRCInfo                       uint   `ini:"InfoCRC"`
		EnableKickUnmatchedDLLVersion bool   `ini:"DisconnectOnInvalidDLLVersion"`
		EnableKickAntiHackBreach      bool   `ini:"AntiHackBreachDisconnectUser"`
		EnableRecvHookProtection      bool   `ini:"RecvHookProtection"`
		PotionDelayTime               int    `ini:"PotionDelayTime"`
		PacketLimit                   int    `ini:"PacketLimit"`
		EnablePacketTimeCheck         bool   `ini:"EnablePacketTimeCheck"`
		PacketTimeMin                 int    `ini:"PacketTimeMinTimeMsec"`
		EnableHackDetectMessage       bool   `ini:"EnableHackDetectMessage"`
		HackDetectMessage             string `ini:"HackDetectMessage"`
		EnableAutoBanHackUser         bool   `ini:"EnableAutoBanAccountForHackUser"`
		EnableBlockAttackInSafeZone   bool   `ini:"EnableAttackBlockInSafeZone"`
	} `ini:"AntiHack"`

	ResetCMD struct {
		Enable                          bool `ini:"Enable"`
		MinLevel                        int  `ini:"MinLevel"`
		Cost                            int  `ini:"Cost"`
		MaxReset                        int  `ini:"MaxReset"`
		EnableResetStats                bool `ini:"IsResetStats"`
		EnableResetMasterLevel          bool `ini:"IsResetMasterLevel"`
		EnableMoveToCharSelectWindow    bool `ini:"MoveToCharSelectWindow"`
		EnableSaveOldStatPoint          bool `ini:"SaveOldStatPoint"`
		EnableRemoveEquipment           bool `ini:"RemoveEquipment"`
		PointPerReset                   int  `ini:"PointPerReset"`
		BlockNormalLevelPointAfterReset int  `ini:"BlockLevelUpPointAfterResetCount"`
		BlockMasterLevelPointAfterReset int  `ini:"BlockMLPointAfterResetCount"`
	} `ini:"ResetCMD"`

	MUHelper struct {
		Enable          bool `ini:"Enable"`
		MinLevel        int  `ini:"MinLevel"`
		Cost            int  `ini:"Cost"`
		NeedVIPLevel    int  `ini:"NeedVIPLevel"`
		AutoDisableTime int  `ini:"AutoDisableTime"`
	} `ini:"MuBot"`

	Guild struct {
		EnableCreate            bool `ini:"GuildCreate"`
		EnableDestroy           bool `ini:"GuildDestroy"`
		CreateLevel             int  `ini:"GuildCreateLevel"`
		MaxMember               int  `ini:"MaxGuildMember"`
		CastleOwnerDestroyLimit bool `ini:"CastleOwnerGuildDestroyLimit"`
		AllianceGuildMinMember  int  `ini:"AllianceMinGuildMember"`
		AllianceMaxGuildCount   int  `ini:"AllianceMaxGuilds"`
	} `ini:"Guilds"`

	GoldenMonster struct {
		GoldenDragonBoxDropCount      int `ini:"GoldenDragonBoxDropCount"`
		GreatGoldenDragonBoxDropCount int `ini:"GreatGoldenDragonBoxCount"`
	} `ini:"GoldenMonster"`

	Acheron struct {
		SpiritMapDropRate         int `ini:"SpiritMapDropRate"`
		SpiritMapDropMonsterLevel int `ini:"SpiritMapMonsterDropLevel"`
	} `ini:"Acheron"`

	EventInventory struct {
		Enable bool   `ini:"IsEventInventoryOpen"`
		Date   string `ini:"date"`
	} `ini:"EventInventory"`

	EggEvent struct {
		RateMoonRabbit      int `ini:"MoonRabbitSpawnRateFromBook"`
		RatePouchOfBlessing int `ini:"PouchBlessingSpawnRateFromBook"`
		RateFireFlameGhost  int `ini:"FireFlameSpawnRateFromBook"`
		RateGoldGoblin      int `ini:"GoldGoblinSpawnRateFromBook"`
	} `ini:"EggEvent"`

	CancelItemSale struct {
		Enable          bool    `ini:"IsCancelItemSale"`
		PriceMultiplier float32 `ini:"PriceMultipler"`
		ItemExpireTime  int     `ini:"ItemExpiryTime"`
	} `ini:"CancelItemSale"`

	SantaVillage struct {
		SantaClauseMinReset         int `ini:"SantaClauseMinReset"`
		SantaClause1stPrizeMaxVisit int `ini:"SantaClause1stPrizeMaxVisit"`
		SantaClause2ndPrizeMaxVisit int `ini:"SantaClause2ndPrizeMaxVisit"`
	} `ini:"SantaVillage"`
}

type configChaosBox struct {
	XMLName       xml.Name `xml:"ChaosBox"`
	CherryBlossom struct {
		CherryBlossomWhiteNeedItem int `xml:"CherryBlossomWhiteNeedItem,attr"`
		CherryBlossomRedNeedItem   int `xml:"CherryBlossomRedNeedItem,attr"`
		CherryBlossomGoldNeedItem  int `xml:"CherryBlossomGoldNeedItem,attr"`
	} `xml:"CherryBlossom"`
	RateCHaosBoxMixs rateCHaosBoxMixs `xml:"RateChaosBoxMixs"`
}

type buff struct {
	// Attack
	AddAttackValue        int `xml:"AddAttackValue,attr"`
	AddMagicAttackValue   int `xml:"AddMagicAttackValue,attr"`
	AddAttackPercent      int `xml:"AddAttackPercent,attr"`
	AddMagicAttackPercent int `xml:"AddMagicAttackPercent,attr"`
	AddAttackSpeed        int `xml:"AddAttackSpeed,attr"`

	// Defense
	AddDefenseValue   int `xml:"AddDefenseValue,attr"`
	AddDefensePercent int `xml:"AddDefensePercent,attr"`

	// Reduce
	ReduceDamageValue   int `xml:"ReduceDamageValue,attr"`
	ReduceDamagePercent int `xml:"ReduceDamagePercent,attr"`

	// HP
	AddHP int `xml:"AddHP,attr"`
}

type configPetRing struct {
	Pets struct {
		Angel       buff `xml:"Angel"`      // (13,0) 小天使/守护天使
		Imp         buff `xml:"Imp"`        // (13,1) 小恶魔
		FenrirGold  buff `xml:"FenrirGold"` // (13,37) 狼
		FenrirRed   buff `xml:"FenrirRed"`
		FenrirBlue  buff `xml:"FenrirBlue"`
		FenrirBlack buff `xml:"FenrirBlack"`
		Demon       buff `xml:"Demon"`       // (13,64) 大恶魔
		SpiritAngel buff `xml:"SpiritAngel"` // (13,65) 大天使
		Panda       buff `xml:"Panda"`       // (13,80) 熊猫
		Unicorn     buff `xml:"Unicorn"`     // (13,106) 兽角
		Skeleton    buff `xml:"Skeleton"`    // (13,123) 召唤骷髅
	} `xml:"Pets"`
	Rings struct {
		WizardRing            buff `xml:"WizardRing"`            // (13,20)
		SkeletonRing          buff `xml:"SkeletonRing"`          // (13,39)
		ChristmasRing         buff `xml:"ChristmasRing"`         // (13,41)
		PandaRing             buff `xml:"PandaRing"`             // (13,76)
		PandaBrownRing        buff `xml:"PandaBrownRing"`        // (13,77)
		PandaPinkRing         buff `xml:"PandaPinkRing"`         // (13,78)
		LethalWizardRing      buff `xml:"LethalWizardRing"`      // (13,107)
		RobotKnightRing       buff `xml:"RobotKnightRing"`       // (13,163)
		MiniRobotRing         buff `xml:"MiniRobotRing"`         // (13,164)
		MageRing              buff `xml:"MageRing"`              // (13,165)
		DecorationRing        buff `xml:"DecorationRing"`        // (13,169)
		DecorationBlessedRing buff `xml:"DecorationBlessedRing"` // (13,170)
		DarkTransformRing     buff `xml:"DarkTransformRing"`     // (13,268)
	} `xml:"Rings"`
}

type CoinType int

const (
	Zen CoinType = iota
	WCoinS8
	WcoinS6E3
	Goblin
)

type configOffTrade struct {
	Enable   bool     `xml:"Enable,attr"`
	CoinType CoinType `xml:"CoinType,attr"`
	Map      []struct {
		Number  int  `xml:"Number,attr"`
		Disable bool `xml:"Disable,attr"`
	} `xml:"Map"`
}

type charClass struct {
	DarkWizard     int `xml:"DarkWizard,attr"`
	DarkKnight     int `xml:"DarkKnight,attr"`
	FairyElf       int `xml:"FairyElf,attr"`
	MagicGladiator int `xml:"MagicGladiator,attr"`
	DarkLord       int `xml:"DarkLord,attr"`
	Summoner       int `xml:"Summoner,attr"`
	Rage           int `xml:"Rage,attr"`
	GrowLancer     int `xml:"GrowLancer,attr"`
}

type configCalcCharacter struct {
	MaxDamageDecreasePercent    int `xml:"MaxDamageDecreasePercent,attr"`
	MaxDamageReflectPercent     int `xml:"MaxDamageReflectPercent,attr"`
	DarkLordPetDamageMultiplier struct {
		DarkHorse float32 `xml:"DarkHorse,attr"`
		DarkRaven float32 `xml:"DarkRaven,attr"`
	} `xml:"DarkLordPetDamageMultiplier"`
	DarkSpiritDamageRate struct {
		PVE int       `xml:"PVE,attr"`
		PVP charClass `xml:"PVP"`
	} `xml:"DarkSpiritDamageRate"`
	DamageRate struct {
		PVP struct {
			DarkWizard     charClass `xml:"DarkWizard"`
			DarkKnight     charClass `xml:"DarkKnight"`
			FairyElf       charClass `xml:"FairyElf"`
			MagicGladiator charClass `xml:"MagicGladiator"`
			DarkLord       charClass `xml:"DarkLord"`
			Summoner       charClass `xml:"Summoner"`
			Rage           charClass `xml:"Rage"`
			GrowLancer     charClass `xml:"GrowLancer"`
		} `xml:"PVP"`
		PVE charClass `xml:"PVE"`
	} `xml:"DamageRate"`
	ElementalDamageRate struct {
		PVP struct {
			DarkWizard     charClass `xml:"DarkWizard"`
			DarkKnight     charClass `xml:"DarkKnight"`
			FairyElf       charClass `xml:"FairyElf"`
			MagicGladiator charClass `xml:"MagicGladiator"`
			DarkLord       charClass `xml:"DarkLord"`
			Summoner       charClass `xml:"Summoner"`
			Rage           charClass `xml:"Rage"`
			GrowLancer     charClass `xml:"GrowLancer"`
		} `xml:"PVP"`
		PVE charClass `xml:"PVE"`
	} `xml:"ElementalDamageRate"`
}

type configPK struct{}

type configItemPrice struct{}
