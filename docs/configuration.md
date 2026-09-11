# Configuration

promptopt resolves configuration from four layers. Each overrides the one
below it:

1. **Command-line flags**: `--model`, `--timeout`, `--temperature`,
   `--output`/`--json`, `--no-color`
2. **Environment variables**
3. **Config file**: `~/.config/promptopt/config.yaml`
4. **Built-in defaults**

To see what a run would actually use:

```sh
promptopt config show
```

```
config file    /Users/you/.config/promptopt/config.yaml
api key        set (gsk_…4b)
model          openai/gpt-oss-120b
base url       https://api.groq.com/openai/v1
output         text
color          auto
temperature    0.30
timeout        1m0s
max retries    2
```

The API key is shown only as "set" or "not set" with a short mask. It is
never printed in full and never written to the config file.

## Environment variables

| Variable | Effect |
| --- | --- |
| `GROQ_API_KEY` | Groq API key. Required for any command that calls the API. |
| `PROMPTOPT_GROQ_API_KEY` | Same, but wins over `GROQ_API_KEY`. Use it to scope a key to promptopt. |
| `PROMPTOPT_MODEL` | Groq model ID. |
| `PROMPTOPT_BASE_URL` | Groq API base URL. Change only if you proxy the API. |
| `PROMPTOPT_OUTPUT` | `text` or `json`. |
| `PROMPTOPT_TIMEOUT` | Per-request timeout as a Go duration: `45s`, `2m`. |
| `PROMPTOPT_MAX_RETRIES` | Retries for transient failures (HTTP 429 and 5xx). |
| `PROMPTOPT_TEMPERATURE` | Sampling temperature, 0–2. |
| `NO_COLOR` | Any non-empty value disables ANSI color. |
| `XDG_CONFIG_HOME` | Moves the config file location to `$XDG_CONFIG_HOME/promptopt/config.yaml`. |

## Config file

Write a commented starter:

```sh
promptopt config init          # writes ~/.config/promptopt/config.yaml
promptopt config path          # prints the path
promptopt config init --force  # overwrite an existing file
```

```yaml
# promptopt configuration

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
```

The parser accepts one `key: value` per line, `#` comments, and quoted
strings. It has no nesting. Unknown keys are ignored, so a config file written
by a newer promptopt will not break an older binary.

## Choosing a model

The default is `openai/gpt-oss-120b`: a large context window and strong
reasoning, which suits prompt analysis and rewriting. Other reasonable
choices on Groq include `openai/gpt-oss-20b` (faster, cheaper) and
`llama-3.3-70b-versatile`. List what is currently available at
<https://console.groq.com/docs/models>.

```sh
promptopt optimize prompt.txt --model openai/gpt-oss-20b
```

If a model does not support structured JSON output well, `promptopt` will
report a "did not return usable JSON" error. Switch to a model that does.

## Keeping the key out of your environment

If you would rather not export `GROQ_API_KEY` globally:

```sh
mkdir -p ~/.config/promptopt
printf '%s' "your-key" > ~/.config/promptopt/api-key
chmod 600 ~/.config/promptopt/api-key
```

Then set `api_key_file: ~/.config/promptopt/api-key` in the config. It is read
only when neither key variable is set.
