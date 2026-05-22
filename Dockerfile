# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY vendor/ vendor/
COPY . .
RUN go build -mod=vendor -o server ./cmd/main.go

# Runtime stage
FROM alpine:3.20
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY assets/ assets/
ENV ENV_FILE_PATH=/app/.env
EXPOSE 8080
CMD ["./server", "serve"]
