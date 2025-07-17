package models

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID uuid.UUID `json:"id" db:"id"`
	QueueID uuid.UUID `json:"queue_id" db:"queue_id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Position int `json:"position" db:"position"`
	Status string `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UserName string `json:"user_name,omitempty" db:"user_name"`
}

const (
	StatusWaiting = "waiting"
	StatusCalled = "called"
	StatusCompleted = "completed"
	StatusSkipped = "skipped"
)