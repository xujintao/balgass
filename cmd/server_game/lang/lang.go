package lang

import (
	"fmt"
	"log"
	"path"

	"github.com/xujintao/balgass/cmd/server_game/conf"
)

type textManager struct {
	msgTexts []string
	mapTexts []string
}

var defalutTextManager textManager

func init() {
	type LangBase struct {
		DefaultLang string `xml:"DefaultLang,attr"`
		Lang        []struct {
			ID       int    `xml:"ID,attr"`
			FileName string `xml:"FileName,attr"`
			Enable   bool   `xml:"Enable,attr"`
			Codepage int    `xml:"Codepage,attr"`
		} `xml:"Lang"`
	}
	var langBase LangBase
	conf.XML(path.Join(conf.PathCommon, "IGC_LangBase.xml"), &langBase)
	var valid bool
	for _, lang := range langBase.Lang {
		if lang.Enable {
			type Language struct {
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
			var language Language
			conf.XML(path.Join(conf.PathCommon, fmt.Sprintf("Langs/%s", lang.FileName)), &language)
			// msg
			defalutTextManager.msgTexts = make([]string, len(language.Message.Msg))
			for _, msg := range language.Message.Msg {
				defalutTextManager.msgTexts[msg.ID-1] = msg.Text
			}

			// map
			defalutTextManager.mapTexts = make([]string, len(language.Map.Msg))
			for _, msg := range language.Map.Msg {
				defalutTextManager.mapTexts[msg.ID] = msg.Text
			}

			valid = true
			break
		}
	}
	if !valid {
		log.Fatalln("lang not specified")
	}
}

type msgText int

func (m msgText) String() string {
	return defalutTextManager.msgTexts[m]
}

func (m msgText) Error() string {
	return m.String()
}

const (
	MsgTextGameServerClosed msgText = iota
	MsgTextGameServerCloseCountdown
	MsgTextRestrengthen      msgText = 270
	MsgTextStrengthenSet     msgText = 273
	MsgTextStrengthenFailed  msgText = 274
	MsgTextStrengthenSuccess msgText = 275
)

type mapText int

func (m mapText) String() string {
	return defalutTextManager.mapTexts[m]
}
