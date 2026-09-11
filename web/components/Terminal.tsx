import type { TerminalLine } from "@/lib/demo";

const toneClass: Record<string, string> = {
  dim: "text-faint",
  ok: "text-ok",
  warn: "text-warn",
  err: "text-err",
  accent: "text-accent",
};

export function TerminalChrome({
  title = "zsh",
  children,
  className = "",
  contentClassName = "p-4 font-mono text-[0.8rem] leading-6 sm:text-[0.83rem]",
}: {
  title?: string;
  children: React.ReactNode;
  className?: string;
  contentClassName?: string;
}) {
  return (
    <div
      className={`overflow-hidden rounded-xl border border-border bg-panel shadow-[0_1px_0_0_rgba(255,255,255,0.03)_inset,0_20px_60px_-30px_rgba(0,0,0,0.7)] ${className}`}
    >
      <div className="flex items-center gap-2 border-b border-border bg-raised/60 px-3.5 py-2.5">
        <span className="flex gap-1.5" aria-hidden>
          <span className="h-2.5 w-2.5 rounded-full bg-borderStrong" />
          <span className="h-2.5 w-2.5 rounded-full bg-borderStrong" />
          <span className="h-2.5 w-2.5 rounded-full bg-borderStrong" />
        </span>
        <span className="ml-1 font-mono text-2xs text-faint">{title}</span>
      </div>
      <div className={contentClassName}>{children}</div>
    </div>
  );
}

export function TerminalLineView({ line }: { line: TerminalLine }) {
  if (line.kind === "gap") return <div className="h-3" aria-hidden />;
  if (line.kind === "input") {
    return (
      <div className="whitespace-pre-wrap break-words text-fg">
        <span className="select-none text-accent">$ </span>
        {line.text}
      </div>
    );
  }
  return (
    <div
      className={`whitespace-pre-wrap break-words ${
        line.tone ? toneClass[line.tone] : "text-fg/90"
      }`}
    >
      {line.text}
    </div>
  );
}

export function TerminalTranscript({ lines }: { lines: TerminalLine[] }) {
  return (
    <div>
      {lines.map((line, i) => (
        <TerminalLineView key={i} line={line} />
      ))}
    </div>
  );
}
