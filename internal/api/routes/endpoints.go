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
			authenticated.POST("/create", handlers.CreateEndpoint)
			authenticated.GET("/list", handlers.ListEndpoints)
			authenticated.PUT("/update/:id", handlers.UpdateEndpoint)
			authenticated.GET("/get/:id", handlers.GetEndpoint)
			authenticated.GET("/get/:id/history", handlers.GetEndpointHistory)
			authenticated.GET("/get/:id/uptime", handlers.GetEndpointUptime)
			authenticated.DELETE("/delete/:id", handlers.DeleteEndpoint)
			authenticated.POST("/check/:uuid", handlers.CheckEndpoint)
		}
	}
}
