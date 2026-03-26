package handler

import (
	"net/http"
	"strings"

	"ticketing-system/internal/service"
	"ticketing-system/internal/utils"

	"github.com/gin-gonic/gin"
)

//
// ======================
// REQUEST STRUCT
// ======================
//

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // optional
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

//
// ======================
// RESPONSE DTO
// ======================
//

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

//
// ======================
// REGISTER
// ======================
//

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// normalize role
	req.Role = strings.ToLower(req.Role)

	// default role
	if req.Role == "" {
		req.Role = "user"
	}

	// 🔥 VALIDASI ROLE (IMPORTANT)
	if req.Role != "user" && req.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	err := service.Register(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered",
	})
}

//
// ======================
// LOGIN
// ======================
//

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	response := AuthResponse{
		Token: token,
		User: gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	}

	c.JSON(http.StatusOK, response)
}
