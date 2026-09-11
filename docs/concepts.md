# How promptopt works

## The model does the work; promptopt does the plumbing

Every operation is one call to a language model on Groq, with a versioned
instruction set that asks for a structured JSON response. promptopt's job is
everything around that call: resolving input, assembling the request,
validating the response, estimating tokens, rendering the result, and turning
failures into messages you can act on.

```
your prompt ──▶ internal/cli ──▶ internal/engine ──▶ internal/groq ──▶ Groq API
                                       │
                              prompts/<op>/v1.md  (the instructions)
                                       │
                                       ▼
                          validated JSON ──▶ pkg/types result ──▶ output
```

The instruction sets live in `prompts/<op>/v1.md` in the repository and are
embedded into the binary with `go:embed`. They are version-controlled, so any
change to how an operation behaves is a reviewable diff. Their YAML front
matter carries a version; a significant change gets a new file (`v2.md`)
rather than an in-place edit.

## Structured output, validated

promptopt asks Groq for a JSON object (`response_format: json_object`) and
parses it into a typed Go struct. It does not regex LLM prose. The parser:

- tolerates a model that wraps the object in ```` ```json ```` fences or emits
  reasoning text before the answer — it extracts the first balanced top-level
  object
- validates the object against the expected shape
- clamps scores into 0–10 and drops empty list entries
- fails with exit code 8 (not a crash, not a guess) if the response cannot be
  parsed

One malformed response is an error you can retry, never a panic.

## Token counting

Token counts are labeled `(estimate)`. promptopt does not ship the exact BPE
tables for every Groq model, so it uses a heuristic that blends a
character-based ratio with a word-and-symbol count. It tracks real tokenizers
to within a few percent for prompt-like text and errs slightly high so budgets
are not underestimated.

The `tokenizer.Estimator` interface exists so a model-exact tokenizer can be
added later without changing any call site. When that lands, the source label
changes from `estimate` to the tokenizer's name, and `TokenReport.source` in
the JSON output reflects it.

Groq also returns real token counts for the request and response in its
`usage` object. promptopt reads those but does not substitute them for the
before/after estimate, because they describe the whole HTTP exchange (system
instructions included), not the prompt text on its own.

## Intent preservation

`optimize`, `compress`, and `transform` are all constrained to keep what the
prompt asks for. The instruction sets say so explicitly, and each command's
report surfaces the risk:

- `optimize` lists every change and warns when the result grew significantly
- `compress` reports a semantic-preservation score and a behavior-risks list
- `transform` warns about anything it had to approximate

`expand` is the exception: its whole job is to add. It handles the risk by
labeling every addition that was not in your original as an assumption.

## Provider isolation

`internal/groq` is the only package that knows about Groq. `internal/engine`
depends on a `groq.Completer` interface with one method. That means:

- every operation is unit-tested with a fake completer and no network
- adding another provider is implementing one interface, not editing business
  logic

## What promptopt does not do

- It does not store your prompts anywhere.
- It does not log prompts or the API key.
- It does not phone home. The only network call is to the Groq API.
- It does not cache results. Each run is a fresh call.
