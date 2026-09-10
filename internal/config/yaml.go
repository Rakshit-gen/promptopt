package config

import (
	"fmt"
	"strconv"
	"strings"
)

// parseYAML reads the small, flat subset of YAML that promptopt's config file
// uses: one `key: value` pair per line, `#` comments, and optionally quoted
// string values. The config has no nesting, lists, or anchors, so a full YAML
// dependency would not earn its place here. Unknown keys are ignored so a
// newer config file does not break an older binary.
func parseYAML(data []byte, cfg *Config) error {
	for i, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(stripComment(raw))
		if line == "" {
			continue
		}
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			return fmt.Errorf("line %d: expected 'key: value', got %q", i+1, raw)
		}
		key = strings.TrimSpace(key)
		val = unquote(strings.TrimSpace(val))

		switch key {
		case "model":
			cfg.Model = val
		case "base_url":
			cfg.BaseURL = val
		case "output":
			cfg.Output = val
		case "color":
			cfg.Color = val
		case "api_key_file":
			cfg.APIKeyFile = val
		case "timeout":
			cfg.TimeoutRaw = val
		case "temperature":
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return fmt.Errorf("line %d: temperature %q is not a number", i+1, val)
			}
			cfg.Temperature = f
		case "max_retries":
			n, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("line %d: max_retries %q is not an integer", i+1, val)
			}
			cfg.MaxRetries = n
		default:
			// ignore unknown keys
		}
	}
	return nil
}

// stripComment removes a trailing `#` comment that is not inside quotes.
func stripComment(line string) string {
	inSingle, inDouble := false, false
	for i, r := range line {
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				return line[:i]
			}
		}
	}
	return line
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
