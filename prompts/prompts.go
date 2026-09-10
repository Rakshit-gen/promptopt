// Package prompts embeds the versioned instruction files that promptopt sends
// to the model for each operation. Keeping them in the binary means the tool
// has no runtime file dependency, and keeping them in version control means
// every change to how an operation behaves is reviewable in a diff.
//
// Each file has YAML front matter (id, version, updated). Load() strips it and
// returns the instruction body.
package prompts

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed optimize/v1.md compress/v1.md expand/v1.md analyze/v1.md transform/v1.md eval/v1.md
var files embed.FS

// Version is the prompt set version reported by the CLI.
const Version = "v1"

// Load returns the instruction body for an operation ("optimize", "compress",
// "expand", "analyze", "transform", "eval"), with front matter removed.
func Load(operation string) (string, error) {
	data, err := files.ReadFile(operation + "/v1.md")
	if err != nil {
		return "", fmt.Errorf("no embedded prompt for operation %q: %w", operation, err)
	}
	return stripFrontMatter(string(data)), nil
}

// MustLoad is Load without the error, for package-level initialization.
func MustLoad(operation string) string {
	s, err := Load(operation)
	if err != nil {
		panic(err)
	}
	return s
}

func stripFrontMatter(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return strings.TrimSpace(s)
	}
	end := strings.Index(s[4:], "\n---\n")
	if end < 0 {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(s[4+end+len("\n---\n"):])
}
