FROM node:20-alpine AS base

# ── deps: install node_modules only ──────────────────────────────────────────
FROM base AS deps
WORKDIR /app/apps/main

# Copy only package.json (no lock file in repo — use legacy-peer-deps)
COPY apps/main/package.json ./package.json

# Remove @ats/shared from deps — it is resolved via tsconfig paths (no npm install needed)
RUN node -e " \
  const fs = require('fs'); \
  const p = JSON.parse(fs.readFileSync('./package.json','utf8')); \
  delete p.dependencies['@ats/shared']; \
  fs.writeFileSync('./package.json', JSON.stringify(p, null, 2)); \
"

RUN npm install --legacy-peer-deps

# ── builder: TypeScript compile + Next.js build ───────────────────────────────
FROM base AS builder
WORKDIR /app

# Mirror monorepo structure so tsconfig paths ("../../packages/shared/src/index.ts") resolve
COPY packages/shared ./packages/shared
COPY apps/main ./apps/main

# Bring in installed node_modules
COPY --from=deps /app/apps/main/node_modules ./apps/main/node_modules

WORKDIR /app/apps/main
# Pass internal API URL at build time so next.config.js rewrites() bakes it in
ARG INTERNAL_API_URL=http://go-api:8080
ENV INTERNAL_API_URL=$INTERNAL_API_URL
RUN npm run build

# ── runner: minimal production image ─────────────────────────────────────────
FROM base AS runner
ENV NODE_ENV=production
WORKDIR /app

COPY --from=builder /app/apps/main/.next/standalone ./
COPY --from=builder /app/apps/main/.next/static ./.next/static
COPY --from=builder /app/apps/main/public ./public

EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME=0.0.0.0
CMD ["node", "server.js"]
