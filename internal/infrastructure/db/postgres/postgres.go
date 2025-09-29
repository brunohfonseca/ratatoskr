package infra

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

var PostgresConn *sql.DB

func ConnectPostgres(uri string) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		log.Fatal().Msgf("❌ erro ao abrir conexão com Postgres: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal().Msgf("❌ não conseguiu conectar no Postgres: %v", err)
	}

	log.Info().Msg("✅ Connected to Postgres (pgx)")
	PostgresConn = db
}

func CheckPostgresHealth() (bool, string, error) {
	if PostgresConn == nil {
		return false, "disconnected", nil
	}

	// timeout curto para health check
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := PostgresConn.PingContext(ctx); err != nil {
		return false, "error", err
	}

	return true, "connected", nil
}

func DisconnectPostgres() {
	if PostgresConn != nil {
		if err := PostgresConn.Close(); err != nil {
			log.Error().Msgf("⚠️ erro ao fechar conexão com Postgres: %v", err)
		} else {
			log.Info().Msg("✅ Disconnected from Postgres")
		}
	}
}
