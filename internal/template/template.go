// Package template provides functionality for rendering secret values
// into user-defined output templates, allowing flexible formatting of
// synced secrets beyond the default .env file format.
package template

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	gotemplate "text/template"
)

// Renderer renders secrets using a Go text/template string.
type Renderer struct {
	tmpl *gotemplate.Template
}

// New parses the given template string and returns a Renderer.
// Returns an error if the template is invalid.
func New(tmplStr string) (*Renderer, error) {
	funcMap := gotemplate.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trim":  strings.TrimSpace,
	}

	t, err := gotemplate.New("vaultpull").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}
	return &Renderer{tmpl: t}, nil
}

// NewFromFile reads a template from the given file path and returns a Renderer.
func NewFromFile(path string) (*Renderer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading template file: %w", err)
	}
	return New(string(data))
}

// Render executes the template with the provided secrets map and returns
// the rendered output as a string.
func (r *Renderer) Render(secrets map[string]string) (string, error) {
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, secrets); err != nil {
		return "", fmt.Errorf("template render error: %w", err)
	}
	return buf.String(), nil
}

// RenderToFile executes the template and writes the result to the given path.
func (r *Renderer) RenderToFile(secrets map[string]string, path string) error {
	out, err := r.Render(secrets)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(out), 0600); err != nil {
		return fmt.Errorf("writing rendered template to %s: %w", path, err)
	}
	return nil
}
