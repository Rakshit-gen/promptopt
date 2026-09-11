// The promptopt wordmark: the cat mascot, then the name.

import { Mascot } from "./Mascot";

export function Logo({ className = "" }: { className?: string }) {
  return (
    <span
      className={`inline-flex items-center gap-2 font-mono text-[0.95rem] font-medium tracking-tight ${className}`}
    >
      <Mascot className="h-6 w-6 shrink-0" />
      <span>
        prompt<span className="text-muted">opt</span>
      </span>
    </span>
  );
}
