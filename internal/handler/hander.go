package handler

import (
	"log/slog"

	"github.com/hel1th/clickstream/internal/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db       *pgxpool.Pool
	producer *kafka.Producer
	logger   *slog.Logger
}

func New(db *pgxpool.Pool, producer *kafka.Producer, logger *slog.Logger) *Handler {
	return &Handler{
		db:       db,
		producer: producer,
		logger:   logger,
	}
}
