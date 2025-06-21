# Frontend Dockerfile (frontend/Dockerfile)
FROM node:20-alpine

# Install serve for lightweight static file serving
RUN npm install -g serve

WORKDIR /app

# Copy package files first for better caching
COPY package*.json ./
COPY vite.config.js svelte.config.js ./

# Install only production dependencies
RUN npm ci --only=production

# Copy source code
COPY src/ ./src/
COPY static/ ./static/

# Build the application
RUN npm run build

# Clean up node_modules and source to reduce size
RUN rm -rf node_modules src static package*.json vite.config.js svelte.config.js

EXPOSE 3000

# Use serve to host the built files
CMD ["serve", "-s", "dist", "-l", "3000"]

# ---