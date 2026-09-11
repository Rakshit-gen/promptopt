"use client";

import { motion, useReducedMotion } from "framer-motion";
import { architecture } from "@/lib/demo";

export function Architecture() {
  const reduce = useReducedMotion();

  return (
    <div className="mx-auto max-w-xl">
      <ol className="space-y-0">
        {architecture.map((row, i) => (
          <li key={row.layer} className="flex flex-col items-center">
            <motion.div
              className="group w-full rounded-lg border border-border bg-panel px-4 py-3 transition-colors hover:border-accent/50 hover:bg-raised"
              initial={reduce ? false : { opacity: 0, y: 10 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-40px" }}
              transition={{ duration: 0.3, delay: i * 0.06 }}
            >
              <div className="flex items-baseline justify-between gap-4">
                <span className="font-mono text-sm text-fg">
                  <span className="text-faint group-hover:text-accent">
                    {String(i + 1).padStart(2, "0")}
                  </span>{" "}
                  {row.layer}
                </span>
                <span className="text-right text-2xs text-faint">{row.note}</span>
              </div>
            </motion.div>
            {i < architecture.length - 1 && (
              <svg
                width="16"
                height="18"
                viewBox="0 0 16 20"
                fill="none"
                aria-hidden
                className="text-borderStrong"
              >
                <path
                  d="M8 1V15M8 15L3 10M8 15L13 10"
                  stroke="currentColor"
                  strokeWidth="1.4"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            )}
          </li>
        ))}
      </ol>
      <p className="mt-5 text-center text-2xs text-faint">
        Dependencies point one way. Every operation is tested with a fake and no
        network.
      </p>
    </div>
  );
}
