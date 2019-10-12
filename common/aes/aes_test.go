package aes

import (
	"bytes"
	"testing"
)

var buf1 = [5]byte{104, 101, 108, 108, 111}
var buf2 = [17]byte{
	122, 84, 253, 52, 133, 151, 96, 4, 121, 51, 80, 230, 132, 236, 126, 210,
	11,
}

var buf3 = [16]byte{101, 120, 97, 109, 112, 108, 101, 112, 108, 97, 105, 110, 116, 101, 120, 116}
var buf4 = [33]byte{
	54, 51, 193, 49, 72, 227, 229, 47, 249, 28, 98, 150, 219, 34, 197, 252,
	191, 115, 144, 50, 157, 103, 197, 75, 219, 192, 185, 180, 143, 103, 51, 113,
	16,
}

func TestEncrypt(t *testing.T) {
	bufenc, err := Encrypt(buf1[:])
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(bufenc, buf2[:]) != 0 {
		t.Error("Encrypt failed")
	}

	bufenc, err = Encrypt(buf3[:])
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(bufenc, buf4[:]) != 0 {
		t.Error("Encrypt failed")
	}
}

func TestDecrypt(t *testing.T) {
	bufdec, err := Decrypt(buf2[:])
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(bufdec, buf1[:]) != 0 {
		t.Error("Decrypt failed")
	}

	bufdec, err = Decrypt(buf4[:])
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(bufdec, buf3[:]) != 0 {
		t.Error("Decrypt failed")
	}
}
