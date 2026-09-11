"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useReducedMotion } from "framer-motion";
import { TerminalChrome, TerminalLineView } from "./Terminal";
import type { TerminalLine } from "@/lib/demo";

// Cycles through a list of scripted terminal runs, forever, on its own. Each
// run types its command character by character, prints the output line by
// line, holds, then the next run begins. The transcript area is a fixed
// height, so the box never resizes as it moves between commands. Under
// prefers-reduced-motion the first run is shown in full and the cycle is off.

type Phase = "typing" | "printing" | "hold";

export function TypingTerminal({ scripts }: { scripts: TerminalLine[][] }) {
  const reduce = useReducedMotion();
  const [si, setSi] = useState(0);
  const [typed, setTyped] = useState("");
  const [visibleOutput, setVisibleOutput] = useState(0);
  const [phase, setPhase] = useState<Phase>("typing");

  const timers = useRef<ReturnType<typeof setTimeout>[]>([]);
  const clearTimers = useCallback(() => {
    timers.current.forEach(clearTimeout);
    timers.current = [];
  }, []);
  const after = useCallback((ms: number, fn: () => void) => {
    timers.current.push(setTimeout(fn, ms));
  }, []);

  const script = scripts[si];
  const commandLine = script.find((l) => l.kind === "input");
  const command =
    commandLine && commandLine.kind === "input" ? commandLine.text : "";
  const outputLines = script.filter((l) => l.kind !== "input");

  // Type the command for the current run.
  useEffect(() => {
    clearTimers();
    setTyped("");
    setVisibleOutput(0);
    setPhase("typing");

    if (reduce) {
      setTyped(command);
      setVisibleOutput(outputLines.length);
      setPhase("hold");
      return;
    }

    let i = 0;
    const typeNext = () => {
      i += 1;
      setTyped(command.slice(0, i));
      if (i < command.length) {
        after(24 + Math.random() * 34, typeNext);
      } else {
        after(360, () => setPhase("printing"));
      }
    };
    after(260, typeNext);

    return clearTimers;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [si, reduce]);

  // Print the output, then advance to the next run.
  useEffect(() => {
    if (phase !== "printing" || reduce) return;
    let n = 0;
    const printNext = () => {
      n += 1;
      setVisibleOutput(n);
      if (n < outputLines.length) {
        after(outputLines[n - 1]?.kind === "gap" ? 45 : 95, printNext);
      } else {
        setPhase("hold");
        after(2600, () => setSi((v) => (v + 1) % scripts.length));
      }
    };
    after(140, printNext);
    return clearTimers;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [phase, reduce]);

  return (
    <div>
      <TerminalChrome title="zsh — promptopt">
        <div
          aria-live="polite"
          className="h-[19rem] overflow-hidden sm:h-[20rem]"
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

      <div className="mt-3 flex items-center gap-3">
        <div className="flex gap-1.5" aria-hidden>
          {scripts.map((_, i) => (
            <span
              key={i}
              className={`h-1 w-1 rounded-full transition-colors ${
                i === si ? "bg-accent" : "bg-borderStrong"
              }`}
            />
          ))}
        </div>
        <span className="ml-auto font-mono text-2xs text-faint">
          scripted — no API call
        </span>
      </div>
    </div>
  );
}
