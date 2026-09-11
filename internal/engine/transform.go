package engine

import (
	"context"
	"sort"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/pkg/types"
)

// Target is a transform output representation.
type Target string

// Every supported transform target.
const (
	TargetMarkdown Target = "markdown"
	TargetXML      Target = "xml"
	TargetJSON     Target = "json"
	TargetSystem   Target = "system"
	TargetTemplate Target = "template"
	TargetAgent    Target = "agent"
)

// Targets lists every supported target, for help text and validation.
func Targets() []string {
	return []string{"markdown", "xml", "json", "system", "template", "agent"}
}

// ValidTarget reports whether s names a supported target.
func ValidTarget(s string) bool {
	for _, t := range Targets() {
		if t == s {
			return true
		}
	}
	return false
}

// TransformInput configures a transform run.
type TransformInput struct {
	Prompt string
	Target Target
}

type transformRaw struct {
	Transformed string   `json:"transformed"`
	Notes       []string `json:"notes"`
	Variables   []string `json:"variables"`
	Warnings    []string `json:"warnings"`
}

// Transform converts a prompt into a different structure or representation
// while preserving what it asks for.
func (e *Engine) Transform(ctx context.Context, in TransformInput) (*types.TransformResult, error) {
	if strings.TrimSpace(in.Prompt) == "" {
		return nil, apperr.Usagef("transform needs a prompt (argument, file, or stdin)")
	}
	if !ValidTarget(string(in.Target)) {
		return nil, apperr.Usagef("--to must be one of: %s", strings.Join(Targets(), ", "))
	}

	user := "PROMPT TO TRANSFORM:\n\n" + in.Prompt + "\n\nTARGET REPRESENTATION: " + string(in.Target)

	var raw transformRaw
	if err := e.call(ctx, "transform", user, &raw); err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.Transformed) == "" {
		return nil, apperr.New("The transformer returned nothing.",
			"Retry the command.").WithCode(apperr.CodeMalformed)
	}

	res := &types.TransformResult{
		Operation: "transform",
		Target:    string(in.Target),
		Result:    strings.TrimSpace(raw.Transformed),
		Notes:     cleanList(raw.Notes),
		Tokens:    e.tokens(in.Prompt, raw.Transformed),
		Warnings:  cleanList(raw.Warnings),
	}

	if in.Target == TargetTemplate {
		res.Variables = dedupeVars(raw.Variables, res.Result)
		if len(res.Variables) == 0 {
			res.Warnings = append(res.Warnings,
				"no variables were introduced; the prompt had no values worth parameterizing")
		}
	}
	return res, nil
}

// dedupeVars merges the model-reported variables with any {{name}} tokens
// actually present in the output, returning a sorted unique list.
func dedupeVars(reported []string, body string) []string {
	set := map[string]struct{}{}
	for _, v := range reported {
		v = strings.Trim(strings.TrimSpace(v), "{}")
		if v != "" {
			set[v] = struct{}{}
		}
	}
	for _, v := range scanTemplateVars(body) {
		set[v] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for v := range set {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

// scanTemplateVars finds {{ name }} tokens in s.
func scanTemplateVars(s string) []string {
	var out []string
	for {
		i := strings.Index(s, "{{")
		if i < 0 {
			break
		}
		j := strings.Index(s[i:], "}}")
		if j < 0 {
			break
		}
		name := strings.TrimSpace(s[i+2 : i+j])
		if name != "" {
			out = append(out, name)
		}
		s = s[i+j+2:]
	}
	return out
}
