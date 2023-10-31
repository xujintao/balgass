package aes

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
		"68656c6c6f",
		"7a54fd3485976004793350e684ec7ed20b",
	},
	{
		"6578616d706c65706c61696e74657874",
		"3633c13148e3e52ff91c6296db22c5fc00",
	},
	{
		"f33c2407b727",
		"7c1021ef7bcf7a84b8b049c130621e660a",
	},
	{
		"f10202",
		"60d231d1f04d93caaf87d9e12323fe7c0d",
	},
}

func Test(t *testing.T) {
	for _, v := range texts {
		// encrypt
		plaintext, err := hex.DecodeString(v.plaintext)
		if err != nil {
			t.Errorf("convert hex string failed, %s", err.Error())
		}
		dstcipher, err := Encrypt(plaintext)
		if err != nil {
			t.Errorf("Encrypt failed, %s", err.Error())
		}
		if hex.EncodeToString(dstcipher) != v.ciphertext {
			t.Errorf("encrypt failed, expected(%s) got(%s)", v.ciphertext, hex.EncodeToString(dstcipher))
		}

		// decrypt
		ciphertext, err := hex.DecodeString(v.ciphertext)
		if err != nil {
			t.Errorf("convert hex string failed, %s", err.Error())
		}
		dstplain, err := Decrypt(ciphertext)
		if err != nil {
			t.Errorf("Decrypt failed, %s", err.Error())
		}
		if hex.EncodeToString(dstplain) != v.plaintext {
			t.Errorf("decrypt failed, expected(%s) got(%s)", v.plaintext, hex.EncodeToString(dstplain))
		}
	}
}
