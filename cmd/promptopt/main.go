// Command promptopt is an AI-powered prompt engineering CLI backed by Groq.
// It analyzes, transforms, optimizes, and evaluates the prompts you send to
// language models.
package main

import (
	"os"

	"github.com/rakshit-gen/promptopt/internal/cli"
)

// These are overridden at build time via -ldflags. See the Makefile and
// .goreleaser.yaml.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.Version = version
	cli.Commit = commit
	cli.Date = date
	os.Exit(cli.Execute())
}
