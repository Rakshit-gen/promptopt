# Prompt compression

This page explains the model behind `promptopt compress`. For the command
reference, see [commands/compress](commands/compress.md).

## Words versus behavior

A prompt is a set of instructions. Its length is a cost — tokens you pay on
every request, and context you spend that could hold something else. But the
length is not the thing you care about; the behavior is.

Compression is the process of removing length without removing behavior. The
failure mode is removing a sentence that looked redundant but was actually
load-bearing:

- "Respond in JSON." … "Return only the JSON, no prose." — the second line
  looks like a restatement but it is doing separate work.
- Two examples that look similar but demonstrate different edge cases.
- A constraint buried in a paragraph of context.

promptopt asks the model to make these cuts and to report its confidence that
behavior survived, plus any specific risks. A compression that drops the score
is telling you it is not sure — and you should look.

## What gets cut

compress targets:

- **Repeated concepts** — the same idea stated more than once.
- **Restatement** — a sentence that says what the previous sentence said.
- **Verbose phrasing** — "in order to" → "to", "it is important that you" →
  an imperative.
- **Duplicate examples** — keeps the clearest one.
- **Meaningless formatting** — decorative separators, headers with one line
  under them.
- **Constraints repeated in different words.**
- **Explanations of things the model already knows** — what JSON is, what a
  REST API is.
- **Contextual prose that does not change the output.**

## What is kept

- Every distinct instruction and constraint.
- The output format.
- At least one example per behavior that examples were teaching.
- Domain terms and named entities.
- Anything load-bearing, even if it is wordy.

## The knobs

### `--target <n>`

A goal, not a guarantee. compress tries to reach it and stops when the next
cut would risk behavior. If it lands short, the report says by how much and
why.

### `--aggressive`

Permission to drop marginal context and collapse examples harder. Use it when
you are close to a hard limit and can tolerate re-testing.

### `--preserve-behavior`

The opposite: refuse any cut the model is not confident about, even if that
misses the target. Use it for a prompt you cannot easily re-validate.

Set both and `preserve-behavior` wins on each individual decision.

## The semantic-preservation score

0–10, the model's confidence that the compressed prompt produces the same
behavior as the original.

| Score | Read it as |
| --- | --- |
| 9–10 | Safe to use. Skim the diff. |
| 7–8 | Probably fine. Read the behavior-risks list. |
| below 7 | Re-test before you rely on it. If the risks list is empty, promptopt adds a warning saying so. |

## A workflow

```sh
# 1. See where the tokens are going
promptopt analyze system.md --json | jq '.scores.efficiency'

# 2. Compress conservatively
promptopt compress system.md --preserve-behavior -o system.min.md

# 3. Confirm behavior held
promptopt eval system.min.md

# 4. If eval is happy and you need more, go aggressive
promptopt compress system.md --target 45 --aggressive -o system.min.md
promptopt eval system.min.md
```

Compression is not a one-shot. It is a cut, a check, and a decision.
