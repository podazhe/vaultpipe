package vault

import (
	"testing"
	"time"
)

func TestSecretCache_SetAndGet(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	data := map[string]string{"API_KEY": "abc123"}
	c.Set("secret/myapp", data)

	got := c.Get("secret/myapp")
	if got == nil {
		t.Fatal("expected cached value, got nil")
	}
	if got["API_KEY"] != "abc123" {
		t.Errorf("expected abc123, got %s", got["API_KEY"])
	}
}

func TestSecretCache_MissingKey(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	if got := c.Get("secret/missing"); got != nil {
		t.Errorf("expected nil for missing key, got %v", got)
	}
}

func TestSecretCache_Expiry(t *testing.T) {
	c := NewSecretCache(50 * time.Millisecond)
	c.Set("secret/short", map[string]string{"K": "V"})

	time.Sleep(100 * time.Millisecond)

	if got := c.Get("secret/short"); got != nil {
		t.Error("expected expired entry to return nil")
	}
}

func TestSecretCache_NoExpiry(t *testing.T) {
	c := NewSecretCache(0)
	c.Set("secret/forever", map[string]string{"K": "V"})
	time.Sleep(20 * time.Millisecond)
	if got := c.Get("secret/forever"); got == nil {
		t.Error("expected entry with TTL=0 to never expire")
	}
}

func TestSecretCache_Invalidate(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	c.Set("secret/myapp", map[string]string{"K": "V"})
	c.Invalidate("secret/myapp")
	if got := c.Get("secret/myapp"); got != nil {
		t.Error("expected nil after invalidation")
	}
}

func TestSecretCache_Flush(t *testing.T) {
	c := NewSecretCache(5 * time.Minute)
	c.Set("secret/a", map[string]string{"A": "1"})
	c.Set("secret/b", map[string]string{"B": "2"})
	c.Flush()
	if c.Get("secret/a") != nil || c.Get("secret/b") != nil {
		t.Error("expected all entries cleared after flush")
	}
}
