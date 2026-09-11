"use client";

import { useEffect, useState } from "react";
import { usePathname } from "next/navigation";
import { Sidebar, type SidebarSection } from "./Sidebar";

export function MobileDocsNav({ nav }: { nav: SidebarSection[] }) {
  const [open, setOpen] = useState(false);
  const pathname = usePathname();

  useEffect(() => {
    setOpen(false);
  }, [pathname]);

  return (
    <div className="lg:hidden">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-expanded={open}
        className="flex w-full items-center justify-between rounded-md border border-border bg-panel px-3 py-2 text-sm text-muted hover:text-fg"
      >
        <span>Documentation menu</span>
        <span aria-hidden className="font-mono text-xs">
          {open ? "close" : "open"}
        </span>
      </button>
      {open && (
        <div className="mt-2 rounded-md border border-border bg-panel p-3">
          <Sidebar nav={nav} />
        </div>
      )}
    </div>
  );
}
