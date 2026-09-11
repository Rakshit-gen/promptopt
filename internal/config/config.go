// Package config resolves promptopt's runtime configuration from, in order of
// precedence: command-line flags (applied by the caller), environment
// variables, a YAML config file, and built-in defaults.
//
// The API key is never stored in the resolved struct's string form beyond
// what is needed to make requests, and is never logged. It is read from the
// environment, or from a file the config points at, but the config file
// itself should hold a reference, not the raw key.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Config is the resolved configuration for a single invocation.
type Config struct {
	// APIKey is the Groq key. Resolved from GROQ_API_KEY or
	// PROMPTOPT_GROQ_API_KEY, or from api_key_file in the config file.
	APIKey string `yaml:"-"`

	Model       string        `yaml:"model"`
	BaseURL     string        `yaml:"base_url"`
	Output      string        `yaml:"output"` // "text" or "json"
	Color       string        `yaml:"color"`  // "auto", "always", "never"
	Temperature float64       `yaml:"temperature"`
	Timeout     time.Duration `yaml:"-"`
	TimeoutRaw  string        `yaml:"timeout"`
	MaxRetries  int           `yaml:"max_retries"`

	// APIKeyFile is an optional path to a file whose contents are the API key.
	APIKeyFile string `yaml:"api_key_file"`

	// source records where the file was loaded from, for diagnostics.
	source string
}

// Defaults returns a Config with built-in defaults and no secrets.
func Defaults() Config {
	return Config{
		Model:       "openai/gpt-oss-120b",
		BaseURL:     "https://api.groq.com/openai/v1",
		Output:      "text",
		Color:       "auto",
		Temperature: 0.3,
		Timeout:     60 * time.Second,
		MaxRetries:  2,
	}
}

// DefaultPath returns the standard config file location, honoring
// XDG_CONFIG_HOME.
func DefaultPath() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "promptopt", "config.yaml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".config", "promptopt", "config.yaml")
	}
	return filepath.Join(home, ".config", "promptopt", "config.yaml")
}

// Load builds the effective config. If path is empty, DefaultPath() is used
// and a missing file is not an error.
func Load(path string) (Config, error) {
	cfg := Defaults()

	explicit := path != ""
	if path == "" {
		path = DefaultPath()
	}

	if data, err := os.ReadFile(path); err == nil {
		if err := parseYAML(data, &cfg); err != nil {
			return cfg, fmt.Errorf("parsing %s: %w", path, err)
		}
		cfg.source = path
		if cfg.TimeoutRaw != "" {
			d, err := time.ParseDuration(cfg.TimeoutRaw)
			if err != nil {
				return cfg, fmt.Errorf("parsing %s: invalid timeout %q", path, cfg.TimeoutRaw)
			}
			cfg.Timeout = d
		}
	} else if explicit {
		return cfg, fmt.Errorf("reading config file %s: %w", path, err)
	} else if !os.IsNotExist(err) {
		return cfg, fmt.Errorf("reading config file %s: %w", path, err)
	}

	applyEnv(&cfg)

	if err := cfg.validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func applyEnv(cfg *Config) {
	if v := firstEnv("PROMPTOPT_GROQ_API_KEY", "GROQ_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if cfg.APIKey == "" && cfg.APIKeyFile != "" {
		if data, err := os.ReadFile(expandHome(cfg.APIKeyFile)); err == nil {
			cfg.APIKey = strings.TrimSpace(string(data))
		}
	}
	if v := os.Getenv("PROMPTOPT_MODEL"); v != "" {
		cfg.Model = v
	}
	if v := os.Getenv("PROMPTOPT_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("PROMPTOPT_OUTPUT"); v != "" {
		cfg.Output = v
	}
	if v := os.Getenv("PROMPTOPT_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Timeout = d
		}
	}
	if v := os.Getenv("PROMPTOPT_MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.MaxRetries = n
		}
	}
	if v := os.Getenv("PROMPTOPT_TEMPERATURE"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.Temperature = f
		}
	}
	if v := os.Getenv("NO_COLOR"); v != "" {
		cfg.Color = "never"
	}
}

func (cfg Config) validate() error {
	switch cfg.Output {
	case "text", "json":
	default:
		return fmt.Errorf("invalid output mode %q (want text or json)", cfg.Output)
	}
	switch cfg.Color {
	case "auto", "always", "never":
	default:
		return fmt.Errorf("invalid color mode %q (want auto, always, or never)", cfg.Color)
	}
	if cfg.Temperature < 0 || cfg.Temperature > 2 {
		return fmt.Errorf("temperature %.2f is out of range 0-2", cfg.Temperature)
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if cfg.MaxRetries < 0 {
		return fmt.Errorf("max_retries cannot be negative")
	}
	return nil
}

// Source returns the path the config file was loaded from, or "" if none.
func (cfg Config) Source() string { return cfg.source }

// Redacted returns a copy safe to print: the API key is replaced with a
// fixed-length mask that reveals nothing.
func (cfg Config) Redacted() Config {
	c := cfg
	if c.APIKey != "" {
		c.APIKey = "****"
	}
	return c
}

func firstEnv(keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return ""
}

func expandHome(p string) string {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}
