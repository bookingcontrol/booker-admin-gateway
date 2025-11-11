FROM golang:1.23-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Build cache invalidation (increment to force rebuild)
ARG CACHE_BUST=1
RUN echo "Cache bust: $CACHE_BUST"

# Copy source code
COPY cmd/ ./cmd/

# Copy internal packages
COPY internal/ ./internal/

# Build
# Note: booker-contracts-go is imported as a dependency, no proto generation needed
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/admin-gateway ./cmd/admin-gateway

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/admin-gateway .
# Create web/dist directory (frontend files are mounted via volume in docker-compose)
RUN mkdir -p ./web/dist

EXPOSE 8080

CMD ["./admin-gateway"]
