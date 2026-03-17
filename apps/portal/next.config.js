const createNextIntlPlugin = require("next-intl/plugin");

const withNextIntl = createNextIntlPlugin("./src/lib/i18n.ts");

/** @type {import("next").NextConfig} */
const nextConfig = {
  output: "standalone",
  reactStrictMode: true,
  async rewrites() {
    // INTERNAL_API_URL: server-only Docker hostname — never sent to browser
    // NEXT_PUBLIC_API_URL: intentionally empty string so Axios uses relative URLs
    const internalUrl = process.env.INTERNAL_API_URL || "http://localhost:8080";
    return [
      {
        source: "/api/v1/:path*",
        destination: internalUrl + "/api/v1/:path*",
      },
    ];
  },
};

module.exports = withNextIntl(nextConfig);