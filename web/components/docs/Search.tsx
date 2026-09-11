"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { useRouter } from "next/navigation";
import type { SearchDoc } from "@/lib/docs";

function route(slug: string): string {
  return slug === "" ? "/docs" : `/docs/${slug}`;
}

export function Search({ docs }: { docs: SearchDoc[] }) {
  const [open, setOpen] = useState(false);
  const [q, setQ] = useState("");
  const [sel, setSel] = useState(0);
  const router = useRouter();
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setOpen((v) => !v);
      } else if (e.key === "Escape") {
        setOpen(false);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  useEffect(() => {
    if (!open) return;
    setQ("");
    setSel(0);
    const t = setTimeout(() => inputRef.current?.focus(), 0);
    return () => clearTimeout(t);
  }, [open]);

  const results = useMemo(() => {
    const term = q.trim().toLowerCase();
    if (!term) return docs.slice(0, 8);
    return docs
      .map((d) => {
        let score = 0;
        if (d.title.toLowerCase().includes(term)) score += 10;
        if (d.section.toLowerCase().includes(term)) score += 2;
        for (const h of d.headings) {
          if (h.toLowerCase().includes(term)) score += 3;
        }
        return { d, score };
      })
      .filter((r) => r.score > 0)
      .sort((a, b) => b.score - a.score)
      .slice(0, 8)
      .map((r) => r.d);
  }, [q, docs]);

  const go = useCallback(
    (d: SearchDoc) => {
      setOpen(false);
      router.push(route(d.slug));
    },
    [router],
  );

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="flex w-full items-center justify-between gap-4 rounded-md border border-border bg-panel px-3 py-2 text-sm text-muted transition-colors hover:text-fg"
      >
        <span>Search docs</span>
        <span className="kbd">⌘K</span>
      </button>

      {open &&
        createPortal(
          <div
            className="fixed inset-0 z-50 flex items-start justify-center bg-bg/80 p-4 pt-[12vh] backdrop-blur-sm"
            onClick={() => setOpen(false)}
            role="presentation"
          >
            <div
              className="w-full max-w-lg overflow-hidden rounded-xl border border-borderStrong bg-panel shadow-2xl"
              onClick={(e) => e.stopPropagation()}
              role="dialog"
              aria-modal="true"
              aria-label="Search documentation"
            >
              <input
                ref={inputRef}
                value={q}
                onChange={(e) => {
                  setQ(e.target.value);
                  setSel(0);
                }}
                onKeyDown={(e) => {
                  if (e.key === "ArrowDown") {
                    e.preventDefault();
                    setSel((s) => Math.min(s + 1, results.length - 1));
                  } else if (e.key === "ArrowUp") {
                    e.preventDefault();
                    setSel((s) => Math.max(s - 1, 0));
                  } else if (e.key === "Enter" && results[sel]) {
                    e.preventDefault();
                    go(results[sel]);
                  }
                }}
                placeholder="Search documentation…"
                className="w-full border-b border-border bg-transparent px-4 py-3 text-sm outline-none placeholder:text-faint"
              />
              <ul className="max-h-80 overflow-y-auto p-2">
                {results.length === 0 && (
                  <li className="px-3 py-6 text-center text-sm text-faint">
                    No matches
                  </li>
                )}
                {results.map((d, i) => (
                  <li key={d.slug || "index"}>
                    <button
                      type="button"
                      onMouseEnter={() => setSel(i)}
                      onClick={() => go(d)}
                      className={`flex w-full flex-col items-start rounded-md px-3 py-2 text-left transition-colors ${
                        i === sel ? "bg-raised" : ""
                      }`}
                    >
                      <span className="text-sm text-fg">{d.title}</span>
                      <span className="font-mono text-2xs text-faint">
                        {d.section}
                      </span>
                    </button>
                  </li>
                ))}
              </ul>
            </div>
          </div>,
          document.body,
        )}
    </>
  );
}
