package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/brunohfonseca/ratatoskr/internal/infrastructure/bootstrap"
	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	redis "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
	"github.com/brunohfonseca/ratatoskr/internal/worker"
	"github.com/rs/zerolog/log"
)

func main() {
	configFile := flag.String("config", "/app/worker-config.yml", "Arquivo de configuração")
	flag.Parse()

	cfg := bootstrap.InitializeWorker(*configFile)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	workerName := cfg.Name + "-worker-" + os.Getenv("HOSTNAME")

	go func() {
		worker.StartWorker(ctx, redis.RedisClient, cfg.Name, workerName)
	}()

	log.Info().Msg("🚀 Worker iniciado! Pressione Ctrl+C para finalizar.")
	<-ctx.Done()

	log.Info().Msg("🛑 Finalizando worker...")
	redis.DisconnectWorkerRedis(cfg.Name, workerName)
	postgres.DisconnectPostgres()
	log.Info().Msg("✅ Worker finalizado com sucesso!")
}
