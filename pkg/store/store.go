package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"neploy.dev/config"
	"neploy.dev/pkg/logger"

	_ "github.com/lib/pq"
)

func NewConnection(cfg config.EnvVar) (Queryable, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSLMode)
	connection, err := sqlx.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to open connection: %v", err)
	}

	if err := connection.Ping(); err != nil {
		logger.Error("Failed to ping connection: %v", err)
	}

	return connection, err
}
