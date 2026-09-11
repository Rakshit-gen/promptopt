// The promptopt wordmark: a terminal prompt caret, then the name. No robot,
// no gradient. The caret is the same ">" a shell prints.

export function Logo({ className = "" }: { className?: string }) {
  return (
    <span
      className={`inline-flex items-center gap-2 font-mono text-[0.95rem] font-medium tracking-tight ${className}`}
    >
      <span
        aria-hidden
        className="inline-flex h-5 w-5 items-center justify-center rounded border border-borderStrong bg-raised text-accent"
      >
        <svg width="10" height="10" viewBox="0 0 10 10" fill="none">
          <path
            d="M1 1.5L4.5 5L1 8.5"
            stroke="currentColor"
            strokeWidth="1.6"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <path
            d="M5.5 8.5H9"
            stroke="currentColor"
            strokeWidth="1.6"
            strokeLinecap="round"
          />
        </svg>
      </span>
      <span>
        prompt<span className="text-muted">opt</span>
      </span>
    </span>
  );
}
