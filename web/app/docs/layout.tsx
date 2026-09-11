import type { ReactNode } from "react";
import { Nav } from "@/components/Nav";
import { Footer } from "@/components/Footer";
import { Sidebar } from "@/components/docs/Sidebar";
import { MobileDocsNav } from "@/components/docs/MobileDocsNav";
import { Search } from "@/components/docs/Search";
import { docsNav, searchIndex } from "@/lib/docs";

export default function DocsLayout({ children }: { children: ReactNode }) {
  const nav = docsNav();
  const search = searchIndex();

  return (
    <>
      <Nav />
      <div className="container-content grid gap-10 py-8 lg:grid-cols-[15rem_minmax(0,1fr)]">
        <aside className="hidden lg:block">
          <div className="sticky top-16 max-h-[calc(100vh-5rem)] space-y-4 overflow-y-auto pb-8 pr-2">
            <Search docs={search} />
            <Sidebar nav={nav} />
          </div>
        </aside>

        <div className="min-w-0">
          <div className="mb-6 space-y-3 lg:hidden">
            <Search docs={search} />
            <MobileDocsNav nav={nav} />
          </div>
          {children}
        </div>
      </div>
      <Footer />
    </>
  );
}
