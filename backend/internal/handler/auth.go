package handler

import (
	"net/http"
	"waitless-backend/internal/models"
	"waitless-backend/internal/services"
	"waitless-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h* AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c,http.StatusBadRequest,err.Error())
		return
	}

	user,err := h.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c,http.StatusBadRequest,err.Error())
		return
	}

	utils.SuccessResponse(c,user)
}