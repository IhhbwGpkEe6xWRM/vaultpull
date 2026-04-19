// Package overlay provides a layered secret resolution strategy that merges
// secrets from multiple sources in priority order. Later layers override
// earlier ones, allowing environment-specific overrides on top of base secrets.
package overlay

import (
	"context"
	"fmt"
)

// Reader is the interface for reading secrets from a single source.
type Reader interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Layer represents a named secret source with an associated path.
type Layer struct {
	Name   string
	Reader Reader
	Path   string
}

// Resolver merges secrets from multiple layers in order, with later layers
// taking precedence over earlier ones.
type Resolver struct {
	layers []Layer
}

// New creates a Resolver from the given layers. Layers are applied in order;
// the last layer has the highest priority.
func New(layers ...Layer) (*Resolver, error) {
	if len(layers) == 0 {
		return nil, fmt.Errorf("overlay: at least one layer is required")
	}
	for i, l := range layers {
		if l.Reader == nil {
			return nil, fmt.Errorf("overlay: layer %d (%q) has nil reader", i, l.Name)
		}
		if l.Path == "" {
			return nil, fmt.Errorf("overlay: layer %d (%q) has empty path", i, l.Name)
		}
	}
	return &Resolver{layers: layers}, nil
}

// Resolve reads secrets from all layers and merges them. Secrets from later
// layers overwrite keys from earlier layers. The returned map is a fresh copy.
func (r *Resolver) Resolve(ctx context.Context) (map[string]string, error) {
	merged := make(map[string]string)
	for _, layer := range r.layers {
		secrets, err := layer.Reader.ReadSecrets(ctx, layer.Path)
		if err != nil {
			return nil, fmt.Errorf("overlay: reading layer %q at %q: %w", layer.Name, layer.Path, err)
		}
		for k, v := range secrets {
			merged[k] = v
		}
	}
	return merged, nil
}

// Sources returns the names of all registered layers in resolution order.
func (r *Resolver) Sources() []string {
	names := make([]string, len(r.layers))
	for i, l := range r.layers {
		names[i] = l.Name
	}
	return names
}
