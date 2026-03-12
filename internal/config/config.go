package config

import "os"

type config struct {
	AppPort string
	DBUrl   string
}

func load() config {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return config{
		AppPort: appPort,
		DBUrl:   os.Getenv("DATABASE_URL"),
	}
}
