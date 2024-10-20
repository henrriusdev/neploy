package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type EnvVar struct {
	Env       string `env:"ENV" envDefault:"local"`
	BaseURL   string `env:"BASE_URL" envDefault:"http://localhost"`
	Port      string `env:"PORT" envDefault:"3000"`
	DBName    string `env:"DB_NAME" envDefault:"neploy"`
	DBUser    string `env:"DB_USER" envDefault:"postgres"`
	DBPass    string `env:"DB_PASS" envDefault:"Reyshell"`
	DBHost    string `env:"DB_HOST" envDefault:"localhost"`
	DBPort    string `env:"DB_PORT" envDefault:"5432"`
	DBSSLMode string `env:"DB_SSL_MODE" envDefault:"disable"`
	JWTSecret string `env:"JWT_SECRET" envDefault:"secret"`
}

var Env EnvVar

func getEnvPath() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, "..", ".env")
}

func LoadEnv() {
	if err := godotenv.Load(getEnvPath()); err != nil {
		log.Println("No .env file found")
	}

	if err := env.Parse(&Env); err != nil {
		log.Fatalf("Failed to parse env: %v", err)
	}
}
