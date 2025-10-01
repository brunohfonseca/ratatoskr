package bootstrap

import (
	"net/http"

	"github.com/brunohfonseca/ratatoskr/internal/api"
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/brunohfonseca/ratatoskr/internal/handlers"
	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	"github.com/rs/zerolog/log"
)

func InitializeAPI(configFile string) (*config.AppConfig, *http.Server) {
	// Logs
	config.SetupLogs()

	// Carrega config
	if _, err := config.LoadConfig(configFile); err != nil {
		log.Fatal().Err(err).Msg("❌ Erro ao carregar config")
	}

	cfg := config.Get()
	if cfg == nil {
		log.Fatal().Msg("❌ Configuração não carregada")
		return nil, nil
	}

	// Migrations
	if err := postgres.Migrate(cfg.Database.PostgresURL); err != nil {
		log.Fatal().Err(err).Msg("❌ Erro ao executar migrations no banco")
	}

	// Keycloak
	if err := handlers.InitKeycloak(); err != nil {
		log.Warn().Err(err).Msg("⚠️ Failed to initialize Keycloak SSO")
	}

	// Inicializa API
	srv := api.ServerStart(cfg)

	return cfg, srv
}
