# Evaluation

This page covers what `promptopt eval` does today and how its API is shaped
for what comes next. For the command reference, see [commands/eval](commands/eval.md).

## Today

eval does one thing: it assesses a prompt and reasons through generated test
cases. One model call produces:

- four sub-scores: clarity, consistency, robustness, output control
- 4–6 generated test cases covering a normal input, an ambiguous input, an
  adversarial/out-of-scope input, and an edge case
- for each case: the input, what a good prompt should do with it, and an
  assessment of what this prompt would actually do, with a coarse pass/fail
- a weaknesses list and a summary

The pass/fail on each case is the model's judgement about the prompt's
wording. eval does **not** run those inputs through a model and check the
output. That distinction matters: eval is a structured review, not a test
harness.

## Why ship this first

A full evaluation platform (datasets, assertions, model comparison,
regression tracking) is a large surface, and most of it is only useful once
you have prompts stable enough to benchmark. The common early need is simpler:
"is this prompt obviously going to break, and where?" eval answers that in one
call.

## The API is built to grow

The result shape (`pkg/types.EvalResult`) already carries the pieces a fuller
system needs:

- `scores` is the shared `ScoreCard`, so a scored eval and a scored analyze
  are comparable.
- `test_cases[].pass` is a `*bool`: nil today when eval only reviewed the
  case, a real verdict once cases are executed.
- `test_cases[].input` and `.expectation` are already the shape of a runnable
  assertion.

Planned additions, in rough order, none of which will break the JSON contract:

1. **Execute generated cases**: run each `input` through the prompt on a
   model and judge the output against `expectation`. `pass` becomes real.
2. **User datasets**: `promptopt eval prompt.txt --dataset cases.jsonl`,
   where each line is `{"input": ..., "expect": ...}`.
3. **Assertions**: declarative checks (`contains`, `matches`, `json_schema`,
   `not_contains`) attached to a case.
4. **Model comparison**: `--models a,b,c`, one column per model.
5. **Regression**: store a baseline, diff against it, fail on a drop.

## Using eval now

The high-value uses today:

```sh
# Before a prompt goes into production
promptopt eval prompt.txt

# After optimize or compress, to confirm behavior held up
promptopt optimize prompt.txt -o prompt.v2.txt
promptopt eval prompt.v2.txt

# Find the cases the prompt handles badly
promptopt eval prompt.txt --json \
  | jq -r '.test_cases[] | select(.pass == false) | "- \(.name): \(.assessment)"'
```

## eval versus analyze

| | analyze | eval |
| --- | --- | --- |
| Approach | Reasons about the prompt text | Probes with concrete inputs |
| Dimensions | 6 (adds specificity, completeness, efficiency) | 4 |
| Output | Findings with severities and IDs | Test cases with assessments |
| Best at | "What is written wrong" | "How it breaks under pressure" |

Run analyze for the checklist, eval for the stress test. They overlap on
robustness and consistency; that overlap is a feature: if both flag the same
thing, it is real.
