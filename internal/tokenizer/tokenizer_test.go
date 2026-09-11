package tokenizer

import "testing"

func TestHeuristicCount(t *testing.T) {
	h := NewHeuristic()

	if got := h.Count(""); got != 0 {
		t.Fatalf("empty string: got %d, want 0", got)
	}
	if got := h.Count("   \n\t "); got != 0 {
		t.Fatalf("whitespace only: got %d, want 0", got)
	}

	// A ~50-word paragraph should land in a plausible token range. Real BPE
	// tokenizers put typical English around 1.3 tokens/word.
	para := "The optimizer inspects a prompt before rewriting it. " +
		"It looks for an unclear objective, missing context, ambiguous language, " +
		"redundant instructions, conflicting requirements, and a missing output format. " +
		"Then it produces an improved prompt that preserves the original intent."
	got := h.Count(para)
	if got < 40 || got > 90 {
		t.Fatalf("paragraph estimate %d outside expected 40-90 range", got)
	}
}

func TestHeuristicMonotonic(t *testing.T) {
	h := NewHeuristic()
	short := h.Count("build a REST API")
	long := h.Count("build a REST API for payments with idempotency keys and webhooks")
	if long <= short {
		t.Fatalf("longer text should estimate more tokens: short=%d long=%d", short, long)
	}
}

func TestReduction(t *testing.T) {
	cases := []struct {
		before, after int
		want          float64
	}{
		{1000, 750, 25},
		{1000, 1000, 0},
		{1000, 1200, -20},
		{0, 100, 0},
	}
	for _, c := range cases {
		if got := Reduction(c.before, c.after); got != c.want {
			t.Errorf("Reduction(%d,%d) = %v, want %v", c.before, c.after, got, c.want)
		}
	}
}

func TestName(t *testing.T) {
	if NewHeuristic().Name() != "heuristic" {
		t.Fatal("unexpected estimator name")
	}
}
