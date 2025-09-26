package routes

import (
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	"github.com/gin-gonic/gin"
)

// setupHealthRoutes configura rotas de health check
func setupHealthRoutes(router *gin.RouterGroup) {
	health := router.Group("/health")
	{
		health.GET("/", handlers.HealthCheck)
		health.GET("/ready", handlers.ReadinessCheck)
		health.GET("/live", handlers.LivenessCheck)
	}
}
