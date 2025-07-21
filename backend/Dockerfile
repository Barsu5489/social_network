# Build stage
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git ca-certificates gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-w -s" -o backend-server ./server.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata sqlite
WORKDIR /root/
COPY --from=builder /app/backend-server .
EXPOSE 3000
CMD ["./backend-server"]
