package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/hel1th/clickstream/internal/domain"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	logger *slog.Logger
}

func NewProducer(brokers, topic string, logger *slog.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:  kafka.TCP(brokers),
		Topic: topic,

		Balancer: &kafka.LeastBytes{},

		RequiredAcks: kafka.RequireOne,

		AllowAutoTopicCreation: true,

		WriteTimeout: 10 * time.Second,
	}

	return &Producer{
		writer: writer,
		logger: logger,
	}
}

func (p *Producer) Publish(ctx context.Context, event domain.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.EventID.String()),
		Value: data,
	})
	if err != nil {
		return fmt.Errorf("write to kafka event: %w", err)
	}

	p.logger.Debug("event published to kafka",
		"event_id", event.EventID,
		"type", event.Type,
	)

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
