package main

import (
	"log/slog"
	"os"
	"subscription-service-go/internal/database"
	"subscription-service-go/internal/environment"
	"subscription-service-go/internal/handlers"
	"subscription-service-go/internal/repository"
	"subscription-service-go/internal/routes"

	"github.com/labstack/echo/v5/middleware"
	"github.com/labstack/echo/v5"
)

func main(){
	environment.LoadEnvVariables()
	db := database.InitDB()
	repo := repository.NewRepository(db)
	
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	e := echo.New()
	
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
    LogStatus:   true,
    LogURI:      true,
    LogMethod:   true,
    LogLatency:  true,
    HandleError: true,
    LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
        	if v.Error != nil {
            	logger.Error("request failed",
                	slog.String("method", v.Method),
                	slog.String("uri", v.URI),
                	slog.Int("status", v.Status),
                	slog.String("err", v.Error.Error()),
            	)
        	} else {
            	logger.Info("request processed",
                	slog.String("method", v.Method),
                	slog.String("uri", v.URI),
                	slog.Int("status", v.Status),
                	slog.Duration("latency", v.Latency),
            	)
        	}
        	return nil
    	},
	}))

	routes.InitSubscriptionRoutes(e, handlers.NewSubscriptionHandler(repo))

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}