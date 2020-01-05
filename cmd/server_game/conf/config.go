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
)

func init() {
	mapINISection("GameServer.ini", "GameServerInfo", &Server)
	mapXML("IGC_ConnectMember.xml", &ConnectMember)
	mapXML("IGC_VipSettings.xml", &VipSystem)
	mapINI("../../config/common/IGCData/IGC_Common.ini", &Common)
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
		log.Fatalf("Failed to read %s, %v", file, err)
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

type chaosBox struct {
	Level10      int    `xml:"Level10,attr"`
	Level11      int    `xml:"Level11,attr"`
	Level12      int    `xml:"Level12,attr"`
	Level13      int    `xml:"Level13,attr"`
	Level14      int    `xml:"Level14,attr"`
	Level15      int    `xml:"Level15,attr"`
	AddLuck      int    `xml:"AddLuck,attr"`
	SocketWeapon string `xml:"SocketWeapon,attr"`
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
			Type              int      `xml:"Type,attr"`
			Name              string   `xml:"Name,attr"`
			MLMonsterMinLevel int      `xml:"ML_MonsterMinLevel,attr"`
			PointPerReset     int      `xml:"PointPerReset,attr"`
			NightStartHour    int      `xml:"NightStartHour,attr"`
			NightStartMinute  int      `xml:"NightStartMinute,attr"`
			NightEndHour      int      `xml:"NightEndHour,attr"`
			NightEndMinute    int      `xml:"NightEndMinute,attr"`
			Day               vipBonus `xml:"Day"`
			Night             vipBonus `xml:"Night"`
			ChaosBoxMixRates  struct {
				Normal    chaosBox `xml:"Normal"`
				Enhanced  chaosBox `xml:"Enhanced"`
				Socket    chaosBox `xml:"Socket"`
				Pentagram chaosBox `xml:"Pentagram"`
				Wing      struct {
					Second          int `xml:"Second,attr"`
					Monster         int `xml:"Monster,attr"`
					Third           int `xml:"Third,attr"`
					Cape            int `xml:"Cape,attr"`
					FeatherOfCondor int `xml:"FeatherOfCondor,attr"`
				} `xml:"Wing"`
			} `xml:"ChaosBoxMixRates"`
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

type configChaosBox struct{}

type configItemPrice struct{}

type configPet struct{}

type configOffTrade struct{}

type configCalcCharacter struct{}

type configPK struct{}
