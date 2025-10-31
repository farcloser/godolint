// Package syntax defines the AST (Abstract Syntax Tree) for Dockerfiles.
// Ported from Language.Docker.Syntax
package syntax

// Ported from Instruction in Language.Docker.Syntax.
type Instruction interface {
	// Name returns the instruction name (FROM, RUN, COPY, etc.)
	Name() string
}

// Ported from InstructionPos in Language.Docker.Syntax.
type InstructionPos struct {
	Instruction Instruction
	LineNumber  int
}

// Field names match Haskell: image, tag, digest, alias, platform (all lowercase).
type BaseImage struct {
	Image    string  // Image name (e.g., "debian", "alpine")
	Tag      *string // Optional tag (e.g., "latest", "3.20")
	Digest   *string // Optional digest (e.g., "sha256:...")
	Alias    *string // Optional alias for multi-stage builds (e.g., "AS builder")
	Platform *string // Optional platform (e.g., "linux/amd64")
}

// Ported from From in Language.Docker.Syntax.
type From struct {
	Image BaseImage
}

func (f *From) Name() string {
	return "FROM"
}

// Ported from Run in Language.Docker.Syntax.
type Run struct {
	Command string   // The shell command to execute
	Flags   []string // RUN instruction flags (e.g., --mount)
}

func (r *Run) Name() string {
	return "RUN"
}

// Ported from Copy in Language.Docker.Syntax.
type Copy struct {
	Source      []string // Source paths
	Destination string   // Destination path
	From        *string  // Optional --from flag for multi-stage
}

func (c *Copy) Name() string {
	return "COPY"
}

// Ported from Add in Language.Docker.Syntax.
type Add struct {
	Source      []string // Source paths or URLs
	Destination string   // Destination path
}

func (a *Add) Name() string {
	return "ADD"
}

// Ported from Env in Language.Docker.Syntax.
type Env struct {
	Pairs []EnvPair // Environment variable key-value pairs
}

type EnvPair struct {
	Key   string
	Value string
}

func (e *Env) Name() string {
	return "ENV"
}

// Ported from Label in Language.Docker.Syntax.
type Label struct {
	Pairs []LabelPair // Label key-value pairs
}

type LabelPair struct {
	Key   string
	Value string
}

func (l *Label) Name() string {
	return "LABEL"
}

// Ported from Workdir in Language.Docker.Syntax.
type Workdir struct {
	Directory string
}

func (w *Workdir) Name() string {
	return "WORKDIR"
}

// Ported from User in Language.Docker.Syntax.
type User struct {
	User string
}

func (u *User) Name() string {
	return "USER"
}

// Ported from Expose in Language.Docker.Syntax.
type Expose struct {
	Ports []string // Port specifications
}

func (e *Expose) Name() string {
	return "EXPOSE"
}

// Ported from Volume in Language.Docker.Syntax.
type Volume struct {
	Volumes []string // Volume mount points
}

func (v *Volume) Name() string {
	return "VOLUME"
}

// Ported from Cmd in Language.Docker.Syntax.
type Cmd struct {
	Arguments []string // Command arguments
	IsJSON    bool     // true if using JSON/exec form, false if shell form
}

func (c *Cmd) Name() string {
	return "CMD"
}

// Ported from Entrypoint in Language.Docker.Syntax.
type Entrypoint struct {
	Arguments []string // Entrypoint arguments
	IsJSON    bool     // true if using JSON/exec form, false if shell form
}

func (e *Entrypoint) Name() string {
	return "ENTRYPOINT"
}

// Ported from Healthcheck in Language.Docker.Syntax.
type Healthcheck struct {
	Command string // Health check command
}

func (h *Healthcheck) Name() string {
	return "HEALTHCHECK"
}

// Haskell: Maintainer !Text (single unnamed field).
type Maintainer struct {
	MaintainerName string // Maintainer name/email
}

func (m *Maintainer) Name() string {
	return "MAINTAINER"
}

// Haskell: Arg !Text !(Maybe Text) (two unnamed fields: name and optional default value).
type Arg struct {
	ArgName string
	Value   *string // Optional default value
}

func (a *Arg) Name() string {
	return "ARG"
}

// Ported from StopSignal in Language.Docker.Syntax.
type StopSignal struct {
	Signal string
}

func (s *StopSignal) Name() string {
	return "STOPSIGNAL"
}

// Ported from Shell in Language.Docker.Syntax.
type Shell struct {
	Arguments []string
}

func (s *Shell) Name() string {
	return "SHELL"
}

// Ported from OnBuild in Language.Docker.Syntax.
type OnBuild struct {
	Inner Instruction // The wrapped instruction
}

func (o *OnBuild) Name() string {
	return "ONBUILD"
}
