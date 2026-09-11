# Security policy

## Reporting a vulnerability

Email **security@promptopt.dev** with a description of the issue and, if you
can, a minimal reproduction. Please do not open a public issue for
security-sensitive reports.

You can expect an acknowledgement within a few days. If the report is valid,
we will work on a fix and coordinate a disclosure timeline with you, and
credit you in the release notes unless you prefer otherwise.

## What promptopt does with your data

- **Prompts are not logged.** They are held in memory for the duration of a
  command and sent to the Groq API for inference. promptopt does not write
  them to disk, and does not send them anywhere except Groq.
- **The API key is not logged.** It is read from the environment (or a file
  you point at) and sent only in the `Authorization` header of requests to
  Groq. `promptopt config show` masks it.
- **No telemetry.** promptopt makes no network calls other than to the Groq
  API. There is no usage reporting, no crash reporting, no update check.
- Groq's processing of prompt data is governed by Groq's own terms and
  privacy policy.

## Scope

The prompt-injection analysis in `promptopt analyze` is a heuristic signal
produced by a language model. It can miss real problems and flag harmless
ones. Do not treat a clean `analyze` result as a security guarantee for a
prompt that will run with elevated privileges or tool access.

## Supported versions

Security fixes are applied to the latest minor release. Older versions are not
maintained.
