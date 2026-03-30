package main

import (
	"log"
	"os"
	"time"

	"ticketing-system/internal/config"
	"ticketing-system/internal/handler"
	"ticketing-system/internal/middleware"
	"ticketing-system/internal/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	if err := config.DB.AutoMigrate(
		&model.User{},
		&model.Ticket{},
		&model.Comment{},
		&model.Attachment{},
	); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "API is running"}) })
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/profile", handler.Profile)

		// Tickets
		auth.POST("/tickets", handler.CreateTicket)
		auth.GET("/tickets", handler.GetTickets)
		auth.PATCH("/tickets/:id/status", handler.UpdateTicketStatus)

		// Comments
		auth.GET("/tickets/:id/comments", handler.GetComments)
		auth.POST("/tickets/:id/comments", handler.AddComment)
		auth.DELETE("/tickets/:id/comments/:comment_id", handler.DeleteComment)

		// Attachments
		auth.POST("/tickets/:id/attachments", handler.UploadAttachment)
		auth.GET("/tickets/:id/attachments", handler.GetAttachments)
		auth.GET("/attachments/:id/download", handler.DownloadAttachment)
		auth.DELETE("/attachments/:id", handler.DeleteAttachment)

		// Admin routes
		admin := auth.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/tickets", handler.GetAllTickets)
			admin.GET("/tickets/:id", handler.GetTicketByID)
			admin.PUT("/tickets/:id", handler.EditTicket)
			admin.DELETE("/tickets/:id", handler.DeleteTicket)
			admin.PATCH("/tickets/:id/assign", handler.AssignTicket)
			admin.PATCH("/tickets/:id/status", handler.UpdateTicketStatus)
			admin.GET("/dashboard", handler.GetDashboard)
			admin.GET("/users", handler.GetAllUsers)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
