import Link from "next/link";
import { Nav } from "@/components/Nav";
import { Footer } from "@/components/Footer";
import { Section } from "@/components/Section";
import { Reveal } from "@/components/Reveal";
import { TypingTerminal } from "@/components/TypingTerminal";
import { CommandExplorer } from "@/components/CommandExplorer";
import { BeforeAfter } from "@/components/BeforeAfter";
import { Workflow } from "@/components/Workflow";
import { InstallBlock } from "@/components/InstallBlock";
import { Architecture } from "@/components/Architecture";
import { ShaderBackground } from "@/components/ShaderBackground";
import { CopyButton } from "@/components/CopyButton";
import { heroScripts, devFacts } from "@/lib/demo";
import { site } from "@/lib/site";

export default function HomePage() {
  return (
    <>
      <Nav />
      <main id="main">
        {/* 1 · Hero */}
        <section className="relative flex min-h-[calc(100svh-3.5rem)] items-center overflow-hidden border-b border-border bg-bg">
          <ShaderBackground />
          <div className="grid-bg pointer-events-none absolute inset-0 opacity-40" aria-hidden />
          <div
            className="pointer-events-none absolute inset-x-0 bottom-0 h-32 bg-gradient-to-b from-transparent to-bg"
            aria-hidden
          />
          <div className="container-content relative grid w-full gap-12 py-16 lg:grid-cols-[1.05fr_1fr] lg:items-center">
            <div>
              <h1 className="text-balance text-4xl font-semibold leading-[1.1] tracking-tight sm:text-5xl">
                Work on your prompts
                <br />
                from the terminal.
              </h1>
              <p className="mt-5 max-w-md text-pretty text-base text-muted sm:text-lg">
                Analyze, optimize, compress, expand, transform, and evaluate the
                prompts you send to language models.
              </p>
              <div className="mt-7 flex flex-col gap-3 sm:flex-row">
                <a
                  href="#install"
                  className="inline-flex items-center justify-center rounded-md border border-accent/40 bg-accent/10 px-4 py-2.5 text-sm font-medium text-fg transition-colors hover:bg-accent/15"
                >
                  Install CLI
                </a>
                <Link
                  href="/docs"
                  className="inline-flex items-center justify-center rounded-md border border-border bg-panel px-4 py-2.5 text-sm text-muted transition-colors hover:text-fg"
                >
                  Read docs
                </Link>
              </div>
              <div className="mt-6 flex items-center gap-3 rounded-md border border-border bg-panel/80 px-3 py-2 font-mono text-2xs backdrop-blur sm:text-xs">
                <code className="truncate text-fg">
                  <span className="text-accent">$ </span>
                  {site.installCommand}
                </code>
                <CopyButton value={site.installCommand} label="" className="ml-auto shrink-0" />
              </div>
            </div>

            <Reveal delay={0.1}>
              <TypingTerminal scripts={heroScripts} />
            </Reveal>
          </div>
        </section>

        {/* 2 · Product / what it is */}
        <Section
          id="product"
          eyebrow="what it is"
          title="A linter and a rewriter for prompts"
          intro="Prompts drift. This is what you run to catch it."
        >
          <div className="grid gap-4 sm:grid-cols-3">
            {[
              {
                h: "Find wasted tokens",
                p: "analyze points at the repeated constraint or the inert persona line.",
              },
              {
                h: "Compress, keep constraints",
                p: "compress cuts a long system prompt and rates its confidence that behavior held.",
              },
              {
                h: "Compose like a Unix tool",
                p: "stdin in, JSON out, stable exit codes.",
              },
            ].map((c, i) => (
              <Reveal
                key={c.h}
                delay={i * 0.07}
                className="group rounded-lg border border-border bg-panel p-5 transition-colors hover:border-borderStrong"
              >
                <h3 className="text-sm font-medium text-fg">
                  <span className="text-accent transition-[margin] group-hover:mr-1">
                    →
                  </span>{" "}
                  {c.h}
                </h3>
                <p className="mt-2 text-sm text-muted">{c.p}</p>
              </Reveal>
            ))}
          </div>
        </Section>

        {/* 3 · Command explorer */}
        <Section
          id="commands"
          eyebrow="six commands"
          title="One for each job"
          intro="The whole public surface. Pick one."
        >
          <CommandExplorer />
        </Section>

        {/* 4 · Before / after */}
        <Section
          eyebrow="optimize, concretely"
          title="Same intent, less noise"
          intro="Toggle to compare."
        >
          <BeforeAfter />
        </Section>

        {/* 5 · Workflow */}
        <Section
          eyebrow="a workflow"
          title="Rough prompt in, dependable prompt out"
        >
          <Workflow />
        </Section>

        {/* 6 · Install */}
        <Section
          id="install"
          eyebrow="install"
          title="One command, no dependencies"
          intro="macOS and Linux, arm64 and amd64. One binary to ~/.local/bin."
        >
          <div className="grid gap-6 lg:grid-cols-[1fr_300px]">
            <InstallBlock />
            <div className="space-y-4 text-sm text-muted">
              <p>
                <code className="text-fg">go install</code> works today; a
                Homebrew tap ships with the first tagged release.
              </p>
              <p>
                You need a Groq API key from{" "}
                <a
                  href="https://console.groq.com/keys"
                  target="_blank"
                  rel="noreferrer"
                  className="text-accent hover:underline"
                >
                  console.groq.com/keys
                </a>
                . promptopt never writes it to disk and never logs it.
              </p>
            </div>
          </div>
        </Section>

        {/* 8 · Architecture */}
        <Section
          id="architecture"
          eyebrow="how it is built"
          title="A thin CLI over a testable engine"
          intro="Everything below the command layer has no idea the CLI exists."
        >
          <Architecture />
        </Section>

        {/* 9 · Developer-first */}
        <Section eyebrow="built for a shell" title="Behaves like the rest of your pipeline">
          <div className="grid gap-px overflow-hidden rounded-lg border border-border bg-border sm:grid-cols-2 lg:grid-cols-4">
            {devFacts.map((f, i) => (
              <Reveal
                key={f.title}
                delay={(i % 4) * 0.05}
                className="bg-panel p-5 transition-colors hover:bg-raised"
              >
                <h3 className="font-mono text-sm text-fg">{f.title}</h3>
                <p className="mt-2 text-sm text-muted">{f.body}</p>
              </Reveal>
            ))}
          </div>
        </Section>

        {/* 10 · Final CTA */}
        <section className="border-t border-border py-24">
          <Reveal className="container-content text-center">
            <h2 className="text-2xl font-semibold tracking-tight sm:text-3xl">
              Install promptopt
            </h2>
            <p className="mx-auto mt-3 max-w-md text-[0.95rem] text-muted">
              One binary, one API key.
            </p>
            <div className="mx-auto mt-6 flex max-w-lg items-center gap-3 rounded-md border border-border bg-panel px-3 py-2.5 font-mono text-xs">
              <code className="truncate text-fg">
                <span className="text-accent">$ </span>
                {site.installCommand}
              </code>
              <CopyButton value={site.installCommand} label="" className="ml-auto shrink-0" />
            </div>
            <div className="mt-6 flex justify-center gap-3">
              <Link
                href="/docs/getting-started"
                className="rounded-md border border-border bg-panel px-4 py-2 text-sm text-muted hover:text-fg"
              >
                Getting started
              </Link>
              <a
                href={site.repo}
                target="_blank"
                rel="noreferrer"
                className="rounded-md border border-border bg-panel px-4 py-2 text-sm text-muted hover:text-fg"
              >
                Star on GitHub
              </a>
            </div>
          </Reveal>
        </section>
      </main>
      <Footer />
    </>
  );
}
