package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/pkg/types"
)

// CompressInput configures a compress run.
type CompressInput struct {
	Prompt string
	// TargetPercent is the desired reduction (0-95). Zero means "as much as is
	// safe".
	TargetPercent int
	// Aggressive allows dropping marginal context and collapsing examples.
	Aggressive bool
	// PreserveBehavior refuses any cut the model is not confident is safe.
	PreserveBehavior bool
}

type compressRaw struct {
	Compressed           string   `json:"compressed"`
	Changes              []string `json:"changes"`
	SemanticPreservation float64  `json:"semantic_preservation"`
	BehaviorRisks        []string `json:"behavior_risks"`
	Warnings             []string `json:"warnings"`
}

// Compress reduces a prompt's token count while trying to preserve its
// behavior. A shorter prompt that changes what the model does is reported as a
// risk, not a success.
func (e *Engine) Compress(ctx context.Context, in CompressInput) (*types.CompressResult, error) {
	if strings.TrimSpace(in.Prompt) == "" {
		return nil, apperr.Usagef("compress needs a prompt (argument, file, or stdin)")
	}
	if in.TargetPercent < 0 || in.TargetPercent > 95 {
		return nil, apperr.Usagef("--target must be between 0 and 95")
	}

	var b strings.Builder
	b.WriteString("PROMPT TO COMPRESS:\n\n")
	b.WriteString(in.Prompt)
	b.WriteString("\n\nREQUEST:\n")
	if in.TargetPercent > 0 {
		fmt.Fprintf(&b, "- target reduction: about %d%%\n", in.TargetPercent)
	} else {
		b.WriteString("- target reduction: as much as is safe\n")
	}
	fmt.Fprintf(&b, "- aggressive: %t\n", in.Aggressive)
	fmt.Fprintf(&b, "- preserve-behavior: %t\n", in.PreserveBehavior)

	var raw compressRaw
	if _, err := e.call(ctx, "compress", b.String(), &raw); err != nil {
		return nil, err
	}
	if strings.TrimSpace(raw.Compressed) == "" {
		return nil, apperr.New("The compressor returned an empty prompt.",
			"Retry the command.").WithCode(apperr.CodeMalformed)
	}

	sem := clampScore(raw.SemanticPreservation)
	res := &types.CompressResult{
		Operation:            "compress",
		Result:               strings.TrimSpace(raw.Compressed),
		Changes:              cleanList(raw.Changes),
		Tokens:               e.tokens(in.Prompt, raw.Compressed),
		SemanticPreservation: sem,
		BehaviorRisks:        cleanList(raw.BehaviorRisks),
		Warnings:             cleanList(raw.Warnings),
	}

	if res.Tokens.ReductionPercent <= 0 {
		res.Warnings = append(res.Warnings,
			"no tokens were saved; the prompt may already be compact")
	}
	if sem < 7 && len(res.BehaviorRisks) == 0 {
		res.BehaviorRisks = append(res.BehaviorRisks,
			"the model reported low confidence that behavior is preserved but gave no specifics; review the diff carefully")
	}
	if in.TargetPercent > 0 {
		got := res.Tokens.ReductionPercent
		if got+5 < float64(in.TargetPercent) {
			res.Warnings = append(res.Warnings, fmt.Sprintf(
				"reached %.1f%% against a %d%% target; deeper cuts were judged unsafe",
				got, in.TargetPercent))
		}
	}
	return res, nil
}
