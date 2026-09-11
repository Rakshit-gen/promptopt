//go:build integration

// Package tests holds integration tests that make real calls to the Groq API.
// They are excluded from `go test ./...` by the build tag and only run when
// explicitly requested:
//
//	PROMPTOPT_INTEGRATION=1 GROQ_API_KEY=... go test -tags=integration ./tests/...
//
// These verify that the structured-output contract still holds against the
// live model, which unit tests with a fake completer cannot catch.
package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rakshit-gen/promptopt/internal/engine"
	"github.com/rakshit-gen/promptopt/internal/groq"
)

func liveEngine(t *testing.T) *engine.Engine {
	t.Helper()
	if os.Getenv("PROMPTOPT_INTEGRATION") == "" {
		t.Skip("set PROMPTOPT_INTEGRATION=1 to run integration tests")
	}
	key := os.Getenv("GROQ_API_KEY")
	if key == "" {
		t.Skip("GROQ_API_KEY not set")
	}
	client, err := groq.New(groq.Config{
		APIKey:     key,
		Model:      envOr("PROMPTOPT_MODEL", groq.DefaultModel),
		Timeout:    60 * time.Second,
		MaxRetries: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	e, err := engine.New(engine.Options{Client: client, Temperature: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

const samplePrompt = `You are helpful. You are a helpful assistant. Help the user with their
request. When they ask for code, write the code. Be helpful. Write good
code. Always help. Do not be unhelpful.`

func TestIntegrationOptimize(t *testing.T) {
	e := liveEngine(t)
	res, err := e.Optimize(context.Background(), engine.OptimizeInput{Prompt: samplePrompt})
	if err != nil {
		t.Fatal(err)
	}
	if res.Result == "" {
		t.Fatal("empty optimized prompt")
	}
	if len(res.Changes) == 0 {
		t.Error("expected the optimizer to report at least one change for a redundant prompt")
	}
	if res.Analysis.Overall <= 0 || res.Analysis.Overall > 10 {
		t.Errorf("overall score out of range: %v", res.Analysis.Overall)
	}
}

func TestIntegrationAnalyze(t *testing.T) {
	e := liveEngine(t)
	res, err := e.Analyze(context.Background(), engine.AnalyzeInput{Prompt: samplePrompt})
	if err != nil {
		t.Fatal(err)
	}
	if res.Scores.Overall <= 0 {
		t.Errorf("no overall score: %+v", res.Scores)
	}
	if len(res.Findings) == 0 {
		t.Error("expected findings for a prompt with repeated instructions")
	}
}

func TestIntegrationTransformXML(t *testing.T) {
	e := liveEngine(t)
	res, err := e.Transform(context.Background(), engine.TransformInput{
		Prompt: "Review this Python code for SQL injection vulnerabilities.",
		Target: engine.TargetXML,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Result == "" {
		t.Fatal("empty transform result")
	}
}

func TestIntegrationEval(t *testing.T) {
	e := liveEngine(t)
	res, err := e.Eval(context.Background(), engine.EvalInput{Prompt: samplePrompt})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.TestCases) == 0 {
		t.Error("expected generated test cases")
	}
}
