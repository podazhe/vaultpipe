package vault

import (
	"testing"
)

func TestRollbackManager_PushAndRollback(t *testing.T) {
	rm := NewRollbackManager(5)
	rm.Push(map[string]string{"KEY": "v1"})
	rm.Push(map[string]string{"KEY": "v2"})

	secrets, err := rm.Rollback(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["KEY"] != "v1" {
		t.Errorf("expected v1, got %s", secrets["KEY"])
	}
}

func TestRollbackManager_VersionNotFound(t *testing.T) {
	rm := NewRollbackManager(5)
	rm.Push(map[string]string{"KEY": "val"})
	_, err := rm.Rollback(99)
	if err == nil {
		t.Fatal("expected error for missing version")
	}
}

func TestRollbackManager_MaxSize(t *testing.T) {
	rm := NewRollbackManager(3)
	for i := 0; i < 5; i++ {
		rm.Push(map[string]string{"I": string(rune('0' + i))})
	}
	h := rm.History()
	if len(h) != 3 {
		t.Errorf("expected 3 entries, got %d", len(h))
	}
	if h[0].Version != 3 {
		t.Errorf("expected oldest version 3, got %d", h[0].Version)
	}
}

func TestRollbackManager_IsolatesCopy(t *testing.T) {
	rm := NewRollbackManager(5)
	secrets := map[string]string{"A": "original"}
	rm.Push(secrets)
	secrets["A"] = "mutated"

	got, err := rm.Rollback(1)
	if err != nil {
		t.Fatal(err)
	}
	if got["A"] != "original" {
		t.Errorf("expected original, got %s", got["A"])
	}
}

func TestRollbackManager_History(t *testing.T) {
	rm := NewRollbackManager(10)
	rm.Push(map[string]string{"X": "1"})
	rm.Push(map[string]string{"X": "2"})
	h := rm.History()
	if len(h) != 2 {
		t.Errorf("expected 2, got %d", len(h))
	}
}
