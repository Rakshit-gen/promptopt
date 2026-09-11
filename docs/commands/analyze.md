# analyze

```
promptopt analyze [prompt|file|-] [flags]
```

Score a prompt and report what is wrong with it, without changing it. Think of
it as a linter for prompts.

## What it does

analyze reads a prompt and returns:

- **Scores (0–10)** for clarity, specificity, completeness, consistency, token
  efficiency, and robustness, plus an overall score that weighs the sub-scores
  and the severity of findings.
- **Findings**, each with a severity, a stable ID, a one-line title, and a
  detail that says what to do about it.
- A short **summary** of the prompt's main strengths and the single most
  important thing to fix.

| Severity | Meaning |
| --- | --- |
| `ERROR` | Likely to cause wrong or unsafe output — conflicting instructions, contradictory constraints, an injection path that overrides the task. |
| `WARNING` | Likely to degrade quality — undefined output format, repeated constraints, ambiguous terms, missing edge cases. |
| `INFO` | Worth knowing, low impact — an inert role line, mild verbosity, a stylistic inconsistency. |

Finding IDs are stable (`P001`–`P014` for the common categories, `P0xx` for
others) so you can reference them in review and, later, suppress them.

## When to use it

- Before building on a prompt you did not write.
- In CI, to catch a prompt regressing before it ships.
- To decide whether a prompt needs `optimize`, `expand`, or nothing.

## When not to use it

- **You already know it needs work and want it fixed.** Go straight to
  [`optimize`](optimize.md), which runs its own inspection.
- **You want to know how it behaves under pressure.** [`eval`](eval.md)
  probes with test cases; analyze reasons about the text.

## How it works

The prompt goes to the model with the `prompts/analyze/v1.md` instruction set,
which asks for a JSON object of scores, findings, and a summary. promptopt
validates it, normalizes severities, and sorts findings most-serious-first.
The overall score is taken from the model but falls back to a weighted average
of the sub-scores if the model omits it.

analyze does not call any other tool or check. It is one model call producing
one structured report.

## Flags

analyze takes the [shared flags](README.md#shared-flags). It has no
operation-specific flags. `--quiet` and `--output` do nothing here — there is
no single "result prompt", only the report. Use `--json` for machine output.

## Examples

```sh
promptopt analyze prompt.txt
cat prompt.txt | promptopt analyze

# Just the overall score
promptopt analyze prompt.txt --json | jq '.scores.overall'

# Fail a CI job if there are any ERROR findings
promptopt analyze prompt.txt --json \
  | jq -e '[.findings[] | select(.severity == "ERROR")] | length == 0' > /dev/null
```

## Expected output

```
$ promptopt analyze prompt.txt

  promptopt  analyze

  overall score  7.8

    clarity          ████████████████·····  8.4
    specificity      █████████████········  6.7
    completeness     ██████████████·······  7.1
    consistency      ██████████████████···  9.2
    efficiency       ████████████·········  6.1
    robustness       ███████████████······  7.5

  findings  0 error, 2 warning, 1 info

  WARNING P003  The expected output format is not explicitly defined
      Add a sentence stating the format, e.g. "Return a JSON object with
      keys `summary` and `action_items` (array of strings)."

  WARNING P008  The same constraint appears three times
      "Do not include an apology" is stated in the intro, the rules list,
      and the closing line. Keep one.

  INFO P012  The role definition is unlikely to affect the requested task
      "You are a world-class expert" adds tokens without changing behavior
      for a straightforward extraction task.

  summary
  The prompt is clear about the task but leaves the output shape open and
  repeats one constraint. Defining the format is the highest-value fix.
```

## Common mistakes

- **Treating the overall score as pass/fail.** It is a relative signal. A 6.5
  with one fixable WARNING is often fine; a 7.5 with an ERROR is not.
- **Expecting the exit code to reflect the score.** It does not. analyze exits
  0 unless something actually failed. Gate on the JSON.
- **Running it once and moving on.** Fix the findings, then run it again — the
  scores should move, and new findings sometimes surface once the loud ones
  are gone.
