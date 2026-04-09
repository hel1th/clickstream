package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hel1th/clickstream/internal/domain"
)

type createEventRequest struct {
	UserID  uuid.UUID        `json:"user_id"`
	Type    domain.EventType `json:"type"`
	Payload map[string]any   `json:"payload"`
}

func (req *createEventRequest) validate() string {
	if req.UserID == uuid.Nil {
		return "user_id is required"
	}
	if !req.Type.IsValid() {
		return "type must be one of: click, view, purchase"
	}
	return ""
}

type createEventResponse struct {
	EventID uuid.UUID `json:"event_id"`
}

func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid json")
	}

	if errMsg := req.validate(); errMsg != "" {
		h.writeError(w, http.StatusBadRequest, errMsg)
		return
	}

	event := domain.Event{
		EventID:   uuid.New(),
		UserID:    req.UserID,
		Type:      req.Type,
		Payload:   req.Payload,
		Timestamp: time.Now().UTC(),
	}

	if err := h.saveEvent(r.Context(), event); err != nil {
		h.logger.Error("failed to save event", "error", err)
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := h.producer.Publish(r.Context(), event); err != nil {
		h.logger.Error("failed to publish event to kafka",
			"error", err,
			"event_id", event.EventID,
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createEventResponse{EventID: event.EventID})
}
