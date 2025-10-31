package shell

import (
	"testing"
)

func TestParseShell(t *testing.T) {
	t.Run("simple command", func(t *testing.T) {
		ps, err := ParseShell("apt-get update")
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		if len(ps.PresentCommands) != 1 {
			t.Fatalf("expected 1 command, got %d", len(ps.PresentCommands))
		}

		cmd := ps.PresentCommands[0]
		if cmd.Name != "apt-get" {
			t.Errorf("expected command name 'apt-get', got '%s'", cmd.Name)
		}

		args := GetArgs(cmd)
		if len(args) != 1 || args[0] != "update" {
			t.Errorf("expected args ['update'], got %v", args)
		}
	})

	t.Run("command with flags", func(t *testing.T) {
		ps, err := ParseShell("apt-get install -y vim")
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		cmd := ps.PresentCommands[0]
		if cmd.Name != "apt-get" {
			t.Errorf("expected command name 'apt-get', got '%s'", cmd.Name)
		}

		if !HasFlag("y", cmd) {
			t.Error("expected command to have -y flag")
		}

		argsNoFlags := GetArgsNoFlags(cmd)
		expectedArgs := []string{"install", "vim"}
		if len(argsNoFlags) != len(expectedArgs) {
			t.Errorf("expected %d args, got %d: %v", len(expectedArgs), len(argsNoFlags), argsNoFlags)
		}
	})

	t.Run("chained commands with &&", func(t *testing.T) {
		ps, err := ParseShell("apt-get update && apt-get install -y vim")
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		if len(ps.PresentCommands) != 2 {
			t.Fatalf("expected 2 commands, got %d", len(ps.PresentCommands))
		}

		names := FindCommandNames(ps)
		if len(names) != 2 || names[0] != "apt-get" || names[1] != "apt-get" {
			t.Errorf("expected ['apt-get', 'apt-get'], got %v", names)
		}
	})

	t.Run("chained commands with semicolon", func(t *testing.T) {
		ps, err := ParseShell("echo foo; echo bar")
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		if len(ps.PresentCommands) != 2 {
			t.Fatalf("expected 2 commands, got %d", len(ps.PresentCommands))
		}
	})

	t.Run("long flag with equals", func(t *testing.T) {
		ps, err := ParseShell("useradd --uid=1000 user")
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		cmd := ps.PresentCommands[0]
		if !HasFlag("uid", cmd) {
			t.Error("expected command to have --uid flag")
		}
	})

	t.Run("multiple short flags", func(t *testing.T) {
		ps, err := ParseShell("tar -xzf archive.tar.gz")
		if err != nil {
			t.Fatalf("failed to parse: %v", err)
		}

		cmd := ps.PresentCommands[0]
		if !HasFlag("x", cmd) || !HasFlag("z", cmd) || !HasFlag("f", cmd) {
			t.Error("expected command to have -x, -z, and -f flags")
		}
	})
}

func TestCmdHasArgs(t *testing.T) {
	ps, _ := ParseShell("apt-get install vim")
	cmd := ps.PresentCommands[0]

	if !CmdHasArgs("apt-get", []string{"install"}, cmd) {
		t.Error("expected CmdHasArgs to return true for apt-get with install")
	}

	if CmdHasArgs("apt-get", []string{"update"}, cmd) {
		t.Error("expected CmdHasArgs to return false for apt-get with update")
	}

	if CmdHasArgs("yum", []string{"install"}, cmd) {
		t.Error("expected CmdHasArgs to return false for different command")
	}
}

func TestCountCommands(t *testing.T) {
	tests := []struct {
		name     string
		script   string
		expected int
	}{
		{"single command", "echo hello", 1},
		{"two commands with &&", "echo foo && echo bar", 2},
		{"three commands with semicolon", "a; b; c", 3},
		{"pipeline", "cat file | grep foo", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps, err := ParseShell(tt.script)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			count := CountCommands(ps)
			if count != tt.expected {
				t.Errorf("expected %d commands, got %d", tt.expected, count)
			}
		})
	}
}
