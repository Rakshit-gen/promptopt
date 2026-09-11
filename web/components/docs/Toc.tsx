"use client";

import { useEffect, useState } from "react";
import type { TocEntry } from "@/lib/docs";

export function Toc({ entries }: { entries: TocEntry[] }) {
  const [active, setActive] = useState("");

  useEffect(() => {
    if (entries.length === 0) return;
    const observer = new IntersectionObserver(
      (items) => {
        const visible = items
          .filter((i) => i.isIntersecting)
          .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top);
        if (visible[0]) setActive(visible[0].target.id);
      },
      { rootMargin: "-80px 0px -70% 0px", threshold: 0 },
    );
    for (const e of entries) {
      const el = document.getElementById(e.id);
      if (el) observer.observe(el);
    }
    return () => observer.disconnect();
  }, [entries]);

  if (entries.length === 0) return null;

  return (
    <nav aria-label="On this page" className="text-sm">
      <div className="mb-2 font-mono text-2xs uppercase tracking-widest text-faint">
        On this page
      </div>
      <ul>
        {entries.map((e) => (
          <li key={e.id}>
            <a
              href={`#${e.id}`}
              className={`-ml-px block border-l-2 py-1 transition-colors ${
                e.depth === 3 ? "pl-6" : "pl-3"
              } ${
                active === e.id
                  ? "border-accent text-fg"
                  : "border-border text-muted hover:border-borderStrong hover:text-fg"
              }`}
            >
              {e.text}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
}
