# promptopt website

The marketing site and documentation for promptopt. Next.js (App Router,
Turbopack), TypeScript, Tailwind CSS, Framer Motion.

```sh
npm install
npm run dev      # http://localhost:3000
npm run lint
npm run build    # static export-friendly production build
```

## How the docs are built

The documentation pages are not written here. They are the Markdown files in
the repository's top-level [`docs/`](../docs) directory — the same files you
read browsing the repo on GitHub. `lib/docs.ts` loads them at build time,
rewrites the relative `.md` links to site routes, extracts a table of
contents, and feeds `app/docs/[[...slug]]/page.tsx`, which statically renders
one HTML page per document.

To add or edit a doc, change the file under `docs/` and, if it is a new page,
add an entry to the `MANIFEST` array in `lib/docs.ts`.

## Data

Everything the interactive components show — terminal transcripts, the
before/after example, the playground's responses — is static data in
`lib/demo.ts`. The site makes no API calls and needs no API key.
