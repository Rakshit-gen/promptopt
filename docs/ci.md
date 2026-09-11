# Using promptopt in CI

promptopt is useful in a pipeline for one thing: catching a prompt regressing
before it ships. You keep prompts in the repo, and a job checks them the same
way a linter checks code.

## Exit codes

promptopt exits non-zero only on a real failure: a bad flag, a missing key, a
provider error. **A low quality score is not a failure.** `analyze` and `eval`
exit 0 with a score of 2.0.

That is deliberate: whether a 6.5 is acceptable depends on the prompt. You set
the threshold, in the pipeline, by reading the JSON.

| Code | Meaning |
| --- | --- |
| 0 | Success |
| 2 | Usage error |
| 3 | Config error |
| 4 | Missing/rejected API key |
| 5 | Rate limited |
| 6 | Prompt exceeds the model's context window |
| 7 | Provider error |
| 8 | Unparseable model response |

## GitHub Actions

Store the key as a repository secret named `GROQ_API_KEY`.

```yaml
name: prompts

on:
  pull_request:
    paths:
      - "prompts/**"

jobs:
  check:
    runs-on: ubuntu-latest
    env:
      GROQ_API_KEY: ${{ secrets.GROQ_API_KEY }}
    steps:
      - uses: actions/checkout@v4

      - name: Install promptopt
        run: curl -fsSL https://promptopt.dev/install.sh | sh

      - name: Analyze prompts
        run: |
          set -euo pipefail
          fail=0
          for f in prompts/*.md; do
            echo "::group::$f"
            result=$("$HOME/.local/bin/promptopt" analyze "$f" --json)
            echo "$result" | jq '{overall: .scores.overall, errors: [.findings[] | select(.severity=="ERROR") | .title]}'

            errors=$(echo "$result" | jq '[.findings[] | select(.severity=="ERROR")] | length')
            overall=$(echo "$result" | jq '.scores.overall')

            if [ "$errors" -gt 0 ]; then
              echo "::error file=$f::$errors ERROR-level finding(s)"
              fail=1
            fi
            if awk "BEGIN { exit !($overall < 7.0) }"; then
              echo "::error file=$f::overall score $overall is below 7.0"
              fail=1
            fi
            echo "::endgroup::"
          done
          exit $fail
```

## Guarding against regression

Commit a baseline score with each prompt and compare on PRs:

```sh
# One-time: record the baseline
promptopt analyze prompts/support.md --json | jq '.scores.overall' > prompts/support.score

# In CI: fail if it dropped by more than 0.5
baseline=$(cat prompts/support.score)
current=$(promptopt analyze prompts/support.md --json | jq '.scores.overall')
awk "BEGIN { exit !($current >= $baseline - 0.5) }" || {
  echo "support.md score dropped: $baseline -> $current"
  exit 1
}
```

## Rate limits

A CI job that checks many prompts can hit Groq's rate limit (exit 5). Options:

- Run the checks serially with a short `sleep` between them.
- Raise `PROMPTOPT_MAX_RETRIES`; promptopt honours `Retry-After`.
- Only check the prompts that changed:

```yaml
      - name: Analyze changed prompts
        run: |
          git fetch origin ${{ github.base_ref }}
          git diff --name-only origin/${{ github.base_ref }}... -- 'prompts/*.md' \
            | while read -r f; do promptopt analyze "$f" --json | jq -e '.scores.overall >= 7'; done
```

## Keeping the cost down

Each check is one model call. `openai/gpt-oss-20b` is faster and cheaper than
the default and is fine for CI gating:

```sh
promptopt analyze prompts/support.md --model openai/gpt-oss-20b --json
```

## What not to do

- Do not run `optimize`/`compress`/`expand` in CI and commit the result
  automatically. Those rewrite the prompt; a human should review the rewrite.
- Do not gate on `eval`'s pass count as if it were a test result: see
  [evaluation](evaluation.md).
