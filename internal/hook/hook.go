// Package hook provides a mechanism for running user-defined shell commands
// before and after a secrets sync operation.
package hook

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Stage indicates when a hook should be executed.
type Stage string

const (
	StagePre  Stage = "pre"
	StagePost Stage = "post"

	defaultTimeout = 30 * time.Second
)

// Hook represents a single shell command to run at a given stage.
type Hook struct {
	Stage   Stage
	Command string
	Timeout time.Duration
}

// Runner executes hooks for a given stage.
type Runner struct {
	execCommand func(ctx context.Context, name string, args ...string) *exec.Cmd
}

// NewRunner returns a Runner using the real exec.CommandContext.
func NewRunner() *Runner {
	return &Runner{execCommand: exec.CommandContext}
}

// newRunnerWithExec returns a Runner with an injected executor (for testing).
func newRunnerWithExec(fn func(ctx context.Context, name string, args ...string) *exec.Cmd) *Runner {
	return &Runner{execCommand: fn}
}

// Run executes all hooks matching the given stage in order.
// It returns the first error encountered.
func (r *Runner) Run(ctx context.Context, hooks []Hook, stage Stage) error {
	for _, h := range hooks {
		if h.Stage != stage {
			continue
		}
		if err := r.runOne(ctx, h); err != nil {
			return fmt.Errorf("hook %q failed: %w", h.Command, err)
		}
	}
	return nil
}

func (r *Runner) runOne(ctx context.Context, h Hook) error {
	timeout := h.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	parts := strings.Fields(h.Command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := r.execCommand(ctx, parts[0], parts[1:]...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
