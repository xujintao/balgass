package conf

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

type serverConfig struct {
	MaxServerCount                   int    `ini:"MAX_SERVER"`
	UseJoinServer                    bool   `ini:"UseJoinServer"`
	UseDataServer                    bool   `ini:"UseDataServer"`
	UseExDataServer                  bool   `ini:"UseExDataServer"`
	JoinServerPort                   int    `ini:"JoinServerPort"`
	DataServerPort                   int    `ini:"DataServerPort"`
	ExDataServerPort                 int    `ini:"ExDataServerPort"`
	WanIP                            string `ini:"WanIP"`
	DataServerOnlyMode               bool   `ini:"DataServerOnlyMode"`
	MapServerInfoPath                string `ini:"MapServerInfoPath"`
	MachineIDConnectionLimitPerGroup int    `ini:"MachineIDConnectionLimitPerGroup"`
	MagicGladiatorCreateMinLevel     int    `ini:"MagicGladiatorCreateMinLevel"`
	DarkLordCreateMinLevel           int    `ini:"DarkLordCreateMinLevel"`
	GrowLancerCreateMinLevel         int    `ini:"GrowLancerCreateMinLevel"`
}

type sqlConfig struct {
	PasswordEncryptType int    `ini:"PasswordEncryptType"`
	MuOnlineDB          string `ini:"MuOnlineDB"`
	MeMuOnlineDB        string `ini:"MeMuOnlineDB"`
	EventDB             string `ini:"EventDB"`
	RankingDB           string `ini:"RankingDB"`
	User                string `ini:"User"`
	Pass                string `ini:"Pass"`
	SQLServerName       string `ini:"SQLServerName"`
}

type gensConfig struct {
	GensRankingUpdateTimeHour int    `ini:"GensRankingUpdateTimeHour"`
	GensRankingPath           string `ini:"GensRankingPath"`
	GensReJoinDaysLimit       int    `ini:"GensReJoinDaysLimit"`
}

type allowedIPListConfig struct {
	XMLName xml.Name `xml:"AllowedIPList"`
	IPs     []struct {
		Address string `xml:"Address,attr"`
	} `xml:"IP"`
}

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

func mapYAML(file string, v interface{}) {
	log.Printf("Load %s", file)
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read %s, %v", file, err)
	}
	if err := yaml.Unmarshal(buf, v); err != nil {
		log.Fatalf("Failed to unmarshal %s, %v", file, err)
	}
}

var (
	Server        serverConfig
	SQL           sqlConfig
	Gens          gensConfig
	AllowedIPList allowedIPListConfig
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mapINI("IGCDS.ini", "SETTINGS", &Server)
	mapINI("IGCDS.ini", "SQL", &SQL)
	mapINI("IGCDS.ini", "GensSystem", &Gens)
	mapXML("IGC_AllowedIPList.xml", &AllowedIPList)
}
