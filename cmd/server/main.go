package main

import (
	"log/slog"
	"os"
	"subscription-service-go/internal/config"
	"subscription-service-go/internal/database"
	"subscription-service-go/internal/environment"
	"subscription-service-go/internal/handlers"
	"subscription-service-go/internal/repository"
	"subscription-service-go/internal/routes"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main(){
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	environment.LoadEnvVariables()
	db := database.InitDB()

	logger.Info("successfully connected to database")

	repo := repository.NewRepository(db)
	
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(config.GetRequestLoggerConfig(logger)))
	routes.InitSubscriptionRoutes(e, handlers.NewSubscriptionHandler(repo, logger))

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}