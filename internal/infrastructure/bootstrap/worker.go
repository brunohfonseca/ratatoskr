package bootstrap

import (
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	"github.com/rs/zerolog/log"
)

func InitializeWorker(configFile string) *config.AppConfig {
	// Logs
	config.SetupLogs()

	// Carrega config
	if _, err := config.LoadConfig(configFile); err != nil {
		log.Fatal().Err(err).Msg("❌ Erro ao carregar config")
	}

	cfg := config.Get()
	if cfg == nil {
		log.Fatal().Msg("❌ Configuração não carregada")
		return nil
	}

	// Migrations
	if err := postgres.Migrate(cfg.Database.PostgresURL); err != nil {
		log.Fatal().Err(err).Msg("❌ Erro ao executar migrations no banco")
	}

	// Keycloak
	if err := handlers.InitKeycloak(); err != nil {
		log.Warn().Err(err).Msg("⚠️ Failed to initialize Keycloak SSO")
	}

	return cfg
}
