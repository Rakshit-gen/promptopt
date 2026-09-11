# promptopt

Understand your prompts, then make them better. `promptopt` is a command-line
tool for working on the prompts you send to language models. It inspects a
prompt, tells you what is weak or wasteful, and rewrites it — while keeping the
job the prompt is meant to do.

You describe the work. promptopt handles the prompt-engineering mechanics:
clarifying objectives, removing redundancy, fixing instruction order, pinning
down output format, and so on. Inference runs on [Groq](https://groq.com).

```
$ promptopt optimize prompt.md

  promptopt  optimize

  tokens  (estimate)
  1,842 → 1,391
  24.5% reduction

  changes
  + clarified the output requirements
  + removed three repeated instructions
  + resolved a conflict between "be concise" and "explain your reasoning"
  + moved the context ahead of the task

  optimized prompt

  ┌──────────────────────────────────────────────
  │ You are a senior backend engineer. Design a
  │ REST API for processing card payments...
  └──────────────────────────────────────────────
```

## Install

```sh
curl -fsSL https://promptopt.dev/install.sh | sh
```

The installer detects your OS and architecture, downloads the release binary
from GitHub Releases, verifies its SHA256 checksum, and installs it into
`~/.local/bin`. It works on macOS and Linux (amd64 and arm64) and needs
nothing but `curl` (or `wget`) and `tar`.

<details>
<summary>Other ways to install</summary>

```sh
# With Go (1.24+)
go install github.com/rakshit-gen/promptopt/cmd/promptopt@latest

# From source
git clone https://github.com/rakshit-gen/promptopt
cd promptopt
make install

# Homebrew (once the tap is published)
brew install rakshit-gen/tap/promptopt
```

Prebuilt binaries and checksums are attached to every
[release](https://github.com/rakshit-gen/promptopt/releases).
</details>

## Quick start

```sh
export GROQ_API_KEY="your-key"   # create one at https://console.groq.com/keys

# Optimize an inline prompt
promptopt optimize "build a REST API for payments"

# Analyze a prompt file, like a linter
promptopt analyze system-prompt.md

# Compress a bloated system prompt, targeting a 30% reduction
promptopt compress system-prompt.md --target 30

# Pipe a prompt in and get machine-readable output back
cat prompt.txt | promptopt eval --json | jq '.scores'
```

## Commands

| Command | What it does |
| --- | --- |
| `promptopt optimize` | Inspect a prompt and rewrite it to be clearer and tighter, preserving intent. |
| `promptopt compress` | Cut tokens without changing what the prompt makes a model do. |
| `promptopt expand` | Turn an underspecified prompt into a detailed one, labeling every assumption. |
| `promptopt analyze` | Score a prompt and report actionable findings with severities and stable IDs. |
| `promptopt transform` | Convert a prompt to markdown, XML, JSON, a system prompt, a template, or an agent brief. |
| `promptopt eval` | Assess a prompt and probe it with generated test cases. |
| `promptopt config` | Inspect the effective configuration and write a starter config file. |

These six operations are the whole public surface. Smaller transformations
(clarify, dedupe, restructure, systemize, and so on) are chosen internally by
each command; they are not commands you have to learn.

Run `promptopt <command> --help` for the full description of a command,
including when to reach for it and when not to.

## Examples

```sh
# Input can be an argument, a file, a path of "-", or piped stdin
promptopt optimize "write unit tests for this function"
promptopt optimize prompt.txt
promptopt optimize - < prompt.txt
cat prompt.txt | promptopt optimize

# Write the result straight to a file (report still goes to stderr)
promptopt optimize prompt.md -o optimized.md

# Only print the resulting prompt, nothing else — good for pipelines
cat draft.txt | promptopt optimize -q > final.txt

# Turn a one-off prompt into a reusable template
promptopt transform "Review this Python code for SQL injection" --to template
#   Review the following {{language}} code for {{issue}}:
#
#   {{code}}

# Expand a rough idea into a production-grade prompt
promptopt expand "build a payment API" --depth production

# Every command speaks JSON
promptopt analyze prompt.txt --json
```

### In CI

`analyze` and `eval` return non-zero only on real errors (bad flags, a
provider failure), not on a low score, so gate on the JSON instead:

```sh
score=$(promptopt analyze prompt.txt --json | jq '.scores.overall')
awk "BEGIN { exit !($score >= 7.5) }" || {
  echo "prompt quality regressed: $score"
  exit 1
}
```

See [docs/ci](docs/ci.md) for a full GitHub Actions example.

## Configuration

Configuration is resolved in this order, each layer overriding the one below:

1. Command-line flags (`--model`, `--timeout`, `--temperature`, ...)
2. Environment variables
3. A config file at `~/.config/promptopt/config.yaml`
4. Built-in defaults

| Variable | Purpose |
| --- | --- |
| `GROQ_API_KEY` | Groq API key. Required. |
| `PROMPTOPT_GROQ_API_KEY` | Alternative key variable; takes precedence over `GROQ_API_KEY`. |
| `PROMPTOPT_MODEL` | Groq model ID. Default: `openai/gpt-oss-120b`. |
| `PROMPTOPT_OUTPUT` | `text` or `json`. |
| `PROMPTOPT_TIMEOUT` | Per-request timeout, e.g. `45s`, `2m`. |
| `PROMPTOPT_MAX_RETRIES` | Retries for transient failures (429, 5xx). |
| `PROMPTOPT_TEMPERATURE` | Sampling temperature, 0–2. |
| `NO_COLOR` | Any value disables ANSI color. |

Write a commented starter file with:

```sh
promptopt config init
promptopt config show   # prints the effective config; the API key is never shown
```

The API key is never written to the config file. Point `api_key_file` at a
file if you do not want it in your environment.

## Architecture

```
cmd/promptopt            main(): wires build info, calls internal/cli
  └─ internal/cli         Cobra commands, flag parsing, input resolution,
       │                  output selection. The only package that imports Cobra.
       ├─ internal/config env + YAML + defaults, with precedence and validation
       ├─ internal/engine the six operations; no CLI or HTTP knowledge
       │    ├─ prompts/          versioned instruction templates, go:embed'd
       │    ├─ internal/groq     the only Groq-aware code; a Completer interface
       │    └─ internal/tokenizer token estimates behind an Estimator interface
       └─ internal/output human renderers and the stable --json encoder
  pkg/types                the result structs; the public --json contract
```

Dependencies point one way: `cli → engine → {groq, prompts, tokenizer}`. The
engine takes a `groq.Completer` interface, so every operation is tested with a
fake and no network. Swapping in another inference provider means implementing
one interface, not touching business logic.

The instructions promptopt sends to the model live in `prompts/<op>/v1.md` and
are embedded into the binary. Every change to how an operation behaves is a
reviewable diff, and the file's front matter carries a version.

## Development

Requires Go 1.24+.

```sh
make build          # build ./bin/promptopt
make test           # go test ./...
make test-race      # with the race detector
make lint           # golangci-lint if present, else gofmt + vet
make run ARGS="analyze prompt.txt"
```

The website lives in `web/` and is isolated from the Go module:

```sh
cd web
npm install
npm run dev
```

## Testing

`go test ./...` runs the full unit suite with no network access — Groq calls
go through a fake `Completer`, and the CLI tests point the client at an
`httptest` server. Coverage spans command parsing, input resolution, config
precedence, token accounting, the Groq client's retry and error handling, the
engine's structured-output parsing, and JSON serialization.

Integration tests that hit the real API are behind a build tag:

```sh
PROMPTOPT_INTEGRATION=1 GROQ_API_KEY=... go test -tags=integration ./tests/...
```

## Releasing

Releases are tag-driven. Pushing a `v*` tag runs
[GoReleaser](https://goreleaser.com) via GitHub Actions, which cross-compiles
for `darwin/{arm64,amd64}` and `linux/{amd64,arm64}`, produces `.tar.gz`
archives with `README`, `LICENSE`, and docs, and publishes `checksums.txt`
(SHA256) alongside them.

```sh
git tag v0.1.0
git push origin v0.1.0
# or, locally, to test the build:
make snapshot
```

## Security and privacy

- Prompts are treated as untrusted input and are **not logged**.
- The API key is **not logged** and is redacted anywhere config is printed.
- **No telemetry.** promptopt makes exactly one kind of network call: to the
  Groq API, using your key.
- Your prompts are sent to Groq for inference and are not stored anywhere by
  promptopt. Groq's handling of that data is governed by Groq's terms.
- `analyze` includes a prompt-injection check. It is a useful signal, not a
  guarantee — an LLM-based review does not make a prompt safe.

See [SECURITY.md](SECURITY.md) for how to report a vulnerability, and
[docs/](docs/) for exactly what leaves your machine.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). In short: `make test lint` should pass,
new behavior needs a test, and changes to how an operation works belong in the
versioned prompt files with a note in the docs.

## License

[MIT](LICENSE).
