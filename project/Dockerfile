# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy shared modules to match replace directive
COPY pkg/models/ ./pkg/models/
COPY pkg/logger/ ./pkg/logger/
COPY pkg/proto/ ./pkg/proto/

# Copy project service
COPY project/ ./project/

# Set working dir and build
WORKDIR /build/project
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o project ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /build/project/project .

EXPOSE 8082 9082

CMD ["./project"]

