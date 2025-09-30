package main

import (
	"flag"
	"os"

	"github.com/brunohfonseca/ratatoskr/internal/config"
	redis "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
	"github.com/brunohfonseca/ratatoskr/internal/worker"
	"github.com/rs/zerolog/log"
)

func main() {
	configFile := flag.String("config", "/app/worker-config.yml", "Arquivo de configuração")
	flag.Parse()

	config.SetupLogs()
	_, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatal().Msgf("❌ Erro ao carregar config: %v", err)
	}

	cfg := config.Get()
	if cfg == nil {
		log.Fatal().Msg("❌ Configuração não carregada")
		return
	}

	workerName := cfg.Name + "-worker-" + os.Getenv("INSTANCE_ID")
	log.Info().Msgf("🚀 Worker %s starting...", workerName)
	// Inicia o worker de health check
	worker.StartHealthCheckWorker(redis.RedisClient, cfg.Name, workerName)
}
