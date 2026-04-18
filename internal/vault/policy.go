package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// PolicyChecker verifies that the current Vault token has the required
// capabilities on a given path before attempting secret reads.
type PolicyChecker struct {
	client *vaultapi.Client
}

// NewPolicyChecker creates a PolicyChecker wrapping the provided Vault client.
func NewPolicyChecker(client *vaultapi.Client) *PolicyChecker {
	return &PolicyChecker{client: client}
}

// Capability represents a Vault capability string.
type Capability string

const (
	CapRead   Capability = "read"
	CapList   Capability = "list"
	CapCreate Capability = "create"
	CapUpdate Capability = "update"
	CapDelete Capability = "delete"
	CapDeny   Capability = "deny"
	CapRoot   Capability = "root"
)

// CheckPath returns the capabilities the current token holds on path.
func (p *PolicyChecker) CheckPath(ctx context.Context, path string) ([]Capability, error) {
	body := map[string]interface{}{"path": path}
	secret, err := p.client.Logical().WriteWithContext(ctx, "sys/capabilities-self", body)
	if err != nil {
		return nil, fmt.Errorf("capabilities-self request failed for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("empty response from capabilities-self for %q", path)
	}

	raw, ok := secret.Data[path]
	if !ok {
		// Vault may return capabilities under the key "capabilities"
		raw, ok = secret.Data["capabilities"]
		if !ok {
			return nil, fmt.Errorf("no capability data returned for %q", path)
		}
	}

	items, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected capability format for %q", path)
	}

	caps := make([]Capability, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			caps = append(caps, Capability(strings.ToLower(s)))
		}
	}
	return caps, nil
}

// HasCapability returns true when the token holds the requested capability on path.
func (p *PolicyChecker) HasCapability(ctx context.Context, path string, want Capability) (bool, error) {
	caps, err := p.CheckPath(ctx, path)
	if err != nil {
		return false, err
	}
	for _, c := range caps {
		if c == CapDeny {
			return false, nil
		}
		if c == want || c == CapRoot {
			return true, nil
		}
	}
	return false, nil
}
