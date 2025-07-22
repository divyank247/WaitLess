package main

import (
	"log"
	"waitless-backend/internal/config"
	"waitless-backend/internal/database"
	"waitless-backend/internal/handler"
	"waitless-backend/internal/middleware"
	"waitless-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migration", err)
	}

	authService := services.NewAuthService(db)
	queueService := services.NewQueueService(db)

	authHandler := handler.NewAuthHandler(authService)
	queueHandler := handler.NewQueueHandler(queueService)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "WaitLess backend is running 🚀",
		})
	})

	api := router.Group("/api") 
	{
		auth := api.Group("/auth") 
		{
			auth.POST("/register",authHandler.Register)
			auth.POST("/login",authHandler.Login)
		}

		queues := api.Group("/queues")
		{
			queues.GET("",queueHandler.GetQueues)
			queues.POST("",middleware.AuthMiddleware(),middleware.AdminMiddleware(),queueHandler.CreateQueue)
			queues.POST("/:id/join", middleware.AuthMiddleware(),queueHandler.JoinQueue)
			queues.DELETE("/:id/leave",middleware.AuthMiddleware(),queueHandler.LeaveQueue)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	router.Run(":" + cfg.Port)
}
