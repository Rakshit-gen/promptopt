// Documentation loader. The canonical docs are the Markdown files in the
// repository's top-level docs/ directory, the same files someone browsing the
// repo on GitHub reads. This module maps them to site routes, resolves the
// relative .md links between them, and extracts a table of contents.

import fs from "node:fs";
import path from "node:path";
import matter from "gray-matter";
import GithubSlugger from "github-slugger";

const REPO_ROOT = path.join(process.cwd(), "..");
const DOCS_DIR = path.join(REPO_ROOT, "docs");
const GITHUB_BLOB = "https://github.com/rakshit-gen/promptopt/blob/master";

export type DocMeta = {
  slug: string; // route segment(s) after /docs, "" for the index
  file: string; // path relative to docs/
  title: string;
  description: string;
  section: string;
};

// Explicit ordering. Titles are read from each file's H1 unless overridden.
const MANIFEST: Array<Omit<DocMeta, "title" | "description">> = [
  { slug: "", file: "README.md", section: "Overview" },
  { slug: "getting-started", file: "getting-started.md", section: "Overview" },
  { slug: "installation", file: "installation.md", section: "Overview" },
  { slug: "configuration", file: "configuration.md", section: "Overview" },

  { slug: "commands", file: "commands/README.md", section: "Commands" },
  { slug: "commands/optimize", file: "commands/optimize.md", section: "Commands" },
  { slug: "commands/compress", file: "commands/compress.md", section: "Commands" },
  { slug: "commands/expand", file: "commands/expand.md", section: "Commands" },
  { slug: "commands/analyze", file: "commands/analyze.md", section: "Commands" },
  { slug: "commands/transform", file: "commands/transform.md", section: "Commands" },
  { slug: "commands/eval", file: "commands/eval.md", section: "Commands" },

  { slug: "concepts", file: "concepts.md", section: "Concepts" },
  { slug: "prompt-analysis", file: "prompt-analysis.md", section: "Concepts" },
  { slug: "prompt-compression", file: "prompt-compression.md", section: "Concepts" },
  { slug: "evaluation", file: "evaluation.md", section: "Concepts" },
  { slug: "json-output", file: "json-output.md", section: "Concepts" },

  { slug: "ci", file: "ci.md", section: "Operating" },
  { slug: "troubleshooting", file: "troubleshooting.md", section: "Operating" },
];

function firstHeading(md: string): string {
  const m = md.match(/^#\s+(.+)$/m);
  return m ? m[1].trim() : "Untitled";
}

function firstParagraph(md: string): string {
  const body = md.replace(/^#\s+.+$/m, "").trim();
  const para = body.split("\n\n")[0]?.replace(/\n/g, " ").trim() ?? "";
  return para.length > 180 ? para.slice(0, 177) + "…" : para;
}

let cache: DocMeta[] | null = null;

export function allDocs(): DocMeta[] {
  if (cache) return cache;
  cache = MANIFEST.map((entry) => {
    const raw = fs.readFileSync(path.join(DOCS_DIR, entry.file), "utf8");
    const { content, data } = matter(raw);
    return {
      ...entry,
      title: (data.title as string) || firstHeading(content),
      description: (data.description as string) || firstParagraph(content),
    };
  });
  return cache;
}

export function docSlugs(): string[][] {
  return allDocs().map((d) => (d.slug === "" ? [] : d.slug.split("/")));
}

export function getDoc(slugParts: string[] | undefined): {
  meta: DocMeta;
  content: string;
  toc: TocEntry[];
  prev: DocMeta | null;
  next: DocMeta | null;
} | null {
  const slug = (slugParts ?? []).join("/");
  const docs = allDocs();
  const idx = docs.findIndex((d) => d.slug === slug);
  if (idx === -1) return null;

  const meta = docs[idx];
  const raw = fs.readFileSync(path.join(DOCS_DIR, meta.file), "utf8");
  const { content } = matter(raw);
  const rewritten = rewriteLinks(content, meta.file);

  return {
    meta,
    content: rewritten,
    toc: buildToc(rewritten),
    prev: idx > 0 ? docs[idx - 1] : null,
    next: idx < docs.length - 1 ? docs[idx + 1] : null,
  };
}

export type TocEntry = { depth: 2 | 3; text: string; id: string };

function buildToc(md: string): TocEntry[] {
  const slugger = new GithubSlugger();
  const out: TocEntry[] = [];
  let inFence = false;
  for (const line of md.split("\n")) {
    if (line.trimStart().startsWith("```")) {
      inFence = !inFence;
      continue;
    }
    if (inFence) continue;
    const m = line.match(/^(#{2,3})\s+(.+?)\s*$/);
    if (!m) continue;
    const depth = m[1].length as 2 | 3;
    const text = m[2].replace(/`/g, "").trim();
    out.push({ depth, text, id: slugger.slug(text) });
  }
  return out;
}

// Map a docs/ file path to its site route.
function fileToRoute(file: string): string | null {
  const meta = allDocs().find((d) => d.file === file);
  if (!meta) return null;
  return meta.slug === "" ? "/docs" : `/docs/${meta.slug}`;
}

// Rewrite Markdown link targets:
//  - relative .md links to other docs -> /docs/... routes
//  - links to files outside docs/ (LICENSE, SECURITY.md) -> GitHub blob URLs
//  - external and pure-anchor links are left alone
function rewriteLinks(md: string, currentFile: string): string {
  const currentDir = path.dirname(currentFile); // relative to docs/
  return md.replace(/\]\(([^)]+)\)/g, (whole, target: string) => {
    const t = target.trim();
    if (/^(https?:|mailto:|#)/.test(t)) return whole;

    const [rawPath, anchor] = t.split("#");
    if (!rawPath) return whole;

    // Resolve relative to the current file's directory, within the repo.
    const fromDocs = path.posix.normalize(
      path.posix.join(currentDir === "." ? "" : currentDir, rawPath),
    );

    if (fromDocs.startsWith("../")) {
      // Points outside docs/: link to the file on GitHub.
      const repoPath = path.posix.normalize(
        path.posix.join("docs", currentDir, rawPath),
      );
      return `](${GITHUB_BLOB}/${repoPath}${anchor ? "#" + anchor : ""})`;
    }

    const route = fileToRoute(fromDocs);
    if (route) {
      return `](${route}${anchor ? "#" + anchor : ""})`;
    }
    return whole;
  });
}

export function docsNav(): Array<{ section: string; items: DocMeta[] }> {
  const sections: string[] = [];
  const bySection = new Map<string, DocMeta[]>();
  for (const d of allDocs()) {
    if (!bySection.has(d.section)) {
      bySection.set(d.section, []);
      sections.push(d.section);
    }
    bySection.get(d.section)!.push(d);
  }
  return sections.map((section) => ({ section, items: bySection.get(section)! }));
}

export type SearchDoc = { title: string; slug: string; section: string; headings: string[] };

export function searchIndex(): SearchDoc[] {
  return allDocs().map((d) => {
    const raw = fs.readFileSync(path.join(DOCS_DIR, d.file), "utf8");
    const { content } = matter(raw);
    return {
      title: d.title,
      slug: d.slug,
      section: d.section,
      headings: buildToc(content).map((t) => t.text),
    };
  });
}
