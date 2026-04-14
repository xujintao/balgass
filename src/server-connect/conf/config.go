package conf

import (
	"encoding/xml"
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

func ENV(v any) {
	err := envconfig.Process("", v)
	if err != nil {
		log.Fatal(err)
	}
	// config log
	var writes []io.Writer
	for _, s := range ServerEnv.LogFile {
		switch s {
		case "-":
			writes = append(writes, os.Stdout)
		default:
			f, err := os.OpenFile(s, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Fatal(err)
			}
			writes = append(writes, f)
		}
	}
	var l slog.Level
	switch ServerEnv.LogLevel {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	}
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				io.MultiWriter(writes...),
				&slog.HandlerOptions{
					Level: l,
					// AddSource: true,
				},
			),
		),
	)
}

func INI(dir, file, section string, v interface{}) {
	file = path.Join(dir, file)
	slog.Info("Load INI", "file", file, "section", section)
	f, err := ini.Load(file)
	if err != nil {
		slog.Error("Failed to load INI file",
			"file", file, "section", section, "error", err)
		os.Exit(1)
	}
	if err := f.Section(section).MapTo(v); err != nil {
		slog.Error("Failed to map INI section",
			"file", file, "section", section, "error", err)
		os.Exit(1)
	}
}

func XML(dir, file string, v interface{}) {
	file = path.Join(dir, file)
	slog.Info("Load XML", "file", file)
	buf, err := os.ReadFile(file)
	if err != nil {
		slog.Error("Failed to read XML file",
			"file", file, "error", err)
		os.Exit(1)
	}
	if err := xml.Unmarshal(buf, v); err != nil {
		slog.Error("Failed to unmarshal XML file",
			"file", file, "error", err)
		os.Exit(1)
	}
}

func YAML(dir, file string, v interface{}) {
	file = path.Join(dir, file)
	slog.Info("Load YAML", "file", file)
	buf, err := os.ReadFile(file)
	if err != nil {
		slog.Error("Failed to read YAML file",
			"file", file, "error", err)
		os.Exit(1)
	}
	if err := yaml.Unmarshal(buf, v); err != nil {
		slog.Error("Failed to unmarshal YAML file",
			"file", file, "error", err)
		os.Exit(1)
	}
}

func init() {
	ENV(&ServerEnv)
	PathConfig = ServerEnv.PathConfig
	INI(PathConfig, "IGCCS.ini", "Config", &Net)
	YAML(PathConfig, "news.yml", &New)
}

var (
	PathConfig string

	// ServerEnv env config for server
	ServerEnv configServerEnv

	// Net net config
	Net NetConfig

	// New new config
	New NewConfig
)

type configServerEnv struct {
	Debug      bool     `envconfig:"DEBUG" default:"false"`
	LogLevel   string   `envconfig:"LOG_LEVEL" default:"info"`
	LogFile    []string `envconfig:"LOG_FILE" default:"-"`
	PathConfig string   `envconfig:"PATH_CONFIG" default:"."`
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
