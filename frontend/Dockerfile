FROM node:18-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/.next ./.next
RUN mkdir -p ./public && (cp -r /app/public/* ./public/ 2>/dev/null || true)
COPY --from=deps /app/node_modules ./node_modules
COPY package*.json ./
EXPOSE 5173
CMD ["npm", "start", "--", "-p", "5173"]
