"use client";

import { useEffect, useMemo, useState } from "react";
import { useReducedMotion } from "framer-motion";
import { TerminalChrome, TerminalLineView } from "./Terminal";
import type { TerminalLine } from "@/lib/demo";

// A terminal that cycles through a list of scripted runs on its own: each run
// types its command, prints the output line by line, holds, then the next
// begins. It is also controllable: the command rail above it jumps straight
// to a run, and hovering the terminal pauses the auto-advance. The transcript
// area is a fixed height so the box never resizes between commands. Under
// prefers-reduced-motion the first run shows in full and nothing animates.
//
// Three effects, each owning its own timer list and clearing only its own:
// a shared bucket previously let the print effect's cleanup cancel the typing
// effect's freshly-scheduled timer, which stalled every run after the first.

type Phase = "typing" | "printing" | "done";

const HOLD_MS = 600;

export function TypingTerminal({ scripts }: { scripts: TerminalLine[][] }) {
  const reduce = useReducedMotion();
  const [si, setSi] = useState(0);
  const [typed, setTyped] = useState("");
  const [visibleOutput, setVisibleOutput] = useState(0);
  const [phase, setPhase] = useState<Phase>("typing");
  const [paused, setPaused] = useState(false);

  const names = useMemo(
    () =>
      scripts.map((s) => {
        const input = s.find((l) => l.kind === "input");
        const text = input && input.kind === "input" ? input.text : "";
        return text.replace(/^promptopt\s+/, "").split(/\s+/)[0] || "run";
      }),
    [scripts],
  );

  const script = scripts[si];
  const commandLine = script.find((l) => l.kind === "input");
  const command =
    commandLine && commandLine.kind === "input" ? commandLine.text : "";
  const outputLines = script.filter((l) => l.kind !== "input");

  // 1 · Type the command for the current run.
  useEffect(() => {
    setTyped("");
    setVisibleOutput(0);
    setPhase("typing");

    if (reduce) {
      setTyped(command);
      setVisibleOutput(outputLines.length);
      setPhase("done");
      return;
    }

    const timers: ReturnType<typeof setTimeout>[] = [];
    const at = (ms: number, fn: () => void) => timers.push(setTimeout(fn, ms));

    let i = 0;
    const typeNext = () => {
      i += 1;
      setTyped(command.slice(0, i));
      if (i < command.length) at(24 + Math.random() * 34, typeNext);
      else at(340, () => setPhase("printing"));
    };
    at(260, typeNext);

    return () => timers.forEach(clearTimeout);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [si, reduce]);

  // 2 · Print the output line by line.
  useEffect(() => {
    if (phase !== "printing" || reduce) return;

    const timers: ReturnType<typeof setTimeout>[] = [];
    const at = (ms: number, fn: () => void) => timers.push(setTimeout(fn, ms));

    let n = 0;
    const printNext = () => {
      n += 1;
      setVisibleOutput(n);
      if (n < outputLines.length) {
        at(outputLines[n - 1]?.kind === "gap" ? 45 : 95, printNext);
      } else {
        setPhase("done");
      }
    };
    at(140, printNext);

    return () => timers.forEach(clearTimeout);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [phase, reduce]);

  // 3 · Hold, then advance, unless paused (pointer is over the terminal).
  useEffect(() => {
    if (phase !== "done" || paused || reduce) return;
    const t = setTimeout(
      () => setSi((v) => (v + 1) % scripts.length),
      HOLD_MS,
    );
    return () => clearTimeout(t);
  }, [phase, paused, reduce, scripts.length]);

  return (
    <div
      onMouseEnter={() => setPaused(true)}
      onMouseLeave={() => setPaused(false)}
    >
      <div className="mb-2 flex flex-wrap gap-1 font-mono text-2xs">
        {names.map((name, i) => (
          <button
            key={name}
            type="button"
            onClick={() => setSi(i)}
            aria-current={i === si ? "true" : undefined}
            className={`rounded px-1.5 py-1 transition-colors ${
              i === si
                ? "bg-accent/10 text-accent"
                : "text-faint hover:text-muted"
            }`}
          >
            {name}
          </button>
        ))}
      </div>

      <TerminalChrome title="zsh · promptopt">
        <div
          aria-live="polite"
          className="h-[16.5rem] overflow-hidden sm:h-[17rem]"
        >
          <div className="whitespace-pre-wrap break-words text-fg">
            <span className="select-none text-accent">$ </span>
            {typed}
            {phase === "typing" && (
              <span className="ml-0.5 inline-block h-4 w-2 translate-y-0.5 bg-fg/70 align-baseline animate-blink" />
            )}
          </div>
          <div className="mt-1">
            {outputLines.map((line, i) => (
              <div key={i} className={i < visibleOutput ? "" : "invisible"}>
                <TerminalLineView line={line} />
              </div>
            ))}
          </div>
        </div>
      </TerminalChrome>

      <div className="mt-3 flex gap-1" aria-hidden>
        {scripts.map((_, i) => (
          <span
            key={i}
            className="h-0.5 flex-1 overflow-hidden rounded-full bg-borderStrong"
          >
            <span
              className={`block h-full bg-accent transition-[width] duration-300 ${
                i < si ? "w-full" : i === si ? "w-1/3" : "w-0"
              }`}
            />
          </span>
        ))}
      </div>
    </div>
  );
}
