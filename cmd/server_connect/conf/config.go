package conf

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/xujintao/balgass/cmd/server_connect/model"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
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
	// Net net config
	Net model.NetConfig
	// New new config
	New model.NewConfig
	// ServerList server config set
	ServerList model.ServerListConfig
	// AutoUpdate auto update config
	AutoUpdate model.AutoUpdateConfig
	mu         sync.RWMutex
)

func init() {
	mapINI("IGCCS.ini", "Config", &Net)
	mapINI("IGCCS.ini", "AutoUpdate", &AutoUpdate)
	if _, err := fmt.Sscanf(AutoUpdate.VerStr, "%d.%d.%d", &AutoUpdate.Ver.Major, &AutoUpdate.Ver.Minor, &AutoUpdate.Ver.Patch); err != nil {
		log.Fatalf("Failed sscanf ver string, %v", err)
	}
	mapXML("IGC_ServerList.xml", &ServerList)
	mapYAML("news.yml", &New)
}
