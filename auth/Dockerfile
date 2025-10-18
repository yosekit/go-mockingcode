# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy shared modules to match replace directive
COPY pkg/models/ ./pkg/models/
COPY pkg/logger/ ./pkg/logger/

# Copy auth service
COPY auth/ ./auth/

# Set working dir and build
WORKDIR /build/auth
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /build/auth/auth .

EXPOSE 8081

CMD ["./auth"]

