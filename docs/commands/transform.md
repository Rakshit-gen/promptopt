# transform

```
promptopt transform [prompt|file|-] --to <target> [flags]
```

Change how a prompt is expressed without changing what it asks for.

## What it does

transform takes a prompt and rewrites it into a different structure or
representation. It preserves intent — the transformed prompt asks for the same
thing — but reorganizes it for a different use.

| `--to` | Result |
| --- | --- |
| `markdown` | Clean Markdown: headings, lists, and a fenced output section. Reorganized for readability. |
| `xml` | Semantic tags — `<role>`, `<task>`, `<context>`, `<constraints>`, `<output_format>`, `<examples>` — for content that exists. |
| `json` | A structured object with keys like `role`, `task`, `context`, `constraints` (array), `output_format`, `examples`. |
| `system` | A standing system prompt in the second person. One-off phrasing becomes a note about what user messages will contain. |
| `template` | `{{variables}}` in place of values that change between runs. |
| `agent` | A goal, available context, and explicit stopping conditions, for an autonomous agent. |

The architecture allows new targets without touching the CLI — see
[CONTRIBUTING.md](../../CONTRIBUTING.md#adding-a-transform-target).

## Template mode

`--to template` finds the values in a prompt that would realistically change
between runs and replaces them with `{{snake_case}}` variables. It
parameterizes inputs, target languages, domains, and file names — not every
noun.

```
Review this Python code for SQL injection vulnerabilities:

<code here>
```

becomes:

```
Review the following {{language}} code for {{security_issue}}:

{{code}}
```

The report lists every variable introduced. promptopt also scans the output
for `{{tokens}}` the model used but forgot to report, and merges both lists.

## When to use it

- Your codebase standardizes on XML-tagged or JSON prompts and you have a
  plain-text one.
- You wrote a one-off prompt and want a reusable template.
- You have a task prompt and need a system prompt for a chat loop.

## When not to use it

- **The prompt has real problems.** transform preserves them faithfully. Run
  [`analyze`](analyze.md) or [`optimize`](optimize.md) first.
- **You want every value parameterized.** template mode deliberately
  under-parameterizes. A template full of variables is hard to use.

## How it works

The prompt and target go to the model with the `prompts/transform/v1.md`
instruction set, which asks for a JSON object: the transformed prompt,
structural notes, the variable list (template only), and any warnings about
content that had to be approximated. promptopt validates it.

## Flags

| Flag | Default | Effect |
| --- | --- | --- |
| `--to <target>` | — (required) | One of the targets above. |
| `--output <file>`, `-o` | — | Write the transformed prompt to a file. |
| `--json` | — | Emit the `TransformResult` object. |
| `--quiet`, `-q` | — | Print only the transformed prompt. |

## Examples

```sh
promptopt transform prompt.txt --to xml
promptopt transform prompt.txt --to template
promptopt transform prompt.txt --to system -o system.txt
cat prompt.txt | promptopt transform --to json | jq .

# See which values got parameterized
promptopt transform prompt.txt --to template --json | jq '.variables'
```

## Expected output

```
$ promptopt transform prompt.txt --to xml

  promptopt  transform → xml

  notes
  · wrapped the persona in <role> and the request in <task>
  · lifted the three bullet constraints into <constraints>
  · no examples in the source, so no <examples> block

  xml

  ┌──────────────────────────────────────────────
  │ <role>You are a senior security engineer.</role>
  │ <task>Review the supplied code for injection flaws.</task>
  │ <constraints>
  │   <constraint>Only report exploitable issues.</constraint>
  │ </constraints>
  └──────────────────────────────────────────────
```

## Common mistakes

- **Forgetting `--to`.** transform will not guess a target; it errors and
  lists the options.
- **Expecting `--to json` to be promptopt's own JSON.** It is not. `--to json`
  produces a JSON *representation of the prompt*. `--json` produces the
  `TransformResult`. You can combine them: `--to json --json` gives you the
  result object with the JSON prompt as a string in `.result`.
- **Over-relying on template variables.** Check the list. If it parameterized
  something that should be fixed, edit it back to a literal.
