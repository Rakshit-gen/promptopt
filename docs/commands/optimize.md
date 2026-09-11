# optimize

```
promptopt optimize [prompt|file|-] [flags]
```

Inspect a prompt for the problems that make language models unreliable, then
return a rewritten version that keeps the same intent.

## What it does

`optimize` first reads the prompt the way a reviewer would, looking for:

- an unclear or missing objective
- context the model would need but does not have
- ambiguous wording and undefined terms
- instructions that repeat, or that contradict each other
- an order that makes the prompt harder to follow than it needs to be
- no stated output format
- verbosity that costs tokens without changing behavior
- weak constraints, inert role descriptions, unclear examples
- edge cases the prompt implies but does not handle

Then it rewrites the prompt to fix what it can and reports what it changed,
the token count before and after, a quick quality read on the result, and any
warnings — for example, intent it had to infer.

## When to use it

- A prompt works "most of the time" and you want it to work more of the time.
- You inherited a prompt and want it tightened before you build on it.
- `promptopt analyze` flagged problems and you want them fixed.

## When not to use it

- **The prompt is underspecified on purpose and you want it filled in.** That
  is [`expand`](expand.md). `optimize` will not invent domain requirements.
- **You only care about token count.** Use [`compress`](compress.md).
  `optimize` improves clarity first; it often makes a terse prompt longer.
- **You want a report, not a rewrite.** Use [`analyze`](analyze.md).

## How it works

promptopt sends the prompt to the model with a versioned instruction set
(`prompts/optimize/v1.md` in the repo, embedded in the binary) that asks for a
structured JSON object: the rewritten prompt, a list of changes, warnings, and
sub-scores. promptopt validates that object, estimates tokens for the original
and the rewrite, and renders the result. If the model does not return usable
JSON, the command fails with exit code 8 rather than guessing.

The rewrite is constrained to preserve intent. The instruction set explicitly
forbids adding requirements the author did not state or clearly imply.

## Flags

| Flag | Default | Effect |
| --- | --- | --- |
| `--goal <text>` | — | Extra guidance about what the prompt should achieve. Use it when the objective is not obvious from the prompt alone. |
| `--output <file>`, `-o` | — | Write the optimized prompt to a file. |
| `--json` | — | Emit the `OptimizeResult` object. |
| `--quiet`, `-q` | — | Print only the optimized prompt. |

Plus the [shared flags](README.md#shared-flags).

## Examples

```sh
# Inline
promptopt optimize "build a REST API for payments"

# A file, writing the result to another file
promptopt optimize system-prompt.md -o system-prompt.optimized.md

# Piped, machine-readable
cat prompt.txt | promptopt optimize --json

# With a hint about the goal
promptopt optimize draft.txt --goal "the model should return only a unified diff"

# In a pipeline: optimize, then compress the result
cat draft.txt | promptopt optimize -q | promptopt compress -q > final.txt
```

## Expected output

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
  + made the response format explicit

  quality  (optimized prompt)
    overall          ███████████████·····  7.9

  optimized prompt

  ┌──────────────────────────────────────────────
  │ You are a senior backend engineer. Design a
  │ REST API for processing card payments.
  │ ...
  └──────────────────────────────────────────────
```

With `--json`:

```json
{
  "operation": "optimize",
  "result": "You are a senior backend engineer. ...",
  "changes": ["clarified the output requirements", "..."],
  "tokens": { "before": 1842, "after": 1391, "reduction_percent": 24.48, "source": "estimate" },
  "analysis": { "overall": 7.9, "clarity": 8.4, "specificity": 7.5, "consistency": 9.0, "efficiency": 7.1, "robustness": 7.5 },
  "warnings": []
}
```

## Common mistakes

- **Expecting new requirements.** If your prompt says "build an API" and you
  want auth, rate limits, and pagination spelled out, run `expand`, not
  `optimize`.
- **Treating a longer result as a failure.** `optimize` trades verbosity for
  clarity. If the token count went up and that matters, pipe the result
  through `compress`. The command prints a warning when the result grows a
  lot.
- **Running it on a prompt that is already good.** It will make small changes
  and tell you the prompt was already in good shape. That is a valid result,
  not a wasted call — but `analyze` is cheaper if you just want confirmation.

## Token counts

Counts are labeled `(estimate)`. promptopt does not ship the exact tokenizer
for every Groq model, so it uses a heuristic that tracks real tokenizers to
within a few percent. See [JSON output](../json-output.md#tokenreport) and
[concepts](../concepts.md#token-counting).
