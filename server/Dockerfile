# --- Build Stage ---
FROM golang:latest AS builder

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
RUN go build -o /app/bin/worker ./cmd/worker/main.go
RUN go build -o /app/bin/migrate ./cmd/migrate/main.go

# --- Final Stage ---
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Run as non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy binaries from builder stage
COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/bin/worker ./worker
COPY --from=builder /app/bin/migrate ./migrate

EXPOSE 8080

CMD ["./api"]
