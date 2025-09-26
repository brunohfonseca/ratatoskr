package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/brunohfonseca/ratatoskr/internal/api"
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/brunohfonseca/ratatoskr/internal/database"
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
	database.ConnectRedis(cfg.Redis.RedisURL)
	database.ConnectMongoDB(cfg.Database.MongoURL)

	// Registrar models e sincronizar automaticamente
	database.RegisterAllModels()
	database.AutoSync(cfg.Database.MongoURL)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	srv := api.ServerStart(cfg)
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

	database.DisconnectRedis()
	database.DisconnectMongoDB()

	log.Info().Msg("✅ Aplicação finalizada com sucesso!")
}
