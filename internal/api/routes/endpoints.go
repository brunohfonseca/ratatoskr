package routes

import (
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

// setupServicesRoutes configura rotas relacionadas aos serviços
func setupServicesRoutes(api *gin.RouterGroup) {
	endpoints := api.Group("/endpoints")
	{
		// CRUD básico de serviços
		endpoints.GET("/", handlers.ListServices)
		endpoints.POST("/", handlers.CreateService)
		endpoints.GET("/:id", handlers.GetService)
		endpoints.PUT("/:id", handlers.UpdateService)
		endpoints.DELETE("/:id", handlers.DeleteService)

		// Health check e status
		endpoints.GET("/:id/status", handlers.GetServiceStatus)
		endpoints.POST("/:id/health-check", handlers.TriggerHealthCheck)

		// Histórico de health checks
		endpoints.GET("/:id/history", handlers.GetServiceHistory)
		endpoints.GET("/:id/uptime", handlers.GetServiceUptime)
	}
}
