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
	println(1)
}
