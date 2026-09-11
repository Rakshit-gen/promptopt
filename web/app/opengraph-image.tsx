import { ImageResponse } from "next/og";
import { site } from "@/lib/site";

// The link-preview card shown when the URL is shared (Twitter/X, Slack,
// iMessage, Discord, etc.). Next.js serves this at /opengraph-image and
// wires it into the metadata automatically; twitter-image falls back to it.

export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

const bg = "#0b0a09";
const panel = "#131110";
const border = "#3a332c";
const fg = "#f2ede6";
const muted = "#a89b8c";
const accent = "#f5a524";

export default async function Image() {
  return new ImageResponse(
    (
      <div
        style={{
          width: "100%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          backgroundColor: bg,
          fontFamily: "sans-serif",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: 28 }}>
          <svg width="96" height="96" viewBox="0 0 80 80" fill="none">
            <path
              d="M58 64C72 60 76 44 70 28"
              stroke={border}
              strokeWidth="5"
              strokeLinecap="round"
            />
            <rect x="66" y="14" width="7" height="12" rx="2" fill={accent} />
            <ellipse cx="38" cy="58" rx="22" ry="18" fill={panel} stroke={border} strokeWidth="2" />
            <circle cx="40" cy="30" r="17" fill={panel} stroke={border} strokeWidth="2" />
            <path d="M24 20L14 2L34 14Z" fill={panel} stroke={border} strokeWidth="2" />
            <path d="M56 20L66 2L46 14Z" fill={panel} stroke={border} strokeWidth="2" />
            <path d="M25 17L20 6L31 13Z" fill={accent} />
            <path d="M55 17L60 6L49 13Z" fill={accent} />
            <circle cx="33" cy="30" r="2.8" fill={fg} />
            <circle cx="47" cy="30" r="2.8" fill={fg} />
            <path d="M38 36L42 36L40 39Z" fill={accent} />
          </svg>
          <div style={{ display: "flex", fontSize: 96, fontWeight: 600, letterSpacing: -2 }}>
            <span style={{ color: fg }}>prompt</span>
            <span style={{ color: muted }}>opt</span>
          </div>
        </div>

        <div style={{ display: "flex", marginTop: 28, fontSize: 32, color: muted }}>
          {site.tagline}
        </div>

        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: 10,
            marginTop: 56,
            padding: "14px 28px",
            borderRadius: 10,
            border: `1px solid ${border}`,
            backgroundColor: panel,
            fontFamily: "monospace",
            fontSize: 28,
          }}
        >
          <span style={{ color: accent }}>$</span>
          <span style={{ color: fg }}>promptopt optimize prompt.txt</span>
        </div>
      </div>
    ),
    size,
  );
}
