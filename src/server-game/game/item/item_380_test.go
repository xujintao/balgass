package item

import "testing"

func TestConfigured380Effects(t *testing.T) {
	var got Item380Effect
	for _, code := range [][2]int{{0, 22}, {7, 29}, {8, 29}, {9, 29}, {10, 29}, {11, 29}} {
		section, index := code[0], code[1]
		if !Item380Manager.Is380Item(section, index) {
			t.Fatalf("configured 380 item %d/%d was not loaded", section, index)
		}
		Item380Manager.Apply380ItemEffect(section, index, &got)
	}
	want := Item380Effect{
		Item380EffectIncPVPAttackRate:  10,
		Item380EffectIncPVPAttack:      200,
		Item380EffectIncPVPDefenseRate: 50,
		Item380EffectIncPVPDefense:     200,
		Item380EffectIncMaxHP:          200,
		Item380EffectIncMaxSD:          700,
		Item380EffectAutoRecoverySD:    true,
		Item380EffectIncRecoverySDRate: 20,
	}
	if got != want {
		t.Fatalf("configured 380 effects = %+v, want %+v", got, want)
	}
}

func Test380MixSuccessRateUsesItemLevelAndCharm(t *testing.T) {
	mix := Item380Manager.MixConfig()
	for _, tc := range []struct{ level, rate int }{
		{6, 50}, {7, 60}, {9, 60}, {10, 70}, {13, 70}, {14, 80},
	} {
		if got := mix.SuccessRate(tc.level, 0); got != tc.rate {
			t.Errorf("SuccessRate(%d, 0) = %d, want %d", tc.level, got, tc.rate)
		}
		if got := mix.SuccessRate(tc.level, 10); got != tc.rate+10 {
			t.Errorf("SuccessRate(%d, 10) = %d, want %d", tc.level, got, tc.rate+10)
		}
	}
}
