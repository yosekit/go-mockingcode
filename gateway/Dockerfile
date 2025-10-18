# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy shared modules
COPY pkg/logger/ ./pkg/logger/
COPY pkg/proto/ ./pkg/proto/

# Copy gateway service
COPY gateway/ ./gateway/

# Set working dir and build
WORKDIR /build/gateway
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /build/gateway/gateway .

EXPOSE 8080

CMD ["./gateway"]

