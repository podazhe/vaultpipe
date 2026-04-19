package vault

import (
	"testing"
)

func TestRollbackManager_Integration_PushRollbackWrite(t *testing.T) {
	rm := NewRollbackManager(10)

	v1 := map[string]string{"DB_PASS": "secret1", "API_KEY": "key1"}
	v2 := map[string]string{"DB_PASS": "secret2", "API_KEY": "key2"}
	v3 := map[string]string{"DB_PASS": "secret3", "API_KEY": "key3"}

	rm.Push(v1)
	rm.Push(v2)
	rm.Push(v3)

	got, err := rm.Rollback(2)
	if err != nil {
		t.Fatalf("rollback failed: %v", err)
	}
	if got["DB_PASS"] != "secret2" {
		t.Errorf("DB_PASS: expected secret2, got %s", got["DB_PASS"])
	}
	if got["API_KEY"] != "key2" {
		t.Errorf("API_KEY: expected key2, got %s", got["API_KEY"])
	}
}

func TestRollbackManager_Integration_HistoryOrdering(t *testing.T) {
	rm := NewRollbackManager(10)
	for i := 1; i <= 4; i++ {
		rm.Push(map[string]string{"V": strconv.Itoa(i)})
	}
	h := rm.History()
	for i, entry := range h {
		if entry.Version != i+1 {
			t.Errorf("position %d: expected version %d, got %d", i, i+1, entry.Version)
		}
	}
}
