package vault

import (
	"testing"
)

// buildPinnedSet is a test helper that pins multiple paths and returns the manager.
func buildPinnedSet(t *testing.T, entries map[string]map[string]string) *PinManager {
	t.Helper()
	pm := NewPinManager()
	version := 1
	for path, data := range entries {
		if err := pm.Pin(path, version, data); err != nil {
			t.Fatalf("buildPinnedSet: pin %q: %v", path, err)
		}
		version++
	}
	return pm
}

func TestBuildPinnedSet_MultipleEntries(t *testing.T) {
	entries := map[string]map[string]string{
		"secret/alpha": {"A": "1"},
		"secret/beta":  {"B": "2"},
		"secret/gamma": {"C": "3"},
	}
	pm := buildPinnedSet(t, entries)
	list := pm.List()
	if len(list) != 3 {
		t.Errorf("expected 3 pins, got %d", len(list))
	}
}

func TestBuildPinnedSet_DataIntact(t *testing.T) {
	entries := map[string]map[string]string{
		"secret/svc": {"TOKEN": "xyz"},
	}
	pm := buildPinnedSet(t, entries)
	p, err := pm.Get("secret/svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Data["TOKEN"] != "xyz" {
		t.Errorf("expected TOKEN=xyz, got %s", p.Data["TOKEN"])
	}
}
