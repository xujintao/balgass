package utils

import (
	"bytes"
	"testing"
)

func TestXor(t *testing.T) {
	data := []byte("https://github.com")
	data1 := Xor(data)
	data2 := Xor(data1)
	if !bytes.Equal(data, data2) {
		t.Error()
	}
}
