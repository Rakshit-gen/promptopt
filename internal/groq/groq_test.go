package groq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rakshit-gen/promptopt/internal/apperr"
)

func testClient(t *testing.T, url string, retries int) *Client {
	t.Helper()
	c, err := New(Config{
		APIKey:     "gsk_test",
		BaseURL:    url,
		Model:      "test-model",
		Timeout:    2 * time.Second,
		MaxRetries: retries,
	})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestNewRequiresAPIKey(t *testing.T) {
	_, err := New(Config{})
	e, ok := apperr.As(err)
	if !ok || e.Code != apperr.CodeAuth {
		t.Fatalf("missing key should give an auth apperr, got %v", err)
	}
}

func TestCompleteSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer gsk_test" {
			t.Errorf("auth header = %q", got)
		}
		var body wireRequest
		json.NewDecoder(r.Body).Decode(&body)
		if body.Model != "test-model" || !body.ResponseFormat.jsonType() {
			t.Errorf("unexpected request body: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"model":"test-model","choices":[{"message":{"content":"{\"ok\":true}"}}],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}`))
	}))
	defer srv.Close()

	resp, err := testClient(t, srv.URL, 0).Complete(context.Background(), Request{
		Messages: []Message{{Role: "user", Content: "hi"}},
		JSON:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text != `{"ok":true}` {
		t.Errorf("text = %q", resp.Text)
	}
	if resp.Usage.TotalTokens != 15 {
		t.Errorf("usage not parsed: %+v", resp.Usage)
	}
}

func (rf *responseFormat) jsonType() bool { return rf != nil && rf.Type == "json_object" }

func TestCompleteRetriesOn5xx(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 3 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.Write([]byte(`{"choices":[{"message":{"content":"done"}}]}`))
	}))
	defer srv.Close()

	c := testClient(t, srv.URL, 3)
	c.cfg.Timeout = time.Second
	resp, err := c.Complete(context.Background(), Request{Messages: []Message{{Role: "user", Content: "x"}}})
	if err != nil {
		t.Fatalf("should have recovered after retries: %v", err)
	}
	if resp.Text != "done" || calls != 3 {
		t.Fatalf("calls=%d text=%q", calls, resp.Text)
	}
}

func TestCompleteAuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"message":"Invalid API Key","type":"invalid_request_error"}}`))
	}))
	defer srv.Close()

	_, err := testClient(t, srv.URL, 0).Complete(context.Background(), Request{Messages: []Message{{Role: "user", Content: "x"}}})
	e, ok := apperr.As(err)
	if !ok || e.Code != apperr.CodeAuth {
		t.Fatalf("want auth error, got %v", err)
	}
}

func TestCompleteRateLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "7")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":{"message":"rate limit"}}`))
	}))
	defer srv.Close()

	_, err := testClient(t, srv.URL, 0).Complete(context.Background(), Request{Messages: []Message{{Role: "user", Content: "x"}}})
	e, ok := apperr.As(err)
	if !ok || e.Code != apperr.CodeRateLimit {
		t.Fatalf("want rate-limit error, got %v", err)
	}
	if !strings.Contains(e.Hint, "7 seconds") {
		t.Errorf("hint should mention Retry-After: %q", e.Hint)
	}
}

func TestCompleteContextLength(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"message":"Please reduce the length of the messages; this model's maximum context length is 8192 tokens"}}`))
	}))
	defer srv.Close()

	_, err := testClient(t, srv.URL, 0).Complete(context.Background(), Request{Messages: []Message{{Role: "user", Content: "x"}}})
	e, ok := apperr.As(err)
	if !ok || e.Code != apperr.CodeContextSize {
		t.Fatalf("want context-size error, got %v", err)
	}
}

func TestCompleteMalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not json at all`))
	}))
	defer srv.Close()

	_, err := testClient(t, srv.URL, 0).Complete(context.Background(), Request{Messages: []Message{{Role: "user", Content: "x"}}})
	e, ok := apperr.As(err)
	if !ok || e.Code != apperr.CodeMalformed {
		t.Fatalf("want malformed error, got %v", err)
	}
}

func TestCompleteEmptyChoices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"choices":[]}`))
	}))
	defer srv.Close()

	_, err := testClient(t, srv.URL, 0).Complete(context.Background(), Request{Messages: []Message{{Role: "user", Content: "x"}}})
	if _, ok := apperr.As(err); !ok {
		t.Fatalf("want apperr, got %v", err)
	}
}

func TestBackoffCap(t *testing.T) {
	if backoff(10) != 8*time.Second {
		t.Fatalf("backoff should cap at 8s, got %v", backoff(10))
	}
}
