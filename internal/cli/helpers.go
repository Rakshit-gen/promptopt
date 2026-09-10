package cli

import (
	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/internal/engine"
	"github.com/rakshit-gen/promptopt/internal/output"
)

// outputWriter is a local alias so command code reads cleanly.
type outputWriter = output.Writer

func apperrUsage(msg string) error { return apperr.Usagef("%s", msg) }

func depthValue(s string) engine.Depth {
	if s == "" {
		return engine.DepthDetailed
	}
	return engine.Depth(s)
}
