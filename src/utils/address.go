package utils

import (
	"net"
	"strings"
)

// ValidAddress reports whether address fits the client's 16-byte,
// zero-padded server address field. Short DNS names are passed to the client.
func ValidAddress(address string) bool {
	if len(address) == 0 || len(address) > 15 {
		return false
	}
	if ip := net.ParseIP(address).To4(); ip != nil && ip.String() == address {
		return true
	}
	if strings.Contains(address, ".") {
		numeric := true
		for i := 0; i < len(address); i++ {
			if address[i] != '.' && (address[i] < '0' || address[i] > '9') {
				numeric = false
				break
			}
		}
		if numeric {
			return false
		}
	}
	for _, label := range strings.Split(address, ".") {
		if len(label) == 0 || !asciiAlphaNumeric(label[0]) || !asciiAlphaNumeric(label[len(label)-1]) {
			return false
		}
		for i := 1; i < len(label)-1; i++ {
			if !asciiAlphaNumeric(label[i]) && label[i] != '-' {
				return false
			}
		}
	}
	return true
}

func asciiAlphaNumeric(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}
