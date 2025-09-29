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
		// Health check e status
		endpoints.GET("/:id/status", handlers.GetServiceStatus)
		endpoints.POST("/:id/health-check", handlers.TriggerHealthCheck)

		authenticated := endpoints.Group("")
		authenticated.Use(middlewares.AuthMiddleware())
		{
			// CRUD básico de endpoints
			endpoints.POST("/", handlers.CreateService)
			endpoints.GET("/", handlers.ListServices)
			endpoints.GET("/:id", handlers.GetService)
			endpoints.PUT("/:id", handlers.UpdateService)
			endpoints.DELETE("/:id", handlers.DeleteService)
			// Histórico de health checks
			endpoints.GET("/:id/history", handlers.GetServiceHistory)
			endpoints.GET("/:id/uptime", handlers.GetServiceUptime)

		}
	}
}
