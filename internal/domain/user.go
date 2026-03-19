package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"id"`
	CreatedAt time.Time `db:"created_at"`
}
