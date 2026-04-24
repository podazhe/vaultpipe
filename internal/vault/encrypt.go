package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// EncryptOptions controls which keys are encrypted and how.
type EncryptOptions struct {
	// Keys is an explicit list of keys to encrypt. If empty, all keys are encrypted.
	Keys []string
	// Passthrough leaves values unmodified (useful for dry-run inspection).
	Passthrough bool
}

// EncryptSecrets encrypts secret values using AES-256-GCM with the supplied
// 32-byte key. Encrypted values are base64-encoded and prefixed with "enc:" so
// they can be identified and decrypted later.
func EncryptSecrets(secrets map[string]string, aesKey []byte, opts EncryptOptions) (map[string]string, error) {
	if len(aesKey) != 32 {
		return nil, errors.New("encrypt: AES key must be exactly 32 bytes")
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("encrypt: create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("encrypt: create GCM: %w", err)
	}

	targets := buildTargetSet(opts.Keys)

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if opts.Passthrough || (!isTarget(k, targets)) {
			out[k] = v
			continue
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, fmt.Errorf("encrypt: generate nonce for key %q: %w", k, err)
		}

		ciphertext := gcm.Seal(nonce, nonce, []byte(v), nil)
		out[k] = "enc:" + base64.StdEncoding.EncodeToString(ciphertext)
	}

	return out, nil
}

// DecryptSecrets reverses EncryptSecrets. Values that do not carry the "enc:"
// prefix are passed through unchanged.
func DecryptSecrets(secrets map[string]string, aesKey []byte) (map[string]string, error) {
	if len(aesKey) != 32 {
		return nil, errors.New("decrypt: AES key must be exactly 32 bytes")
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt: create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("decrypt: create GCM: %w", err)
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if len(v) < 4 || v[:4] != "enc:" {
			out[k] = v
			continue
		}

		raw, err := base64.StdEncoding.DecodeString(v[4:])
		if err != nil {
			return nil, fmt.Errorf("decrypt: base64 decode key %q: %w", k, err)
		}

		if len(raw) < gcm.NonceSize() {
			return nil, fmt.Errorf("decrypt: ciphertext too short for key %q", k)
		}

		nonce, ct := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
		plain, err := gcm.Open(nil, nonce, ct, nil)
		if err != nil {
			return nil, fmt.Errorf("decrypt: open key %q: %w", k, err)
		}

		out[k] = string(plain)
	}

	return out, nil
}

func buildTargetSet(keys []string) map[string]struct{} {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}

func isTarget(key string, targets map[string]struct{}) bool {
	if targets == nil {
		return true // no filter → all keys are targets
	}
	_, ok := targets[key]
	return ok
}
