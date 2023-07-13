package maps

import (
	"testing"
)

func TestGetMapPos(t *testing.T) {
	x, y := MapManager.GetMapRegenPos(Lorencia)
	rect := MapManager[Lorencia].regenRect
	if !(x >= rect.left &&
		x <= rect.right &&
		y >= rect.top &&
		y <= rect.bottom) {
		t.Error()
	}
}
