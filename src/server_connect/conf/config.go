package conf

import (
	"encoding/xml"
	"io/ioutil"
	"log"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

func INI(file, section string, v interface{}) {
	log.Printf("Load %s:%s", file, section)
	cfg, err := ini.Load(file)
	if err != nil {
		log.Fatalf("Failed to load %s", file)
	}
	cfg.Section(section).MapTo(v)
}

func XML(file string, v interface{}) {
	log.Printf("Load %s", file)
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read %s, %v", file, err)
	}
	if err := xml.Unmarshal(buf, v); err != nil {
		log.Fatalf("Failed to unmarshal %s, %v", file, err)
	}
}

func YAML(file string, v interface{}) {
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
	Net NetConfig

	// New new config
	New NewConfig
)

func init() {
	INI("IGCCS.ini", "Config", &Net)

	YAML("news.yml", &New)
}

// NetConfig info about listen and connect restriction
type NetConfig struct {
	TCPPort              int    `ini:"TCP_PORT"`
	UDPPort              int    `ini:"UDP_PORT"`
	MaxConnectionsPerIP  int    `ini:"MaxConnectionsPerIP"`
	MaxPacketsPerSecond  int    `ini:"MaxPacketsPerSecond"`
	LauncherProxyWhiteIP string `ini:"LauncherProxyWhiteListIP"`
}

// NewConfig represents some message sent to client
type NewConfig struct {
	Title string `yaml:"title"`
	Infos []struct {
		Index  int    `yaml:"index"`
		DateR  int    `yaml:"dateR"`
		DateG  int    `yaml:"dateG"`
		DateB  int    `yaml:"dateB"`
		TitleR int    `yaml:"titleR"`
		TitleG int    `yaml:"titleG"`
		TitleB int    `yaml:"titleB"`
		TextR  int    `yaml:"textR"`
		TextG  int    `yaml:"textG"`
		TextB  int    `yaml:"textB"`
		Day    int    `yaml:"day"`
		Month  int    `yaml:"month"`
		Year   int    `yaml:"year"`
		Title  string `yaml:"title"`
		Text   string `yaml:"text"`
	} `yaml:"news"`
}
