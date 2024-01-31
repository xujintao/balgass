package utils

import "bytes"

func TrimStr(s []byte) string {
	s1, _, _ := bytes.Cut(s, []byte{0})
	return string(s1)
}
