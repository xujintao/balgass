package conf

import (
	"encoding/xml"
	"io/ioutil"
	"log"

	"gopkg.in/ini.v1"
)

func mapINI(file, section string, v interface{}) {
	log.Printf("Load %s:%s", file, section)
	cfg, err := ini.Load(file)
	if err != nil {
		log.Fatalf("Failed to load %s", file)
	}
	cfg.Section(section).MapTo(v)
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

var (
	// Server server config
	Server serverConfig

	// ConnectMember connect memeber config
	ConnectMember connectMemberConfig

	// VipSystem vip system config
	VipSystem vipSystemConfig
)

func init() {
	mapINI("GameServer.ini", "GameServerInfo", &Server)
	mapXML("IGC_ConnectMember.xml", &ConnectMember)
	mapXML("IGC_VipSettings.xml", &VipSystem)
}

type serverConfig struct {
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

type connectMemberConfig struct {
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

type vipSystemConfig struct {
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
