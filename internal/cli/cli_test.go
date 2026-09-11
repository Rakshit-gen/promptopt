package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// runCLI executes the root command with args and a fake Groq backend, and
// returns stdout, stderr, and the resulting error.
func runCLI(t *testing.T, groqBody string, groqStatus int, stdin string, args ...string) (string, string, error) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if groqStatus != 0 && groqStatus != http.StatusOK {
			w.WriteHeader(groqStatus)
		}
		w.Write([]byte(groqBody))
	}))
	t.Cleanup(srv.Close)

	t.Setenv("GROQ_API_KEY", "gsk_test_key")
	t.Setenv("PROMPTOPT_BASE_URL", srv.URL)
	t.Setenv("PROMPTOPT_MODEL", "test-model")
	t.Setenv("NO_COLOR", "1")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir()) // isolate from a real config file

	root := newRootCmd()
	var out, errb bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errb)
	root.SetIn(strings.NewReader(stdin))
	root.SetArgs(args)

	err := root.Execute()
	return out.String(), errb.String(), err
}

func groqJSON(payload string) string {
	b, _ := json.Marshal(payload)
	return `{"model":"test-model","choices":[{"message":{"content":` + string(b) + `}}],"usage":{"prompt_tokens":50,"completion_tokens":40,"total_tokens":90}}`
}

func TestCLIOptimizeJSON(t *testing.T) {
	model := `{"optimized":"You are a senior engineer. Build a payments REST API with idempotency.","changes":["clarified objective"],"warnings":[],"scores":{"clarity":9,"specificity":8,"completeness":7,"consistency":9,"efficiency":8,"robustness":7,"overall":8.1}}`
	out, _, err := runCLI(t, groqJSON(model), 200, "", "optimize", "build a payments api", "--json")
	if err != nil {
		t.Fatal(err)
	}
	var res struct {
		Operation string `json:"operation"`
		Result    string `json:"result"`
		Tokens    struct {
			Before int `json:"before"`
		} `json:"tokens"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("not json: %v\n%s", err, out)
	}
	if res.Operation != "optimize" || !strings.Contains(res.Result, "senior engineer") {
		t.Errorf("unexpected result: %+v", res)
	}
	if res.Tokens.Before == 0 {
		t.Error("token count not populated")
	}
}

func TestCLIAnalyzeFromStdin(t *testing.T) {
	model := `{"scores":{"clarity":6,"specificity":5,"completeness":6,"consistency":8,"efficiency":7,"robustness":6,"overall":6.2},"findings":[{"id":"P003","severity":"WARNING","title":"no output format","detail":"add one"}],"summary":"vague"}`
	out, _, err := runCLI(t, groqJSON(model), 200, "review this code\n", "analyze")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"analyze", "6.2", "P003", "no output format"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

func TestCLIMissingAPIKey(t *testing.T) {
	t.Setenv("GROQ_API_KEY", "")
	t.Setenv("PROMPTOPT_GROQ_API_KEY", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	root := newRootCmd()
	var out, errb bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errb)
	root.SetArgs([]string{"optimize", "hello"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected an error without an API key")
	}
	code := reportError(&errb, err, false)
	if code != 4 {
		t.Errorf("exit code = %d, want 4 (auth)", code)
	}
	if !strings.Contains(errb.String(), "GROQ_API_KEY") {
		t.Errorf("error should tell the user to set GROQ_API_KEY:\n%s", errb.String())
	}
}

func TestCLITransformNeedsTarget(t *testing.T) {
	_, _, err := runCLI(t, groqJSON("{}"), 200, "", "transform", "some prompt")
	if err == nil {
		t.Fatal("transform without --to should fail")
	}
	if !strings.Contains(err.Error(), "--to") {
		t.Errorf("error should mention --to: %v", err)
	}
}

func TestCLIRateLimitExitCode(t *testing.T) {
	body := `{"error":{"message":"Rate limit reached for model"}}`
	_, _, err := runCLI(t, body, http.StatusTooManyRequests, "", "eval", "a prompt")
	if err == nil {
		t.Fatal("expected rate-limit error")
	}
	var errb bytes.Buffer
	if code := reportError(&errb, err, false); code != 5 {
		t.Errorf("rate-limit exit code = %d, want 5", code)
	}
}

func TestCLIHelpListsSixOperations(t *testing.T) {
	root := newRootCmd()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"--help"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, op := range []string{"optimize", "compress", "expand", "analyze", "transform", "eval"} {
		if !strings.Contains(out.String(), op) {
			t.Errorf("root help missing %q", op)
		}
	}
}

func TestCLIQuietPrintsOnlyResult(t *testing.T) {
	model := `{"compressed":"Tight prompt.","changes":["cut repetition"],"semantic_preservation":9,"behavior_risks":[],"warnings":[]}`
	out, _, err := runCLI(t, groqJSON(model), 200, "a long prompt that repeats itself a lot and repeats\n", "compress", "-q")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out) != "Tight prompt." {
		t.Errorf("quiet output = %q", out)
	}
}
