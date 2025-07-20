package handler

import (
	"net/http"
	"waitless-backend/internal/models"
	"waitless-backend/internal/services"
	"waitless-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QueueHandler struct {
	queueService *services.QueueService
}

func NewQueueHandler(queueService *services.QueueService) *QueueHandler {
	return &QueueHandler{
		queueService: queueService,
	}
}

func (h *QueueHandler) CreateQueue(c *gin.Context) {
	var req models.CreateQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c,http.StatusBadRequest,err.Error())
		return 
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c,http.StatusUnauthorized,"User not authorised")
		return
	}

	queue,err := h.queueService.CreateQueue(&req, adminID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c,http.StatusInternalServerError,err.Error())
		return
	}
	utils.SuccessResponse(c,queue)
}