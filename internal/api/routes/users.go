package routes

import (
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

func setupUsersRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("/register", handlers.CreateUser)
		users.POST("/login", handlers.Login)
		//users.GET("/", handlers.ListUsers)
		//users.GET("/:id", handlers.GetUser)
		//users.PUT("/:id", handlers.UpdateUser)
		//users.DELETE("/:id", handlers.DeleteUser)
	}
}
