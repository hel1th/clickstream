package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeClick    EventType = "click"
	EventTypeView     EventType = "view"
	EventTypePurchase EventType = "purchase"
)

func (t EventType) IsValid() bool {
	switch t {
	case EventTypeClick, EventTypeView, EventTypePurchase:
		return true
	default:
		return false
	}
}

type Event struct {
	EventID   uuid.UUID      `json:"event_id"`
	UserID    uuid.UUID      `json:"user_id"`
	Type      EventType      `json:"type"`
	Payload   map[string]any `json:"payload"`
	Timestamp time.Time      `json:"timestamp"`
}
