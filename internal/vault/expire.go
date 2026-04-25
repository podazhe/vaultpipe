package vault

import (
	"errors"
	"fmt"
	"time"
)

// ExpiryPolicy defines how expiry is evaluated for a secret.
type ExpiryPolicy struct {
	// WarnBefore is the duration before expiry at which a warning is emitted.
	WarnBefore time.Duration
	// ErrorBefore is the duration before expiry at which an error is returned.
	ErrorBefore time.Duration
}

// ExpiryResult holds the outcome of an expiry check for a single key.
type ExpiryResult struct {
	Key       string
	ExpiresAt time.Time
	Expired   bool
	Warning   bool
	Message   string
}

// ExpiryReport is the full result of CheckExpiry.
type ExpiryReport struct {
	Results []ExpiryResult
	Expired []string
	Warning []string
}

// Summary returns a human-readable summary of the expiry report.
func (r *ExpiryReport) Summary() string {
	return fmt.Sprintf("%d expired, %d expiring soon (of %d checked)",
		len(r.Expired), len(r.Warning), len(r.Results))
}

// CheckExpiry evaluates expiry metadata stored as RFC3339 timestamps in the
// secrets map under keys with the suffix "_EXPIRES_AT". It returns an
// ExpiryReport describing which secrets are expired or nearing expiry.
func CheckExpiry(secrets map[string]string, policy ExpiryPolicy, now time.Time) (*ExpiryReport, error) {
	if secrets == nil {
		return nil, errors.New("CheckExpiry: secrets map must not be nil")
	}
	report := &ExpiryReport{}
	const suffix = "_EXPIRES_AT"
	for k, v := range secrets {
		if len(k) <= len(suffix) || k[len(k)-len(suffix):] != suffix {
			continue
		}
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return nil, fmt.Errorf("CheckExpiry: invalid timestamp for key %q: %w", k, err)
		}
		res := ExpiryResult{Key: k, ExpiresAt: t}
		switch {
		case !now.Before(t):
			res.Expired = true
			res.Message = fmt.Sprintf("%s expired at %s", k, t.Format(time.RFC3339))
			report.Expired = append(report.Expired, k)
		case policy.ErrorBefore > 0 && t.Sub(now) <= policy.ErrorBefore:
			res.Warning = true
			res.Message = fmt.Sprintf("%s expires in %s", k, t.Sub(now).Round(time.Second))
			report.Warning = append(report.Warning, k)
		case policy.WarnBefore > 0 && t.Sub(now) <= policy.WarnBefore:
			res.Warning = true
			res.Message = fmt.Sprintf("%s expires in %s", k, t.Sub(now).Round(time.Second))
			report.Warning = append(report.Warning, k)
		}
		report.Results = append(report.Results, res)
	}
	return report, nil
}
