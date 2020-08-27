package model

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"log"
)

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

// ServerListConfig info from config
type ServerListConfig struct {
	XMLName xml.Name `xml:"ServerList"`
	Servers []struct {
		Code    int    `xml:"Code,attr"`
		IP      string `xml:"IP,attr"`
		Port    int    `xml:"Port,attr"`
		Visible int    `xml:"Visible,attr"`
		Name    string `xml:"Name,attr"`
	} `xml:"Server"`
}

// ServerListInfo list info return to client
type ServerListInfo struct {
	Code     uint16 // little endian
	Percent  uint8
	PlayType uint8
}

// ServerListRes response getServerList
type ServerListRes struct {
	Size int16
	Slis []*ServerListInfo
}

// Marshal marshal struct to byte slice
func (r *ServerListRes) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, r.Size); err != nil {
		return nil, err
	}
	for _, sli := range r.Slis {
		if err := binary.Write(buf, binary.LittleEndian, sli); err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// ServerConnectInfo contains ip and port
type ServerConnectInfo struct {
	IP   string // 16 bytes
	Port uint16 // little endian
}

// Marshal marshal struct to byte slice
func (info *ServerConnectInfo) Marshal() ([]byte, error) {
	var ip [16]byte
	copy(ip[:], info.IP)
	buf := bytes.NewBuffer(ip[:])
	if err := binary.Write(buf, binary.LittleEndian, info.Port); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ServerRegisterInfo register info with which game server register itself
// request
type ServerRegisterInfo struct {
	Sli ServerListInfo
	// Sci          ServerConnectInfo
	UserCount    uint16
	AccountCount uint16
	MaxUserCount uint16
}

// Unmarshal unmarshal byte slice to struct
func (info *ServerRegisterInfo) Unmarshal(buf []byte) error {
	return binary.Read(bytes.NewReader(buf), binary.LittleEndian, info)
}
