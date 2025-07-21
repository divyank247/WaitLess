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
	return &QueueService{db: db}
}

func (s *QueueService) CreateQueue(req *models.CreateQueueRequest, adminID uuid.UUID) (*models.Queue, error) {
	queue := &models.Queue{
		Name:        req.Name,
		Description: req.Description,
		AdminID:     adminID,
		IsActive:    true,
	}

	query := `INSERT INTO queues(name,description,admin_id)
	VALUES ($1,$2,$3)
	RETURNING id, created_at`

	err := s.db.QueryRow(query, queue.Name, queue.Description, queue.AdminID).Scan(&queue.ID, &queue.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}

	return queue, nil
}

func (s *QueueService) GetQueues() ([]models.Queue, error) {
	query :=
		`SELECT q.id,q.name,q.description,q.admin_id,q.is_active,q.created_at,u.name as admin_name,
	COALESCE(COUNT(t.id),0) as ticket_count
	FROM queues q
	LEFT JOIN users u ON q.admin_id = u.id
	LEFT JOIN tickets t ON q.id = t.queue_id AND t.status = 'waiting'
	WHERE q.is_active = true
	GROUP BY q.id,q.name,q.description,q.admin_id,q.is_active,q.created_at,u.name
	ORDER BY q.created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get queues: %w", err)
	}
	defer rows.Close()

	var queues []models.Queue
	for rows.Next() {
		var queue models.Queue
		err := rows.Scan(&queue.ID, &queue.Name, &queue.Description, &queue.AdminID,&queue.IsActive, &queue.CreatedAt, &queue.AdminName, &queue.TicketCount)

		if err != nil {
			return nil,fmt.Errorf("failed to scan quotes: %w",err)
		}
		queues = append(queues, queue)
	}

	return queues,nil
}

func (s *QueueService) JoinQueue(queueID,userID uuid.UUID) (*models.Ticket,error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tickets WHERE queue_id = $1 AND user_id = $2 AND status ='waiting')",queueID,userID).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("failed to check ticket existence: %w",err)
	}
	if exists {
		return nil,fmt.Errorf("user already in queue")
	}

	var nextPosition int

	err = s.db.QueryRow("SELECT COALESCE(MAX(position),0) + 1 FROM tickets WHERE queue_id = $1", queueID).Scan(&nextPosition)

	if err != nil {
		return nil,fmt.Errorf("failed to get next position: %w",err)
	}

	ticket := &models.Ticket{
		QueueID: queueID,
		UserID: userID,
		Position: nextPosition,
		Status: models.StatusWaiting,
	}

	query := `
	INSERT INTO tickets (queue_id,user_id,position,status)
	VALUES($1,$2,$3,$4)
	RETURNING id, created_at`

	err = s.db.QueryRow(query,ticket.QueueID,ticket.UserID,ticket.Position,ticket.Status).Scan(&ticket.ID,&ticket.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w",err)
	}

	return ticket, nil
}
