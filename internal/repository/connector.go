package repository

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func GetDBConnection() {
	var err error
	Db, err = gorm.Open(postgres.Open(os.Getenv("GOOSE_DBSTRING")), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database")
	}
}