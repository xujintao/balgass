package lang

import (
	"fmt"
	"log"
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
)

func init() {
	type langBaseConfig struct {
		DefaultLang string `xml:"DefaultLang,attr"`
		Lang        []struct {
			ID       int    `xml:"ID,attr"`
			FileName string `xml:"FileName,attr"`
			Enable   bool   `xml:"Enable,attr"`
			Codepage int    `xml:"Codepage,attr"`
		} `xml:"Lang"`
	}
	var langBase langBaseConfig
	conf.XML(path.Join(conf.PathCommon, "IGC_LangBase.xml"), &langBase)
	var valid bool
	for _, lang := range langBase.Lang {
		if lang.Enable {
			type languageConfig struct {
				Message struct {
					Msg []struct {
						ID   int    `xml:"ID,attr"`
						Text string `xml:"Text,attr"`
					} `xml:"Msg"`
				} `xml:"Message"`
				Map struct {
					Msg []struct {
						ID   int    `xml:"ID,attr"`
						Text string `xml:"Text,attr"`
					} `xml:"Msg"`
				} `xml:"Map"`
			}
			var language languageConfig
			conf.XML(path.Join(conf.PathCommon, fmt.Sprintf("Langs/%s", lang.FileName)), &language)
			// msg
			textManagerDefault.textMsgs = make(map[int]string)
			for _, msg := range language.Message.Msg {
				textManagerDefault.textMsgs[msg.ID] = msg.Text
			}

			// map
			textManagerDefault.textMaps = make(map[int]string)
			for _, msg := range language.Map.Msg {
				textManagerDefault.textMaps[msg.ID] = msg.Text
			}

			valid = true
			break
		}
	}
	if !valid {
		log.Fatalln("lang not specified")
	}
}

var textManagerDefault textManager

type textManager struct {
	textMsgs map[int]string
	textMaps map[int]string
}

func (m *textManager) getMsg(index int) string {
	msg, ok := m.textMsgs[index]
	if !ok {
		return "unknown message"
	}
	return msg
}

func (m *textManager) getMap(index int) string {
	msg, ok := m.textMsgs[index]
	if !ok {
		return "unknown map"
	}
	return msg
}

type msgText int

func (m msgText) String() string {
	return textManagerDefault.getMsg(int(m))
}

func (m msgText) Error() string {
	return m.String()
}

type mapText int

func (m mapText) String() string {
	return textManagerDefault.getMap(int(m))
}

const (
	MsgGameServerClosed msgText = 1 + iota
	MsgGameServerCloseCountdown
	MsgRestrengthen      msgText = 270
	MsgStrengthenSet     msgText = 273
	MsgStrengthenFailed  msgText = 274
	MsgStrengthenSuccess msgText = 275
)
