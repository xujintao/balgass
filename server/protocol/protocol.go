package protocol

import (
	"encoding/binary"

	"github.com/xujintao/balgass/server/xor"

	"github.com/xujintao/balgass/server/aes"
)

type Message struct {
	Flag byte
	Len  int
	Code uint8
	Body []byte
}

// client发往connect服务器的帧没有xor编码，
// client发往game服务器的帧进行了xor编码
func parseFrame(frame []byte, needxor bool) (*Message, error) {
	// aes
	flag := frame[0]
	var plaintext []byte
	switch flag {
	case 0xC3:
		dst, err := aes.Decrypt(frame[2:])
		if err != nil {
			return nil, err
		}
		plaintext = append(plaintext, 0xC1)
		plaintext = append(plaintext, byte(len(dst)+2))
		plaintext = append(plaintext, dst...)
	case 0xC4:
		dst, err := aes.Decrypt(frame[3:])
		if err != nil {
			return nil, err
		}
		plaintext = append(plaintext, 0xC2)
		var lenC2 [2]byte
		binary.BigEndian.PutUint16(lenC2[:], uint16(len(dst)+3))
		plaintext = append(plaintext, lenC2[:]...)
		plaintext = append(plaintext, dst...)
	}

	// xor
	if needxor {
		switch flag {
		case 0xC1:
			xor.Dec(plaintext, 3, len(plaintext)-1)
		case 0xC2:
			xor.Dec(plaintext, 4, len(plaintext)-1)
		}
	}

	// Message
	offset := 1
	len := 0
	switch flag {
	case 0xC1:
		len = int(plaintext[offset])
		offset++
	case 0xC2:
		len = int(binary.BigEndian.Uint16(plaintext[offset:]))
		offset += 2
	}
	code := plaintext[offset]
	offset++
	body := plaintext[offset:]
	return &Message{flag, len, code, body}, nil
}

func createFrame(msg *Message) []byte {
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
