package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAuthTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/approle/login":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"auth":{"client_token":"test-token","renewable":true,"lease_duration":3600}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestAuthenticator_TokenMethod(t *testing.T) {
	c, err := NewClient("http://127.0.0.1", "tok")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	auth := NewAuthenticator(c)
	res, err := auth.Authenticate(context.Background(), AuthConfig{
		Method: AuthToken,
		Token:  "mytoken",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Token != "mytoken" {
		t.Errorf("expected token mytoken, got %s", res.Token)
	}
}

func TestAuthenticator_TokenMethod_Missing(t *testing.T) {
	c, _ := NewClient("http://127.0.0.1", "tok")
	auth := NewAuthenticator(c)
	_, err := auth.Authenticate(context.Background(), AuthConfig{Method: AuthToken})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestAuthenticator_AppRole_MissingFields(t *testing.T) {
	c, _ := NewClient("http://127.0.0.1", "tok")
	auth := NewAuthenticator(c)
	_, err := auth.Authenticate(context.Background(), AuthConfig{
		Method: AuthAppRole,
		RoleID: "role",
	})
	if err == nil {
		t.Fatal("expected error for missing secret_id")
	}
}

func TestAuthenticator_UnsupportedMethod(t *testing.T) {
	c, _ := NewClient("http://127.0.0.1", "tok")
	auth := NewAuthenticator(c)
	_, err := auth.Authenticate(context.Background(), AuthConfig{Method: "ldap"})
	if err == nil {
		t.Fatal("expected error for unsupported method")
	}
}

func TestAuthenticator_Kubernetes_MissingRole(t *testing.T) {
	c, _ := NewClient("http://127.0.0.1", "tok")
	auth := NewAuthenticator(c)
	_, err := auth.Authenticate(context.Background(), AuthConfig{Method: AuthKubernetes})
	if err == nil {
		t.Fatal("expected error for missing role")
	}
}
