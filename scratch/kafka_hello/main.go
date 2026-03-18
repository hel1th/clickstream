package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	produce()
	time.Sleep(time.Second)
	consume()
}

func produce() {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "events",

		Balancer: &kafka.Hash{},
	})

	defer w.Close()

	messages := []kafka.Message{
		{Key: []byte("user-123"), Value: []byte(`{"event_id":"g1","type":"click","user_id":"user-123"}`)},
		{Key: []byte("user-123"), Value: []byte(`{"event_id":"g2","type":"view","user_id":"user-123"}`)},
		{Key: []byte("user-123"), Value: []byte(`{"event_id":"g3","type":"purchase","user_id":"user-123"}`)},
		{Key: []byte("user-456"), Value: []byte(`{"event_id":"g4","type":"click","user_id":"user-456"}`)},
		{Key: []byte("user-456"), Value: []byte(`{"event_id":"g5","type":"view","user_id":"user-456"}`)},
	}

	err := w.WriteMessages(context.Background(), messages...)
	if err != nil {
		log.Fatal("write:", err)
	}

	fmt.Printf("Sent %d messages\n", len(messages))
}

func consume() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "events",
		GroupID: "go-hello-group",

		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer r.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}
			fmt.Println("read msg err:", err)
			break
		}

		fmt.Printf(
			"partition=%d offset=%d key=%-10s value=%s\n",
			m.Partition, m.Offset, string(m.Key), string(m.Value),
		)
	}

	fmt.Println("Done!")
}
