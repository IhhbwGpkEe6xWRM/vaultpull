package vault

import (
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with namespace-aware secret fetching.
type Client struct {
	api       *vaultapi.Client
	namespace string
}

// NewClient creates a new Vault client configured with the given address, token, and namespace.
func NewClient(address, token, namespace string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	api.SetToken(token)

	return &Client{
		api:       api,
		namespace: strings.Trim(namespace, "/"),
	}, nil
}

// ReadSecrets reads key-value secrets from the given path, prepending the
// configured namespace if one is set.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	fullPath := c.resolvePath(path)

	secret, err := c.api.Logical().Read(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets at %q: %w", fullPath, err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secrets found at path %q", fullPath)
	}

	// Support both KV v1 (data at top level) and KV v2 (data nested under "data").
	data, ok := secret.Data["data"]
	if ok {
		if nested, ok := data.(map[string]interface{}); ok {
			return toStringMap(nested), nil
		}
	}

	return toStringMap(secret.Data), nil
}

// resolvePath prepends the namespace to the secret path when a namespace is configured.
func (c *Client) resolvePath(path string) string {
	path = strings.TrimPrefix(path, "/")
	if c.namespace == "" {
		return path
	}
	return c.namespace + "/" + path
}

// toStringMap converts map[string]interface{} to map[string]string, skipping non-string values.
func toStringMap(in map[string]interface{}) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
}
