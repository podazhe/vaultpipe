package vault

import (
	"testing"
)

func TestScopeManager_RegisterAndList(t *testing.T) {
	sm := NewScopeManager()
	if err := sm.Register("prod", "PROD_", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := sm.Register("dev", "DEV_", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := sm.List()
	if len(names) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(names))
	}
}

func TestScopeManager_DuplicateRegister(t *testing.T) {
	sm := NewScopeManager()
	_ = sm.Register("prod", "PROD_", nil)
	err := sm.Register("prod", "PROD_", nil)
	if err == nil {
		t.Fatal("expected error for duplicate scope name")
	}
}

func TestScopeManager_EmptyName(t *testing.T) {
	sm := NewScopeManager()
	err := sm.Register("", "X_", nil)
	if err == nil {
		t.Fatal("expected error for empty scope name")
	}
}

func TestScopeManager_Resolve_LongestPrefix(t *testing.T) {
	sm := NewScopeManager()
	_ = sm.Register("infra", "INFRA_", nil)
	_ = sm.Register("infra_db", "INFRA_DB_", nil)

	scope, ok := sm.Resolve("INFRA_DB_PASSWORD")
	if !ok {
		t.Fatal("expected a scope match")
	}
	if scope.Name != "infra_db" {
		t.Errorf("expected infra_db scope, got %q", scope.Name)
	}
}

func TestScopeManager_Resolve_NoMatch(t *testing.T) {
	sm := NewScopeManager()
	_ = sm.Register("prod", "PROD_", nil)

	_, ok := sm.Resolve("DEV_SECRET")
	if ok {
		t.Error("expected no match for DEV_SECRET")
	}
}

func TestScopeManager_Partition(t *testing.T) {
	sm := NewScopeManager()
	_ = sm.Register("prod", "PROD_", nil)
	_ = sm.Register("dev", "DEV_", nil)

	secrets := map[string]string{
		"PROD_DB_PASS": "secret1",
		"DEV_API_KEY":  "secret2",
		"GLOBAL_TOKEN": "secret3",
	}

	buckets := sm.Partition(secrets)

	if len(buckets["prod"]) != 1 {
		t.Errorf("expected 1 prod secret, got %d", len(buckets["prod"]))
	}
	if len(buckets["dev"]) != 1 {
		t.Errorf("expected 1 dev secret, got %d", len(buckets["dev"]))
	}
	if len(buckets["_default"]) != 1 {
		t.Errorf("expected 1 default secret, got %d", len(buckets["_default"]))
	}
	if buckets["_default"]["GLOBAL_TOKEN"] != "secret3" {
		t.Errorf("unexpected default value: %q", buckets["_default"]["GLOBAL_TOKEN"])
	}
}
