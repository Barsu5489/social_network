# Builder stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY pkg/ ./pkg/
RUN go build -o /backend-server ./pkg/handlers/auth_handler.go

# Production stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /backend-server /backend-server
EXPOSE 8080
CMD ["/backend-server"]
