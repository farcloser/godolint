package pragma_test

import (
	"testing"

	"github.com/forkcloser/godolint/internal/pragma"
)

func TestParseShell(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want string // "": not a shell pragma
	}{
		{"plain", "hadolint shell=powershell", "powershell"},
		{"spaces", "  hadolint \t shell = powershell.exe  ", "powershell.exe"},
		{"inline comment", "hadolint shell=powershell # foobar", "powershell"},
		{"a path", "hadolint shell=/bin/bash", "/bin/bash"},
		{"trailing text is not a comment", "hadolint shell=powershell foo", ""},
		{"empty name", "hadolint shell=", ""},
		{"no space after hadolint", "hadolintshell=powershell", ""},
		{"an ignore pragma", "hadolint ignore=DL3002", ""},
		{"a comment", "shell=powershell", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := pragma.ParseShell(tt.text)
			if ok != (tt.want != "") || got != tt.want {
				t.Fatalf("ParseShell(%q) = %q, %v; want %q, %v", tt.text, got, ok, tt.want, tt.want != "")
			}
		})
	}
}
