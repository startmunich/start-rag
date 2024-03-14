# STAGE 1: A container with pnpm and python3 is required
FROM node:21-alpine as pnpm_base

WORKDIR /app
RUN npm i --global --no-update-notifier --no-fund pnpm@8
RUN apk add --no-cache g++ make py3-pip libc6-compat

# STAGE 2: fetch deps into the pnpm store
# We run pnpm fetch in a separate step to avoid re-fetching deps on every code change
# fetch is a pnpm command that downloads all dependencies to the local store
# You could remove or skip this step if using npm or yarn (but make sure to copy your lock file)
FROM pnpm_base as fetched_deps
WORKDIR /app
COPY pnpm-lock.yaml ./
RUN pnpm config set store-dir /workdir/.pnpm-store
RUN pnpm fetch

# STAGE 3: Copy the application code and install all deps from cache into the application
FROM fetched_deps as with_all_deps
COPY . ./
RUN pnpm i

# STAGE 4: Build the NextJS app
# Here we use pnpm filtering to only build the frontend app
# Then we use pnpm deploy command to prune the dependencies
FROM with_all_deps as builder
RUN pnpm build

# STAGE 5: Create a clean production image - only take pruned assets
FROM node:18-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 app
RUN adduser --system --uid 1001 app
USER app
COPY --chown=app:app --from=builder /app/.next/standalone src/
COPY --chown=app:app --from=builder /app/public src/public
COPY --chown=app:app --from=builder /app/.next/static src/.next/static
ENV PORT 3000
EXPOSE 3000
CMD ["node", "src/apps/frontend/server.js"]