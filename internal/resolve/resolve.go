// Package resolve provides path resolution utilities for mapping
// logical secret names to their full Vault paths, supporting both
// KV v1 and KV v2 mount structures.
package resolve

import (
	"fmt"
	"strings"
)

// MountVersion represents the KV secrets engine version.
type MountVersion int

const (
	KVv1 MountVersion = 1
	KVv2 MountVersion = 2
)

// Resolver builds full Vault API paths from a mount and secret path.
type Resolver struct {
	mount   string
	version MountVersion
}

// New returns a Resolver for the given mount point and KV version.
func New(mount string, version MountVersion) (*Resolver, error) {
	mount = strings.Trim(mount, "/")
	if mount == "" {
		return nil, fmt.Errorf("resolve: mount must not be empty")
	}
	if version != KVv1 && version != KVv2 {
		return nil, fmt.Errorf("resolve: unsupported KV version %d", version)
	}
	return &Resolver{mount: mount, version: version}, nil
}

// DataPath returns the API path used to read a secret.
func (r *Resolver) DataPath(secretPath string) string {
	secretPath = strings.Trim(secretPath, "/")
	if r.version == KVv2 {
		return fmt.Sprintf("%s/data/%s", r.mount, secretPath)
	}
	return fmt.Sprintf("%s/%s", r.mount, secretPath)
}

// MetadataPath returns the API path used to list or inspect secret metadata (KV v2 only).
func (r *Resolver) MetadataPath(secretPath string) string {
	secretPath = strings.Trim(secretPath, "/")
	if r.version == KVv2 {
		return fmt.Sprintf("%s/metadata/%s", r.mount, secretPath)
	}
	// KV v1 has no separate metadata path; return data path.
	return fmt.Sprintf("%s/%s", r.mount, secretPath)
}

// Version returns the configured MountVersion.
func (r *Resolver) Version() MountVersion { return r.version }

// Mount returns the configured mount point.
func (r *Resolver) Mount() string { return r.mount }
