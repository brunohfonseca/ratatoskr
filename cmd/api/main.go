package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/brunohfonseca/ratatoskr/internal/infrastructure"
	"github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/mongodb"
	"github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/redis"
	"github.com/rs/zerolog/log"
)

func main() {
	configFile := flag.String("config", "/app/config.yml", "Arquivo de configuração")
	flag.Parse()

	_, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatal().Msgf("❌ Erro ao carregar config: %v", err)
	}

	cfg := config.Get()
	if cfg == nil {
		log.Fatal().Msg("❌ Configuração não carregada")
		return
	}

	config.SetupLogs()
	log.Info().Msgf("🚀 Iniciando o serviço com o arquivo de configuração: %s", *configFile)
	redis.ConnectRedis(cfg.Redis.RedisURL)
	mongodb.ConnectMongoDB(cfg.Database.MongoURL)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	srv := infrastructure.ServerStart(cfg)
	// inicia servidor em goroutine
	go func() {
		if cfg.Server.SSL.Enabled {
			if err := srv.ListenAndServeTLS(cfg.Server.SSL.Cert, cfg.Server.SSL.Key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Msgf("❌ Erro ao iniciar servidor SSL: %v", err)
			}

		} else {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Msgf("❌ Erro ao iniciar servidor: %v", err)
			}
		}
	}()

	log.Info().Msg("🚀 Servidor iniciado! Pressione Ctrl+C para finalizar.")

	// Aguardar sinal de parada
	<-c
	fmt.Println("") //Quebra de Linha no CTRL+C
	log.Info().Msg("🛑 Sinal de parada recebido. Finalizando aplicação...")

	redis.DisconnectRedis()
	mongodb.DisconnectMongoDB()

	log.Info().Msg("✅ Aplicação finalizada com sucesso!")
}
