package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/brunohfonseca/ratatoskr/internal/config"
	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	workerName := cfg.Name + "-worker-" + os.Getenv("HOSTNAME")

	go func() {
		// Inicia o worker de health check
		worker.StartHealthCheckWorker(redis.RedisClient, cfg.Name, workerName)
	}()

	log.Info().Msg("🚀 Servidor iniciado! Pressione Ctrl+C para finalizar.")

	<-c
	fmt.Println("") //Quebra de Linha no CTRL+C
	log.Info().Msg("🛑 Sinal de parada recebido. Finalizando aplicação...")

	redis.DisconnectRedis()
	postgres.DisconnectPostgres()

	log.Info().Msg("✅ Aplicação finalizada com sucesso!")
}
