package repository

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("GOOSE_DBSTRING")), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database")
	}

	return db
}