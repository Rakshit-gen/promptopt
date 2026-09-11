# compress

```
promptopt compress [prompt|file|-] [flags]
```

Reduce a prompt's token count without changing what the prompt makes a model
do.

## What it does

compress is useful when a prompt has accumulated instructions, examples, and
context over time and is now paying for tokens it does not need. It removes
repeated concepts, sentences that restate their neighbours, verbose phrasing,
duplicate examples, meaningless formatting, and prose that does not affect the
output.

The important distinction is between removing words and removing behavior. A
shorter prompt that changes what the model does is a failed compression, not a
successful one. promptopt measures the result against the original intent
rather than treating a shorter string as automatically better, and the report
includes:

- original and compressed token estimates, and the reduction percentage
- a **semantic preservation** score (0–10): the model's confidence that
  behavior is unchanged
- **behavior risks**: specific ways the compressed prompt might behave
  differently, when the model is not fully confident

## When to use it

- A system prompt has grown past a few thousand tokens and most requests do
  not need all of it.
- You are hitting the context window and need headroom.
- You are paying per token at scale and the prompt is the fixed cost.

## When not to use it

- **The prompt is already tight.** compress will tell you no tokens were
  saved. Nothing is broken; there was just nothing to cut.
- **You want it clearer, not shorter.** Use [`optimize`](optimize.md).
- **You want to remove a specific instruction.** Edit the prompt. compress
  removes redundancy, not requirements.

## How it works

The prompt and your options (target, aggressive, preserve-behavior) go to the
model with the `prompts/compress/v1.md` instruction set, which asks for a JSON
object: the compressed prompt, the categories of cut it made, a
semantic-preservation score, and behavior risks. promptopt validates it and
estimates tokens for both versions.

If you set `--target` and the model reaches less than that, promptopt reports
the gap and says deeper cuts were judged unsafe. If semantic preservation
comes back below 7 with no risks listed, promptopt adds a warning telling you
to review the diff yourself.

## Flags

| Flag | Default | Effect |
| --- | --- | --- |
| `--target <n>` | 0 | Target reduction percentage, 0–95. `0` means "as much as is safe". |
| `--aggressive` | false | Allow dropping marginal context and collapsing examples harder. |
| `--preserve-behavior` | false | Refuse any cut not confidently safe, even if it misses the target. |
| `--output <file>`, `-o` | none | Write the compressed prompt to a file. |
| `--json` | none | Emit the `CompressResult` object. |
| `--quiet`, `-q` | none | Print only the compressed prompt. |

`--aggressive` and `--preserve-behavior` pull in opposite directions. If you
set both, `preserve-behavior` wins on any individual cut.

## Examples

```sh
promptopt compress system-prompt.md

# Aim for a 30% reduction
promptopt compress system.md --target 30

# Squeeze harder, accept some risk
promptopt compress system.md --target 50 --aggressive

# Only cuts the model is sure about
promptopt compress system.md --preserve-behavior

# Machine-readable, check the preservation score
promptopt compress system.md --json | jq '.semantic_preservation'
```

## Expected output

```
$ promptopt compress system-prompt.md --target 30

  promptopt  compress

  tokens  (estimate)
  2,431 → 1,487
  38.8% reduction

  changes
  + removed a constraint stated three times in different words
  + collapsed two near-identical examples into one
  + cut a paragraph explaining what JSON is

  semantic preservation
    preserved        ████████████████····  8.4

  compressed prompt

  ┌──────────────────────────────────────────────
  │ You are a support agent for Acme Cloud. ...
  └──────────────────────────────────────────────
```

## Common mistakes

- **Chasing the target number.** `--target 60` does not mean you will get 60%.
  It means promptopt tries, and stops when the next cut would risk behavior.
  The gap is reported, not hidden.
- **Ignoring the behavior risks list.** It is short and specific for a reason.
  If it says "the compressed prompt no longer tells the model to cite
  sources", that is a real change. Put it back if you need it.
- **Compressing before optimizing.** If the prompt is also unclear, optimize
  first. A clear prompt compresses better because redundancy is easier to
  spot.
