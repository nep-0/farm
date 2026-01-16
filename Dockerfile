# Build stage
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

# Create non-root user
RUN addgroup -g 1000 farm && \
    adduser -D -u 1000 -G farm farm

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/server .

# Change ownership to non-root user
RUN chown -R farm:farm /app

# Switch to non-root user
USER farm

# Expose default port
EXPOSE 8080

# Run the application
CMD ["./server"]