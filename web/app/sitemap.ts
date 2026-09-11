import type { MetadataRoute } from "next";
import { site } from "@/lib/site";
import { allDocs } from "@/lib/docs";

export default function sitemap(): MetadataRoute.Sitemap {
  const now = new Date();
  const docRoutes = allDocs().map((d) => ({
    url: d.slug === "" ? `${site.url}/docs` : `${site.url}/docs/${d.slug}`,
    lastModified: now,
    changeFrequency: "monthly" as const,
    priority: d.slug === "" ? 0.7 : 0.5,
  }));

  return [
    { url: site.url, lastModified: now, changeFrequency: "monthly", priority: 1 },
    ...docRoutes,
  ];
}
