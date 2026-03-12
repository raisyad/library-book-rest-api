package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppEnv  string
	DBUrl   string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("warning: .env file not found, using system environment variables")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	return Config{
		AppPort: appPort,
		AppEnv:  appEnv,
		DBUrl:   os.Getenv("DATABASE_URL"),
	}
}
