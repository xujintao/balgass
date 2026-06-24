package monster

import "testing"

func TestIsHiddenMonsterClass(t *testing.T) {
	for _, tt := range []struct {
		class int
		want  bool
	}{
		{class: 99, want: false},
		{class: 100, want: true},
		{class: 110, want: true},
		{class: 111, want: false},
		{class: 247, want: false},
		{class: 249, want: false},
		{class: 523, want: true},
		{class: 689, want: true},
	} {
		if got := isHiddenMonsterClass(tt.class); got != tt.want {
			t.Errorf("isHiddenMonsterClass(%d) = %t, want %t", tt.class, got, tt.want)
		}
	}
}
