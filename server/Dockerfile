# --- Build Stage ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy Go module manifests
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build target binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker ./cmd/worker/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migrate ./cmd/migrate/main.go

# --- Final Stage ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Run as non-root user for security compliance
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy binaries from builder stage
COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/bin/worker ./worker
COPY --from=builder /app/bin/migrate ./migrate

EXPOSE 8080

CMD ["./api"]
