# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy shared models to match replace directive
COPY pkg/models/ ./pkg/models/

# Copy data service
COPY data/ ./data/

# Set working dir and build
WORKDIR /build/data
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o data ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /build/data/data .

EXPOSE 8083

CMD ["./data"]

