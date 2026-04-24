package vault

import (
	"testing"
)

func TestNamespaceManager_FullNamespace_WithSub(t *testing.T) {
	nm, err := NewNamespaceManager("admin", "team-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := nm.FullNamespace(); got != "admin/team-a" {
		t.Errorf("expected admin/team-a, got %s", got)
	}
}

func TestNamespaceManager_FullNamespace_RootOnly(t *testing.T) {
	nm, err := NewNamespaceManager("admin", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := nm.FullNamespace(); got != "admin" {
		t.Errorf("expected admin, got %s", got)
	}
}

func TestNamespaceManager_EmptyRoot_Errors(t *testing.T) {
	_, err := NewNamespaceManager("", "team-a")
	if err == nil {
		t.Fatal("expected error for empty root, got nil")
	}
}

func TestNamespaceManager_QualifyPath(t *testing.T) {
	nm, _ := NewNamespaceManager("admin", "team-a")
	got := nm.QualifyPath("secret/db")
	if got != "admin/team-a/secret/db" {
		t.Errorf("expected admin/team-a/secret/db, got %s", got)
	}
}

func TestNamespaceManager_QualifyPath_LeadingSlash(t *testing.T) {
	nm, _ := NewNamespaceManager("admin", "team-a")
	got := nm.QualifyPath("/secret/db")
	if got != "admin/team-a/secret/db" {
		t.Errorf("expected admin/team-a/secret/db, got %s", got)
	}
}

func TestNamespaceManager_StripNamespace(t *testing.T) {
	nm, _ := NewNamespaceManager("admin", "team-a")
	got := nm.StripNamespace("admin/team-a/secret/db")
	if got != "secret/db" {
		t.Errorf("expected secret/db, got %s", got)
	}
}

func TestNamespaceManager_StripNamespace_NoPrefix(t *testing.T) {
	nm, _ := NewNamespaceManager("admin", "team-a")
	got := nm.StripNamespace("other/secret/db")
	if got != "other/secret/db" {
		t.Errorf("expected path unchanged, got %s", got)
	}
}

func TestNamespaceManager_QualifyAndStripSecrets_RoundTrip(t *testing.T) {
	nm, _ := NewNamespaceManager("root", "env")
	secrets := map[string]string{
		"db/password": "s3cr3t",
		"api/key":     "abc123",
	}
	qualified := nm.QualifySecrets(secrets)
	stripped := nm.StripSecrets(qualified)
	for k, v := range secrets {
		if stripped[k] != v {
			t.Errorf("round-trip failed for key %s: got %s, want %s", k, stripped[k], v)
		}
	}
}
