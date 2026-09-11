// Site-wide constants. The repo and domain do not necessarily exist yet;
// these are the single place to change them.

export const site = {
  name: "promptopt",
  // A concise technical positioning statement, not a slogan.
  description:
    "A command-line tool for prompt engineering. Inspect a prompt, see where it is weak or wasteful, and rewrite it without changing what it asks a model to do. Inference runs on Groq.",
  tagline: "Understand your prompts, then make them better.",
  url: "https://promptopt.dev",
  repo: "https://github.com/rakshit-gen/promptopt",
  repoShort: "rakshit-gen/promptopt",
  installCommand: "curl -fsSL https://promptopt.dev/install.sh | sh",
  keyExport: 'export GROQ_API_KEY="..."',
} as const;

export const nav = [
  { label: "Product", href: "/#product" },
  { label: "Commands", href: "/#commands" },
  { label: "Docs", href: "/docs" },
  { label: "GitHub", href: site.repo, external: true },
  { label: "Install", href: "/#install" },
] as const;

export const commands = [
  "optimize",
  "compress",
  "expand",
  "analyze",
  "transform",
  "eval",
] as const;

export type CommandName = (typeof commands)[number];
