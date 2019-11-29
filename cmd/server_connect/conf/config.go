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
	log.Println("Load IGCCS.ini")
	cfg, err := ini.Load("IGCCS.ini")
	if err != nil {
		log.Fatalf("Failed to load IGCCS.ini, %v", err)
	}
	cfg.Section("Config").MapTo(&Net)
	cfg.Section("AutoUpdate").MapTo(&AutoUpdate)
	_, err = fmt.Sscanf(AutoUpdate.VerStr, "%d.%d.%d", &AutoUpdate.Ver.Major, &AutoUpdate.Ver.Minor, &AutoUpdate.Ver.Patch)
	if err != nil {
		log.Fatalf("Failed sscanf ver string, %v", err)
	}

	log.Println("Load IGC_ServerList.xml")
	buf, err := ioutil.ReadFile("IGC_ServerList.xml")
	if err != nil {
		log.Fatalf("Failed to read IGC_ServerList.xml, %v", err)
	}
	if err := xml.Unmarshal(buf, &ServerList); err != nil {
		log.Fatalf("Failed to unmarshal IGC_ServerList.xml, %v", err)
	}

	log.Println("Load news.yml")
	buf, err = ioutil.ReadFile("news.yml")
	if err != nil {
		log.Fatalf("Failed to read news.yml, %v", err)
	}
	if err := yaml.Unmarshal(buf, &New); err != nil {
		log.Fatalf("Failed to unmarshal news.yml, %v", err)
	}
}
