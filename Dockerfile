# Build Stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# Runtime Stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server ./server

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy scripts
COPY --from=builder /app/scripts ./scripts

# Make scripts executable
RUN chmod +x ./scripts/wait-for-db.sh

# Expose Port
EXPOSE 3000

# Run migrations then start server
CMD ["sh", "./scripts/wait-for-db.sh"]


