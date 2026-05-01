package database

import (
	"log"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("GOOSE_DBSTRING")), &gorm.Config{ 
		TranslateError: true,
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database")
	}

	slog.Info("successfully connected to database")

	return db
}