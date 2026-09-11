// The promptopt mascot: a cat — the `cat` command, sitting at the terminal.
// Its tail ends in a blinking cursor. No robot, no gradient.

export function Mascot({ className = "" }: { className?: string }) {
  return (
    <svg
      viewBox="0 0 80 80"
      fill="none"
      role="img"
      aria-label="promptopt mascot: a cat"
      className={className}
    >
      {/* tail, ending in a blinking terminal cursor */}
      <path
        d="M58 64C72 60 76 44 70 28"
        className="stroke-borderStrong"
        strokeWidth="5"
        strokeLinecap="round"
      />
      <rect x="66" y="14" width="7" height="12" rx="2" className="fill-accent animate-blink" />

      {/* body */}
      <ellipse cx="38" cy="58" rx="22" ry="18" className="fill-panel stroke-borderStrong" strokeWidth="2" />

      {/* head */}
      <circle cx="40" cy="30" r="17" className="fill-panel stroke-borderStrong" strokeWidth="2" />

      {/* ears */}
      <path d="M24 20L14 2L34 14Z" className="fill-panel stroke-borderStrong" strokeWidth="2" strokeLinejoin="round" />
      <path d="M56 20L66 2L46 14Z" className="fill-panel stroke-borderStrong" strokeWidth="2" strokeLinejoin="round" />
      <path d="M25 17L20 6L31 13Z" className="fill-accent" />
      <path d="M55 17L60 6L49 13Z" className="fill-accent" />

      {/* face */}
      <circle cx="33" cy="30" r="2.8" className="fill-fg" />
      <circle cx="47" cy="30" r="2.8" className="fill-fg" />
      <path d="M38 36L42 36L40 39Z" className="fill-accent" />
      <path
        d="M35 40Q40 44 45 40"
        className="stroke-faint"
        strokeWidth="1.6"
        strokeLinecap="round"
      />
      <path d="M18 34L28 33M17 38L28 37M18 42L28 41" className="stroke-faint" strokeWidth="1" />
      <path d="M62 34L52 33M63 38L52 37M62 42L52 41" className="stroke-faint" strokeWidth="1" />
    </svg>
  );
}
