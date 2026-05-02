package database

import (
	"log/slog"
	"os"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db_logger = slogGorm.New(
    slogGorm.WithHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
        	if a.Key == "file" {
            	return slog.Attr{} 
        	}
        	return a
    	},
	})), 
    slogGorm.WithTraceAll(),
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("GOOSE_DBSTRING")), &gorm.Config{ 
		TranslateError: true,
		Logger: db_logger,
	})

	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
    	os.Exit(1)
	}

	slog.Info("successfully connected to database")

	return db
}