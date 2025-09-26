package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// setupRoutes define todas as rotas da aplicação
func setupRoutes(router *gin.Engine) {
	sslRoutes := router.Group("/ssl")
	{
		sslRoutes.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Success",
			})
		})
		sslRoutes.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "route test",
			})
		})
		//sslRoutes.GET("/certs", handlers.GetSSLHandler)
		//sslRoutes.POST("/certs", handlers.PostSSLHandler)
		//sslRoutes.DELETE("/certs/:id", handlers.DeleteSSLHandler)
		//sslRoutes.PATCH("/certs/:id", handlers.PatchSSLHandler)
		//sslRoutes.GET("/certs/refresh", handlers.RefreshSSLHandler)
	}

	// Futuras rotas podem ser adicionadas aqui
	// Exemplo:
	// apiRoutes := router.Group("/api/v1")
	// healthRoutes := router.Group("/health")
}
