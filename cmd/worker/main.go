package main

import (
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	config.SetupLogs()
	log.Info().Msg("starting worker")
}
