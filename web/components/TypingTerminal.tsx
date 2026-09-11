"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useReducedMotion } from "framer-motion";
import { TerminalChrome, TerminalLineView } from "./Terminal";
import type { TerminalLine } from "@/lib/demo";

// Plays a scripted terminal transcript: the command is "typed" character by
// character, then the output prints line by line. A Run button replays it.
// With prefers-reduced-motion, the whole transcript renders immediately.

type Phase = "idle" | "typing" | "printing" | "done";

export function TypingTerminal({ script }: { script: TerminalLine[] }) {
  const reduce = useReducedMotion();
  const [typed, setTyped] = useState("");
  const [visibleOutput, setVisibleOutput] = useState(0);
  const [phase, setPhase] = useState<Phase>("idle");
  const timers = useRef<ReturnType<typeof setTimeout>[]>([]);

  const commandLine = script.find((l) => l.kind === "input");
  const command = commandLine && commandLine.kind === "input" ? commandLine.text : "";
  const outputLines = script.filter((l) => l.kind !== "input");

  const clearTimers = () => {
    timers.current.forEach(clearTimeout);
    timers.current = [];
  };

  const showAll = useCallback(() => {
    clearTimers();
    setTyped(command);
    setVisibleOutput(outputLines.length);
    setPhase("done");
  }, [command, outputLines.length]);

  const run = useCallback(() => {
    clearTimers();
    setTyped("");
    setVisibleOutput(0);
    setPhase("typing");

    let i = 0;
    const typeNext = () => {
      i += 1;
      setTyped(command.slice(0, i));
      if (i < command.length) {
        timers.current.push(setTimeout(typeNext, 26 + Math.random() * 34));
      } else {
        timers.current.push(setTimeout(() => setPhase("printing"), 320));
      }
    };
    timers.current.push(setTimeout(typeNext, 240));
  }, [command]);

  // Print output lines once typing finishes.
  useEffect(() => {
    if (phase !== "printing") return;
    let n = 0;
    const printNext = () => {
      n += 1;
      setVisibleOutput(n);
      if (n < outputLines.length) {
        const delay = outputLines[n - 1]?.kind === "gap" ? 40 : 90;
        timers.current.push(setTimeout(printNext, delay));
      } else {
        setPhase("done");
      }
    };
    timers.current.push(setTimeout(printNext, 120));
    return clearTimers;
  }, [phase, outputLines]);

  // Autoplay on mount (once), or show everything if reduced motion.
  useEffect(() => {
    if (reduce) {
      showAll();
      return;
    }
    run();
    return clearTimers;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const running = phase === "typing" || phase === "printing";

  return (
    <div>
      <TerminalChrome title="zsh — promptopt">
        <div aria-live="polite">
          <div className="whitespace-pre-wrap break-words text-fg">
            <span className="select-none text-accent">$ </span>
            {typed}
            {(phase === "typing" || phase === "idle") && (
              <span className="ml-0.5 inline-block h-4 w-2 translate-y-0.5 bg-fg/70 align-baseline animate-blink" />
            )}
          </div>
          <div className="mt-1">
            {outputLines.slice(0, visibleOutput).map((line, i) => (
              <TerminalLineView key={i} line={line} />
            ))}
          </div>
        </div>
      </TerminalChrome>

      <div className="mt-3 flex items-center gap-3">
        <button
          type="button"
          onClick={run}
          disabled={running}
          className="inline-flex items-center gap-1.5 rounded-md border border-borderStrong bg-raised px-2.5 py-1.5 font-mono text-2xs text-muted transition-colors hover:border-accent/50 hover:text-fg disabled:opacity-50"
        >
          <svg width="11" height="11" viewBox="0 0 11 11" fill="none" aria-hidden>
            <path d="M2.5 1.5L9 5.5L2.5 9.5V1.5Z" fill="currentColor" />
          </svg>
          {phase === "done" ? "Replay" : "Running…"}
        </button>
        <span className="font-mono text-2xs text-faint">
          scripted demo — deterministic data, no API call
        </span>
      </div>
    </div>
  );
}
