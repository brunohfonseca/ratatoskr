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
	"github.com/brunohfonseca/ratatoskr/internal/utils/logger"
	"github.com/brunohfonseca/ratatoskr/internal/worker"
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

	logger.InfoLog("🚀 Worker iniciado! Pressione Ctrl+C para finalizar.")
	<-ctx.Done()

	logger.InfoLog("🛑 Finalizando worker...")
	redis.DisconnectWorkerRedis(cfg.Name, workerName)
	postgres.DisconnectPostgres()
	logger.InfoLog("✅ Worker finalizado com sucesso!")
}
