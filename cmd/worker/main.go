package main

import (
	"flag"

	"github.com/brunohfonseca/ratatoskr/internal/config"
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
	log.Info().Msg("starting worker")
}
