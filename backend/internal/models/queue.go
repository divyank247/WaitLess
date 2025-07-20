package models

import (
	"time"

	"github.com/google/uuid"
)

type Queue struct {
	ID uuid.UUID `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	AdminID uuid.UUID `json:"admin_id" db:"admin_id"`
	IsActive bool `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	AdminName string `json:"admin_name,omitempty" db:"admin_name"`
	TicketCount int `json:"ticket_count,omitempty" db:"ticket_count"`
}

type CreateQueueRequest struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type QueueWithTickets struct {
	Queue Queue `json:"queue"`
	Tickets []Ticket `json:"tickets"`
}