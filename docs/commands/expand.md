# expand

```
promptopt expand [prompt|file|-] [flags]
```

Turn a short, underspecified prompt into a detailed one a model can execute
without guessing.

## What it does

expand takes a prompt like `build a payment API` and fills in the parts a
model would otherwise have to invent: objective, context, constraints,
assumptions, inputs, outputs and their format, edge cases, failure handling,
evaluation criteria, and at least one worked example.

The part that matters is the line between **what you asked for** and **what
expand added**. Everything expand adds that is not in your original is labeled
as an assumption — both in the report and inside the expanded prompt — so you
can correct it. expand does not invent facts, numbers, names, or domain rules;
where a detail is needed and missing, it makes a reasonable default explicit
rather than silently baking it in.

## When to use it

- You have a one-line idea and want a prompt you can actually run.
- You are drafting a prompt for a system and need the edge cases and failure
  behavior thought through.
- A prompt keeps producing plausible-but-wrong output because it leaves too
  much open.

## When not to use it

- **The prompt is already specified and just needs tightening.** Use
  [`optimize`](optimize.md). expand will pad a prompt that did not need
  padding.
- **You want the model to answer the prompt.** expand rewrites the prompt; it
  does not execute it.
- **You need exact domain requirements.** expand marks its guesses as
  assumptions precisely because it does not know your domain. Review and
  replace them.

## Depth

| `--depth` | What you get |
| --- | --- |
| `concise` | Only the gaps that block execution are filled. Stays tight. |
| `detailed` | (default) Adds structure, constraints, an output format, and examples. |
| `production` | A prompt you could hand to a model in a real system: inputs, outputs, edge cases, failure handling, and evaluation criteria. |

## How it works

The prompt and depth go to the model with the `prompts/expand/v1.md`
instruction set, which asks for a JSON object: the expanded prompt, the list
of components added, and an array of assumptions (each with a field, a value,
and a note on how to change it). promptopt validates it and shows the
assumptions prominently, because they are the thing you most need to check.

## Flags

| Flag | Default | Effect |
| --- | --- | --- |
| `--depth <level>` | `detailed` | `concise`, `detailed`, or `production`. |
| `--output <file>`, `-o` | — | Write the expanded prompt to a file. |
| `--json` | — | Emit the `ExpandResult` object. |
| `--quiet`, `-q` | — | Print only the expanded prompt. |

## Examples

```sh
promptopt expand "build a payment API"

# A production-grade prompt
promptopt expand "build a payment API" --depth production

# Just enough to run
promptopt expand idea.txt --depth concise

# Check what it assumed
promptopt expand "summarize support tickets" --json | jq '.assumptions'
```

## Expected output

```
$ promptopt expand "build a payment API" --depth production

  promptopt  expand · production

  tokens  (estimate)
  6 → 512
  8433.3% larger

  added
  + explicit objective
  + input and output contracts
  + idempotency and retry behavior
  + error taxonomy and failure responses
  + evaluation criteria

  assumptions  (correct these if wrong)
  ~ payment processor: Stripe-style API with a PaymentIntent object
      Change to your processor's model if different.
  ~ currency: amounts are integer minor units (cents)
      Standard for card APIs; adjust if you use decimal major units.

  expanded prompt

  ┌──────────────────────────────────────────────
  │ ## Objective
  │ Design the HTTP API for a service that ...
  │
  │ ## Assumptions
  │ - Payment processor: Stripe-style ...
  └──────────────────────────────────────────────
```

## Common mistakes

- **Shipping the assumptions unread.** They are guesses. If expand assumed
  Postgres and you use DynamoDB, the expanded prompt is now subtly wrong until
  you fix that line.
- **Using `production` for a throwaway.** It produces a lot of prompt. For a
  quick experiment, `concise` is usually enough.
- **Expecting real API contracts.** expand writes a plausible contract shape,
  not your actual endpoints. Treat it as a skeleton.
