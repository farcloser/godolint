// Package shell provides shell script parsing for analyzing RUN instructions.
package shell

import (
	"fmt"
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

// CmdPart represents a part of a command (argument or flag).
// Ported from Hadolint.Shell.CmdPart.
type CmdPart struct {
	Arg string // The argument text
	ID  int    // Position/ID for tracking which args belong to which flags
}

// Command represents a parsed shell command.
// Ported from Hadolint.Shell.Command.
type Command struct {
	Name      string    // Command name (e.g., "apt-get")
	Arguments []CmdPart // All arguments including flags
	Flags     []CmdPart // Extracted flags only
}

// ParsedShell represents a parsed shell script.
// Ported from Hadolint.Shell.ParsedShell.
type ParsedShell struct {
	Original        string    // Original script text
	PresentCommands []Command // Extracted commands
}

// ParseShell parses a shell script and extracts commands.
// Ported from Hadolint.Shell.parseShell.
func ParseShell(script string) (*ParsedShell, error) {
	// Add shebang to help parser
	fullScript := "#!/bin/bash\n" + script

	// Parse the script
	r := strings.NewReader(fullScript)
	file, err := syntax.NewParser().Parse(r, "")
	if err != nil {
		return nil, fmt.Errorf("failed to parse shell script: %w", err)
	}

	// Extract commands
	commands := extractCommands(file)

	return &ParsedShell{
		Original:        script,
		PresentCommands: commands,
	}, nil
}

// extractCommands walks the AST and extracts all commands.
func extractCommands(file *syntax.File) []Command {
	var commands []Command

	syntax.Walk(file, func(node syntax.Node) bool {
		// Look for CallExpr nodes (simple commands)
		if call, ok := node.(*syntax.CallExpr); ok {
			if cmd := extractCommand(call); cmd != nil {
				commands = append(commands, *cmd)
			}
		}
		return true
	})

	return commands
}

// extractCommand extracts a Command from a CallExpr node.
func extractCommand(call *syntax.CallExpr) *Command {
	if len(call.Args) == 0 {
		return nil
	}

	// Get command name from first argument
	nameWord := call.Args[0]
	name := wordToString(nameWord)
	if name == "" {
		return nil
	}

	// Extract all arguments with IDs
	var allArgs []CmdPart
	for i, arg := range call.Args[1:] {
		argStr := wordToString(arg)
		allArgs = append(allArgs, CmdPart{
			Arg: argStr,
			ID:  i,
		})
	}

	// Extract flags from arguments
	flags := extractFlags(allArgs)

	return &Command{
		Name:      name,
		Arguments: allArgs,
		Flags:     flags,
	}
}

// wordToString converts a Word node to a string.
// Similar to ShellCheck.ASTLib.oversimplify.
func wordToString(word *syntax.Word) string {
	var sb strings.Builder

	for _, part := range word.Parts {
		switch p := part.(type) {
		case *syntax.Lit:
			sb.WriteString(p.Value)
		case *syntax.SglQuoted:
			sb.WriteString(p.Value)
		case *syntax.DblQuoted:
			// Recursively process quoted parts
			for _, qp := range p.Parts {
				if lit, ok := qp.(*syntax.Lit); ok {
					sb.WriteString(lit.Value)
				} else {
					// Variables, expansions, etc. - simplified as ${VAR}
					sb.WriteString("${VAR}")
				}
			}
		case *syntax.ParamExp:
			sb.WriteString("${VAR}")
		case *syntax.CmdSubst:
			sb.WriteString("${VAR}")
		case *syntax.ArithmExp:
			sb.WriteString("${VAR}")
		default:
			// Other expansions simplified
			sb.WriteString("${VAR}")
		}
	}

	return sb.String()
}

// extractFlags extracts flag arguments from a list of arguments.
// Ported from Hadolint.Shell.getAllFlags.
func extractFlags(args []CmdPart) []CmdPart {
	var flags []CmdPart

	for _, arg := range args {
		// Skip special cases
		if arg.Arg == "--" || arg.Arg == "-" {
			continue
		}

		// Long flags: --flag or --flag=value
		if strings.HasPrefix(arg.Arg, "--") {
			flagName := strings.TrimPrefix(arg.Arg, "--")
			// Remove =value part if present
			if idx := strings.IndexByte(flagName, '='); idx != -1 {
				flagName = flagName[:idx]
			}
			flags = append(flags, CmdPart{
				Arg: flagName,
				ID:  arg.ID,
			})
			continue
		}

		// Short flags: -abc becomes three flags: a, b, c
		if strings.HasPrefix(arg.Arg, "-") {
			flagChars := strings.TrimPrefix(arg.Arg, "-")
			for _, ch := range flagChars {
				flags = append(flags, CmdPart{
					Arg: string(ch),
					ID:  arg.ID,
				})
			}
		}
	}

	return flags
}

// FindCommandNames returns all command names in the parsed shell.
// Ported from Hadolint.Shell.findCommandNames.
func FindCommandNames(ps *ParsedShell) []string {
	names := make([]string, len(ps.PresentCommands))
	for i, cmd := range ps.PresentCommands {
		names[i] = cmd.Name
	}
	return names
}

// CmdHasArgs checks if a command has a specific name and contains specific arguments.
// Ported from Hadolint.Shell.cmdHasArgs.
func CmdHasArgs(expectedName string, expectedArgs []string, cmd Command) bool {
	if cmd.Name != expectedName {
		return false
	}

	// Check if any of the expected args are present
	for _, expected := range expectedArgs {
		for _, arg := range cmd.Arguments {
			if arg.Arg == expected {
				return true
			}
		}
	}

	return false
}

// HasFlag checks if a command has a specific flag.
// Ported from Hadolint.Shell.hasFlag.
func HasFlag(flag string, cmd Command) bool {
	for _, f := range cmd.Flags {
		if f.Arg == flag {
			return true
		}
	}
	return false
}

// GetArgs returns all argument strings (including flags).
// Ported from Hadolint.Shell.getArgs.
func GetArgs(cmd Command) []string {
	args := make([]string, len(cmd.Arguments))
	for i, arg := range cmd.Arguments {
		args[i] = arg.Arg
	}
	return args
}

// GetArgsNoFlags returns arguments that are not flags or flag values.
// Ported from Hadolint.Shell.getArgsNoFlags.
func GetArgsNoFlags(cmd Command) []string {
	// Get IDs of all flags
	flagIDs := make(map[int]bool)
	for _, flag := range cmd.Flags {
		flagIDs[flag.ID] = true
	}

	// Return arguments whose IDs are not in flag IDs
	var result []string
	for _, arg := range cmd.Arguments {
		if !flagIDs[arg.ID] {
			result = append(result, arg.Arg)
		}
	}

	return result
}

// CountCommands returns the number of commands in a parsed shell script.
// Used by DL3059 to detect chained commands.
func CountCommands(ps *ParsedShell) int {
	return len(ps.PresentCommands)
}

// HasAnyFlag checks if a command has any of the specified flags.
// Ported from Hadolint.Shell.hasAnyFlag.
func HasAnyFlag(flags []string, cmd Command) bool {
	for _, flag := range flags {
		if HasFlag(flag, cmd) {
			return true
		}
	}
	return false
}

// CountFlag counts how many times a specific flag appears.
// Ported from Hadolint.Shell.countFlag.
func CountFlag(flag string, cmd Command) int {
	count := 0
	for _, f := range cmd.Flags {
		if f.Arg == flag {
			count++
		}
	}
	return count
}

// HasArg checks if a command has a specific argument.
// Ported from Hadolint.Shell.hasArg.
func HasArg(arg string, cmd Command) bool {
	for _, a := range cmd.Arguments {
		if a.Arg == arg {
			return true
		}
	}
	return false
}

// UsingProgram checks if any command uses a specific program.
// Ported from Hadolint.Shell.usingProgram.
func UsingProgram(program string, ps *ParsedShell) bool {
	for _, cmd := range ps.PresentCommands {
		if cmd.Name == program {
			return true
		}
	}
	return false
}
