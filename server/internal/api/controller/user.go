package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
type AuthenticateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserHandler struct {
	userService UserService // Inyectamos la interfaz
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Lógica delegada al servicio
	user, err := h.userService.RegisterUser(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not register new user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "User was successfully created.",
		"username":  user.Username,
		"createdAt": user.CreatedAt,
	})
}

func (h *UserHandler) AuthenticateUser(c *gin.Context) {
	var req AuthenticateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	authenticated, err := h.userService.AuthenticateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		// Diferenciar error de sistema vs usuario no encontrado si es necesario
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed", "details": err.Error()})
		return
	}

	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully authenticated.",
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {

}

func (h *UserHandler) UpdateUser(c *gin.Context) {

}
