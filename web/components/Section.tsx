export function Section({
  id,
  eyebrow,
  title,
  intro,
  children,
  className = "",
}: {
  id?: string;
  eyebrow?: string;
  title?: string;
  intro?: React.ReactNode;
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <section id={id} className={`border-t border-border py-16 sm:py-20 ${className}`}>
      <div className="container-content">
        {(eyebrow || title) && (
          <div className="max-w-prose">
            {eyebrow && (
              <div className="mb-2 font-mono text-2xs uppercase tracking-widest text-accent">
                {eyebrow}
              </div>
            )}
            {title && (
              <h2 className="text-2xl font-semibold tracking-tight sm:text-3xl">
                {title}
              </h2>
            )}
            {intro && <p className="mt-3 text-[0.95rem] text-muted">{intro}</p>}
          </div>
        )}
        <div className={eyebrow || title ? "mt-8" : ""}>{children}</div>
      </div>
    </section>
  );
}
