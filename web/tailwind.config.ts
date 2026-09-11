import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./lib/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Surfaces, darkest to lightest. Warm near-black, not blue-black —
        // pairs with the amber accent instead of fighting it.
        bg: "#0b0a09",
        panel: "#131110",
        raised: "#1a1714",
        border: "#2a2521",
        borderStrong: "#3a332c",
        // Text.
        fg: "#f2ede6",
        muted: "#a89b8c",
        faint: "#6f6459",
        // One accent. Used sparingly.
        accent: "#f5a524",
        accentDim: "#c9820f",
        // Semantic.
        ok: "#4ade80",
        warn: "#eab308",
        err: "#f87171",
      },
      fontFamily: {
        sans: ["var(--font-sans)", "ui-sans-serif", "system-ui", "sans-serif"],
        mono: [
          "var(--font-mono)",
          "ui-monospace",
          "SFMono-Regular",
          "Menlo",
          "monospace",
        ],
      },
      fontSize: {
        "2xs": ["0.6875rem", { lineHeight: "1rem" }],
      },
      maxWidth: {
        content: "72rem",
        prose: "44rem",
      },
      keyframes: {
        "fade-up": {
          from: { opacity: "0", transform: "translateY(8px)" },
          to: { opacity: "1", transform: "translateY(0)" },
        },
        blink: {
          "0%, 49%": { opacity: "1" },
          "50%, 100%": { opacity: "0" },
        },
      },
      animation: {
        "fade-up": "fade-up 0.4s ease both",
        blink: "blink 1s step-end infinite",
      },
    },
  },
  plugins: [],
};

export default config;
