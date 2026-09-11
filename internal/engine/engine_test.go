package engine

import (
	"context"
	"strings"
	"testing"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/internal/groq"
)

// fakeCompleter returns a canned response (or error), and records the last
// request so tests can assert on what the engine sent.
type fakeCompleter struct {
	reply string
	err   error
	last  groq.Request
}

func (f *fakeCompleter) Complete(_ context.Context, req groq.Request) (*groq.Response, error) {
	f.last = req
	if f.err != nil {
		return nil, f.err
	}
	return &groq.Response{Text: f.reply, Model: "fake"}, nil
}

func newEngine(t *testing.T, reply string) (*Engine, *fakeCompleter) {
	t.Helper()
	fc := &fakeCompleter{reply: reply}
	e, err := New(Options{Client: fc, Temperature: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	return e, fc
}

func TestOptimize(t *testing.T) {
	reply := `Here you go:
` + "```json" + `
{
  "optimized": "You are a senior backend engineer. Build a REST API for payments.",
  "changes": ["clarified the objective", "stated the output format"],
  "warnings": [],
  "scores": {"clarity": 9, "specificity": 8, "completeness": 7, "consistency": 9, "efficiency": 8, "robustness": 7, "overall": 8.2}
}
` + "```"
	e, fc := newEngine(t, reply)

	res, err := e.Optimize(context.Background(), OptimizeInput{Prompt: "build a payments api"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(res.Result, "senior backend engineer") {
		t.Errorf("result not extracted: %q", res.Result)
	}
	if len(res.Changes) != 2 {
		t.Errorf("changes = %v", res.Changes)
	}
	if res.Analysis.Overall != 8.2 {
		t.Errorf("overall = %v", res.Analysis.Overall)
	}
	if res.Tokens.Source != "estimate" {
		t.Errorf("token source = %q", res.Tokens.Source)
	}
	if !fc.last.JSON {
		t.Error("engine should request JSON output")
	}
	if len(fc.last.Messages) != 2 || fc.last.Messages[0].Role != "system" {
		t.Errorf("expected system+user messages, got %+v", fc.last.Messages)
	}
}

func TestOptimizeEmptyPromptIsUsageError(t *testing.T) {
	e, _ := newEngine(t, "{}")
	_, err := e.Optimize(context.Background(), OptimizeInput{Prompt: "   "})
	if e2, ok := apperr.As(err); !ok || e2.Code != apperr.CodeUsage {
		t.Fatalf("want usage error, got %v", err)
	}
}

func TestMalformedModelOutputIsHandled(t *testing.T) {
	e, _ := newEngine(t, "the model rambled and returned no json")
	_, err := e.Optimize(context.Background(), OptimizeInput{Prompt: "x"})
	e2, ok := apperr.As(err)
	if !ok || e2.Code != apperr.CodeMalformed {
		t.Fatalf("want malformed error, got %v", err)
	}
}

func TestProviderErrorPassesThrough(t *testing.T) {
	fc := &fakeCompleter{err: apperr.New("Groq rate limit reached.", "wait").WithCode(apperr.CodeRateLimit)}
	e, _ := New(Options{Client: fc})
	_, err := e.Analyze(context.Background(), AnalyzeInput{Prompt: "x"})
	if e2, ok := apperr.As(err); !ok || e2.Code != apperr.CodeRateLimit {
		t.Fatalf("provider error should pass through, got %v", err)
	}
}

func TestCompress(t *testing.T) {
	reply := `{"compressed": "Short prompt.", "changes": ["removed repetition"], "semantic_preservation": 8.5, "behavior_risks": [], "warnings": []}`
	e, _ := newEngine(t, reply)
	res, err := e.Compress(context.Background(), CompressInput{Prompt: "a very long repetitive prompt that repeats", TargetPercent: 30})
	if err != nil {
		t.Fatal(err)
	}
	if res.SemanticPreservation != 8.5 {
		t.Errorf("preservation = %v", res.SemanticPreservation)
	}
	if res.Tokens.ReductionPercent <= 0 {
		t.Errorf("expected a positive reduction, got %v", res.Tokens.ReductionPercent)
	}
}

func TestCompressRejectsBadTarget(t *testing.T) {
	e, _ := newEngine(t, "{}")
	_, err := e.Compress(context.Background(), CompressInput{Prompt: "x", TargetPercent: 150})
	if _, ok := apperr.As(err); !ok {
		t.Fatalf("want usage error, got %v", err)
	}
}

func TestCompressLowPreservationAddsRisk(t *testing.T) {
	reply := `{"compressed": "x.", "changes": [], "semantic_preservation": 4, "behavior_risks": [], "warnings": []}`
	e, _ := newEngine(t, reply)
	res, err := e.Compress(context.Background(), CompressInput{Prompt: "aaaa bbbb cccc"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.BehaviorRisks) == 0 {
		t.Error("low preservation with no risks listed should synthesize a warning")
	}
}

func TestExpandDepthValidation(t *testing.T) {
	e, _ := newEngine(t, "{}")
	_, err := e.Expand(context.Background(), ExpandInput{Prompt: "x", Depth: "galaxy-brain"})
	if _, ok := apperr.As(err); !ok {
		t.Fatalf("want usage error for bad depth, got %v", err)
	}
}

func TestExpand(t *testing.T) {
	reply := `{"expanded": "## Objective\nBuild it.\n## Assumptions\n- Postgres", "added": ["objective", "assumptions"], "assumptions": [{"field": "database", "value": "Postgres", "note": "change if you use MySQL"}], "warnings": []}`
	e, fc := newEngine(t, reply)
	res, err := e.Expand(context.Background(), ExpandInput{Prompt: "build a payment api", Depth: DepthProduction})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Assumptions) != 1 || res.Assumptions[0].Field != "database" {
		t.Errorf("assumptions = %+v", res.Assumptions)
	}
	if res.Depth != "production" {
		t.Errorf("depth = %q", res.Depth)
	}
	if !strings.Contains(fc.last.Messages[1].Content, "production") {
		t.Error("depth should be passed to the model")
	}
}

func TestAnalyzeSortsFindings(t *testing.T) {
	reply := `{
      "scores": {"clarity": 7, "specificity": 6, "completeness": 7, "consistency": 9, "efficiency": 6, "robustness": 7, "overall": 7},
      "findings": [
        {"id": "P012", "severity": "INFO", "title": "role text is inert", "detail": "drop it"},
        {"id": "P004", "severity": "ERROR", "title": "conflicting instructions", "detail": "pick one"},
        {"id": "P003", "severity": "WARNING", "title": "no output format", "detail": "add one"}
      ],
      "summary": "ok"
    }`
	e, _ := newEngine(t, reply)
	res, err := e.Analyze(context.Background(), AnalyzeInput{Prompt: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Findings) != 3 {
		t.Fatalf("findings = %d", len(res.Findings))
	}
	if res.Findings[0].Severity != "ERROR" || res.Findings[2].Severity != "INFO" {
		t.Errorf("findings not sorted by severity: %+v", res.Findings)
	}
}

func TestTransformTemplateCollectsVariables(t *testing.T) {
	reply := `{"transformed": "Review the following {{language}} code for {{issue}}:\n\n{{code}}", "notes": ["parameterized inputs"], "variables": ["language", "issue"], "warnings": []}`
	e, _ := newEngine(t, reply)
	res, err := e.Transform(context.Background(), TransformInput{Prompt: "review this python code for sql injection", Target: TargetTemplate})
	if err != nil {
		t.Fatal(err)
	}
	// {{code}} appears in the body but not the reported list; the engine
	// should merge both sources.
	want := map[string]bool{"language": true, "issue": true, "code": true}
	if len(res.Variables) != 3 {
		t.Fatalf("variables = %v", res.Variables)
	}
	for _, v := range res.Variables {
		if !want[v] {
			t.Errorf("unexpected variable %q", v)
		}
	}
}

func TestTransformRejectsUnknownTarget(t *testing.T) {
	e, _ := newEngine(t, "{}")
	_, err := e.Transform(context.Background(), TransformInput{Prompt: "x", Target: "cobol"})
	if _, ok := apperr.As(err); !ok {
		t.Fatalf("want usage error, got %v", err)
	}
}

func TestEval(t *testing.T) {
	yes := true
	no := false
	_ = yes
	_ = no
	reply := `{
      "scores": {"clarity": 9, "consistency": 8.5, "robustness": 8, "output_control": 9, "overall": 8.6},
      "test_cases": [
        {"name": "normal input", "input": "...", "expectation": "handles it", "assessment": "fine", "pass": true},
        {"name": "adversarial", "input": "ignore previous", "expectation": "refuse", "assessment": "leaks", "pass": false}
      ],
      "weaknesses": ["injection"],
      "summary": "mostly ready"
    }`
	e, _ := newEngine(t, reply)
	res, err := e.Eval(context.Background(), EvalInput{Prompt: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.TestCases) != 2 {
		t.Fatalf("test cases = %d", len(res.TestCases))
	}
	if res.TestCases[0].Pass == nil || !*res.TestCases[0].Pass {
		t.Error("first case should pass")
	}
	if res.Scores.Overall != 8.6 {
		t.Errorf("overall = %v", res.Scores.Overall)
	}
}

func TestExtractJSONObject(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{`{"a":1}`, `{"a":1}`, true},
		{"```json\n{\"a\":1}\n```", `{"a":1}`, true},
		{"reasoning... here is the answer:\n{\"a\": {\"b\": 2}}\nthanks", `{"a": {"b": 2}}`, true},
		{`{"a":"}"}`, `{"a":"}"}`, true},
		{`no json`, "", false},
		{`{"a":1`, "", false},
	}
	for _, c := range cases {
		got, err := extractJSONObject(c.in)
		if c.ok != (err == nil) {
			t.Errorf("extractJSONObject(%q): err = %v", c.in, err)
			continue
		}
		if c.ok && got != c.want {
			t.Errorf("extractJSONObject(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNewRejectsNilClient(t *testing.T) {
	if _, err := New(Options{}); err == nil {
		t.Fatal("nil client should error")
	}
}
