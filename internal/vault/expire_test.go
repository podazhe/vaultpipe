package vault

import (
	"testing"
	"time"
)

func baseExpireSecrets() map[string]string {
	now := time.Now().UTC()
	return map[string]string{
		"DB_PASSWORD":            "s3cr3t",
		"DB_PASSWORD_EXPIRES_AT": now.Add(72 * time.Hour).Format(time.RFC3339),
		"API_KEY_EXPIRES_AT":     now.Add(10 * time.Minute).Format(time.RFC3339),
		"OLD_TOKEN_EXPIRES_AT":   now.Add(-1 * time.Hour).Format(time.RFC3339),
	}
}

func TestCheckExpiry_DetectsExpired(t *testing.T) {
	secrets := baseExpireSecrets()
	policy := ExpiryPolicy{WarnBefore: 30 * time.Minute, ErrorBefore: 5 * time.Minute}
	report, err := CheckExpiry(secrets, policy, time.Now().UTC())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Expired) != 1 || report.Expired[0] != "OLD_TOKEN_EXPIRES_AT" {
		t.Errorf("expected OLD_TOKEN_EXPIRES_AT expired, got %v", report.Expired)
	}
}

func TestCheckExpiry_DetectsWarning(t *testing.T) {
	secrets := baseExpireSecrets()
	policy := ExpiryPolicy{WarnBefore: 30 * time.Minute}
	report, err := CheckExpiry(secrets, policy, time.Now().UTC())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Warning) != 1 || report.Warning[0] != "API_KEY_EXPIRES_AT" {
		t.Errorf("expected API_KEY_EXPIRES_AT in warning, got %v", report.Warning)
	}
}

func TestCheckExpiry_IgnoresNonExpiryKeys(t *testing.T) {
	secrets := map[string]string{"DB_PASSWORD": "secret", "API_KEY": "key"}
	policy := ExpiryPolicy{WarnBefore: time.Hour}
	report, err := CheckExpiry(secrets, policy, time.Now().UTC())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Results) != 0 {
		t.Errorf("expected no results, got %d", len(report.Results))
	}
}

func TestCheckExpiry_NilSecrets_Errors(t *testing.T) {
	_, err := CheckExpiry(nil, ExpiryPolicy{}, time.Now())
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestCheckExpiry_InvalidTimestamp_Errors(t *testing.T) {
	secrets := map[string]string{"BAD_EXPIRES_AT": "not-a-date"}
	_, err := CheckExpiry(secrets, ExpiryPolicy{}, time.Now())
	if err == nil {
		t.Fatal("expected error for invalid timestamp")
	}
}

func TestCheckExpiry_Summary(t *testing.T) {
	secrets := baseExpireSecrets()
	policy := ExpiryPolicy{WarnBefore: 30 * time.Minute}
	report, err := CheckExpiry(secrets, policy, time.Now().UTC())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := report.Summary()
	if sum == "" {
		t.Error("expected non-empty summary")
	}
}
