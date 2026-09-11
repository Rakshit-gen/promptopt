// Flat ESLint config. `next lint` was removed in Next 16, so the project runs
// ESLint directly (see the "lint" script in package.json).
import next from "eslint-config-next/core-web-vitals";

const config = [
  { ignores: [".next/**", "node_modules/**", "out/**"] },
  ...next,
  {
    rules: {
      // react-hooks 7 flags any setState in an effect body. We use it
      // deliberately for one-shot syncs (reset a dialog when it opens, close
      // the mobile menu on route change) where the cascading-render cost is nil.
      "react-hooks/set-state-in-effect": "off",
    },
  },
];

export default config;
