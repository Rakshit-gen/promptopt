// Package apperr defines the error type promptopt uses to turn failures into
// useful terminal messages. An Error carries a short message, an optional
// longer hint with a suggested fix, and an exit code. The CLI's top-level
// handler prints the hint for users and only shows stack-like detail in
// verbose mode.
package apperr

import (
	"errors"
	"fmt"
)

// Error is a user-facing error with an actionable hint.
type Error struct {
	// Message is a single line describing what went wrong.
	Message string
	// Hint is optional multi-line guidance, e.g. the command to run next.
	Hint string
	// Code is the process exit code (1 for generic, others for categories).
	Code int
	// Cause is the wrapped underlying error, shown only in verbose mode.
	Cause error
}

func (e *Error) Error() string { return e.Message }

func (e *Error) Unwrap() error { return e.Cause }

// New builds an Error with exit code 1.
func New(message, hint string) *Error {
	return &Error{Message: message, Hint: hint, Code: 1}
}

// Wrap attaches a cause to a new Error.
func Wrap(cause error, message, hint string) *Error {
	return &Error{Message: message, Hint: hint, Code: 1, Cause: cause}
}

// WithCode sets a specific exit code and returns the error for chaining.
func (e *Error) WithCode(code int) *Error {
	e.Code = code
	return e
}

// As reports whether err is an *Error and returns it.
func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// Common exit codes. These are deliberately small and stable so scripts can
// branch on them.
const (
	CodeGeneric     = 1
	CodeUsage       = 2
	CodeConfig      = 3
	CodeAuth        = 4
	CodeRateLimit   = 5
	CodeContextSize = 6
	CodeProvider    = 7
	CodeMalformed   = 8
)

// Usagef builds a usage error (exit code 2).
func Usagef(format string, args ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, args...), Code: CodeUsage}
}
