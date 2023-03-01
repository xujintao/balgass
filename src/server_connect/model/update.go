package model

import (
	"bytes"
	"encoding/binary"
)

// Version semantic version
type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

// Unmarshal unmarshal byte slice to struct
func (v *Version) Unmarshal(buf []byte) error {
	return binary.Read(bytes.NewReader(buf), binary.LittleEndian, v)
}

// AutoUpdateConfig update info return to client when client version unmatched
// 1, config
// 2, response
type AutoUpdateConfig struct {
	Ver         Version `ini:"-"`
	VerStr      string  `ini:"Version"`
	HostURL     string  `ini:"HostURL"`
	FTPPort     int     `ini:"FTPPort"`
	FTPLogin    string  `ini:"FTPLogin"`
	FTPPasswd   string  `ini:"FTPPasswd"`
	VersionFile string  `ini:"VersionFile"`
}

// Marshal marshal struct to byte slice
func (c *AutoUpdateConfig) Marshal() ([]byte, error) {
	// ver
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, &c.Ver); err != nil {
		return nil, err
	}
	// host
	var host [100]byte
	copy(host[:], c.HostURL)
	if _, err := buf.Write(host[:]); err != nil {
		return nil, err
	}
	// port
	if err := binary.Write(buf, binary.LittleEndian, uint16(c.FTPPort)); err != nil {
		return nil, err
	}
	// login
	var login [20]byte
	copy(login[:], c.FTPLogin)
	if _, err := buf.Write(login[:]); err != nil {
		return nil, err
	}
	// passwd
	var passwd [20]byte
	copy(passwd[:], c.FTPPasswd)
	if _, err := buf.Write(passwd[:]); err != nil {
		return nil, err
	}
	// file
	var file [20]byte
	copy(file[:], c.VersionFile)
	if _, err := buf.Write(file[:]); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
