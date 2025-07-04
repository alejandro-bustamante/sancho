package controller

import (
	"database/sql"
	"net/http"

	model "github.com/alejandro-bustamante/sancho/server/internal/model"
	db "github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	queries *db.Queries
}

func NewUserHandler(q *db.Queries) *UserHandler {
	return &UserHandler{
		queries: q,
	}

}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	ctx := c.Request.Context()

	userParams := db.InsertUserParams{
		Username:     req.Username,
		PasswordHash: hashPassword(req.Password),
		Email:        sql.NullString{String: req.Email, Valid: req.Email != ""},
	}
	userDB, err := h.queries.InsertUser(ctx, userParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not register new user in the database", "details": err.Error()})
		return
	}
	user := model.UserFromDB(userDB)

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

	ctx := c.Request.Context()

	userDB, err := h.queries.GetUserByUsername(ctx, req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error searching user in the database", "details": err.Error()})
		return
	}
	authenticated, err := userAuthenticated(userDB.PasswordHash, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error authenticating user", "details": err.Error()})
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

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// en producción deberías loguear esto, no panickear
		panic("failed to hash password: " + err.Error())
	}
	return string(hash)
}

func userAuthenticated(passwordHash, passwordPlain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordPlain))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (h *UserHandler) DeleteUser(c *gin.Context) {

}

func (h *UserHandler) UpdateUser(c *gin.Context) {

}
