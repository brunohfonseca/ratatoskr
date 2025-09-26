package api

import (
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

// setupRoutes define todas as rotas da aplicação
func setupRoutes(router *gin.Engine) {
	// Health check routes (sem autenticação)
	setupHealthRoutes(router)

	// API v1 routes
	setupAPIv1Routes(router)
}

// setupHealthRoutes configura rotas de health check
func setupHealthRoutes(router *gin.Engine) {
	health := router.Group("/health")
	{
		health.GET("/", handlers.HealthCheck)
		health.GET("/ready", handlers.ReadinessCheck)
		health.GET("/live", handlers.LivenessCheck)
	}
}

// setupAPIv1Routes configura todas as rotas da API v1
func setupAPIv1Routes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Services routes - monitoramento de serviços
		setupServicesRoutes(api)

		// Alerts routes - configuração de alertas
		setupAlertsRoutes(api)
	}
}

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

// setupAlertsRoutes configura rotas de alertas
func setupAlertsRoutes(api *gin.RouterGroup) {
	alerts := api.Group("/alerts")
	{
		alerts.Group("/channels")
		{

		}
		alerts.Group("/groups")
		{

		}
	}
}
