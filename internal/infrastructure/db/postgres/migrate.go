package infra

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

func Migrate(uri string) error {
	log.Info().Msg("🚀 Iniciando migrations no banco de dados")
	// Usa o driver pgx (já importado no pacote via blank import em `postgres.go`)
	db, err := sql.Open("pgx", uri)
	if err != nil {
		log.Error().Msgf("❌ erro ao abrir conexão com Postgres: %v", err)
		return err
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Error().Msgf("⚠️ erro ao fechar conexão de migrations: %v", cerr)
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error().Msgf("❌ erro ao inicializar driver de migrations (postgres): %v", err)
		return err
	}

	// Usa caminho relativo do projeto. "file:///migrations" aponta para "/migrations", que
	// não existe no ambiente local quando executado via `go run`. O correto é relativo.
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Error().Msgf("❌ erro ao criar instância de migrations: %v", err)
		return err
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Info().Msg("ℹ️ Nenhuma migration pendente (ErrNoChange)")
		} else {
			log.Error().Msgf("❌ erro ao executar migrations no banco de dados: %v", err)
			return err
		}
	}
	log.Info().Msg("✅ Migrations aplicadas com sucesso!")
	return nil
}
