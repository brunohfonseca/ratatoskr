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

	"github.com/brunohfonseca/ratatoskr/internal/infrastructure/bootstrap"
	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	redis "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
	"github.com/rs/zerolog/log"
)

func main() {
	configFile := flag.String("config", "/app/api-config.yml", "Arquivo de configuração")
	flag.Parse()

	cfg, srv := bootstrap.InitializeAPI(*configFile)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// start server
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
