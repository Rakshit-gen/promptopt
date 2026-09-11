import Link from "next/link";
import type { DocMeta } from "@/lib/docs";

function href(meta: DocMeta): string {
  return meta.slug === "" ? "/docs" : `/docs/${meta.slug}`;
}

export function Pager({
  prev,
  next,
}: {
  prev: DocMeta | null;
  next: DocMeta | null;
}) {
  return (
    <nav
      aria-label="Pagination"
      className="mt-14 grid gap-3 border-t border-border pt-6 sm:grid-cols-2"
    >
      {prev ? (
        <Link
          href={href(prev)}
          className="rounded-lg border border-border bg-panel/50 p-3 transition-colors hover:border-borderStrong"
        >
          <div className="font-mono text-2xs uppercase tracking-widest text-faint">
            Previous
          </div>
          <div className="mt-1 text-sm text-fg">{prev.title}</div>
        </Link>
      ) : (
        <span />
      )}
      {next ? (
        <Link
          href={href(next)}
          className="rounded-lg border border-border bg-panel/50 p-3 text-right transition-colors hover:border-borderStrong"
        >
          <div className="font-mono text-2xs uppercase tracking-widest text-faint">
            Next
          </div>
          <div className="mt-1 text-sm text-fg">{next.title}</div>
        </Link>
      ) : (
        <span />
      )}
    </nav>
  );
}
