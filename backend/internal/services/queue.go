package services

import (
	"database/sql"
	"fmt"
	"waitless-backend/internal/models"

	"github.com/google/uuid"
)

type QueueService struct {
	db *sql.DB
}

func NewQueueService(db *sql.DB) *QueueService {
	return &QueueService{db:db}
}

func (s *QueueService) CreateQueue(req *models.CreateQueueRequest, adminID uuid.UUID) (*models.Queue,error) {
	queue := &models.Queue{
		Name: req.Name,
		Description: req.Description,
		AdminID: adminID,
		IsActive: true,
	}

	query := `INSERT INTO queues(name,description,admin_id)
	VALUES ($1,$2,$3)
	RETURNING id, created_at`

	err := s.db.QueryRow(query,queue.Name,queue.Description,queue.AdminID).Scan(&queue.ID,&queue.CreatedAt)

	if err != nil {
		return nil,fmt.Errorf("failed to create queue: %w",err)
	}

	return queue,nil
}