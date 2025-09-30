package handlers

import (
	"net/http"

	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/brunohfonseca/ratatoskr/internal/utils"
	"github.com/brunohfonseca/ratatoskr/internal/utils/responses"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		responses.Error(c, http.StatusBadRequest, err)
		return
	}
	v7, err := uuid.NewV7()
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db := postgres.PostgresConn
	sql := "INSERT INTO users (uuid, full_name, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING id"
	err = db.QueryRow(sql,
		v7,
		user.FullName,
		user.Email,
		hashedPassword,
	).Scan(&user.ID)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err)
		return
	}
	responses.Success(c, http.StatusCreated, gin.H{
		"message": "ok",
		"id":      user.ID,
		"uuid":    v7,
		"email":   user.Email,
	})
}

func Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&loginRequest); err != nil {
		responses.ErrorMsg(c, http.StatusUnauthorized, "Email and password are required")
		return
	}

	db := postgres.PostgresConn

	// Buscar usuário pelo email
	var user models.User
	var passwordHash string
	sql := "SELECT id, uuid, email, full_name, password_hash, enabled FROM users WHERE email = $1"
	err := db.QueryRow(sql, loginRequest.Email).Scan(
		&user.ID,
		&user.UUID,
		&user.Email,
		&user.FullName,
		&passwordHash,
		&user.Enabled,
	)
	if err != nil {
		responses.ErrorMsg(c, http.StatusUnauthorized, "Email or password is invalid")
		return
	}

	// Verificar se o usuário está ativo
	if !user.Enabled {
		responses.ErrorMsg(c, http.StatusUnauthorized, "User is disabled")
		return
	}

	// Validar senha
	isValid, err := utils.VerifyPassword(loginRequest.Password, passwordHash)
	if err != nil {
		responses.ErrorMsg(c, http.StatusUnauthorized, "Error in password validation")
		return
	}

	if !isValid {
		responses.ErrorMsg(c, http.StatusUnauthorized, "Email or password is invalid")
		return
	}

	// Gerar JWT token
	token, err := utils.GenerateJWT(user.ID, user.UUID, user.Email)
	if err != nil {
		ErrMsg := "Error in token generation: " + err.Error()
		responses.ErrorMsg(c, http.StatusInternalServerError, ErrMsg)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"uuid":      user.UUID,
			"email":     user.Email,
			"full_name": user.FullName,
		},
	})
}
