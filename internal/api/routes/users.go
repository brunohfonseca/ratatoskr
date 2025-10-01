package routes

import (
	"github.com/brunohfonseca/ratatoskr/internal/api/middlewares"
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

func setupUsersRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		autenticated := users.Group("/")
		autenticated.Use(middlewares.AuthMiddleware())
		autenticated.Use(middlewares.AuditMiddleware())
		{
			autenticated.POST("/register", handlers.CreateUser)
		}
	}
}
