// A fixed film-grain overlay. Static SVG noise, blended over the whole page at
// a very low opacity. Adds texture and depth without motion or JS. Sits below
// the nav (z-40) and any modal (z-50) and never intercepts pointer events.
const NOISE =
  "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='160' height='160'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.82' numOctaves='2' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E";

export function Grain() {
  return (
    <div
      aria-hidden
      className="pointer-events-none fixed inset-0 z-[1] opacity-[0.022] mix-blend-overlay"
      style={{ backgroundImage: `url("${NOISE}")` }}
    />
  );
}
