package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
)

// maxPromptBytes caps prompt input so a runaway pipe or a huge file does not
// exhaust memory. 2 MiB is far larger than any real prompt.
const maxPromptBytes = 2 << 20

// promptSource is the assembled prompt plus a label describing where it came
// from ("argument", "stdin", or a file path).
type promptSource struct {
	Text  string
	Label string
}

// resolvePrompt figures out where the prompt text comes from, in this order:
//
//  1. positional arg "-"            -> read stdin
//  2. positional arg is a file path -> read the file
//  3. positional arg (anything else) -> the literal prompt
//  4. no arg, stdin is piped        -> read stdin
//  5. otherwise                     -> usage error
func resolvePrompt(cmdName string, args []string, stdin io.Reader, stdinPiped bool) (promptSource, error) {
	if len(args) > 1 {
		return promptSource{}, apperr.Usagef(
			"expected one prompt (a string or a file path), got %d arguments", len(args))
	}

	if len(args) == 1 {
		arg := args[0]
		if arg == "-" {
			return readAll(stdin, "stdin")
		}
		if looksLikePath(arg) {
			if info, err := os.Stat(arg); err == nil && !info.IsDir() {
				data, err := os.ReadFile(arg)
				if err != nil {
					return promptSource{}, apperr.Wrap(err, "Could not read "+arg+".", "")
				}
				return finalize(string(data), arg)
			}
		}
		return finalize(arg, "argument")
	}

	if stdinPiped {
		return readAll(stdin, "stdin")
	}

	return promptSource{}, apperr.Usagef(
		"no prompt given. Pass a string, a file path, or pipe one in:\n\n"+
			"  promptopt %s \"your prompt\"\n"+
			"  promptopt %s prompt.txt\n"+
			"  cat prompt.txt | promptopt %s", cmdName, cmdName, cmdName)
}

func readAll(r io.Reader, label string) (promptSource, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxPromptBytes+1))
	if err != nil {
		return promptSource{}, apperr.Wrap(err, "Could not read "+label+".", "")
	}
	if len(data) > maxPromptBytes {
		return promptSource{}, apperr.Usagef("prompt from %s is larger than %d bytes", label, maxPromptBytes)
	}
	return finalize(string(data), label)
}

func finalize(raw, label string) (promptSource, error) {
	text := normalize(raw)
	if strings.TrimSpace(text) == "" {
		return promptSource{}, apperr.Usagef("the prompt from %s is empty", label)
	}
	return promptSource{Text: text, Label: label}, nil
}

func normalize(s string) string {
	return strings.TrimRight(strings.ReplaceAll(s, "\r\n", "\n"), "\n \t")
}

// looksLikePath is a cheap heuristic to decide whether to stat the arg. A
// string with a newline, or very long, or containing spaces without a
// separator, is treated as a literal prompt.
func looksLikePath(s string) bool {
	if strings.ContainsAny(s, "\n\r") || len(s) > 400 {
		return false
	}
	if strings.Contains(s, "/") || strings.Contains(s, string(os.PathSeparator)) {
		return true
	}
	return !strings.Contains(s, " ")
}

// stdinIsPipe reports whether stdin is connected to a pipe or file rather than
// an interactive terminal.
func stdinIsPipe() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}

// inputIsPiped decides whether to read a prompt from stdin when no argument is
// given. It returns true when stdin is a real pipe/file, or when a caller
// (a test) has swapped the command's input stream for something other than
// the process stdin.
func inputIsPiped(in io.Reader) bool {
	if f, ok := in.(*os.File); !ok || f != os.Stdin {
		return true
	}
	return stdinIsPipe()
}

// writeOutputFile writes result text to path, creating parent dirs.
func writeOutputFile(path, content string) error {
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return apperr.Wrap(err, fmt.Sprintf("Could not write %s.", path), "")
	}
	return nil
}
