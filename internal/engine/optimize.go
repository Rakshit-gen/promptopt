package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/pkg/types"
)

// OptimizeInput configures an optimize run.
type OptimizeInput struct {
	Prompt string
	// Instruction is optional extra guidance from the user (e.g. --goal).
	Instruction string
}

type optimizeRaw struct {
	Optimized string   `json:"optimized"`
	Changes   []string `json:"changes"`
	Warnings  []string `json:"warnings"`
	Scores    struct {
		Clarity      float64 `json:"clarity"`
		Specificity  float64 `json:"specificity"`
		Completeness float64 `json:"completeness"`
		Consistency  float64 `json:"consistency"`
		Efficiency   float64 `json:"efficiency"`
		Robustness   float64 `json:"robustness"`
		Overall      float64 `json:"overall"`
	} `json:"scores"`
}

// Optimize inspects a prompt and returns an improved version with a summary of
// what changed. It preserves intent and does not add domain requirements.
func (e *Engine) Optimize(ctx context.Context, in OptimizeInput) (*types.OptimizeResult, error) {
	if strings.TrimSpace(in.Prompt) == "" {
		return nil, apperr.Usagef("optimize needs a prompt (argument, file, or stdin)")
	}

	user := "PROMPT TO OPTIMIZE:\n\n" + in.Prompt
	if s := strings.TrimSpace(in.Instruction); s != "" {
		user += "\n\nADDITIONAL GUIDANCE FROM THE AUTHOR:\n" + s
	}

	var raw optimizeRaw
	usage, err := e.call(ctx, "optimize", user, &raw)
	if err != nil {
		return nil, err
	}
	_ = usage

	if strings.TrimSpace(raw.Optimized) == "" {
		return nil, apperr.New("The optimizer returned an empty prompt.",
			"Retry the command.").WithCode(apperr.CodeMalformed)
	}

	res := &types.OptimizeResult{
		Operation: "optimize",
		Result:    strings.TrimSpace(raw.Optimized),
		Changes:   cleanList(raw.Changes),
		Warnings:  cleanList(raw.Warnings),
		Tokens:    e.tokens(in.Prompt, raw.Optimized),
		Analysis: types.ScoreCard{
			Overall:      clampScore(raw.Scores.Overall),
			Clarity:      clampScore(raw.Scores.Clarity),
			Specificity:  clampScore(raw.Scores.Specificity),
			Completeness: clampScore(raw.Scores.Completeness),
			Consistency:  clampScore(raw.Scores.Consistency),
			Efficiency:   clampScore(raw.Scores.Efficiency),
			Robustness:   clampScore(raw.Scores.Robustness),
		},
	}

	if len(res.Changes) == 0 {
		res.Changes = []string{"no changes needed; the prompt was already in good shape"}
	}
	if res.Tokens.ReductionPercent < -25 {
		res.Warnings = append(res.Warnings, fmt.Sprintf(
			"the optimized prompt is %.0f%% longer; run `promptopt compress` if token budget matters",
			-res.Tokens.ReductionPercent))
	}
	return res, nil
}

// cleanList trims entries and drops empties.
func cleanList(in []string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}
