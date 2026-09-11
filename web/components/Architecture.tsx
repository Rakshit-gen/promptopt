import { architecture } from "@/lib/demo";

export function Architecture() {
  return (
    <div className="mx-auto max-w-xl">
      <ol className="space-y-0">
        {architecture.map((row, i) => (
          <li key={row.layer} className="flex flex-col items-center">
            <div className="w-full rounded-lg border border-border bg-panel px-4 py-3">
              <div className="flex items-baseline justify-between gap-4">
                <span className="font-mono text-sm text-fg">{row.layer}</span>
                <span className="text-right text-2xs text-faint">{row.note}</span>
              </div>
            </div>
            {i < architecture.length - 1 && (
              <svg
                width="16"
                height="20"
                viewBox="0 0 16 20"
                fill="none"
                aria-hidden
                className="text-borderStrong"
              >
                <path
                  d="M8 1V15M8 15L3 10M8 15L13 10"
                  stroke="currentColor"
                  strokeWidth="1.4"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            )}
          </li>
        ))}
      </ol>
      <p className="mt-5 text-center text-2xs text-faint">
        Dependencies point one way. The engine talks to a one-method interface,
        so every operation is tested with a fake and no network, and another
        provider is one implementation away.
      </p>
    </div>
  );
}
