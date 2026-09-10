// Package tokenizer provides token count estimates for prompts.
//
// promptopt does not ship the exact BPE tables for every Groq model, so by
// default it reports an estimate. When the Groq API returns a usage object
// with real prompt/completion counts, callers should prefer those and mark
// the source as "provider". The Estimator interface exists so a model-exact
// tokenizer can be dropped in later without touching call sites.
package tokenizer

import (
	"strings"
	"unicode"
)

// Estimator returns an approximate token count for a piece of text.
type Estimator interface {
	// Count returns the estimated number of tokens.
	Count(text string) int
	// Name identifies the estimator for display ("heuristic", "cl100k", ...).
	Name() string
}

// Heuristic is a tokenizer-free estimate. It blends a character-based ratio
// with a word-based ratio, which tracks real BPE tokenizers to within a few
// percent for typical English prompts and prompt-like text (code, markdown,
// XML). It intentionally errs slightly high so token budgets are not
// underestimated.
type Heuristic struct{}

// NewHeuristic returns the default estimator.
func NewHeuristic() Heuristic { return Heuristic{} }

// Name implements Estimator.
func (Heuristic) Name() string { return "heuristic" }

// Count implements Estimator.
func (Heuristic) Count(text string) int {
	if strings.TrimSpace(text) == "" {
		return 0
	}

	runes := []rune(text)
	chars := len(runes)

	// Word/segment count: runs of letters or digits, with punctuation and
	// symbols each counting as their own segment (BPE tends to split these).
	words := 0
	symbols := 0
	inWord := false
	for _, r := range runes {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if !inWord {
				words++
				inWord = true
			}
		case unicode.IsSpace(r):
			inWord = false
		default:
			inWord = false
			symbols++
		}
	}

	// Character model: ~4 chars per token is the common rule of thumb.
	charEst := float64(chars) / 4.0
	// Word model: ~1.3 tokens per word for English, plus symbols roughly 1:1.
	wordEst := float64(words)*1.3 + float64(symbols)

	est := (charEst*0.5 + wordEst*0.5)
	// Nudge up by 3% to stay on the safe side of budgets.
	est *= 1.03

	n := int(est + 0.5)
	if n < 1 {
		n = 1
	}
	return n
}

// Reduction returns the percentage reduction from before to after. Positive
// means the text got shorter. Returns 0 when before is 0.
func Reduction(before, after int) float64 {
	if before == 0 {
		return 0
	}
	return (float64(before-after) / float64(before)) * 100
}
