package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// setupRoutes define todas as rotas da aplicação
func setupRoutes(router *gin.Engine) {
	servicesRoutes := router.Group("/services")
	{
		servicesRoutes.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Success",
			})
		})
		servicesRoutes.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "route test",
			})
		})
	}

	// Futuras rotas podem ser adicionadas aqui
	// Exemplo:
	// apiRoutes := router.Group("/api/v1")
	// healthRoutes := router.Group("/health")
}
