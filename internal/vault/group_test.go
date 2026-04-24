package vault_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/vault"
)

func TestGroupSecrets_ByPrefix(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"APP_NAME":    "vaultpipe",
		"APP_VERSION": "1.0.0",
		"LOG_LEVEL":   "info",
	}

	groups := vault.GroupSecrets(secrets, vault.GroupByPrefix("_"))

	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB secrets, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 2 {
		t.Errorf("expected 2 APP secrets, got %d", len(groups["APP"]))
	}
	if len(groups["LOG"]) != 1 {
		t.Errorf("expected 1 LOG secret, got %d", len(groups["LOG"]))
	}
}

func TestGroupSecrets_NoMatchFallsToOther(t *testing.T) {
	secrets := map[string]string{
		"PLAIN": "value",
		"DB_HOST": "localhost",
	}

	groups := vault.GroupSecrets(secrets, vault.GroupByPrefix("_"))

	if _, ok := groups["other"]; !ok {
		t.Error("expected 'other' group for unmatched keys")
	}
	if groups["other"]["PLAIN"] != "value" {
		t.Errorf("expected PLAIN in 'other' group, got %v", groups["other"])
	}
}

func TestGroupSecrets_Empty(t *testing.T) {
	groups := vault.GroupSecrets(map[string]string{}, vault.GroupByPrefix("_"))
	if len(groups) != 0 {
		t.Errorf("expected empty groups map, got %d entries", len(groups))
	}
}

func TestGroupNames_Sorted(t *testing.T) {
	groups := map[string]map[string]string{
		"zebra": {"Z_KEY": "val"},
		"alpha": {"A_KEY": "val"},
		"middle": {"M_KEY": "val"},
	}

	names := vault.GroupNames(groups)

	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	if names[0] != "alpha" || names[1] != "middle" || names[2] != "zebra" {
		t.Errorf("expected sorted names, got %v", names)
	}
}

func TestGroupSecrets_ValuesPreserved(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "db.example.com",
		"DB_PASS": "s3cr3t",
	}

	groups := vault.GroupSecrets(secrets, vault.GroupByPrefix("_"))

	if groups["DB"]["DB_HOST"] != "db.example.com" {
		t.Errorf("unexpected value for DB_HOST: %s", groups["DB"]["DB_HOST"])
	}
	if groups["DB"]["DB_PASS"] != "s3cr3t" {
		t.Errorf("unexpected value for DB_PASS: %s", groups["DB"]["DB_PASS"])
	}
}

func TestGroupSecrets_SingleCharPrefix(t *testing.T) {
	secrets := map[string]string{
		"A_ONE": "1",
		"A_TWO": "2",
		"B_ONE": "3",
	}

	groups := vault.GroupSecrets(secrets, vault.GroupByPrefix("_"))

	if len(groups["A"]) != 2 {
		t.Errorf("expected 2 entries in group A, got %d", len(groups["A"]))
	}
	if len(groups["B"]) != 1 {
		t.Errorf("expected 1 entry in group B, got %d", len(groups["B"]))
	}
}
