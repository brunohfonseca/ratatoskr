package api

import (
	"net/http"
	"strconv"

	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

func ServerStart(cfg *config.AppConfig) *http.Server {
	router := setupRouter(cfg)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Server.Port),
		Handler: router,
	}

	if cfg.Server.SSL.Enabled {
		srv.Addr = ":" + strconv.Itoa(cfg.Server.SSL.Port)
		srv.TLSConfig = nil
	}

	log.Info().Msgf("Starting REST API on port %s", srv.Addr)
	return srv
}
