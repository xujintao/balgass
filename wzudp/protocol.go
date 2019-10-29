package wzudp

import (
	"encoding/binary"
)

type Message struct {
	Flag byte
	Len  int
	Code uint8
	Body []byte
}

func parseFrame(frame []byte) (*Message, error) {
	// aes
	flag := frame[0]

	// Message
	offset := 1
	len := 0
	switch flag {
	case 0xC1:
		len = int(frame[offset])
		offset++
	case 0xC2:
		len = int(binary.BigEndian.Uint16(frame[offset:]))
		offset += 2
	}
	code := frame[offset]
	offset++
	body := frame[offset:]
	return &Message{flag, len, code, body}, nil
}

func newFrame(msg *Message) []byte {
	var frame []byte
	flag := msg.Flag
	frame = append(frame, flag)
	switch flag {
	case 0xC1:
		frame = append(frame, byte(msg.Len))
	case 0xC2:
		var lenC2 [2]byte
		binary.BigEndian.PutUint16(lenC2[:], uint16(msg.Len))
		frame = append(frame, lenC2[:]...)
	}
	frame = append(frame, msg.Code)
	frame = append(frame, msg.Body...)
	return nil
}
