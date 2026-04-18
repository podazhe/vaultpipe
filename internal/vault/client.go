package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	vc *vaultapi.Client
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
	RoleID  string
	SecretID string
}

// NewClient creates a new Vault client from the given config.
// If Token is empty, it falls back to the VAULT_TOKEN environment variable.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()
	if cfg.Address != "" {
		vcfg.Address = cfg.Address
	}

	vc, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to create client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}

	if token != "" {
		vc.SetToken(token)
	} else if cfg.RoleID != "" && cfg.SecretID != "" {
		t, err := appRoleLogin(vc, cfg.RoleID, cfg.SecretID)
		if err != nil {
			return nil, err
		}
		vc.SetToken(t)
	} else {
		return nil, fmt.Errorf("vault: no authentication method provided")
	}

	return &Client{vc: vc}, nil
}

// ReadSecrets reads key/value secrets at the given path and returns them as a map.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.vc.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to read path %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault: no secret found at path %q", path)
	}

	data, ok := secret.Data["data"]
	if !ok {
		// KV v1 — data is at the top level
		data = secret.Data
	}

	raw, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("vault: unexpected data format at path %q", path)
	}

	result := make(map[string]string, len(raw))
	for k, v := range raw {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}

func appRoleLogin(vc *vaultapi.Client, roleID, secretID string) (string, error) {
	data := map[string]interface{}{"role_id": roleID, "secret_id": secretID}
	secret, err := vc.Logical().Write("auth/approle/login", data)
	if err != nil {
		return "", fmt.Errorf("vault: approle login failed: %w", err)
	}
	if secret == nil || secret.Auth == nil {
		return "", fmt.Errorf("vault: approle login returned no auth info")
	}
	return secret.Auth.ClientToken, nil
}
