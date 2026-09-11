// The promptopt mascot: the terminal tile from the favicon, given a face.
// Two eyes and a ">_" grin — still the exact glyphs a shell prints, now a
// little guy. No robot, no gradient. Inherits currentColor for the accent.

export function Mascot({ className = "" }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 64 64"
      fill="none"
      role="img"
      aria-label="promptopt mascot"
      className={className}
    >
      {/* blinking prompt bar / antenna */}
      <rect x="29.5" y="3" width="5" height="8" rx="2.5" className="fill-accent" />
      {/* terminal tile */}
      <rect
        x="6"
        y="11"
        width="52"
        height="44"
        rx="11"
        className="fill-panel stroke-borderStrong"
        strokeWidth="2"
      />
      {/* title-bar dots */}
      <circle cx="15" cy="20" r="1.7" className="fill-borderStrong" />
      <circle cx="21" cy="20" r="1.7" className="fill-borderStrong" />
      <circle cx="27" cy="20" r="1.7" className="fill-borderStrong" />
      {/* eyes */}
      <circle cx="25" cy="33" r="3" className="fill-fg" />
      <circle cx="39" cy="33" r="3" className="fill-fg" />
      {/* grin: ">_" */}
      <path
        d="M21 41.5 L29 46.5 L21 51.5"
        className="stroke-accent"
        strokeWidth="3.6"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <rect x="33" y="49" width="10" height="3.4" rx="1.7" className="fill-accent" />
    </svg>
  );
}
