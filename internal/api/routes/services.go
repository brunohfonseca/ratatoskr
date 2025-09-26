package routes

import (
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

// setupServicesRoutes configura rotas relacionadas aos serviços
func setupServicesRoutes(api *gin.RouterGroup) {
	services := api.Group("/services")
	{
		// CRUD básico de serviços
		services.GET("/", handlers.ListServices)
		services.POST("/", handlers.CreateService)
		services.GET("/:id", handlers.GetService)
		services.PUT("/:id", handlers.UpdateService)
		services.DELETE("/:id", handlers.DeleteService)

		// Health check e status
		services.GET("/:id/status", handlers.GetServiceStatus)
		services.POST("/:id/health-check", handlers.TriggerHealthCheck)

		// Histórico de health checks
		services.GET("/:id/history", handlers.GetServiceHistory)
		services.GET("/:id/uptime", handlers.GetServiceUptime)
	}
}
