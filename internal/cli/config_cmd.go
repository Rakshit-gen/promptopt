package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rakshit-gen/promptopt/internal/apperr"
	"github.com/rakshit-gen/promptopt/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect and initialize promptopt configuration",
		Long: `config helps you see what promptopt would use for a run and set up a config
file. The API key is never printed, only whether one was found and where.`,
	}
	cmd.AddCommand(newConfigShowCmd(), newConfigPathCmd(), newConfigInitCmd())
	return cmd
}

func newConfigShowCmd() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Print the effective configuration (secrets redacted)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return apperr.Wrap(err, "Configuration error.", "").WithCode(apperr.CodeConfig)
			}
			out := cmd.OutOrStdout()
			src := cfg.Source()
			if src == "" {
				src = "(no config file; defaults + environment)"
			}
			keyState := "not set"
			if cfg.APIKey != "" {
				keyState = "set (" + maskKey(cfg.APIKey) + ")"
			}
			fmt.Fprintf(out, "config file    %s\n", src)
			fmt.Fprintf(out, "api key        %s\n", keyState)
			fmt.Fprintf(out, "model          %s\n", cfg.Model)
			fmt.Fprintf(out, "base url       %s\n", cfg.BaseURL)
			fmt.Fprintf(out, "output         %s\n", cfg.Output)
			fmt.Fprintf(out, "color          %s\n", cfg.Color)
			fmt.Fprintf(out, "temperature    %.2f\n", cfg.Temperature)
			fmt.Fprintf(out, "timeout        %s\n", cfg.Timeout)
			fmt.Fprintf(out, "max retries    %d\n", cfg.MaxRetries)
			return nil
		},
	}
	cmd.Flags().StringVar(&configPath, "config", "", "path to a config file")
	return cmd
}

func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the default config file path",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), config.DefaultPath())
			return nil
		},
	}
}

func newConfigInitCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Write a commented starter config file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			path := config.DefaultPath()
			if _, err := os.Stat(path); err == nil && !force {
				return apperr.New(path+" already exists.",
					"Pass --force to overwrite it.").WithCode(apperr.CodeConfig)
			}
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return apperr.Wrap(err, "Could not create the config directory.", "")
			}
			if err := os.WriteFile(path, []byte(starterConfig), 0o644); err != nil {
				return apperr.Wrap(err, "Could not write the config file.", "")
			}
			fmt.Fprintln(cmd.OutOrStdout(), "wrote "+path)
			fmt.Fprintln(cmd.OutOrStdout(), "set your Groq API key with:  export GROQ_API_KEY=\"your-key\"")
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing config file")
	return cmd
}

func maskKey(k string) string {
	k = strings.TrimSpace(k)
	if len(k) <= 8 {
		return "****"
	}
	return k[:4] + "…" + k[len(k)-2:]
}

const starterConfig = `# promptopt configuration
# Location: ~/.config/promptopt/config.yaml (honors XDG_CONFIG_HOME)
#
# The API key is NOT stored here. promptopt reads it from the environment:
#   export GROQ_API_KEY="your-key"
# or from a file you point at with api_key_file below.

# Groq model. See https://console.groq.com/docs/models
model: openai/gpt-oss-120b

# Groq API base URL. Change only if you proxy the API.
base_url: https://api.groq.com/openai/v1

# Default output mode: text or json
output: text

# Color: auto (TTY only), always, or never
color: auto

# Sampling temperature, 0-2. Lower is more deterministic.
temperature: 0.3

# Per-request timeout (Go duration: 45s, 2m, ...)
timeout: 60s

# Retries for transient failures (429, 5xx)
max_retries: 2

# Optional: path to a file whose contents are the API key.
# Used only if GROQ_API_KEY / PROMPTOPT_GROQ_API_KEY are unset.
# api_key_file: ~/.config/promptopt/api-key
`
