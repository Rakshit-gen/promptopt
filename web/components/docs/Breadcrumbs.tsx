import Link from "next/link";
import type { DocMeta } from "@/lib/docs";

function titleCase(segment: string): string {
  return segment
    .split("-")
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(" ");
}

export function Breadcrumbs({ meta }: { meta: DocMeta }) {
  const crumbs: Array<{ label: string; href: string }> = [
    { label: "Docs", href: "/docs" },
  ];

  if (meta.slug !== "") {
    const parts = meta.slug.split("/");
    let acc = "";
    parts.forEach((part, i) => {
      acc = acc ? `${acc}/${part}` : part;
      const isLast = i === parts.length - 1;
      crumbs.push({
        label: isLast ? meta.title : titleCase(part),
        href: `/docs/${acc}`,
      });
    });
  }

  return (
    <nav aria-label="Breadcrumb" className="mb-6">
      <ol className="flex flex-wrap items-center gap-1.5 font-mono text-2xs text-faint">
        {crumbs.map((c, i) => {
          const isLast = i === crumbs.length - 1;
          return (
            <li key={c.href} className="flex items-center gap-1.5">
              {isLast ? (
                <span className="text-muted">{c.label}</span>
              ) : (
                <Link href={c.href} className="hover:text-fg">
                  {c.label}
                </Link>
              )}
              {!isLast && <span aria-hidden>/</span>}
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
