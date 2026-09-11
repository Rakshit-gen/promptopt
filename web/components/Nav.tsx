"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Logo } from "./Logo";
import { nav, site } from "@/lib/site";

export function Nav() {
  const [open, setOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 8);
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  return (
    <header
      className={`sticky top-0 z-40 border-b transition-colors ${
        scrolled
          ? "border-border bg-bg/85 backdrop-blur"
          : "border-transparent bg-transparent"
      }`}
    >
      <div className="container-content flex h-14 items-center justify-between">
        <Link href="/" aria-label="promptopt home" className="shrink-0">
          <Logo />
        </Link>

        <nav className="hidden items-center gap-1 md:flex" aria-label="Primary">
          {nav.map((item) =>
            "external" in item && item.external ? (
              <a
                key={item.label}
                href={item.href}
                target="_blank"
                rel="noreferrer"
                className="rounded-md px-3 py-1.5 text-sm text-muted transition-colors hover:text-fg"
              >
                {item.label}
              </a>
            ) : (
              <Link
                key={item.label}
                href={item.href}
                className="rounded-md px-3 py-1.5 text-sm text-muted transition-colors hover:text-fg"
              >
                {item.label}
              </Link>
            ),
          )}
        </nav>

        <div className="hidden items-center gap-2 md:flex">
          <Link
            href="/docs"
            className="rounded-md px-3 py-1.5 text-sm text-muted transition-colors hover:text-fg"
          >
            Read docs
          </Link>
          <a
            href="#install"
            className="rounded-md border border-borderStrong bg-raised px-3 py-1.5 text-sm font-medium text-fg transition-colors hover:border-accent/50"
          >
            Install CLI
          </a>
        </div>

        <button
          type="button"
          className="rounded-md border border-border p-2 md:hidden"
          aria-expanded={open}
          aria-label="Toggle menu"
          onClick={() => setOpen((v) => !v)}
        >
          <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden>
            {open ? (
              <path
                d="M4 4L14 14M14 4L4 14"
                stroke="currentColor"
                strokeWidth="1.6"
                strokeLinecap="round"
              />
            ) : (
              <path
                d="M3 5H15M3 9H15M3 13H15"
                stroke="currentColor"
                strokeWidth="1.6"
                strokeLinecap="round"
              />
            )}
          </svg>
        </button>
      </div>

      {open && (
        <div className="border-t border-border bg-bg md:hidden">
          <nav className="container-content flex flex-col py-3" aria-label="Mobile">
            {nav.map((item) => (
              <a
                key={item.label}
                href={item.href}
                target={"external" in item && item.external ? "_blank" : undefined}
                rel={"external" in item && item.external ? "noreferrer" : undefined}
                className="rounded-md px-2 py-2.5 text-sm text-muted hover:text-fg"
                onClick={() => setOpen(false)}
              >
                {item.label}
              </a>
            ))}
            <a
              href="#install"
              onClick={() => setOpen(false)}
              className="mt-2 rounded-md border border-borderStrong bg-raised px-2 py-2.5 text-center text-sm font-medium text-fg"
            >
              Install CLI
            </a>
          </nav>
        </div>
      )}
    </header>
  );
}

export function siteRepo() {
  return site.repo;
}
