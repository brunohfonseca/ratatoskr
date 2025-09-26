package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// setupNotificationsRoutes configura rotas de alertas
func setupNotificationsRoutes(api *gin.RouterGroup) {
	alerts := api.Group("/alerts")
	{
		// Rotas de canais de alertas
		channels := alerts.Group("/channels")
		{
			channels.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"total":   0,
					"message": "canais",
				})
			})
		}

		// Rotas de grupos de alertas
		groups := alerts.Group("/groups")
		{
			groups.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"total":   0,
					"message": "grupos de canais",
				})
			})
		}
	}
}
