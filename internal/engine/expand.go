package engine

import (
	"context"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/pkg/types"
)

// Depth controls how far expand goes.
type Depth string

const (
	DepthConcise    Depth = "concise"
	DepthDetailed   Depth = "detailed"
	DepthProduction Depth = "production"
)

// ValidDepth reports whether s is a known depth.
func ValidDepth(s string) bool {
	switch Depth(s) {
	case DepthConcise, DepthDetailed, DepthProduction:
		return true
	}
	return false
}

// ExpandInput configures an expand run.
type ExpandInput struct {
	Prompt string
	Depth  Depth
}

type expandRaw struct {
	Expanded    string `json:"expanded"`
	Added       []string `json:"added"`
	Assumptions []struct {
		Field string `json:"field"`
		Value string `json:"value"`
		Note  string `json:"note"`
	} `json:"assumptions"`
	Warnings []string `json:"warnings"`
}

// Expand turns an underspecified prompt into a detailed one. Additions that
// are not in the original are labeled as assumptions the author can correct.
func (e *Engine) Expand(ctx context.Context, in ExpandInput) (*types.ExpandResult, error) {
	if strings.TrimSpace(in.Prompt) == "" {
		return nil, apperr.Usagef("expand needs a prompt (argument, file, or stdin)")
	}
	if in.Depth == "" {
		in.Depth = DepthDetailed
	}
	if !ValidDepth(string(in.Depth)) {
		return nil, apperr.Usagef("--depth must be concise, detailed, or production")
	}

	user := "PROMPT TO EXPAND:\n\n" + in.Prompt + "\n\nDEPTH LEVEL: " + string(in.Depth)

	var raw expandRaw
	if _, err := e.call(ctx, "expand", user, &raw); err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.Expanded) == "" {
		return nil, apperr.New("The expander returned an empty prompt.",
			"Retry the command.").WithCode(apperr.CodeMalformed)
	}

	res := &types.ExpandResult{
		Operation: "expand",
		Result:    strings.TrimSpace(raw.Expanded),
		Depth:     string(in.Depth),
		Added:     cleanList(raw.Added),
		Tokens:    e.tokens(in.Prompt, raw.Expanded),
		Warnings:  cleanList(raw.Warnings),
	}
	for _, a := range raw.Assumptions {
		if strings.TrimSpace(a.Value) == "" {
			continue
		}
		res.Assumptions = append(res.Assumptions, types.Assumption{
			Field: strings.TrimSpace(a.Field),
			Value: strings.TrimSpace(a.Value),
			Note:  strings.TrimSpace(a.Note),
		})
	}
	return res, nil
}
