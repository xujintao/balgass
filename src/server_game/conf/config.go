package conf

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
	"path"

	"gopkg.in/ini.v1"
)

var (
	PathConfig string
	PathCommon string

	// SeasonX represents protocol compatibility with seasonX
	SeasonX bool

	// Server server config
	Server configServer

	// ConnectMember connect memeber config
	ConnectMember configConnectMember

	// VipSystem vip system config
	VipSystem configVipSystem

	// Common represents common config
	Common configCommon

	CommonServer configCommonServer

	// ChaosBox represents chaosBox mix rate
	ChaosBox configChaosBox

	// PetRing represents pet and ring
	PetRing configPetRing

	// OffTrade represents personal shop works when player is offline
	OffTrade configOffTrade

	// CalcChar represents calculate percent config
	CalcChar configCalcCharacter

	// PK represents PlayerKillSystem config
	PK configPK

	// Price represents item price
	Price configItemPrice

	// Events represents event config for every server
	Events configEvents

	// MapServer
	MapServers configMapServer
)

func init() {
	PathConfig = os.Getenv("CONFIG_PATH")
	PathCommon = os.Getenv("COMMON_PATH")
	log.Printf("[PWD]%s", os.Getenv("PWD"))
	if PathConfig == "" {
		PathConfig = "."
		log.Printf("$CONFIG_PATH is %q, use default %q", "", PathConfig)
	}
	if PathCommon == "" {
		PathCommon = "../../config/server_game_common"
		log.Printf("$COMMON_PATH is %q, use default %q", "", PathCommon)
	}
	INI(PathConfig, "GameServer.ini", &Server)
	XML(PathConfig, "IGC_ConnectMember.xml", &ConnectMember)
	XML(PathConfig, "IGC_VipSettings.xml", &VipSystem)
	INI(path.Join(PathCommon, "Data"), "CommonServer.cfg", &CommonServer)
	PathCommon = path.Join(PathCommon, "IGCData")
	INI(PathCommon, "IGC_Common.ini", &Common)
	XML(PathCommon, "IGC_ChaosBox.xml", &ChaosBox)
	XML(PathCommon, "IGC_PetSettings.xml", &PetRing)
	XML(PathCommon, "IGC_OffTrade.xml", &OffTrade)
	XML(PathCommon, "IGC_CalcCharacter.xml", &CalcChar)
	XML(PathCommon, "IGC_PlayerKillSystem.xml", &PK)
	INI(PathCommon, "IGC_PriceSettings.ini", &Price)
	XML(PathCommon, "events.xml", &Events)
	XML(PathCommon, "IGC_MapServerInfo.xml", &MapServers)
}

func INI(dir, file string, v interface{}) {
	file = path.Join(dir, file)
	log.Printf("Load %s", file)
	f, err := ini.Load(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := f.MapTo(v); err != nil {
		log.Fatalln(err)
	}
}

func XML(dir, file string, v interface{}) {
	file = path.Join(dir, file)
	log.Printf("Load %s", file)
	buf, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := xml.Unmarshal(buf, v); err != nil {
		log.Fatalf("Failed to unmarshal %s, %v", file, err)
	}
}

func JSON(dir, file string, v interface{}) {
	file = path.Join(dir, file)
	log.Printf("Load %s", file)
	// os.ReadFile(file)
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.NewDecoder(f).Decode(v)
	if err != nil {
		log.Fatalln(err)
	}
}

type configServer struct {
	GameServerInfo struct {
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
		MaxPlayerCount        int    `ini:"PlayerCount"`
		MaxMonsterCount       int    `ini:"MonsterCount"`
		MaxSummonMonsterCount int    `ini:"SummonMonsterCount"`
		MaxObjectItemCount    int    `ini:"MapItemCount"`
		HTTPPort              int    `ini:"HTTPPort"`
		DBName                string `int:"DBName"`
		DBUser                string `int:"DBUser"`
		DBPassword            string `int:"DBPassword"`
		DBHost                string `int:"DBHost"`
		DBPort                int    `int:"DBPort"`
	} `ini:"GameServerInfo"`
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
		SeasonX                       bool
		MaxLevelNormal                int     `ini:"MaxNormalLevel"`
		MaxLevelMaster                int     `ini:"MaxMasterLevel"`
		MasterPointPerLevel           int     `ini:"MasterPointPerLevel"`
		MinMonsterLevelForMasterExp   int     `ini:"MonsterMinLevelForMLExp"`
		ZenDropMultiplier             float64 `ini:"ZenDropMultipler"`
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
		EnableUseSetHarmonyItem       bool    `ini:"CanUseAnciHarmonyItem"`
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

type configCommonServer struct {
	GameServerInfo struct {
		ServerType                     int     `ini:"ServerType"`
		StalkProtocolEnable            bool    `ini:"StalkProtocol"`
		StalkProtocolID                string  `ini:"StalkProtocolId"`
		DisableTrashLog                bool    `ini:"DisableTrashLog"`
		CheckSumEnable                 bool    `ini:"CheckSumCheck"`
		ServerGroupGuildChatEnable     bool    `ini:"ServerGroupGuildChatting"`
		ServerGroupUnionChatEnable     bool    `ini:"ServerGroupUnionChatting"`
		CharacterCreateEnable          bool    `ini:"CreateCharacter"`
		TradeEnable                    bool    `ini:"Trade"`
		MultiWarehouseEnable           bool    `ini:"IsMultiWareHouse"`
		WarehouseCount                 int     `ini:"MultiWarehouseCount"`
		PostCommandMinLevel            int     `ini:"PostCommandMinLvl"`
		PostCommandMoneyRequire        int     `ini:"PostCommandMoneyReq"`
		QuestNPCTeleportTime           int     `ini:"QuestNPCTeleportTime"`
		ItemDropPercent                int     `ini:"ItemDropPer"`
		ExcelItemDropPercent           int     `ini:"ExcItemDropPer"`
		ItemLuckyDropPercent           int     `ini:"ItemLuckDrop"`
		ItemSkillDropPercent           int     `ini:"ItemSkillDrop"`
		ExcelItemLuckyDropPercent      int     `ini:"ExcItemLuckDrop"`
		ExcelItemSkillDropPercent      int     `ini:"ExcItemSkillDrop"`
		ZenDurationTime                int     `ini:"ZenDurationTime"`
		ItemDurationTime               int     `ini:"ItemDisapearTime"`
		LootingTime                    int     `ini:"LootingTime"`
		MaxStrength                    int     `ini:"MaxStrength"`
		MaxAgility                     int     `ini:"MaxAgility"`
		MaxVitality                    int     `ini:"MaxVitality"`
		MaxEnergy                      int     `ini:"MaxEnergy"`
		MaxCommand                     int     `ini:"MaxCommand"`
		ItemSerialCheckEnable          bool    `ini:"ItemSerialCheck"`
		ItemSerial0CheckEnable         bool    `ini:"ItemSerialZeroCheck"`
		HackUserKickEnable             bool    `ini:"DisconnectHackUser"`
		HackUserKickCount              int     `ini:"DetectedHackKickCount"`
		PacketHackCheckDisable         bool    `ini:"IsIgnorePacketHackDetect"`
		PenetrationSkillCheckEnable    bool    `ini:"EnableCheckPenetrationSkill"`
		SpeedHackCheckEnable           bool    `ini:"CheckSpeedHack"`
		SpeedHackPenaltyEnable         bool    `ini:"SpeedHackPenalty"`
		SkillDistanceCheckEnable       bool    `ini:"SkillDistanceCheck"`
		SkillDistanceCheckTemp         int     `ini:"SkillDistanceCheckTemp"`
		SkillDistanceKickEnable        bool    `ini:"SkillDistanceKick"`
		SKillDistanceKickCount         int     `ini:"SkillDistanceKickCount"`
		SkillDistanceKickCheckTime     int     `ini:"SkillDistanceKickCheckTime"`
		LevelPoint5                    int     `ini:"LevelUpPointNormal"`
		LevelPoint7                    int     `ini:"LevelUpPointMGDL"`
		AutoRecuperationEnable         bool    `ini:"UseCharacterAutoRecuperationSystem"`
		AutoRecuperationLevelLimit     int     `ini:"CharacterRecuperationMaxLevel"`
		PersonalShopEnable             bool    `ini:"PersonalShopOpen"`
		ItemRingTransformDropEnable    bool    `ini:"IsItemDropRingOfTransform"`
		ItemRingTransformDropRate      int     `ini:"ItemDropRingOfTransform"`
		BattleSoccerEnable             bool    `ini:"EnableBattleSoccer"`
		LuckyCoinDropRate              int     `ini:"LuckyCoinDrop"`
		Party2ExpBonus                 int     `ini:"NormalParty2ExpBonus"`
		Party3ExpBonus                 int     `ini:"NormalParty3ExpBonus"`
		Party4ExpBonus                 int     `ini:"NormalParty4ExpBonus"`
		Party5ExpBonus                 int     `ini:"NormalParty5ExpBonus"`
		Party2ExpBonusSet              int     `ini:"SetParty2ExpBonus"`
		Party3ExpBonusSet              int     `ini:"SetParty3ExpBonus"`
		Party4ExpBonusSet              int     `ini:"SetParty4ExpBonus"`
		Party5ExpBonusSet              int     `ini:"SetParty5ExpBonus"`
		ShadowPhantomLevelLimit        int     `ini:"ShadowPhantomMaxLevel"`
		ItemFenrirDropEnable           bool    `ini:"FenrirStuffItemDrop"`
		ItemFenrir01DropLevelMin       int     `ini:"FenrirStuff_01_DropLv_Min"`
		ItemFenrir01DropLevelMax       int     `ini:"FenrirStuff_01_DropLv_Max"`
		ItemFenrir01DropMap            int     `ini:"FenrirStuff_01_DropMap"`
		ItemFenrir01DropRate           int     `ini:"FenrirStuff_01_DropRate"`
		ItemFenrir02DropLevelMin       int     `ini:"FenrirStuff_02_DropLv_Min"`
		ItemFenrir02DropLevelMax       int     `ini:"FenrirStuff_02_DropLv_Max"`
		ItemFenrir02DropMap            int     `ini:"FenrirStuff_02_DropMap"`
		ItemFenrir02DropRate           int     `ini:"FenrirStuff_02_DropRate"`
		ItemFenrir03DropLevelMin       int     `ini:"FenrirStuff_03_DropLv_Min"`
		ItemFenrir03DropLevelMax       int     `ini:"FenrirStuff_03_DropLv_Max"`
		ItemFenrir03DropMap            int     `ini:"FenrirStuff_03_DropMap"`
		ItemFenrir03DropRate           int     `ini:"FenrirStuff_03_DropRate"`
		ItemFenrirRepairRate           int     `ini:"FenrirRepairRate"`
		ItemFenrirDefaultMaxDurSmall   int     `ini:"FenrirDefaultMaxDurSmall"`
		ItemFenrirElfMaxDurSmall       int     `ini:"FenrirElfMaxDurSmall"`
		ItemFenrir01LevelMinRate       int     `ini:"Fenrir_01Level_MixRate"`
		ItemFenrir02LevelMinRate       int     `ini:"Fenrir_02Level_MixRate"`
		ItemFenrir03LevelMinRate       int     `ini:"Fenrir_03Level_MixRate"`
		ItemDarkLordDropEnable         bool    `ini:"IsDropDarkLordItem"`
		ItemDarkSpiritAddExperience    float32 `ini:"DarkSpiritAddExperience"`
		ItemLochFeatherDropRate        int     `ini:"SleeveOfLordDropRate"` // 洛克之羽
		ItemLochFeatherDropLevel       int     `ini:"SleeveOfLordDropLevel"`
		ItemSpiritOfDarkHorseDropRate  int     `ini:"SoulOfDarkHorseDropRate"` // 黑马王
		ItemSpiritOfDarkHorseDropLevel int     `ini:"SoulOfDarkHorseDropLevel"`
		ItemSpiritOfDarkRavenDropRate  int     `ini:"SoulOfDarkSpiritDropRate"` // 天鹰
		ItemSpiritOfDarkRavenDropLevel int     `ini:"SoulOfDarkSpiritDropLevel"`

		// SD
		ShieldSystemEnable          bool `ini:"ShieldSystemOn"` // SD is short for Shield
		DamageDivideSDRate          int  `ini:"DamageDevideToSD"`
		DamageDivideHPRate          int  `ini:"DamageDevideToHP"`
		AttackSuccessRateOption     int  `ini:"SuccessAttackRateOption"`
		SDChargingOption            int  `ini:"SDChargingOption"`
		ConstNumberOfShieldPoint    int  `ini:"ConstNumberOfShieldPoint"`
		SDAutoRefillEnable          bool `ini:"ShieldAutoRefillOn"`
		SDAutoRefillSafeZoneEnable  bool `ini:"ShieldAutoRefillOnSafeZone"`
		CompoundPotionDropEnable    bool `ini:"CompoundPotionDropOn"` // 生命圣水
		CompoundPotion1DropRate     int  `ini:"CompoundPotionLv1DropRate"`
		CompoundPotion1DropLevel    int  `ini:"CompoundPotionLv1DropLevel"`
		CompoundPotion2DropRate     int  `ini:"CompoundPotionLv2DropRate"`
		CompoundPotion2DropLevel    int  `ini:"CompoundPotionLv2DropLevel"`
		CompoundPotion3DropRate     int  `ini:"CompoundPotionLv3DropRate"`
		CompoundPotion3DropLevel    int  `ini:"CompoundPotionLv3DropLevel"`
		ShieldComboMissOptionEnable bool `ini:"ShieldComboMissOptionOn"`       // ?
		SDPotion1MixRate            int  `ini:"ShiledPotionLv1MixSuccessRate"` // 防护药水
		SDPotion1MixMoney           int  `ini:"ShieldPotionLv1MixMoney"`
		SDPotion2MixRate            int  `ini:"ShiledPotionLv2MixSuccessRate"`
		SDPotion2MixMoney           int  `ini:"ShieldPotionLv2MixMoney"`
		SDPotion3MixRate            int  `ini:"ShiledPotionLv3MixSuccessRate"`
		SDPotion3MixMoney           int  `ini:"ShieldPotionLv3MixMoney"`
		SDGageConstA                int  `ini:"ShieldGageConstA"` // point*1.2
		SDGageConstB                int  `ini:"ShieldGageConstB"` // level*level/30

		// Kalima event
		KundunRefillHPSec   int `ini:"KundunRefillHPSec"`
		KundunRefillHP      int `ini:"KundunRefillHP"`
		KundunRefillHPTime  int `ini:"KundunRefillTime"`
		KundunHPLogSaveTime int `ini:"KundunHPLogSaveTime"`
		KundunMarkDropRate  int `ini:"KundunMarkDropRate"`

		// Red Dragon Invasion event
		Event1ItemDropTodayMax     int `ini:"Event1ItemDropTodayMax"`
		Event1ItemDropTodayPercent int `ini:"Event1ItemDropTodayPercent"`

		// Fire Cracker event
		FireCrackerEventEnable bool `ini:"FireCrackerEvent"`
		FireCrackerDropRate    int  `ini:"FireCrackerDropRate"`

		// Medals event
		MedalEventEnable    bool `ini:"MedalEvent"`
		MedalSilverDropRate int  `ini:"SilverMedalDropRate"`
		MedalGoldDropRate   int  `ini:"GoldMedalDropRate"`

		// Christmas event
		NPCChristmasEnable bool `ini:"MerryXMasTalkNpc"`
		NPCNewYearEnable   bool `ini:"HappyNewYearTalkNpc"`

		// Halloween event 万圣节
		HalloweenEventEnable bool `ini:"HallowinEventOn"`
		LuckyPumpkinDropRate int  `ini:"HallowinEventPumpkinOfLuckDropRate"`

		// Heart Of love event
		HeartOfLoveEventEnable bool `ini:"HeartOfLoveEvent"`
		HeartOfLoveDropRate    int  `ini:"HeartOfLoveDropRate"`

		// DarkLord heart event
		CondorFlameDropRate int `ini:"CondorFlameDropRate"` // 神鹰火种

		// CherryBlossom event 樱花活动
		CherryBlossomEventEnable bool `ini:"CherryBlossomEventOn"`
		CherryBlossomBoxDropRate int  `ini:"CherryBlossomEventItemDropRate"`

		// CandyBox event
		CandyBoxEventEnable   bool `ini:"CandyBoxEvent"`
		CandyPinkDropLevelMin int  `ini:"LightPurpleCandyBoxDropLv_Min"`
		CandyPinkDropLevelMax int  `ini:"LightPurpleCandyBoxDropLv_Max"`
		CandyPinkDropRate     int  `ini:"LightPurpleCandyBoxDropRate"`
		CandyRedDropLevelMin  int  `ini:"VermilionCandyBoxDropLv_Min"`
		CandyRedDropLevelMax  int  `ini:"VermilionCandyBoxDropLv_Max"`
		CandyRedDropRate      int  `ini:"VermilionCandyBoxDropRate"`
		CandyBlueDropLevelMin int  `ini:"DeepBlueCandyBoxDropLv_Min"`
		CandyBlueDropLevelMax int  `ini:"DeepBlueCandyBoxDropLv_Max"`
		CandyBlueDropRate     int  `ini:"DeepBlueCandyBoxDropRate"`

		// Box Drop Rate
		BoxSilverDropRate         int `ini:"SilverBoxDropRate"`
		BoxGoldDropRate           int `ini:"GoldBoxDropRate"`
		SecretGemDropRate1        int `ini:"MysteriouseBeadDropRate1"`
		SecretGemDropRate2        int `ini:"MysteriouseBeadDropRate2"`
		HiddenTreasureBoxDropRate int `ini:"HiddenTreasureBoxOfflineRate"`

		// Christmas ribbon // 圣诞箱
		RibbonBoxEventEnable       bool `ini:"RibbonBoxEvent"`
		RibbonBoxRedDropLevelMin   int  `ini:"RedRibbonBoxDropLv_Min"`
		RibbonBoxRedDropLevelMax   int  `ini:"RedRibbonBoxDropLv_Max"`
		RibbonBoxRedDropRate       int  `ini:"RedRibbonBoxDropRate"`
		RibbonBoxGreenDropLevelMin int  `ini:"GreenRibbonBoxDropLv_Min"`
		RibbonBoxGreenDropLevelMax int  `ini:"GreenRibbonBoxDropLv_Max"`
		RibbonBoxGreenDropRate     int  `ini:"GreenRibbonBoxDropRate"`
		RibbonBoxBlueDropLevelMin  int  `ini:"BlueRibbonBoxDropLv_Min"`
		RibbonBoxBlueDropLevelMax  int  `ini:"BlueRibbonBoxDropLv_Max"`
		RibbonBoxBlueDropRate      int  `ini:"BlueRibbonBoxDropRate"`

		// Chip event
		ChipEventEnable  bool `ini:"EventChipEvent"`
		BoxLuckyDropRate int  `ini:"BoxOfGoldDropRate"`

		// Loren deep
		GuardianJewelDropEnable bool `ini:"IsDropGemOfDefend"`
		GuardianJewelDropLevel  int  `ini:"GemOfDefendDropLevel"`
		GuardianJewelDropRate   int  `ini:"GemOfDefendDropRate"`

		// Rena/LordMark
		RenaDropRate int `ini:"MarkOfTheLord"`
	} `ini:"GameServerInfo"`
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

type PKLevel struct {
	Level                         int    `xml:"Level,attr"`
	ReqPoint                      int    `xml:"ReqPoint,attr"`
	ItemLoseRateOnMonsterKill     int    `xml:"ItemLoseRateOnMonsterKill,attr"`
	ItemLoseRateOnPlayerKill      int    `xml:"ItemLoseRateOnPlayerKill,attr"`
	ZenDeductionRateOnMonsterKill int    `xml:"ZenDeductionRateOnMonsterKill,attr"`
	ZenDeductionRateOnPlayerKill  int    `xml:"ZenDeductionRateOnPlayerKill,attr"`
	WarpChargeMultiplier          int    `xml:"WarpChargeMultiplier,attr"`
	CanSummonByDarkLord           bool   `xml:"CanSummonByDarkLord,attr"`
	Description                   string `xml:"Description,attr"`
	LevelRange                    []struct {
		Start         int `xml:"Start,attr"`
		End           int `xml:"End,attr"`
		DeductionRate int `xml:"DeductionRate,attr"`
	} `xml:"LevelRange"`
}

type configPK struct {
	DisablePKLevelIncrease bool `xml:"DisablePKLevelIncrease,attr"`
	DisablePenalty         bool `xml:"DisablePenalty,attr"`
	PKCanUseshops          bool `xml:"PKCanUseshops,attr"`
	DropExpensiveItems     bool `xml:"DropExpensiveItems,attr"`
	MaxItemLevelDrop       int  `xml:"MaxItemLevelDrop,attr"`
	PointDeductionDivisor  int  `xml:"PointDeductionDivisor,attr"`
	MurdererPointIncrease  int  `xml:"MurdererPointIncrease,attr"`
	PKClearCommand         struct {
		Enable                  bool `xml:"Enable,attr"`
		Cost                    int  `xml:"Cost,attr"`
		CostMultiplyByKillCount int  `xml:"CostMultiplyByKillCount,attr"`
	} `xml:"PKClearCommand"`
	General struct {
		PKLevels []PKLevel `xml:"PK"`
	} `xml:"General"`
	ExpDeduction struct {
		PKLevels []PKLevel `xml:"PK"`
	} `xml:"ExpDeduction"`
}

type configItemPrice struct {
	Value struct {
		ItemSellPriceDivisor int `ini:"ItemSellPriceDivisor"`
		JewelOfBlessPrice    int `ini:"JewelOfBlessPrice"`
		JewelOfSoulPrice     int `ini:"JewelOfSoulPrice"`
		JewelOfChaosPrice    int `ini:"JewelOfChaosPrice"`
		JewelOfLifePrice     int `ini:"JewelOfLifePrice"`
		JewelOfCreationPrice int `ini:"JewelOfCreationPrice"`
		CrestOfMonarchPrice  int `ini:"CrestOfMonarchPrice"`
		LochFeatherPrice     int `ini:"LochFeatherPrice"`
		JewelOfGuardianPrice int `ini:"JewelOfGuardianPrice"`
		WereRabbitEggPrice   int `ini:"WereRabbitEggPrice"`
	} `ini:"Value"`
}

type eventServer struct {
	Type   int    `xml:"type,attr"`
	Name   string `xml:"name,attr"`
	Enable bool   `xml:"enable,attr"`
}

type configEvents struct {
	BloodCastle struct {
		Servers []eventServer `xml:"server"`
	} `xml:"BloodCastle"`
	DevilSquare struct {
		Servers []eventServer `xml:"server"`
	} `xml:"DevilSquare"`
	DevilSquareSurival struct {
		Servers []eventServer `xml:"server"`
	} `xml:"DevilSquareSurival"`
	ChaosCastle struct {
		Servers []eventServer `xml:"server"`
	} `xml:"ChaosCastle"`
	ChaosCastleSurvival struct {
		Servers []eventServer `xml:"server"`
	} `xml:"ChaosCastleSurvival"`
	IllusionTemple struct {
		Servers []eventServer `xml:"server"`
	} `xml:"IllusionTemple"`
	CastleSiege struct {
		Servers []eventServer `xml:"server"`
	} `xml:"CastleSiege"`
	LorenDeep struct {
		Servers []eventServer `xml:"server"`
	} `xml:"LorenDeep"`
	Crywolf struct {
		Servers []eventServer `xml:"server"`
	} `xml:"Crywolf"`
	Kanturu struct {
		Servers []eventServer `xml:"server"`
	} `xml:"Kanturu"`
	Raklion struct {
		Servers []eventServer `xml:"server"`
	} `xml:"Raklion"`
	DoppelGanger struct {
		Servers []eventServer `xml:"server"`
	} `xml:"DoppelGanger"`
	ImperialGuardian struct {
		Servers []eventServer `xml:"server"`
	} `xml:"ImperialGuardian"`
	RingAttack struct {
		Servers []eventServer `xml:"server"`
	} `xml:"RingAttack"`
	ChristmasAttack struct {
		Servers []eventServer `xml:"server"`
	} `xml:"ChristmasAttack"`
	ArcaBattle struct {
		Servers []eventServer `xml:"server"`
	} `xml:"ArcaBattle"`
	AcheronGuardian struct {
		Servers []eventServer `xml:"server"`
	} `xml:"AcheronGuardian"`
	LastManStanding struct {
		Servers []eventServer `xml:"server"`
	} `xml:"LastManStanding"`
}

// MapServer was generated 2023-07-10 22:54:22 by https://xml-to-go.github.io/ in Ukraine.
type configMapServer struct {
	XMLName    xml.Name `xml:"MapServer"`
	Text       string   `xml:",chardata"`
	ServerInfo struct {
		Text    string `xml:",chardata"`
		Version string `xml:"Version,attr"`
		Serial  string `xml:"Serial,attr"`
	} `xml:"ServerInfo"`
	ServerList struct {
		Text   string `xml:",chardata"`
		Server []struct {
			Text       string `xml:",chardata"`
			Code       string `xml:"Code,attr"`
			Group      string `xml:"Group,attr"`
			Initiation string `xml:"Initiation,attr"`
			IP         string `xml:"IP,attr"`
			Port       string `xml:"Port,attr"`
			Name       string `xml:"Name,attr"`
		} `xml:"Server"`
	} `xml:"ServerList"`
	ServerMapping struct {
		Text   string `xml:",chardata"`
		Server []struct {
			Text           string `xml:",chardata"`
			Code           string `xml:"Code,attr"`
			MoveAble       string `xml:"MoveAble,attr"`
			MapNumber      string `xml:"MapNumber,attr"`
			DestServerCode string `xml:"DestServerCode,attr"`
			Name           string `xml:"Name,attr"`
		} `xml:"Server"`
	} `xml:"ServerMapping"`
}
