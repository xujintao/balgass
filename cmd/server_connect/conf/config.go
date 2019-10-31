package conf

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type config struct {
	News struct {
		Title string `mapstructure:"title"`
		News  []struct {
			Index  int    `mapstructure:"index"`
			DateR  int    `mapstructure:"dateR"`
			DateG  int    `mapstructure:"dateG"`
			DateB  int    `mapstructure:"dateB"`
			TitleR int    `mapstructure:"titleR"`
			TitleG int    `mapstructure:"titleG"`
			TitleB int    `mapstructure:"titleB"`
			TextR  int    `mapstructure:"textR"`
			TextG  int    `mapstructure:"textG"`
			TextB  int    `mapstructure:"textB"`
			Day    int    `mapstructure:"Day"`
			Month  int    `mapstructure:"Month"`
			Year   int    `mapstructure:"Year"`
			Title  string `mapstructure:"title"`
			Text   string `mapstructure:"text"`
		} `mapstructure:"news"`
	} `mapstructure:"news"`
	Net struct {
		TCPPort              int    `mapstructure:"tcp_port"`
		UDPPort              int    `mapstructure:"udp_port"`
		MaxConnectionsPerIP  int    `mapstructure:"max_connections_per_ip"`
		MaxPacketsPerSecond  int    `mapstructure:"max_packets_per_second"`
		LauncherProxyWhiteIP string `mapstructure:"launcher_proxy_white_ip"`
	} `mapstructure:"net"`
	Update struct {
		Version     string `mapstructure:"version"`
		VersionFile string `mapstructure:"version_file"`
		HostURL     string `mapstructure:"host_url"`
		FtpLogin    string `mapstructure:"ftp_login"`
		FtpPasswd   string `mapstructure:"ftp_passwd"`
		FtpPort     int    `mapstructure:"ftp_port"`
	} `mapstructure:"update"`
	Servers []struct {
		Code    int    `mapstructure:"code"`
		IP      string `mapstructure:"ip"`
		Port    int    `mapstructure:"port"`
		Visible int    `mapstructure:"visible"`
		Name    string `mapstructure:"name"`
	} `mapstructure:"servers"`
}

var (
	Config config
	mu     sync.RWMutex
)

func init() {
	viper.SetConfigFile("config.yml")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatal(err)
	}
}
