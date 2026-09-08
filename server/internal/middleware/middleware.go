package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/redis/go-redis/v9"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// GetRequestID returns the request ID from context if present
func GetRequestID(ctx context.Context) string {
	if val := ctx.Value(requestIDKey); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

// generateUUID formats a random UUID-like string using crypto/rand
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// RequestID injects X-Request-ID header and stores it in request context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = generateUUID()
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ResponseWriter wrapper to capture status code and bytes written
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

// Logger logs incoming requests using slog JSON logger
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(rw, r)
		
		latency := time.Since(start)
		reqID := GetRequestID(r.Context())

		slog.Info("Request completed",
			slog.String("request_id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.statusCode),
			slog.Duration("latency", latency),
			slog.String("ip", r.RemoteAddr),
		)
	})
}

// Recovery recovers from panics and returns 500 JSON response
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				reqID := GetRequestID(r.Context())
				slog.Error("Request panic recovered",
					slog.String("request_id", reqID),
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"success":false,"data":null,"message":"Internal Server Error","errors":[]}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORS handles Cross-Origin Resource Sharing based on configurations
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				allowed := false
				for _, o := range allowedOrigins {
					if o == "*" || o == origin {
						allowed = true
						break
					}
				}
				if allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders sets standard HTTP security headers on all responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// RateLimiter holds Redis client and rate limits requests
type RateLimiter struct {
	rdb *redis.Client
}

func NewRateLimiter(rdb *redis.Client) *RateLimiter {
	return &RateLimiter{rdb: rdb}
}

func (rl *RateLimiter) Limit(limit int, period time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rl.rdb == nil {
				// Graceful bypass when Redis is not configured
				next.ServeHTTP(w, r)
				return
			}

			// Perform Ping check with short timeout to ensure Redis is alive
			ctx, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
			defer cancel()
			if err := rl.rdb.Ping(ctx).Err(); err != nil {
				slog.Warn("Redis rate limiter ping failed, bypassing checks", slog.String("error", err.Error()))
				next.ServeHTTP(w, r)
				return
			}

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			// Rate limit key structure
			key := fmt.Sprintf("ratelimit:%s:%s", ip, r.URL.Path)

			// Simple sliding window counter using Redis multi-transaction
			pipe := rl.rdb.TxPipeline()
			now := time.Now().UnixNano()
			clearBefore := time.Now().Add(-period).UnixNano()

			pipe.ZRemRangeByScore(r.Context(), key, "0", fmt.Sprintf("%d", clearBefore))
			pipe.ZAdd(r.Context(), key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})
			pipe.ZCard(r.Context(), key)
			pipe.Expire(r.Context(), key, period)

			cmds, err := pipe.Exec(r.Context())
			if err != nil {
				slog.Error("Redis rate limiter execution failed", slog.String("error", err.Error()))
				next.ServeHTTP(w, r)
				return
			}

			zcardCmd, ok := cmds[2].(*redis.IntCmd)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			currentCount, err := zcardCmd.Result()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if int(currentCount) > limit {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"success":false,"data":null,"message":"Too many requests. Please slow down.","errors":[]}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
