package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brunohfonseca/ratatoskr/internal/api"
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	redis "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
	"github.com/rs/zerolog/log"
)

func main() {
	configFile := flag.String("config", "/app/api-config.yml", "Arquivo de configuração")
	flag.Parse()

	config.SetupLogs()
	if _, err := config.LoadConfig(*configFile); err != nil {
		log.Fatal().Err(err).Msg("❌ Erro ao carregar config")
	}

	cfg := config.Get()
	if cfg == nil {
		log.Fatal().Msg("❌ Configuração não carregada")
		return
	}

	if err := postgres.Migrate(cfg.Database.PostgresURL); err != nil {
		log.Fatal().Err(err).Msg("❌ Erro ao executar migrations no banco")
	}

	if err := handlers.InitKeycloak(); err != nil {
		log.Warn().Err(err).Msg("⚠️ Failed to initialize Keycloak SSO")
	}

	srv := api.ServerStart(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		var err error
		if cfg.Server.SSL.Enabled {
			err = srv.ListenAndServeTLS(cfg.Server.SSL.Cert, cfg.Server.SSL.Key)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("❌ Erro ao iniciar servidor")
		}
	}()

	log.Info().Msg("🚀 API iniciada! Pressione Ctrl+C para finalizar.")

	<-ctx.Done()
	log.Info().Msg("🛑 Parando API...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("⚠️ Erro ao finalizar servidor")
	}

	redis.DisconnectRedis()
	postgres.DisconnectPostgres()
	log.Info().Msg("✅ API finalizada com sucesso!")
}
