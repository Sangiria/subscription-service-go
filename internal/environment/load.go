package environment

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load(".env")
	
	if err != nil {
		slog.Info("no .env file found, relying on system environment variables")
    }

	slog.Info("successfully loaded environment variables")
}