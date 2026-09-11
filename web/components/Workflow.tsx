"use client";

import { useEffect, useState } from "react";
import { AnimatePresence, motion, useReducedMotion } from "framer-motion";
import { workflow } from "@/lib/demo";

// A pipeline stepper. The steps sit in a horizontal track connected by arrows;
// hovering, focusing, or clicking one shows its detail below. Left alone, it
// advances on its own.

const ADVANCE_MS = 3200;

export function Workflow() {
  const reduce = useReducedMotion();
  const [active, setActive] = useState(0);
  const [held, setHeld] = useState(false);

  useEffect(() => {
    if (reduce || held) return;
    const t = setInterval(
      () => setActive((v) => (v + 1) % workflow.length),
      ADVANCE_MS,
    );
    return () => clearInterval(t);
  }, [reduce, held]);

  const step = workflow[active];

  return (
    <div onMouseLeave={() => setHeld(false)}>
      <ol className="flex flex-nowrap items-center overflow-x-auto pb-1">
        {workflow.map((s, i) => {
          const on = i === active;
          const done = i < active;
          return (
            <li key={s.label} className="flex shrink-0 items-center">
              <button
                type="button"
                onMouseEnter={() => {
                  setHeld(true);
                  setActive(i);
                }}
                onFocus={() => setActive(i)}
                onClick={() => setActive(i)}
                aria-current={on ? "step" : undefined}
                className={`flex items-center gap-2 rounded-full border px-3 py-1.5 font-mono text-xs transition-colors ${
                  on
                    ? "border-transparent bg-accent font-medium text-bg"
                    : done
                      ? "border-border bg-raised text-muted hover:text-fg"
                      : "border-border bg-panel text-faint hover:border-borderStrong hover:text-muted"
                }`}
              >
                <span
                  className={`flex h-4 w-4 items-center justify-center rounded-full text-[0.6rem] ${
                    on
                      ? "bg-bg/25 text-bg"
                      : done
                        ? "bg-accentDim text-bg"
                        : "bg-raised text-faint"
                  }`}
                >
                  {i + 1}
                </span>
                {s.label}
              </button>
              {i < workflow.length - 1 && (
                <svg
                  width="22"
                  height="12"
                  viewBox="0 0 22 12"
                  fill="none"
                  aria-hidden
                  className={`mx-0.5 shrink-0 transition-colors ${
                    i < active ? "text-accentDim" : "text-borderStrong"
                  }`}
                >
                  <path
                    d="M1 6H20M20 6L15 1M20 6L15 11"
                    stroke="currentColor"
                    strokeWidth="1.3"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              )}
            </li>
          );
        })}
      </ol>

      <div className="relative mt-5 overflow-hidden rounded-lg border border-border border-l-2 border-l-accent bg-panel p-5">
        <AnimatePresence mode="wait">
          <motion.div
            key={active}
            initial={reduce ? false : { opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={reduce ? undefined : { opacity: 0, y: -8 }}
            transition={{ duration: 0.2 }}
          >
            <div className="font-mono text-sm text-fg">
              <span className="text-faint">
                {String(active + 1).padStart(2, "0")}
              </span>{" "}
              {step.title}
            </div>
            <p className="mt-2 max-w-xl text-sm text-muted">{step.body}</p>
          </motion.div>
        </AnimatePresence>
        {!reduce && !held && (
          <motion.span
            key={`bar-${active}`}
            className="absolute inset-x-0 bottom-0 h-0.5 origin-left bg-accent"
            initial={{ scaleX: 0 }}
            animate={{ scaleX: 1 }}
            transition={{ duration: ADVANCE_MS / 1000, ease: "linear" }}
          />
        )}
      </div>
    </div>
  );
}
