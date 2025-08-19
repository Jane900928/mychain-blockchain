# MyChain Blockchain Docker Image
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates jq

# Create app user
RUN addgroup -g 1000 mychain && \
    adduser -D -s /bin/sh -u 1000 -G mychain mychain

# Set working directory
WORKDIR /home/mychain

# Copy binary from builder stage
COPY --from=builder /app/build/mychaind /usr/local/bin/

# Copy scripts
COPY --from=builder /app/scripts/ ./scripts/

# Make scripts executable
RUN chmod +x ./scripts/*.sh

# Create data directory
RUN mkdir -p .mychain && chown -R mychain:mychain .mychain

# Switch to app user
USER mychain

# Expose ports
EXPOSE 26656 26657 1317 9090 9091

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:26657/health || exit 1

# Default command
CMD ["mychaind", "start"]
