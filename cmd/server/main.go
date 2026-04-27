package main

import (
	"subscription-service-go/internal/database"
	"subscription-service-go/internal/environment"
	"subscription-service-go/internal/handlers"
	"subscription-service-go/internal/repository"
	"subscription-service-go/internal/routes"

	"github.com/labstack/echo/v5"
)

func main(){
	environment.LoadEnvVariables()
	db := database.InitDB()
	repo := repository.NewSubscriptionRepository(db)
	
	e := echo.New()
	routes.InitSubscriptionRoutes(e, handlers.NewSubscriptionHandler(repo))

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}