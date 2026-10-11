package utils

import "testing"

func TestValidAddress(t *testing.T) {
	for _, address := range []string{"127.0.0.1", "255.255.255.255", "game.r2f2.com", "a-b.example"} {
		if !ValidAddress(address) {
			t.Errorf("rejected %q", address)
		}
	}
	for _, address := range []string{"", "game.example.com", "https://game", "game:56900", "bad..host", "-bad.example", "game/host", "::1", "999.999.999.999"} {
		if ValidAddress(address) {
			t.Errorf("accepted %q", address)
		}
	}
}
