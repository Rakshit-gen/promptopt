/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  poweredByHeader: false,
  // Docs content is read from ../docs at build time; nothing dynamic at runtime.
  images: {
    // No remote images; the site uses inline SVG for all graphics.
    unoptimized: true,
  },
  async redirects() {
    return [
      { source: "/docs/commands/index", destination: "/docs/commands", permanent: true },
    ];
  },
};

export default nextConfig;
