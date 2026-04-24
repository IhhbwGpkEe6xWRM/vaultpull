// Package envflow provides a sequential pipeline that applies a series of
// named transformation stages to a secret map, collecting stage results and
// any errors encountered along the way.
package envflow

import "fmt"

// Stage is a single named step in the pipeline.
type Stage struct {
	Name    string
	Apply   func(map[string]string) (map[string]string, error)
}

// Result holds the outcome of a single pipeline stage.
type Result struct {
	Stage  string
	Keys   int
	Err    error
}

// Pipeline executes an ordered list of stages.
type Pipeline struct {
	stages []Stage
}

// New returns a Pipeline containing the given stages.
// It returns an error if any stage has an empty name or nil Apply func.
func New(stages []Stage) (*Pipeline, error) {
	for i, s := range stages {
		if s.Name == "" {
			return nil, fmt.Errorf("stage %d: name must not be empty", i)
		}
		if s.Apply == nil {
			return nil, fmt.Errorf("stage %q: Apply func must not be nil", s.Name)
		}
	}
	return &Pipeline{stages: stages}, nil
}

// Run executes each stage in order, passing the output of one stage as the
// input to the next. It stops on the first error and returns all results
// collected up to that point together with the final secret map.
func (p *Pipeline) Run(secrets map[string]string) (map[string]string, []Result, error) {
	current := copyMap(secrets)
	results := make([]Result, 0, len(p.stages))

	for _, s := range p.stages {
		out, err := s.Apply(current)
		r := Result{Stage: s.Name, Err: err}
		if err == nil {
			r.Keys = len(out)
			current = out
		}
		results = append(results, r)
		if err != nil {
			return current, results, fmt.Errorf("stage %q: %w", s.Name, err)
		}
	}
	return current, results, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
