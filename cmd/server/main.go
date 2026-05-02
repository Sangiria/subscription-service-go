package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"subscription-service-go/internal/config"
	"subscription-service-go/internal/database"
	"subscription-service-go/internal/environment"
	"subscription-service-go/internal/handlers"
	"subscription-service-go/internal/repository"
	"subscription-service-go/internal/routes"
	"syscall"
	"time"

	_ "subscription-service-go/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	go func(){
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			slog.Error("server startup failed", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	<-ctx.Done()
	
	slog.Info("shutting down server")

	ctx_shd, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := e.Shutdown(ctx_shd); err != nil {
		slog.Error("shutting down with error", "error", err)
	} else {
		slog.Info("shut down complete")
	}
}