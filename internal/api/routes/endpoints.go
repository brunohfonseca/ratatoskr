package routes

import (
	"github.com/brunohfonseca/ratatoskr/internal/api/middlewares"
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

// setupEndpointsRoutes configura rotas relacionadas aos serviços
func setupEndpointsRoutes(api *gin.RouterGroup) {

	endpoints := api.Group("/endpoints")
	{
		// Health check e status (rotas públicas para monitoramento)
		endpoints.GET("/:id/status", handlers.GetEndpointStatus)
		endpoints.POST("/:id/health-check", handlers.TriggerHealthCheck)

		// Rotas protegidas - requerem autenticação JWT
		authenticated := endpoints.Group("")
		authenticated.Use(middlewares.AuthMiddleware())
		{
			// CRUD básico de endpoints
			authenticated.POST("/", handlers.CreateEndpoint)
			authenticated.GET("/", handlers.ListEndpoints)
			authenticated.GET("/:id", handlers.GetEndpoint)
			authenticated.PUT("/:id", handlers.UpdateEndpoint)
			authenticated.DELETE("/:id", handlers.DeleteEndpoint)
			// Histórico de health checks
			authenticated.GET("/:id/history", handlers.GetEndpointHistory)
			authenticated.GET("/:id/uptime", handlers.GetEndpointUptime)
		}
	}
}
