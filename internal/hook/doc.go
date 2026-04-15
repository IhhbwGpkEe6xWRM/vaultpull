// Package hook provides pre- and post-sync lifecycle hooks for vaultpull.
//
// Hooks allow users to run arbitrary shell commands before secrets are pulled
// from Vault and after the local .env file has been written. This is useful
// for tasks such as reloading a process, notifying a service, or validating
// the resulting environment.
//
// Usage:
//
//	runner := hook.NewRunner()
//
//	hooks := []hook.Hook{
//		{Stage: hook.StagePre,  Command: "echo starting sync"},
//		{Stage: hook.StagePost, Command: "systemctl reload myapp"},
//	}
//
//	if err := runner.Run(ctx, hooks, hook.StagePre); err != nil {
//		log.Fatal(err)
//	}
//
// Hook commands are split on whitespace and executed directly without a shell.
// Each hook has an optional Timeout; if unset, a default of 30 seconds applies.
package hook
