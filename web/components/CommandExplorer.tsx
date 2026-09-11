"use client";

import { useState } from "react";
import { AnimatePresence, motion, useReducedMotion } from "framer-motion";
import { TerminalChrome, TerminalTranscript } from "./Terminal";
import { CopyButton } from "./CopyButton";
import { commandDemos } from "@/lib/demo";
import { commands, type CommandName } from "@/lib/site";
import Link from "next/link";

export function CommandExplorer() {
  const [active, setActive] = useState<CommandName>("optimize");
  const reduce = useReducedMotion();
  const demo = commandDemos[active];

  return (
    <div className="grid gap-6 lg:grid-cols-[200px_1fr]">
      <div
        className="flex gap-1.5 overflow-x-auto lg:flex-col lg:overflow-visible"
        role="tablist"
        aria-label="Commands"
      >
        {commands.map((name) => {
          const selected = name === active;
          return (
            <button
              key={name}
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(name)}
              className={`shrink-0 rounded-md border px-3 py-2 text-left font-mono text-sm transition-colors ${
                selected
                  ? "border-borderStrong bg-raised text-fg"
                  : "border-transparent text-muted hover:bg-panel hover:text-fg"
              }`}
            >
              <span className="text-faint">promptopt </span>
              {name}
            </button>
          );
        })}
      </div>

      <div>
        <AnimatePresence mode="wait">
          <motion.div
            key={active}
            initial={reduce ? false : { opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={reduce ? undefined : { opacity: 0, y: -8 }}
            transition={{ duration: 0.22 }}
          >
            <div className="min-h-[4.5rem] sm:min-h-[4rem]">
              <p className="text-[0.95rem] text-fg/90">{demo.summary}</p>
              <p className="mt-2 text-sm text-muted">
                <span className="text-faint">When: </span>
                {demo.when}
              </p>
            </div>

            <div className="mt-4 grid gap-4 lg:grid-cols-2">
              <div className="flex flex-col overflow-hidden rounded-lg border border-border bg-panel">
                <div className="flex items-center justify-between border-b border-border px-3 py-2">
                  <span className="font-mono text-2xs text-faint">
                    input · {demo.input.label}
                  </span>
                </div>
                <pre className="h-64 overflow-auto p-3 font-mono text-[0.78rem] leading-6 text-fg/80">
                  {demo.input.body}
                </pre>
              </div>

              <div className="flex flex-col gap-2">
                <div className="flex min-h-[2.75rem] items-center justify-between gap-2 rounded-lg border border-border bg-raised px-3 py-2">
                  <code className="break-all font-mono text-[0.78rem] text-fg">
                    <span className="text-accent">$ </span>
                    {demo.command}
                  </code>
                  <CopyButton value={demo.command} label="" className="shrink-0" />
                </div>
                <TerminalChrome title="output">
                  <div className="h-64 overflow-y-auto">
                    <TerminalTranscript lines={demo.output} />
                  </div>
                </TerminalChrome>
              </div>
            </div>

            <Link
              href={`/docs/commands/${active}`}
              className="mt-4 inline-flex items-center gap-1 text-sm text-accent hover:underline"
            >
              {active} reference
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden>
                <path
                  d="M3 9L9 3M9 3H4M9 3V8"
                  stroke="currentColor"
                  strokeWidth="1.4"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </Link>
          </motion.div>
        </AnimatePresence>
      </div>
    </div>
  );
}
