package routes

import (
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

func setupUsersRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		// Rotas públicas (sem autenticação)
		users.POST("/register", handlers.CreateUser)
		users.POST("/login", handlers.Login)

		// Rotas protegidas (requerem autenticação JWT)
		// Para proteger uma rota, adicione o middleware: middlewares.AuthMiddleware()
		// Exemplo:
		// authenticated := users.Group("")
		// authenticated.Use(middlewares.AuthMiddleware())
		// {
		//     authenticated.GET("/", handlers.ListUsers)
		//     authenticated.GET("/:id", handlers.GetUser)
		//     authenticated.PUT("/:id", handlers.UpdateUser)
		//     authenticated.DELETE("/:id", handlers.DeleteUser)
		// }
	}
}
