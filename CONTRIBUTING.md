# Contributing to promptopt

Thanks for taking the time. This is a focused tool, and the goal is to keep it
that way.

## Setup

```sh
git clone https://github.com/rakshit-gen/promptopt
cd promptopt
make build
make test
```

Requires Go 1.24 or newer. For the website, `cd web && npm install`.

## Before you open a PR

```sh
make fmt-check   # gofmt-clean
make vet
make test        # unit tests, no network
make lint        # golangci-lint if installed
cd web && npm run lint && npm run build   # if you touched web/
```

CI runs the same checks on Linux and macOS.

## What changes go where

- **How an operation behaves** (what the model is told to do) lives in
  `prompts/<op>/v1.md`. These are embedded into the binary. If a change is
  significant, bump the version file (`v2.md`) rather than editing `v1.md` in
  place, and update `prompts.Load`.
- **How results are shaped** is `pkg/types`. That is the public `--json`
  contract; adding fields is fine, renaming or removing them is a breaking
  change.
- **Provider behavior** is `internal/groq` and nowhere else. If you are adding
  a second provider, implement `groq.Completer` (consider renaming the
  interface to a neutral package) and keep `internal/engine` untouched.
- **Terminal output** is `internal/output`. Never print anything on the
  `--json` path except the JSON object.

## Adding a transform target

`transform` targets are a good first contribution:

1. Add the constant and name to `internal/engine/transform.go` (`Target`,
   `Targets()`).
2. Describe it in `prompts/transform/v1.md` under "Targets".
3. Add a case to the docs at `docs/commands/transform.md` and `web/`'s command
   docs.
4. Add a test in `internal/engine/engine_test.go`.

No other wiring is needed — the CLI reads `engine.Targets()` for validation
and help.

## Tests

New behavior needs a test. Unit tests must not touch the network: use the fake
`Completer` in `internal/engine/engine_test.go` or the `httptest` pattern in
`internal/cli/cli_test.go`. Tests that genuinely need the API go in `tests/`
behind the `integration` build tag.

## Commit style

Small, focused commits with imperative subject lines. Conventional-commit
prefixes (`feat:`, `fix:`, `docs:`) are used for the changelog but not
required for every commit.

## Reporting bugs

Use the issue templates. A prompt that reproduces the problem (redact anything
sensitive) and the output of `promptopt --version` go a long way.

## Code of conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md).
