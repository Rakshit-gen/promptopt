// Package cli wires promptopt's commands to Cobra. It is the only package that
// imports Cobra; every command delegates the actual work to internal/engine,
// which has no CLI dependency and is tested on its own.
package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/internal/config"
	"github.com/rakshit-gen/promptopt/internal/engine"
	"github.com/rakshit-gen/promptopt/internal/groq"
	"github.com/rakshit-gen/promptopt/internal/output"
	"github.com/rakshit-gen/promptopt/internal/tokenizer"
	"github.com/rakshit-gen/promptopt/prompts"
	"github.com/spf13/cobra"
)

// Build metadata, set via -ldflags at release time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// globalFlags are shared by every operation command.
type globalFlags struct {
	model       string
	jsonOut     bool
	noColor     bool
	quiet       bool
	outputFile  string
	configPath  string
	timeout     time.Duration
	temperature float64
	verbose     bool
}

func (g *globalFlags) bind(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringVar(&g.model, "model", "", "Groq model to use (overrides config and PROMPTOPT_MODEL)")
	f.BoolVar(&g.jsonOut, "json", false, "emit a stable JSON object instead of formatted text")
	f.BoolVar(&g.noColor, "no-color", false, "disable ANSI color (also respects NO_COLOR)")
	f.BoolVarP(&g.quiet, "quiet", "q", false, "print only the resulting prompt, no report")
	f.StringVarP(&g.outputFile, "output", "o", "", "write the resulting prompt to a file")
	f.StringVar(&g.configPath, "config", "", "path to a config file (default ~/.config/promptopt/config.yaml)")
	f.DurationVar(&g.timeout, "timeout", 0, "per-request timeout (e.g. 45s, 2m)")
	f.Float64Var(&g.temperature, "temperature", -1, "sampling temperature 0-2 (default from config)")
	f.BoolVarP(&g.verbose, "verbose", "v", false, "show underlying errors and diagnostics")
}

// runContext is everything a command handler needs.
type runContext struct {
	cfg    config.Config
	eng    *engine.Engine
	writer *output.Writer
	flags  *globalFlags
	stdin  io.Reader
}

// setup resolves config, builds the Groq client, engine, and output writer.
func (g *globalFlags) setup(cmd *cobra.Command) (*runContext, error) {
	cfg, err := config.Load(g.configPath)
	if err != nil {
		return nil, apperr.Wrap(err, "Configuration error.",
			"Check "+config.DefaultPath()+" or pass --config.").WithCode(apperr.CodeConfig)
	}

	if g.model != "" {
		cfg.Model = g.model
	}
	if g.timeout > 0 {
		cfg.Timeout = g.timeout
	}
	if g.temperature >= 0 {
		cfg.Temperature = g.temperature
	}
	if cmd.Flags().Changed("json") && g.jsonOut {
		cfg.Output = "json"
	}

	color := resolveColor(cfg.Color, g.noColor, cmd.OutOrStdout())

	client, err := groq.New(groq.Config{
		APIKey:     cfg.APIKey,
		BaseURL:    cfg.BaseURL,
		Model:      cfg.Model,
		Timeout:    cfg.Timeout,
		MaxRetries: cfg.MaxRetries,
	})
	if err != nil {
		return nil, err
	}

	eng, err := engine.New(engine.Options{
		Client:      client,
		Tokenizer:   tokenizer.NewHeuristic(),
		Model:       cfg.Model,
		Temperature: cfg.Temperature,
	})
	if err != nil {
		return nil, err
	}

	w := output.New(output.Options{
		Out:   cmd.OutOrStdout(),
		Err:   cmd.ErrOrStderr(),
		Color: color,
		Quiet: g.quiet,
	})

	return &runContext{cfg: cfg, eng: eng, writer: w, flags: g, stdin: cmd.InOrStdin()}, nil
}

func resolveColor(mode string, noColorFlag bool, out io.Writer) bool {
	if noColorFlag || mode == "never" || os.Getenv("NO_COLOR") != "" {
		return false
	}
	if mode == "always" {
		return true
	}
	// auto: color only when stdout is a terminal.
	if f, ok := out.(*os.File); ok {
		info, err := f.Stat()
		if err == nil && (info.Mode()&os.ModeCharDevice) != 0 {
			return true
		}
	}
	return false
}

// Execute is the CLI entry point.
func Execute() int {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		return reportError(root.ErrOrStderr(), err, rootVerbose(root))
	}
	return 0
}

func rootVerbose(cmd *cobra.Command) bool {
	// Walk executed command tree for a --verbose flag that was set.
	for _, c := range append([]*cobra.Command{cmd}, cmd.Commands()...) {
		if f := c.Flags().Lookup("verbose"); f != nil && f.Changed {
			return true
		}
	}
	return false
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promptopt",
		Short: "Understand your prompts, then make them better.",
		Long: `promptopt is a command-line tool for working on the prompts you send to
language models. It inspects a prompt, tells you what is weak or wasteful,
and rewrites it while keeping the job the prompt is meant to do.

You describe the work. promptopt handles the prompt-engineering mechanics:
clarifying objectives, removing redundancy, fixing instruction order,
pinning down output format, and so on.

Six operations:

  optimize    inspect a prompt and rewrite it to be clearer and tighter
  compress    cut tokens without changing what the prompt makes a model do
  expand      turn an underspecified prompt into a detailed one
  analyze     score a prompt and report actionable findings, like a linter
  transform   convert a prompt to markdown, XML, JSON, a system prompt, a
              template, or an agent brief
  eval        assess a prompt and probe it with generated test cases

Input comes from an argument, a file, or stdin. Every command supports
--json for a stable machine-readable object, so results compose in a shell
pipeline.

  promptopt optimize "build a REST API for payments"
  promptopt compress system-prompt.md --target 30
  cat prompt.txt | promptopt analyze --json

Configuration is read from environment variables and an optional file at
~/.config/promptopt/config.yaml. The only network call promptopt makes is to
Groq, using GROQ_API_KEY. It does not send telemetry and does not store your
prompts anywhere.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       fmt.Sprintf("%s (commit %s, built %s, prompts %s)", Version, Commit, Date, prompts.Version),
	}

	cmd.SetVersionTemplate("promptopt {{.Version}}\n")
	cmd.AddCommand(
		newOptimizeCmd(),
		newCompressCmd(),
		newExpandCmd(),
		newAnalyzeCmd(),
		newTransformCmd(),
		newEvalCmd(),
		newConfigCmd(),
	)
	return cmd
}

func baseContext() context.Context {
	return context.Background()
}
