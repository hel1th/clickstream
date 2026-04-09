package handler

import (
	"context"
	"fmt"

	"github.com/hel1th/clickstream/internal/domain"
)

func (h *Handler) saveEvent(ctx context.Context, e domain.Event) error {
	_, err := h.db.Exec(ctx, `
        INSERT INTO events (event_id, user_id, type, payload, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `, e.EventID, e.UserID, string(e.Type), e.Payload, e.Timestamp)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	return nil
}
