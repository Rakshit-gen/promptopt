"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { TerminalLineView } from "./Terminal";
import { commandDemos } from "@/lib/demo";
import type { TerminalLine } from "@/lib/demo";
import { commands } from "@/lib/site";

const HELP: TerminalLine[] = [
  { kind: "output", text: "promptopt — understand your prompts, then make them better", tone: "accent" },
  { kind: "gap" },
  { kind: "output", text: "Usage:", tone: "dim" },
  { kind: "output", text: "  promptopt <command> [prompt|file|-] [flags]" },
  { kind: "gap" },
  { kind: "output", text: "Commands:", tone: "dim" },
  { kind: "output", text: "  optimize    inspect a prompt and rewrite it to be clearer and tighter" },
  { kind: "output", text: "  compress    cut tokens without changing what the prompt makes a model do" },
  { kind: "output", text: "  expand      turn an underspecified prompt into a detailed one" },
  { kind: "output", text: "  analyze     score a prompt and report actionable findings" },
  { kind: "output", text: "  transform   convert a prompt to markdown, xml, json, system, template, agent" },
  { kind: "output", text: "  eval        assess a prompt and probe it with generated test cases" },
  { kind: "output", text: "  config      inspect configuration and write a starter file" },
  { kind: "gap" },
  { kind: "output", text: "Flags:", tone: "dim" },
  { kind: "output", text: "  --json         emit a stable JSON object instead of formatted text" },
  { kind: "output", text: "  --quiet, -q    print only the resulting prompt" },
  { kind: "output", text: "  --output, -o   write the result to a file" },
  { kind: "output", text: "  --model        Groq model to use" },
  { kind: "output", text: "  --no-color     disable ANSI color" },
  { kind: "gap" },
  { kind: "output", text: "This is a browser simulation. Try: analyze, compress, transform, help, clear", tone: "dim" },
];

const INTRO: TerminalLine[] = [
  { kind: "output", text: "promptopt playground — a simulation of the real CLI, running in your browser.", tone: "dim" },
  { kind: "output", text: "Type 'help' to start. Commands use the same example data as the docs.", tone: "dim" },
  { kind: "gap" },
];

type Entry = { input: string; output: TerminalLine[] };

export function Playground() {
  const [entries, setEntries] = useState<Entry[]>([]);
  const [value, setValue] = useState("");
  const [history, setHistory] = useState<string[]>([]);
  const [histIdx, setHistIdx] = useState(-1);
  const scrollRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const known = useMemo(
    () => ["help", "clear", ...commands, "config", "version"],
    [],
  );

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight });
  }, [entries]);

  function resolve(raw: string): TerminalLine[] | "clear" {
    const cmd = raw.trim().replace(/^promptopt\s+/, "");
    const [name] = cmd.split(/\s+/);

    if (!name) return [];
    if (name === "clear") return "clear";
    if (name === "help" || name === "--help" || name === "-h") return HELP;
    if (name === "version" || name === "--version") {
      return [{ kind: "output", text: "promptopt 0.1.0 (prompts v1)" }];
    }
    if (name === "config") {
      return [
        { kind: "output", text: "config file    (none; defaults + environment)", tone: "dim" },
        { kind: "output", text: "api key        not set" },
        { kind: "output", text: "model          openai/gpt-oss-120b" },
        { kind: "output", text: "output         text" },
        { kind: "output", text: "timeout        1m0s" },
      ];
    }
    if ((commands as readonly string[]).includes(name)) {
      const demo = commandDemos[name as (typeof commands)[number]];
      return [
        { kind: "output", text: `  (simulation) using example input: ${demo.input.label}`, tone: "dim" },
        { kind: "gap" },
        ...demo.output,
      ];
    }
    return [
      { kind: "output", text: `promptopt: unknown command "${name}"`, tone: "err" },
      { kind: "output", text: "run 'help' to see the six commands", tone: "dim" },
    ];
  }

  function submit(e: React.FormEvent) {
    e.preventDefault();
    const raw = value;
    if (!raw.trim()) return;

    const result = resolve(raw);
    setHistory((h) => [...h, raw]);
    setHistIdx(-1);
    setValue("");

    if (result === "clear") {
      setEntries([]);
      return;
    }
    setEntries((prev) => [...prev, { input: raw, output: result }]);
  }

  function onKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === "ArrowUp") {
      e.preventDefault();
      if (!history.length) return;
      const next = histIdx === -1 ? history.length - 1 : Math.max(0, histIdx - 1);
      setHistIdx(next);
      setValue(history[next]);
    } else if (e.key === "ArrowDown") {
      e.preventDefault();
      if (histIdx === -1) return;
      const next = histIdx + 1;
      if (next >= history.length) {
        setHistIdx(-1);
        setValue("");
      } else {
        setHistIdx(next);
        setValue(history[next]);
      }
    }
  }

  return (
    <div
      className="overflow-hidden rounded-xl border border-border bg-panel"
      onClick={() => inputRef.current?.focus()}
    >
      <div className="flex items-center gap-2 border-b border-border bg-raised/60 px-3.5 py-2.5">
        <span className="flex gap-1.5" aria-hidden>
          <span className="h-2.5 w-2.5 rounded-full bg-borderStrong" />
          <span className="h-2.5 w-2.5 rounded-full bg-borderStrong" />
          <span className="h-2.5 w-2.5 rounded-full bg-borderStrong" />
        </span>
        <span className="ml-1 font-mono text-2xs text-faint">promptopt — playground</span>
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            setEntries([]);
          }}
          className="ml-auto font-mono text-2xs text-faint hover:text-fg"
        >
          clear
        </button>
      </div>

      <div
        ref={scrollRef}
        className="h-[340px] overflow-y-auto p-4 font-mono text-[0.8rem] leading-6"
      >
        {INTRO.map((l, i) => (
          <TerminalLineView key={`intro-${i}`} line={l} />
        ))}
        {entries.map((entry, i) => (
          <div key={i}>
            <div className="whitespace-pre-wrap break-words text-fg">
              <span className="select-none text-accent">$ </span>
              {entry.input}
            </div>
            <div className="mb-3 mt-0.5">
              {entry.output.map((l, j) => (
                <TerminalLineView key={j} line={l} />
              ))}
            </div>
          </div>
        ))}

        <form onSubmit={submit} className="flex items-center">
          <span className="select-none text-accent" aria-hidden>
            ${" "}
          </span>
          <input
            ref={inputRef}
            value={value}
            onChange={(e) => setValue(e.target.value)}
            onKeyDown={onKeyDown}
            spellCheck={false}
            autoComplete="off"
            autoCapitalize="off"
            aria-label="playground command input"
            className="ml-1 flex-1 bg-transparent text-fg outline-none placeholder:text-faint"
            placeholder="analyze"
          />
        </form>
      </div>

      <div className="flex flex-wrap gap-1.5 border-t border-border px-4 py-2.5">
        {known.slice(0, 8).map((c) => (
          <button
            key={c}
            type="button"
            onClick={() => {
              setValue(c);
              inputRef.current?.focus();
            }}
            className="rounded border border-border px-2 py-0.5 font-mono text-2xs text-muted hover:border-borderStrong hover:text-fg"
          >
            {c}
          </button>
        ))}
      </div>
    </div>
  );
}
