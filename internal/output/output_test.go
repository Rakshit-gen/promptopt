package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/rakshit-gen/promptopt/pkg/types"
)

func newTestWriter(color bool) (*Writer, *bytes.Buffer, *bytes.Buffer) {
	var out, errb bytes.Buffer
	w := New(Options{Out: &out, Err: &errb, Color: color})
	return w, &out, &errb
}

func TestJSONHasNoDecoration(t *testing.T) {
	w, out, _ := newTestWriter(true)
	res := &types.CompressResult{
		Operation:            "compress",
		Result:               "short prompt",
		Tokens:               types.TokenReport{Before: 1842, After: 1103, ReductionPercent: 40.12, Source: "estimate"},
		SemanticPreservation: 8.5,
	}
	if err := w.JSON(res); err != nil {
		t.Fatal(err)
	}
	var back types.CompressResult
	if err := json.Unmarshal(out.Bytes(), &back); err != nil {
		t.Fatalf("json output is not valid json: %v\n%s", err, out.String())
	}
	if strings.ContainsAny(out.String(), "\x1b") {
		t.Error("json output contains ANSI escape codes")
	}
	if back.Tokens.ReductionPercent != 40.12 {
		t.Errorf("round-tripped reduction = %v", back.Tokens.ReductionPercent)
	}
}

func TestHumanOutputNoColorHasNoEscapes(t *testing.T) {
	w, out, _ := newTestWriter(false)
	w.Analyze(&types.AnalyzeResult{
		Operation: "analyze",
		Scores:    types.ScoreCard{Overall: 7.8, Clarity: 8.4, Consistency: 9.2},
		Findings: []types.Finding{
			{ID: "P003", Severity: types.SeverityWarning, Title: "no output format", Detail: "add one"},
		},
		Summary: "decent prompt",
	})
	s := out.String()
	if strings.Contains(s, "\x1b[") {
		t.Errorf("no-color output still has ANSI escapes:\n%q", s)
	}
	for _, want := range []string{"analyze", "7.8", "P003", "no output format", "decent prompt"} {
		if !strings.Contains(s, want) {
			t.Errorf("output missing %q\n%s", want, s)
		}
	}
}

func TestHumanIntGrouping(t *testing.T) {
	cases := map[int]string{0: "0", 42: "42", 1842: "1,842", 1000000: "1,000,000", -1391: "-1,391"}
	for in, want := range cases {
		if got := humanInt(in); got != want {
			t.Errorf("humanInt(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestWrapText(t *testing.T) {
	lines := wrapText("the quick brown fox jumps over the lazy dog again and again", 20)
	for _, ln := range lines {
		if len(ln) > 20 {
			t.Errorf("line exceeds width: %q", ln)
		}
	}
	if strings.Join(lines, " ") != "the quick brown fox jumps over the lazy dog again and again" {
		t.Errorf("wrap lost content: %v", lines)
	}
}

func TestQuietPlainResult(t *testing.T) {
	w, out, _ := newTestWriter(false)
	w.PlainResult("just the prompt\n\n")
	if out.String() != "just the prompt\n" {
		t.Errorf("PlainResult = %q", out.String())
	}
}
