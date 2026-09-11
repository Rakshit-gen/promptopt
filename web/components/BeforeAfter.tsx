"use client";

import { useState } from "react";
import { beforeAfter } from "@/lib/demo";

const changeStyle: Record<string, { dot: string; label: string }> = {
  cut: { dot: "bg-err", label: "removed" },
  fix: { dot: "bg-warn", label: "fixed" },
  add: { dot: "bg-ok", label: "added" },
};

export function BeforeAfter() {
  const [view, setView] = useState<"before" | "after">("after");
  const { before, after, tokensBefore, tokensAfter, changes } = beforeAfter;
  const body = view === "before" ? before : after;
  const tokens = view === "before" ? tokensBefore : tokensAfter;
  const reduction = Math.round(((tokensBefore - tokensAfter) / tokensBefore) * 100);

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_320px]">
      <div className="rounded-xl border border-border bg-panel">
        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-4 py-3">
          <div
            className="inline-flex rounded-md border border-border p-0.5"
            role="tablist"
            aria-label="Before or after"
          >
            {(["before", "after"] as const).map((v) => (
              <button
                key={v}
                role="tab"
                aria-selected={view === v}
                onClick={() => setView(v)}
                className={`rounded px-3 py-1 text-sm capitalize transition-colors ${
                  view === v ? "bg-raised text-fg" : "text-muted hover:text-fg"
                }`}
              >
                {v}
              </button>
            ))}
          </div>
          <div className="flex items-center gap-3 font-mono text-2xs">
            <span className="text-faint">tokens (est.)</span>
            <span className="text-fg">{tokens}</span>
            {view === "after" && (
              <span className="rounded bg-ok/10 px-1.5 py-0.5 text-ok">
                −{reduction}%
              </span>
            )}
          </div>
        </div>
        <pre className="h-[22rem] overflow-auto p-4 font-mono text-[0.8rem] leading-6 text-fg/90">
          {body}
        </pre>
      </div>

      <div>
        <h3 className="text-sm font-medium text-fg">What optimize changed</h3>
        <ul className="mt-3 space-y-3">
          {changes.map((c) => {
            const s = changeStyle[c.kind];
            return (
              <li key={c.text} className="flex gap-2.5 text-sm text-muted">
                <span className={`mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full ${s.dot}`} aria-hidden />
                <span>
                  <span className="text-2xs uppercase tracking-wide text-faint">
                    {s.label}
                  </span>
                  <br />
                  {c.text}
                </span>
              </li>
            );
          })}
        </ul>
        <p className="mt-4 text-2xs text-faint">
          Intent is unchanged — the optimized prompt still asks for a customer
          reply. It just stops repeating itself and says what &quot;done&quot;
          means.
        </p>
      </div>
    </div>
  );
}
