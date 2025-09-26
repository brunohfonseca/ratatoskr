package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	ssls := router.Group("/ssl")
	{
		ssls.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Success",
			})
		})
		//ssls.GET("/certs", handlers.GetSSLHandler)
		//ssls.POST("/certs", handlers.PostSSLHandler)
		//ssls.DELETE("/certs/:id", handlers.DeleteSSLHandler)
		//ssls.PATCH("/certs/:id", handlers.PatchSSLHandler)
		//ssls.GET("/certs/refresh", handlers.RefreshSSLHandler)
	}

	return router
}

func Server(port int) {
	log.Printf("Starting REST API on port %d", port)
	router := setupRouter()
	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		panic("Erro ao iniciar o servidor: " + err.Error())
	}
}
