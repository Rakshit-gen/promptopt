import { Children, isValidElement, type ReactNode } from "react";
import Link from "next/link";
import ReactMarkdown, { type Components } from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeSlug from "rehype-slug";
import rehypeHighlight from "rehype-highlight";
import { CodeBlock } from "./CodeBlock";

// Reconstruct the plain text of a rendered node tree. After rehype-highlight
// runs, a code block's children are nested <span> elements, so the raw source
// for the clipboard has to be walked back out of them.
function nodeText(node: ReactNode): string {
  if (node == null || typeof node === "boolean") return "";
  if (typeof node === "string" || typeof node === "number") return String(node);
  if (Array.isArray(node)) return node.map(nodeText).join("");
  if (isValidElement(node)) {
    return nodeText((node.props as { children?: ReactNode }).children);
  }
  return "";
}

const components: Components = {
  a({ href, children }) {
    const url = href ?? "";
    if (url.startsWith("/")) {
      return <Link href={url}>{children}</Link>;
    }
    const external = /^https?:/.test(url);
    return (
      <a href={url} {...(external ? { target: "_blank", rel: "noreferrer" } : {})}>
        {children}
      </a>
    );
  },
  pre({ children }) {
    let raw = "";
    let language: string | undefined;
    Children.forEach(children, (child) => {
      if (!isValidElement(child)) return;
      const props = child.props as { className?: string; children?: ReactNode };
      raw = nodeText(props.children).replace(/\n$/, "");
      const m = /language-([\w-]+)/.exec(props.className ?? "");
      if (m) language = m[1];
    });
    return (
      <CodeBlock raw={raw} language={language}>
        {children}
      </CodeBlock>
    );
  },
};

export function Markdown({ children }: { children: string }) {
  return (
    <div className="prose-doc">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[
          rehypeSlug,
          [rehypeHighlight, { detect: true, ignoreMissing: true }],
        ]}
        components={components}
      >
        {children}
      </ReactMarkdown>
    </div>
  );
}
