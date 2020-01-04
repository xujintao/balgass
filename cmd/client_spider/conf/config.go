package conf

import (
	"encoding/xml"
	"io/ioutil"
	"log"

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

type Task struct {
	Scheme string   `yaml:"scheme"`
	Host   string   `yaml:"host"`
	Note   string   `yaml:"note"`
	Hosts  []string `yaml:"hosts"`
	Paths  []string `yaml:"paths"`
	Proxy  string   `yaml:"proxy"`
}
type config struct {
	Tasks []*Task `yaml:"tasks"`
}

var (
	Tasks []*Task
)

func init() {
	c := config{}
	mapYAML("config.yaml", &c)
	Tasks = c.Tasks
}
