# godolint

> Same shit, different language

![godolint logo](logo.png)

## godolint? hadolint?

This project is a (loving and) friendly port of hadolint, a venerable and wildly used dockerfile linter
written in haskell.

While there is nothing wrong with hadolint in itself, it does not fit quite well inside the
go devtools ecosytem:
- it is a foreign dependency, that has to be handled outside the familiar go toolchain
- it comes as a binary, that you have to shell out to, and parse its output
- modifying it or extending it requires haskell knowledge (which might be a foreign proposition for some gophers),
and maintaining a fork

godolint is a pure go port of hadolint, filling an age-old gap in the container tooling landscape:
the ability to lint dockerfiles in go.

It comes with an example binary that mimicks hadolint behavior, and can also be used (of course,
and primarily) as a library, integrated into other go projects and devtools.

## Installation

### As a CLI tool

```bash
go install github.com/farcloser/godolint/cmd/godolint@latest
```

### As a library

```bash
go get github.com/farcloser/godolint
```

## Quick Start

### CLI Usage

```bash
# Lint a Dockerfile
godolint Dockerfile

# Output: JSON array of violations
# [
#   {
#     "Code": "DL3007",
#     "Severity": "warning",
#     "Message": "Using latest is prone to errors if the image will ever update",
#     "Line": 1
#   }
# ]
```

### Library Usage

```go
package main

import (
    "fmt"
    "os"

    "github.com/farcloser/godolint/internal/parser"
    "github.com/farcloser/godolint/internal/process"
    "github.com/farcloser/godolint/internal/rule"
    "github.com/farcloser/godolint/internal/rules"
)

func main() {
    // Read Dockerfile
    content, _ := os.ReadFile("Dockerfile")

    // Parse
    p := parser.NewBuildkitParser()
    instructions, _ := p.Parse(content)

    // Configure rules
    myRules := []rule.Rule{
        rules.DL3007(), // No latest tag
        rules.DL3000(), // Absolute WORKDIR
        rules.DL3020(), // Use COPY not ADD
        rules.DL4000(), // No MAINTAINER
        // ... 26 more rules available

        // Optional: Add shellcheck for RUN validation
        // shell.NewShellcheckRule(shell.NewBinaryShellchecker())
    }

    // Lint
    processor := process.NewProcessor(myRules)
    failures := processor.Run(instructions)

    // Report
    for _, f := range failures {
        fmt.Printf("%s (line %d): %s\n", f.Code, f.Line, f.Message)
    }
}
```

## Current Status

### Implementation Progress

- **Total rules**: 65 (matching hadolint)
- **Fully implemented**: 30 rules + ShellCheck integration
- **Test coverage**: 27 test files with 200+ test cases
- **Test pass rate**: 100% (all tests passing)

### Implemented Rules (30)

**Dockerfile best practices:**
- DL3000 - Use absolute WORKDIR
- DL3001 - Pipe into commands not supported
- DL3002 - Last USER should not be root
- DL3003 - Use WORKDIR instead of cd
- DL3004 - Do not use sudo
- DL3006 - Always tag image versions explicitly
- DL3007 - Using latest is prone to errors
- DL3011 - Valid UNIX ports range 0-65535
- DL3012 - Multiple HEALTHCHECK instructions
- DL3014 - Use -y switch for apt-get
- DL3015 - Avoid additional packages with yum
- DL3018 - Pin versions in apk add
- DL3019 - Use --no-cache in apk add
- DL3020 - Use COPY instead of ADD for files
- DL3021 - COPY with more than 2 arguments requires slash
- DL3022 - COPY --from needs valid reference
- DL3023 - COPY --from cannot reference itself
- DL3024 - FROM alias must be unique
- DL3025 - Use JSON notation for CMD/ENTRYPOINT
- DL3027 - Do not use apt
- DL3029 - Do not use --platform in FROM
- DL3043 - ONBUILD requires FROM
- DL3045 - COPY --chmod invalid format
- DL3048 - Invalid DNS option
- DL3059 - Multiple consecutive RUN instructions
- DL3061 - Multiple instructions for same port
- DL4000 - MAINTAINER is deprecated
- DL4003 - Multiple CMD instructions
- DL4004 - Multiple ENTRYPOINT instructions

**Shell validation:**
- DL3008 - Pin versions in apt-get install (requires shellcheck)
- ShellCheck integration - Full support for shell script validation in RUN instructions

### Output Format

godolint outputs JSON arrays of violations:

```json
[
  {
    "Code": "DL3007",
    "Severity": "warning",
    "Message": "Using latest is prone to errors if the image will ever update",
    "Line": 1
  },
  {
    "Code": "DL3000",
    "Severity": "error",
    "Message": "Use absolute WORKDIR",
    "Line": 5
  }
]
```

Exit codes:
- `0`: No violations
- `1`: Violations found

## Architecture

### Parser

godolint uses [moby/buildkit](https://github.com/moby/buildkit) for parsing Dockerfiles, ensuring compatibility with Docker's official parser.

The parser interface is pluggable:

```go
type Parser interface {
    Parse(dockerfile []byte) ([]syntax.InstructionPos, error)
}
```

### Shell Script Validation

godolint includes full shellcheck integration for validating shell commands in RUN instructions:

- **Stateful tracking** - Tracks ENV, ARG, and SHELL instructions across the Dockerfile
- **Multi-stage support** - Correctly resets state on FROM instructions
- **Smart skipping** - Automatically skips non-POSIX shells (PowerShell, cmd)
- **Complete context** - Constructs scripts with proper shebang and environment exports

ShellCheck violations are reported with SC#### codes alongside DL#### rules.

### Rule Engine

Rules implement a stateful interface allowing cross-instruction validation:

```go
type Rule interface {
    Code() RuleCode              // "DL3007"
    Severity() Severity          // Error, Warning, Info, Style
    Message() string             // Human-readable description
    InitialState() State         // Initialize rule state
    Check(line int, state State, instruction Instruction) State
    Finalize(state State) State  // Process final state
}
```

Rules can maintain state across instructions for complex validations (multi-stage builds, shell context tracking, etc.).

### Instruction Types

The AST covers all Dockerfile instructions:
- `FROM` (image, tag, digest, platform, alias)
- `RUN`, `CMD`, `ENTRYPOINT` (commands)
- `COPY`, `ADD` (sources, destination, flags)
- `ENV`, `ARG`, `LABEL` (key-value pairs)
- `WORKDIR`, `USER`, `EXPOSE`, `VOLUME`
- `HEALTHCHECK`, `STOPSIGNAL`, `SHELL`
- `MAINTAINER`, `ONBUILD`

## Rules & Tests

Rules and tests are (partly) generated from the hadolint codebase.

As such, godolint passes hadolint's internal test suite, (hopefully) guaranteeing
a seamless transition.

Further, as we do not plan on fragmenting the ecosystem for no good reason, we use exactly the same rules,
and plan on keeping up with hadolint's updated/new rules.

### Code Generation

Rule stubs and tests are auto-generated from hadolint's Haskell source:

```bash
# Generate rule stubs and tests (runs both generators)
go generate ./internal/rules
```

This single command:
- Extracts rule metadata (code, severity, message) from Haskell source
- Generates Go function stubs for unimplemented rules
- Converts Haskell test cases to Go tests for implemented rules
- Reports implementation status

## Development

### Build

```bash
# Build CLI binary
make build
# Output: ./bin/godolint

# Or directly
go build -o godolint ./cmd/godolint
```

### Test

```bash
# Run all tests
make test

# Run specific rule tests
go test ./internal/rules -run DL3007
```

### Lint

```bash
# Install dev tools
make install-dev-tools

# Run linters
make lint
```

### Adding New Rules

1. Check if rule stub exists in `internal/rules/dlXXXX.go`
2. If not, run `go generate ./internal/rules` to generate stub
3. Implement the `checkDLXXXX` function
4. Add rule to `cmd/godolint/main.go` registration
5. Run tests: `go test ./internal/rules -run DLXXXX`

Example:

```go
// internal/rules/dl3007.go
func checkDL3007(instruction syntax.Instruction) bool {
    from, ok := instruction.(*syntax.From)
    if !ok {
        return true // Not a FROM instruction, passes
    }

    // Fail if using :latest tag
    if from.Image.Tag != nil && *from.Image.Tag == "latest" {
        return false
    }

    return true
}
```

## Roadmap

### Completed ✓
- [x] Stateful rule processor for cross-instruction validation
- [x] ShellCheck integration for shell script validation
- [x] 30 hadolint rules fully implemented
- [x] 100% test pass rate

### Short-term
- [ ] Add shellcheck rule to default CLI rule set
- [ ] Implement remaining shell-dependent rules (DL3009, DL3013, DL3016, etc.)
- [ ] Add CLI flags (verbosity, rule selection, output format)
- [ ] Configuration file support (YAML/TOML)
- [ ] Pragma support (`# godolint ignore=DL3007`)

### Medium-term
- [ ] Alternative output formats (SARIF, human-readable TTY)
- [ ] Performance optimization and benchmarking
- [ ] Documentation improvements
- [ ] GitHub Actions integration example

### Long-term
- [ ] Plugin architecture for custom rules
- [ ] Integration with container build tools
- [ ] Language-agnostic rule repository
- [ ] IDE extensions (VS Code, GoLand)

## Limitations

### Current
- **No configuration file support** (hardcoded rules)
- **No pragma/inline ignore directives**
- **JSON output only**
- **No CLI flags** (except file path)

### Parser Differences
- Quote handling differs slightly from hadolint (buildkit includes quotes in values)
- Some escape sequences handled differently
- Affects ~30% of auto-generated test cases

### Rule Coverage
- **30 rules** fully implemented and tested
- **35 rules** remaining (mix of simple rules and shell-dependent)
- **ShellCheck integration** complete and ready for shell-dependent rules

## Contributing

Contributions welcome! Current priorities:
- Implementing remaining rules using shellcheck integration
- Adding CLI flags and configuration file support
- Alternative output formats (SARIF, GitHub Actions annotations)
- Performance optimization
- Documentation and examples

See `SHELLCHECK_IMPLEMENTATION.md` for details on the shell validation architecture.

## License

See LICENSE file.

