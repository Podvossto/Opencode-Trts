FROM node:20-alpine AS base

# ── deps: install node_modules only ──────────────────────────────────────────
FROM base AS deps
WORKDIR /app/apps/portal

# Copy only package.json
COPY apps/portal/package.json ./package.json

# Remove @ats/shared from deps — resolved via tsconfig paths
RUN node -e " \
  const fs = require('fs'); \
  const p = JSON.parse(fs.readFileSync('./package.json','utf8')); \
  if (p.dependencies && p.dependencies['@ats/shared']) delete p.dependencies['@ats/shared']; \
  fs.writeFileSync('./package.json', JSON.stringify(p, null, 2)); \
"

RUN npm install --legacy-peer-deps

# ── builder: TypeScript compile + Next.js build ───────────────────────────────
FROM base AS builder
WORKDIR /app

# Mirror monorepo structure so tsconfig paths ("../../packages/shared/src/index.ts") resolve
COPY packages/shared ./packages/shared
COPY apps/portal ./apps/portal

# Bring in installed node_modules
COPY --from=deps /app/apps/portal/node_modules ./apps/portal/node_modules

WORKDIR /app/apps/portal
# Pass internal API URL at build time so next.config.js rewrites() bakes it in
ARG INTERNAL_API_URL=http://go-api:8080
ENV INTERNAL_API_URL=$INTERNAL_API_URL
RUN npm run build

# ── runner: minimal production image ─────────────────────────────────────────
FROM base AS runner
ENV NODE_ENV=production
WORKDIR /app

COPY --from=builder /app/apps/portal/.next/standalone ./
COPY --from=builder /app/apps/portal/.next/static ./.next/static
COPY --from=builder /app/apps/portal/public ./public

EXPOSE 3001
ENV PORT=3001
ENV HOSTNAME=0.0.0.0
CMD ["node", "server.js"]
