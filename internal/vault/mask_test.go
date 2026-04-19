package vault

import (
	"strings"
	"testing"
)

func TestMaskSecrets_DefaultOptions(t *testing.T) {
	secrets := map[string]string{
		"API_KEY": "supersecretvalue",
	}
	masked := MaskSecrets(secrets, nil, nil)
	v := masked["API_KEY"]
	if !strings.HasSuffix(v, "alue") {
		t.Errorf("expected suffix 'alue', got %q", v)
	}
	if !strings.HasPrefix(v, "********") {
		t.Errorf("expected mask prefix, got %q", v)
	}
}

func TestMaskSecrets_SensitiveKeyFullyMasked(t *testing.T) {
	secrets := map[string]string{
		"PASSWORD": "mypassword123",
	}
	masked := MaskSecrets(secrets, []string{"PASSWORD"}, nil)
	if masked["PASSWORD"] != "********" {
		t.Errorf("expected full mask, got %q", masked["PASSWORD"])
	}
}

func TestMaskSecrets_ShortValue(t *testing.T) {
	secrets := map[string]string{
		"TOKEN": "abc",
	}
	masked := MaskSecrets(secrets, nil, nil)
	if masked["TOKEN"] != "********" {
		t.Errorf("short value should be fully masked, got %q", masked["TOKEN"])
	}
}

func TestMaskSecrets_CustomMaskChar(t *testing.T) {
	secrets := map[string]string{
		"KEY": "abcdefghij",
	}
	opts := &MaskOptions{ShowChars: 2, MaskChar: "#"}
	masked := MaskSecrets(secrets, nil, opts)
	if !strings.HasPrefix(masked["KEY"], "########") {
		t.Errorf("expected '#' mask char, got %q", masked["KEY"])
	}
	if !strings.HasSuffix(masked["KEY"], "ij") {
		t.Errorf("expected suffix 'ij', got %q", masked["KEY"])
	}
}

func TestMaskSecrets_EmptyMap(t *testing.T) {
	masked := MaskSecrets(map[string]string{}, nil, nil)
	if len(masked) != 0 {
		t.Errorf("expected empty result, got %d entries", len(masked))
	}
}

func TestMaskSecrets_CaseInsensitiveSensitiveKey(t *testing.T) {
	secrets := map[string]string{
		"db_password": "topsecret99",
	}
	masked := MaskSecrets(secrets, []string{"DB_PASSWORD"}, nil)
	if masked["db_password"] != "********" {
		t.Errorf("expected full mask for case-insensitive key, got %q", masked["db_password"])
	}
}
