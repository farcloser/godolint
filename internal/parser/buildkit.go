package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/farcloser/godolint/internal/syntax"
)

// BuildkitParser implements Parser using moby/buildkit's Dockerfile parser.
type BuildkitParser struct{}

// NewBuildkitParser creates a new buildkit-based parser.
func NewBuildkitParser() *BuildkitParser {
	return &BuildkitParser{}
}

// Parse parses a Dockerfile using moby/buildkit and converts to our AST format.
func (p *BuildkitParser) Parse(dockerfile []byte) ([]syntax.InstructionPos, error) {
	// Parse using buildkit
	result, err := parser.Parse(bytes.NewReader(dockerfile))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Dockerfile: %w", err)
	}

	// Convert buildkit AST to our AST format
	var instructions []syntax.InstructionPos

	for _, child := range result.AST.Children {
		instr, err := convertNode(child)
		if err != nil {
			// Skip instructions we can't convert (continue linting rest)
			continue
		}

		if instr != nil {
			instructions = append(instructions, syntax.InstructionPos{
				Instruction: instr,
				LineNumber:  child.StartLine,
			})
		}
	}

	return instructions, nil
}

// convertNode converts a buildkit AST node to our Instruction type.
func convertNode(node *parser.Node) (syntax.Instruction, error) {
	switch strings.ToLower(node.Value) {
	case "from":
		return convertFrom(node)
	case "run":
		return convertRun(node)
	case "copy":
		return convertCopy(node)
	case "add":
		return convertAdd(node)
	case "env":
		return convertEnv(node)
	case "label":
		return convertLabel(node)
	case "workdir":
		return convertWorkdir(node)
	case "user":
		return convertUser(node)
	case "expose":
		return convertExpose(node)
	case "volume":
		return convertVolume(node)
	case "cmd":
		return convertCmd(node)
	case "entrypoint":
		return convertEntrypoint(node)
	case "healthcheck":
		return convertHealthcheck(node)
	case "maintainer":
		return convertMaintainer(node)
	case "arg":
		return convertArg(node)
	case "stopsignal":
		return convertStopSignal(node)
	case "shell":
		return convertShell(node)
	case "onbuild":
		return convertOnBuild(node)
	default:
		return nil, fmt.Errorf("unknown instruction: %s", node.Value)
	}
}

// Helper to get next child value.
func nextValue(node *parser.Node) string {
	if node.Next != nil {
		return node.Next.Value
	}

	return ""
}

// Helper to collect all remaining values.
func collectValues(node *parser.Node) []string {
	var values []string

	for n := node.Next; n != nil; n = n.Next {
		values = append(values, n.Value)
	}

	return values
}

func convertFrom(node *parser.Node) (*syntax.From, error) {
	if node.Next == nil {
		return nil, errors.New("FROM missing image")
	}

	// Collect all values to handle: image:tag@digest AS alias
	values := collectValues(node)
	if len(values) == 0 {
		return nil, errors.New("FROM missing image")
	}

	baseImage := syntax.BaseImage{
		Image: values[0],
	}

	// Extract --platform flag if present
	for _, flag := range node.Flags {
		if strings.HasPrefix(flag, "--platform=") {
			platform := strings.TrimPrefix(flag, "--platform=")
			baseImage.Platform = &platform
		}
	}

	// Check for AS alias (last token after AS keyword)
	for i := 1; i < len(values); i++ {
		if values[i] == "AS" || values[i] == "as" {
			if i+1 < len(values) {
				alias := values[i+1]
				baseImage.Alias = &alias
			}

			break
		}
	}

	// Parse image reference (image:tag@digest)
	// The buildkit parser gives us the raw string, we need to parse it
	imagePart := values[0]

	// Check for digest (@sha256:...)
	if atIdx := indexOf(imagePart, '@'); atIdx != -1 {
		digest := imagePart[atIdx+1:]
		baseImage.Digest = &digest
		imagePart = imagePart[:atIdx]
	}

	// Check for tag (:tag)
	if colonIdx := indexOf(imagePart, ':'); colonIdx != -1 {
		tag := imagePart[colonIdx+1:]
		baseImage.Tag = &tag
		baseImage.Image = imagePart[:colonIdx]
	}

	return &syntax.From{
		Image: baseImage,
	}, nil
}

// Helper to find index of character in string.
func indexOf(s string, ch rune) int {
	for i, c := range s {
		if c == ch {
			return i
		}
	}

	return -1
}

func convertRun(node *parser.Node) (*syntax.Run, error) {
	// Collect all parts of the RUN command
	command := ""

	for n := node.Next; n != nil; n = n.Next {
		if command != "" {
			command += " "
		}

		command += n.Value
	}

	return &syntax.Run{
		Command: command,
		Flags:   node.Flags,
	}, nil
}

func convertCopy(node *parser.Node) (*syntax.Copy, error) {
	values := collectValues(node)
	if len(values) < 2 {
		return nil, errors.New("COPY requires at least source and destination")
	}

	copy := &syntax.Copy{
		Source:      values[:len(values)-1],
		Destination: values[len(values)-1],
	}

	// Check for --from flag
	// Buildkit parser includes flags in node.Flags
	for _, flag := range node.Flags {
		if strings.HasPrefix(flag, "--from=") {
			fromValue := strings.TrimPrefix(flag, "--from=")
			copy.From = &fromValue
		}
	}

	return copy, nil
}

func convertAdd(node *parser.Node) (*syntax.Add, error) {
	values := collectValues(node)
	if len(values) < 2 {
		return nil, errors.New("ADD requires at least source and destination")
	}

	return &syntax.Add{
		Source:      values[:len(values)-1],
		Destination: values[len(values)-1],
	}, nil
}

func convertEnv(node *parser.Node) (*syntax.Env, error) {
	values := collectValues(node)

	var pairs []syntax.EnvPair

	// ENV can be "ENV key=value" or "ENV key value"
	for i := 0; i < len(values); i++ {
		if i+1 < len(values) {
			pairs = append(pairs, syntax.EnvPair{
				Key:   values[i],
				Value: values[i+1],
			})
			i++ // Skip next value
		}
	}

	return &syntax.Env{
		Pairs: pairs,
	}, nil
}

func convertLabel(node *parser.Node) (*syntax.Label, error) {
	values := collectValues(node)

	var pairs []syntax.LabelPair

	for i := 0; i < len(values); i++ {
		if i+1 < len(values) {
			pairs = append(pairs, syntax.LabelPair{
				Key:   values[i],
				Value: values[i+1],
			})
			i++
		}
	}

	return &syntax.Label{
		Pairs: pairs,
	}, nil
}

func convertWorkdir(node *parser.Node) (*syntax.Workdir, error) {
	return &syntax.Workdir{
		Directory: nextValue(node),
	}, nil
}

func convertUser(node *parser.Node) (*syntax.User, error) {
	return &syntax.User{
		User: nextValue(node),
	}, nil
}

func convertExpose(node *parser.Node) (*syntax.Expose, error) {
	return &syntax.Expose{
		Ports: collectValues(node),
	}, nil
}

func convertVolume(node *parser.Node) (*syntax.Volume, error) {
	return &syntax.Volume{
		Volumes: collectValues(node),
	}, nil
}

func convertCmd(node *parser.Node) (*syntax.Cmd, error) {
	// Check if using JSON notation (exec form)
	isJSON := node.Attributes != nil && node.Attributes["json"]

	return &syntax.Cmd{
		Arguments: collectValues(node),
		IsJSON:    isJSON,
	}, nil
}

func convertEntrypoint(node *parser.Node) (*syntax.Entrypoint, error) {
	// Check if using JSON notation (exec form)
	isJSON := node.Attributes != nil && node.Attributes["json"]

	return &syntax.Entrypoint{
		Arguments: collectValues(node),
		IsJSON:    isJSON,
	}, nil
}

func convertHealthcheck(node *parser.Node) (*syntax.Healthcheck, error) {
	command := ""

	for n := node.Next; n != nil; n = n.Next {
		if command != "" {
			command += " "
		}

		command += n.Value
	}

	return &syntax.Healthcheck{
		Command: command,
	}, nil
}

func convertMaintainer(node *parser.Node) (*syntax.Maintainer, error) {
	return &syntax.Maintainer{
		MaintainerName: nextValue(node),
	}, nil
}

func convertArg(node *parser.Node) (*syntax.Arg, error) {
	value := nextValue(node)
	if value == "" {
		return nil, errors.New("ARG missing name")
	}

	// Parse ARG name or ARG name=value
	arg := &syntax.Arg{}

	// Check for = separator
	if eqIdx := indexOf(value, '='); eqIdx != -1 {
		arg.ArgName = value[:eqIdx]
		defaultVal := value[eqIdx+1:]
		arg.Value = &defaultVal
	} else {
		arg.ArgName = value
	}

	return arg, nil
}

func convertStopSignal(node *parser.Node) (*syntax.StopSignal, error) {
	return &syntax.StopSignal{
		Signal: nextValue(node),
	}, nil
}

func convertShell(node *parser.Node) (*syntax.Shell, error) {
	return &syntax.Shell{
		Arguments: collectValues(node),
	}, nil
}

func convertOnBuild(node *parser.Node) (*syntax.OnBuild, error) {
	// ONBUILD wraps another instruction
	// buildkit stores the full instruction in Original field (e.g., "ONBUILD FROM debian")
	// We need to parse the inner instruction from the Original string
	original := strings.TrimSpace(node.Original)

	// Remove "ONBUILD " prefix
	if !strings.HasPrefix(strings.ToUpper(original), "ONBUILD ") {
		return nil, errors.New("invalid ONBUILD instruction")
	}

	innerText := strings.TrimSpace(original[8:]) // len("ONBUILD ") = 8
	if innerText == "" {
		return nil, errors.New("ONBUILD missing instruction")
	}

	// Parse the inner instruction as a complete Dockerfile
	innerResult, err := parser.Parse(bytes.NewReader([]byte(innerText)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ONBUILD inner instruction: %w", err)
	}

	if len(innerResult.AST.Children) == 0 {
		return nil, errors.New("ONBUILD has no inner instruction")
	}

	// Convert the first (and only) child instruction
	inner, err := convertNode(innerResult.AST.Children[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert ONBUILD inner instruction: %w", err)
	}

	return &syntax.OnBuild{
		Inner: inner,
	}, nil
}
