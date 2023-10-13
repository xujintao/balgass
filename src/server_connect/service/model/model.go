package model

import (
	"bytes"
	"encoding/binary"

	"github.com/xujintao/balgass/src/server_connect/service/utils"
)

type MsgConnectReply struct {
	Result int
}

func (msg *MsgConnectReply) Marshal() ([]byte, error) {
	return []byte{1}, nil
}

type Version struct {
	Major int
	Minor int
	Patch int
}

type AutoUpdateConfig struct {
	Ver         Version `ini:"-"`
	VerStr      string  `ini:"Version"`
	HostURL     string  `ini:"HostURL"`
	FTPPort     int     `ini:"FTPPort"`
	FTPLogin    string  `ini:"FTPLogin"`
	FTPPasswd   string  `ini:"FTPPasswd"`
	VersionFile string  `ini:"VersionFile"`
}

type MsgCheckVersion struct {
	Version
}

func (msg *MsgCheckVersion) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// major
	major, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Major = int(major)

	// minor
	minor, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Minor = int(minor)

	patch, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Patch = int(patch)

	return nil
}

type MsgCheckVersionSuccess struct {
	Result bool
}

func (msg *MsgCheckVersionSuccess) Marshal() ([]byte, error) {
	return []byte{1}, nil
}

type MsgCheckVersionFailed struct {
	*AutoUpdateConfig
}

func (msg *MsgCheckVersionFailed) Marshal() ([]byte, error) {
	bw := new(bytes.Buffer)

	// Ver
	if err := bw.WriteByte(byte(msg.Ver.Major)); err != nil {
		return nil, err
	}
	if err := bw.WriteByte(byte(msg.Ver.Minor)); err != nil {
		return nil, err
	}
	if err := bw.WriteByte(byte(msg.Ver.Patch)); err != nil {
		return nil, err
	}

	// host
	var host [100]byte
	copy(host[:], msg.HostURL)
	if _, err := bw.Write(utils.Xor(host[:])); err != nil {
		return nil, err
	}

	// port
	if err := binary.Write(bw, binary.LittleEndian, uint16(msg.FTPPort)); err != nil {
		return nil, err
	}

	// login
	var login [20]byte
	copy(login[:], msg.FTPLogin)
	if _, err := bw.Write(utils.Xor(login[:])); err != nil {
		return nil, err
	}

	// passwd
	var passwd [20]byte
	copy(passwd[:], msg.FTPPasswd)
	if _, err := bw.Write(utils.Xor(passwd[:])); err != nil {
		return nil, err
	}

	// file
	var file [20]byte
	copy(file[:], msg.VersionFile)
	if _, err := bw.Write(utils.Xor(file[:])); err != nil {
		return nil, err
	}

	return bw.Bytes(), nil
}

type ServerState struct {
	Code    int
	Percent int
	Type_   int
}

type ServerConfig struct {
	Code    int    `xml:"Code,attr"`
	IP      string `xml:"IP,attr"`
	Port    int    `xml:"Port,attr"`
	Visible bool   `xml:"Visible,attr"`
	Name    string `xml:"Name,attr"`
}

type MsgRegister struct {
	ServerState
}

func (msg *MsgRegister) Unmarshal(buf []byte) error {
	br := bytes.NewReader(buf)

	// code
	br.ReadByte() // omit 1 byte to align
	var code uint16
	binary.Read(br, binary.LittleEndian, &code)
	msg.Code = int(code)

	// percent
	percent, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Percent = int(percent)

	// type
	type_, err := br.ReadByte()
	if err != nil {
		return err
	}
	msg.Type_ = int(type_)

	return nil
}

type MsgUnregister struct {
	Code int
}

type MsgGetServerList struct {
}

func (msg *MsgGetServerList) Unmarshal(buf []byte) error {
	return nil
}

type MsgGetServerListReply struct {
	ServerList []*ServerState
}

func (msg *MsgGetServerListReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer

	// array len
	binary.Write(&bw, binary.BigEndian, uint16(len(msg.ServerList)))

	// array elements
	for _, s := range msg.ServerList {
		if err := binary.Write(&bw, binary.LittleEndian, uint16(s.Code)); err != nil {
			return nil, err
		}
		if err := bw.WriteByte(byte(s.Percent)); err != nil {
			return nil, err
		}
		if err := bw.WriteByte(byte(s.Type_)); err != nil {
			return nil, err
		}
	}

	return bw.Bytes(), nil
}

type MsgGetServer struct {
	Code int
}

func (msg *MsgGetServer) Unmarshal(buf []byte) error {
	bw := bytes.NewReader(buf)

	var code uint16
	if err := binary.Read(bw, binary.LittleEndian, &code); err != nil {
		return err
	}
	msg.Code = int(code)
	return nil
}

type MsgGetServerReply struct {
	Server ServerConfig
}

func (msg *MsgGetServerReply) Marshal() ([]byte, error) {
	var bw bytes.Buffer

	//  ip
	var ip [16]byte
	copy(ip[:], msg.Server.IP)
	if _, err := bw.Write(ip[:]); err != nil {
		return nil, err
	}

	//  port
	if err := binary.Write(&bw, binary.LittleEndian, uint16(msg.Server.Port)); err != nil {
		return nil, err
	}

	return bw.Bytes(), nil
}
