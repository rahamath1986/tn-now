# --- Build Stage ---
FROM golang:alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

# Allow Go to auto-download the required toolchain version if needed
ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

# Copy Go module manifests first (better Docker layer caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build target binaries
RUN go build -o /app/bin/api ./cmd/api/main.go
RUN go build -o /app/bin/migrate ./cmd/migrate/main.go

# --- Final Stage ---
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Run as non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Memory limit optimization for low-memory environments (e.g. Render Free Tier 512MB):
# GOMEMLIMIT directs the Go runtime GC to aggressively reclaim heap before reaching 220MB,
# strictly preventing Linux cgroup OOM kills with more than 50% safety headroom.
ENV GOMEMLIMIT=220MiB
ENV GOGC=30

# Copy binaries from builder stage
COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/bin/migrate ./migrate

# Copy migrations folder so migrate binary can find SQL files
COPY --from=builder /app/migrations ./migrations

# Copy startup script
COPY --from=builder /app/start.sh ./start.sh

RUN chmod +x ./start.sh ./api ./migrate

USER appuser

EXPOSE 8080

# Run migrations then start API
CMD ["./start.sh"]
