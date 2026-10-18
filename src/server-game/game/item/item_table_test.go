package item

import (
	"fmt"
	"strings"
	"testing"
)

func TestItemTableLookup(t *testing.T) {
	for _, id := range [][2]int{{-1, 0}, {16, 0}, {0, -1}, {0, 512}, {0, 511}} {
		if base, err := ItemTable.GetItemBase(id[0], id[1]); err == nil || base != nil {
			t.Errorf("lookup %v: base=%v err=%v", id, base, err)
		}
	}
}

func TestNewItemRejectsUndefinedItem(t *testing.T) {
	if it := NewItem(14, 15); it.ItemBase == nil || it.Section != 14 || it.Index != 15 {
		t.Fatalf("valid item = %+v", it)
	}
	defer func() {
		p := recover()
		if p == nil || !strings.Contains(fmt.Sprint(p), "NewItem(-1, 0)") {
			t.Fatalf("panic = %v", p)
		}
	}()
	NewItem(-1, 0)
}
