package xor

import (
	"encoding/hex"
	"testing"
)

type text struct {
	plaintext  string
	ciphertext string
}

var texts = [...]text{
	{
		"c104f40c",
		"c104f406",
	},
	{
		"c10400ff",
		"c1040001",
	},
	{
		"c110fa044b15d2a7d92493a24d752580",
		"c110fa00536572766572204e65777300",
	},
	{
		"c109fa118bff558bec",
		"c109fa15865acae2c4",
	},
}

func Test(t *testing.T) {
	for _, v := range texts {
		ciphertext, err := hex.DecodeString(v.plaintext)
		if err != nil {
			t.Errorf("convert hex string failed, %s", err.Error())
		}
		Enc(ciphertext, 3, len(ciphertext)-1)
		if hex.EncodeToString(ciphertext) != v.ciphertext {
			t.Errorf("encrypt failed, expected(%s) got(%s)", v.ciphertext, hex.EncodeToString(ciphertext))
		}

		plaintext, err := hex.DecodeString(v.ciphertext)
		if err != nil {
			t.Errorf("convert hex string failed, %s", err.Error())
		}
		Dec(plaintext, 3, len(plaintext)-1)
		if hex.EncodeToString(plaintext) != v.plaintext {
			t.Errorf("encrypt failed, expected(%s) got(%s)", v.plaintext, hex.EncodeToString(plaintext))
		}
	}
}
