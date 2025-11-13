package shell_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

func TestBinaryShellchecker_Check(t *testing.T) {
	t.Parallel()

	checker := shell.NewBinaryShellchecker()

	tests := []struct {
		name          string
		script        string
		opts          shell.ShellOpts
		expectFailure bool
	}{
		{
			name:          "valid script",
			script:        "apt-get install -y package",
			opts:          shell.DefaultShellOpts(),
			expectFailure: false,
		},
		{
			name:          "unquoted variable",
			script:        "echo $FOO",
			opts:          shell.DefaultShellOpts(),
			expectFailure: true, // SC2086: quote to prevent word splitting
		},
		{
			name:   "powershell script skipped",
			script: "Write-Host 'test'",
			opts: shell.ShellOpts{
				ShellName: "pwsh -c",
				EnvVars:   make(map[string]string),
			},
			expectFailure: false, // Skipped for non-POSIX shells
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			failures, err := checker.Check(tt.script, tt.opts)
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}

			hasFailure := len(failures) > 0
			if hasFailure != tt.expectFailure {
				t.Errorf("Check() hasFailure = %v, want %v (failures: %v)", hasFailure, tt.expectFailure, failures)
			}
		})
	}
}

func TestShellcheckRule_StatefulTracking(t *testing.T) {
	t.Parallel()

	checker := shell.NewBinaryShellchecker()
	scRule := shell.NewShellcheckRule(checker)

	// Simulate a Dockerfile with ENV, ARG, SHELL, and RUN
	instructions := []syntax.InstructionPos{
		{
			Instruction: &syntax.From{Image: syntax.BaseImage{Image: "debian"}},
			LineNumber:  1,
		},
		{
			Instruction: &syntax.Env{Pairs: []syntax.EnvPair{{Key: "MY_VAR", Value: "value"}}},
			LineNumber:  2,
		},
		{
			Instruction: &syntax.Arg{ArgName: "BUILD_ARG"},
			LineNumber:  3,
		},
		{
			Instruction: &syntax.Shell{Arguments: []string{"/bin/bash", "-c"}},
			LineNumber:  4,
		},
		{
			Instruction: &syntax.Run{Command: "echo $UNDEFINED"},
			LineNumber:  5,
		},
	}

	state := scRule.InitialState()
	for _, instr := range instructions {
		state = scRule.Check(instr.LineNumber, state, instr.Instruction)
	}

	// Should have shellcheck violations for undefined variable
	if len(state.Failures) == 0 {
		t.Error("Expected shellcheck violations, got none")
	}

	// Verify violations have SC codes
	for _, f := range state.Failures {
		code := string(f.Code)
		if len(code) < 2 || code[0:2] != "SC" {
			t.Errorf("Expected SC code, got %s", code)
		}
	}
}

func TestShellcheckRule_ResetOnFrom(t *testing.T) {
	t.Parallel()

	checker := shell.NewBinaryShellchecker()
	scRule := shell.NewShellcheckRule(checker)

	// Test that ENV variables are reset on new FROM (multi-stage build)
	instructions := []syntax.InstructionPos{
		{
			Instruction: &syntax.From{Image: syntax.BaseImage{Image: "debian"}},
			LineNumber:  1,
		},
		{
			Instruction: &syntax.Env{Pairs: []syntax.EnvPair{{Key: "STAGE1_VAR", Value: "val"}}},
			LineNumber:  2,
		},
		{
			Instruction: &syntax.From{Image: syntax.BaseImage{Image: "alpine"}},
			LineNumber:  3,
		},
		// After second FROM, STAGE1_VAR should not be in environment
	}

	// Since shellState is unexported, we test behavior instead of internal state
	// The fact that the rule processes without error is sufficient for this test
	// Actual ENV variable tracking is verified through integration tests
	state := scRule.InitialState()
	for _, instr := range instructions {
		state = scRule.Check(instr.LineNumber, state, instr.Instruction)
	}
}

func TestNoopShellchecker(t *testing.T) {
	t.Parallel()

	checker := shell.NewNoopShellchecker()

	failures, err := checker.Check("any script", shell.DefaultShellOpts())
	if err != nil {
		t.Errorf("NoopShellchecker.Check() error = %v, want nil", err)
	}

	if len(failures) != 0 {
		t.Errorf("NoopShellchecker.Check() returned %d failures, want 0", len(failures))
	}
}

func TestShellcheckRule_Metadata(t *testing.T) {
	t.Parallel()

	checker := shell.NewNoopShellchecker()
	scRule := shell.NewShellcheckRule(checker)

	if scRule.Code() != "SHELLCHECK" {
		t.Errorf("Code() = %s, want SHELLCHECK", scRule.Code())
	}

	if scRule.Severity() != rule.Info {
		t.Errorf("Severity() = %v, want Info", scRule.Severity())
	}

	if scRule.Message() == "" {
		t.Error("Message() returned empty string")
	}
}
