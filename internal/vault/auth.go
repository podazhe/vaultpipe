package vault

import (
	"context"
	"fmt"
	"time"
)

// AuthMethod represents a Vault authentication method.
type AuthMethod string

const (
	AuthToken    AuthMethod = "token"
	AuthAppRole  AuthMethod = "approle"
	AuthKubernetes AuthMethod = "kubernetes"
)

// AuthConfig holds configuration for a Vault auth method.
type AuthConfig struct {
	Method     AuthMethod
	Token      string
	RoleID     string
	SecretID   string
	Role       string // kubernetes role
	JWTPath    string // path to kubernetes JWT
	MountPath  string // custom mount path, defaults to method name
}

// AuthResult holds the result of a successful authentication.
type AuthResult struct {
	Token     string
	Renewable bool
	LeaseDur  time.Duration
}

// Authenticator performs authentication against Vault.
type Authenticator struct {
	client *Client
}

// NewAuthenticator creates a new Authenticator.
func NewAuthenticator(c *Client) *Authenticator {
	return &Authenticator{client: c}
}

// Authenticate authenticates using the provided AuthConfig and returns an AuthResult.
func (a *Authenticator) Authenticate(ctx context.Context, cfg AuthConfig) (*AuthResult, error) {
	switch cfg.Method {
	case AuthToken:
		if cfg.Token == "" {
			return nil, fmt.Errorf("auth: token method requires a token")
		}
		return &AuthResult{Token: cfg.Token, Renewable: false}, nil
	case AuthAppRole:
		return a.appRole(ctx, cfg)
	case AuthKubernetes:
		return a.kubernetes(ctx, cfg)
	default:
		return nil, fmt.Errorf("auth: unsupported method %q", cfg.Method)
	}
}

func (a *Authenticator) appRole(ctx context.Context, cfg AuthConfig) (*AuthResult, error) {
	if cfg.RoleID == "" || cfg.SecretID == "" {
		return nil, fmt.Errorf("auth: approle requires role_id and secret_id")
	}
	mount := cfg.MountPath
	if mount == "" {
		mount = "approle"
	}
	secret, err := a.client.vault.Auth().Login(ctx, nil)
	_ = secret
	if err != nil {
		return nil, fmt.Errorf("auth: approle login failed: %w", err)
	}
	return &AuthResult{Token: "", Renewable: true}, nil
}

func (a *Authenticator) kubernetes(ctx context.Context, cfg AuthConfig) (*AuthResult, error) {
	if cfg.Role == "" {
		return nil, fmt.Errorf("auth: kubernetes requires a role")
	}
	return nil, fmt.Errorf("auth: kubernetes not yet implemented")
}
