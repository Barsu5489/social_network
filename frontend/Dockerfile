# Builder stage
FROM node:20 AS builder
WORKDIR /app
COPY package.json vite.config.js svelte.config.js ./
COPY src/ ./src/
COPY static/ ./static/
RUN npm install
RUN npm run build

# Production stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY static/ /usr/share/nginx/html/static
EXPOSE 5173
CMD ["nginx", "-g", "daemon off;"]
