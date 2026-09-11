# Getting started

## 1. Install

```sh
curl -fsSL https://promptopt.dev/install.sh | sh
```

This installs the `promptopt` binary into `~/.local/bin`. If that directory is
not on your `PATH`, the installer prints the line to add. See
[Installation](installation.md) for other methods.

Check it:

```sh
promptopt --version
```

## 2. Set your Groq API key

promptopt runs inference on [Groq](https://console.groq.com). Create a key at
<https://console.groq.com/keys> and export it:

```sh
export GROQ_API_KEY="your-key"
```

Put that line in your shell profile (`~/.zshrc`, `~/.bashrc`) so it persists.
promptopt never writes the key to disk and never logs it.

## 3. Run a command

Every command takes a prompt as an argument, a file, or on stdin.

```sh
# Inline
promptopt analyze "write a function that reverses a string"

# From a file
promptopt analyze prompt.txt

# From stdin
cat prompt.txt | promptopt analyze
```

A first useful workflow is to analyze a prompt you already have, then optimize
it:

```sh
promptopt analyze system-prompt.md
promptopt optimize system-prompt.md -o system-prompt.optimized.md
```

`analyze` tells you what is wrong. `optimize` fixes what it can and reports
what it changed. The `-o` flag writes the rewritten prompt to a file; the
report still goes to your terminal.

## 4. Use it in a pipeline

Every command supports `--json`, which prints a single stable object and
nothing else — no color, no headings.

```sh
promptopt analyze prompt.txt --json | jq '.scores.overall'
```

`--quiet` prints only the resulting prompt, which is what you want when
chaining:

```sh
cat draft.txt | promptopt optimize -q | promptopt compress -q > final.txt
```

## Next

- [Commands overview](commands/README.md) — what each of the six does and when
- [Configuration](configuration.md) — defaults, models, timeouts
- [Using promptopt in CI](ci.md) — gate a prompt's quality in a pipeline
