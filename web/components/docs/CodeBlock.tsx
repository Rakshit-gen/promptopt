"use client";

import { useRef, useState } from "react";

// A <pre> with a copy button. The rendered, syntax-highlighted markup is
// passed as children; the raw text (for the clipboard) is passed separately
// because it is extracted from the markdown AST at render time.
export function CodeBlock({
  raw,
  language,
  children,
}: {
  raw: string;
  language?: string;
  children: React.ReactNode;
}) {
  const [copied, setCopied] = useState(false);
  const ref = useRef<HTMLPreElement>(null);

  async function copy() {
    try {
      await navigator.clipboard.writeText(raw);
      setCopied(true);
      setTimeout(() => setCopied(false), 1600);
    } catch {
      // ignore
    }
  }

  return (
    <div className="group relative my-4">
      {language && (
        <span className="absolute left-3 top-2.5 select-none font-mono text-2xs text-faint">
          {language}
        </span>
      )}
      <button
        type="button"
        onClick={copy}
        aria-label={copied ? "Copied" : "Copy code"}
        className="absolute right-2 top-2 z-10 rounded border border-border bg-raised px-1.5 py-1 font-mono text-2xs text-muted opacity-0 transition-opacity hover:text-fg focus-visible:opacity-100 group-hover:opacity-100"
      >
        {copied ? "copied" : "copy"}
      </button>
      <pre ref={ref} className={language ? "pt-8" : undefined}>
        {children}
      </pre>
    </div>
  );
}
