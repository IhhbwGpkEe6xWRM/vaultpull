package vault

import (
	"context"
	"fmt"
	"net/http"

	"github.com/your-org/vaultpull/internal/retry"
)

// retryableStatusCodes are HTTP status codes from Vault that warrant a retry.
var retryableStatusCodes = map[int]bool{
	http.StatusTooManyRequests:     true,
	http.StatusServiceUnavailable:  true,
	http.StatusGatewayTimeout:      true,
	http.StatusInternalServerError: true,
}

// ReadSecretsWithRetry wraps ReadSecrets with exponential-backoff retry logic.
// HTTP 403/404 responses from Vault are treated as permanent (non-retryable)
// to avoid hammering the server on permission or path errors.
func (c *Client) ReadSecretsWithRetry(ctx context.Context, path string, cfg retry.Config) (map[string]string, error) {
	var result map[string]string

	err := retry.Do(ctx, cfg, func() error {
		secrets, err := c.ReadSecrets(ctx, path)
		if err != nil {
			if isPermanentVaultError(err) {
				return retry.Permanent(err)
			}
			return err
		}
		result = secrets
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", path, err)
	}
	return result, nil
}

// isPermanentVaultError returns true for errors that should not be retried.
// Currently treats 403 Forbidden and 404 Not Found as permanent failures.
func isPermanentVaultError(err error) bool {
	if err == nil {
		return false
	}
	var ve *VaultHTTPError
	if ok := errorAs(err, &ve); ok {
		switch ve.StatusCode {
		case http.StatusForbidden, http.StatusNotFound:
			return true
		}
	}
	return false
}

// VaultHTTPError represents an HTTP-level error returned by the Vault API.
type VaultHTTPError struct {
	StatusCode int
	Message    string
}

func (e *VaultHTTPError) Error() string {
	return fmt.Sprintf("vault http %d: %s", e.StatusCode, e.Message)
}

// errorAs is a thin wrapper around errors.As for testability.
func errorAs(err error, target interface{ Unwrap() error }) bool {
	type asInterface interface {
		As(interface{}) bool
	}
	if a, ok := err.(asInterface); ok {
		return a.As(target)
	}
	import_errors_as(err, target)
	return false
}
