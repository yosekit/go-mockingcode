# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy gateway service
COPY gateway/ ./

# Download deps and build
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/gateway .

EXPOSE 8080

CMD ["./gateway"]

