package conf

import "testing"

func TestParseServerList(t *testing.T) {
	raw := `'[{"code":0,"ip":"127.0.0.1","port":56900,"visible":true,"name":"Regular"}]'`
	for _, input := range []string{raw, raw[1 : len(raw)-1]} {
		entries, err := parseServerList(input)
		if err != nil || len(entries) != 1 || entries[0].Code != 0 || !entries[0].Visible {
			t.Fatalf("parseServerList(%q) = %v, %v", input, entries, err)
		}
	}
}

func TestParseServerListRejectsInvalidEntries(t *testing.T) {
	for _, input := range []string{
		`not json`,
		`[]`,
		`[{"code":0,"ip":"127.0.0.1","port":0,"name":"Regular"}]`,
		`[{"code":0,"ip":"127.0.0.1","port":56900,"name":"Regular"},{"code":0,"ip":"127.0.0.2","port":56901,"name":"Other"}]`,
	} {
		if _, err := parseServerList(input); err == nil {
			t.Errorf("parseServerList(%q) accepted invalid list", input)
		}
	}
}
