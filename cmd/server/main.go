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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "subscription-service-go/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Subscription Service API
// @version 1.0
// @description This is a documentation for subscription service API
// @contact.email lisaosadchenko@gmail.com
// @host localhost:1323
// @BasePath /

func main(){
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	environment.LoadEnvVariables()
	db := database.InitDB()

	repo := repository.NewRepository(db)
	
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(config.GetRequestLoggerConfig()))
	routes.InitSubscriptionRoutes(e, handlers.NewSubscriptionHandler(repo))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}