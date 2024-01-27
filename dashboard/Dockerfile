FROM node:18-alpine AS base

FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /dashboard

COPY dashboard/package.json dashboard/package-lock.json* ./
RUN npm ci

FROM base AS builder
WORKDIR /dashboard
COPY --from=deps /dashboard/node_modules ./node_modules
COPY dashboard .
RUN npm run build

FROM base AS runner
WORKDIR /dashboard
ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
COPY --from=builder /dashboard/public ./public
RUN mkdir .next
RUN chown nextjs:nodejs .next
COPY --from=builder --chown=nextjs:nodejs /dashboard/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /dashboard/.next/static ./.next/static

USER nextjs
ENV PORT 3000
ENV HOSTNAME "0.0.0.0"
EXPOSE 3000

CMD ["node", "server.js"]
