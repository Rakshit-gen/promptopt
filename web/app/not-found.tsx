import Link from "next/link";
import { Nav } from "@/components/Nav";
import { Footer } from "@/components/Footer";

export default function NotFound() {
  return (
    <>
      <Nav />
      <main id="main" className="container-content flex min-h-[60vh] flex-col items-center justify-center py-24 text-center">
        <div className="font-mono text-2xs uppercase tracking-widest text-accent">
          404
        </div>
        <h1 className="mt-3 text-2xl font-semibold tracking-tight">
          No page here
        </h1>
        <p className="mt-2 max-w-sm text-sm text-muted">
          That route does not exist. The command you were looking for is
          probably in the docs.
        </p>
        <div className="mt-6 flex gap-3">
          <Link
            href="/"
            className="rounded-md border border-border bg-panel px-4 py-2 text-sm text-muted hover:text-fg"
          >
            Home
          </Link>
          <Link
            href="/docs"
            className="rounded-md border border-border bg-panel px-4 py-2 text-sm text-muted hover:text-fg"
          >
            Docs
          </Link>
        </div>
      </main>
      <Footer />
    </>
  );
}
