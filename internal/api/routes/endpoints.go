package routes

import (
	"github.com/brunohfonseca/ratatoskr/internal/api/middlewares"
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

// setupServicesRoutes configura rotas relacionadas aos serviços
func setupEndpointsRoutes(api *gin.RouterGroup) {

	endpoints := api.Group("/endpoints")
	{
		// Health check e status (rotas públicas para monitoramento)
		endpoints.GET("/:id/status", handlers.GetServiceStatus)
		endpoints.POST("/:id/health-check", handlers.TriggerHealthCheck)

		// Rotas protegidas - requerem autenticação JWT
		authenticated := endpoints.Group("")
		authenticated.Use(middlewares.AuthMiddleware())
		{
			// CRUD básico de endpoints
			authenticated.POST("/", handlers.CreateService)
			authenticated.GET("/", handlers.ListServices)
			authenticated.GET("/:id", handlers.GetService)
			authenticated.PUT("/:id", handlers.UpdateService)
			authenticated.DELETE("/:id", handlers.DeleteService)
			// Histórico de health checks
			authenticated.GET("/:id/history", handlers.GetServiceHistory)
			authenticated.GET("/:id/uptime", handlers.GetServiceUptime)
		}
	}
}
