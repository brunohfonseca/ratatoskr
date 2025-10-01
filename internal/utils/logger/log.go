package logger

import (
	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/rs/zerolog/log"
)

func ErrLog(msg string, err error) {
	log.Error().Err(err).Msg(msg)
}

func InfoLog(msg string) {
	log.Info().Msg(msg)
}

func DebugLog(msg string) {
	cfg := config.Get()

	if cfg.Environment == "development" {
		log.Debug().Msg(msg)
	}
}

func WarnLog(msg string) {
	log.Warn().Msg(msg)
}

func FatalLog(msg string, err error) {
	log.Fatal().Err(err).Msg(msg)
}

func FatalStrLog(msg string, key string, value string) {
	log.Fatal().Str(key, value).Msg(msg)
}
