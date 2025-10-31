![godolint logo](logo.png)

# godolint

> Same shit, different language

## godolint? hadolint?

This project is a friendly port of [hadolint](https://github.com/hadolint/hadolint) (a venerable and wildly used dockerfile linter
written in haskell).

While there is nothing wrong with hadolint in itself, it does not fit well inside the
go devtools ecosystem (*):
- it is a foreign dependency, that has to be handled outside the familiar go toolchain
- it comes as a binary, that you have to shell out to, and parse the output of
- modifying it requires haskell knowledge

godolint is a pure go port of hadolint, filling an age-old gap in the container tooling landscape:
the ability to lint dockerfiles in go.

It comes with an example binary that does mimic a small subset of the hadolint binary behavior,
and can also be used (of course, and primarily) as a library in golang projects.

(*) see [Why?](#why) section below for more details

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

### SDK Usage

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/farcloser/godolint/sdk"
)

func main() {
    // Read Dockerfile
    content, _ := os.ReadFile("Dockerfile")

    // Create linter with default configuration (all rules)
    linter := sdk.New()

    // Lint the Dockerfile
    result, err := linter.Lint(context.Background(), content)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Report violations
    for _, v := range result.Violations {
        fmt.Printf("[%s] %s (line %d): %s\n",
            v.Severity, v.Code, v.Line, v.Message)
    }

    // Exit with error code if violations found
    if !result.Passed {
        os.Exit(1)
    }
}
```

### Advanced SDK Usage

```go
// Use a specific rule set
linter := sdk.New(sdk.WithRuleSet(sdk.RuleSetRecommended))

// Disable specific rules
linter := sdk.New(sdk.WithDisabledRules("DL3007", "DL3008"))

// Enable shellcheck integration
linter := sdk.New(sdk.WithShellcheck())

// Check for specific severity levels
if result.HasErrors() {
    fmt.Println("Critical issues found!")
}

counts := result.CountBySeverity()
fmt.Printf("Errors: %d, Warnings: %d\n",
    counts[sdk.SeverityError],
    counts[sdk.SeverityWarning])
```

## Why?

> What is wrong with depending on (hadolint) binaries?

### TL;DR

> copy-paste into a powerpoint for your CISO (don't forget to add ✅ and ❌)

| Aspect | External Binary | Pure Go |
|--------|----------------|---------|
| **Supply Chain Risk** | HIGH - Unverified provenance | LOW - Checksummed source |
| **Compromise Scenarios** | Account takeover, build pipeline, MITM | Requires both GitHub + sum.golang.org |
| **Cryptographic Verification** | Manual, often missing | Automatic via go.sum |
| **Reproducibility** | Poor - binary availability not guaranteed | Excellent - source-based |
| **Multi-arch Support** | Complex per-platform logic | `GOARCH=arm64 go build` |
| **Installation Complexity** | Bash scripts, sudo, platform detection | `go install` |
| **CI/CD Integration** | Manual download + caching | Native Go tooling |
| **Container Image Size** | +50-100MB (deps + binary) | +5-10MB (static binary) |
| **SBOM Generation** | Manual, incomplete | `go list -m all` |
| **Vulnerability Scanning** | Opaque binary, hard to scan | `govulncheck ./...` |
| **License Compliance** | Manual transitive analysis | `go-licenses report` |
| **TCB Size** | Large (GitHub, CDN, toolchain) | Small (Go compiler, go.sum) |
| **Offline Builds** | Impossible | `go mod vendor` |
| **Dependency Updates** | Manual version bumps | `go get -u` + Dependabot |

### With more words

Supply chain concerns:
* binary provenance is unverifiable
* no cryptographic verification
* version pinning is fragile
* SBOM is incomplete or missing
* license compliance is opaque / hardcoded
* trusted computing base is WIDE
* scanning binaries is difficult

Operational complexity:
* multi-architecture is hell
* installation script is hell with fries
* CI/CD compounds hell further into full-blown hell
* availability is not guaranteed
* os-level dependencies if not compiled statically (which you have to figure out)

You COULD (should!) instead build from source (that is, if you like Gentoo).
At least you would get proper git commit pinning, scannable / auditable source code, SBOM.
But then, you have to maintain a build environment for it, in a language and with tools
that does not match the rest of your stack, and that may also raise additional supply
chain concerns, time and dedicated expertise.

So, yes, there is nothing wrong with hadolint-the-project-and-tool, per-se.
But if what you do is more than noodling around, if you are not a haskell shop,
or if you take your security seriously, depending on external binaries is definitely not
something you should do.

This is why we wrote godolint, for people doing go.

## Relationship with hadolint

godolint is a faithful port of hadolint to Go. Rule declarations and tests are automatically generated from hadolint's Haskell source.

**Current status:**
- All unit tests are automatically converted to Go tests
- Rule metadata (code, severity, message) matches hadolint exactly
- Ongoing work to implement remaining rules

**Project philosophy:**
- hadolint is the upstream for rule definitions and baseline tests
- We have no intention of fragmenting the community
- New **core** rule requests should go through hadolint
- We will regularly sync with hadolint's evolution

## SDK API

### Core Types

```go
// Linter performs Dockerfile linting
type Linter struct { ... }

// Result contains linting results
type Result struct {
    Violations []Violation
    Passed     bool
}

// Violation represents a single rule violation
type Violation struct {
    Code     string   // Rule code (e.g., "DL3000")
    Severity Severity // error, warning, info, style
    Message  string   // Human-readable description
    Line     int      // Line number (1-indexed)
}
```

### Configuration Options

```go
// WithRuleSet - Use a predefined rule set
sdk.New(sdk.WithRuleSet(sdk.RuleSetRecommended))

// WithDisabledRules - Disable specific rules
sdk.New(sdk.WithDisabledRules("DL3007", "DL3008"))

// WithShellcheck - Enable shellcheck integration
sdk.New(sdk.WithShellcheck())
```

### Rule Sets

- `RuleSetAll` - All implemented rules (default)
- `RuleSetRecommended` - Only Error and Warning severity rules
- `RuleSetStrict` - Same as All (for compatibility)

### Result Methods

```go
result.HasErrors()           // Check for error-severity violations
result.HasWarnings()         // Check for warning-severity violations
result.CountBySeverity()     // Get violation counts by severity
```

### Typed Errors

```go
*sdk.ParseError    // Dockerfile parsing failed
*sdk.RuleError     // Rule execution failed
```

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

Rule stubs and tests are auto-generated from hadolint's source:

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

### Short-term
- [ ] Implement remaining hadolint rules
- [ ] Add CLI flags (verbosity, rule selection, output format)
- [ ] Pragma support (`# godolint ignore=DL3007`)
- [ ] Configuration file support (YAML/TOML)

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
- **No pragma/inline ignore directives** - Coming soon
- **CLI limited** - JSON output only, minimal flags (SDK has full flexibility)
- **No configuration file support** - Use SDK for programmatic configuration

### Parser Differences
- Quote handling differs slightly from hadolint (buildkit includes quotes in values)
- Some escape sequences handled differently
- Affects ~30% of auto-generated test cases

## Contributing

Contributions welcome! Current priorities:
- Implementing remaining hadolint rules
- Adding CLI flags and pragma support
- Alternative output formats (SARIF, GitHub Actions annotations)
- Performance optimization
- Documentation and examples

See `SHELLCHECK_IMPLEMENTATION.md` for details on the shell validation architecture.

## License

See LICENSE file.

