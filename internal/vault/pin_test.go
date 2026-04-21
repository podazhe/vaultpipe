package vault

import (
	"testing"
)

func TestPinManager_PinAndGet(t *testing.T) {
	pm := NewPinManager()
	data := map[string]string{"API_KEY": "abc123"}

	if err := pm.Pin("secret/app", 3, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, err := pm.Get("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Version != 3 {
		t.Errorf("expected version 3, got %d", p.Version)
	}
	if p.Data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %s", p.Data["API_KEY"])
	}
}

func TestPinManager_GetMissing(t *testing.T) {
	pm := NewPinManager()
	_, err := pm.Get("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing pin")
	}
}

func TestPinManager_InvalidVersion(t *testing.T) {
	pm := NewPinManager()
	err := pm.Pin("secret/app", 0, map[string]string{})
	if err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestPinManager_EmptyPath(t *testing.T) {
	pm := NewPinManager()
	err := pm.Pin("", 1, map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestPinManager_Unpin(t *testing.T) {
	pm := NewPinManager()
	_ = pm.Pin("secret/app", 1, map[string]string{"X": "y"})
	pm.Unpin("secret/app")
	_, err := pm.Get("secret/app")
	if err == nil {
		t.Fatal("expected error after unpin")
	}
}

func TestPinManager_List(t *testing.T) {
	pm := NewPinManager()
	_ = pm.Pin("secret/a", 1, map[string]string{})
	_ = pm.Pin("secret/b", 2, map[string]string{})
	list := pm.List()
	if len(list) != 2 {
		t.Errorf("expected 2 pinned paths, got %d", len(list))
	}
}

func TestPinManager_IsolatesCopy(t *testing.T) {
	pm := NewPinManager()
	data := map[string]string{"KEY": "val"}
	_ = pm.Pin("secret/app", 1, data)
	data["KEY"] = "mutated"
	p, _ := pm.Get("secret/app")
	if p.Data["KEY"] != "val" {
		t.Error("pin should store an isolated copy of data")
	}
}
