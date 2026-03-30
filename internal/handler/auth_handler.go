package handler

import (
	"net/http"

	"ticketing-system/internal/service"
	"ticketing-system/internal/utils"

	"github.com/gin-gonic/gin"
)

// ======================
// REQUEST STRUCTS
// ======================

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ======================
// RESPONSE DTO
// ======================

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// ======================
// REGISTER
// ======================
// Role selalu "user" — admin hanya bisa dibuat manual lewat DB atau endpoint khusus admin.
// Ini mencegah privilege escalation.

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.Register(req.Name, req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// ======================
// LOGIN
// ======================

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.Login(req.Email, req.Password)
	if err != nil {
		// Pesan generik — jangan kasih tau "user not found" vs "wrong password"
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
