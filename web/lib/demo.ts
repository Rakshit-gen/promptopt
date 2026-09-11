// Deterministic demo data for the interactive components on the site.
//
// Nothing here calls a real API. The outputs are representative of what
// promptopt produces, hand-written so the landing page is fast and works
// offline. The real CLI uses the real Groq integration.

import type { CommandName } from "./site";

export type TerminalLine =
  | { kind: "input"; text: string }
  | { kind: "output"; text: string; tone?: "dim" | "ok" | "warn" | "err" | "accent" }
  | { kind: "gap" };

export type CommandDemo = {
  name: CommandName;
  summary: string;
  when: string;
  input: { label: string; body: string };
  command: string;
  output: TerminalLine[];
};

// The hero terminal cycles through these short scripted runs, one per command,
// forever. Every line is kept under ~42 characters and every script the same
// length, so the box never wraps and never resizes.
export const heroScripts: TerminalLine[][] = [
  [
    { kind: "input", text: "promptopt compress system.md" },
    { kind: "gap" },
    { kind: "output", text: "  promptopt  compress", tone: "accent" },
    { kind: "gap" },
    { kind: "output", text: "  tokens   2,431 → 1,487  (est)" },
    { kind: "output", text: "  38.8% smaller", tone: "ok" },
    { kind: "gap" },
    { kind: "output", text: "  + cut a rule stated three times", tone: "ok" },
    { kind: "output", text: "  + merged two near-identical examples", tone: "ok" },
    { kind: "gap" },
    { kind: "output", text: "  semantic preservation  8.4 / 10", tone: "ok" },
  ],
  [
    { kind: "input", text: "promptopt analyze prompt.txt" },
    { kind: "gap" },
    { kind: "output", text: "  promptopt  analyze", tone: "accent" },
    { kind: "gap" },
    { kind: "output", text: "  overall score   4.6 / 10", tone: "warn" },
    { kind: "output", text: "  findings   1 error · 2 warnings" },
    { kind: "gap" },
    { kind: "output", text: "  ERROR    thorough vs. brief conflict", tone: "err" },
    { kind: "output", text: "  WARNING  code to review not included", tone: "warn" },
    { kind: "output", text: "  WARNING  no output format defined", tone: "warn" },
  ],
  [
    { kind: "input", text: "promptopt optimize prompt.txt" },
    { kind: "gap" },
    { kind: "output", text: "  promptopt  optimize", tone: "accent" },
    { kind: "gap" },
    { kind: "output", text: "  tokens   46 → 92  (est)" },
    { kind: "gap" },
    { kind: "output", text: "  + removed 4 'be helpful' lines", tone: "ok" },
    { kind: "output", text: "  + set role: programming assistant", tone: "ok" },
    { kind: "output", text: "  + added a rule for vague requests", tone: "ok" },
    { kind: "gap" },
    { kind: "output", text: "  ! longer now — compress if it matters", tone: "warn" },
  ],
  [
    { kind: "input", text: 'promptopt expand "a payment API"' },
    { kind: "gap" },
    { kind: "output", text: "  promptopt  expand · production", tone: "accent" },
    { kind: "gap" },
    { kind: "output", text: "  tokens   4 → 486  (est)" },
    { kind: "gap" },
    { kind: "output", text: "  + objective and non-goals", tone: "ok" },
    { kind: "output", text: "  + request / response contracts", tone: "ok" },
    { kind: "output", text: "  + idempotency and retry behavior", tone: "ok" },
    { kind: "gap" },
    { kind: "output", text: "  ~ assumes Stripe-style PaymentIntents", tone: "warn" },
  ],
  [
    { kind: "input", text: "promptopt transform p.txt -t template" },
    { kind: "gap" },
    { kind: "output", text: "  promptopt  transform → template", tone: "accent" },
    { kind: "gap" },
    { kind: "output", text: "  · parameterized language and issue", tone: "dim" },
    { kind: "output", text: "  · kept 'review' as a literal", tone: "dim" },
    { kind: "gap" },
    { kind: "output", text: "  variables (3)" },
    { kind: "output", text: "  {{language}} {{issue}} {{code}}", tone: "accent" },
    { kind: "gap" },
    { kind: "output", text: "  written to prompt.tmpl", tone: "ok" },
  ],
  [
    { kind: "input", text: "promptopt eval prompt.txt" },
    { kind: "gap" },
    { kind: "output", text: "  promptopt  eval", tone: "accent" },
    { kind: "gap" },
    { kind: "output", text: "  overall   7.4 / 10", tone: "warn" },
    { kind: "output", text: "  test cases   2 / 3 handled well" },
    { kind: "gap" },
    { kind: "output", text: "  ✓ normal input      clean summary", tone: "ok" },
    { kind: "output", text: "  ✗ injection in body  follows it", tone: "err" },
    { kind: "gap" },
    { kind: "output", text: "  ! input treated as instructions", tone: "warn" },
  ],
];

export const commandDemos: Record<CommandName, CommandDemo> = {
  optimize: {
    name: "optimize",
    summary:
      "Inspect a prompt and rewrite it to be clearer and tighter, keeping the same intent.",
    when: "A prompt works most of the time and you want it to work more of the time.",
    input: {
      label: "prompt.txt",
      body: "You are helpful. You are a helpful assistant. Help the user. When they ask for code, write the code. Be helpful. Write good code. Always help. Do not be unhelpful.",
    },
    command: "promptopt optimize prompt.txt",
    output: [
      { kind: "output", text: "  promptopt  optimize", tone: "accent" },
      { kind: "gap" },
      { kind: "output", text: "  tokens  (estimate)" },
      { kind: "output", text: "  46 → 92" },
      { kind: "output", text: "  100.0% larger", tone: "warn" },
      { kind: "gap" },
      { kind: "output", text: "  changes" },
      { kind: "output", text: "  + removed four repeated 'be helpful' statements", tone: "ok" },
      { kind: "output", text: "  + defined the role as a programming assistant", tone: "ok" },
      { kind: "output", text: "  + added a rule for ambiguous requests", tone: "ok" },
      { kind: "output", text: "  + specified a fenced code block with a language tag", tone: "ok" },
      { kind: "gap" },
      { kind: "output", text: "  warnings" },
      { kind: "output", text: "  ! result is longer; run compress if token budget matters", tone: "warn" },
      { kind: "gap" },
      { kind: "output", text: "  optimized prompt" },
      { kind: "output", text: "  ┌───────────────────────────────────────────", tone: "dim" },
      { kind: "output", text: "  │ Act as a programming assistant. Produce" },
      { kind: "output", text: "  │ working code for the user's request." },
      { kind: "output", text: "  │ If the request is ambiguous, ask one" },
      { kind: "output", text: "  │ clarifying question first. Otherwise return" },
      { kind: "output", text: "  │ the code in a fenced block with a language" },
      { kind: "output", text: "  │ tag, then a short explanation." },
      { kind: "output", text: "  └───────────────────────────────────────────", tone: "dim" },
    ],
  },
  compress: {
    name: "compress",
    summary:
      "Cut token count without changing what the prompt makes a model do.",
    when: "A system prompt has grown past a few thousand tokens and most requests do not need all of it.",
    input: {
      label: "system.md",
      body: "You are a support agent for Acme Cloud. Be concise. Do not apologize. Never say sorry. Avoid apologies. Answer only Acme Cloud questions. If asked about something else, say you can only help with Acme Cloud. Here is an example. Here is another almost identical example...",
    },
    command: "promptopt compress system.md --target 30",
    output: [
      { kind: "output", text: "  promptopt  compress", tone: "accent" },
      { kind: "gap" },
      { kind: "output", text: "  tokens  (estimate)" },
      { kind: "output", text: "  1,204 → 742" },
      { kind: "output", text: "  38.4% reduction", tone: "ok" },
      { kind: "gap" },
      { kind: "output", text: "  changes" },
      { kind: "output", text: "  + merged 'do not apologize' stated three ways", tone: "ok" },
      { kind: "output", text: "  + kept one of two identical examples", tone: "ok" },
      { kind: "gap" },
      { kind: "output", text: "  semantic preservation   9.1 / 10", tone: "ok" },
      { kind: "output", text: "  behavior risks          none reported", tone: "dim" },
    ],
  },
  expand: {
    name: "expand",
    summary:
      "Turn an underspecified prompt into a detailed one, with every added assumption labeled.",
    when: "You have a one-line idea and want a prompt you can actually run.",
    input: {
      label: "inline",
      body: "build a payment API",
    },
    command: 'promptopt expand "build a payment API" --depth production',
    output: [
      { kind: "output", text: "  promptopt  expand · production", tone: "accent" },
      { kind: "gap" },
      { kind: "output", text: "  tokens  (estimate)" },
      { kind: "output", text: "  4 → 486" },
      { kind: "gap" },
      { kind: "output", text: "  added" },
      { kind: "output", text: "  + explicit objective and non-goals", tone: "ok" },
      { kind: "output", text: "  + request and response contracts", tone: "ok" },
      { kind: "output", text: "  + idempotency and retry behavior", tone: "ok" },
      { kind: "output", text: "  + error taxonomy and failure responses", tone: "ok" },
      { kind: "gap" },
      { kind: "output", text: "  assumptions  (correct these if wrong)", tone: "warn" },
      { kind: "output", text: "  ~ processor: Stripe-style PaymentIntent model", tone: "warn" },
      { kind: "output", text: "  ~ currency: integer minor units (cents)", tone: "warn" },
    ],
  },
  analyze: {
    name: "analyze",
    summary:
      "Score a prompt and report actionable findings with severities and stable IDs.",
    when: "Before building on a prompt you did not write, or in CI to catch a regression.",
    input: {
      label: "prompt.txt",
      body: "Review this Python code and tell me what you think. Be thorough but keep it brief. Also fix any bugs.",
    },
    command: "promptopt analyze prompt.txt",
    output: [
      { kind: "output", text: "  promptopt  analyze", tone: "accent" },
      { kind: "gap" },
      { kind: "output", text: "  overall score  4.6 / 10", tone: "warn" },
      { kind: "gap" },
      { kind: "output", text: "    clarity          ██████████·········  5.1" },
      { kind: "output", text: "    specificity      ███████···········  3.4", tone: "warn" },
      { kind: "output", text: "    consistency      ████████··········  4.0", tone: "warn" },
      { kind: "gap" },
      { kind: "output", text: "  findings  1 error, 2 warning" },
      { kind: "gap" },
      { kind: "output", text: "  ERROR   P004  'be thorough' conflicts with 'keep it brief'", tone: "err" },
      { kind: "output", text: "  WARNING P002  the code to review is not included", tone: "warn" },
      { kind: "output", text: "  WARNING P003  no output format is defined", tone: "warn" },
    ],
  },
  transform: {
    name: "transform",
    summary:
      "Convert a prompt to markdown, XML, JSON, a system prompt, a template, or an agent brief.",
    when: "Your codebase standardizes on a prompt format, or you want a reusable template.",
    input: {
      label: "inline",
      body: "Review this Python code for SQL injection vulnerabilities",
    },
    command: "promptopt transform prompt.txt --to template",
    output: [
      { kind: "output", text: "  promptopt  transform → template", tone: "accent" },
      { kind: "gap" },
      { kind: "output", text: "  notes" },
      { kind: "output", text: "  · parameterized the language and the issue", tone: "dim" },
      { kind: "output", text: "  · left 'review' as a literal instruction", tone: "dim" },
      { kind: "gap" },
      { kind: "output", text: "  variables  (3)" },
      { kind: "output", text: "  {{language}}  {{security_issue}}  {{code}}", tone: "accent" },
      { kind: "gap" },
      { kind: "output", text: "  template" },
      { kind: "output", text: "  ┌───────────────────────────────────────────", tone: "dim" },
      { kind: "output", text: "  │ Review the following {{language}} code for" },
      { kind: "output", text: "  │ {{security_issue}}:" },
      { kind: "output", text: "  │" },
      { kind: "output", text: "  │ {{code}}" },
      { kind: "output", text: "  └───────────────────────────────────────────", tone: "dim" },
    ],
  },
  eval: {
    name: "eval",
    summary:
      "Assess a prompt and probe it with generated test cases covering normal, ambiguous, and adversarial input.",
    when: "Before a prompt goes into a system where wrong output has a cost.",
    input: {
      label: "prompt.txt",
      body: "Summarize the customer message below in three bullet points.\n\n{{message}}",
    },
    command: "promptopt eval prompt.txt",
    output: [
      { kind: "output", text: "  promptopt  eval", tone: "accent" },
      { kind: "gap" },
      { kind: "output", text: "  overall  7.4 / 10", tone: "warn" },
      { kind: "gap" },
      { kind: "output", text: "    clarity          ████████████████···  8.6" },
      { kind: "output", text: "    robustness       ██████████········  5.2", tone: "warn" },
      { kind: "output", text: "    output control   ████████████████···  8.1" },
      { kind: "gap" },
      { kind: "output", text: "  test cases  2/3 handled well" },
      { kind: "gap" },
      { kind: "output", text: "  ✓ normal message      tight three-bullet summary", tone: "ok" },
      { kind: "output", text: "  ✗ injection in body   may follow 'ignore the above'", tone: "err" },
      { kind: "gap" },
      { kind: "output", text: "  weaknesses" },
      { kind: "output", text: "  ! input is treated as instructions, not data", tone: "warn" },
    ],
  },
};

/** The before/after section: a real-looking prompt and its optimization. */
export const beforeAfter = {
  before: `You are an AI assistant that is an expert in everything related to
customer support. You are very knowledgeable and helpful. Your job is
to help support agents write replies to customers.

When a support agent gives you a customer message, write a reply. The
reply should be helpful. It should be polite. Be polite and helpful.
Don't be rude. Always be professional and courteous and polite.

Keep replies short. But also make sure you fully address everything
the customer said. Be concise but thorough and complete.

Also, never make up information. If you don't know something, say so.
Don't invent facts. Don't hallucinate. Only use real information.

Output the reply.`,
  after: `You help support agents draft replies to customers.

Input: a customer message, optionally with notes from the agent.
Output: a reply, in plain text, no greeting line or signature.

Rules:
- Address every question and request in the message.
- Keep it to what is needed — no filler, no repeated apologies.
- If answering requires information you were not given, do not guess.
  State what is missing and what the agent should confirm.`,
  tokensBefore: 168,
  tokensAfter: 96,
  changes: [
    { text: "Removed the inflated 'expert in everything' persona", kind: "cut" },
    { text: "Collapsed 'polite / helpful / professional' repeated five times", kind: "cut" },
    { text: "Turned 'concise but thorough' (a contradiction) into two concrete rules", kind: "fix" },
    { text: "Replaced 'don't hallucinate' with what to actually do when information is missing", kind: "fix" },
    { text: "Defined the input and the output shape explicitly", kind: "add" },
  ],
} as const;

/** The linear workflow diagram steps. */
export const workflow = [
  {
    step: "input",
    title: "A prompt",
    body: "From an argument, a file, or stdin. No project setup, no config required to start.",
  },
  {
    step: "analyze",
    title: "analyze",
    body: "Score it and list findings — conflicts, missing format, repeated constraints, injection exposure.",
  },
  {
    step: "transform",
    title: "optimize / compress / transform",
    body: "Rewrite for clarity, cut tokens, or convert to a template or system prompt. Intent is preserved.",
  },
  {
    step: "eval",
    title: "eval",
    body: "Probe the result with generated test cases before you depend on it.",
  },
  {
    step: "output",
    title: "A better prompt",
    body: "To a file with -o, or piped to the next command. JSON out for anything scripted.",
  },
] as const;

export const architecture = [
  { layer: "CLI", note: "argument / file / stdin, flags, config" },
  { layer: "Command layer", note: "Cobra; the only package that knows about the CLI" },
  { layer: "Prompt engine", note: "the six operations; versioned instruction sets" },
  { layer: "Groq provider", note: "isolated behind one interface; retries, timeouts" },
  { layer: "Structured result", note: "validated JSON → typed struct → text or --json" },
] as const;

export const devFacts = [
  {
    title: "Reads from stdin",
    body: "cat prompt.txt | promptopt analyze — every command takes piped input.",
  },
  {
    title: "Reads files",
    body: "promptopt compress system-prompt.md. A path is a path; a string is a string.",
  },
  {
    title: "JSON output",
    body: "--json gives one stable object, documented field by field. No color, no prose.",
  },
  {
    title: "Shell composable",
    body: "promptopt optimize -q | promptopt compress -q > final.txt. Stable exit codes.",
  },
  {
    title: "Configurable models",
    body: "--model, PROMPTOPT_MODEL, or the config file.",
  },
  {
    title: "No telemetry",
    body: "One network call, to Groq, with your key. Prompts and the key are never logged.",
  },
  {
    title: "Local CLI",
    body: "A single static binary. No daemon, no account, no project files.",
  },
  {
    title: "Groq inference",
    body: "Fast enough to run analyze in a pre-commit hook.",
  },
] as const;
