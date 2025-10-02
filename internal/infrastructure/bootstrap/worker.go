package bootstrap

import (
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/brunohfonseca/ratatoskr/internal/utils/logger"
)

func InitializeWorker(configFile string) *config.AppConfig {
	// Logs
	config.SetupLogs()

	// Carrega config
	if _, err := config.LoadConfig(configFile); err != nil {
		logger.FatalLog("❌ Erro ao carregar config", err)
	}

	cfg := config.Get()
	if cfg == nil {
		logger.FatalLog("❌ Configuração não carregada", nil)
		return nil
	}

	return cfg
}
