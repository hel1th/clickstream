package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hel1th/clickstream/internal/config"
	"github.com/hel1th/clickstream/internal/handler"
	"github.com/hel1th/clickstream/internal/kafka"
	"github.com/hel1th/clickstream/internal/middleware"
	"github.com/hel1th/clickstream/internal/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	))

	slog.SetDefault(logger)

	cfg, err := config.New()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pgPool, err := postgres.New(ctx, cfg.PostgresDSN)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pgPool.Close()

	logger.Info("connected to database")

	mux := http.NewServeMux()

	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic, logger)
	defer producer.Close()

	logger.Info("kafka producer initialized", "brokers", cfg.KafkaBrokers)

	h := handler.New(pgPool, producer, logger)
	h.RegisterRoutes(mux)

	srv := http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: middleware.Logging(logger)(mux),

		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		logger.Info("starting gateway", "port", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}

	logger.Info("gateway stopped")
}
