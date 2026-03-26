package main

import (
	"log"
	"os"

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
	if err := config.DB.AutoMigrate(
		&model.User{},
		&model.Ticket{},
	); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// ======================
	// INIT GIN (PRO MODE)
	// ======================
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// ======================
	// ENV PORT (FLEXIBLE)
	// ======================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ======================
	// PUBLIC ROUTES
	// ======================
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
		})
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

		// update status
		auth.PATCH("/tickets/:id/status", handler.UpdateTicketStatus)

		// ======================
		// ADMIN ROUTES
		// ======================
		admin := auth.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/tickets", handler.GetAllTickets)
			admin.PATCH("/tickets/:id/assign", handler.AssignTicket)

			// 🔥 BONUS: dashboard endpoint (next step)
			admin.GET("/dashboard", handler.GetDashboard)
		}
	}

	// ======================
	// RUN SERVER
	// ======================
	log.Println("Server running on port:", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
