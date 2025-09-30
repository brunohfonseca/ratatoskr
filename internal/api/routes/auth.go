package routes

import (
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

// setupAuthRoutes configura rotas de autenticação
func setupAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		// Login tradicional (email/senha)
		auth.POST("/login", handlers.Login)

		// SSO Keycloak
		auth.GET("/keycloak", handlers.KeycloakLogin)
		auth.GET("/keycloak/callback", handlers.KeycloakCallback)
	}
}
