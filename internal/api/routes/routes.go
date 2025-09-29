package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes define todas as rotas da aplicação
func SetupRoutes(router *gin.Engine) {
	// API v1 routes
	setupAPIv1Routes(router)
}

// setupAPIv1Routes configura todas as rotas da API v1
func setupAPIv1Routes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Services routes - monitoramento de serviços
		setupEndpointsRoutes(api)
		// Alerts routes - configuração de alertas
		setupNotificationsRoutes(api)
		// Health routes - health check
		setupHealthRoutes(api)
		// Users routes - configuração de usuários
		setupUsersRoutes(api)
	}
}
