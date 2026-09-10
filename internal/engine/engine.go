// Package engine implements the six promptopt operations: optimize, compress,
// expand, analyze, transform, and eval.
//
// The engine has no knowledge of Cobra, terminals, or JSON output formatting.
// It takes a prompt and options, calls a groq.Completer, validates the
// model's structured response, and returns a typed result from pkg/types.
// This keeps every operation testable with a fake Completer and no CLI.
package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/internal/groq"
	"github.com/rakshit-gen/promptopt/internal/tokenizer"
	"github.com/rakshit-gen/promptopt/pkg/types"
	"github.com/rakshit-gen/promptopt/prompts"
)

// Engine runs operations against a completion provider.
type Engine struct {
	client      groq.Completer
	tok         tokenizer.Estimator
	model       string
	temperature float64
}

// Options configures a new Engine.
type Options struct {
	Client      groq.Completer
	Tokenizer   tokenizer.Estimator
	Model       string
	Temperature float64
}

// New builds an Engine. Client is required.
func New(opts Options) (*Engine, error) {
	if opts.Client == nil {
		return nil, fmt.Errorf("engine: nil completion client")
	}
	tok := opts.Tokenizer
	if tok == nil {
		tok = tokenizer.NewHeuristic()
	}
	return &Engine{
		client:      opts.Client,
		tok:         tok,
		model:       opts.Model,
		temperature: opts.Temperature,
	}, nil
}

// tokens builds a TokenReport from estimated counts of the before/after text.
// Provider usage counts describe the whole request and response, not the
// prompt text itself, so promptopt reports the estimate here and labels it
// as such. The tokenizer package is where a model-exact count would land.
func (e *Engine) tokens(before, after string) types.TokenReport {
	b := e.tok.Count(before)
	a := e.tok.Count(after)
	return types.TokenReport{
		Before:           b,
		After:            a,
		ReductionPercent: round2(tokenizer.Reduction(b, a)),
		Source:           "estimate",
	}
}

func round2(f float64) float64 {
	return float64(int64(f*100+0.5*sign(f))) / 100
}

func sign(f float64) float64 {
	if f < 0 {
		return -1
	}
	return 1
}

// call runs one structured completion: system = embedded op instructions,
// user = the assembled request. It requests a JSON object, extracts it from
// whatever the model returned, and unmarshals into out.
func (e *Engine) call(ctx context.Context, op, userContent string, out any) (groq.Usage, error) {
	system, err := prompts.Load(op)
	if err != nil {
		return groq.Usage{}, apperr.Wrap(err, "Internal error: missing prompt template.", "")
	}

	resp, err := e.client.Complete(ctx, groq.Request{
		Model:       e.model,
		Temperature: e.temperature,
		JSON:        true,
		Messages: []groq.Message{
			{Role: "system", Content: system},
			{Role: "user", Content: userContent},
		},
	})
	if err != nil {
		return groq.Usage{}, err
	}

	obj, err := extractJSONObject(resp.Text)
	if err != nil {
		return resp.Usage, apperr.Wrap(err,
			"The model did not return usable JSON for this operation.",
			"Retry the command. If it persists, try --model with a model that\n"+
				"supports structured output, such as openai/gpt-oss-120b.").
			WithCode(apperr.CodeMalformed)
	}

	dec := json.NewDecoder(strings.NewReader(obj))
	dec.DisallowUnknownFields()
	if err := dec.Decode(out); err != nil {
		// Retry once more leniently: unknown fields are not fatal.
		if err2 := json.Unmarshal([]byte(obj), out); err2 != nil {
			return resp.Usage, apperr.Wrap(err2,
				"The model's JSON response did not match the expected shape.",
				"Retry the command. This is usually transient.").
				WithCode(apperr.CodeMalformed)
		}
	}
	return resp.Usage, nil
}

// extractJSONObject pulls the first balanced top-level JSON object out of a
// model response. It tolerates surrounding prose, ```json fences, and
// reasoning text that some models emit before the answer.
func extractJSONObject(s string) (string, error) {
	s = strings.TrimSpace(s)

	// Strip a fenced block if present.
	if i := strings.Index(s, "```"); i >= 0 {
		rest := s[i+3:]
		rest = strings.TrimPrefix(rest, "json")
		rest = strings.TrimPrefix(rest, "JSON")
		if j := strings.Index(rest, "```"); j >= 0 {
			candidate := strings.TrimSpace(rest[:j])
			if strings.HasPrefix(candidate, "{") {
				s = candidate
			}
		}
	}

	start := strings.Index(s, "{")
	if start < 0 {
		return "", fmt.Errorf("no JSON object found in response")
	}

	depth := 0
	inStr := false
	escaped := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if inStr {
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == '"':
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1], nil
			}
		}
	}
	return "", fmt.Errorf("unbalanced JSON object in response")
}

// clampScore keeps a model-provided score inside 0-10.
func clampScore(f float64) float64 {
	if f < 0 {
		return 0
	}
	if f > 10 {
		return 10
	}
	return f
}
