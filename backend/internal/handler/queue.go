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

func (h *QueueHandler) GetQueues(c *gin.Context) {
	queues, err := h.queueService.GetQueues()
	if err != nil {
		utils.ErrorResponse(c,http.StatusInternalServerError,err.Error())
		return
	}

	utils.SuccessResponse(c,queues)
}

func (h *QueueHandler) JoinQueue(c* gin.Context) {
	queueIDStr := c.Param("id")
	queueID,err := uuid.Parse(queueIDStr)
	if err != nil {
		utils.ErrorResponse(c,http.StatusBadRequest,"Invalid queue ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c,http.StatusUnauthorized,"User not authenticated")
		return
	}

	ticket,err := h.queueService.JoinQueue(queueID, userID.(uuid.UUID))

	if err != nil {
		utils.ErrorResponse(c,http.StatusBadRequest,err.Error())
		return
	}

	utils.SuccessResponse(c,ticket)
}

func (h *QueueHandler) LeaveQueue(c *gin.Context) {
	queueIDStr := c.Param("id")
	queueID, err := uuid.Parse(queueIDStr)
	if err != nil {
		utils.ErrorResponse(c,http.StatusBadRequest,"Invalid Queue ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c,http.StatusUnauthorized,"User not Authenticated")
		return
	}
	
	err = h.queueService.LeaveQueue(queueID,userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c,http.StatusBadRequest,err.Error())
		return
	}
	utils.MessageResponse(c,"Left queue successfully")
}