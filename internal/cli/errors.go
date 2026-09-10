package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
)

// reportError prints a failure the way a developer wants to read it: one clear
// line, an optional hint with the exact next command, and the underlying
// cause only when --verbose is set. It returns the process exit code.
func reportError(w io.Writer, err error, verbose bool) int {
	e, ok := apperr.As(err)
	if !ok {
		// Unexpected error: still don't dump a stack for normal users.
		fmt.Fprintln(w, "promptopt: "+err.Error())
		if verbose {
			fmt.Fprintf(w, "\n%+v\n", err)
		}
		return apperr.CodeGeneric
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "  "+e.Message)
	if strings.TrimSpace(e.Hint) != "" {
		fmt.Fprintln(w)
		for _, ln := range strings.Split(e.Hint, "\n") {
			fmt.Fprintln(w, "  "+ln)
		}
	}
	if verbose && e.Cause != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  underlying error:")
		fmt.Fprintln(w, "  "+e.Cause.Error())
	}
	fmt.Fprintln(w)

	if e.Code == 0 {
		return apperr.CodeGeneric
	}
	return e.Code
}
