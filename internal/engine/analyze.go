package engine

import (
	"context"
	"sort"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/pkg/types"
)

// AnalyzeInput configures an analyze run.
type AnalyzeInput struct {
	Prompt string
}

type analyzeRaw struct {
	Scores struct {
		Clarity      float64 `json:"clarity"`
		Specificity  float64 `json:"specificity"`
		Completeness float64 `json:"completeness"`
		Consistency  float64 `json:"consistency"`
		Efficiency   float64 `json:"efficiency"`
		Robustness   float64 `json:"robustness"`
		Overall      float64 `json:"overall"`
	} `json:"scores"`
	Findings []struct {
		ID       string `json:"id"`
		Severity string `json:"severity"`
		Title    string `json:"title"`
		Detail   string `json:"detail"`
	} `json:"findings"`
	Summary string `json:"summary"`
}

var severityRank = map[types.Severity]int{
	types.SeverityError:   0,
	types.SeverityWarning: 1,
	types.SeverityInfo:    2,
}

// Analyze runs the static analyzer over a prompt and returns scores plus
// ordered, actionable findings.
func (e *Engine) Analyze(ctx context.Context, in AnalyzeInput) (*types.AnalyzeResult, error) {
	if strings.TrimSpace(in.Prompt) == "" {
		return nil, apperr.Usagef("analyze needs a prompt (argument, file, or stdin)")
	}

	var raw analyzeRaw
	if _, err := e.call(ctx, "analyze", "PROMPT TO ANALYZE:\n\n"+in.Prompt, &raw); err != nil {
		return nil, err
	}

	res := &types.AnalyzeResult{
		Operation: "analyze",
		Summary:   strings.TrimSpace(raw.Summary),
		Tokens:    e.tokens(in.Prompt, in.Prompt),
		Scores: types.ScoreCard{
			Overall:      clampScore(raw.Scores.Overall),
			Clarity:      clampScore(raw.Scores.Clarity),
			Specificity:  clampScore(raw.Scores.Specificity),
			Completeness: clampScore(raw.Scores.Completeness),
			Consistency:  clampScore(raw.Scores.Consistency),
			Efficiency:   clampScore(raw.Scores.Efficiency),
			Robustness:   clampScore(raw.Scores.Robustness),
		},
	}

	for _, f := range raw.Findings {
		sev := normalizeSeverity(f.Severity)
		title := strings.TrimSpace(f.Title)
		if title == "" {
			continue
		}
		res.Findings = append(res.Findings, types.Finding{
			ID:       strings.TrimSpace(f.ID),
			Severity: sev,
			Title:    title,
			Detail:   strings.TrimSpace(f.Detail),
		})
	}

	sort.SliceStable(res.Findings, func(i, j int) bool {
		return severityRank[res.Findings[i].Severity] < severityRank[res.Findings[j].Severity]
	})

	if res.Scores.Overall == 0 {
		res.Scores.Overall = round2(averageScore(res.Scores))
	}
	return res, nil
}

func normalizeSeverity(s string) types.Severity {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "ERROR", "CRITICAL", "HIGH":
		return types.SeverityError
	case "WARNING", "WARN", "MEDIUM":
		return types.SeverityWarning
	default:
		return types.SeverityInfo
	}
}

func averageScore(s types.ScoreCard) float64 {
	vals := []float64{s.Clarity, s.Specificity, s.Completeness, s.Consistency, s.Efficiency, s.Robustness}
	sum, n := 0.0, 0.0
	for _, v := range vals {
		if v > 0 {
			sum += v
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return sum / n
}
