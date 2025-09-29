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
		"user_id": user.ID,
		"uuid":    v7,
		"email":   user.Email,
	})
}

func Login(c *gin.Context) {

}
