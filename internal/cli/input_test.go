package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rakshit-gen/promptopt/internal/apperr"
)

func TestResolvePromptFromArgument(t *testing.T) {
	src, err := resolvePrompt("optimize", []string{"build a REST API for payments"}, strings.NewReader(""), false)
	if err != nil {
		t.Fatal(err)
	}
	if src.Text != "build a REST API for payments" || src.Label != "argument" {
		t.Fatalf("got %+v", src)
	}
}

func TestResolvePromptFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "prompt.txt")
	os.WriteFile(path, []byte("system prompt body\n\n"), 0o644)

	src, err := resolvePrompt("compress", []string{path}, strings.NewReader(""), false)
	if err != nil {
		t.Fatal(err)
	}
	if src.Text != "system prompt body" {
		t.Fatalf("file body not trimmed: %q", src.Text)
	}
	if src.Label != path {
		t.Fatalf("label = %q", src.Label)
	}
}

func TestResolvePromptFromStdin(t *testing.T) {
	src, err := resolvePrompt("analyze", nil, strings.NewReader("piped prompt\n"), true)
	if err != nil {
		t.Fatal(err)
	}
	if src.Text != "piped prompt" || src.Label != "stdin" {
		t.Fatalf("got %+v", src)
	}
}

func TestResolvePromptDashMeansStdin(t *testing.T) {
	src, err := resolvePrompt("analyze", []string{"-"}, strings.NewReader("dash prompt"), false)
	if err != nil {
		t.Fatal(err)
	}
	if src.Text != "dash prompt" {
		t.Fatalf("got %+v", src)
	}
}

func TestResolvePromptNoInput(t *testing.T) {
	_, err := resolvePrompt("optimize", nil, strings.NewReader(""), false)
	e, ok := apperr.As(err)
	if !ok || e.Code != apperr.CodeUsage {
		t.Fatalf("want usage error, got %v", err)
	}
	if !strings.Contains(e.Message, "promptopt optimize") {
		t.Errorf("usage message should show the command name: %q", e.Message)
	}
}

func TestResolvePromptEmptyIsError(t *testing.T) {
	_, err := resolvePrompt("optimize", []string{"   "}, strings.NewReader(""), false)
	if _, ok := apperr.As(err); !ok {
		t.Fatalf("blank argument should be a usage error, got %v", err)
	}
}

func TestResolvePromptTooManyArgs(t *testing.T) {
	_, err := resolvePrompt("optimize", []string{"a", "b"}, strings.NewReader(""), false)
	if _, ok := apperr.As(err); !ok {
		t.Fatalf("want usage error, got %v", err)
	}
}

func TestResolvePromptLiteralThatLooksLikeMissingFile(t *testing.T) {
	// A short no-space string that is not a real file should be treated as a
	// literal prompt, not a "file not found" error.
	src, err := resolvePrompt("optimize", []string{"summarize"}, strings.NewReader(""), false)
	if err != nil {
		t.Fatal(err)
	}
	if src.Text != "summarize" {
		t.Fatalf("got %+v", src)
	}
}

func TestLooksLikePath(t *testing.T) {
	yes := []string{"prompt.txt", "./x", "dir/file.md", "noextension"}
	no := []string{"build an api for me", "line one\nline two"}
	for _, s := range yes {
		if !looksLikePath(s) {
			t.Errorf("looksLikePath(%q) = false, want true", s)
		}
	}
	for _, s := range no {
		if looksLikePath(s) {
			t.Errorf("looksLikePath(%q) = true, want false", s)
		}
	}
}

func TestNormalize(t *testing.T) {
	if got := normalize("a\r\nb\r\n\n"); got != "a\nb" {
		t.Fatalf("normalize = %q", got)
	}
}
