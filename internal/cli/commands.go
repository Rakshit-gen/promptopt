package cli

import (
	"strings"

	"github.com/rakshit-gen/promptopt/internal/engine"
	"github.com/rakshit-gen/promptopt/pkg/types"
	"github.com/spf13/cobra"
)

// operate resolves the prompt, builds the run context, calls fn, and hands the
// typed result to output. Every operation command is a thin wrapper around
// this: describe the command, add its own flags, call the right engine method.
func operate(g *globalFlags, cmd *cobra.Command, args []string, fn func(rc *runContext, prompt string) (any, string, error)) error {
	name := cmd.Name()
	src, err := resolvePrompt(name, args, cmd.InOrStdin(), inputIsPiped(cmd.InOrStdin()))
	if err != nil {
		return err
	}
	rc, err := g.setup(cmd)
	if err != nil {
		return err
	}
	value, plain, err := fn(rc, src.Text)
	if err != nil {
		return err
	}

	if g.outputFile != "" && plain != "" {
		if err := writeOutputFile(g.outputFile, plain); err != nil {
			return err
		}
		rc.writer.Errln("wrote " + g.outputFile)
	}

	if rc.cfg.Output == "json" {
		return rc.writer.JSON(value)
	}
	if g.quiet && plain != "" {
		rc.writer.PlainResult(plain)
		return nil
	}
	renderHuman(rc.writer, value)
	return nil
}

func renderHuman(w *outputWriter, value any) {
	switch v := value.(type) {
	case *types.OptimizeResult:
		w.Optimize(v)
	case *types.CompressResult:
		w.Compress(v)
	case *types.ExpandResult:
		w.Expand(v)
	case *types.AnalyzeResult:
		w.Analyze(v)
	case *types.TransformResult:
		w.Transform(v)
	case *types.EvalResult:
		w.Eval(v)
	}
}

func newOptimizeCmd() *cobra.Command {
	g := &globalFlags{}
	var goal string
	cmd := &cobra.Command{
		Use:   "optimize [prompt|file|-]",
		Short: "Inspect a prompt and rewrite it to be clearer and tighter",
		Long: `optimize reads a prompt, checks it for the problems that make language
models unreliable — unclear objective, missing context, ambiguous wording,
repeated or conflicting instructions, bad ordering, no stated output format,
dead weight — and returns a rewritten version.

It preserves intent. optimize will not add domain requirements the prompt
did not state; that is what expand is for. If the prompt is already good it
makes small changes and says so.

The report shows token counts before and after, a summary of what changed,
a quick quality read on the result, and any warnings (for example, intent it
had to guess).`,
		Example: `  promptopt optimize "build a REST API for payments"
  promptopt optimize prompt.txt
  cat prompt.txt | promptopt optimize --json
  promptopt optimize prompt.md -o optimized.md`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return operate(g, cmd, args, func(rc *runContext, prompt string) (any, string, error) {
				res, err := rc.eng.Optimize(baseContext(), engine.OptimizeInput{Prompt: prompt, Instruction: goal})
				if err != nil {
					return nil, "", err
				}
				return res, res.Result, nil
			})
		},
	}
	g.bind(cmd)
	cmd.Flags().StringVar(&goal, "goal", "", "extra guidance about what the prompt should achieve")
	return cmd
}

func newCompressCmd() *cobra.Command {
	g := &globalFlags{}
	var (
		target           int
		aggressive       bool
		preserveBehavior bool
	)
	cmd := &cobra.Command{
		Use:   "compress [prompt|file|-]",
		Short: "Cut tokens without changing what the prompt makes a model do",
		Long: `compress is useful when a prompt has accumulated instructions, examples, and
context over time and is now paying for tokens it does not need. It removes
repetition, verbose phrasing, duplicate examples, and prose that does not
affect the output.

The distinction that matters is between removing words and removing
behavior. A shorter prompt that changes what the model does is a failed
compression, so the report includes a semantic-preservation estimate and a
list of any behavior risks the model flagged.`,
		Example: `  promptopt compress system-prompt.md
  promptopt compress system.md --target 30
  promptopt compress prompt.txt --aggressive
  promptopt compress prompt.txt --preserve-behavior --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return operate(g, cmd, args, func(rc *runContext, prompt string) (any, string, error) {
				res, err := rc.eng.Compress(baseContext(), engine.CompressInput{
					Prompt:           prompt,
					TargetPercent:    target,
					Aggressive:       aggressive,
					PreserveBehavior: preserveBehavior,
				})
				if err != nil {
					return nil, "", err
				}
				return res, res.Result, nil
			})
		},
	}
	g.bind(cmd)
	f := cmd.Flags()
	f.IntVar(&target, "target", 0, "target reduction percentage (0-95); default is as much as is safe")
	f.BoolVar(&aggressive, "aggressive", false, "allow dropping marginal context and collapsing examples")
	f.BoolVar(&preserveBehavior, "preserve-behavior", false, "refuse any cut not confidently safe, even if it misses the target")
	return cmd
}

func newExpandCmd() *cobra.Command {
	g := &globalFlags{}
	var depth string
	cmd := &cobra.Command{
		Use:   "expand [prompt|file|-]",
		Short: "Turn an underspecified prompt into a detailed one",
		Long: `expand takes a short prompt and fills in the parts a model would otherwise
have to guess: objective, context, constraints, inputs, outputs and their
format, edge cases, failure handling, and evaluation criteria.

It separates what you asked for from what it added. Anything it had to
assume is labeled as an assumption, both in the report and inside the
expanded prompt, so you can correct it. expand does not invent facts,
numbers, or domain rules.`,
		Example: `  promptopt expand "build a payment API"
  promptopt expand "build a payment API" --depth production
  promptopt expand idea.txt --depth concise
  cat idea.txt | promptopt expand --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if depth != "" && !engine.ValidDepth(depth) {
				return usageErr(cmd, "--depth must be concise, detailed, or production")
			}
			return operate(g, cmd, args, func(rc *runContext, prompt string) (any, string, error) {
				res, err := rc.eng.Expand(baseContext(), engine.ExpandInput{
					Prompt: prompt,
					Depth:  depthValue(depth),
				})
				if err != nil {
					return nil, "", err
				}
				return res, res.Result, nil
			})
		},
	}
	g.bind(cmd)
	cmd.Flags().StringVar(&depth, "depth", "detailed", "how far to expand: concise, detailed, or production")
	return cmd
}

func newAnalyzeCmd() *cobra.Command {
	g := &globalFlags{}
	cmd := &cobra.Command{
		Use:   "analyze [prompt|file|-]",
		Short: "Score a prompt and report actionable findings, like a linter",
		Long: `analyze reads a prompt and reports what is wrong with it without changing
it. It scores clarity, specificity, completeness, consistency, token
efficiency, and robustness, then lists findings with a severity:

  ERROR    likely to cause wrong or unsafe output
  WARNING  likely to degrade quality
  INFO     worth knowing, low impact

Each finding has a stable ID (P003, P008, ...) and says what to do about it.
Use analyze in CI to catch regressions in a prompt before they ship.`,
		Example: `  promptopt analyze prompt.txt
  cat prompt.txt | promptopt analyze
  promptopt analyze prompt.txt --json | jq '.scores.overall'`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return operate(g, cmd, args, func(rc *runContext, prompt string) (any, string, error) {
				res, err := rc.eng.Analyze(baseContext(), engine.AnalyzeInput{Prompt: prompt})
				if err != nil {
					return nil, "", err
				}
				return res, "", nil
			})
		},
	}
	g.bind(cmd)
	return cmd
}

func newTransformCmd() *cobra.Command {
	g := &globalFlags{}
	var to string
	cmd := &cobra.Command{
		Use:   "transform [prompt|file|-]",
		Short: "Convert a prompt to a different structure or representation",
		Long: `transform changes how a prompt is expressed without changing what it asks
for. Targets:

  markdown   clean headings, lists, and a fenced output section
  xml        semantic tags (<task>, <context>, <constraints>, ...)
  json       a structured object with role, task, constraints, and so on
  system     a standing system prompt in the second person
  template   {{variables}} in place of values that change between runs
  agent      a goal, available context, and explicit stopping conditions

For template, transform parameterizes the values that would realistically
vary — inputs, target languages, file names — not every noun.`,
		Example: `  promptopt transform prompt.txt --to xml
  promptopt transform prompt.txt --to template
  promptopt transform prompt.txt --to system -o system.txt
  cat prompt.txt | promptopt transform --to json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if to == "" {
				return usageErr(cmd, "transform needs a target: --to "+strings.Join(engine.Targets(), "|"))
			}
			if !engine.ValidTarget(to) {
				return usageErr(cmd, "--to must be one of: "+strings.Join(engine.Targets(), ", "))
			}
			return operate(g, cmd, args, func(rc *runContext, prompt string) (any, string, error) {
				res, err := rc.eng.Transform(baseContext(), engine.TransformInput{
					Prompt: prompt,
					Target: engine.Target(to),
				})
				if err != nil {
					return nil, "", err
				}
				return res, res.Result, nil
			})
		},
	}
	g.bind(cmd)
	cmd.Flags().StringVar(&to, "to", "", "target representation: "+strings.Join(engine.Targets(), ", "))
	_ = cmd.RegisterFlagCompletionFunc("to", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return engine.Targets(), cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}

func newEvalCmd() *cobra.Command {
	g := &globalFlags{}
	cmd := &cobra.Command{
		Use:   "eval [prompt|file|-]",
		Short: "Assess a prompt and probe it with generated test cases",
		Long: `eval gives you a read on whether a prompt is ready to depend on. It scores
clarity, consistency, robustness, and output control, then generates test
cases that probe the prompt's weak points — a normal input, an ambiguous
one, an out-of-scope or adversarial one, and an edge case — and reports how
the prompt would handle each.

This is a foundation, not a full benchmark platform. The result shape is
designed so datasets, assertions, and model comparisons can be added later
without breaking the JSON contract.`,
		Example: `  promptopt eval prompt.txt
  cat prompt.txt | promptopt eval
  promptopt eval prompt.txt --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return operate(g, cmd, args, func(rc *runContext, prompt string) (any, string, error) {
				res, err := rc.eng.Eval(baseContext(), engine.EvalInput{Prompt: prompt})
				if err != nil {
					return nil, "", err
				}
				return res, "", nil
			})
		},
	}
	g.bind(cmd)
	return cmd
}

func usageErr(cmd *cobra.Command, msg string) error {
	cmd.SilenceUsage = false
	return apperrUsage(msg)
}
