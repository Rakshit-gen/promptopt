// Package groq is a small, isolated client for the Groq chat completions API.
//
// It is deliberately the only place in promptopt that knows about Groq. The
// engine packages depend on the Completer interface, not on this type, so a
// different provider can be added without touching business logic.
//
// Reference: https://console.groq.com/docs/api-reference (OpenAI-compatible
// chat completions at POST https://api.groq.com/openai/v1/chat/completions).
package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rakshit-gen/promptopt/internal/apperr"
)

// DefaultBaseURL is the Groq OpenAI-compatible endpoint root.
const DefaultBaseURL = "https://api.groq.com/openai/v1"

// DefaultModel is a current Groq production model with a large context window
// and reasoning behavior, which suits prompt analysis and rewriting. It can
// be overridden per invocation with --model or PROMPTOPT_MODEL.
const DefaultModel = "openai/gpt-oss-120b"

// Message is a single chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Completer is the interface the engine depends on. Tests substitute a fake.
type Completer interface {
	// Complete sends a chat request and returns the assistant text plus usage.
	Complete(ctx context.Context, req Request) (*Response, error)
}

// Request is a provider-agnostic completion request.
type Request struct {
	Messages    []Message
	Model       string
	Temperature float64
	MaxTokens   int
	// JSON asks the provider to return a single JSON object.
	JSON bool
}

// Usage mirrors the provider usage object. Zero values mean "not reported".
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Response is the result of a completion.
type Response struct {
	Text  string
	Model string
	Usage Usage
}

// Config configures the client.
type Config struct {
	APIKey     string
	BaseURL    string
	Model      string
	Timeout    time.Duration
	MaxRetries int
	// HTTPClient is optional; a sane default is used when nil.
	HTTPClient *http.Client
}

// Client talks to Groq.
type Client struct {
	cfg  Config
	http *http.Client
}

// New builds a client. It returns an actionable error when the API key is
// missing so the CLI can tell the user exactly what to set.
func New(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, apperr.New(
			"No Groq API key configured.",
			"Set it in your shell:\n\n  export GROQ_API_KEY=\"your-key\"\n\n"+
				"Create one at https://console.groq.com/keys",
		).WithCode(apperr.CodeAuth)
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	if cfg.Model == "" {
		cfg.Model = DefaultModel
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second
	}
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 0
	}
	hc := cfg.HTTPClient
	if hc == nil {
		hc = &http.Client{}
	}
	return &Client{cfg: cfg, http: hc}, nil
}

// Model returns the model this client will use by default.
func (c *Client) Model() string { return c.cfg.Model }

// wire types for the Groq request/response.

type wireRequest struct {
	Model               string          `json:"model"`
	Messages            []Message       `json:"messages"`
	Temperature         float64         `json:"temperature"`
	MaxCompletionTokens int             `json:"max_completion_tokens,omitempty"`
	Stream              bool            `json:"stream"`
	ResponseFormat      *responseFormat `json:"response_format,omitempty"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type wireResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

type wireError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Complete implements Completer with timeout and bounded retries for
// transient failures (429 and 5xx). It classifies common failures into
// apperr.Error values with useful hints.
func (c *Client) Complete(ctx context.Context, req Request) (*Response, error) {
	model := req.Model
	if model == "" {
		model = c.cfg.Model
	}
	body := wireRequest{
		Model:       model,
		Messages:    req.Messages,
		Temperature: req.Temperature,
		Stream:      false,
	}
	if req.MaxTokens > 0 {
		body.MaxCompletionTokens = req.MaxTokens
	}
	if req.JSON {
		body.ResponseFormat = &responseFormat{Type: "json_object"}
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, apperr.Wrap(err, "Could not encode the request.", "")
	}

	var lastErr error
	for attempt := 0; attempt <= c.cfg.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := backoff(attempt)
			select {
			case <-ctx.Done():
				return nil, ctxErr(ctx.Err())
			case <-time.After(delay):
			}
		}

		resp, retry, err := c.do(ctx, payload, model)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, lastErr
}

func (c *Client) do(ctx context.Context, payload []byte, model string) (*Response, bool, error) {
	reqCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodPost,
		c.cfg.BaseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, false, apperr.Wrap(err, "Could not build the HTTP request.", "")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	httpReq.Header.Set("User-Agent", "promptopt")

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(reqCtx.Err(), context.DeadlineExceeded) {
			return nil, true, apperr.Wrap(err,
				fmt.Sprintf("Groq request timed out after %s.", c.cfg.Timeout),
				"Increase the timeout with --timeout, or check your connection.")
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil, false, ctxErr(ctx.Err())
		}
		return nil, true, apperr.Wrap(err, "Could not reach the Groq API.",
			"Check your network connection and https://groqstatus.com.")
	}
	defer func() { _ = httpResp.Body.Close() }()

	raw, err := io.ReadAll(io.LimitReader(httpResp.Body, 8<<20))
	if err != nil {
		return nil, true, apperr.Wrap(err, "Could not read the Groq response.", "")
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, retryable(httpResp.StatusCode), classify(httpResp, raw, model)
	}

	var wr wireResponse
	if err := json.Unmarshal(raw, &wr); err != nil {
		return nil, false, apperr.Wrap(err,
			"Groq returned a response that could not be parsed as JSON.",
			"This is usually transient. Retry, and report it if it persists.").
			WithCode(apperr.CodeMalformed)
	}
	if len(wr.Choices) == 0 || strings.TrimSpace(wr.Choices[0].Message.Content) == "" {
		return nil, true, apperr.New(
			"Groq returned an empty completion.",
			"Retry the command. If it keeps happening, try a different --model.").
			WithCode(apperr.CodeMalformed)
	}

	return &Response{
		Text:  wr.Choices[0].Message.Content,
		Model: wr.Model,
		Usage: wr.Usage,
	}, false, nil
}

func classify(resp *http.Response, raw []byte, model string) error {
	var we wireError
	_ = json.Unmarshal(raw, &we)
	msg := strings.TrimSpace(we.Error.Message)
	if msg == "" {
		msg = strings.TrimSpace(string(raw))
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return apperr.New("Groq rejected the API key.",
			"Check GROQ_API_KEY is set to a valid key from\n"+
				"https://console.groq.com/keys").WithCode(apperr.CodeAuth)

	case http.StatusTooManyRequests:
		hint := "Wait and retry."
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(strings.TrimSpace(ra)); err == nil {
				hint = fmt.Sprintf("Retry after about %d seconds (the provider's suggested interval).", secs)
			}
		}
		return apperr.New("Groq rate limit reached.", hint).WithCode(apperr.CodeRateLimit)

	case http.StatusRequestEntityTooLarge, http.StatusBadRequest:
		if isContextLength(msg) {
			return apperr.New(
				fmt.Sprintf("The prompt is longer than the context window of %q.", model),
				"Try `promptopt compress` first, split the prompt, or pick a model\n"+
					"with a larger context window via --model.").WithCode(apperr.CodeContextSize)
		}
		return apperr.New("Groq rejected the request as invalid.",
			detailHint(msg)).WithCode(apperr.CodeProvider)

	case http.StatusNotFound:
		return apperr.New(
			fmt.Sprintf("Groq does not recognize the model %q.", model),
			"List available models at https://console.groq.com/docs/models\n"+
				"and set one with --model or PROMPTOPT_MODEL.").WithCode(apperr.CodeProvider)

	default:
		if resp.StatusCode >= 500 {
			return apperr.New(
				fmt.Sprintf("Groq returned a server error (HTTP %d).", resp.StatusCode),
				"This is on the provider's side. Retry shortly; check https://groqstatus.com.").
				WithCode(apperr.CodeProvider)
		}
		return apperr.New(
			fmt.Sprintf("Groq request failed (HTTP %d).", resp.StatusCode),
			detailHint(msg)).WithCode(apperr.CodeProvider)
	}
}

func detailHint(msg string) string {
	if msg == "" {
		return ""
	}
	if len(msg) > 400 {
		msg = msg[:400] + "…"
	}
	return "Provider message:\n\n  " + msg
}

func isContextLength(msg string) bool {
	m := strings.ToLower(msg)
	return strings.Contains(m, "context") &&
		(strings.Contains(m, "length") || strings.Contains(m, "window") || strings.Contains(m, "maximum"))
}

func retryable(status int) bool {
	return status == http.StatusTooManyRequests || status >= 500
}

func backoff(attempt int) time.Duration {
	// 400ms, 800ms, 1.6s, ... capped at 8s.
	d := time.Duration(math.Pow(2, float64(attempt-1))) * 400 * time.Millisecond
	if d > 8*time.Second {
		d = 8 * time.Second
	}
	return d
}

func ctxErr(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return apperr.New("The operation timed out.", "Increase --timeout and retry.")
	}
	return apperr.New("The operation was canceled.", "")
}
