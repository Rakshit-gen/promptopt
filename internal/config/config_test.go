package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"GROQ_API_KEY", "PROMPTOPT_GROQ_API_KEY", "PROMPTOPT_MODEL",
		"PROMPTOPT_BASE_URL", "PROMPTOPT_OUTPUT", "PROMPTOPT_TIMEOUT",
		"PROMPTOPT_MAX_RETRIES", "PROMPTOPT_TEMPERATURE", "NO_COLOR",
		"XDG_CONFIG_HOME",
	} {
		t.Setenv(k, "")
		os.Unsetenv(k)
	}
}

func TestDefaults(t *testing.T) {
	clearEnv(t)
	cfg, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("explicit missing config file should error")
	}
	_ = cfg

	cfg, err = Load("") // default path, likely missing -> not an error
	if err != nil {
		t.Fatalf("default load: %v", err)
	}
	if cfg.Model != "openai/gpt-oss-120b" {
		t.Errorf("default model = %q", cfg.Model)
	}
	if cfg.Output != "text" || cfg.Color != "auto" {
		t.Errorf("unexpected defaults: output=%q color=%q", cfg.Output, cfg.Color)
	}
}

func TestEnvOverrides(t *testing.T) {
	clearEnv(t)
	t.Setenv("GROQ_API_KEY", "gsk_test_key_value")
	t.Setenv("PROMPTOPT_MODEL", "llama-3.3-70b-versatile")
	t.Setenv("PROMPTOPT_TIMEOUT", "90s")
	t.Setenv("PROMPTOPT_MAX_RETRIES", "5")
	t.Setenv("NO_COLOR", "1")

	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "gsk_test_key_value" {
		t.Errorf("api key not read from env")
	}
	if cfg.Model != "llama-3.3-70b-versatile" {
		t.Errorf("model = %q", cfg.Model)
	}
	if cfg.Timeout != 90*time.Second {
		t.Errorf("timeout = %v", cfg.Timeout)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("retries = %d", cfg.MaxRetries)
	}
	if cfg.Color != "never" {
		t.Errorf("NO_COLOR should force color=never, got %q", cfg.Color)
	}
}

func TestFileParsing(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	body := `# promptopt config
model: "openai/gpt-oss-20b"   # inline comment
output: json
temperature: 0.7
timeout: 2m
max_retries: 1
color: always
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Model != "openai/gpt-oss-20b" {
		t.Errorf("model = %q", cfg.Model)
	}
	if cfg.Output != "json" {
		t.Errorf("output = %q", cfg.Output)
	}
	if cfg.Temperature != 0.7 {
		t.Errorf("temperature = %v", cfg.Temperature)
	}
	if cfg.Timeout != 2*time.Minute {
		t.Errorf("timeout = %v", cfg.Timeout)
	}
	if cfg.Source() != path {
		t.Errorf("source = %q", cfg.Source())
	}
}

func TestEnvBeatsFile(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	os.WriteFile(path, []byte("model: from-file\n"), 0o644)
	t.Setenv("PROMPTOPT_MODEL", "from-env")

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Model != "from-env" {
		t.Errorf("env should win over file, got %q", cfg.Model)
	}
}

func TestValidation(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	os.WriteFile(path, []byte("output: yaml\n"), 0o644)
	if _, err := Load(path); err == nil {
		t.Fatal("invalid output mode should fail validation")
	}
}

func TestRedacted(t *testing.T) {
	cfg := Config{APIKey: "gsk_secret"}
	if cfg.Redacted().APIKey != "****" {
		t.Fatal("redacted config still exposes the key")
	}
}

func TestAPIKeyFile(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "key")
	os.WriteFile(keyPath, []byte("  gsk_from_file\n"), 0o600)
	cfgPath := filepath.Join(dir, "config.yaml")
	os.WriteFile(cfgPath, []byte("api_key_file: "+keyPath+"\n"), 0o644)

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "gsk_from_file" {
		t.Errorf("api key from file = %q", cfg.APIKey)
	}
}
