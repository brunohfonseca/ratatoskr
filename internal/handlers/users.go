package handlers

import (
	"net/http"

	"github.com/brunohfonseca/ratatoskr/internal/entities"
	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	"github.com/brunohfonseca/ratatoskr/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(c *gin.Context) {
	var user entities.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v7, err := uuid.NewV7()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "ok",
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email e senha são obrigatórios"})
		return
	}

	db := postgres.PostgresConn

	// Buscar usuário pelo email
	var user entities.User
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
		return
	}

	// Verificar se o usuário está ativo
	if !user.Enabled {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário desativado"})
		return
	}

	// Validar senha
	isValid, err := utils.VerifyPassword(loginRequest.Password, passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao validar senha"})
		return
	}

	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha inválidos"})
		return
	}

	// Gerar JWT token
	token, err := utils.GenerateJWT(user.UUID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
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
