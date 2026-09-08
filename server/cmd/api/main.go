package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"tn-now/server/internal/admin"
	"tn-now/server/internal/auth"
	"tn-now/server/internal/compliance"
	"tn-now/server/internal/config"
	"tn-now/server/internal/content"
	"tn-now/server/internal/cron"
	"tn-now/server/internal/db"
	"tn-now/server/internal/middleware"
	"tn-now/server/internal/moderation"
	"tn-now/server/internal/portal"
	"tn-now/server/internal/reputation"
	"tn-now/server/internal/search"
	"tn-now/server/internal/social"
	"tn-now/server/internal/swagger"
	"tn-now/server/internal/video"
)

func main() {
	// Initialize structured slog JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration variables
	cfg := config.Load()

	logger.Info("Starting TN NOW API server",
		slog.String("env", cfg.AppEnv),
		slog.String("port", cfg.Port),
	)

	// Connect to PostgreSQL Database with connection pool
	var err error
	var pool *pgxpool.Pool
	ctxConn, cancelConn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConn()
	poolConfig, parseErr := pgxpool.ParseConfig(cfg.DatabaseURL)
	if parseErr != nil {
		logger.Error("Failed to parse PostgreSQL database URL", slog.String("error", parseErr.Error()))
	} else {
		poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
		poolConfig.MaxConns = 25
		poolConfig.MinConns = 5
		poolConfig.MaxConnLifetime = 1 * time.Hour
		poolConfig.MaxConnIdleTime = 30 * time.Minute
		pool, err = pgxpool.NewWithConfig(ctxConn, poolConfig)
		if err != nil {
			logger.Error("Failed to connect to PostgreSQL database pool, starting with pending connection", slog.String("error", err.Error()))
		} else {
			logger.Info("Connected to PostgreSQL pool successfully.")
			defer pool.Close()
		}
	}

	// Initialize Redis connection for Rate Limiter
	var rdb *redis.Client
	if cfg.RedisURL != "" {
		opt, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			logger.Warn("Failed to parse Redis URL, rate limiting will be bypassed", slog.String("error", err.Error()))
		} else {
			rdb = redis.NewClient(opt)
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			if err := rdb.Ping(ctx).Err(); err != nil {
				logger.Warn("Redis server ping failed, rate limiter will gracefully fallback to bypass", slog.String("error", err.Error()))
				_ = rdb.Close()
				rdb = nil
			} else {
				logger.Info("Connected to Redis successfully.")
			}
		}
	}

	// Base HTTP Mux Router
	mux := http.NewServeMux()

	// Register OpenAPI static endpoint
	swagger.RegisterRoutes(mux)

	// Initialize database accessors & auth modules if connection succeeded
	if pool != nil {
		queries := db.New(pool)
		authService := auth.NewService(queries, pool, cfg)
		authHandler := auth.NewHandler(authService)
		authHandler.RegisterRoutes(mux)

		// Home handler registration
		homeHandler := content.NewHomeHandler(queries, pool)
		homeHandler.RegisterRoutes(mux)

		// Content handler registration
		contentHandler := content.NewContentHandler(queries, pool)
		contentHandler.RegisterRoutes(mux)

		// Video handler registration
		videoHandler := video.NewHandler()
		videoHandler.RegisterRoutes(mux)

		// Submission handler registration
		submissionHandler := content.NewSubmissionHandler(queries, pool)
		submissionHandler.RegisterRoutes(mux)

		// Moderation handler registration
		moderationHandler := moderation.NewHandler(queries, pool)
		moderationHandler.RegisterRoutes(mux)

		// Reputation handler registration
		reputationHandler := reputation.NewHandler(queries, pool)
		reputationHandler.RegisterRoutes(mux)

		// Social handler registration
		socialHandler := social.NewHandler(queries, pool)
		socialHandler.RegisterRoutes(mux)

		// Compliance handler registration
		complianceHandler := compliance.NewHandler(queries, pool)
		complianceHandler.RegisterRoutes(mux)

		// Search handler registration
		searchHandler := search.NewHandler(queries)
		searchHandler.RegisterRoutes(mux)

		// Cron Background Scheduler registration
		cronScheduler := cron.NewScheduler(pool)
		cronScheduler.Start()
		defer cronScheduler.Stop()

		// Admin console & analytics handler registration
		adminHandler := admin.NewHandler(queries, pool, rdb, cronScheduler, cfg)
		adminHandler.RegisterRoutes(mux)

		// Public news & editorial portal handler registration
		portalHandler := portal.NewPortalHandler(queries, pool)
		portalHandler.RegisterRoutes(mux)
	}

	// Core application routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(portal.RenderPortalPage()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"data":{"service":"TN NOW API","status":"online","docs":"/swagger/doc.json","portal":"/portal","admin":"/admin"},"message":"Welcome to TN NOW API Server","errors":[]}`))
	})

	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"data":{"status":"live"},"message":null,"errors":[]}`))
	})

	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		redisReady := "disconnected"
		if rdb != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
			defer cancel()
			if err := rdb.Ping(ctx).Err(); err == nil {
				redisReady = "ready"
			}
		}

		dbReady := "disconnected"
		if pool != nil {
			ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
			defer cancel()
			if err := pool.Ping(ctx); err == nil {
				dbReady = "ready"
			}
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"data":{"status":"ready","redis":"` + redisReady + `","db":"` + dbReady + `"},"message":null,"errors":[]}`))
	})

	// Unified Rate Limiting struct
	rateLimiter := middleware.NewRateLimiter(rdb)

	// Chain middlewares: Recovery -> RequestID -> Logger -> CORS -> SecurityHeaders -> RateLimit
	var handler http.Handler = mux
	handler = rateLimiter.Limit(60, time.Minute)(handler) // 60 requests per minute limit
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.CORS(cfg.CORSAllowedOrigins)(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.Recovery(handler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for system signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan

		logger.Info("Shutdown signal received, shutting down server gracefully...")
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Error("Graceful shutdown timed out, forcing exit.")
				os.Exit(1)
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server shutdown failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
		
		if rdb != nil {
			_ = rdb.Close()
		}
		
		serverStopCtx()
	}()

	// Start HTTP Server
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}

	<-serverCtx.Done()
	logger.Info("Server stopped cleanly.")
}
