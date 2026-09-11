package output

import (
	"fmt"
	"strings"

	"github.com/rakshit-gen/promptopt/pkg/types"
)

// header prints the "promptopt" line and a one-line operation label.
func (w *Writer) header(label string) {
	w.nl()
	w.line("  " + w.theme.Title.Render("promptopt") + w.theme.Dim.Render("  "+label))
	w.nl()
}

func (w *Writer) tokenBlock(t types.TokenReport) {
	w.line("  " + w.theme.Heading.Render("tokens") + w.theme.Dim.Render("  ("+t.Source+")"))
	arrow := fmt.Sprintf("  %s → %s", humanInt(t.Before), humanInt(t.After))
	w.line(w.theme.Number.Render(arrow))
	switch {
	case t.ReductionPercent > 0.05:
		w.line("  " + w.theme.Good.Render(fmt.Sprintf("%.1f%% reduction", t.ReductionPercent)))
	case t.ReductionPercent < -0.05:
		w.line("  " + w.theme.Warn.Render(fmt.Sprintf("%.1f%% larger", -t.ReductionPercent)))
	default:
		w.line("  " + w.theme.Dim.Render("no change"))
	}
	w.nl()
}

func (w *Writer) changeList(title string, changes []string) {
	if len(changes) == 0 {
		return
	}
	w.line("  " + w.theme.Heading.Render(title))
	for _, c := range changes {
		w.line("  " + w.theme.Add.Render("+ ") + c)
	}
	w.nl()
}

func (w *Writer) warnings(ws []string) {
	if len(ws) == 0 {
		return
	}
	w.line("  " + w.theme.Warn.Render("warnings"))
	for _, x := range ws {
		w.line("  " + w.theme.Warn.Render("! ") + x)
	}
	w.nl()
}

// promptBox prints the title as a heading and the body as plain text, with
// nothing else on those lines, so selecting the body and pasting it
// elsewhere reproduces it exactly. It used to sit inside a bordered box,
// which looked nicer on screen but meant every copy-paste picked up the
// border characters too, and the border broke across lines wider than the
// terminal.
func (w *Writer) promptBox(title, body string) {
	w.line("  " + w.theme.Heading.Render(title))
	w.nl()
	w.line(strings.TrimRight(body, "\n"))
	w.nl()
}

// Optimize renders an OptimizeResult.
func (w *Writer) Optimize(r *types.OptimizeResult) {
	w.header("optimize")
	w.tokenBlock(r.Tokens)
	w.changeList("changes", r.Changes)
	if r.Analysis.Overall > 0 {
		w.line("  " + w.theme.Heading.Render("quality") + w.theme.Dim.Render("  (optimized prompt)"))
		w.line(w.scoreBar("overall", r.Analysis.Overall))
		w.nl()
	}
	w.warnings(r.Warnings)
	if !w.quiet {
		w.promptBox("optimized prompt", r.Result)
	} else {
		w.line(r.Result)
	}
}

// Compress renders a CompressResult.
func (w *Writer) Compress(r *types.CompressResult) {
	w.header("compress")
	w.tokenBlock(r.Tokens)
	w.changeList("changes", r.Changes)
	w.line("  " + w.theme.Heading.Render("semantic preservation"))
	w.line(w.scoreBar("preserved", r.SemanticPreservation))
	w.nl()
	if len(r.BehaviorRisks) > 0 {
		w.line("  " + w.theme.Warn.Render("behavior risks"))
		for _, risk := range r.BehaviorRisks {
			w.line("  " + w.theme.Warn.Render("! ") + risk)
		}
		w.nl()
	}
	w.warnings(r.Warnings)
	if !w.quiet {
		w.promptBox("compressed prompt", r.Result)
	} else {
		w.line(r.Result)
	}
}

// Expand renders an ExpandResult.
func (w *Writer) Expand(r *types.ExpandResult) {
	w.header("expand · " + r.Depth)
	w.tokenBlock(r.Tokens)
	w.changeList("added", r.Added)
	if len(r.Assumptions) > 0 {
		w.line("  " + w.theme.Heading.Render("assumptions") + w.theme.Dim.Render("  (correct these if wrong)"))
		for _, a := range r.Assumptions {
			label := a.Field
			if label == "" {
				label = "assumption"
			}
			w.line("  " + w.theme.Warn.Render("~ ") + w.theme.Heading.Render(label) + ": " + a.Value)
			if a.Note != "" {
				w.line("    " + w.theme.Dim.Render(a.Note))
			}
		}
		w.nl()
	}
	w.warnings(r.Warnings)
	if !w.quiet {
		w.promptBox("expanded prompt", r.Result)
	} else {
		w.line(r.Result)
	}
}

// Analyze renders an AnalyzeResult.
func (w *Writer) Analyze(r *types.AnalyzeResult) {
	w.header("analyze")
	w.line("  " + w.theme.Heading.Render("overall score") + "  " +
		w.theme.Number.Render(fmt.Sprintf("%.1f/10", r.Scores.Overall)))
	w.nl()
	s := r.Scores
	for _, row := range []struct {
		name string
		val  float64
	}{
		{"clarity", s.Clarity},
		{"specificity", s.Specificity},
		{"completeness", s.Completeness},
		{"consistency", s.Consistency},
		{"efficiency", s.Efficiency},
		{"robustness", s.Robustness},
	} {
		if row.val > 0 {
			w.line(w.scoreBar(row.name, row.val))
		}
	}
	w.nl()

	if len(r.Findings) == 0 {
		w.line("  " + w.theme.Good.Render("no findings"))
		w.nl()
	} else {
		counts := map[types.Severity]int{}
		for _, f := range r.Findings {
			counts[f.Severity]++
		}
		w.line("  " + w.theme.Heading.Render("findings") + w.theme.Dim.Render(fmt.Sprintf(
			"  %d error, %d warning, %d info",
			counts[types.SeverityError], counts[types.SeverityWarning], counts[types.SeverityInfo])))
		w.nl()
		for _, f := range r.Findings {
			w.line("  " + w.severityTag(f.Severity) + " " + w.theme.Dim.Render(f.ID) + "  " + w.theme.Heading.Render(f.Title))
			if f.Detail != "" {
				for _, ln := range wrapText(f.Detail, 72) {
					w.line("      " + ln)
				}
			}
			w.nl()
		}
	}
	if r.Summary != "" {
		w.line("  " + w.theme.Heading.Render("summary"))
		for _, ln := range wrapText(r.Summary, 72) {
			w.line("  " + ln)
		}
		w.nl()
	}
}

// Transform renders a TransformResult.
func (w *Writer) Transform(r *types.TransformResult) {
	w.header("transform → " + r.Target)
	if len(r.Notes) > 0 {
		w.line("  " + w.theme.Heading.Render("notes"))
		for _, n := range r.Notes {
			w.line("  " + w.theme.Dim.Render("· ") + n)
		}
		w.nl()
	}
	if len(r.Variables) > 0 {
		w.line("  " + w.theme.Heading.Render("variables") + w.theme.Dim.Render(fmt.Sprintf("  (%d)", len(r.Variables))))
		w.line("  " + strings.Join(prefixEach("{{", r.Variables, "}}"), "  "))
		w.nl()
	}
	w.warnings(r.Warnings)
	if !w.quiet {
		w.promptBox(r.Target, r.Result)
	} else {
		w.line(r.Result)
	}
}

// Eval renders an EvalResult.
func (w *Writer) Eval(r *types.EvalResult) {
	w.header("eval")
	w.line("  " + w.theme.Heading.Render("overall") + "  " +
		w.theme.Number.Render(fmt.Sprintf("%.1f/10", r.Scores.Overall)))
	w.nl()
	s := r.Scores
	for _, row := range []struct {
		name string
		val  float64
	}{
		{"clarity", s.Clarity},
		{"consistency", s.Consistency},
		{"robustness", s.Robustness},
		{"output control", s.OutputControl},
	} {
		if row.val > 0 {
			w.line(w.scoreBar(row.name, row.val))
		}
	}
	w.nl()

	if len(r.TestCases) > 0 {
		pass := 0
		for _, tc := range r.TestCases {
			if tc.Pass != nil && *tc.Pass {
				pass++
			}
		}
		w.line("  " + w.theme.Heading.Render("test cases") + w.theme.Dim.Render(fmt.Sprintf("  %d/%d handled well", pass, len(r.TestCases))))
		w.nl()
		for _, tc := range r.TestCases {
			mark := w.theme.Warn.Render("~")
			if tc.Pass != nil {
				if *tc.Pass {
					mark = w.theme.Good.Render("✓")
				} else {
					mark = w.theme.Err.Render("✗")
				}
			}
			w.line("  " + mark + " " + w.theme.Heading.Render(tc.Name))
			if tc.Input != "" {
				w.line("    " + w.theme.Dim.Render("input: ") + truncate(tc.Input, 100))
			}
			if tc.Assessment != "" {
				for _, ln := range wrapText(tc.Assessment, 72) {
					w.line("    " + ln)
				}
			}
			w.nl()
		}
	}

	if len(r.Weaknesses) > 0 {
		w.line("  " + w.theme.Heading.Render("weaknesses"))
		for _, x := range r.Weaknesses {
			w.line("  " + w.theme.Warn.Render("! ") + x)
		}
		w.nl()
	}
	if r.Summary != "" {
		w.line("  " + w.theme.Heading.Render("summary"))
		for _, ln := range wrapText(r.Summary, 72) {
			w.line("  " + ln)
		}
		w.nl()
	}
}

func (w *Writer) severityTag(s types.Severity) string {
	switch s {
	case types.SeverityError:
		return w.theme.Err.Render("ERROR  ")
	case types.SeverityWarning:
		return w.theme.Warn.Render("WARNING")
	default:
		return w.theme.Info.Render("INFO   ")
	}
}

// PlainResult writes only the transformed/optimized text, for --quiet piping.
func (w *Writer) PlainResult(s string) { w.line(strings.TrimRight(s, "\n")) }

func prefixEach(prefix string, items []string, suffix string) []string {
	out := make([]string, len(items))
	for i, s := range items {
		out[i] = prefix + s + suffix
	}
	return out
}

func truncate(s string, n int) string {
	s = strings.ReplaceAll(strings.TrimSpace(s), "\n", " ")
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

// wrapText wraps s to width columns on word boundaries.
func wrapText(s string, width int) []string {
	var lines []string
	for _, para := range strings.Split(strings.TrimSpace(s), "\n") {
		words := strings.Fields(para)
		if len(words) == 0 {
			continue
		}
		cur := words[0]
		for _, wd := range words[1:] {
			if len(cur)+1+len(wd) > width {
				lines = append(lines, cur)
				cur = wd
			} else {
				cur += " " + wd
			}
		}
		lines = append(lines, cur)
	}
	return lines
}
