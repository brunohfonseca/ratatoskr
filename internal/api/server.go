package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bhfonseca/ratatoskr/internal/config"
	"github.com/gin-gonic/gin"
)

func setupRouter(cfg *config.AppConfig) *gin.Engine {
	// Configurar modo de produção para reduzir logs de debug
	if !cfg.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configurar proxies confiáveis para remover warning
	err := router.SetTrustedProxies(cfg.Server.TrustedProxies)
	if err != nil {
		return nil
	}
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

func ServerStart(cfg *config.AppConfig) {
	log.Printf("Starting REST API on port %d", cfg.Server.Port)
	router := setupRouter(cfg)

	tls := cfg.Server.SSL.Enabled
	switch tls {
	case true:
		err := router.RunTLS(":"+strconv.Itoa(cfg.Server.SSL.Port), cfg.Server.SSL.Cert, cfg.Server.SSL.Key)
		if err != nil {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	case false:
		err := router.Run(":" + strconv.Itoa(cfg.Server.Port))
		if err != nil {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	default:
		log.Fatalf("Erro ao iniciar o servidor: Opção inválida")
	}
}
