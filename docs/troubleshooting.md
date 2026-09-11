# Troubleshooting

Run any command with `--verbose` to see the underlying error and diagnostics.

## "No Groq API key configured." (exit 4)

promptopt did not find a key. It checks, in order: `PROMPTOPT_GROQ_API_KEY`,
`GROQ_API_KEY`, then `api_key_file` in the config.

```sh
export GROQ_API_KEY="your-key"      # from https://console.groq.com/keys
promptopt config show               # confirm it now shows "api key  set"
```

If you set it and promptopt still does not see it, you probably set it in a
different shell. Add the export to `~/.zshrc` or `~/.bashrc`.

## "Groq rejected the API key." (exit 4)

The key was sent but Groq refused it. It is wrong, revoked, or has no quota.
Generate a fresh one at <https://console.groq.com/keys>.

## "Groq rate limit reached." (exit 5)

You hit Groq's request or token rate limit. The message includes the
provider's suggested retry interval when it sends one. Options:

- Wait and retry.
- Raise retries: `export PROMPTOPT_MAX_RETRIES=4` (promptopt honours
  `Retry-After`).
- Use a lighter model: `--model openai/gpt-oss-20b`.

## "The prompt is longer than the context window..." (exit 6)

The prompt plus promptopt's instructions exceeds the model's context.

```sh
promptopt compress big-prompt.md -o smaller.md   # then retry on smaller.md
# or use a larger-context model
promptopt analyze big-prompt.md --model openai/gpt-oss-120b
```

## "The model did not return usable JSON..." (exit 8)

The model's response could not be parsed into the expected structure. Usually
transient, so retry. If it persists on a specific model, that model handles
structured output poorly; switch:

```sh
promptopt analyze prompt.txt --model openai/gpt-oss-120b
```

Report a reproducible case at
<https://github.com/rakshit-gen/promptopt/issues>.

## "Groq request timed out..." 

The request took longer than the timeout (default 60s). Large prompts with a
slower model can exceed it.

```sh
promptopt optimize prompt.txt --timeout 3m
# or persist it
export PROMPTOPT_TIMEOUT=180s
```

## "Groq does not recognize the model..." (exit 7)

The model ID is wrong or has been retired. List current models at
<https://console.groq.com/docs/models> and set a valid one with `--model` or
`PROMPTOPT_MODEL`.

## "no prompt given" (exit 2)

promptopt got no argument and stdin was not piped. Pass a prompt:

```sh
promptopt analyze "the prompt"
promptopt analyze prompt.txt
cat prompt.txt | promptopt analyze
```

If you are calling promptopt from a script and expect stdin, make sure you are
actually piping into it and not just redirecting a closed descriptor.

## Colors are showing up in a file or a pipe

They should not: promptopt disables color when stdout is not a terminal. If a
tool is capturing output through a pseudo-terminal, force it off:

```sh
promptopt analyze prompt.txt --no-color
# or
export NO_COLOR=1
```

## `go install` gives `--version` as "dev"

`go install` does not pass build flags, so the version string is not stamped.
Use a release binary, the curl installer, or `make install` from a checkout.

## The installer says "unsupported architecture"

promptopt ships for `amd64` and `arm64` on macOS and Linux. For anything else
(32-bit ARM, Windows, BSD), build from source with Go:

```sh
go install github.com/rakshit-gen/promptopt/cmd/promptopt@latest
```

## Getting more detail

```sh
promptopt <command> ... --verbose
```

This prints the wrapped error. For the Groq exchange itself, promptopt does
not log request bodies (they contain your prompt); if you need to see the wire
traffic, set `PROMPTOPT_BASE_URL` to a local logging proxy you control.
