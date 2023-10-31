package c1c2

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/xujintao/balgass/src/c1c2/aes"
	"github.com/xujintao/balgass/src/c1c2/xor"
)

func Hex2string(src []byte) string {
	dst := make([]byte, 2*len(src))
	hex.Encode(dst, src)
	buf := new(bytes.Buffer)
	for len(dst) > 0 {
		buf.Write(dst[:2])
		buf.WriteByte(' ')
		dst = dst[2:]
	}
	return buf.String()
}

// Request represents an request received by a server
type Request struct {
	// Flag uint8
	// Size int
	// Code    uint8
	Body    []uint8
	Encrypt bool
}

// Response represents the response from a request.
type Response struct {
	flag uint8
	body []uint8
}

// WriteHead write flag and code
func (r *Response) WriteHead(flag uint8) *Response {
	r.flag = flag
	return r
}

// Write write body
func (r *Response) Write(buf []byte) *Response {
	r.body = buf
	return r
}

func (r *Response) Body() []byte {
	return r.body
}

// readFrame read frame from bufio.Reader
func readFrame(br *bufio.Reader) ([]byte, error) {
	for {
		// peek 1 byte
		frameHead, err := br.Peek(1)
		if err != nil {
			return nil, err
		}
		flag := frameHead[0]
		switch flag {
		case 0xC1, 0xC2, 0xC3, 0xC4:
		default:
			br.ReadByte() // discard invalid frame guide flag
			continue      // try next
		}

		// peek 3 bytes
		frameHead, err = br.Peek(3)
		if err != nil {
			return nil, err
		}
		flag = frameHead[0]
		size := 0
		switch flag {
		case 0xC1, 0xC3:
			size = int(frameHead[1])
		case 0xC2, 0xC4:
			size = int(binary.BigEndian.Uint16(frameHead[1:]))
		}

		// peek size bytes
		if _, err := br.Peek(size); err != nil {
			return nil, err
		}
		// now we get a valid frame and read
		frame := make([]byte, size)
		if _, err := br.Read(frame); err != nil {
			return nil, err
		}
		return frame, nil
	}
}

// parseFrame parse frame to a Request
func parseFrame(frame []byte, needxor bool) (*Request, error) {
	flag := frame[0]
	// var plaintext []byte
	plaintext := new(bytes.Buffer)
	var enc bool

	// aes: begin with code
	switch flag {
	case 0xC3:
		dst, err := aes.Decrypt(frame[2:])
		if err != nil {
			return nil, fmt.Errorf("aes decrypt failed, %s, frame: %v", err.Error(), Hex2string(frame))
		}
		plaintext.WriteByte(0xC1)
		plaintext.WriteByte(byte(len(dst) + 2))
		plaintext.Write(dst)
		flag = 0xC1
		enc = true
		frame = plaintext.Bytes()
	case 0xC4:
		dst, err := aes.Decrypt(frame[3:])
		if err != nil {
			return nil, fmt.Errorf("aes decrypt failed, %s, frame: %s", err.Error(), Hex2string(frame))
		}
		plaintext.WriteByte(0xC2)
		binary.Write(plaintext, binary.BigEndian, uint16(len(dst)+3))
		plaintext.Write(dst)
		flag = 0xC2
		enc = true
		frame = plaintext.Bytes()
	}

	// xor: begin with subcode or body
	if needxor {
		switch flag {
		case 0xC1:
			xor.Dec(frame, 3, len(frame)-1)
		case 0xC2:
			xor.Dec(frame, 4, len(frame)-1)
		}
	}

	// Message
	offset := 1
	// size := 0
	switch flag {
	case 0xC1:
		// size = int(frame[offset])
		offset++
	case 0xC2:
		// size = int(binary.BigEndian.Uint16(frame[offset:]))
		offset += 2
	}
	// code := frame[offset]
	// offset++
	body := frame[offset:]
	req := &Request{
		// Flag: flag,
		// Size:    size,
		// Code:    code,
		Body:    body,
		Encrypt: enc,
	}

	return req, nil
}

// newFrame unmarshal Response to a frame
func newFrame(res *Response) ([]byte, error) {
	var head [3]byte
	var buf *bytes.Buffer
	size := 0

	// flag
	flag := res.flag
	head[0] = flag
	size++

	// size placehold
	switch flag {
	case 0xC1, 0xC3:
		buf = bytes.NewBuffer(head[:2])
		size++
	case 0xC2, 0xC4:
		buf = bytes.NewBuffer(head[:3])
		size += 2
	default:
		return nil, fmt.Errorf("invalid flag: %02x", flag)
	}

	// body with aes
	var body []byte
	var err error
	switch flag {
	case 0xC1, 0xC2:
		body = res.body
	case 0xC3, 0xC4:
		body, err = aes.Encrypt(res.body)
		if err != nil {
			return nil, err
		}
	}
	buf.Write(body)
	size += len(body)

	// size
	frame := buf.Bytes()
	switch flag {
	case 0xC1, 0xC3:
		frame[1] = byte(size)
	case 0xC2, 0xC4:
		binary.BigEndian.PutUint16(frame[1:3], uint16(size))
	}

	return frame, nil
}
