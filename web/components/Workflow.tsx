"use client";

import { motion, useReducedMotion } from "framer-motion";
import { workflow } from "@/lib/demo";

export function Workflow() {
  const reduce = useReducedMotion();

  return (
    <ol className="relative space-y-3">
      <div
        aria-hidden
        className="absolute left-[15px] top-2 bottom-2 w-px bg-gradient-to-b from-border via-border to-transparent sm:left-[19px]"
      />
      {workflow.map((s, i) => (
        <motion.li
          key={s.step}
          className="relative flex gap-4 rounded-lg border border-border bg-panel/60 p-4 sm:gap-5"
          initial={reduce ? false : { opacity: 0, x: -10 }}
          whileInView={{ opacity: 1, x: 0 }}
          viewport={{ once: true, margin: "-40px" }}
          transition={{ duration: 0.35, delay: i * 0.05 }}
        >
          <span className="relative z-10 flex h-8 w-8 shrink-0 items-center justify-center rounded-full border border-borderStrong bg-raised font-mono text-2xs text-accent sm:h-10 sm:w-10">
            {i + 1}
          </span>
          <div>
            <div className="font-mono text-sm text-fg">{s.title}</div>
            <p className="mt-1 text-sm text-muted">{s.body}</p>
          </div>
        </motion.li>
      ))}
    </ol>
  );
}
