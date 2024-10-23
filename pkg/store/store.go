package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"neploy.dev/config"

	_ "github.com/lib/pq"
)

func NewConnection(cfg config.EnvVar) (Queryable, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSLMode)
	connection, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open connection")
	}

	if err := connection.Ping(); err != nil {
		log.Error().Err(err).Msg("Failed to ping connection")
	}

	return connection, err
}
