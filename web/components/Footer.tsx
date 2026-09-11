import Link from "next/link";
import { Logo } from "./Logo";
import { site } from "@/lib/site";

export function Footer() {
  return (
    <footer className="border-t border-border">
      <div className="container-content grid gap-8 py-12 sm:grid-cols-2 lg:grid-cols-4">
        <div className="space-y-3">
          <Logo />
          <p className="max-w-xs text-sm text-muted">{site.tagline}</p>
        </div>

        <FooterCol
          title="Product"
          links={[
            { label: "Overview", href: "/#product" },
            { label: "Commands", href: "/#commands" },
            { label: "Install", href: "/#install" },
            { label: "Architecture", href: "/#architecture" },
          ]}
        />
        <FooterCol
          title="Docs"
          links={[
            { label: "Getting started", href: "/docs/getting-started" },
            { label: "Commands", href: "/docs/commands" },
            { label: "Configuration", href: "/docs/configuration" },
            { label: "JSON output", href: "/docs/json-output" },
            { label: "CI", href: "/docs/ci" },
          ]}
        />
        <FooterCol
          title="Project"
          links={[
            { label: "GitHub", href: site.repo, external: true },
            { label: "Releases", href: `${site.repo}/releases`, external: true },
            { label: "Issues", href: `${site.repo}/issues`, external: true },
            { label: "License (MIT)", href: `${site.repo}/blob/master/LICENSE`, external: true },
            { label: "Security", href: `${site.repo}/blob/master/SECURITY.md`, external: true },
          ]}
        />
      </div>
      <div className="container-content flex flex-col gap-2 border-t border-border py-6 text-2xs text-faint sm:flex-row sm:items-center sm:justify-between">
        <span>
          MIT licensed. Inference by{" "}
          <a href="https://groq.com" className="text-muted hover:text-fg" target="_blank" rel="noreferrer">
            Groq
          </a>
          .
        </span>
        <span>No telemetry. Your prompts and key are never logged.</span>
      </div>
    </footer>
  );
}

function FooterCol({
  title,
  links,
}: {
  title: string;
  links: Array<{ label: string; href: string; external?: boolean }>;
}) {
  return (
    <div className="space-y-2.5">
      <h3 className="text-2xs font-medium uppercase tracking-wide text-faint">{title}</h3>
      <ul className="space-y-2">
        {links.map((l) => (
          <li key={l.label}>
            {l.external ? (
              <a
                href={l.href}
                target="_blank"
                rel="noreferrer"
                className="text-sm text-muted transition-colors hover:text-fg"
              >
                {l.label}
              </a>
            ) : (
              <Link href={l.href} className="text-sm text-muted transition-colors hover:text-fg">
                {l.label}
              </Link>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}
