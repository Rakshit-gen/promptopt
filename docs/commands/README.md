# Commands

promptopt has six operations. They are the whole public surface.

| Command | Use it to |
| --- | --- |
| [optimize](optimize.md) | Rewrite a prompt so a model follows it more reliably, without changing what it asks for. |
| [compress](compress.md) | Reduce token count while keeping the prompt's behavior. |
| [expand](expand.md) | Fill in the parts an underspecified prompt leaves a model to guess. |
| [analyze](analyze.md) | Get a scored report of a prompt's problems, without changing it. |
| [transform](transform.md) | Convert a prompt to markdown, XML, JSON, a system prompt, a template, or an agent brief. |
| [eval](eval.md) | Check whether a prompt is ready to depend on, with generated test cases. |

Smaller transformations (clarify, dedupe, reorder, systemize, parameterize)
are chosen internally by each command. They are not commands you invoke.

## Shared behavior

### Input

Every command reads a prompt from one of:

```sh
promptopt optimize "the prompt text"     # argument
promptopt optimize prompt.txt            # file path
promptopt optimize - < prompt.txt        # explicit stdin
cat prompt.txt | promptopt optimize      # piped stdin
```

An argument with no spaces and no `/` that is not an existing file is treated
as a literal prompt, not a missing-file error.

### Shared flags

| Flag | Effect |
| --- | --- |
| `--model <id>` | Override the Groq model for this run. |
| `--json` | Emit a single stable JSON object, nothing else. |
| `--quiet`, `-q` | Print only the resulting prompt (not for `analyze`/`eval`). |
| `--output <file>`, `-o` | Write the resulting prompt to a file; the report still prints. |
| `--config <path>` | Use a specific config file. |
| `--timeout <dur>` | Per-request timeout, e.g. `45s`. |
| `--temperature <n>` | Sampling temperature 0–2. |
| `--no-color` | Disable ANSI color (also respects `NO_COLOR` and non-TTY output). |
| `--verbose`, `-v` | Show underlying errors and diagnostics. |

### Output

Default output is formatted for a terminal: a short header, token accounting
where relevant, a summary of changes or findings, and then the resulting
prompt in a box. Color is used sparingly and turns itself off when stdout is
not a terminal.

`--json` output is exactly the result struct documented in
[JSON output](../json-output.md), with no decorative text. The two modes never
mix.

### Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Success |
| 1 | Generic error |
| 2 | Usage error (bad flag, no input) |
| 3 | Configuration error |
| 4 | Missing or rejected API key |
| 5 | Rate limited |
| 6 | Prompt exceeds the model's context window |
| 7 | Other provider error |
| 8 | Model returned a response promptopt could not parse |

A low quality score is **not** an error. `analyze` and `eval` exit 0 with a
low score; gate on the JSON if you want a pipeline to fail. See
[Using promptopt in CI](../ci.md).
