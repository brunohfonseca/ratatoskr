package api

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
)

func setupRouter(cfg *config.AppConfig) *gin.Engine {
	if !cfg.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Desabilitar logs padrão do Gin redirecionando para discard
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	// Criar router sem middleware padrão
	router := gin.New()

	// Adicionar middleware de recovery personalizado
	router.Use(gin.Recovery())
	// Adicionar nosso middleware de logging com zerolog
	router.Use(ZerologMiddleware())

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

	// Criar logger personalizado que descarta logs de erro TLS
	silentLogger := log.New(io.Discard, "", 0)

	srv := &http.Server{
		Addr:     ":" + strconv.Itoa(cfg.Server.Port),
		Handler:  router,
		ErrorLog: silentLogger,
	}

	if cfg.Server.SSL.Enabled {
		srv.Addr = ":" + strconv.Itoa(cfg.Server.SSL.Port)
		srv.TLSConfig = nil
	}

	zlog.Info().Msgf("Starting REST API on port %s", srv.Addr)
	return srv
}
