import { Reveal } from "./Reveal";

export function Section({
  id,
  eyebrow,
  title,
  intro,
  children,
  surface,
  className = "",
}: {
  id?: string;
  eyebrow?: string;
  title?: string;
  intro?: React.ReactNode;
  children: React.ReactNode;
  // Optional ground treatment so the page isn't one flat slab: "tint" is a
  // faint raised panel, "grid" adds the dotted schematic grid on top of it.
  surface?: "tint" | "grid";
  className?: string;
}) {
  return (
    <section
      id={id}
      className={`relative overflow-hidden border-t border-border py-16 sm:py-24 ${
        surface ? "bg-panel/30" : ""
      } ${className}`}
    >
      {surface === "grid" && (
        <div
          className="grid-bg pointer-events-none absolute inset-0 opacity-50"
          aria-hidden
        />
      )}
      <div className="container-content relative">
        {(eyebrow || title) && (
          <Reveal className="max-w-prose">
            {eyebrow && (
              <div className="mb-3 flex items-center gap-2 font-mono text-2xs uppercase tracking-widest text-accent">
                <span className="h-px w-6 bg-accent/50" />
                {eyebrow}
              </div>
            )}
            {title && (
              <h2 className="text-2xl font-semibold tracking-tight sm:text-3xl">
                {title}
              </h2>
            )}
            {intro && <p className="mt-3 text-[0.95rem] text-muted">{intro}</p>}
          </Reveal>
        )}
        <div className={eyebrow || title ? "mt-10" : ""}>{children}</div>
      </div>
    </section>
  );
}
