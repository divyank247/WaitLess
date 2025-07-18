package middleware

import (
	"net/http"
	"strings"
	"waitless-backend/internal/config"
	"waitless-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c,http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader,"Bearer ")
		if tokenString == authHeader {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token format")
			c.Abort()
			return
		}
		cfg := config.Load()
		claims, err := utils.ValidateToken(tokenString,cfg.JWTSecret)
		if err != nil {
			utils.ErrorResponse(c,http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id",claims.UserID)
		c.Set("user_email",claims.Email)
		c.Set("user_role",claims.Role)
		c.Next()
	}
}