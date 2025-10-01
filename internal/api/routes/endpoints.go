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
		// Rotas protegidas - requerem autenticação JWT + Audit
		authenticated := endpoints.Group("")
		authenticated.Use(middlewares.AuthMiddleware())
		authenticated.Use(middlewares.AuditMiddleware())
		{
			// CRUD básico de endpoints
			authenticated.POST("/", handlers.CreateEndpoint)
			authenticated.GET("/", handlers.ListEndpoints)
			authenticated.GET("/show/:id", handlers.GetEndpoint)
			authenticated.PUT("/update/:id", handlers.UpdateEndpoint)
			authenticated.DELETE("/delete/:id", handlers.DeleteEndpoint)
			// Histórico de health checks
			authenticated.GET("/:id/history", handlers.GetEndpointHistory)
			authenticated.GET("/:id/uptime", handlers.GetEndpointUptime)
			// Roda o Check do endpoint
			authenticated.POST("/check/", handlers.CheckEndpoint)
		}
	}
}
