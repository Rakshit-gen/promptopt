"use client";

import { useState } from "react";
import { CopyButton } from "./CopyButton";
import { site } from "@/lib/site";

const tabs = {
  macOS: [
    "# works on Apple silicon or Intel, detected automatically",
    site.installCommand,
    "",
    "# then set your key (from https://console.groq.com/keys)",
    site.keyExport,
  ].join("\n"),
  Linux: [
    "# works on amd64 or arm64, detected automatically",
    site.installCommand,
    "",
    "# add ~/.local/bin to PATH if the installer asks, then",
    site.keyExport,
  ].join("\n"),
  "go install": [
    "go install github.com/rakshit-gen/promptopt/cmd/promptopt@latest",
    "",
    site.keyExport,
  ].join("\n"),
};

type Tab = keyof typeof tabs;

export function InstallBlock() {
  const [tab, setTab] = useState<Tab>("macOS");

  return (
    <div className="rounded-xl border border-border bg-panel">
      <div className="flex items-center justify-between border-b border-border px-3 py-2">
        <div className="flex gap-1" role="tablist" aria-label="Install method">
          {(Object.keys(tabs) as Tab[]).map((t) => (
            <button
              key={t}
              role="tab"
              aria-selected={tab === t}
              onClick={() => setTab(t)}
              className={`rounded px-2.5 py-1 font-mono text-2xs transition-colors ${
                tab === t ? "bg-raised text-fg" : "text-muted hover:text-fg"
              }`}
            >
              {t}
            </button>
          ))}
        </div>
        <CopyButton value={tabs[tab]} label="Copy" />
      </div>
      <pre className="overflow-x-auto p-4 font-mono text-[0.82rem] leading-6">
        {tabs[tab].split("\n").map((line, i) => (
          <div
            key={i}
            className={line.startsWith("#") ? "text-faint" : "text-fg"}
          >
            {line.startsWith("#") || line === "" ? (
              line || " "
            ) : (
              <>
                <span className="select-none text-accent">$ </span>
                {line}
              </>
            )}
          </div>
        ))}
      </pre>
      <div className="border-t border-border px-4 py-2.5 text-2xs text-faint">
        The installer checks the SHA256 checksum before installing. It needs
        only <code className="text-muted">curl</code> and{" "}
        <code className="text-muted">tar</code>, no Python, Node, Go, Homebrew,
        or Docker.
      </div>
    </div>
  );
}
