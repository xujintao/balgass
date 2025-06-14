package model

import (
	"bytes"
	"encoding/binary"
)

type MsgServerRegister struct {
	Code    int
	Percent int
	Type_   int
}

func (msg *MsgServerRegister) Marshal() ([]byte, error) {
	var bw bytes.Buffer

	// code
	bw.WriteByte(0x00)
	err := binary.Write(&bw, binary.LittleEndian, uint16(msg.Code))
	if err != nil {
		return nil, err
	}

	// percent
	err = bw.WriteByte(byte(msg.Percent))
	if err != nil {
		return nil, err
	}

	// type
	err = bw.WriteByte(byte(msg.Type_))
	if err != nil {
		return nil, err
	}

	return bw.Bytes(), nil
}
