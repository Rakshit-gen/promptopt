package engine

import (
	"context"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/pkg/types"
)

// EvalInput configures an eval run.
type EvalInput struct {
	Prompt string
}

type evalRaw struct {
	Scores struct {
		Clarity       float64 `json:"clarity"`
		Consistency   float64 `json:"consistency"`
		Robustness    float64 `json:"robustness"`
		OutputControl float64 `json:"output_control"`
		Overall       float64 `json:"overall"`
	} `json:"scores"`
	TestCases []struct {
		Name        string `json:"name"`
		Input       string `json:"input"`
		Expectation string `json:"expectation"`
		Assessment  string `json:"assessment"`
		Pass        *bool  `json:"pass"`
	} `json:"test_cases"`
	Weaknesses []string `json:"weaknesses"`
	Summary    string   `json:"summary"`
}

// Eval assesses a prompt and probes it with generated test cases. It is a
// foundation for a fuller evaluation system: the result shape leaves room for
// datasets and model comparisons without a breaking change.
func (e *Engine) Eval(ctx context.Context, in EvalInput) (*types.EvalResult, error) {
	if strings.TrimSpace(in.Prompt) == "" {
		return nil, apperr.Usagef("eval needs a prompt (argument, file, or stdin)")
	}

	var raw evalRaw
	if _, err := e.call(ctx, "eval", "PROMPT TO EVALUATE:\n\n"+in.Prompt, &raw); err != nil {
		return nil, err
	}

	res := &types.EvalResult{
		Operation:  "eval",
		Summary:    strings.TrimSpace(raw.Summary),
		Weaknesses: cleanList(raw.Weaknesses),
		Scores: types.ScoreCard{
			Overall:       clampScore(raw.Scores.Overall),
			Clarity:       clampScore(raw.Scores.Clarity),
			Consistency:   clampScore(raw.Scores.Consistency),
			Robustness:    clampScore(raw.Scores.Robustness),
			OutputControl: clampScore(raw.Scores.OutputControl),
		},
	}

	for _, tc := range raw.TestCases {
		name := strings.TrimSpace(tc.Name)
		if name == "" && strings.TrimSpace(tc.Input) == "" {
			continue
		}
		res.TestCases = append(res.TestCases, types.TestCase{
			Name:        name,
			Input:       strings.TrimSpace(tc.Input),
			Expectation: strings.TrimSpace(tc.Expectation),
			Assessment:  strings.TrimSpace(tc.Assessment),
			Pass:        tc.Pass,
		})
	}

	if res.Scores.Overall == 0 {
		s := res.Scores
		vals := []float64{s.Clarity, s.Consistency, s.Robustness, s.OutputControl}
		sum, n := 0.0, 0.0
		for _, v := range vals {
			if v > 0 {
				sum += v
				n++
			}
		}
		if n > 0 {
			res.Scores.Overall = round2(sum / n)
		}
	}
	return res, nil
}
