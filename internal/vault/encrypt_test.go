package vault

import (
	"strings"
	"testing"
)

var testKey = []byte("12345678901234567890123456789012") // 32 bytes

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
	}

	encrypted, err := EncryptSecrets(secrets, testKey, EncryptOptions{})
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	for k, v := range encrypted {
		if !strings.HasPrefix(v, "enc:") {
			t.Errorf("key %q: expected enc: prefix, got %q", k, v)
		}
	}

	decrypted, err := DecryptSecrets(encrypted, testKey)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	for k, want := range secrets {
		if got := decrypted[k]; got != want {
			t.Errorf("key %q: want %q, got %q", k, want, got)
		}
	}
}

func TestEncryptSecrets_SelectedKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "secret",
		"LOG_LEVEL":   "info",
	}

	encrypted, err := EncryptSecrets(secrets, testKey, EncryptOptions{Keys: []string{"DB_PASSWORD"}})
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	if !strings.HasPrefix(encrypted["DB_PASSWORD"], "enc:") {
		t.Errorf("DB_PASSWORD should be encrypted")
	}
	if encrypted["LOG_LEVEL"] != "info" {
		t.Errorf("LOG_LEVEL should be unchanged, got %q", encrypted["LOG_LEVEL"])
	}
}

func TestEncryptSecrets_Passthrough(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abc"}

	out, err := EncryptSecrets(secrets, testKey, EncryptOptions{Passthrough: true})
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if out["TOKEN"] != "abc" {
		t.Errorf("passthrough: expected unchanged value, got %q", out["TOKEN"])
	}
}

func TestDecryptSecrets_PlaintextPassthrough(t *testing.T) {
	secrets := map[string]string{"LOG_LEVEL": "debug"}

	out, err := DecryptSecrets(secrets, testKey)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if out["LOG_LEVEL"] != "debug" {
		t.Errorf("expected unchanged plaintext, got %q", out["LOG_LEVEL"])
	}
}

func TestEncryptSecrets_BadKeyLength(t *testing.T) {
	_, err := EncryptSecrets(map[string]string{"K": "v"}, []byte("short"), EncryptOptions{})
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestDecryptSecrets_BadKeyLength(t *testing.T) {
	_, err := DecryptSecrets(map[string]string{"K": "enc:abc"}, []byte("short"))
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestEncryptDecrypt_EmptyMap(t *testing.T) {
	out, err := EncryptSecrets(map[string]string{}, testKey, EncryptOptions{})
	if err != nil {
		t.Fatalf("encrypt empty: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
