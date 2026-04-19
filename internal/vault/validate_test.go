package vault

import (
	"strings"
	"testing"
)

func TestValidateSecrets_ValidKeys(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	res := ValidateSecrets(secrets, ValidateOptions{})
	if !res.OK() {
		t.Fatalf("expected OK, got errors: %v", res.Errors)
	}
}

func TestValidateSecrets_InvalidKey(t *testing.T) {
	secrets := map[string]string{"123BAD": "value"}
	res := ValidateSecrets(secrets, ValidateOptions{})
	if res.OK() {
		t.Fatal("expected error for invalid key")
	}
	if !strings.Contains(res.Errors[0], "invalid key") {
		t.Errorf("unexpected error: %s", res.Errors[0])
	}
}

func TestValidateSecrets_NoEmpty(t *testing.T) {
	secrets := map[string]string{"KEY": ""}
	res := ValidateSecrets(secrets, ValidateOptions{NoEmpty: true})
	if res.OK() {
		t.Fatal("expected error for empty value")
	}
}

func TestValidateSecrets_WarnLong(t *testing.T) {
	secrets := map[string]string{"TOKEN": strings.Repeat("x", 300)}
	res := ValidateSecrets(secrets, ValidateOptions{WarnLong: 100})
	if !res.OK() {
		t.Fatalf("expected no errors, got: %v", res.Errors)
	}
	if len(res.Warnings) == 0 {
		t.Fatal("expected a warning for long value")
	}
}

func TestValidateSecrets_RequiredKeys(t *testing.T) {
	secrets := map[string]string{"EXISTING": "val"}
	res := ValidateSecrets(secrets, ValidateOptions{RequiredKeys: []string{"EXISTING", "MISSING"}})
	if res.OK() {
		t.Fatal("expected error for missing required key")
	}
	if !strings.Contains(res.Errors[0], "MISSING") {
		t.Errorf("unexpected error: %s", res.Errors[0])
	}
}

func TestValidateSecrets_Summary(t *testing.T) {
	secrets := map[string]string{"BAD KEY": ""}
	res := ValidateSecrets(secrets, ValidateOptions{NoEmpty: true})
	summary := res.Summary()
	if !strings.Contains(summary, "ERROR:") {
		t.Errorf("summary missing ERROR prefix: %s", summary)
	}
}
