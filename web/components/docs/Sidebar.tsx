"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import type { DocMeta } from "@/lib/docs";

export type SidebarSection = { section: string; items: DocMeta[] };

export function Sidebar({ nav }: { nav: SidebarSection[] }) {
  const pathname = usePathname();

  return (
    <nav aria-label="Documentation" className="text-sm">
      {nav.map((group) => (
        <div key={group.section} className="mb-6">
          <div className="mb-2 px-2 font-mono text-2xs uppercase tracking-widest text-faint">
            {group.section}
          </div>
          <ul className="space-y-0.5">
            {group.items.map((item) => {
              const href = item.slug === "" ? "/docs" : `/docs/${item.slug}`;
              const active = pathname === href;
              return (
                <li key={href}>
                  <Link
                    href={href}
                    aria-current={active ? "page" : undefined}
                    className={`block rounded-md px-2 py-1.5 transition-colors ${
                      active
                        ? "bg-raised text-fg"
                        : "text-muted hover:bg-raised/50 hover:text-fg"
                    }`}
                  >
                    {item.title}
                  </Link>
                </li>
              );
            })}
          </ul>
        </div>
      ))}
    </nav>
  );
}
