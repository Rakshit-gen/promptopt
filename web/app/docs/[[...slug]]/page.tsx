import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { docSlugs, getDoc } from "@/lib/docs";
import { site } from "@/lib/site";
import { Markdown } from "@/components/docs/Markdown";
import { Toc } from "@/components/docs/Toc";
import { Breadcrumbs } from "@/components/docs/Breadcrumbs";
import { Pager } from "@/components/docs/Pager";

type Props = { params: Promise<{ slug?: string[] }> };

export function generateStaticParams() {
  return docSlugs().map((slug) => ({ slug }));
}

export const dynamicParams = false;

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const doc = getDoc(slug);
  if (!doc) return {};

  const url =
    doc.meta.slug === "" ? `${site.url}/docs` : `${site.url}/docs/${doc.meta.slug}`;
  const title = doc.meta.title;
  const description = doc.meta.description;

  return {
    title,
    description,
    alternates: { canonical: url },
    openGraph: { title, description, url, type: "article" },
    twitter: { card: "summary_large_image", title, description },
  };
}

export default async function DocPage({ params }: Props) {
  const { slug } = await params;
  const doc = getDoc(slug);
  if (!doc) notFound();

  return (
    <div className="xl:grid xl:grid-cols-[minmax(0,1fr)_13rem] xl:gap-10">
      <article className="min-w-0">
        <Breadcrumbs meta={doc.meta} />
        <Markdown>{doc.content}</Markdown>
        <Pager prev={doc.prev} next={doc.next} />
      </article>

      <aside className="hidden xl:block">
        <div className="sticky top-16 max-h-[calc(100vh-5rem)] overflow-y-auto pb-8">
          <Toc entries={doc.toc} />
        </div>
      </aside>
    </div>
  );
}
