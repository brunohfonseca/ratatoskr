package api

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/brunohfonseca/ratatoskr/internal/api/middlewares"
	"github.com/brunohfonseca/ratatoskr/internal/api/routes"
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
)

// setupRouter configura o router com middlewares e configurações básicas
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
	router.Use(middlewares.ZerologMiddleware())

	// Configurar proxies confiáveis para remover warning
	err := router.SetTrustedProxies(cfg.Server.TrustedProxies)
	if err != nil {
		return nil
	}

	// Configurar as rotas
	routes.SetupRoutes(router)

	return router
}

func ServerStart(cfg *config.AppConfig) *http.Server {
	router := setupRouter(cfg)

	// Criar logger personalizado que descarta logs de erro TLS
	silentLogger := log.New(io.Discard, "", 0)

	msg := "🔓 Servidor iniciado na porta: %s"

	srv := &http.Server{
		Addr:     ":" + strconv.Itoa(cfg.Server.Port),
		Handler:  router,
		ErrorLog: silentLogger,
	}

	if cfg.Server.SSL.Enabled {
		srv.Addr = ":" + strconv.Itoa(cfg.Server.SSL.Port)
		srv.TLSConfig = nil
		msg = "🔒 Servidor iniciado com SSL na porta %s"
	}
	zlog.Info().Msgf(msg, srv.Addr)
	return srv
}
