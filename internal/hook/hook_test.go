package hook_test

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/hook"
)

func realExec(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}

func TestRun_NoMatchingStage(t *testing.T) {
	r := hook.NewRunner()
	hooks := []hook.Hook{
		{Stage: hook.StagePost, Command: "false"},
	}
	if err := r.Run(context.Background(), hooks, hook.StagePre); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRun_SuccessfulCommand(t *testing.T) {
	r := hook.NewRunner()
	hooks := []hook.Hook{
		{Stage: hook.StagePre, Command: "echo hello"},
	}
	if err := r.Run(context.Background(), hooks, hook.StagePre); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_FailingCommand(t *testing.T) {
	r := hook.NewRunner()
	hooks := []hook.Hook{
		{Stage: hook.StagePre, Command: "false"},
	}
	if err := r.Run(context.Background(), hooks, hook.StagePre); err == nil {
		t.Fatal("expected error from failing command")
	}
}

func TestRun_EmptyCommand(t *testing.T) {
	r := hook.NewRunner()
	hooks := []hook.Hook{
		{Stage: hook.StagePre, Command: ""},
	}
	if err := r.Run(context.Background(), hooks, hook.StagePre); err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestRun_Timeout(t *testing.T) {
	r := hook.NewRunner()
	hooks := []hook.Hook{
		{
			Stage:   hook.StagePre,
			Command: "sleep 10",
			Timeout: 50 * time.Millisecond,
		},
	}
	if err := r.Run(context.Background(), hooks, hook.StagePre); err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestRun_MultipleHooks_StopsOnFirstError(t *testing.T) {
	called := 0
	r := hook.NewRunner()
	hooks := []hook.Hook{
		{Stage: hook.StagePre, Command: "false"},
		{Stage: hook.StagePre, Command: "echo second"},
	}
	_ = called
	if err := r.Run(context.Background(), hooks, hook.StagePre); err == nil {
		t.Fatal("expected error from first failing hook")
	}
}

func TestRun_PostHooksRunSeparately(t *testing.T) {
	r := hook.NewRunner()
	hooks := []hook.Hook{
		{Stage: hook.StagePre, Command: "false"},
		{Stage: hook.StagePost, Command: "echo ok"},
	}
	if err := r.Run(context.Background(), hooks, hook.StagePost); err != nil {
		t.Fatalf("post hook should succeed independently: %v", err)
	}
}
