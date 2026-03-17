package main

import (
	"ticketing-system/internal/config"
	"ticketing-system/internal/handler"
	"ticketing-system/internal/middleware"
	"ticketing-system/internal/model"

	"github.com/gin-gonic/gin"
)

func main() {
	// ======================
	// CONNECT DATABASE
	// ======================
	config.ConnectDB()

	// ======================
	// AUTO MIGRATE
	// ======================
	config.DB.AutoMigrate(
		&model.User{},
		&model.Ticket{},
	)

	// ======================
	// INIT GIN
	// ======================
	r := gin.Default()

	// ======================
	// PUBLIC ROUTES
	// ======================
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API is running"})
	})

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	// ======================
	// AUTH ROUTES
	// ======================
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// ======================
		// USER ROUTES
		// ======================
		auth.GET("/profile", handler.Profile)
		auth.POST("/tickets", handler.CreateTicket)
		auth.GET("/tickets", handler.GetTickets)

		// 🔥 NEW: update status ticket
		auth.PATCH("/tickets/:id/status", handler.UpdateTicketStatus)

		// ======================
		// ADMIN ROUTES
		// ======================
		admin := auth.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/tickets", handler.GetAllTickets)
			admin.PATCH("/tickets/:id/assign", handler.AssignTicket)
		}
	}

	// ======================
	// RUN SERVER
	// ======================
	r.Run(":8080")
}
