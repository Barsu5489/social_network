# Backend Dockerfile (backend/Dockerfile)
FROM golang:1.22-alpine

# Install git for dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with optimizations for smaller binary
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o backend-server ./cmd/main.go

# Clean up build dependencies to reduce image size
RUN go clean -cache -modcache -testcache

EXPOSE 8080

CMD ["./backend-server"]
