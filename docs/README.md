# promptopt documentation

promptopt is a command-line tool for working on the prompts you send to
language models. It inspects a prompt, reports what is weak or wasteful, and
rewrites it — while keeping the job the prompt is meant to do.

## Start here

- [Getting started](getting-started.md) — install, set a key, run the first command
- [Installation](installation.md) — every install method and how the installer works
- [Configuration](configuration.md) — env vars, the config file, precedence

## Commands

- [Overview](commands/README.md)
- [optimize](commands/optimize.md) — rewrite a prompt to be clearer and tighter
- [compress](commands/compress.md) — cut tokens without changing behavior
- [expand](commands/expand.md) — turn an underspecified prompt into a detailed one
- [analyze](commands/analyze.md) — score a prompt and report findings
- [transform](commands/transform.md) — convert a prompt to another representation
- [eval](commands/eval.md) — assess a prompt and probe it with test cases

## Concepts

- [How promptopt works](concepts.md)
- [Prompt analysis](prompt-analysis.md)
- [Prompt compression](prompt-compression.md)
- [Evaluation](evaluation.md)
- [JSON output](json-output.md)

## Operating

- [Using promptopt in CI](ci.md)
- [Troubleshooting](troubleshooting.md)

## What leaves your machine

promptopt makes exactly one kind of network call: to the Groq API, using your
`GROQ_API_KEY`, to run inference on the prompt you passed it. It does not log
prompts, does not log the key, and sends no telemetry. See
[SECURITY.md](../SECURITY.md).
