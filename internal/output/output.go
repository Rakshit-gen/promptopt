// Package output renders operation results to a terminal or as JSON.
//
// The two modes never mix: JSON output is exactly the pkg/types struct with
// no decorative text, and human output never appears on the JSON path. All
// styling routes through a Theme so --no-color and non-TTY output degrade to
// plain text.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Writer renders results.
type Writer struct {
	out   io.Writer
	err   io.Writer
	theme Theme
	// ShowPrompt controls whether the input prompt text is echoed. Off by
	// default: promptopt does not print user prompts unless asked.
	quiet bool
}

// Options configures a Writer.
type Options struct {
	Out   io.Writer
	Err   io.Writer
	Color bool
	Quiet bool
}

// New builds a Writer.
func New(opts Options) *Writer {
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}
	errw := opts.Err
	if errw == nil {
		errw = os.Stderr
	}
	return &Writer{
		out:   out,
		err:   errw,
		theme: newTheme(opts.Color),
		quiet: opts.Quiet,
	}
}

// JSON writes v as indented JSON followed by a newline.
func (w *Writer) JSON(v any) error {
	enc := json.NewEncoder(w.out)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

// Theme holds the lipgloss styles promptopt uses. When color is off every
// style is the identity function.
type Theme struct {
	Title   lipgloss.Style
	Heading lipgloss.Style
	Dim     lipgloss.Style
	Add     lipgloss.Style
	Warn    lipgloss.Style
	Err     lipgloss.Style
	Info    lipgloss.Style
	Good    lipgloss.Style
	Number  lipgloss.Style
	Box     lipgloss.Style
}

func newTheme(color bool) Theme {
	if color {
		lipgloss.SetColorProfile(termenv.ANSI256)
	} else {
		lipgloss.SetColorProfile(termenv.Ascii)
	}
	accent := lipgloss.Color("39")  // cyan-blue
	dim := lipgloss.Color("244")    // grey
	green := lipgloss.Color("35")   //
	yellow := lipgloss.Color("214") //
	red := lipgloss.Color("203")    //
	return Theme{
		Title:   lipgloss.NewStyle().Bold(true).Foreground(accent),
		Heading: lipgloss.NewStyle().Bold(true),
		Dim:     lipgloss.NewStyle().Foreground(dim),
		Add:     lipgloss.NewStyle().Foreground(green),
		Warn:    lipgloss.NewStyle().Foreground(yellow),
		Err:     lipgloss.NewStyle().Foreground(red).Bold(true),
		Info:    lipgloss.NewStyle().Foreground(dim),
		Good:    lipgloss.NewStyle().Foreground(green),
		Number:  lipgloss.NewStyle().Bold(true),
		Box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dim).
			Padding(0, 1),
	}
}

// helpers shared by the render_*.go files.

func (w *Writer) p(format string, args ...any) {
	fmt.Fprintf(w.out, format, args...)
}

func (w *Writer) nl() { fmt.Fprintln(w.out) }

func (w *Writer) line(s string) { fmt.Fprintln(w.out, s) }

// Errln writes a line to the error stream.
func (w *Writer) Errln(s string) { fmt.Fprintln(w.err, s) }

// scoreBar renders a 0-10 score as a short bar plus the number.
func (w *Writer) scoreBar(label string, score float64) string {
	const width = 20
	filled := int((score/10)*float64(width) + 0.5)
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("·", width-filled)
	col := w.theme.Good
	switch {
	case score < 5:
		col = w.theme.Err
	case score < 7:
		col = w.theme.Warn
	}
	return fmt.Sprintf("  %-16s %s  %s", label, col.Render(bar), w.theme.Number.Render(fmt.Sprintf("%.1f", score)))
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// humanInt formats an int with thousands separators.
func humanInt(n int) string {
	neg := n < 0
	if neg {
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	var out []byte
	for i, c := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, c)
	}
	if neg {
		return "-" + string(out)
	}
	return string(out)
}
